package internal

import (
	"context"
	"encoding/base64"
	"fmt"
	"google.golang.org/api/gmail/v1"
)

var charities = []string{
	"Фонд “ВЕРА”",
	"Благотворительный фонд \"Дом с маяком\"",
	"БФ «Банк еды «Русь»",
}

func FilterMessages(ctx context.Context, messages []*gmail.Message, maxMessagesNbr int) []*gmail.Message {
	resultMessages := make([]*gmail.Message, 0, maxMessagesNbr)
	for _, message := range messages {
		body := message.Payload.Body.Data
		decodedBody, err := base64.URLEncoding.DecodeString(body)
		if err != nil {
			fmt.Println("что-то пошло не так при декодировании сообщения с id", message.Id)
			continue
		}
		if multiPartString(decodedBody).containsAny(charities) {
			resultMessages = append(resultMessages, message)
		}
	}
	return resultMessages
}
