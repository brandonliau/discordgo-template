package command

import (
	"fmt"

	"DiscordTemplate/pkg/database"
	"DiscordTemplate/pkg/shared"

	"github.com/bwmarrin/discordgo"
)

type addCommand struct {
	db database.Database
}

func NewAddCommand(db database.Database) *addCommand {
	return &addCommand{
		db: db,
	}
}

func (c *addCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "add",
		Description: "Add a secret to the database.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "data",
				Description: "data",
				Required:    true,
			},
		},
	}
}

func (c *addCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponseData, error) {
	opts := ParseInteractionOptions(args.Interaction.ApplicationCommandData())
	secret := opts["data"]
	err := c.db.Exec("INSERT INTO userdata (userID, secret) VALUES (?, ?)", args.UserID, secret)
	if err != nil {
		return nil, err
	}
	rsp := shared.EphemeralContentResponse(fmt.Sprintf("Added `%s` to database.", secret))
	return rsp, nil
}
