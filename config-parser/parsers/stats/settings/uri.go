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
	"fmt"
)

type URI struct {
	Prefix  string
	Comment string
}

func (u *URI) Parse(parts []string, comment string) error {
	if len(parts) < 3 {
		return fmt.Errorf("not enough params")
	}

	if comment != "" {
		u.Comment = comment
	}
	u.Prefix = parts[2]
	return nil
}

func (u *URI) String() string {
	if u.Prefix != "" {
		return fmt.Sprint("uri ", u.Prefix)
	}
	return "uri"
}

func (u *URI) GetComment() string {
	return u.Comment
}
