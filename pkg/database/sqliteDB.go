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

func (db *sqliteDB) Close() {
	err := db.readDB.Close()
	if err != nil {
		db.logger.Error("Failed to close read database connection: %v", err)
	}
	err = db.writeDB.Close()
	if err != nil {
		db.logger.Error("Failed to close write database connection: %v", err)
	}
}

func (db *sqliteDB) Query(query string, args ...any) *sql.Rows {
	rows, err := db.readDB.Query(query, args...)
	if err != nil {
		db.logger.Info("Query: %s", query)
		db.logger.Error("Failed to query database: %v", err)
		return nil
	}
	return rows
}

func (db *sqliteDB) Exec(query string, args ...any) {
	_, err := db.writeDB.Exec(query, args...)
	if err != nil {
		db.logger.Info("Query: %s", query)
		db.logger.Error("Failed to execute query: %v", err)
	}
}

func (db *sqliteDB) Prepare(query string) *sql.Stmt {
	stmt, err := db.writeDB.Prepare(query)
	if err != nil {
		db.logger.Info("Query: %s", query)
		db.logger.Error("Failed to prepare query: %v", err)
		return nil
	}
	return stmt
}

func (db *sqliteDB) Begin() {
	_, err := db.writeDB.Exec("BEGIN IMMEDIATE")
	if err != nil {
		db.logger.Error("Failed to begin transaction: %v", err)
	}
}

func (db *sqliteDB) Commit() {
	_, err := db.writeDB.Exec("COMMIT")
	if err != nil {
		db.logger.Error("Failed to commit transaction: %v", err)
		db.Rollback()
	}
}

func (db *sqliteDB) Rollback() {
	_, err := db.writeDB.Exec("ROLLBACK")
	if err != nil {
		db.logger.Error("Failed to rollback transaction: %v", err)
	}
}
