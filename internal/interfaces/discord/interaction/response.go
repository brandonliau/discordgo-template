package interaction

import "github.com/bwmarrin/discordgo"

type ResponseOption func(*discordgo.InteractionResponseData)

func Response(responseType discordgo.InteractionResponseType, opts ...ResponseOption) *discordgo.InteractionResponse {
	data := &discordgo.InteractionResponseData{}
	for _, opt := range opts {
		opt(data)
	}
	return &discordgo.InteractionResponse{Type: responseType, Data: data}
}

func InitialResponse(opts ...ResponseOption) *discordgo.InteractionResponse {
	return Response(discordgo.InteractionResponseChannelMessageWithSource, opts...)
}

func UpdateResponse(opts ...ResponseOption) *discordgo.InteractionResponse {
	return Response(discordgo.InteractionResponseUpdateMessage, opts...)
}

func WithContent(content string) ResponseOption {
	return func(data *discordgo.InteractionResponseData) {
		data.Content = content
	}
}

func WithEmbeds(embeds ...*discordgo.MessageEmbed) ResponseOption {
	return func(data *discordgo.InteractionResponseData) {
		data.Embeds = embeds
	}
}

func WithComponents(components ...discordgo.MessageComponent) ResponseOption {
	return func(data *discordgo.InteractionResponseData) {
		data.Components = []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: components},
		}
	}
}

func WithEphemeral() ResponseOption {
	return func(data *discordgo.InteractionResponseData) {
		data.Flags |= discordgo.MessageFlagsEphemeral
	}
}
