package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"discord-template/pkg/config"
	"discord-template/pkg/database"
	"discord-template/pkg/logger"

	"discord-template/internal/authenticator"
	"discord-template/internal/command"
	"discord-template/internal/component"
	"discord-template/internal/manager"
	"discord-template/internal/notifier"
	"discord-template/internal/repository"
	"discord-template/internal/service"

	"github.com/bwmarrin/discordgo"
	_ "modernc.org/sqlite"
)

func main() {
	// Create logger, config, database, and service
	logger := logger.NewStdLogger(logger.LevelDebug)
	cfg := config.NewDiscordConfig("./config/config.yml", logger)
	db := database.NewSqliteDB("./database.db", logger)
	defer db.Close()

	// Create new discord session
	s, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		logger.Fatal("Failed to create discord session : %v", err)
	}
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// Create repository, notifier, authenticator, services, and session manager
	repo := repository.NewCommandRepository()
	notifier := notifier.NewDiscordNotifier(s)
	authenticator := authenticator.NewDiscordAuthenticator(cfg, s)
	service := service.NewExampleService(db, logger)
	m := manager.NewDiscordManager(s, repo, notifier, authenticator, logger)

	// Add event handlers
	s.AddHandler(m.CommandInteractionHandler)
	s.AddHandler(m.ComponentInteractionHandler)
	s.AddHandler(m.ReadyHandler)
	s.AddHandler(m.ResumedHandler)
	s.AddHandler(m.RateLimitHandler)

	// Establish websocket connection
	err = s.Open()
	if err != nil {
		logger.Fatal("Failed to establish websocket connection : %v", err)
	}
	defer s.Close()

	// Start service
	err = service.Start()
	if err != nil {
		logger.Fatal("Failed to start service: %v", err)
	}

	// Register application commands
	m.RegisterCommand(command.NewPingCommand())
	m.RegisterCommand(command.NewUptimeCommand(time.Now().Unix()))
	m.RegisterCommand(command.NewAddCommand(db))
	m.RegisterCommand(command.NewClearCommand(db))
	m.RegisterCommand(command.NewRetrieveCommand(db, logger))
	m.RegisterCommand(command.NewButtonCommand())
	m.RegisterCommand(command.NewCleanCommand(notifier))

	// Register application components
	m.RegisterComponent(component.NewPingButton())

	// Bot online
	logger.Info("Bot running")

	// Create stop channel and block execution until a stop signal is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Stop service
	err = service.Stop()
	if err != nil {
		logger.Error("Failed to stop service: %v", err)
	}

	// Remove application commands
	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", nil)
	if err != nil {
		logger.Error("Failed to delete application commands")
	}
	logger.Info("Bot shut down")
}
