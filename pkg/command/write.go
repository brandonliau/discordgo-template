package command

import (
	"fmt"

	"DiscordTemplate/pkg/database"

	"github.com/bwmarrin/discordgo"
)

type writeCommand struct {
	db database.Database
}

func NewWriteCommand(db database.Database) *writeCommand {
	return &writeCommand{
		db: db,
	}
}

func (c *writeCommand) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "write",
		Description: "Write to database.",
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

func (c *writeCommand) Execute(args *CmdArgs) (*discordgo.InteractionResponseData, error) {
	opts := ParseInteractionOptions(args.Interaction.ApplicationCommandData())
	secret := opts["data"]
	err := c.db.Write(args.UserID, secret)
	if err != nil {
		return nil, err
	}
	rsp := EphemeralContentResponse(fmt.Sprintf("Wrote `%s` to database.", secret))
	return rsp, nil
}
