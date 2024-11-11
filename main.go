package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"DiscordTemplate/pkg/command"
	"DiscordTemplate/pkg/config"
	"DiscordTemplate/pkg/database"
	"DiscordTemplate/pkg/manager"
	_ "DiscordTemplate/pkg/service"

	"github.com/bwmarrin/discordgo"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.NewYamlConfig("config.yml")
	db := database.NewSqliteDB()
	defer db.Close()
	
	// Create new discord session
	s, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatalf("[FATAL] Failed to create discord session : %v", err)
	}
	m := manager.NewSessionManager(s)

	// Identify intents and add handlers
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	s.AddHandler(m.InteractionHandler)

	// Establish websocket connection
	err = s.Open()
	if err != nil {
		log.Fatalf("[FATAL] Failed to establish websocket connection : %v", err)
	}
	defer s.Close()

	// Register application commands
	m.RegisterCommand(command.NewPingCommand())
	m.RegisterCommand(command.NewUptimeCommand(time.Now().Unix()))
	m.RegisterCommand(command.NewWriteCommand(db))
	m.RegisterCommand(command.NewRemoveCommand(db))
	s.UpdateCustomStatus("👁️‍🗨️ Monitoring...")
	log.Println("[INFO] Bot running")

	// Create stop channel
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Remove application commands
	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", nil)
	if err != nil {
		log.Printf("[ERROR] Failed to delete application commands")
	}
	log.Println("[INFO] Bot shut down")
}
