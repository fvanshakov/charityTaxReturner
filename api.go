package main

import (
	"fmt"
	"google.golang.org/api/gmail/v1"
	"time"
)

const user = "me" // константа пользователя, подразумевается что при таком значении пользователь определяется по токену

func getMessages(client *gmail.Service, maxMessages int, dateAfter string) ([]*gmail.Message, error) {

	var messages []*gmail.Message
	var fullMessages []*gmail.Message
	var pageToken string

	for {
		call := client.Users.Messages.List("me").MaxResults(500).Q(dateAfter)
		response, err := call.Do()
		if err != nil {
			fmt.Printf("не получилось сделать листинг сообщений")
			return nil, err
		}
		messages = append(messages, response.Messages...)
		pageToken = response.NextPageToken
		if pageToken == "" || len(messages) >= maxMessages {
			break
		}
	}

	for _, message := range messages {
		call := client.Users.Messages.Get(user, message.Id)
		response, err := call.Do()
		if err != nil {
			fmt.Println("Что-то пошло не так с получением письма c id ", message.Id)
			continue
		}
		fullMessages = append(fullMessages, response)

		// чтобы не ддосить АПИ почты
		time.Sleep(100 * time.Millisecond)

	}
	return fullMessages, nil
}
