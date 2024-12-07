package command

import (
	"fmt"
	"time"

	"DiscordTemplate/internal/shared"
	"DiscordTemplate/pkg/database"
	"DiscordTemplate/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type retrieveCommand struct {
	auth   bool
	db     database.Database
	logger logger.Logger
}

func NewRetrieveCommand(db database.Database, logger logger.Logger) *retrieveCommand {
	return &retrieveCommand{
		auth:   false,
		db:     db,
		logger: logger,
	}
}

func (c *retrieveCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "retrieve",
		Description: "Retrieve secrets from database.",
	}
}

func (c *retrieveCommand) Auth() bool {
	return c.auth
}

func (c *retrieveCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
	rows, _ := c.db.Query("SELECT secret FROM userdata WHERE userID = ?", args.UserID)
	defer rows.Close()
	var secret string
	var secrets []string
	for rows.Next() {
		err := rows.Scan(&secret)
		if err != nil {
			c.logger.Error("Failed to scan row: %v", err)
		}
		secrets = append(secrets, secret)
	}
	embed := c.retrieveEmbed(secrets...)
	rsp := shared.EphemeralContentResponse("No secrets found for user.")
	if embed != nil {
		rsp = shared.EphemeralEmbedResponse(c.retrieveEmbed(secrets...))
	}
	return rsp, nil
}

func (c *retrieveCommand) retrieveEmbed(secrets ...string) *discordgo.MessageEmbed {
	var description string
	for _, secret := range secrets {
		description += fmt.Sprintf("%s\n", secret)
	}
	if description == "" {
		return nil
	}
	description = description[:len(description)-1]
	return &discordgo.MessageEmbed{
		Title:       "Secrets",
		Description: description,
		Color:       shared.Blue,
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
