package handlers

import (
	"charityTax/internal"
	"charityTax/internal/gmail"
	"charityTax/internal/user"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserRequest struct {
	password string
	email    string
	authCode string
}

type Option func(*UserRequest)

func WithAuthCode(authCode string) Option {
	return func(r *UserRequest) { r.authCode = authCode }
}

func WithEmail(email string) Option {
	return func(r *UserRequest) { r.email = email }
}

func WithPassword(password string) Option {
	return func(r *UserRequest) { r.password = password }
}

func NewUserRequest(
	opts ...Option,
) *UserRequest {
	ur := UserRequest{}

	for _, opt := range opts {
		opt(&ur)
	}
	return &ur
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

func NewCreateUserHandler(
	ctx context.Context,
	ginCtx *gin.Context,
	clientFile []byte,
	date time.Time,
) {
	var userRequest = NewUserRequest()
	err := ginCtx.ShouldBindJSON(&userRequest)
	if err != nil {
		internal.RespondWithError(ginCtx, http.StatusBadRequest, err.Error())
		return
	}

	userConfig := user.NewUserRepositoryConfig()
	userConfig.Load()

	user, err := user.NewUser(userRequest.getPassword(), userRequest.getEmail(), userConfig)

	service, err := gmail.NewGmailClient(ctx, userRequest.getAuthCode(), user, clientFile)
	if err != nil {
		internal.RespondWithError(ginCtx, http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		internal.RespondWithError(ginCtx, http.StatusBadRequest, err.Error())
		return
	}
	receivedMessages, err := service.GetMessages(ctx, maxMessagesNbr, date)
	if err != nil {
		internal.RespondWithError(ginCtx, http.StatusBadRequest, err.Error())
		return
	}
	filteredMessages := internal.FilterMessages(ctx, receivedMessages, maxMessagesNbr)
	println(filteredMessages)
}
