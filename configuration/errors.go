package configuration

import (
	"errors"

	oaerrors "github.com/go-openapi/errors"
)

func CompositeTransactionError(e ...error) *oaerrors.CompositeError {
	return &oaerrors.CompositeError{Errors: append([]error{}, e...)}
}

var ErrVersionMismatch error = errors.New("Version mismatch")
var ErrBothVersionTransaction error = errors.New("Cannot have both transactionID and version")
var ErrNoVersionTransaction error = errors.New("Must have either transactionID and version")
var ErrNoParentSpecified error = errors.New("Parent not specified when parentType is")
var ErrObjectDoesNotExist error = errors.New("Requested object does not exist")
var ErrTransactionNotFound error = errors.New("Given transaction does not exist")
var ErrNotExists error = errors.New("Object does not exist")
var ErrCannotReadConfFile error = errors.New("Cannot read configuration file")
var ErrCannotReadVersion error = errors.New("Cannot read version from the configuration file")
var ErrCannotIncrementVersion error = errors.New("Cannot increment the version in the configuration file")
