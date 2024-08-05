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

package sorter

import (
	"sort"

	parserErrors "github.com/haproxytech/client-native/v6/config-parser/errors"
)

type Section struct {
	Name string
	From string
}

// Sort sections based on rules:
//   - all dependencies must exist
//   - there must be no circular dependencies
//   - sort by Name
//   - if there is a dependency, section is moved after the one it depends on
func Sort(sections []Section) error {
	// first check that all dependencies exist
	for _, section := range sections {
		if section.From == "" {
			continue
		}
		depOK := false
		for _, cmp := range sections {
			if section.From == cmp.Name {
				depOK = true
				break
			}
		}
		if !depOK {
			return parserErrors.ErrFromDefaultsSectionMissing
		}
	}
	// first check for circular dependencies
	for _, section := range sections {
		if circularDependency(map[string]struct{}{section.Name: {}}, section.From, sections) {
			return parserErrors.ErrCircularDependency
		}
	}

	// first sort by name
	sort.SliceStable(sections, func(i, j int) bool {
		return sections[i].Name < sections[j].Name
	})
	// then go through list and check for circular dependencies
	sortByFrom(0, sections)
	// done
	return nil
}

func circularDependency(visited map[string]struct{}, current string, sections []Section) bool {
	_, alreadyVisited := visited[current]
	if alreadyVisited {
		return true
	}
	if current == "" {
		return false
	}
	for _, next := range sections {
		if next.Name == current {
			visited[next.Name] = struct{}{}
			return circularDependency(visited, next.From, sections)
		}
	}
	return false
}

func sortByFrom(index int, sections []Section) {
	// if section has from, move it until
	if index >= len(sections) {
		return
	}
	if sections[index].From == "" {
		sortByFrom(index+1, sections)
		return
	}
	// we check if from is before, if it is, its ok
	for i := 0; i < index; i++ {
		if sections[i].Name == sections[index].From {
			sortByFrom(index+1, sections)
			return
		}
	}
	// we have a from, find that from and move this one after that one
	hasChange := false
	for i := index + 1; i < len(sections); i++ {
		hasChange = true
		sections[i-1], sections[i] = sections[i], sections[i-1]
		if sections[i-1].Name == sections[i].From {
			break
		}
	}
	if sections[index].From != "" && hasChange {
		sortByFrom(index, sections)
		return
	}
	sortByFrom(index+1, sections)
}
