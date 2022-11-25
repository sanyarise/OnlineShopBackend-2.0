package repository

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

type Pgrepo struct {
	db *sql.DB
	logger  *zap.Logger
}

func NewPgrepo(dns string, logger *zap.Logger) (*Pgrepo, error) {
	logger.Debug("Enter in NewPgrepo()")
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Pgrepo{db: db, logger: logger}, nil
}
