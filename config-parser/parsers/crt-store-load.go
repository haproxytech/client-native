/*
Copyright 2024 HAProxy Technologies

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
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type LoadCert struct {
	data        []types.LoadCert
	preComments []string // comments that appear before the actual line
}

func (p *LoadCert) parseError(line string) *errors.ParseError {
	return &errors.ParseError{Parser: "LoadCert", Line: line}
}

func (p *LoadCert) parse(line string, parts []string, comment string) (*types.LoadCert, error) {
	if len(parts) < 3 {
		return nil, p.parseError(line)
	}
	if parts[0] != "load" {
		return nil, p.parseError(line)
	}

	load := new(types.LoadCert)

	for i := 1; i < len(parts); i++ {
		element := parts[i]
		switch element {
		case "crt":
			CheckParsePair(parts, &i, &load.Certificate)
		case "acme":
			CheckParsePair(parts, &i, &load.Acme)
		case "alias":
			CheckParsePair(parts, &i, &load.Alias)
		case "domains":
			CheckParsePair(parts, &i, &load.Domains)
		case "key":
			CheckParsePair(parts, &i, &load.Key)
		case "ocsp":
			CheckParsePair(parts, &i, &load.Ocsp)
		case "issuer":
			CheckParsePair(parts, &i, &load.Issuer)
		case "sctl":
			CheckParsePair(parts, &i, &load.Sctl)
		case "ocsp-update":
			i++
			if i >= len(parts) {
				return nil, p.parseError(line)
			}
			load.OcspUpdate = new(bool)
			if parts[i] == "on" {
				*load.OcspUpdate = true
			} else if parts[i] != "off" {
				return nil, p.parseError(line)
			}
		}
	}
	load.Comment = comment

	// crt is mandatory
	if load.Certificate == "" {
		return nil, p.parseError(line)
	}

	return load, nil
}

func (p *LoadCert) Result() ([]common.ReturnResultLine, error) {
	if len(p.data) == 0 {
		return nil, errors.ErrFetch
	}

	result := make([]common.ReturnResultLine, len(p.data))
	sb := new(strings.Builder)

	for i, load := range p.data {
		sb.Reset()
		sb.WriteString("load")
		CheckWritePair(sb, "crt", load.Certificate)
		CheckWritePair(sb, "acme", load.Acme)
		CheckWritePair(sb, "alias", load.Alias)
		CheckWritePair(sb, "domains", load.Domains)
		CheckWritePair(sb, "key", load.Key)
		CheckWritePair(sb, "ocsp", load.Ocsp)
		CheckWritePair(sb, "issuer", load.Issuer)
		CheckWritePair(sb, "sctl", load.Sctl)
		CheckWritePair(sb, "ocsp-update", fmtOnOff(load.OcspUpdate))

		result[i] = common.ReturnResultLine{
			Data:    sb.String(),
			Comment: load.Comment,
		}
	}

	return result, nil
}

func fmtOnOff(b *bool) string {
	if b == nil {
		return ""
	}
	if *b {
		return "on"
	}
	return "off"
}
