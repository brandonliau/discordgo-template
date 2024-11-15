package shared

import (
	"github.com/bwmarrin/discordgo"
)

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

func AddComponent(rsp *discordgo.InteractionResponseData, c ...discordgo.MessageComponent) *discordgo.InteractionResponseData {
	rsp.Components = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: c,
		},
	}
	return rsp
}
