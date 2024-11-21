package database

import (
	"database/sql"

	"DiscordTemplate/pkg/logger"
)

type sqliteDB struct {
	readDB  *sql.DB
	writeDB *sql.DB
	logger  logger.Logger
}

func NewSqliteDB(logger logger.Logger) *sqliteDB {
	readDB, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		logger.Fatal("Failed to open read database connection: %v", err)
	}
	err = applyPerformanceOptions(readDB, 4)
	if err != nil {
		logger.Fatal("Failed to apply database performance options: %v", err)
	}

	writeDB, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		logger.Fatal("Failed to open write database connection: %v", err)
	}
	err = applyPerformanceOptions(readDB, 1)
	if err != nil {
		logger.Fatal("Failed to apply database performance options: %v", err)
	}

	sqliteDB := &sqliteDB{
		readDB:  readDB,
		writeDB: writeDB,
		logger:  logger,
	}

	err = sqliteDB.Exec("CREATE TABLE IF NOT EXISTS userdata (userID TEXT, secret TEXT)")
	if err != nil {
		logger.Error("Failed to initialize database: %v", err)
	}

	return sqliteDB
}

func applyPerformanceOptions(db *sql.DB, maxOpenConns int) error {
	db.SetMaxOpenConns(maxOpenConns)
	_, err := db.Exec(`
		PRAGMA synchronous = NORMAL;
		PRAGMA journal_mode = WAL;
		PRAGMA foreign_keys = OFF;
		PRAGMA busy_timeout = 5000;
	`)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqliteDB) Close() {
	err := s.readDB.Close()
	if err != nil {
		s.logger.Error("Failed to close read database connection: %v", err)
	}
	err = s.writeDB.Close()
	if err != nil {
		s.logger.Error("Failed to close write database connection: %v", err)
	}
}

func (s *sqliteDB) Query(query string, args ...any) (*sql.Rows, error) {
	rows, err := s.readDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *sqliteDB) Exec(query string, args ...any) error {
	_, err := s.writeDB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqliteDB) Prepare(query string) (*sql.Stmt, error) {
	stmt, err := s.writeDB.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (s *sqliteDB) Begin() error {
	_, err := s.writeDB.Exec("BEGIN IMMEDIATE")
	if err != nil {
		return err
	}
	return nil
}

func (s *sqliteDB) Commit() error {
	_, err := s.writeDB.Exec("COMMIT")
	if err != nil {
		return err
	}
	return nil
}

func (s *sqliteDB) Rollback() error {
	_, err := s.writeDB.Exec("ROLLBACK")
	if err != nil {
		return err
	}
	return nil
}
