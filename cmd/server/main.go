package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Sujon37/school-management-system/cmd/seed"
	"github.com/Sujon37/school-management-system/db"
	"github.com/Sujon37/school-management-system/internal/auth"
	"github.com/Sujon37/school-management-system/internal/config"
	"github.com/Sujon37/school-management-system/internal/database"
	"github.com/Sujon37/school-management-system/internal/httpx"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.MustLoad()
	ctx := context.Background()

	sqlDB, err := database.NewDB(ctx, cfg.StoragePath)
	if err != nil {
		return fmt.Errorf("database initialization failed: %w", err)
	}
	defer sqlDB.Close()

	log.Println("Database connected successfully!")

	if err := database.MigrateUp(ctx, sqlDB, db.Migrations); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Database migrations completed successfully!")

	userRepo := auth.NewUserRepository(sqlDB)

	err = seed.SeedAdmin(ctx, userRepo, cfg.Admin.Email, cfg.Admin.Password)

	if err != nil {
		return fmt.Errorf("Failed to seed admin: %s", err)
	}

	router := httpx.NewRouter()

	server := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	log.Printf("HTTP server listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server failed: %w", err)
	}

	return nil
}
