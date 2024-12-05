package command

import (
	"DiscordTemplate/internal/shared"

	"github.com/bwmarrin/discordgo"
)

const (
	blue  = 0x5865f2
	green = 0x2dcc70
	red   = 0xe74d3b
)

type Command interface {
	Command() *discordgo.ApplicationCommand
	Auth() bool
	Execute(args *shared.CmdArgs) (*discordgo.InteractionResponseData, error)
}

func ParseInteractionOptions(cid discordgo.ApplicationCommandInteractionData) map[string]string {
	opts := make(map[string]string)
	for _, opt := range cid.Options {
		opts[opt.Name] = opt.Value.(string)
	}
	return opts
}
