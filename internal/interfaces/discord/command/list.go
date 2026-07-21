package command

import (
	"discordgo-skeleton/internal/application/usecase"
	"discordgo-skeleton/internal/interfaces/discord/interaction"
	"discordgo-skeleton/internal/interfaces/discord/presentation"

	"github.com/bwmarrin/discordgo"
)

func ListDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "list",
		Description: "Show the weather for your pinned zip codes.",
	}
}

func ListHandler(weatherService *usecase.WeatherService) interaction.HandleFunc {
	return func(_ *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
		views, err := weatherService.List(interaction.GetUserID(i))
		if err != nil {
			return nil, err
		}

		if len(views) == 0 {
			embed := presentation.NoticeEmbed("Invalid Request!", "You have no pinned locations.", presentation.Red)
			return interaction.InitialResponse(
				interaction.WithEmbeds(embed),
				interaction.WithEphemeral(),
			)
		}

		return interaction.InitialResponse(
			interaction.WithEmbeds(presentation.PinsEmbed(views)),
			interaction.WithEphemeral(),
		)
	}
}
