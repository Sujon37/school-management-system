package main

import (
	"context"
	"log"

	"github.com/Sujon37/school-management-system/db"
	"github.com/Sujon37/school-management-system/internal/config"
	"github.com/Sujon37/school-management-system/internal/database"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	sqlDB, err := database.NewDB(ctx, cfg.StoragePath)
	if err != nil {
		log.Fatal("Database initialization failed: ", err)
	}
	defer sqlDB.Close()

	log.Println("Database connected successfully!")

	if err := database.MigrateUp(ctx, sqlDB, db.Migrations); err != nil {
		log.Fatal("Migration failed: ", err)
	}

	log.Println("Database migrations completed successfully!")
}
