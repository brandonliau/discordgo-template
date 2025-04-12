package database

import (
	"database/sql"
	"os"

	"discord-template/pkg/logger"
)

type SqliteDB struct {
	readDB  *sql.DB
	writeDB *sql.DB
	logger  logger.Logger
}

func NewSqliteDB(file string, logger logger.Logger) *SqliteDB {
	readDB, err := sql.Open("sqlite", file)
	if err != nil {
		logger.Fatal("Failed to open read database connection: %v", err)
	}
	writeDB, err := sql.Open("sqlite", file)
	if err != nil {
		logger.Fatal("Failed to open write database connection: %v", err)
	}
	sqliteDB := &SqliteDB{
		readDB:  readDB,
		writeDB: writeDB,
		logger:  logger,
	}
	err = sqliteDB.applyPerformanceOptions(readDB, 4, "./pkg/database/migrations/perf.sql")
	if err != nil {
		logger.Fatal("Failed to apply database performance options on readDB: %v", err)
	}
	err = sqliteDB.applyPerformanceOptions(writeDB, 1, "./pkg/database/migrations/perf.sql")
	if err != nil {
		logger.Fatal("Failed to apply database performance options on writeDB: %v", err)
	}
	return sqliteDB
}

func (db *SqliteDB) applyPerformanceOptions(sdb *sql.DB, maxOpenConns int, file string) error {
	sdb.SetMaxOpenConns(maxOpenConns)
	err := db.ExecSQLFile(file)
	if err != nil {
		return err
	}
	return nil
}

func (db *SqliteDB) Close() error {
	err := db.readDB.Close()
	if err != nil {
		db.logger.Error("Failed to close read database connection: %v", err)
		return err
	}
	err = db.writeDB.Close()
	if err != nil {
		db.logger.Error("Failed to close write database connection: %v", err)
		return err
	}
	return nil
}

func (db *SqliteDB) ExecSQLFile(file string) error {
	sqlContent, err := os.ReadFile(file)
	if err != nil {
		db.logger.Info("File: %s", file)
		db.logger.Error("Failed to read file: %v", err)
		return err
	}
	_, err = db.writeDB.Exec(string(sqlContent))
	if err != nil {
		db.logger.Info("Query: %s", sqlContent)
		db.logger.Error("Failed to execute query: %v", err)
		return err
	}
	return nil
}

func (db *SqliteDB) Exec(query string, args ...any) error {
	_, err := db.writeDB.Exec(query, args...)
	if err != nil {
		db.logger.Info("Query: %s", query)
		db.logger.Error("Failed to execute query: %v", err)
		return err
	}
	return nil
}

func (db *SqliteDB) Query(query string, args ...any) (*sql.Rows, error) {
	rows, err := db.readDB.Query(query, args...)
	if err != nil {
		db.logger.Info("Query: %s", query)
		db.logger.Error("Failed to query database: %v", err)
		return nil, err
	}
	return rows, nil
}

func (db *SqliteDB) QueryRow(query string, args ...any) (*sql.Row, error) {
	row := db.readDB.QueryRow(query, args...)
	return row, nil
}

func (db *SqliteDB) PrepareExec(query string) (*sql.Stmt, error) {
	stmt, err := db.writeDB.Prepare(query)
	if err != nil {
		db.logger.Info("Query: %s", query)
		db.logger.Error("Failed to prepare write query: %v", err)
		return nil, err
	}
	return stmt, nil
}

func (db *SqliteDB) PrepareQuery(query string) (*sql.Stmt, error) {
	stmt, err := db.readDB.Prepare(query)
	if err != nil {
		db.logger.Info("Query: %s", query)
		db.logger.Error("Failed to prepare read query: %v", err)
		return nil, err
	}
	return stmt, nil
}

func (db *SqliteDB) Begin() error {
	_, err := db.writeDB.Exec("BEGIN IMMEDIATE")
	if err != nil {
		db.logger.Error("Failed to begin transaction: %v", err)
		return err
	}
	return nil
}

func (db *SqliteDB) Commit() error {
	_, err := db.writeDB.Exec("COMMIT")
	if err != nil {
		db.logger.Error("Failed to commit transaction: %v", err)
		db.Rollback()
		return err
	}
	return nil
}

func (db *SqliteDB) Rollback() error {
	_, err := db.writeDB.Exec("ROLLBACK")
	if err != nil {
		db.logger.Error("Failed to rollback transaction: %v", err)
		return err
	}
	return nil
}
