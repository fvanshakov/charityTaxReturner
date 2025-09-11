package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseManager struct {
	Pool *pgxpool.Pool
}

func NewDatabase(
	ctx context.Context,
) (*DatabaseManager, error) {
	config := NewDatabaseConfig()
	err := config.Load()
	if err != nil {
		return nil, err
	}
	pgxConfig, err := pgxpool.ParseConfig(config.ConnectionString())
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &DatabaseManager{
		Pool: pool,
	}, nil
}

func (db *DatabaseManager) Close() {
	db.Pool.Close()
}
