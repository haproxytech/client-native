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

package stats

import (
	"errors"
	"fmt"
)

type Refresh struct {
	Delay   string
	Comment string
}

func (r *Refresh) Parse(parts []string, comment string) error {
	if len(parts) < 3 {
		return errors.New("not enough params")
	}

	if comment != "" {
		r.Comment = comment
	}
	r.Delay = parts[2]
	return nil
}

func (r *Refresh) String() string {
	if r.Delay != "" {
		return fmt.Sprint("refresh ", r.Delay)
	}
	return "refresh"
}

func (r *Refresh) GetComment() string {
	return r.Comment
}
