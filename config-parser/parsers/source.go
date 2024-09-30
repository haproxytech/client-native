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
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type Source struct {
	data        *types.Source
	preComments []string // comments that appear before the actual line
}

func (s *Source) Parse(line string, parts []string, comment string) (string, error) { //nolint: gocognit
	if len(parts) < 2 {
		return "", &errors.ParseError{Parser: "Source", Line: line}
	}
	if parts[0] != "source" {
		return "", &errors.ParseError{Parser: "Source", Line: line}
	}
	s.data = &types.Source{
		Comment: comment,
	}
	if strings.Contains(parts[1], ":") {
		addressAndPort := strings.Split(parts[1], ":")
		s.data.Address = addressAndPort[0]
		if port, err := strconv.ParseInt(addressAndPort[1], 10, 64); err == nil {
			s.data.Port = port
		}
	} else {
		s.data.Address = parts[1]
	}
	for i := 2; i < len(parts); i++ {
		element := parts[i]
		switch element {
		case "usesrc":
			s.data.UseSrc = true
			i++
			if i >= len(parts) {
				s.data = nil
				return "", &errors.ParseError{Parser: "Source", Line: line}
			}
			switch {
			case strings.HasPrefix(parts[i], "clientip"):
				s.data.ClientIP = true
			case strings.HasPrefix(parts[i], "client"):
				s.data.Client = true
			case strings.HasPrefix(parts[i], "hdr_ip"):
				s.data.HdrIP = true
				param := strings.TrimPrefix(parts[i], "hdr_ip(")
				param = strings.TrimRight(param, ")")
				if strings.Contains(param, ",") {
					HdrAndOcc := strings.Split(param, ",")
					s.data.Hdr = HdrAndOcc[0]
					s.data.Occ = HdrAndOcc[1]
				} else {
					s.data.Hdr = param
				}
			default:
				if strings.Contains(parts[i], ":") {
					addressAndPort := strings.Split(parts[i], ":")
					s.data.AddressSecond = addressAndPort[0]
					if port, err := strconv.ParseInt(addressAndPort[1], 10, 64); err == nil {
						s.data.PortSecond = port
					}
				} else {
					s.data.AddressSecond = parts[i]
				}
			}
		case "interface":
			i++
			if i >= len(parts) {
				s.data = nil
				return "", &errors.ParseError{Parser: "Source", Line: line}
			}
			s.data.Interface = parts[i]
		}
	}
	return "", nil
}

func (s *Source) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	var sb strings.Builder
	sb.WriteString("source")
	sb.WriteString(" ")
	sb.WriteString(s.data.Address)
	if s.data.Port > 0 {
		sb.WriteString(":")
		sb.WriteString(strconv.FormatInt(s.data.Port, 10))
	}
	if s.data.UseSrc {
		sb.WriteString(" ")
		sb.WriteString("usesrc")
		sb.WriteString(" ")
		if s.data.Client {
			sb.WriteString("client")
		}
		if s.data.ClientIP {
			sb.WriteString("clientip")
		}
		if s.data.HdrIP {
			sb.WriteString("hdr_ip(")
			sb.WriteString(s.data.Hdr)
			if s.data.Occ != "" {
				sb.WriteString(",")
				sb.WriteString(s.data.Occ)
			}
			sb.WriteString(")")
		}
		if s.data.AddressSecond != "" {
			sb.WriteString(s.data.AddressSecond)
			if s.data.PortSecond > 0 {
				sb.WriteString(":")
				sb.WriteString(strconv.FormatInt(s.data.PortSecond, 10))
			}
		}
	}
	if s.data.Interface != "" {
		sb.WriteString(" ")
		sb.WriteString("interface")
		sb.WriteString(" ")
		sb.WriteString(s.data.Interface)
	}
	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: s.data.Comment,
		},
	}, nil
}
