package authenticator

import (
	"slices"

	"discord-template/internal/command"
	"discord-template/pkg/config"

	"github.com/bwmarrin/discordgo"
)

type discordAuthenticator struct {
	config  *config.DiscordConfig
	session *discordgo.Session
}

func NewDiscordAuthenticator(cfg config.Config, session *discordgo.Session) *discordAuthenticator {
	return &discordAuthenticator{
		config:  cfg.(*config.DiscordConfig),
		session: session,
	}
}

func (a *discordAuthenticator) Authenticate(cmd command.Command, userID string) bool {
	member, _ := a.session.State.Member(a.config.Guild, userID)
	if !cmd.Auth() {
		return true
	}
	return cmd.Auth() && slices.Contains(member.Roles, a.config.AuthorizedRole)
}
