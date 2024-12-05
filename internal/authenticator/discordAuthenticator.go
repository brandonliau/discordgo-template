package authenticator

import (
	"slices"

	"DiscordTemplate/internal/command"
	"DiscordTemplate/internal/shared"
	"DiscordTemplate/pkg/config"
)

type discordAuthenticator struct {
	config *config.DiscordConfig
}

func NewDiscordAuthenticator(cfg config.Config) *discordAuthenticator {
	return &discordAuthenticator{
		config: cfg.(*config.DiscordConfig),
	}
}

func (a *discordAuthenticator) Authenticate(cmd command.Command, cmdArgs *shared.CmdArgs) bool {
	member, _ := cmdArgs.Session.State.Member(a.config.Guild, cmdArgs.UserID)
	if !cmd.Auth() {
		return true
	}
	return cmd.Auth() && slices.Contains(member.Roles, a.config.AuthorizedRole)
}
