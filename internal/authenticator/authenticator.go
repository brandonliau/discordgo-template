package authenticator

import (
	"DiscordTemplate/internal/command"
	"DiscordTemplate/internal/shared"
)

type Authenticator interface {
	Authenticate(cmd command.Command, cmdArgs *shared.CmdArgs) bool
}
