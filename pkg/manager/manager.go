package manager

import (
	"DiscordTemplate/pkg/command"

	"github.com/bwmarrin/discordgo"
)

type Manager interface {
	RegisterCommand(c command.Command)
	InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate)
}
