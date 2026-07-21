package main

import (
	"os"
	"os/signal"
	"syscall"

	"discordgo-skeleton/internal/application/usecase"
	"discordgo-skeleton/internal/config"
	"discordgo-skeleton/internal/infrastructure/openmeteo"
	"discordgo-skeleton/internal/infrastructure/persistence/sqlite"
	"discordgo-skeleton/internal/infrastructure/zippopotam"
	"discordgo-skeleton/internal/interfaces/discord"
	"discordgo-skeleton/internal/interfaces/discord/command"
	"discordgo-skeleton/internal/interfaces/discord/component"

	"discordgo-skeleton/pkg/database"
	"discordgo-skeleton/pkg/logger"

	"github.com/bwmarrin/discordgo"
	_ "modernc.org/sqlite"
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
	db, err := database.NewSqliteDB("./discordgo-skeleton.db")
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

	// Create infrastructure repositories
	pinRepository := sqlite.NewPinRepository(db)

	// Create ports
	geocoder := zippopotam.NewGeocoder()
	weatherFeed := openmeteo.NewWeatherFeed()

	// Create application usecases
	weatherService := usecase.NewWeatherService(geocoder, weatherFeed, pinRepository)
	pinService := usecase.NewPinService(pinRepository, geocoder)

	// Create application gateway
	discordGateway := discord.NewGateway(s, cfg.Discord.ApplicationID, cfg.Discord.GuildID, logger)

	// Register application commands
	err = discordGateway.RegisterCommand(command.RandomDefinition(), command.RandomHandler(weatherService))
	if err != nil {
		logger.Fatal("Failed to register command: %v", err)
	}
	err = discordGateway.RegisterCommand(command.SearchDefinition(), command.SearchHandler(weatherService))
	if err != nil {
		logger.Fatal("Failed to register command: %v", err)
	}
	err = discordGateway.RegisterCommand(command.AddDefinition(), command.AddHandler(pinService))
	if err != nil {
		logger.Fatal("Failed to register command: %v", err)
	}
	err = discordGateway.RegisterCommand(command.RemoveDefinition(), command.RemoveHandler(pinService))
	if err != nil {
		logger.Fatal("Failed to register command: %v", err)
	}
	err = discordGateway.RegisterCommand(command.ListDefinition(), command.ListHandler(weatherService))
	if err != nil {
		logger.Fatal("Failed to register command: %v", err)
	}

	// Register application components
	err = discordGateway.RegisterComponent(component.RefreshDefinition(), component.RefreshHandler(weatherService))
	if err != nil {
		logger.Fatal("Failed to register component: %v", err)
	}

	// Sync commands with discord
	err = discordGateway.Sync()
	if err != nil {
		logger.Fatal("Failed to sync commands: %v", err)
	}

	// Add event handlers
	s.AddHandler(discordGateway.InteractionHandler)

	// Start discord gateway
	err = discordGateway.Start()
	if err != nil {
		logger.Fatal("Failed to start gateway: %v", err)
	}

	// Bot online
	logger.Info("Bot running")

	// Create stop channel and block execution until a stop signal is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Stop discord gateway
	err = discordGateway.Stop()
	if err != nil {
		logger.Fatal("Failed to stop gateway: %v", err)
	}

	logger.Info("Bot shut down")
}
