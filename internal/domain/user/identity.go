package user

import (
	"github.com/google/uuid"
)

type Provider string

const (
	ProviderEmail   Provider = "email"
	ProviderDiscord Provider = "discord"
)

type Identity struct {
	Provider   Provider
	ExternalID string
}

type IdentityResolver interface {
	Resolve(provider Provider, externalID string) (uuid.UUID, error)
	Link(userID uuid.UUID, provider Provider, externalID string) error
}
