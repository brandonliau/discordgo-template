package component

import (
	"DiscordTemplate/pkg/command"

	"github.com/bwmarrin/discordgo"
)

type Component interface {
	CustomID() string
	Component() discordgo.MessageComponent
	Execute(args *command.CmdArgs) (*discordgo.InteractionResponseData, error)
}
