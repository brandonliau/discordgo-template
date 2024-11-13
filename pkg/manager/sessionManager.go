package manager

import (
	"DiscordTemplate/pkg/command"
	"DiscordTemplate/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type sessionManager struct {
	session  *discordgo.Session
	commands map[string]command.Command
	logger   logger.Logger
}

func NewSessionManager(s *discordgo.Session, logger logger.Logger) *sessionManager {
	return &sessionManager{
		session:  s,
		logger:   logger,
		commands: make(map[string]command.Command),
	}
}

func (m *sessionManager) RegisterCommand(c command.Command) {
	cname := c.Command().Name
	if _, ok := m.commands[cname]; ok {
		m.logger.Error("Application command %s already registered", cname)
	}
	_, err := m.session.ApplicationCommandCreate(m.session.State.User.ID, "", c.Command())
	if err != nil {
		m.logger.Error("Failed to add application command %s : %v", cname, err)
	}
	m.commands[cname] = c
}

func (m *sessionManager) SendResponse(i *discordgo.InteractionCreate, rd *discordgo.InteractionResponseData) error {
	err := m.session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: rd,
	})
	if err != nil {
		return err
	}
	return nil
}
