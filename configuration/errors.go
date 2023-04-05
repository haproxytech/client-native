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
	"errors"
	"fmt"

	oaerrors "github.com/go-openapi/errors"
)

var (
	// General error, unknown cause
	ErrGeneralError = errors.New("unknown error")

	// Errors regarding configurations
	ErrNoParentSpecified      = errors.New("unspecified parent")
	ErrParentDoesNotExist     = errors.New("missing parent")
	ErrBothVersionTransaction = errors.New("both version and transaction specified, specify only one")
	ErrNoVersionTransaction   = errors.New("version or transaction not specified")
	ErrValidationError        = errors.New("validation error")
	ErrVersionMismatch        = errors.New("version mismatch")

	ErrTransactionDoesNotExist  = errors.New("transaction does not exist")
	ErrTransactionAlreadyExists = errors.New("transaction already exist")
	ErrCannotParseTransaction   = errors.New("failed to parse transaction")

	ErrObjectDoesNotExist    = errors.New("missing object")
	ErrObjectAlreadyExists   = errors.New("object already exists")
	ErrObjectIndexOutOfRange = errors.New("index out of range")

	ErrErrorChangingConfig = errors.New("invalid change")
	ErrCannotReadConfFile  = errors.New("failed to read configuration")
	ErrCannotReadVersion   = errors.New("cannot read version")
	ErrCannotSetVersion    = errors.New("cannot set version")

	ErrCannotFindHAProxy = errors.New("failed to find HAProxy")

	ErrClientDoesNotExists = errors.New("client does not exist")
)

// ConfError general configuration client error
type ConfError struct {
	err    error
	reason string // optional
}

func (e *ConfError) Err() error {
	return e.err
}

// Error implementation for ConfError
func (e *ConfError) Error() string {
	if e.reason == "" {
		return e.err.Error()
	}
	return fmt.Sprintf("%s: %s", e.err.Error(), e.reason)
}

func (e *ConfError) Is(target error) bool {
	return e.err == target
}

// NewConfError constructor for ConfError
func NewConfError(err error, reason string) *ConfError {
	return &ConfError{err: err, reason: reason}
}

// CompositeTransactionError helper function to aggregate multiple errors
// when calling multiple operations in transactions.
func CompositeTransactionError(e ...error) *oaerrors.CompositeError {
	return &oaerrors.CompositeError{Errors: append([]error{}, e...)}
}
