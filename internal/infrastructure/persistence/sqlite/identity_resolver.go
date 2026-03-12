package sqlite

import (
	"discordgo-template/internal/domain/user"

	"discordgo-template/pkg/database"

	"github.com/google/uuid"
)

type IdentityResolverImpl struct {
	db *database.SqliteDB
}

func NewIdentityResolver(db *database.SqliteDB) *IdentityResolverImpl {
	return &IdentityResolverImpl{db: db}
}

func (r *IdentityResolverImpl) Resolve(provider user.Provider, externalID string) (uuid.UUID, error) {
	row, err := r.db.QueryRow(
		`SELECT id
		 FROM users
		 WHERE provider = ? AND external_id = ?`,
		provider, externalID,
	)
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	err = row.Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Nil, nil
}

func (r *IdentityResolverImpl) Link(userID uuid.UUID, provider user.Provider, externalID string) error {
	return r.db.Exec(
		`INSERT INTO identities (user_id, provider, external_id)
		 VALUES (?, ?, ?)`,
		userID, provider, externalID,
	)
}
