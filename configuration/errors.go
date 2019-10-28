// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package configuration

import (
	"fmt"

	oaerrors "github.com/go-openapi/errors"
)

const (
	// General error, unknown cause
	ErrGeneralError = 0

	// Errors regarding configurations
	ErrNoParentSpecified      = 10
	ErrParentDoesNotExist     = 11
	ErrBothVersionTransaction = 12
	ErrNoVersionTransaction   = 13
	ErrValidationError        = 14
	ErrVersionMismatch        = 15

	ErrTransactionDoesNotExist  = 20
	ErrTransactionAlreadyExists = 21
	ErrCannotParseTransaction   = 22

	ErrObjectDoesNotExist    = 30
	ErrObjectAlreadyExists   = 31
	ErrObjectIndexOutOfRange = 32

	ErrErrorChangingConfig = 40
	ErrCannotReadConfFile  = 41
	ErrCannotReadVersion   = 42
	ErrCannotSetVersion    = 43

	ErrCannotFindHAProxy = 50
)

// ConfError general configuration client error
type ConfError struct {
	code int
	msg  string
}

// Error implementation for ConfError
func (e *ConfError) Error() string {
	return fmt.Sprintf("%v: %s", e.code, e.msg)
}

// Code returns ConfError code, which is one of constants in this package
func (e *ConfError) Code() int {
	return e.code
}

// NewConfError constructor for ConfError
func NewConfError(code int, msg string) *ConfError {
	return &ConfError{code: code, msg: msg}
}

// CompositeTransactionError helper function to aggregate multiple errors
// when calling multiple operations in transactions.
func CompositeTransactionError(e ...error) *oaerrors.CompositeError {
	return &oaerrors.CompositeError{Errors: append([]error{}, e...)}
}
