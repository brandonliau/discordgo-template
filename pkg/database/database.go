package database

type Database interface {
	Close()
	Migrate() error
	Write(userID string, secret string) error
	Remove(userID string) error
}
