package gmail

import (
	"context"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"project/internal"
	"time"
)

const ultimateMaxMessages = 500
const user = "me" // константа пользователя, подразумевается что при таком значении пользователь определяется по токену

func GetMessages(ctx context.Context, client *gmail.Service, maxMessages int, dateAfter time.Time) ([]*gmail.Message, error) {

	messages := make([]*gmail.Message, 0, maxMessages)
	fullMessages := make([]*gmail.Message, 0, maxMessages)
	var pageToken string

	dateAfterString := "after:" + internal.FormatDateYYYYMMDD(dateAfter)

	maxResults := maxMessages
	if maxMessages > ultimateMaxMessages {
		maxResults = ultimateMaxMessages
	}

	for {
		call := client.Users.Messages.List("me").MaxResults(int64(maxResults)).Q(dateAfterString)
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
