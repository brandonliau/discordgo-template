package command

import (
	"DiscordTemplate/pkg/database"

	"github.com/bwmarrin/discordgo"
)

type removeCommand struct {
	db database.Database
}

func NewRemoveCommand(db database.Database) *removeCommand {
	return &removeCommand{
		db: db,
	}
}

func (c *removeCommand) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "remove",
		Description: "Remove secrets from database.",
	}
}

func (c *removeCommand) Execute(args *CmdArgs) (*discordgo.InteractionResponseData, error) {
	err := c.db.Remove(args.UserID)
	if err != nil {
		return nil, err
	}
	rsp := EphemeralContentResponse("Removed all secrets from database.")
	return rsp, nil
}
