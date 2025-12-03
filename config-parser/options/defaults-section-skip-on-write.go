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

type defaultSectionsSkipOnWrite struct {
	sectionNames map[string]struct{}
}

func (o defaultSectionsSkipOnWrite) Set(p *Parser) error {
	p.DefaultSectionsSkipOnWrite = o.sectionNames
	return nil
}

// sectionNames is a list of defaults section that the parser will not serialize even if existing in the configuration
func DefaultSectionsSkipOnWrite(sectionNames []string) ParserOption {
	names := make(map[string]struct{}, len(sectionNames))
	for _, s := range sectionNames {
		names[s] = struct{}{}
	}
	return defaultSectionsSkipOnWrite{
		sectionNames: names,
	}
}
