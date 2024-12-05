package manager

import (
	"github.com/bwmarrin/discordgo"
)

func (m *discordManager) ConnectHandler(s *discordgo.Session, c *discordgo.Connect) {
	m.logger.Info("Connect event")
}

func (m *discordManager) DisconnectHandler(s *discordgo.Session, d *discordgo.Disconnect) {
	m.logger.Info("Disconnect event")
}

func (m *discordManager) ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	m.logger.Info("Ready event")
}

func (m *discordManager) ResumedHandler(s *discordgo.Session, r *discordgo.Resumed) {
	m.logger.Info("Resumed event")
}

func (m *discordManager) RateLimitHandler(s *discordgo.Session, r *discordgo.RateLimit) {
	m.logger.Info("Rate limit event")
}
