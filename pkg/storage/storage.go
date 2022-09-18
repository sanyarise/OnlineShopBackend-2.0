package storage

import (
	"context"
	"fmt"
	"net"
	"time"

	"online_shop_backend/pkg/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Pool *pgxpool.Pool
}

func New() *Storage {
	return &Storage{}
}

func configurePool(conf *pgxpool.Config) (err error) {
	// add cofiguration
	conf.MaxConns = int32(10)
	conf.MinConns = int32(5)

	conf.HealthCheckPeriod = 1 * time.Minute
	conf.MaxConnLifetime = 24 * time.Hour
	conf.MaxConnIdleTime = 30 * time.Minute
	conf.ConnConfig.ConnectTimeout = 1 * time.Second
	conf.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: conf.HealthCheckPeriod,
		Timeout:   conf.ConnConfig.ConnectTimeout,
	}).DialContext
	return nil
}

func (pg *Storage) Start(ctx context.Context, config config.Config) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("done with context")
	default:
		conf, err := pgxpool.ParseConfig(config.DSN)
		if err != nil {
			// z.Errorf("can't init storage: %s", err)
			return fmt.Errorf("can't init storage: %w", err)
		}
		err = configurePool(conf)
		if err != nil {
			// z.Errorf("can't configure pool %s", err)
			return fmt.Errorf("can't configure pool %w", err)
		}

		dbPool, err := pgxpool.NewWithConfig(ctx, conf)
		if err != nil {
			// z.Errorf("can't create pool %s", err)
			return fmt.Errorf("can't create pool %w", err)
		}
		pg.Pool = dbPool
		return nil
	}
}

func (pg *Storage) ShutDown(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("done with context")
	default:
		pg.Pool.Close()
		return nil
	}
}

func (pg *Storage) GetName() string {
	return "storage"
}
