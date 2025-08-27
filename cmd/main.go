package main

import (
	"charityTax/internal"
	"context"
	"log"
	"os"
	"time"
)

func main() {
	authCode := os.Getenv("AUTH_CODE")
	ctx := context.Background()

	const maxMessagesNbr = 1000
	date := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)

	clientFile, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatal("не удалось распарсить файл с секретами клиента")
	}

	service, err := internal.GetGmailService(ctx, authCode, clientFile)
	if err != nil {
		log.Fatal("не удалось создать http сервис")
	}
	receivedMessages, err := internal.GetMessages(ctx, service, maxMessagesNbr, date)
	if err != nil {
		return
	}
	filteredMessages := internal.FilterMessages(ctx, receivedMessages, maxMessagesNbr)
	println(filteredMessages)
}
