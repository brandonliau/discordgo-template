package command

import (
	"github.com/bwmarrin/discordgo"
)

const (
	blue  = 0x5865f2
	green = 0x2dcc70
	red   = 0xe74d3b
)

type Command interface {
	GetCommand() *discordgo.ApplicationCommand
	Execute(args *CmdArgs) (*discordgo.InteractionResponseData, error)
}

type CmdArgs struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	UserID      string
}

func ContentResponse(c string) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Content: c,
	}
}

func EphemeralContentResponse(c string) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Flags:   discordgo.MessageFlagsEphemeral,
		Content: c,
	}
}

func EmbedResponse(e *discordgo.MessageEmbed) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{e},
	}
}

func EphemeralEmbedResponse(e *discordgo.MessageEmbed) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Flags:  discordgo.MessageFlagsEphemeral,
		Embeds: []*discordgo.MessageEmbed{e},
	}
}

func ParseInteractionOptions(opts discordgo.ApplicationCommandInteractionData) map[string]string {
	icd := make(map[string]string)
	for _, opt := range opts.Options {
		icd[opt.Name] = opt.Value.(string)
	}
	return icd
}
