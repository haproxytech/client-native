// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package configuration

import (
	"fmt"

	parser "github.com/haproxytech/client-native/v6/config-parser"
)

type Parser interface {
	Parser() parser.Parser
	HasParser(transactionID string) bool
	GetParser(transactionID string) (parser.Parser, error)
	AddParser(transactionID string) error
	DeleteParser(transactionID string) error
	CommitParser(transactionID string) error
	GetVersion(transactionID string) (int64, error)
	IncrementVersion() error
	LoadData(filename string) error
	Save(transactionFile, transactionID string) error
}

func getParserFromParent(attribute, parentType, parentName string) (parser.Section, string, error) {
	switch attribute {
	case "http-request", "http-response", "http-after-response", "tcp-request", "quic-initial":
		switch parentType {
		case BackendParentName:
			return parser.Backends, parentName, nil
		case FrontendParentName:
			return parser.Frontends, parentName, nil
		case DefaultsParentName:
			if parentName == "" {
				parentName = parser.DefaultSectionName
			}
			return parser.Defaults, parentName, nil

		default:
			return "", "", fmt.Errorf("unsupported parent: %s", parentType)
		}
	case "acl":
		switch parentType {
		case BackendParentName:
			return parser.Backends, parentName, nil
		case FrontendParentName:
			return parser.Frontends, parentName, nil
		case FCGIAppParentName:
			return parser.FCGIApp, parentName, nil
		case DefaultsParentName:
			if parentName == "" {
				parentName = parser.DefaultSectionName
			}
			return parser.Defaults, parentName, nil

		default:
			return "", "", fmt.Errorf("unsupported parent: %s", parentType)
		}
	case "tcp-response":
		switch parentType {
		case BackendParentName:
			return parser.Backends, parentName, nil
		case DefaultsParentName:
			if parentName == "" {
				parentName = parser.DefaultSectionName
			}
			return parser.Defaults, parentName, nil
		default:
			return "", "", fmt.Errorf("unsupported parent: %s", parentType)
		}
	}

	return "", "", fmt.Errorf("unsupported attribute: %s", attribute)
}
