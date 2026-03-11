package user

import (
	"github.com/google/uuid"
)

// todo: implement identity to map internal uuids to external (discord) user IDs
type Identity struct {
	Provider   string // "discord", "email", etc.
	ExternalID string
}

type IdentityResolver interface {
	Resolve(provider string, externalID string) (uuid.UUID, error)
	// question: [userID string] or [id string]
	// question: should i make link on user so that the interfaces layer handles linking?
	Link(userID uuid.UUID, provider string, externalID string) error
}
