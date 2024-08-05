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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/parsers/actions"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

// http-check send [meth <method>] [{ uri <uri> | uri-lf <fmt> }>] [ver <version>]
//
//	[hdr <name> <fmt>]* [{ body <string> | body-lf <fmt> }]
//	[comment <msg>]
type CheckSend struct {
	Method        string
	URI           string
	URILogFormat  string
	Version       string
	Header        []CheckSendHeader
	Body          string
	BodyLogFormat string
	CheckComment  string
	Comment       string
}

type CheckSendHeader struct {
	Name   string
	Format string
}

func (c *CheckSend) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		c.Comment = comment
	}

	if len(parts) < 3 {
		return fmt.Errorf("not enough params")
	}

	for i := 2; i < len(parts); i++ {
		el := parts[i]
		switch el {
		case "meth":
			actions.CheckParsePair(parts, &i, &c.Method)
		case "uri":
			actions.CheckParsePair(parts, &i, &c.URI)
		case "uri-lf":
			actions.CheckParsePair(parts, &i, &c.URILogFormat)
		case "ver":
			actions.CheckParsePair(parts, &i, &c.Version)
		// NOTE: HAProxy config supports header values containing spaces,
		// which config-parser normally would support with `\ `.
		// However, because parts is split by spaces and hdr can be provided
		// multiple times with other directives surrounding it, it's
		// impossible to read ahead to put the pieces together.
		// As such, header values with spaces are not supported.
		case "hdr":
			if (i+1) < len(parts) && (i+2) < len(parts) {
				c.Header = append(c.Header, CheckSendHeader{Name: parts[i+1], Format: parts[i+2]})
				i++
			}
		case "body":
			actions.CheckParsePair(parts, &i, &c.Body)
		case "body-lf":
			actions.CheckParsePair(parts, &i, &c.BodyLogFormat)
		case "comment":
			actions.CheckParsePair(parts, &i, &c.CheckComment)
		}
	}

	return nil
}

func (c *CheckSend) String() string {
	sb := &strings.Builder{}

	sb.WriteString("send")

	actions.CheckWritePair(sb, "meth", c.Method)
	actions.CheckWritePair(sb, "uri", c.URI)
	actions.CheckWritePair(sb, "uri-lf", c.URILogFormat)
	actions.CheckWritePair(sb, "ver", c.Version)
	for _, h := range c.Header {
		actions.CheckWritePair(sb, "hdr", h.Name)
		actions.CheckWritePair(sb, "", h.Format)
	}
	actions.CheckWritePair(sb, "body", c.Body)
	actions.CheckWritePair(sb, "body-lf", c.BodyLogFormat)
	actions.CheckWritePair(sb, "comment", c.CheckComment)

	return sb.String()
}

func (c *CheckSend) GetComment() string {
	return c.Comment
}
