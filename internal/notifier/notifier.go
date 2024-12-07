package notifier

import (
	"github.com/bwmarrin/discordgo"
)

type Notifier interface {
	SendResponse(i *discordgo.InteractionCreate, ir *discordgo.InteractionResponse) error
	SendComplexMessage(channelID string, data *discordgo.MessageSend) error
	CreateDMChannel(userID string) (string, error)
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
