package command

import (
	"DiscordTemplate/pkg/database"
	"DiscordTemplate/pkg/shared"

	"github.com/bwmarrin/discordgo"
)

type clearCommand struct {
	db database.Database
}

func NewClearCommand(db database.Database) *clearCommand {
	return &clearCommand{
		db: db,
	}
}

func (c *clearCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "clear",
		Description: "Clear secrets from database.",
	}
}

func (c *clearCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponseData, error) {
	err := c.db.Exec("DELETE FROM userdata WHERE userID = ?", args.UserID)
	if err != nil {
		return nil, err
	}
	rsp := shared.EphemeralContentResponse("Removed all secrets from database.")
	return rsp, nil
}
