package gmail

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"net/http"
	userLogic "project/internal/user"
)

func saveTokenToDB(
	ctx context.Context,
	token *oauth2.Token,
	user *userLogic.User,
	rep *userLogic.UserRepository,
) error {
	if err := rep.SaveOauthToken(ctx, user, token); err != nil {
		return err
	}
	return nil
}

func getTokenFromWeb(ctx context.Context, authCode string, config *oauth2.Config) (*oauth2.Token, error) {
	token, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func getClient(
	ctx context.Context,
	authCode string,
	config *oauth2.Config,
	rep *userLogic.UserRepository,
	user *userLogic.User,
) (*http.Client, error) {

	token, err := getTokenFromWeb(ctx, authCode, config)
	if err != nil {
		fmt.Printf("не удалось получить токен из сетевого запроса")
		return nil, err
	}
	err = saveTokenToDB(ctx, token, user, rep)
	if err != nil {
		fmt.Printf("не удалось  сохранить токен в БД")
		return nil, err
	}
	return config.Client(context.Background(), token), nil
}

func GetGmailService(
	ctx context.Context,
	authCode string,
	clientFile []byte,
	rep *userLogic.UserRepository,
	user *userLogic.User,
) (*gmail.Service, error) {
	clientConfig, err := google.ConfigFromJSON(clientFile, gmail.GmailReadonlyScope)

	client, err := getClient(ctx, authCode, clientConfig, rep, user)
	if err != nil {
		log.Fatal("не удалось создать клиент")
	}

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatal("не удалось создать http сервис")
		return nil, err
	}
	return service, nil
}
