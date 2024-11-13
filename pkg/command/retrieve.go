package command

import (
	"fmt"
	"time"

	"DiscordTemplate/pkg/database"

	"github.com/bwmarrin/discordgo"
)

type retrieveCommand struct {
	db database.Database
}

func NewRetrieveCommand(db database.Database) *retrieveCommand {
	return &retrieveCommand{
		db: db,
	}
}

func (c *retrieveCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "retrieve",
		Description: "Retrieve secrets from database.",
	}
}

func (c *retrieveCommand) Execute(args *CmdArgs) (*discordgo.InteractionResponseData, error) {
	rows, err := c.db.Query("SELECT secret FROM userdata WHERE userID = ?", args.UserID)
	if err != nil {
		return nil, err
	}
	var secret string
	var secrets []string
	for rows.Next() {
		rows.Scan(&secret)
		secrets = append(secrets, secret)
	}
	embed := c.retrieveEmbed(secrets...)
	rsp := EphemeralContentResponse("No secrets found for user.")
	if embed != nil {
		rsp = EphemeralEmbedResponse(c.retrieveEmbed(secrets...))
	}
	return rsp, nil
}

func (c *retrieveCommand) retrieveEmbed(secrets ...string) *discordgo.MessageEmbed {
	var desc string
	for _, secret := range secrets {
		desc += fmt.Sprintf("%s\n", secret)
	}
	if desc == "" {
		return nil
	}
	desc = desc[:len(desc)-1]
	return &discordgo.MessageEmbed{
		Title: "Secrets",
		Description: desc,
		Color: blue,
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
