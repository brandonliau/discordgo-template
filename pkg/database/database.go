package database

import (
	"database/sql"
)

type Database interface {
	Close()
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) error
	Prepare(query string) (*sql.Stmt, error)
	Begin() error
	Commit() error
	Rollback() error
}
