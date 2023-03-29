package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	Pool *pgxpool.Pool
}

func NewStore(dns string) (*Store, error) {
	config, err := pgxpool.ParseConfig(dns)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &Store{Pool: conn}, nil
}
