package component

import (
	"discord-template/internal/shared"

	"github.com/bwmarrin/discordgo"
)

type githubButton struct{}

func NewGithubButton() *githubButton {
	return &githubButton{}
}

func (c *githubButton) Component() discordgo.MessageComponent {
	return discordgo.Button{
		Label: "GitHub",
		Style: discordgo.LinkButton,
		URL:   "https://github.com/",
	}
}

func (c *githubButton) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
	return nil, nil
}
