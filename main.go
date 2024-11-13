package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"DiscordTemplate/pkg/command"
	"DiscordTemplate/pkg/config"
	"DiscordTemplate/pkg/database"
	"DiscordTemplate/pkg/logger"
	"DiscordTemplate/pkg/manager"
	_ "DiscordTemplate/pkg/service"

	"github.com/bwmarrin/discordgo"
	_ "modernc.org/sqlite"
)

func main() {
	logger := logger.NewStdLogger(logger.LevelDebug)
	cfg := config.NewYamlConfig("config.yml", logger)
	db := database.NewSqliteDB(logger)
	defer db.Close()

	// Create new discord session
	s, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		logger.Fatal("Failed to create discord session : %v", err)
	}
	m := manager.NewSessionManager(s, logger)

	// Identify intents and add handlers
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	s.AddHandler(m.InteractionHandler)
	s.AddHandler(m.ConnectHandler)
	s.AddHandler(m.ReadyHandler)
	s.AddHandler(m.ResumedHandler)
	s.AddHandler(m.RateLimitHandler)

	// Establish websocket connection
	err = s.Open()
	if err != nil {
		logger.Fatal("Failed to establish websocket connection : %v", err)
	}
	defer s.Close()

	// Register application commands
	m.RegisterCommand(command.NewPingCommand())
	m.RegisterCommand(command.NewUptimeCommand(time.Now().Unix()))
	m.RegisterCommand(command.NewAddCommand(db))
	m.RegisterCommand(command.NewClearCommand(db))
	m.RegisterCommand(command.NewRetrieveCommand(db))
	s.UpdateCustomStatus("👁️‍🗨️ Monitoring...")
	logger.Info("Bot running")

	// Create stop channel
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Remove application commands
	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", nil)
	if err != nil {
		logger.Error("Failed to delete application commands")
	}
	logger.Info("Bot shut down")
}
