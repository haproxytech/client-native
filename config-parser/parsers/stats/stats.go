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

package stats

import (
	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	stats "github.com/haproxytech/client-native/v6/config-parser/parsers/stats/settings"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type Stats struct {
	Name        string
	Mode        string
	data        []types.StatsSettings
	preComments []string // comments that appear before the actual line
}

func (s *Stats) Init() {
	s.Name = "stats"
	s.data = []types.StatsSettings{}
}

func (s *Stats) ParseStats(stats types.StatsSettings, parts []string, comment string) error {
	err := stats.Parse(parts, comment)
	if err != nil {
		return &errors.ParseError{Parser: "Stats", Line: "", Message: err.Error()}
	}
	s.data = append(s.data, stats)
	return nil
}

func (s *Stats) Parse(line string, parts []string, comment string) (string, error) {
	var err error
	if parts[0] != "stats" || len(parts) < 2 {
		return "", &errors.ParseError{Parser: "Stats", Line: line}
	}

	switch parts[1] {
	case "admin":
		if s.Mode == "defaults" {
			return "", &errors.ParseError{Parser: "Stats", Line: line}
		}
		err = s.ParseStats(&stats.Admin{}, parts, comment)
	case "auth":
		err = s.ParseStats(&stats.Auth{}, parts, comment)
	case "bind-process":
		err = s.ParseStats(&stats.BindProcess{}, parts, comment)
	case "enable", "hide-version", "show-legends", "show-modules":
		err = s.ParseStats(&stats.OneWord{}, parts, comment)
	case "maxconn":
		err = s.ParseStats(&stats.MaxConn{}, parts, comment)
	case "realm":
		err = s.ParseStats(&stats.Realm{}, parts, comment)
	case "refresh":
		err = s.ParseStats(&stats.Refresh{}, parts, comment)
	case "scope":
		err = s.ParseStats(&stats.Scope{}, parts, comment)
	case "show-desc":
		err = s.ParseStats(&stats.ShowDesc{}, parts, comment)
	case "show-node":
		err = s.ParseStats(&stats.ShowNode{}, parts, comment)
	case "uri":
		err = s.ParseStats(&stats.URI{}, parts, comment)

	case "http-request":
		if s.Mode == "defaults" || s.Mode == "frontend" {
			return "", &errors.ParseError{Parser: "Stats", Line: line}
		}
		err = s.ParseStats(&stats.HTTPRequest{}, parts, comment)
	default:
		return "", &errors.ParseError{Parser: "Stats", Line: line}
	}
	if err != nil {
		return "", err
	}
	return "", nil
}

func (s *Stats) Result() ([]common.ReturnResultLine, error) {
	if len(s.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(s.data))
	for index, stats := range s.data {
		result[index] = common.ReturnResultLine{
			Data:    "stats " + stats.String(),
			Comment: stats.GetComment(),
		}
	}
	return result, nil
}
