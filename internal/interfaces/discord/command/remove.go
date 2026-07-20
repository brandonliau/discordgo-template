package command

import (
	"errors"
	"fmt"

	"discordgo-skeleton/internal/application/usecase"
	"discordgo-skeleton/internal/interfaces/discord/interaction"
	"discordgo-skeleton/internal/interfaces/discord/presentation"

	"github.com/bwmarrin/discordgo"
)

type removeCommand struct {
	pinService *usecase.PinService
}

func RemoveDefinition() *discordgo.ApplicationCommand {
	minLength := 5
	return &discordgo.ApplicationCommand{
		Name:        "remove",
		Description: "Remove a pinned US zip code.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "zip",
				Description: "5 digit US zip code.",
				Required:    true,
				MinLength:   &minLength,
				MaxLength:   5,
			},
		},
	}
}

func RemoveHandler(pinService *usecase.PinService) interaction.HandleFunc {
	c := &removeCommand{
		pinService: pinService,
	}
	return c.execute
}

func (c *removeCommand) execute(_ *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	opts := interaction.ParseInteractionOptions(i)
	userID := interaction.GetUserID(i)
	zip := opts["zip"].StringValue()

	loc, err := c.pinService.Remove(userID, zip)
	switch {
	case err == nil:
		embed := presentation.NoticeEmbed("Pin removed", fmt.Sprintf("Removed the pin for **%s, %s** (`%s`).", loc.City, loc.State, loc.Zip), presentation.Green)
		return interaction.InitialResponse(
			interaction.WithEmbeds(embed),
			interaction.WithEphemeral(),
		)
	case errors.Is(err, usecase.ErrRemoveNotFound):
		embed := presentation.NoticeEmbed("Pin not found", fmt.Sprintf("`%s` is not pinned.", zip), presentation.Red)
		return interaction.InitialResponse(
			interaction.WithEmbeds(embed),
			interaction.WithEphemeral(),
		)
	default:
		return nil, err
	}
}
