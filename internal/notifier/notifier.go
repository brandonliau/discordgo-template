package notifier

import (
	"github.com/bwmarrin/discordgo"
)

type Notifier interface {
	SendResponse(i *discordgo.InteractionCreate, rd *discordgo.InteractionResponseData) error
	SendComplexMessage(userID string, data *discordgo.MessageSend) error
}

func MessageSend(c string, e *discordgo.MessageEmbed, mc ...discordgo.MessageComponent) *discordgo.MessageSend {
	if len(mc) == 0 {
		return &discordgo.MessageSend{
			Content: c,
			Embeds:  []*discordgo.MessageEmbed{e},
		}
	}
	return &discordgo.MessageSend{
		Content: c,
		Embeds:  []*discordgo.MessageEmbed{e},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: mc,
			},
		},
	}
}
