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

package params

import "fmt"

type ErrParseBindOption interface {
	Error() string
}

type ErrParseServerOption interface {
	Error() string
}

// NotFoundError struct for creating parse errors
type NotFoundError struct {
	Have string
	Want string
}

func (e *NotFoundError) Error() string {
	if e.Want == "" {
		return fmt.Sprintf("error: have [%s] ", e.Have)
	}
	return fmt.Sprintf("error: have [%s] want [%s]", e.Have, e.Want)
}

// NotEnoughParamsError struct for creating parse errors
type NotEnoughParamsError struct{}

func (e *NotEnoughParamsError) Error() string {
	return "error: not enough params"
}

// NotAllowedValuesError struct for allowed values errors
type NotAllowedValuesError struct {
	Have string
	Want []string
}

func (e *NotAllowedValuesError) Error() string {
	return "error: values not allowed"
}
