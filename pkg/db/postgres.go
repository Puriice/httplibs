package db

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabase() (*pgxpool.Pool, error) {
	connectionString := os.Getenv("DB_URL")

	if connectionString == "" {
		return nil, errors.New("Invalid connection string.")
	}

	db, err := pgxpool.New(context.Background(), connectionString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
