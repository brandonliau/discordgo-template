package repository

import (
	"github.com/bwmarrin/discordgo"
)

type commandRepo struct {
	registeredCommands   map[string]string
	registeredComponents map[string]discordgo.MessageComponent
}

func NewCommandRepository() *commandRepo {
	return &commandRepo{
		registeredCommands:   make(map[string]string),
		registeredComponents: make(map[string]discordgo.MessageComponent),
	}
}

// command registration
func (repo *commandRepo) RegisterCommand(name string, id string) {
	repo.registeredCommands[name] = id
}

func (repo *commandRepo) RetrieveCommands(names ...string) []string {
	commands := make([]string, len(names))
	for i, name := range names {
		commands[i] = repo.registeredCommands[name]
	}
	return commands
}

// component registration
func (repo *commandRepo) RegisterComponent(name string, component discordgo.MessageComponent) {
	repo.registeredComponents[name] = component
}

func (repo *commandRepo) RetrieveComponents(names ...string) []discordgo.MessageComponent {
	components := make([]discordgo.MessageComponent, len(names))
	for i, name := range names {
		components[i] = repo.registeredComponents[name]
	}
	return components
}
