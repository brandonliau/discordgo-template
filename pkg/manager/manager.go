package manager

import (
	"DiscordTemplate/pkg/command"
	"DiscordTemplate/pkg/component"

	"github.com/bwmarrin/discordgo"
)

type Manager interface {
	RegisterCommand(c command.Command)
	RegisterComponent(c component.Component)
	InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate)
}
