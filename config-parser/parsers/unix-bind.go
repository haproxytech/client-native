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

package parsers

import (
	"io"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type UnixBind struct {
	data        *types.UnixBind
	preComments []string // comments that appear before the actual line
}

func (p *UnixBind) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) < 3 {
		return "", &errors.ParseError{Parser: "UnixBind", Line: line}
	}
	if parts[0] != "unix-bind" {
		return "", &errors.ParseError{Parser: "UnixBind", Line: line}
	}

	data := &types.UnixBind{}

	for i := 1; i < len(parts); i++ {
		element := parts[i]
		switch element {
		case "prefix":
			CheckParsePair(parts, &i, &data.Prefix)
		case "mode":
			CheckParsePair(parts, &i, &data.Mode)
		case "user":
			CheckParsePair(parts, &i, &data.User)
		case "uid":
			CheckParsePair(parts, &i, &data.UID)
		case "group":
			CheckParsePair(parts, &i, &data.Group)
		case "gid":
			CheckParsePair(parts, &i, &data.GID)
		}
	}
	p.data = data
	return "", nil
}

func (p *UnixBind) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}

	sb := &strings.Builder{}

	sb.WriteString("unix-bind")

	CheckWritePair(sb, "prefix", p.data.Prefix)
	CheckWritePair(sb, "mode", p.data.Mode)
	CheckWritePair(sb, "user", p.data.User)
	CheckWritePair(sb, "uid", p.data.UID)
	CheckWritePair(sb, "group", p.data.Group)
	CheckWritePair(sb, "gid", p.data.GID)

	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: p.data.Comment,
		},
	}, nil
}

func CheckParsePair(parts []string, i *int, str *string) {
	if (*i + 1) < len(parts) {
		*str = parts[*i+1]
		*i++
	}
}

func CheckWritePair(sb io.StringWriter, key string, value string) {
	if value != "" {
		_, _ = sb.WriteString(" ")
		if key != "" {
			_, _ = sb.WriteString(key)
			_, _ = sb.WriteString(" ")
		}
		_, _ = sb.WriteString(value)
	}
}
