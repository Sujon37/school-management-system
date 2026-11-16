package auth

import (
	"context"
	"database/sql"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Role         string
	IsActive     bool
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (
			id,
			email,
			password_hash,
			role,
			is_active
		)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
		boolToInt(user.IsActive),
	)

	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	var isActive int

	query := `
		SELECT id, email, password_hash, role, is_active
		FROM users
		WHERE email = ?
	`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&isActive,
	)

	if err != nil {
        return nil, err
    }

	user.IsActive = isActive == 1

	return &user, nil
}

func boolToInt(value bool) int {
	if value {return 1}
	return 0
}