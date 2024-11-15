package component

import (
	"DiscordTemplate/pkg/shared"

	"github.com/bwmarrin/discordgo"
)

type Component interface {
	CustomID() string
	Component() discordgo.MessageComponent
	Execute(args *shared.CmdArgs) (*discordgo.InteractionResponseData, error)
}
