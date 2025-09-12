package user

import "project/internal"

type UserRepositoryConfig struct {
	EncryptionKey      []byte
	OAuthEncryptionKey []byte
	HMACKey            []byte
}

func NewUserRepositoryConfig() *UserRepositoryConfig {
	return &UserRepositoryConfig{
		EncryptionKey:      make([]byte, 32),
		OAuthEncryptionKey: make([]byte, 32),
		HMACKey:            make([]byte, 32),
	}
}

func (config *UserRepositoryConfig) Load() {
	config.EncryptionKey = []byte(internal.GetString("USER_MAIL_ENCRYPTION_KEY", ""))
	config.HMACKey = []byte(internal.GetString("USER_MAIL_HMAC_KEY", ""))
	config.OAuthEncryptionKey = []byte(internal.GetString("USER_OAUTH_TOKEN_ENCRYPTION_KEY", ""))
}
