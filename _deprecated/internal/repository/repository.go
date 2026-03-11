package repository

import (
	"github.com/bwmarrin/discordgo"
)

type Repository interface {
	// command registration
	RegisterCommand(name string, id string)
	RetrieveCommands(names ...string) []string

	// component registration
	RegisterComponent(name string, component discordgo.MessageComponent)
	RetrieveComponents(names ...string) []discordgo.MessageComponent
}
