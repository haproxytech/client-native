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

package actions

import (
	"fmt"

	"github.com/haproxytech/client-native/v5/config-parser/types"
)

// http-check send-state
type CheckSendState struct {
	Comment string
}

func (c *CheckSendState) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		c.Comment = comment
	}

	if len(parts) > 2 {
		return fmt.Errorf("too many params")
	}

	return nil
}

func (c *CheckSendState) String() string {
	return "send-state"
}

func (c *CheckSendState) GetComment() string {
	return c.Comment
}
