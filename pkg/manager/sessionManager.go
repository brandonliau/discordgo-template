package manager

import (
	"DiscordTemplate/pkg/command"
	"DiscordTemplate/pkg/logger"
	"DiscordTemplate/pkg/notifier"

	"github.com/bwmarrin/discordgo"
)

type sessionManager struct {
	session  *discordgo.Session
	commands map[string]command.Command
	logger   logger.Logger
	notifier notifier.Notifier
}

func NewSessionManager(s *discordgo.Session, logger logger.Logger, notifier notifier.Notifier) *sessionManager {
	return &sessionManager{
		session:  s,
		logger:   logger,
		notifier: notifier,
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

func (m *sessionManager) InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else {
		userID = i.User.ID
	}
	cmdArgs := &command.CmdArgs{
		Session:     s,
		Interaction: i,
		UserID:      userID,
	}
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if command, ok := m.commands[i.ApplicationCommandData().Name]; ok {
			rd, err := command.Execute(cmdArgs)
			if err != nil {
				m.logger.Error("Failed to execute %s: %v", command.Command().Name, err)
			}
			err = m.notifier.SendResponse(i, rd)
			if err != nil {
				m.logger.Warn("Failed to respond to user %s: %v", userID, err)
			}
		}
		m.logger.Debug("%s executed %s", cmdArgs.UserID, i.ApplicationCommandData().Name)
	}
}
