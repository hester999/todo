package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database interface {
	Connection() error
	Close() error
	Exec(query string, args ...interface{}) error
	Query(query string, dest interface{}, args ...interface{}) error
	QueryRow(query string, dest interface{}, args ...interface{}) error
}

type PostgresDatabase struct {
	// TODO Заменить на sqlx
	db *sqlx.DB
}

func (p *PostgresDatabase) Connection() error {
	// TODO Не должно лежать в открытом виде
	db, err := sqlx.Open("postgres", "postgres://admin:123@localhost:5432/postgres?sslmode=disable")
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

func (p *PostgresDatabase) Query(query string, dest interface{}, args ...interface{}) error {
	err := p.db.Select(dest, query, args...)
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}
	return err
}

func (p *PostgresDatabase) QueryRow(query string, dest interface{}, args ...interface{}) error {
	err := p.db.Get(dest, query, args...)
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}
	return nil
}
