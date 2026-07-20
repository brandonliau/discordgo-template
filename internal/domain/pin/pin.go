package pin

type Pin struct {
	UserID string
	Zip    string
}

func New(userID string, zip string) Pin {
	return Pin{
		UserID: userID,
		Zip:    zip,
	}
}
