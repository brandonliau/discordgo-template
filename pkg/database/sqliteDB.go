package database

import (
	"database/sql"

	"DiscordTemplate/pkg/logger"
)

type sqliteDB struct {
	db     *sql.DB
	logger logger.Logger
}

func NewSqliteDB(logger logger.Logger) *sqliteDB {
	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	db.SetMaxOpenConns(1)
	sqlitedb := &sqliteDB{
		db:     db,
		logger: logger,
	}
	err = sqlitedb.Migrate()
	if err != nil {
		logger.Warn("Failed to initialize database: %v", err)
	}
	return sqlitedb
}

func (s *sqliteDB) Close() {
	err := s.db.Close()
	if err != nil {
		s.logger.Error("Failed to close database connection")
	}
}

func (s *sqliteDB) Migrate() error {
	query := "CREATE TABLE IF NOT EXISTS userdata (userID TEXT, secret TEXT)"
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqliteDB) Query(query string, args ...any) (*sql.Rows, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *sqliteDB) Exec(query string, args ...any) error {
	_, err := s.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}
