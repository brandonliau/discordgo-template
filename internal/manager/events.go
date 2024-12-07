package manager

import (
	"github.com/bwmarrin/discordgo"
)

func (m *discordManager) ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	m.logger.Info("Ready event")
}

func (m *discordManager) ResumedHandler(s *discordgo.Session, r *discordgo.Resumed) {
	m.logger.Info("Resumed event")
}

func (m *discordManager) RateLimitHandler(s *discordgo.Session, r *discordgo.RateLimit) {
	m.logger.Info("Rate limit event")
}
