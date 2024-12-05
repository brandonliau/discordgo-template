package command

import (
	"fmt"

	"DiscordTemplate/internal/shared"
	"DiscordTemplate/pkg/database"

	"github.com/bwmarrin/discordgo"
)

type addCommand struct {
	auth bool
	db   database.Database
}

func NewAddCommand(db database.Database) *addCommand {
	return &addCommand{
		auth: false,
		db:   db,
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

func (c *addCommand) Auth() bool {
	return c.auth
}

func (c *addCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponseData, error) {
	opts := ParseInteractionOptions(args.Interaction.ApplicationCommandData())
	secret := opts["data"]
	c.db.Exec("INSERT INTO userdata (userID, secret) VALUES (?, ?)", args.UserID, secret)
	rsp := shared.EphemeralContentResponse(fmt.Sprintf("Added `%s` to database.", secret))
	return rsp, nil
}
