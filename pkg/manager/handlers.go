package manager

import (
	"github.com/bwmarrin/discordgo"
)

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
