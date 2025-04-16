package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connection() (*sqlx.DB, error) {
	// TODO Не должно лежать в открытом виде
	db, err := sqlx.Open("postgres", "postgres://admin:123@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres database: %w", err)
	}
	return db, nil
}
