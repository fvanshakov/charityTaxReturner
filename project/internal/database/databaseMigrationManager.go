package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

type Migrator struct {
	migrate *migrate.Migrate
}

func NewMigrator(
	migrationsPath string,
) (*Migrator, error) {
	config := NewDatabaseConfig()
	err := config.Load()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", config.ConnectionString())

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return nil, err
	}

	return &Migrator{migrate: m}, nil
}

func (m *Migrator) Migrate(version uint) error {
	err := m.Migrate(version)
	if err != nil {
		return err
	}
	return nil
}

func (m *Migrator) Close() error {
	_, err := m.migrate.Close()
	if err != nil {
		return err
	}
	return nil
}
