package command

import (
	"DiscordTemplate/internal/shared"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Command() *discordgo.ApplicationCommand
	Auth() bool
	Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error)
}

func ParseInteractionOptions(cid discordgo.ApplicationCommandInteractionData) map[string]string {
	opts := make(map[string]string)
	for _, opt := range cid.Options {
		opts[opt.Name] = opt.Value.(string)
	}
	return opts
}
