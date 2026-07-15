package presentation

import (
	"discordgo-template/internal/application/usecase"

	"github.com/bwmarrin/discordgo"
)

func SampleEmbed(result *usecase.SampleResult) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       result.Title,
		Description: result.Message,
		Color:       0x5865F2,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "State is encoded in the button custom ID.",
		},
	}
}
