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
	"net/url"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type MonitorURI struct {
	data        *types.MonitorURI
	preComments []string
}

func (p *MonitorURI) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) == 2 && parts[0] == "monitor-uri" {
		if _, err := url.Parse(parts[1]); err != nil {
			return "", &errors.ParseError{Parser: "monitor-uri", Line: line, Message: err.Error()}
		}
		p.data = &types.MonitorURI{URI: parts[1]}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "monitor-uri", Line: line}
}

func (p *MonitorURI) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data: fmt.Sprintf("monitor-uri %s", p.data.URI),
		},
	}, nil
}
