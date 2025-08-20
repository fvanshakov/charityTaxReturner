package main

import (
	"log"
	"os"
)

func main() {
	authCode := os.Getenv("AUTH_CODE")

	clientFile, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatal("не удалось распарсить файл с секретами клиента")
	}

	service, err := getGmailService(authCode, clientFile)
	if err != nil {
		log.Fatal("не удалось создать http сервис")
	}
	receivedMessages, err := getMessages(service, 1000, "after:2025/08/19")
	if err != nil {
		return
	}
	filteredMessages := filterMessages(receivedMessages)
	println(filteredMessages)
}
