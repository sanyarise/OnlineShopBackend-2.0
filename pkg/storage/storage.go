package storage

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Pool *pgxpool.Pool
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

func New(ctx context.Context) (*Storage, error) {
	pg := Storage{}
	conf, err := pgxpool.ParseConfig(config.GetURI())
	if err != nil {
		z.Errorf("can't init storage: %s", err)
		return nil, fmt.Errorf("can't init storage: %s", err)
	}
	err = configurePool(conf)
	if err != nil {
		z.Errorf("can't configure pool %s", err)
		return nil, fmt.Errorf("can't configure pool %s", err)
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		z.Errorf("can't create pool %s", err)
		return nil, fmt.Errorf("can't create pool %s", err)
	}
	pg.Pool = dbPool
	return pg, nil
}
