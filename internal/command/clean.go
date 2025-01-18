package command

import (
	"time"

	"DiscordTemplate/internal/notifier"
	"DiscordTemplate/internal/shared"

	"github.com/bwmarrin/discordgo"
)

type cleanCommand struct {
	notifier notifier.Notifier
	auth     bool
}

func NewCleanCommand(notifier notifier.Notifier) *cleanCommand {
	return &cleanCommand{
		notifier: notifier,
		auth:     true,
	}
}

func (c *cleanCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "clean",
		Description: "Clear bot DMs.",
	}
}

func (c *cleanCommand) Auth() bool {
	return c.auth
}

func (c *cleanCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
	defer func() {
		channel, _ := c.notifier.CreateDMChannel(args.UserID)
		messages, _ := args.Session.ChannelMessages(channel, 100, "", "", "")
		for len(messages) != 0 {
			for _, msg := range messages {
				args.Session.ChannelMessageDelete(channel, msg.ID)
				time.Sleep(700 * time.Millisecond)
			}
			messages, _ = args.Session.ChannelMessages(channel, 100, "", "", "")
		}
	}()
	rsp := shared.EphemeralContentResponse("Clearing all DMs...")
	return rsp, nil
}
