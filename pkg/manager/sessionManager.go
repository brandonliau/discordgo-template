package manager

import (
	"DiscordTemplate/pkg/command"
	"DiscordTemplate/pkg/component"
	"DiscordTemplate/pkg/logger"
	"DiscordTemplate/pkg/notifier"
	"DiscordTemplate/pkg/shared"

	"github.com/bwmarrin/discordgo"
)

type sessionManager struct {
	session    *discordgo.Session
	logger     logger.Logger
	notifier   notifier.Notifier
	commands   map[string]command.Command
	components map[string]component.Component
}

func NewSessionManager(s *discordgo.Session, logger logger.Logger, notifier notifier.Notifier) *sessionManager {
	return &sessionManager{
		session:    s,
		logger:     logger,
		notifier:   notifier,
		commands:   make(map[string]command.Command),
		components: make(map[string]component.Component),
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

func (m *sessionManager) RegisterComponent(c component.Component) {
	cname := c.CustomID()
	if _, ok := m.components[cname]; ok {
		m.logger.Error("Application component %s already registered", cname)
	}
	m.components[cname] = c
}

func (m *sessionManager) InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else {
		userID = i.User.ID
	}
	cmdArgs := &shared.CmdArgs{
		Session:     s,
		Interaction: i,
		UserID:      userID,
	}
	var rd *discordgo.InteractionResponseData
	var err error
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if command, ok := m.commands[i.ApplicationCommandData().Name]; ok {
			rd, err = command.Execute(cmdArgs)
			if err != nil {
				m.logger.Error("Failed to execute %s: %v", command.Command().Name, err)
			}
			m.logger.Info("%s executed %s", cmdArgs.UserID, i.ApplicationCommandData().Name)
		}
	case discordgo.InteractionMessageComponent:
		if component, ok := m.components[i.MessageComponentData().CustomID]; ok {
			rd, err = component.Execute(cmdArgs)
			if err != nil {
				m.logger.Error("Failed to execute %s: %v", component.CustomID(), err)
			}
			m.logger.Info("%s executed %s", cmdArgs.UserID, component.CustomID())
		}
	}
	err = m.notifier.SendResponse(i, rd)
	if err != nil {
		m.logger.Error("Failed to respond to user %s: %v", userID, err)
	}
}
