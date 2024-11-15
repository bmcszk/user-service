package db

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

func InitDB(ctx context.Context, postgresUrl string) (*pgx.Conn, error) {
	err := migrateUp(postgresUrl)
	if err != nil {
		return nil, err
	}
	conn, err := pgx.Connect(ctx, postgresUrl)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func migrateUp(postgresUrl string) error {
	m, err := migrate.New(
		"file://db/migrations",
		postgresUrl)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
