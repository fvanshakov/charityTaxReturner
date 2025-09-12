package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"project/internal"
	"project/internal/gmail"
	"project/internal/user"
	"time"
)

type UserRequest struct {
	password string
	email    string
	authCode string
}

func (r UserRequest) getPassword() string {
	return r.password
}

func (r UserRequest) getEmail() string {
	return r.email
}

func (r UserRequest) getAuthCode() string {
	return r.authCode
}

const maxMessagesNbr = 1000

func BuildCreateUserHandler(
	ctx context.Context,
	repo *user.UserRepository,
	clientFile []byte,
	date time.Time,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userRequest UserRequest
		err := json.NewDecoder(r.Body).Decode(&userRequest)
		if err != nil {
			internal.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		userConfig := user.NewUserRepositoryConfig()
		userConfig.Load()

		user, err := user.NewUser(userRequest.getPassword(), userRequest.getEmail(), userConfig)

		service, err := gmail.GetGmailService(ctx, userRequest.getAuthCode(), clientFile, repo, user)
		if err != nil {
			internal.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err != nil {
			internal.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		receivedMessages, err := gmail.GetMessages(ctx, service, maxMessagesNbr, date)
		if err != nil {
			internal.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		filteredMessages := internal.FilterMessages(ctx, receivedMessages, maxMessagesNbr)
		println(filteredMessages)

	}
}
