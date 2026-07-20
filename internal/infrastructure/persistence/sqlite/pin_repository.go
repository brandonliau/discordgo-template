package sqlite

import (
	"errors"

	"discordgo-skeleton/internal/domain/pin"

	"discordgo-skeleton/pkg/database"

	"modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"
)

var _ pin.Repository = (*PinRepository)(nil)

type PinRepository struct {
	db *database.SqliteDB
}

func NewPinRepository(db *database.SqliteDB) *PinRepository {
	return &PinRepository{
		db: db,
	}
}

func (r *PinRepository) Create(pinned pin.Pin) error {
	err := r.db.Exec(
		`INSERT INTO saved_locations (user_id, zip) VALUES (?, ?)`,
		pinned.UserID,
		pinned.Zip,
	)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE {
			return pin.ErrPinDuplicate
		}
		return err
	}
	return nil
}

func (r *PinRepository) Delete(userID, zip string) error {
	affected, err := r.db.ExecAffected(
		`DELETE FROM saved_locations WHERE user_id = ? AND zip = ?`,
		userID,
		zip,
	)
	if err != nil {
		return err
	}
	if affected == 0 {
		return pin.ErrPinNotFound
	}
	return nil
}

func (r *PinRepository) ListByUser(userID string) ([]*pin.Pin, error) {
	rows, err := r.db.Query(
		`SELECT user_id, zip FROM saved_locations WHERE user_id = ? ORDER BY zip ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pins []*pin.Pin
	for rows.Next() {
		var pinned pin.Pin
		if err := rows.Scan(&pinned.UserID, &pinned.Zip); err != nil {
			return nil, err
		}
		pins = append(pins, &pinned)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pins, nil
}
