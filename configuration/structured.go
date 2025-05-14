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
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/configuration/options"
)

type Structured interface {
	StructuredGlobal
	StructuredDefaults
	StructuredFrontend
	StructuredBackend
	StructuredPeerSection
	StructuredResolver
	StructuredUserlist
	StructuredRing
	StructuredLogProfile
	StructuredLogForward
	StructuredFCGIApp
	StructuredMailersSection
	StructuredTraces
	StructuredAcmeProvider
}

type StructuredToParserArgs struct {
	TID         string
	Parser      *parser.Parser
	Options     *options.ConfigurationOptions
	HandleError func(id, parentType, parentName, transactionID string, implicit bool, err error) error
}
