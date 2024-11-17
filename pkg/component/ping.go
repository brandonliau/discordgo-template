package component

import (
	"fmt"

	"DiscordTemplate/pkg/shared"

	"github.com/bwmarrin/discordgo"
)

type pingComponent struct{}

func NewPingComponent() *pingComponent {
	return &pingComponent{}
}

func (c *pingComponent) CustomID() string {
	return c.Component().(discordgo.Button).CustomID
}

func (c *pingComponent) Component() discordgo.MessageComponent {
	return discordgo.Button{
		Label:    "Ping",
		Style:    discordgo.PrimaryButton,
		CustomID: "component_ping",
	}
}

func (c *pingComponent) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponseData, error) {
	rsp := shared.EphemeralContentResponse(fmt.Sprintf("Pong! `%d ms`", args.Session.HeartbeatLatency().Milliseconds()))
	return rsp, nil
}
