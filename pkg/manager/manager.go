package manager

import (
	"DiscordTemplate/pkg/command"
	"DiscordTemplate/pkg/component"
)

type Manager interface {
	RegisterCommand(c command.Command)
	RegisterComponent(c component.Component)
}
