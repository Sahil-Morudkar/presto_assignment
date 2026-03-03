package db

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB() (*pgxpool.Pool, error) {

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	// Parse pool configuration
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	// Verify DB connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	log.Println("✅ Database connected successfully")

	return pool, nil
}