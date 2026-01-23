package errors

//nolint:revive
import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrGeneral       = errors.New("general error")
)
