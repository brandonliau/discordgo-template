package database

import (
	"database/sql"
)

type Database interface {
	Close()
	Migrate() error
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) error
}
