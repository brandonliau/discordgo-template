package command

import (
	"time"

	"DiscordTemplate/pkg/component"
	"DiscordTemplate/pkg/shared"

	"github.com/bwmarrin/discordgo"
)

type buttonCommand struct{}

func NewButtonCommand() *buttonCommand {
	return &buttonCommand{}
}

func (c *buttonCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "button",
		Description: "Get example buttons.",
	}
}

func (c *buttonCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponseData, error) {
	rsp := shared.EphemeralEmbedResponse(c.buttonEmbed())
	rsp = shared.AddComponent(rsp, component.NewPingComponent().Component(), component.NewLinkComponent().Component())
	return rsp, nil
}

func (c *buttonCommand) buttonEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Buttons",
		Description: "Here are some sample buttons",
		Color:       blue,
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
