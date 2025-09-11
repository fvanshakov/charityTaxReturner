package user

import (
	"encoding/base64"
	"encoding/json"
	"project/internal"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type User struct {
	ID           int64
	Email        string
	EmailHMAC    string
	PasswordHash string
	OAuthToken   []byte
}

func NewUser(
	password string,
	email string,
	config *UserRepositoryConfig,
) (*User, error) {
	user := User{}
	err := user.SetPassword(password)
	if err != nil {
		return nil, err
	}
	err = user.SetEmail(email, config.EncryptionKey, config.HMACKey)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) SetEmail(email string, encryptionKey []byte, hmacKey []byte) error {
	encryptedEmail, err := internal.EncryptData([]byte(email), encryptionKey)
	if err != nil {
		return err
	}

	u.Email = base64.StdEncoding.EncodeToString(encryptedEmail)
	u.EmailHMAC = internal.ComputeHMAC(email, hmacKey)
	return nil
}

func (u *User) GetEmail(encryptionKey []byte) (string, error) {
	encryptedData, err := base64.StdEncoding.DecodeString(u.Email)
	if err != nil {
		return "", err
	}

	decryptedData, err := internal.DecryptData(encryptedData, encryptionKey)
	if err != nil {
		return "", err
	}

	return string(decryptedData), nil
}

func (u *User) CheckEmail(email string, hmacKey []byte) bool {
	return internal.VerifyHMAC(email, u.EmailHMAC, hmacKey)
}

func (u *User) SetOAuthToken(token *oauth2.Token, encryptionKey []byte) error {
	tokenData, err := json.Marshal(token)
	if err != nil {
		return err
	}

	encryptedToken, err := internal.EncryptData(tokenData, encryptionKey)
	if err != nil {
		return err
	}

	u.OAuthToken = encryptedToken
	return nil
}

func (u *User) GetOAuthToken(encryptionKey []byte) (*oauth2.Token, error) {
	decryptedData, err := internal.DecryptData(u.OAuthToken, encryptionKey)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal(decryptedData, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
