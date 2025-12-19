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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type ACL struct {
	data        []models.ACL
	preComments []string // comments that appear before the actual line
}

func (h *ACL) parse(line string, parts []string, comment string) (*models.ACL, error) {
	if len(parts) >= 3 {
		data := &models.ACL{
			ACLName:   parts[1],
			Criterion: parts[2],
			Value:     strings.Join(parts[3:], " "),
		}
		if comment != "" {
			data.Metadata = misc.ParseMetadata(comment)
		}
		return data, nil
	}
	return nil, &errors.ParseError{Parser: "ACLLines", Line: line}
}

func (h *ACL) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, req := range h.data {
		var sb strings.Builder
		sb.WriteString("acl ")
		sb.WriteString(req.ACLName)
		sb.WriteString(" ")
		sb.WriteString(req.Criterion)
		sb.WriteString(" ")
		sb.WriteString(req.Value)

		comment, err := misc.SerializeMetadata(req.Metadata)
		if err != nil {
			comment = ""
		}
		result[index] = common.ReturnResultLine{
			Data:    sb.String(),
			Comment: comment,
		}
	}
	return result, nil
}
