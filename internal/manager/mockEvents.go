package manager

import (
	"github.com/bwmarrin/discordgo"
)

func (m *mockManager) ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	m.session.UpdateCustomStatus("👁️‍🗨️ Monitoring...")
	m.application = r.Application
	m.logger.Info("Ready event")
}

func (m *mockManager) ResumedHandler(s *discordgo.Session, r *discordgo.Resumed) {
	m.logger.Info("Resumed event")
}

func (m *mockManager) RateLimitHandler(s *discordgo.Session, r *discordgo.RateLimit) {
	m.logger.Info("Rate limit event")
}
