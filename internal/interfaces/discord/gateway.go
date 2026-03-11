package discord

import (
	"strings"

	"discordgo-template/internal/application/usecase"
	"discordgo-template/internal/config"
	"discordgo-template/internal/interfaces/discord/interaction"

	"discordgo-template/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type gateway struct {
	session       *discordgo.Session
	applicationID string
	guildID       string
	userService   *usecase.UserService
	handleFuncs   map[string]interaction.HandleFunc
	cfg           *config.DiscordConfig
	logger        logger.Logger
}

func NewGateway(session *discordgo.Session, applicationID string, guildID string, userService *usecase.UserService, cfg *config.DiscordConfig, logger logger.Logger) *gateway {
	return &gateway{
		session:       session,
		applicationID: applicationID,
		guildID:       guildID,
		userService:   userService,
		handleFuncs:   make(map[string]interaction.HandleFunc),
		cfg:           cfg,
		logger:        logger,
	}
}

func (g *gateway) Start() {
	g.logger.Info("Started discord gateway for application %s", g.applicationID)
}

func (g *gateway) Stop() {
	g.logger.Info("Stopped discord gateway")
}

func (g *gateway) RegisterCommand(def *discordgo.ApplicationCommand, handleFunc interaction.HandleFunc) {
	if _, ok := g.handleFuncs[def.Name]; ok {
		g.logger.Warn("Command %s already registered", def.Name)
		return
	}

	_, err := g.session.ApplicationCommandCreate(g.applicationID, "", def)
	if err != nil {
		g.logger.Error("Failed to register command %s: %v", def.Name, err)
		return
	}

	g.handleFuncs[def.Name] = handleFunc
	g.logger.Info("Registered command %s", def.Name)
}

func (g *gateway) RegisterComponent(def discordgo.MessageComponent, handleFunc interaction.HandleFunc) {
	var customID string
	switch v := def.(type) {
	case discordgo.Button:
		customID = v.CustomID
	case discordgo.SelectMenu:
		customID = v.CustomID
	}

	routingKey, _, _ := strings.Cut(customID, "?")
	if _, ok := g.handleFuncs[routingKey]; ok {
		g.logger.Warn("Component %s already registered", routingKey)
		return
	}

	g.handleFuncs[customID] = handleFunc
	g.logger.Info("Registered component %s", routingKey)
}

func (g *gateway) RegisterModal(def *discordgo.InteractionResponseData, handleFunc interaction.HandleFunc) {
	if _, ok := g.handleFuncs[def.CustomID]; ok {
		g.logger.Warn("Modal %s already registered", def.CustomID)
		return
	}

	g.handleFuncs[def.CustomID] = handleFunc
	g.logger.Info("Registered modal %s", def.CustomID)
}
