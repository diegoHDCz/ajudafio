package shared

import "errors"

type EntityNotFoundError struct {
	Entity string
}

func (e *EntityNotFoundError) Error() string {
	return e.Entity + " not found"
}

func ErrNotFound(entity string) error {
	return &EntityNotFoundError{Entity: entity}
}

func IsNotFound(err error) bool {
	var target *EntityNotFoundError
	return errors.As(err, &target)
}

var ErrUserAlreadyExists = errors.New("user already exists")
