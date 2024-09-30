/*
Copyright 2019 HAProxy Technologies

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package errors

import (
	"errors"
	"fmt"
)

// ParseError struct for creating parse errors
type ParseError struct {
	Parser  string
	Line    string
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s:%s: Parse error on %s", e.Parser, e.Message, e.Line)
}

var ErrAttributeNotFound = errors.New("attribute not found")

var ErrFetch = errors.New("no data")

var ErrIndexOutOfRange = errors.New("index out of range")

var ErrInvalidData = errors.New("invalid data")

var ErrParserMissing = errors.New("parser missing")

var ErrSectionAlreadyExists = errors.New("section already exists")

var ErrSectionMissing = errors.New("section missing")

var ErrSectionTypeMissing = errors.New("section type missing")

var ErrSectionTypeNotAllowed = errors.New("section type not allowed")

var ErrScopeMissing = errors.New("scope missing")

var ErrScopeAlreadyExists = errors.New("scope already exists")

var ErrFromDefaultsSectionMissing = errors.New("defaults section specified in from does not exist")

var ErrCircularDependency = errors.New("circular dependency detected")
