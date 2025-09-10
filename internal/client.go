package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

func retrieveTokenFromFile() (*oauth2.Token, error) {
	file, err := os.Open("token.json")
	if err != nil {
		fmt.Println("не удалось достать токен из файла")
		return nil, err
	}
	defer file.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)
	return token, err
}

func saveTokenToFile(token *oauth2.Token) error {
	file, err := os.Create("token.json")
	if err != nil {
		fmt.Println("не удалось записать токен из файла")
		return err
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(token)
	return err
}

func getTokenFromWeb(ctx context.Context, authCode string, config *oauth2.Config) (*oauth2.Token, error) {
	token, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func getClient(ctx context.Context, authCode string, config *oauth2.Config) (*http.Client, error) {

	token, err := retrieveTokenFromFile()
	if err != nil || token == nil {
		fmt.Printf("не удалось достать токен из файла")
		token, err = getTokenFromWeb(ctx, authCode, config)
		if err != nil {
			fmt.Printf("не удалось получить токен из сетевого запроса")
			return nil, err
		}
		err := saveTokenToFile(token)
		if err != nil {
			fmt.Printf("не удалось получить сохранить токен в файл")
		}
	}
	return config.Client(context.Background(), token), nil
}

func GetGmailService(ctx context.Context, authCode string, clientFile []byte) (*gmail.Service, error) {
	clientConfig, err := google.ConfigFromJSON(clientFile, gmail.GmailReadonlyScope)

	client, err := getClient(ctx, authCode, clientConfig)
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
