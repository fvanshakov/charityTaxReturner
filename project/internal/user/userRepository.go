package user

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"project/internal"
	"project/internal/database"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	db     *pgxpool.Pool
	config *UserRepositoryConfig
}

func NewUserRepository(db *database.DatabaseManager) *UserRepository {
	config := NewUserRepositoryConfig()
	config.Load()
	return &UserRepository{
		db:     db.Pool,
		config: config,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, password, email string) (*User, error) {
	user, err := NewUser(password, email, r.config)
	if err != nil {
		return nil, err
	}
	query := `
		INSERT INTO users (email, email_hmac, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err = r.db.QueryRow(ctx, query,
		user.Email,
		user.EmailHMAC,
		user.PasswordHash,
	).Scan(&user.ID)

	return user, err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, email, email_hmac, password_hash, oauth_token
		FROM users WHERE id = $1
	`

	var user User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.EmailHMAC,
		&user.PasswordHash,
		&user.OAuthToken,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, email_hmac, password_hash, oauth_token, created_at
		FROM users WHERE email_hmac = $1
	`

	var user User
	err := r.db.QueryRow(ctx, query, internal.ComputeHMAC(email, r.config.HMACKey)).Scan(
		&user.ID,
		&user.Email,
		&user.EmailHMAC,
		&user.PasswordHash,
		&user.OAuthToken,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *User) error {
	query := `
		UPDATE users 
		SET email = $2, email_hmac = $3, password_hash = $4, oauth_token = $5
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Email,
		user.EmailHMAC,
		user.PasswordHash,
		user.OAuthToken,
		time.Now(),
	)

	return err
}

func (r *UserRepository) SaveOauthToken(ctx context.Context, user *User, token *oauth2.Token) error {
	err := user.SetOAuthToken(token, r.config.OAuthEncryptionKey)
	if err != nil {
		return err
	}
	err = r.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
