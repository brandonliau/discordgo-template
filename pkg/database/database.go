package database

type Database interface {
	InitDB() error
	Close()
	Write(userID string, secret string) error
	Remove(userID string) error
}
