/*
Copyright 2021 HAProxy Technologies

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

package options

import (
	"io"

	"github.com/haproxytech/go-logger"
)

type Parser struct {
	Path                       string
	Reader                     io.Reader
	Logger                     logger.Format // we always will have p.Options.LogPrefix
	UseV2HTTPCheck             bool
	UseMd5Hash                 bool
	UseListenSectionParsers    bool
	DisableUnProcessed         bool
	Log                        bool
	LogPrefix                  string
	NoNamedDefaultsFrom        bool
	DefaultSectionsSkipOnWrite map[string]struct{}

	// AllowMissingDefaults specifies default sections allowed to be missing in the parser's configuration.
	// note: it is an array instead of a map because we don't expect to have many elements (in this case an array is more performant than a map)
	AllowMissingDefaults []string
}

type ParserOption interface {
	Set(p *Parser) error
}
