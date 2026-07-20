package component

import (
	"discordgo-skeleton/internal/application/usecase"
	"discordgo-skeleton/internal/interfaces/discord/interaction"
	"discordgo-skeleton/internal/interfaces/discord/presentation"

	"github.com/bwmarrin/discordgo"
)

const RefreshRoutingKey = "weather.refresh"

type refreshComponent struct {
	weatherService *usecase.WeatherService
}

func RefreshDefinition() discordgo.Button {
	return discordgo.Button{
		CustomID: "refresh",
		Style:    discordgo.PrimaryButton,
		Emoji:    &discordgo.ComponentEmoji{
			Name: "🔄",
		},
	}
}

func RefreshHandler(weatherService *usecase.WeatherService) interaction.HandleFunc {
	c := &refreshComponent{
		weatherService: weatherService,
	}
	return c.execute
}

func (c *refreshComponent) execute(_ *discordgo.Session, _ *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	wv, err := c.weatherService.Random()
	if err != nil {
		return nil, err
	}
	return interaction.UpdateResponse(
		interaction.WithEmbeds(presentation.WeatherEmbed(wv)),
		interaction.WithComponents(RefreshDefinition()),
	)
}
