package database

import (
	"database/sql"
)

type Database interface {
	Close()
	Query(query string, args ...any) *sql.Rows
	Exec(query string, args ...any)
	Prepare(query string) *sql.Stmt
	Begin()
	Commit()
	Rollback()
}
