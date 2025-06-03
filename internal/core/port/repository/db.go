package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	GetPool() *pgxpool.Pool
	Ping(ctx context.Context) error
}
