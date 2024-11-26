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
	sqliteDB.migrate()

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

func (s *sqliteDB) migrate() {
	s.Exec("CREATE TABLE IF NOT EXISTS userdata (userID TEXT, secret TEXT)")
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

func (s *sqliteDB) Query(query string, args ...any) *sql.Rows {
	rows, err := s.readDB.Query(query, args...)
	if err != nil {
		s.logger.Debug("Query: %s", query)
		s.logger.Error("Failed to query database: %v", err)
		return nil
	}
	return rows
}

func (s *sqliteDB) Exec(query string, args ...any) {
	_, err := s.writeDB.Exec(query, args...)
	if err != nil {
		s.logger.Debug("Query: %s", query)
		s.logger.Error("Failed to execute query: %v", err)
	}
}

func (s *sqliteDB) Prepare(query string) *sql.Stmt {
	stmt, err := s.writeDB.Prepare(query)
	if err != nil {
		s.logger.Debug("Query: %s", query)
		s.logger.Error("Failed to prepare query: %v", err)
		return nil
	}
	return stmt
}

func (s *sqliteDB) Begin() {
	_, err := s.writeDB.Exec("BEGIN IMMEDIATE")
	if err != nil {
		s.logger.Error("Failed to begin transaction: %v", err)
	}
}

func (s *sqliteDB) Commit() {
	_, err := s.writeDB.Exec("COMMIT")
	if err != nil {
		s.logger.Error("Failed to commit transaction: %v", err)
		s.Rollback()
	}
}

func (s *sqliteDB) Rollback() {
	_, err := s.writeDB.Exec("ROLLBACK")
	if err != nil {
		s.logger.Error("Failed to rollback transaction: %v", err)
	}
}
