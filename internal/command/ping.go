package command

import (
	"fmt"

	"DiscordTemplate/internal/shared"

	"github.com/bwmarrin/discordgo"
)

type pingCommand struct {
	auth bool
}

func NewPingCommand() *pingCommand {
	return &pingCommand{
		auth: false,
	}
}

func (c *pingCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Check bot latency.",
	}
}

func (c *pingCommand) Auth() bool {
	return c.auth
}

func (c *pingCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponseData, error) {
	rsp := shared.EphemeralContentResponse(fmt.Sprintf("Pong! `%d ms`", args.Session.HeartbeatLatency().Milliseconds()))
	return rsp, nil
}
