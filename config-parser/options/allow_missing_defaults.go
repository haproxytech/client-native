/*
Copyright 2026 HAProxy Technologies

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
	"slices"
	"strings"
)

// allowMissingDefaults is ParserOption that permits specified default sections to be missing in the configuration.
// the use case is to be able to write fragments of configuration that are not complete but still valid when assembling with the rest of the configuration.
type allowMissingDefaults struct {
	names []string
}

func (o allowMissingDefaults) Set(p *Parser) error {
	p.AllowMissingDefaults = o.names
	return nil
}

// AllowMissingDefaults returns a ParserOption that permits specified default sections to be missing in the configuration.
func AllowMissingDefaults(sectionNames []string) ParserOption {
	names := []string{}
	for _, name := range sectionNames {
		satanized := strings.Trim(name, " ")
		if satanized != "" && !slices.Contains(names, satanized) {
			names = append(names, satanized)
		}
	}

	return allowMissingDefaults{
		names: names,
	}
}
