package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//type Database interface {
//	Connection() error
//	Close() error
//	GetCursor() *sqlx.DB
//}
//
//type PostgresDatabase struct {
//	// TODO Заменить на sqlx
//	db *sqlx.DB
//}
//
//func (p *PostgresDatabase) GetCursor() *sqlx.DB {
//	return p.db
//}

//func (p *PostgresDatabase) Connection() error {
//	// TODO Не должно лежать в открытом виде
//	db, err := sqlx.Open("postgres", "postgres://admin:123@localhost:5432/postgres?sslmode=disable")
//	if err != nil {
//		return fmt.Errorf("could not connect to postgres database: %w", err)
//	}
//
//	p.db = db
//	return nil
//}
//func (p *PostgresDatabase) Close() error {
//	return p.db.Close()
//}

func Connection() (*sqlx.DB, error) {
	// TODO Не должно лежать в открытом виде
	db, err := sqlx.Open("postgres", "postgres://admin:123@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres database: %w", err)
	}
	return db, nil
}
