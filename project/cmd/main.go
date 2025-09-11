package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"project/internal"
	"project/internal/database"
	"project/internal/handlers"
	"project/internal/user"
	"time"
)

func main() {
	ctx := context.Background()

	const maxMessagesNbr = 1000
	date := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)

	clientFile, err := os.ReadFile("../client_secret.json")
	dbManager, err := database.NewDatabase(ctx)
	defer dbManager.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	repo := user.NewUserRepository(dbManager)

	var migrator *database.Migrator
	migrator, err = database.NewMigrator("./migrations")
	err = migrator.Migrate(uint(internal.GetInt("MIGRATION_VERSION", 20250911112905)))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer migrator.Close()

	http.HandleFunc("/api/status", handlers.BuildCreateUserHandler(ctx, repo, clientFile, date))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
