package authenticator

import (
	"discord-template/internal/command"
)

type Authenticator interface {
	Authenticate(cmd command.Command, userID string) bool
}
