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
	"strings"
)

type ShowDesc struct {
	Desc    string
	Comment string
}

func (s *ShowDesc) Parse(parts []string, comment string) error {
	if len(parts) < 3 {
		return errors.New("not enough params")
	}

	if comment != "" {
		s.Comment = comment
	}
	s.Desc = strings.Join(parts[2:], " ")
	return nil
}

func (s *ShowDesc) String() string {
	if s.Desc != "" {
		return fmt.Sprint("show-desc ", s.Desc)
	}
	return "show-desc"
}

func (s *ShowDesc) GetComment() string {
	return s.Comment
}
