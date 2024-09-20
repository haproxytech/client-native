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
	stderrors "errors"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/parsers/actions"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

// tcp-check send-binary-lf <hexfmt> [comment <msg>]
type CheckSendBinaryLf struct {
	HexFmt       string
	CheckComment string
	Comment      string
}

func (c *CheckSendBinaryLf) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		c.Comment = comment
	}
	if len(parts) < 3 {
		return stderrors.New("not enough params")
	}
	c.HexFmt = parts[2]
	for i := 3; i < len(parts); i++ {
		el := parts[i]
		if el == "comment" {
			actions.CheckParsePair(parts, &i, &c.CheckComment)
		}
	}
	return nil
}

func (c *CheckSendBinaryLf) String() string {
	sb := &strings.Builder{}
	sb.WriteString("send-binary-lf")
	sb.WriteString(" ")
	sb.WriteString(c.HexFmt)
	actions.CheckWritePair(sb, "comment", c.CheckComment)
	return sb.String()
}

func (c *CheckSendBinaryLf) GetComment() string {
	return c.Comment
}
