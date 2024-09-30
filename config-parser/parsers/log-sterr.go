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
	"fmt"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type LogStdErr struct {
	preComments []string
	data        []types.LogStdErr
}

func (p *LogStdErr) parse(parts []string) (*types.LogStdErr, error) {
	protectedKeywords := map[string]string{
		"sample": "",
		"format": "",
		"len":    "",
	}

	getValue := func(part []string, index *int) (string, error) {
		var value string
		if *index+1 == len(parts) {
			return "", fmt.Errorf("missing attribute value")
		}

		value = part[*index+1]

		_, ok := protectedKeywords[value]
		if ok {
			return "", fmt.Errorf("missing value for attribute")
		}

		defer func() {
			*index += 2
		}()

		return parts[*index+1], nil
	}

	data := &types.LogStdErr{}

	for index := 0; index < len(parts); {
		part := parts[index]

		switch {
		case part == "sample":
			v, err := getValue(parts, &index)
			if err != nil {
				return nil, err
			}

			sampleData := strings.Split(v, ":")

			if len(sampleData) != 2 || sampleData[0] == "" {
				return nil, fmt.Errorf("sample size is malformed")
			}

			data.SampleRange = sampleData[0]

			if data.SampleSize, err = strconv.ParseInt(sampleData[1], 10, 64); err != nil {
				return nil, fmt.Errorf("expected integer value")
			}
		case part == "format":
			v, err := getValue(parts, &index)
			if err != nil {
				return nil, err
			}

			data.Format = v
		case part == "len":
			v, err := getValue(parts, &index)
			if err != nil {
				return nil, err
			}

			if data.Length, err = strconv.ParseInt(v, 10, 64); err != nil {
				return nil, fmt.Errorf("expected integer value")
			}
		case index == len(parts)-3: // <facility> <level> <minlevel>
			data.Facility = part
			data.Level = parts[index+1]
			data.MinLevel = parts[index+2]

			return data, nil
		case index == len(parts)-2: // <facility> <level>
			data.Facility = part
			data.Level = parts[index+1]

			return data, nil
		case index == len(parts)-1: // <facility>
			data.Facility = part

			return data, nil
		default:
			index++
		}
	}

	return data, nil
}

func (p *LogStdErr) Parse(line string, parts []string, comment string) (string, error) {
	var err error
	if len(parts) < 2 {
		return "", &errors.ParseError{Parser: "LogStdErr", Line: line, Message: "Missing required attributes"}
	}

	if parts[1] == "global" {
		if len(parts) > 2 {
			return "", &errors.ParseError{Parser: "LogStdErr", Line: line, Message: "global attribute is exclusive"}
		}

		p.data = append(p.data, types.LogStdErr{
			Global: true,
		})

		return "", nil
	}

	var data *types.LogStdErr

	if data, err = p.parse(parts[2:]); err != nil {
		return "", &errors.ParseError{Parser: "LogStdErr", Line: line, Message: err.Error()}
	}

	data.Address = parts[1]

	if !data.Global && len(data.Facility) == 0 {
		return "", &errors.ParseError{Parser: "LogStdErr", Line: line, Message: "Missing required facility"}
	}

	p.data = append(p.data, *data)

	return "", nil
}

func (p *LogStdErr) Result() ([]common.ReturnResultLine, error) {
	if len(p.data) == 0 {
		return nil, errors.ErrFetch
	}

	lines := make([]common.ReturnResultLine, 0, len(p.data))

	for _, data := range p.data {
		var sb strings.Builder

		sb.WriteString("log-stderr ")

		if data.Global {
			sb.WriteString("global")

			lines = append(lines, common.ReturnResultLine{
				Data:    sb.String(),
				Comment: data.Comment,
			})

			continue
		}

		sb.WriteString(data.Address)

		if data.Length > 0 {
			sb.WriteString(" len ")
			sb.WriteString(strconv.FormatInt(data.Length, 10))
		}

		if len(data.Format) > 0 {
			sb.WriteString(" format ")
			sb.WriteString(data.Format)
		}

		if data.SampleSize > 0 && len(data.SampleRange) > 0 {
			sb.WriteString(" sample ")
			sb.WriteString(data.SampleRange)
			sb.WriteString(":")
			sb.WriteString(strconv.FormatInt(data.SampleSize, 10))
		}

		if len(data.Facility) > 0 {
			sb.WriteString(" ")
			sb.WriteString(data.Facility)
		}

		if len(data.Level) > 0 {
			sb.WriteString(" ")
			sb.WriteString(data.Level)
		}

		if len(data.MinLevel) > 0 {
			sb.WriteString(" ")
			sb.WriteString(data.MinLevel)
		}

		if comment := data.Comment; len(comment) > 0 {
			sb.WriteString(" # ")
			sb.WriteString(comment)
		}

		lines = append(lines, common.ReturnResultLine{
			Data:    sb.String(),
			Comment: data.Comment,
		})
	}

	return lines, nil
}
