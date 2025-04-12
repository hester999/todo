package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Database interface {
	Connection() error
	Close() error
	Exec(query string, args ...interface{}) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type PostgresDatabase struct {
	// TODO Заменить на sqlx
	db *sql.DB
}

func (p *PostgresDatabase) Connection() error {
	// TODO Не должно лежать в открытом виде
	db, err := sql.Open("postgres", "postgres://admin:123@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return fmt.Errorf("could not connect to postgres database: %w", err)
	}

	p.db = db
	return nil
}
func (p *PostgresDatabase) Close() error {
	return p.db.Close()
}

func (p *PostgresDatabase) Exec(query string, args ...interface{}) error {
	_, err := p.db.Exec(query, args...)

	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}
	return nil
}

func (p *PostgresDatabase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	res, err := p.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	return res, err
}
