package manager

import (
	"DiscordTemplate/pkg/command"
	"DiscordTemplate/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type sessionManager struct {
	session  *discordgo.Session
	commands map[string]command.Command
	logger logger.Logger
}

func NewSessionManager(s *discordgo.Session, logger logger.Logger) *sessionManager {
	return &sessionManager{
		session:  s,
		logger: logger,
		commands: make(map[string]command.Command),
	}
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
				m.logger.Error("Failed to execute %s: %v", command.GetCommand().Name, err)
			}
			err = m.SendResponse(i, rd)
			if err != nil {
				m.logger.Warn("Failed to message user %s: %v", userID, err)
			}
		}
		m.logger.Debug("%s executed %s", cmdArgs.UserID, i.ApplicationCommandData().Name)
	}
}

func (m *sessionManager) RegisterCommand(c command.Command) {
	cname := c.GetCommand().Name
	if _, ok := m.commands[cname]; ok {
		m.logger.Error("Application command %s already registered", cname)
	}
	_, err := m.session.ApplicationCommandCreate(m.session.State.User.ID, "", c.GetCommand())
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
