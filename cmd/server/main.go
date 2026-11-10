package main

import (
	"context"
	"log"

	"github.com/Sujon37/school-management-system/internal/config"
	"github.com/Sujon37/school-management-system/internal/database"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()
	
	db, err := database.NewDB(ctx, cfg.StoragePath)
	if err != nil {
		log.Fatal("Database initialization failed: ", err)
	}
	defer db.Close()

	log.Println("Database connected successfully!")
}
