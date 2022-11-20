package repository

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Pgrepo struct {
	db *sql.DB
}

func NewPgrepo(dns string) (*Pgrepo, error) {
	log.Println("Enter in NewPgrepo()")
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Pgrepo{db: db}, nil
}
