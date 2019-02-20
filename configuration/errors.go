package configuration

import (
	"fmt"

	oaerrors "github.com/go-openapi/errors"
)

const (
	// General error, unknown cause
	ErrGeneralError = 0

	// Errors regarding configurations
	ErrSyntaxWrong             = 1
	ErrNoParentSpecified       = 10
	ErrTransactionDoesNotExist = 20
	ErrObjectAlreadyExists     = 23
	ErrObjectDoesNotExist      = 22
	ErrObjectIndexOutOfRange   = 26
	ErrErrorChangingConfig     = 25
	ErrCannotReadConfFile      = 30
	ErrCannotReadVersion       = 31
	ErrCannotSetVersion        = 32
	ErrVersionMismatch         = 33
	ErrBothVersionTransaction  = 34
	ErrNoVersionTransaction    = 35
	ErrCannotParseTransaction  = 40
	ErrValidationError         = 50
)

// ConfError general configuration client error
type ConfError struct {
	code int
	msg  string
}

// LBCTLError error when executing lbctl, embeds ConfError
type LBCTLError struct {
	ConfError
	action string
	cmd    string
}

// Error implementation for ConfError
func (e *ConfError) Error() string {
	return fmt.Sprintf("%v: %s", e.code, e.msg)
}

// Code returns ConfError code, which is one of constants in this package
func (e *ConfError) Code() int {
	return e.code
}

// Error implementation for LBCTLError
func (e *LBCTLError) Error() string {
	return fmt.Sprintf("%v: Error executing LBCTL: %s, %s. Output: %s", e.code, e.cmd, e.action, e.msg)
}

// NewLBCTLError contstructor for LBCTLError
func NewLBCTLError(code int, cmd, action, msg string) *LBCTLError {
	return &LBCTLError{ConfError: ConfError{code: code, msg: msg}, action: action, cmd: cmd}
}

// NewConfError contstructor for ConfError
func NewConfError(code int, msg string) *ConfError {
	return &ConfError{code: code, msg: msg}
}

// CompositeTransactionError helper function to aggregate multiple errors
// when calling multiple operations in transactions.
func CompositeTransactionError(e ...error) *oaerrors.CompositeError {
	return &oaerrors.CompositeError{Errors: append([]error{}, e...)}
}
