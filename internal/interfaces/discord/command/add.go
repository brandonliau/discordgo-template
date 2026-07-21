package command

import (
	"errors"
	"fmt"

	"discordgo-skeleton/internal/application/usecase"
	"discordgo-skeleton/internal/interfaces/discord/interaction"
	"discordgo-skeleton/internal/interfaces/discord/presentation"

	"github.com/bwmarrin/discordgo"
)

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
	return func(_ *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
		options := interaction.ParseOptions(i)
		userID := interaction.GetUserID(i)
		zip := options["zip"].StringValue()

		loc, err := pinService.Add(userID, zip)
		switch {
		case err == nil:
			embed := presentation.NoticeEmbed("Successfully Added Pin!", fmt.Sprintf("Pinned **%s, %s** (`%s`).", loc.City, loc.State, loc.Zip), presentation.Green)
			return interaction.InitialResponse(
				interaction.WithEmbeds(embed),
				interaction.WithEphemeral(),
			)
		case errors.Is(err, usecase.ErrAddDuplicate):
			embed := presentation.NoticeEmbed("Invalid Request!", fmt.Sprintf("`%s` is already pinned.", options["zip"].StringValue()), presentation.Red)
			return interaction.InitialResponse(
				interaction.WithEmbeds(embed),
				interaction.WithEphemeral(),
			)
		default:
			return nil, err
		}
	}
}
