package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strconv"
	"strings"
)

type Migration struct {
	Version int64
	Name    string
	File    string
}

func MigrateUp(ctx context.Context, db *sql.DB, migrations fs.FS) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("create migration table: %w", err)
	}

	files, err := fs.Glob(migrations, "migrations/*.up.sql")
	if err != nil {
		return fmt.Errorf("find migration files: %w", err)
	}

	list := make([]Migration, 0, len(files))
	for _, file := range files {
		migration, err := parseMigration(file)
		if err != nil {
			return err
		}
		list = append(list, migration)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Version < list[j].Version
	})

	for i := 1; i < len(list); i++ {
		if list[i].Version == list[i-1].Version {
			return fmt.Errorf("duplicate migration version: %d", list[i].Version)
		}
	}

	for _, migration := range list {
		var count int
		err := db.QueryRowContext(
			ctx,
			"SELECT COUNT(*) FROM schema_migrations WHERE version = ?",
			migration.Version,
		).Scan(&count)
		if err != nil {
			return fmt.Errorf("check migration %d: %w", migration.Version, err)
		}

		if count > 0 {
			continue
		}

		if err := applyMigration(ctx, db, migrations, migration); err != nil {
			return err
		}
	}

	return nil
}

func applyMigration(ctx context.Context, db *sql.DB, migrations fs.FS, migration Migration) error {
	content, err := fs.ReadFile(migrations, migration.File)
	if err != nil {
		return fmt.Errorf("read migration file %s: %w", migration.File, err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, string(content)); err != nil {
		return fmt.Errorf("execute migration %s: %w", migration.File, err)
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO schema_migrations (version, name)
		 VALUES (?, ?)`,
		migration.Version,
		migration.Name,
	)
	if err != nil {
		return fmt.Errorf("record migration %s: %w", migration.File, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", migration.File, err)
	}

	return nil
}

func MigrateDown(ctx context.Context, db *sql.DB, migrations fs.FS) error {
	var version int64
	var name string

	err := db.QueryRowContext(
		ctx,
		`SELECT version, name
		 FROM schema_migrations
		 ORDER BY version DESC
		 LIMIT 1`,
	).Scan(&version, &name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("find latest migration: %w", err)
	}

	file := fmt.Sprintf("migrations/%06d_%s.down.sql", version, name)

	content, err := fs.ReadFile(migrations, file)
	if err != nil {
		return fmt.Errorf("read down migration %s: %w", file, err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin rollback transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, string(content)); err != nil {
		return fmt.Errorf("execute down migration %s: %w", file, err)
	}

	_, err = tx.ExecContext(
		ctx,
		"DELETE FROM schema_migrations WHERE version = ?",
		version,
	)
	if err != nil {
		return fmt.Errorf("remove migration record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit rollback: %w", err)
	}

	return nil
}

func parseMigration(file string) (Migration, error) {
	base := path.Base(file)
	filename := strings.TrimSuffix(base, ".up.sql")

	parts := strings.SplitN(filename, "_", 2)
	if len(parts) != 2 || parts[1] == "" {
		return Migration{}, fmt.Errorf("invalid migration filename: %s", file)
	}

	version, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return Migration{}, fmt.Errorf("invalid migration version %s: %w", file, err)
	}

	if version <= 0 {
		return Migration{}, fmt.Errorf("migration version must be positive: %s", file)
	}

	return Migration{
		Version: version,
		Name:    parts[1],
		File:    file,
	}, nil
}