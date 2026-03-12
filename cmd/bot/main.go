package main

import (
	"os"
	"os/signal"
	"syscall"

	"discordgo-template/internal/application/usecase"
	"discordgo-template/internal/config"
	"discordgo-template/internal/infrastructure/external"
	"discordgo-template/internal/infrastructure/persistence/sqlite"
	"discordgo-template/internal/interfaces/discord"
	"discordgo-template/internal/interfaces/discord/command"

	"discordgo-template/pkg/database"
	"discordgo-template/pkg/logger"

	"github.com/bwmarrin/discordgo"
	_ "modernc.org/sqlite"
)

func main() {
	// Create logger
	logger := logger.NewStdLogger(logger.LevelInfo)

	// Create config
	cfg, err := config.NewYamlConfig("./config/config.yml")
	if err != nil {
		logger.Fatal("Failed to create config: %v", err)
	}

	// Create database
	db, err := database.NewSqliteDB("./database.db")
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
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// Create infrastructure repositories
	userRepository := sqlite.NewUserRepository(db)
	identityResolver := sqlite.NewIdentityResolver(db)

	// Create ports
	systemMonitor := external.NewSystemMonitor()

	// Create application usescases
	userService := usecase.NewUserService(userRepository)
	systemService := usecase.NewSystemService(systemMonitor)

	// Create application gateways
	discordGateway := discord.NewGateway(s, cfg.Discord.ApplicationID, cfg.Discord.GuildID, userService, identityResolver, cfg.Discord, logger)

	// Register application commands
	discordGateway.RegisterCommand(command.StatusCommandDefinition(), command.StatusCommandHandler(systemService))

	// Add event handlers
	s.AddHandler(discordGateway.InteractionHandler)
	s.AddHandler(discordGateway.ReadyHandler)
	s.AddHandler(discordGateway.ResumedHandler)
	s.AddHandler(discordGateway.RateLimitHandler)
	s.AddHandler(discordGateway.MemberJoinHandler)
	s.AddHandler(discordGateway.MemberLeaveHandler)

	// Start gateway
	discordGateway.Start()

	// Establish websocket connection
	err = s.Open()
	if err != nil {
		logger.Fatal("Failed to establish websocket connection : %v", err)
	}
	defer s.Close()

	// Bot online
	logger.Info("Bot running")

	// Create stop channel and block execution until a stop signal is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Stop application gateways
	discordGateway.Stop()

	// Remove application commands
	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", nil)
	if err != nil {
		logger.Error("Failed to delete application commands")
	}
	logger.Info("Bot shut down")
}
