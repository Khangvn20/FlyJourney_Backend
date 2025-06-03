package repository

import (
	"context"
	"fmt"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/port/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewPgxDatabase() (repository.Database, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_User"),
		os.Getenv("DB_Pass"),
		os.Getenv("DB_Host"),
		os.Getenv("DB_Port"),
		os.Getenv("DB_Name"),
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database failed: %w", err)
	}
	return &DB{pool: pool}, nil
}
func (db *DB) GetPool() *pgxpool.Pool {
	return db.pool
}

func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}
