package component

import (
	"DiscordTemplate/pkg/shared"

	"github.com/bwmarrin/discordgo"
)

type linkComponent struct{}

func NewLinkComponent() *linkComponent {
	return &linkComponent{}
}

func (c *linkComponent) CustomID() string {
	return "github"
}

func (c *linkComponent) Component() discordgo.MessageComponent {
	return discordgo.Button{
		Label: "GitHub",
		Style: discordgo.LinkButton,
		URL:   "https://github.com/",
	}
}

func (c *linkComponent) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponseData, error) {
	return nil, nil
}
