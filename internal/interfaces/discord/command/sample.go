package command

import (
	"discordgo-skeleton/internal/application/usecase"
	"discordgo-skeleton/internal/interfaces/discord/component"
	"discordgo-skeleton/internal/interfaces/discord/interaction"
	"discordgo-skeleton/internal/interfaces/discord/presentation"

	"github.com/bwmarrin/discordgo"
)

func SampleDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "sample",
		Description: "Show the template's stateful interaction sample.",
	}
}

func SampleHandler(service *usecase.SampleService) interaction.HandleFunc {
	return func(_ *discordgo.Session, _ *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
		result, err := service.Get(0)
		if err != nil {
			return nil, err
		}
		button, err := component.SampleDefinition(0)
		if err != nil {
			return nil, err
		}
		return interaction.InitialResponse(
			interaction.WithEmbeds(presentation.SampleEmbed(result)),
			interaction.WithComponents(button),
		), nil
	}
}
