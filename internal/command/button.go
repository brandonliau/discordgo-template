package command

import (
	"time"

	"DiscordTemplate/internal/component"
	"DiscordTemplate/internal/shared"

	"github.com/bwmarrin/discordgo"
)

type buttonCommand struct {
	auth bool
}

func NewButtonCommand() *buttonCommand {
	return &buttonCommand{
		auth: true,
	}
}

func (c *buttonCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "button",
		Description: "Get example buttons.",
	}
}

func (c *buttonCommand) Auth() bool {
	return c.auth
}

func (c *buttonCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponseData, error) {
	rsp := shared.EphemeralEmbedResponse(c.buttonEmbed())
	rsp = shared.AddComponent(rsp, component.NewPingButton().Component(), component.NewGithubButton().Component())
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
