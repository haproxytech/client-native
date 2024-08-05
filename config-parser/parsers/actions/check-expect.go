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
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

// http-check/tcp-check expect
//
//	[min-recv <int>] [comment <msg>]
//	[ok-status <st>] [error-status <st>] [tout-status <st>]
//	[on-success <fmt>] [on-error <fmt>] [status-code <expr>]
//	[!] <match> <pattern>
type CheckExpect struct {
	MinRecv         *int64
	CheckComment    string
	OKStatus        string
	ErrorStatus     string
	TimeoutStatus   string
	OnSuccess       string
	OnError         string
	StatusCode      string
	ExclamationMark bool
	Match           string
	Pattern         string
	Comment         string
}

func (c *CheckExpect) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		c.Comment = comment
	}

	if len(parts) < 3 {
		return fmt.Errorf("not enough params")
	}

	var i int
LoopExpect:
	for i = 2; i < len(parts); i++ {
		el := parts[i]
		switch el {
		case "min-recv":
			if (i + 1) < len(parts) {
				// a *int64 is used as opposed to just int64 since 0 has a special meaning for min-recv
				minRecv, err := strconv.ParseInt(parts[i+1], 10, 64)
				if err != nil {
					return err
				}
				c.MinRecv = &minRecv
				i++
			}
		case "comment":
			CheckParsePair(parts, &i, &c.CheckComment)
		case "ok-status":
			CheckParsePair(parts, &i, &c.OKStatus)
		case "error-status":
			CheckParsePair(parts, &i, &c.ErrorStatus)
		case "tout-status":
			CheckParsePair(parts, &i, &c.TimeoutStatus)
		case "on-success":
			CheckParsePair(parts, &i, &c.OnSuccess)
		case "on-error":
			CheckParsePair(parts, &i, &c.OnError)
		case "status-code":
			CheckParsePair(parts, &i, &c.StatusCode)
		case "!":
			c.ExclamationMark = true
		default:
			break LoopExpect
		}
	}

	// if we broke out of the loop, whatever is leftover should be
	// `<match> <pattern>`. Prevent panics with bounds checks for safety.
	if i >= len(parts) {
		if parserType == types.HTTP {
			return &errors.ParseError{Parser: "HttpCheck", Message: "http-check expect: match not provided"}
		}
		return &errors.ParseError{Parser: "TcpCheck", Message: "tcp-check expect: match not provided"}
	}
	c.Match = parts[i]

	if i+1 >= len(parts) {
		if parserType == types.HTTP {
			return &errors.ParseError{Parser: "HttpCheck", Message: "http-check expect: pattern not provided"}
		}
		return &errors.ParseError{Parser: "TcpCheck", Message: "tcp-check expect: pattern not provided"}
	}
	// Since pattern is always the last option provided, we can safely join
	// the remainder as part of the pattern.
	pattern := strings.Join(parts[i+1:], " ")
	c.Pattern = pattern

	return nil
}

func (c *CheckExpect) String() string {
	sb := &strings.Builder{}

	sb.WriteString("expect")

	if c.MinRecv != nil {
		CheckWritePair(sb, "min-recv", strconv.Itoa(int(*c.MinRecv)))
	}
	CheckWritePair(sb, "comment", c.CheckComment)
	CheckWritePair(sb, "ok-status", c.OKStatus)
	CheckWritePair(sb, "error-status", c.ErrorStatus)
	CheckWritePair(sb, "tout-status", c.TimeoutStatus)
	CheckWritePair(sb, "on-success", c.OnSuccess)
	CheckWritePair(sb, "on-error", c.OnError)
	CheckWritePair(sb, "status-code", c.StatusCode)

	if c.ExclamationMark {
		CheckWritePair(sb, "", "!")
	}
	CheckWritePair(sb, "", c.Match)
	CheckWritePair(sb, "", c.Pattern)

	return sb.String()
}

func (c *CheckExpect) GetComment() string {
	return c.Comment
}
