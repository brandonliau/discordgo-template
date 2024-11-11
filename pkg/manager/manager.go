package manager

import (
	"DiscordTemplate/pkg/command"

	"github.com/bwmarrin/discordgo"
)

type Manager interface {
	InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate)
	RegisterCommand(c command.Command)
	SendResponse(i *discordgo.InteractionCreate, rd *discordgo.InteractionResponseData) error
}
