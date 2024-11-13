package manager

import (
	"DiscordTemplate/pkg/command"

	"github.com/bwmarrin/discordgo"
)

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
			err = m.SendResponse(i, rd)
			if err != nil {
				m.logger.Warn("Failed to message user %s: %v", userID, err)
			}
		}
		m.logger.Debug("%s executed %s", cmdArgs.UserID, i.ApplicationCommandData().Name)
	}
}

func (m *sessionManager) ConnectHandler(s *discordgo.Session, c *discordgo.Connect) {
	m.logger.Debug("Connect event")
}

func (m *sessionManager) DisconnectHandler(s *discordgo.Session, d *discordgo.Disconnect) {
	m.logger.Debug("Disconnect event")
}

func (m *sessionManager) ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	m.logger.Debug("Ready event")
}

func (m *sessionManager) ResumedHandler(s *discordgo.Session, r *discordgo.Resumed) {
	m.logger.Debug("Resumed event")
}

func (m *sessionManager) RateLimitHandler(s *discordgo.Session, r *discordgo.RateLimit) {
	m.logger.Debug("Rate limit event")
}
