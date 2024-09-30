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

type ShowNode struct {
	Name    string
	Comment string
}

func (s *ShowNode) Parse(parts []string, comment string) error {
	if len(parts) < 2 {
		return fmt.Errorf("not enough params")
	}

	if comment != "" {
		s.Comment = comment
	}
	if len(parts) > 2 {
		s.Name = parts[2]
	}
	return nil
}

func (s *ShowNode) String() string {
	if s.Name != "" {
		return fmt.Sprint("show-node ", s.Name)
	}
	return "show-node"
}

func (s *ShowNode) GetComment() string {
	return s.Comment
}
