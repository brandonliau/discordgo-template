package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type pingCommand struct{}

func NewPingCommand() *pingCommand {
	return &pingCommand{}
}

func (c *pingCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Check bot latency.",
	}
}

func (c *pingCommand) Execute(args *CmdArgs) (*discordgo.InteractionResponseData, error) {
	rsp := EphemeralContentResponse(fmt.Sprintf("Pong! `%d ms`", args.Session.HeartbeatLatency().Milliseconds()))
	return rsp, nil
}
