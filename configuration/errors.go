package configuration

import (
	"errors"

	oaerrors "github.com/go-openapi/errors"
)

// CompositeTransactionError helper function to aggregate multiple errors
// when calling multiple operations in transactions.
func CompositeTransactionError(e ...error) *oaerrors.CompositeError {
	return &oaerrors.CompositeError{Errors: append([]error{}, e...)}
}

// ErrVersionMismatch returned when configured and given version do not match
var ErrVersionMismatch = errors.New("Version mismatch")

// ErrBothVersionTransaction returned when you send both transaction and version
var ErrBothVersionTransaction = errors.New("Cannot have both transactionID and version")

// ErrNoVersionTransaction returned when you don't send transaction or version
var ErrNoVersionTransaction = errors.New("Must have either transactionID and version")

// ErrNoParentSpecified returned when parent is required but not given
var ErrNoParentSpecified = errors.New("Parent not specified when parentType is")

// ErrObjectDoesNotExist returned when requested object does not exist in configuration
var ErrObjectDoesNotExist = errors.New("Requested object does not exist")

// ErrTransactionNotFound when given transaction does not exist
var ErrTransactionNotFound = errors.New("Given transaction does not exist")

// ErrCannotReadConfFile when configuration file could not be read
var ErrCannotReadConfFile = errors.New("Cannot read configuration file")

// ErrCannotReadVersion when there is no version comment in configuration
var ErrCannotReadVersion = errors.New("Cannot read version from the configuration file")

// ErrCannotIncrementVersion when version isn't incremented correctly
var ErrCannotIncrementVersion = errors.New("Cannot increment the version in the configuration file")
