package service

import (
	"DiscordTemplate/pkg/database"
	"DiscordTemplate/pkg/logger"
)

type exampleServices struct {
	db     database.Database
	logger logger.Logger
}

func NewExampleService(db database.Database, logger logger.Logger) *exampleServices {
	exampleService := &exampleServices{
		db:     db,
		logger: logger,
	}
	exampleService.migrate()
	return exampleService
}

func (s *exampleServices) migrate() {
	err := s.db.ExecSQLFile("./pkg/database/migrations/tables.sql")
	if err != nil {
		s.logger.Fatal("Failed to migrate tables: %v", err)
	}
}
