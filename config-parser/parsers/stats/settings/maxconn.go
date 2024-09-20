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
	"errors"

	parsers "github.com/haproxytech/client-native/v6/config-parser/parsers"
)

type MaxConn struct {
	Maxconn *parsers.MaxConn
}

func (m *MaxConn) Parse(parts []string, comment string) error {
	if len(parts) < 3 {
		return errors.New("not enough params")
	}

	m.Maxconn = &parsers.MaxConn{}
	_, err := m.Maxconn.Parse("", parts[1:], comment)
	if err != nil {
		return errors.New("error parsing maxconn")
	}
	return nil
}

func (m *MaxConn) String() string {
	if res, _ := m.Maxconn.Result(); len(res) != 0 {
		return res[0].Data
	}
	return "maxconn"
}

func (m *MaxConn) GetComment() string {
	if res, _ := m.Maxconn.Result(); len(res) != 0 {
		return res[0].Comment
	}
	return ""
}
