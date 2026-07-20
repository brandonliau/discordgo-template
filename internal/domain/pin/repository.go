package pin

type Repository interface {
	Create(pin Pin) error
	Delete(userID, zip string) error
	ListByUser(userID string) ([]*Pin, error)
}
