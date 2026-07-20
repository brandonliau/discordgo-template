package command

import (
	"errors"
	"fmt"

	"discordgo-skeleton/internal/application/usecase"
	"discordgo-skeleton/internal/interfaces/discord/interaction"
	"discordgo-skeleton/internal/interfaces/discord/presentation"

	"github.com/bwmarrin/discordgo"
)

type searchCommand struct {
	weatherService *usecase.WeatherService
}

func SearchDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "search",
		Description: "Show the weather for a US zip code.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "zip",
				Description: "5 digit US zip code.",
				Required:    true,
				MinLength:   new(5),
				MaxLength:   5,
			},
		},
	}
}

func SearchHandler(weatherService *usecase.WeatherService) interaction.HandleFunc {
	c := &searchCommand{
		weatherService: weatherService,
	}
	return c.execute
}

func (c *searchCommand) execute(_ *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	opts := interaction.ParseInteractionOptions(i)
	zip := opts["zip"].StringValue()

	w, err := c.weatherService.Search(zip)
	switch {
	case err == nil:
		embed := presentation.WeatherEmbed(w)
		return interaction.InitialResponse(
			interaction.WithEmbeds(embed),
			interaction.WithEphemeral(),
		)
	case errors.Is(err, usecase.ErrSearchZipInvalid):
		embed := presentation.NoticeEmbed("Invalid zip", fmt.Sprintf("`%s` is not a valid US zip code.", zip), presentation.Red)
		return interaction.InitialResponse(
			interaction.WithEmbeds(embed),
			interaction.WithEphemeral(),
		)
	default:
		return nil, err
	}
}
