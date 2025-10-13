package main

import (
	"charityTax/internal/handlers"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
)

func main() {
	ctx := context.Background()

	const maxMessagesNbr = 1000
	date := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)

	clientFile, err := os.ReadFile("./client_secret.json")

	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/api/status", func(ginCtx *gin.Context) {
		handlers.NewCreateUserHandler(ctx, ginCtx, clientFile, date)
	})

	err = r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
