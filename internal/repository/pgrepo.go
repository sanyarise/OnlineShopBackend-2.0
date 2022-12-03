package repository

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

type PGres struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

// Returns new empty connection
func NewPgxStorage(logger *zap.SugaredLogger) *PGres {
	return &PGres{
		logger: logger,
	}
}

var _ Storage = (*PGres)(nil)

// used to create configuration for connection from config
func configurePool(conf *pgxpool.Config) (err error) {
	// add cofiguration
	conf.MaxConns = int32(15)
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

// Configurates connection to get ready for work
func (pg *PGres) InitStorage(ctx context.Context, dns string) (Storage, error) {

	conf, err := pgxpool.ParseConfig(dns)
	if err != nil {
		pg.logger.Errorf("can't init storage: %s", err)
		return nil, fmt.Errorf("can't init storage: %w", err)
	}
	err = configurePool(conf)
	if err != nil {
		pg.logger.Errorf("can't configure pool %s", err)
		return nil, fmt.Errorf("can't configure pool %w", err)
	}

	dbPool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		pg.logger.Errorf("can't create pool %s", err)
		return nil, fmt.Errorf("can't create pool %w", err)
	}
	pg.pool = dbPool
	// pg.logger = pg.logger
	return pg, nil
}

// Returns pool to make queries
func (pg *PGres) GetPool() *pgxpool.Pool {
	return pg.pool
}
