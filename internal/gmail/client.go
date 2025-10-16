package gmail

import (
	"charityTax/internal"
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"net/http"
	"time"
)

const ultimateMaxMessages = 500
const user = "me" // константа пользователя, подразумевается что при таком значении пользователь определяется по токену

type GmailClient struct {
	gmailService *gmail.Service
}

type GmailClientParamsProvider struct {
	token  *oauth2.Token
	client *http.Client
}

func (paramsProvider *GmailClientParamsProvider) GetToken() *oauth2.Token {
	return paramsProvider.token
}

func (paramsProvider *GmailClientParamsProvider) GetClient() *http.Client {
	return paramsProvider.client
}

func NewGmailClientParamsProvider(
	ctx context.Context,
	authCode string,
	token *oauth2.Token,
	clientFile []byte,
) (*GmailClientParamsProvider, error) {
	resultToken := token
	config, err := google.ConfigFromJSON(clientFile, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}
	if token == nil && authCode != "" {
		resultToken, err = config.Exchange(ctx, authCode)
	}
	client := config.Client(context.Background(), token)
	if err != nil {
		return nil, err
	}
	return &GmailClientParamsProvider{token: resultToken, client: client}, nil
}

func NewGmailClient(
	ctx context.Context,
	client *http.Client,
) (*GmailClient, error) {
	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return &GmailClient{gmailService: service}, nil
}

func (client *GmailClient) GetMessages(ctx context.Context, maxMessages int, dateAfter time.Time) ([]*gmail.Message, error) {

	messages := make([]*gmail.Message, 0, maxMessages)
	fullMessages := make([]*gmail.Message, 0, maxMessages)
	var pageToken string

	dateAfterString := "after:" + internal.FormatDateYYYYMMDD(dateAfter)

	maxResults := maxMessages
	if maxMessages > ultimateMaxMessages {
		maxResults = ultimateMaxMessages
	}

	for {
		call := client.gmailService.Users.Messages.List(user).MaxResults(int64(maxResults)).Q(dateAfterString)
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
		call := client.gmailService.Users.Messages.Get(user, message.Id)
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
