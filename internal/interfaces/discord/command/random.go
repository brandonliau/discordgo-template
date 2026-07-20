package command

import (
	"discordgo-skeleton/internal/application/usecase"
	"discordgo-skeleton/internal/interfaces/discord/component"
	"discordgo-skeleton/internal/interfaces/discord/interaction"
	"discordgo-skeleton/internal/interfaces/discord/presentation"

	"github.com/bwmarrin/discordgo"
)

type randomCommand struct {
	weatherService *usecase.WeatherService
}

func RandomDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "random",
		Description: "Show the weather for a random US location.",
	}
}

func RandomHandler(weatherService *usecase.WeatherService) interaction.HandleFunc {
	c := &randomCommand{
		weatherService: weatherService,
	}
	return c.execute
}

func (c *randomCommand) execute(_ *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	w, err := c.weatherService.Random()
	if err != nil {
		return nil, err
	}

	return interaction.InitialResponse(
		interaction.WithEmbeds(presentation.WeatherEmbed(w)),
		interaction.WithComponents(component.RefreshDefinition()),
	)
}
