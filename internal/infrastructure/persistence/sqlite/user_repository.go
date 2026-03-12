package sqlite

import (
	"database/sql"
	"errors"

	"discordgo-template/internal/domain/user"

	"discordgo-template/pkg/database"

	"github.com/google/uuid"
	"modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"
)

var _ user.UserRepository = (*UserRepositoryImpl)(nil)

type UserRepositoryImpl struct {
	db *database.SqliteDB
}

func NewUserRepository(db *database.SqliteDB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) Create(usr *user.User) error {
	err := r.db.Exec(
		`INSERT INTO users (id)
		 VALUES (?)`,
		usr.ID,
	)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE {
			return user.ErrUserDuplicate
		}
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) Save(usr *user.User) error {
	return r.db.Exec(
		`UPDATE users
		 WHERE id = ?`,
		usr.ID,
	)
}

func (r *UserRepositoryImpl) Delete(usr *user.User) error {
	return r.db.Exec(
		`DELETE FROM users
		 WHERE id = ?`,
		usr.ID,
	)
}

func (r *UserRepositoryImpl) Get(id uuid.UUID) (*user.User, error) {
	row, err := r.db.QueryRow(
		`SELECT id
		 FROM users
		 WHERE id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}

	var usr user.User
	err = row.Scan(&usr.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func (r *UserRepositoryImpl) GetAll() ([]*user.User, error) {
	rows, err := r.db.Query(
		`SELECT id
		 FROM users`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var usr user.User
		err = rows.Scan(&usr.ID)
		if err != nil {
			return nil, err
		}
		users = append(users, &usr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
