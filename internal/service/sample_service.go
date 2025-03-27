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

func (s *exampleServices) migrate() error {
	err := s.db.ExecSQLFile("./pkg/database/migrations/tables.sql")
	return err
}

func (s *exampleServices) Start() error {
	err := s.migrate()
	if err != nil {
		return err
	}
	return nil
}

func (s *exampleServices) Stop() error {
	return nil
}
