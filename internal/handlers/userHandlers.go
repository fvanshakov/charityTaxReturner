package handlers

import (
	"charityTax/internal"
	"charityTax/internal/database"
	"charityTax/internal/gmail"
	userPackage "charityTax/internal/user"
	"context"
	"github.com/gin-gonic/gin"
	"log"
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

func UserCreateHandler(
	ctx context.Context,
	ginCtx *gin.Context,
	clientFile []byte,
	date time.Time,
) {

	// 1. достали данные по пользователю из запроса
	var userRequest = NewUserRequest()
	err := ginCtx.ShouldBindJSON(&userRequest)
	if err != nil {
		internal.RespondWithError(ginCtx, http.StatusBadRequest, err.Error())
		return
	}
	userConfig := userPackage.NewUserRepositoryConfig()
	userConfig.Load()

	// 2. инициализировали репозиторий
	dbManager, err := database.NewDatabase(ctx)
	defer dbManager.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	repo := userPackage.NewUserRepository(dbManager)

	// 3. проверили наличие юзера и создали, если его нет
	draftUser, err := userPackage.NewUser(userRequest.getPassword(), userRequest.getEmail(), userConfig)
	existingUser, err := repo.GetUserByEmail(ctx, draftUser.EmailHMAC)
	if err != nil {
		log.Fatal(err)
	}
	if existingUser != nil {
		return
	}
	createdUser, err := repo.CreateUser(ctx, draftUser)
	if err != nil {
		return
	}

	// 4. получили токен-холдер
	clientParamsProvider, err := gmail.NewGmailClientParamsProvider(ctx, userRequest.getAuthCode(), nil, clientFile)
	if err != nil {
		return
	}
	if err = repo.SaveOauthToken(ctx, createdUser, clientParamsProvider.GetToken()); err != nil {
		log.Fatal(err)
		return
	}

	// 5. инициализируем клиент
	service, err := gmail.NewGmailClient(ctx, clientParamsProvider.GetClient())

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
