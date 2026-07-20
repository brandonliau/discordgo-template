package discord

import (
	"fmt"
	"strings"
	
	"discordgo-skeleton/internal/interfaces/discord/interaction"

	"discordgo-skeleton/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type gateway struct {
	session       *discordgo.Session
	applicationID string
	guildID       string
	handleFuncs   map[string]interaction.HandleFunc
	logger        logger.Logger
}

func NewGateway(session *discordgo.Session, applicationID string, guildID string, logger logger.Logger) *gateway {
	return &gateway{
		session:       session,
		applicationID: applicationID,
		guildID:       guildID,
		handleFuncs:   make(map[string]interaction.HandleFunc),
		logger:        logger,
	}
}

func (g *gateway) Start() error {
	err := g.session.Open()
	if err != nil {
		return err
	}
	g.logger.Info("Started discord gateway for application %s", g.applicationID)
	return nil
}

func (g *gateway) Stop() error {
	err := g.session.Close()
	if err != nil {
		return err
	}
	g.logger.Info("Stopped discord gateway")
	return nil
}

func (g *gateway) RegisterCommand(c *discordgo.ApplicationCommand, handleFunc interaction.HandleFunc) error {
	if _, ok := g.handleFuncs[c.Name]; ok {
		return fmt.Errorf("Command %s already registered", c.Name)
	}

	_, err := g.session.ApplicationCommandCreate(g.applicationID, "", c)
	if err != nil {
		return err
	}

	g.handleFuncs[c.Name] = handleFunc
	g.logger.Info("Registered command %s", c.Name)
	return nil
}

func (g *gateway) RegisterComponent(c discordgo.MessageComponent, handleFunc interaction.HandleFunc) error {
	var customID string
	switch v := c.(type) {
	case discordgo.Button:
		customID = v.CustomID
	case *discordgo.Button:
		customID = v.CustomID
	case discordgo.SelectMenu:
		customID = v.CustomID
	case *discordgo.SelectMenu:
		customID = v.CustomID
	default:
		return fmt.Errorf("Unsupported component type %T", c)
	}

	routingKey, _, _ := strings.Cut(customID, "?")
	if _, ok := g.handleFuncs[routingKey]; ok {
		return fmt.Errorf("Component %s already registered", routingKey)
	}

	g.handleFuncs[customID] = handleFunc
	g.logger.Info("Registered component %s", routingKey)
	return nil
}

func (g *gateway) RegisterModal(modal *discordgo.InteractionResponseData, handleFunc interaction.HandleFunc) error {
	if _, ok := g.handleFuncs[modal.CustomID]; ok {
		return fmt.Errorf("Modal %s already registered", modal.CustomID)
	}

	g.handleFuncs[modal.CustomID] = handleFunc
	g.logger.Info("Registered modal %s", modal.CustomID)
	return nil
}
