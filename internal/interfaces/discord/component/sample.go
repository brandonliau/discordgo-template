package component

import (
	"fmt"
	"net/url"
	"strconv"

	"discordgo-skeleton/internal/application/usecase"
	"discordgo-skeleton/internal/interfaces/discord/interaction"
	"discordgo-skeleton/internal/interfaces/discord/presentation"

	"github.com/bwmarrin/discordgo"
)

const SampleRoutingKey = "sample.increment"

func SampleDefinition(count int) (discordgo.Button, error) {
	customID, err := interaction.EncodeCustomID(
		SampleRoutingKey,
		url.Values{"count": {strconv.Itoa(count)}},
	)
	if err != nil {
		return discordgo.Button{}, err
	}
	return discordgo.Button{
		Label:    "Increment",
		Style:    discordgo.PrimaryButton,
		CustomID: customID,
	}, nil
}

func SampleHandler(service *usecase.SampleService) interaction.HandleFunc {
	return func(_ *discordgo.Session, event *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
		_, params, err := interaction.DecodeCustomID(event.MessageComponentData().CustomID)
		if err != nil {
			return nil, err
		}
		count, err := strconv.Atoi(params.Get("count"))
		if err != nil {
			return nil, fmt.Errorf("parse sample count: %w", err)
		}
		count++

		result, err := service.Get(count)
		if err != nil {
			return nil, err
		}
		button, err := SampleDefinition(count)
		if err != nil {
			return nil, err
		}
		return interaction.UpdateResponse(
			interaction.WithEmbeds(presentation.SampleEmbed(result)),
			interaction.WithComponents(button),
		), nil
	}
}
