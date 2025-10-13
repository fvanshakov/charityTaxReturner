package gmail

import (
	"charityTax/internal"
	"charityTax/internal/database"
	userLogic "charityTax/internal/user"
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"time"
)

const ultimateMaxMessages = 500
const user = "me" // константа пользователя, подразумевается что при таком значении пользователь определяется по токену

type GmailClient struct {
	rep          *userLogic.UserRepository
	token        *oauth2.Token
	gmailService *gmail.Service
}

func NewGmailClient(
	ctx context.Context,
	authCode string,
	user *userLogic.User,
	clientFile []byte,
) (*GmailClient, error) {
	config, err := google.ConfigFromJSON(clientFile, gmail.GmailReadonlyScope)
	token, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, err
	}

	dbManager, err := database.NewDatabase(ctx)
	defer dbManager.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	repo := userLogic.NewUserRepository(dbManager)
	if err != nil {
		log.Fatal("не удалось создать http сервис")
		return nil, err
	}
	gmailClient := &GmailClient{
		repo,
		token,
		nil,
	}
	client := config.Client(context.Background(), token)

	existingUser, err := repo.GetUserByEmail(ctx, user.EmailHMAC)
	if err != nil {
		log.Fatal(err)
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}
	createdUser, err := repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	if err = repo.SaveOauthToken(ctx, createdUser, token); err != nil {
		return nil, err
	}
	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	gmailClient.gmailService = service
	return gmailClient, nil
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
		call := client.gmailService.Users.Messages.List("me").MaxResults(int64(maxResults)).Q(dateAfterString)
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
