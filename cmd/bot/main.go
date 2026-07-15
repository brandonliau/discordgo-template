package main

import (
	"os"
	"os/signal"
	"syscall"

	"discordgo-template/internal/application/usecase"
	"discordgo-template/internal/application/worker"
	"discordgo-template/internal/config"
	"discordgo-template/internal/infrastructure/persistence/sqlite"
	"discordgo-template/internal/interfaces/discord"
	"discordgo-template/internal/interfaces/discord/command"
	"discordgo-template/internal/interfaces/discord/component"

	"discordgo-template/pkg/database"
	"discordgo-template/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Create logger
	logger := logger.NewStdLogger(logger.LevelInfo)

	// Create config
	cfg, err := config.Load("./config/config.yml")
	if err != nil {
		logger.Fatal("Failed to create config: %v", err)
	}

	// Create database
	db, err := database.NewSqliteDB(cfg.Database.Path)
	if err != nil {
		logger.Fatal("Failed to create database: %v", err)
	}
	defer db.Close()

	// Perform database migrations
	err = sqlite.Migrate(db)
	if err != nil {
		logger.Fatal("Failed to migrate database: %v", err)
	}

	// Create discord session
	s, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		logger.Fatal("Failed to create discord session: %v", err)
	}
	s.Identify.Intents = discordgo.IntentsGuilds

	// Create application usecases
	sampleService := usecase.NewSampleService()

	// Create application gateway
	discordGateway := discord.NewGateway(
		s,
		cfg.Discord.ApplicationID,
		cfg.Discord.GuildID,
		logger,
	)

	// Register application commands
	err = discordGateway.RegisterCommand(
		command.SampleDefinition(),
		command.SampleHandler(sampleService),
	)
	if err != nil {
		logger.Fatal("Failed to register application command: %v", err)
	}

	// Register application components
	sampleButton, err := component.SampleDefinition(0)
	if err != nil {
		logger.Fatal("Failed to create sample component: %v", err)
	}
	err = discordGateway.RegisterComponent(
		sampleButton,
		component.SampleHandler(sampleService),
	)
	if err != nil {
		logger.Fatal("Failed to register application component: %v", err)
	}

	// Register application services
	orchestrator := worker.NewOrchestrator(logger)
	orchestrator.RegisterWorker(
		"sample",
		worker.NewSampleWorker(cfg.SampleWorker.Interval.Duration, logger),
	)

	// Start gateway
	err = discordGateway.Start()
	if err != nil {
		logger.Fatal("Failed to start discord gateway: %v", err)
	}

	// Start workers
	err = orchestrator.StartAll()
	if err != nil {
		logger.Fatal("Failed to start services: %v", err)
	}

	// Bot online
	logger.Info("Bot running")

	// Create stop channel and block execution until a stop signal is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Stop workers
	err = orchestrator.StopAll()
	if err != nil {
		logger.Error("Failed to stop services: %v", err)
	}

	// Stop gateway
	err = discordGateway.Stop()
	if err != nil {
		logger.Error("Failed to stop discord gateway: %v", err)
	}

	logger.Info("Bot shut down")
}
