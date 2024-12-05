package manager

import (
	"DiscordTemplate/internal/command"
	"DiscordTemplate/internal/component"
)

type Manager interface {
	RegisterCommand(c command.Command)
	RegisterComponent(c component.Component)
}
