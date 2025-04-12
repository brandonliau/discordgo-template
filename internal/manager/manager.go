package manager

import (
	"discord-template/internal/command"
	"discord-template/internal/component"
)

type Manager interface {
	RegisterCommand(c command.Command)
	RegisterComponent(c component.Component)
}
