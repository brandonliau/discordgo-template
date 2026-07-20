package pin

import "errors"

var (
	ErrPinDuplicate = errors.New("pin already exists")
	ErrPinNotFound  = errors.New("pin does not exist")
)
