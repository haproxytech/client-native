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

package filters

import (
	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type Filters struct {
	Name        string
	data        []types.Filter
	preComments []string // comments that appear before the actual line
}

func (h *Filters) Init() {
	h.data = []types.Filter{}
	h.Name = "filter"
}

func (h *Filters) ParseFilter(filter types.Filter, parts []string, comment string) error {
	if err := filter.Parse(parts, ""); err != nil {
		return &errors.ParseError{Parser: "FilterLines", Line: ""}
	}
	h.data = append(h.data, filter)
	return nil
}

func (h *Filters) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) >= 2 && parts[0] == "filter" {
		var err error
		switch parts[1] {
		case "trace":
			err = h.ParseFilter(&Trace{}, parts, comment)
		case "compression":
			err = h.ParseFilter(&Compression{}, parts, comment)
		case "cache":
			err = h.ParseFilter(&Cache{}, parts, comment)
		case "spoe":
			err = h.ParseFilter(&Spoe{}, parts, comment)
		case "fcgi-app":
			err = h.ParseFilter(&FcgiApp{}, parts, comment)
		case "opentracing":
			err = h.ParseFilter(&Opentracing{}, parts, comment)
		case "bwlim-in", "bwlim-out":
			err = h.ParseFilter(&BandwidthLimit{}, parts, comment)
		default:
			return "", &errors.ParseError{Parser: "FilterLines", Line: line}
		}
		if err != nil {
			return "", err
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "FilterLines", Line: line}
}

func (h *Filters) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, req := range h.data {
		result[index] = req.Result()
	}
	return result, nil
}
