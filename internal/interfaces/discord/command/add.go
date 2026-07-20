package command

import (
	"errors"
	"fmt"

	"discordgo-skeleton/internal/application/usecase"
	"discordgo-skeleton/internal/interfaces/discord/interaction"
	"discordgo-skeleton/internal/interfaces/discord/presentation"

	"github.com/bwmarrin/discordgo"
)

type addCommand struct {
	pinService *usecase.PinService
}

func AddDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "add",
		Description: "Pin a US zip code.",
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

func AddHandler(pinService *usecase.PinService) interaction.HandleFunc {
	c := &addCommand{
		pinService: pinService,
	}
	return c.execute
}

func (c *addCommand) execute(_ *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	options := interaction.ParseInteractionOptions(i)
	userID := interaction.GetUserID(i)
	zip := options["zip"].StringValue()

	loc, err := c.pinService.Add(userID, zip)
	switch {
	case err == nil:
		embed := presentation.NoticeEmbed("Pin added", fmt.Sprintf("Pinned **%s, %s** (`%s`).", loc.City, loc.State, loc.Zip), presentation.Green)
		return interaction.InitialResponse(
			interaction.WithEmbeds(embed),
			interaction.WithEphemeral(),
		)
	case errors.Is(err, usecase.ErrAddDuplicate):
		embed := presentation.NoticeEmbed("Duplicate pin", fmt.Sprintf("`%s` is already pinned.", options["zip"].StringValue()), presentation.Red)
		return interaction.InitialResponse(
			interaction.WithEmbeds(embed),
			interaction.WithEphemeral(),
		)
	default:
		return nil, err
	}
}
