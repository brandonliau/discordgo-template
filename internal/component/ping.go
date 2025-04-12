package component

import (
	"fmt"

	"discord-template/internal/shared"

	"github.com/bwmarrin/discordgo"
)

type pingButton struct{}

func NewPingButton() *pingButton {
	return &pingButton{}
}

func (c *pingButton) CustomID() string {
	return c.Component().(discordgo.Button).CustomID
}

func (c *pingButton) Component() discordgo.MessageComponent {
	return discordgo.Button{
		Label:    "Ping",
		Style:    discordgo.PrimaryButton,
		CustomID: "ping",
	}
}

func (c *pingButton) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
	rsp := shared.EphemeralContentResponse(fmt.Sprintf("Pong! `%d ms`", args.Session.HeartbeatLatency().Milliseconds()))
	return rsp, nil
}
