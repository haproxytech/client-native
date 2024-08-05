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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
)

type HTTPRequest struct {
	Type     string
	Cond     string
	CondTest string
	Comment  string
}

func (h *HTTPRequest) Parse(parts []string, comment string) error {
	if len(parts) < 3 {
		return fmt.Errorf("not enough params")
	}

	if comment != "" {
		h.Comment = comment
	}

	switch parts[2] {
	case "allow", "deny":
		command, condition := common.SplitRequest(parts[2:])
		if len(command) != 1 {
			return fmt.Errorf("error parsing http-request")
		}
		h.Type = command[0]
		if len(condition) > 1 {
			h.Cond = condition[0]
			h.CondTest = strings.Join(condition[1:], " ")
		}
		return nil
	case "auth":
		return h.parseAuth(parts)
	default:
		return fmt.Errorf("error parsing http-request")
	}
}

func (h *HTTPRequest) parseAuth(parts []string) error {
	command, condition := common.SplitRequest(parts[2:])
	switch len(command) {
	case 1:
		h.Type = strings.Join(command, " ")
	case 3:
		if command[1] != "realm" {
			return fmt.Errorf("error parsing http-request")
		}
		h.Type = strings.Join(command, " ")
	default:
		return fmt.Errorf("error parsing http-request")
	}
	if len(condition) > 1 {
		h.Cond = condition[0]
		h.CondTest = strings.Join(condition[1:], " ")
	}
	return nil
}

func (h *HTTPRequest) String() string {
	condition := ""
	if h.Cond != "" {
		condition = fmt.Sprintf(" %s %s", h.Cond, h.CondTest)
	}
	return fmt.Sprintf("http-request %s%s", h.Type, condition)
}

func (h *HTTPRequest) GetComment() string {
	return h.Comment
}
