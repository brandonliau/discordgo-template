package manager

import (
	"github.com/bwmarrin/discordgo"
)

func (m *sessionManager) ConnectHandler(s *discordgo.Session, c *discordgo.Connect) {
	m.logger.Info("Connect event")
}

func (m *sessionManager) DisconnectHandler(s *discordgo.Session, d *discordgo.Disconnect) {
	m.logger.Info("Disconnect event")
}

func (m *sessionManager) ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	m.logger.Info("Ready event")
}

func (m *sessionManager) ResumedHandler(s *discordgo.Session, r *discordgo.Resumed) {
	m.logger.Info("Resumed event")
}

func (m *sessionManager) RateLimitHandler(s *discordgo.Session, r *discordgo.RateLimit) {
	m.logger.Info("Rate limit event")
}
