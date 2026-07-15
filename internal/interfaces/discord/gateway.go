package discord

import (
	"errors"
	"fmt"
	"sync"

	"discordgo-template/internal/interfaces/discord/interaction"

	"discordgo-template/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type Gateway struct {
	session       *discordgo.Session
	applicationID string
	guildID       string
	logger        logger.Logger

	mu          sync.RWMutex
	handlers    map[string]interaction.HandleFunc
	commands    []*discordgo.ApplicationCommand
	started     bool
	removeEvent func()
}

func NewGateway(
	session *discordgo.Session,
	applicationID string,
	guildID string,
	logger logger.Logger,
) *Gateway {
	return &Gateway{
		session:       session,
		applicationID: applicationID,
		guildID:       guildID,
		logger:        logger,
		handlers:      make(map[string]interaction.HandleFunc),
	}
}

func (g *Gateway) RegisterCommand(
	definition *discordgo.ApplicationCommand,
	handler interaction.HandleFunc,
) error {
	if definition == nil || definition.Name == "" {
		return errors.New("command definition must have a name")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.started {
		return errors.New("cannot register a command after gateway start")
	}
	if _, exists := g.handlers[definition.Name]; exists {
		return fmt.Errorf("interaction routing key %q is already registered", definition.Name)
	}
	g.handlers[definition.Name] = handler
	g.commands = append(g.commands, definition)
	return nil
}

func (g *Gateway) RegisterComponent(
	definition discordgo.MessageComponent,
	handler interaction.HandleFunc,
) error {
	var customID string
	switch component := definition.(type) {
	case discordgo.Button:
		customID = component.CustomID
	case *discordgo.Button:
		customID = component.CustomID
	case discordgo.SelectMenu:
		customID = component.CustomID
	case *discordgo.SelectMenu:
		customID = component.CustomID
	default:
		return fmt.Errorf("unsupported component definition %T", definition)
	}
	routingKey, err := interaction.RoutingKey(customID)
	if err != nil {
		return fmt.Errorf("component definition: %w", err)
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	if g.started {
		return errors.New("cannot register a component after gateway start")
	}
	if _, exists := g.handlers[routingKey]; exists {
		return fmt.Errorf("interaction routing key %q is already registered", routingKey)
	}
	g.handlers[routingKey] = handler
	return nil
}

func (g *Gateway) Start() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.started {
		return nil
	}

	g.removeEvent = g.session.AddHandler(g.handleInteraction)
	if err := g.session.Open(); err != nil {
		g.removeEvent()
		g.removeEvent = nil
		return fmt.Errorf("open Discord gateway: %w", err)
	}
	if _, err := g.session.ApplicationCommandBulkOverwrite(
		g.applicationID,
		g.guildID,
		g.commands,
	); err != nil {
		g.removeEvent()
		g.removeEvent = nil
		return errors.Join(
			fmt.Errorf("register Discord commands: %w", err),
			g.session.Close(),
		)
	}
	g.started = true
	g.logger.Info("Discord gateway started with %d commands", len(g.commands))
	return nil
}

func (g *Gateway) Stop() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.started {
		return nil
	}
	if g.removeEvent != nil {
		g.removeEvent()
		g.removeEvent = nil
	}
	g.started = false
	if err := g.session.Close(); err != nil {
		return fmt.Errorf("close Discord gateway: %w", err)
	}
	g.logger.Info("Discord gateway stopped")
	return nil
}

func (g *Gateway) Dispatch(
	session *discordgo.Session,
	event *discordgo.InteractionCreate,
) (*discordgo.InteractionResponse, error) {
	customID, err := interactionID(event)
	if err != nil {
		return nil, err
	}
	routingKey, err := interaction.RoutingKey(customID)
	if err != nil {
		return nil, err
	}

	g.mu.RLock()
	handler, exists := g.handlers[routingKey]
	g.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("no handler registered for %q", routingKey)
	}
	return handler(session, event)
}

func (g *Gateway) handleInteraction(
	session *discordgo.Session,
	event *discordgo.InteractionCreate,
) {
	response, err := g.Dispatch(session, event)
	if err != nil {
		g.logger.Error("interaction failed: %v", err)
		response = interaction.InitialResponse(
			interaction.WithContent("Something went wrong."),
			interaction.WithEphemeral(),
		)
	}
	if response == nil {
		return
	}
	if err := session.InteractionRespond(event.Interaction, response); err != nil {
		g.logger.Error("send interaction response: %v", err)
	}
}

func interactionID(event *discordgo.InteractionCreate) (string, error) {
	switch event.Type {
	case discordgo.InteractionApplicationCommand:
		return event.ApplicationCommandData().Name, nil
	case discordgo.InteractionMessageComponent:
		return event.MessageComponentData().CustomID, nil
	case discordgo.InteractionModalSubmit:
		return event.ModalSubmitData().CustomID, nil
	default:
		return "", fmt.Errorf("unsupported interaction type %d", event.Type)
	}
}
