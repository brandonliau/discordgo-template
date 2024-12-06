package manager

import (
	"DiscordTemplate/internal/authenticator"
	"DiscordTemplate/internal/command"
	"DiscordTemplate/internal/component"
	"DiscordTemplate/internal/notifier"
	"DiscordTemplate/internal/shared"
	"DiscordTemplate/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type discordManager struct {
	session       *discordgo.Session
	logger        logger.Logger
	authenticator authenticator.Authenticator
	notifier      notifier.Notifier
	commands      map[string]command.Command
	components    map[string]component.Component
}

func NewDiscordManager(s *discordgo.Session, logger logger.Logger, authenticator authenticator.Authenticator, notifier notifier.Notifier) *discordManager {
	return &discordManager{
		session:       s,
		logger:        logger,
		authenticator: authenticator,
		notifier:      notifier,
		commands:      make(map[string]command.Command),
		components:    make(map[string]component.Component),
	}
}

func (m *discordManager) RegisterCommand(c command.Command) {
	cname := c.Command().Name
	if _, ok := m.commands[cname]; ok {
		m.logger.Error("Application command %s already registered", cname)
	}
	_, err := m.session.ApplicationCommandCreate(m.session.State.User.ID, "", c.Command())
	if err != nil {
		m.logger.Error("Failed to add application command %s : %v", cname, err)
	}
	m.logger.Debug("Registered command %v", cname)
	m.commands[cname] = c
}

func (m *discordManager) RegisterComponent(c component.Component) {
	cname := c.CustomID()
	if _, ok := m.components[cname]; ok {
		m.logger.Error("Application component %s already registered", cname)
	}
	m.logger.Debug("Registered component %v", cname)
	m.components[cname] = c
}

func (m *discordManager) InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			if !m.authenticator.Authenticate(command, userID) {
				rd = shared.EphemeralContentResponse("Authorized users only!")
				m.logger.Info("Unauthorized user %s attempted to execute /%s", userID, i.ApplicationCommandData().Name)
				break
			}
			rd, err = command.Execute(cmdArgs)
			if err != nil {
				m.logger.Error("Failed to execute /%s: %v", command.Command().Name, err)
			}
			m.logger.Info("%s executed /%s", userID, i.ApplicationCommandData().Name)
		}
	case discordgo.InteractionMessageComponent:
		if component, ok := m.components[i.MessageComponentData().CustomID]; ok {
			rd, err = component.Execute(cmdArgs)
			if err != nil {
				m.logger.Error("Failed to execute /%s: %v", component.CustomID(), err)
			}
			m.logger.Info("%s executed /%s", userID, component.CustomID())
		}
	}
	err = m.notifier.SendResponse(i, rd)
	if err != nil {
		m.logger.Error("Failed to respond to user %s: %v", userID, err)
	}
}
