package database

import (
	"fmt"
	"project/internal"
)

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxConnections  int
	MinConnections  int
	MaxConnLifetime int
	MaxConnIdleTime int
	EncryptionKey   string
	HMACKey         string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "postgres",
		Name:            "charity",
		SSLMode:         "disable",
		MaxConnections:  10,
		MinConnections:  2,
		MaxConnLifetime: 3600,
		MaxConnIdleTime: 1800,
		EncryptionKey:   "default-32-byte-encryption-key-here!",
		HMACKey:         "default-hmac-key-for-email-indexing",
	}
}

func (c *DatabaseConfig) Load() error {
	c.Host = internal.GetString("DB_HOST", c.Host)
	c.Port = internal.GetInt("DB_PORT", c.Port)
	c.User = internal.GetString("DB_USER", c.User)
	c.Password = internal.GetString("DB_PASSWORD", c.Password)
	c.Name = internal.GetString("DB_NAME", c.Name)
	c.SSLMode = internal.GetString("DB_SSL_MODE", c.SSLMode)
	c.MaxConnections = internal.GetInt("DB_MAX_CONNECTIONS", c.MaxConnections)
	c.MinConnections = internal.GetInt("DB_MIN_CONNECTIONS", c.MinConnections)
	c.MaxConnLifetime = internal.GetInt("DB_MAX_CONN_LIFETIME", c.MaxConnLifetime)
	c.MaxConnIdleTime = internal.GetInt("DB_MAX_CONN_IDLE_TIME", c.MaxConnIdleTime)
	c.EncryptionKey = internal.GetString("DB_ENCRYPTION_KEY", c.EncryptionKey)
	c.HMACKey = internal.GetString("DB_HMAC_KEY", c.HMACKey)

	return nil
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}
