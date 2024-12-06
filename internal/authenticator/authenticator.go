package authenticator

import (
	"DiscordTemplate/internal/command"
)

type Authenticator interface {
	Authenticate(cmd command.Command, userID string) bool
}
