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

package parser

import (
	"fmt"
	"sort"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/sorter"
)

func (p *configParser) lock() {
	p.mutex.Lock()
}

func (p *configParser) unLock() {
	p.mutex.Unlock()
}

// Get get attribute from defaults section
func (p *configParser) Get(sectionType Section, sectionName string, attribute string, createIfNotExist ...bool) (common.ParserData, error) {
	p.lock()
	defer p.unLock()
	st, ok := p.Parsers[sectionType]
	if !ok {
		return nil, errors.ErrSectionMissing
	}
	section, ok := st[sectionName]
	if !ok {
		return nil, errors.ErrSectionMissing
	}
	return section.Get(attribute, len(createIfNotExist) > 0 && createIfNotExist[0])
}

// GetResult get attribute lines from section
func (p *configParser) GetResult(sectionType Section, sectionName string, attribute string) ([]common.ReturnResultLine, error) {
	p.lock()
	defer p.unLock()
	st, ok := p.Parsers[sectionType]
	if !ok {
		return nil, errors.ErrSectionMissing
	}
	section, ok := st[sectionName]
	if !ok {
		return nil, errors.ErrSectionMissing
	}
	return section.GetResult(attribute)
}

// GetPreComments get attribute from section
func (p *configParser) GetPreComments(sectionType Section, sectionName string, attribute string) ([]string, error) {
	p.lock()
	defer p.unLock()
	st, ok := p.Parsers[sectionType]
	if !ok {
		return nil, errors.ErrSectionMissing
	}
	section, ok := st[sectionName]
	if !ok {
		return nil, errors.ErrSectionMissing
	}
	return section.GetPreComments(attribute)
}

// GetOne get attribute from section
func (p *configParser) GetOne(sectionType Section, sectionName string, attribute string, index ...int) (common.ParserData, error) {
	p.lock()
	defer p.unLock()
	setIndex := -1
	if len(index) > 0 && index[0] > -1 {
		setIndex = index[0]
	}
	st, ok := p.Parsers[sectionType]
	if !ok {
		return nil, errors.ErrSectionMissing
	}
	section, ok := st[sectionName]
	if !ok {
		return nil, errors.ErrSectionMissing
	}
	return section.GetOne(attribute, setIndex)
}

// SectionsGet lists all sections of certain type
func (p *configParser) SectionsGet(sectionType Section) ([]string, error) {
	p.lock()
	defer p.unLock()
	st, ok := p.Parsers[sectionType]
	if !ok {
		return nil, errors.ErrSectionMissing
	}
	result := make([]string, len(st))
	index := 0
	for sectionName := range st {
		result[index] = sectionName
		index++
	}
	return result, nil
}

// SectionsDelete deletes one section of sectionType
func (p *configParser) SectionsDelete(sectionType Section, sectionName string) error {
	p.lock()
	defer p.unLock()
	if _, ok := p.Parsers[sectionType]; !ok {
		return errors.ErrSectionMissing
	}
	delete(p.Parsers[sectionType], sectionName)
	return nil
}

// SectionsCreate creates one section of sectionType
func (p *configParser) SectionsCreate(sectionType Section, sectionName string) error {
	p.lock()
	defer p.unLock()
	st, ok := p.Parsers[sectionType]
	if !ok {
		return errors.ErrSectionMissing
	}
	_, ok = st[sectionName]
	if ok {
		return errors.ErrSectionAlreadyExists
	}

	parsers := ConfiguredParsers{
		State:    "",
		Active:   p.Parsers[Comments][CommentsSectionName],
		Comments: p.Parsers[Comments][CommentsSectionName],
		Global:   p.Parsers[Global][GlobalSectionName],
	}

	parts := []string{string(sectionType), sectionName}
	comment := ""
	p.ProcessLine(fmt.Sprintf("%s %s", sectionType, sectionName), parts, comment, parsers)
	return nil
}

// SectionsDefaultsFromGet returns dedicated defaults section for particular section.
// in configuration:
//
//	defaults DefaultsName
//	  mode tcp
//
//	frontend FrontendName from DefaultsName
//	  ...
//
//	SectionsDefaultsFromGet(Frontends,"FrontendName") => "DefaultsName",nil
func (p *configParser) SectionsDefaultsFromGet(sectionType Section, sectionName string) (string, error) {
	p.lock()
	defer p.unLock()
	st, ok := p.Parsers[sectionType]
	if !ok {
		return "", errors.ErrSectionTypeMissing
	}
	section, ok := st[sectionName]
	if !ok {
		return "", errors.ErrSectionMissing
	}

	return section.DefaultSectionName, nil
}

// SectionsDefaultsFromSet set default section for section.
// in configuration:
//
//	defaults DefaultsName1
//	  mode tcp
//
//	defaults DefaultsName2
//	  mode tcp
//
//	frontend FrontendName from DefaultsName1
//	  ...
//
//	SectionsDefaultsFromSet(Frontends,"FrontendName", "DefaultsName2") => nil
//
// ONLY defaults, frontend, backend and listen sections can be used here
func (p *configParser) SectionsDefaultsFromSet(sectionType Section, sectionName, defaultsSection string) error {
	p.lock()
	defer p.unLock()
	switch sectionType { //nolint:exhaustive
	case Defaults:
		// for defaults we need to check for circular dependencies
		// and whether that section exists or not due to sorting method
		sections := p.Parsers[Defaults]
		listDefaults := make([]sorter.Section, len(sections))
		i := 0
		for k, v := range sections {
			s := sorter.Section{
				Name: k,
				From: v.DefaultSectionName,
			}
			if s.Name == sectionName {
				s.From = defaultsSection
			}
			listDefaults[i] = s
			i++
		}
		err := sorter.Sort(listDefaults)
		if err != nil {
			return err
		}
	case Frontends, Backends, Listen:
	default:
		// catch all other sections
		return errors.ErrSectionTypeNotAllowed
	}

	st, ok := p.Parsers[sectionType]
	if !ok {
		return errors.ErrSectionTypeMissing
	}
	section, ok := st[sectionName]
	if !ok {
		return errors.ErrSectionMissing
	}
	// if not set to empty, do extra validation
	if defaultsSection != "" {
		stDef, ok := p.Parsers[Defaults]
		if !ok {
			return errors.ErrSectionTypeMissing
		}
		_, ok = stDef[defaultsSection]
		if !ok {
			return errors.ErrFromDefaultsSectionMissing
		}
	}
	section.DefaultSectionName = defaultsSection
	return nil
}

// Set sets attribute from section, can be nil to disable/remove
func (p *configParser) Set(sectionType Section, sectionName string, attribute string, data common.ParserData, index ...int) error {
	p.lock()
	defer p.unLock()
	setIndex := -1
	if len(index) > 0 && index[0] > -1 {
		setIndex = index[0]
	}
	st, ok := p.Parsers[sectionType]
	if !ok {
		return errors.ErrSectionMissing
	}
	section, ok := st[sectionName]
	if !ok {
		return errors.ErrSectionMissing
	}
	return section.Set(attribute, data, setIndex)
}

// Set sets line comment before attribute from section, can be nil to disable/remove
func (p *configParser) SetPreComments(sectionType Section, sectionName string, attribute string, preComment []string) error {
	p.lock()
	defer p.unLock()
	st, ok := p.Parsers[sectionType]
	if !ok {
		return errors.ErrSectionMissing
	}
	section, ok := st[sectionName]
	if !ok {
		return errors.ErrSectionMissing
	}
	return section.SetPreComments(attribute, preComment)
}

// Delete remove attribute on defined index, in case of single attributes, index is ignored
func (p *configParser) Delete(sectionType Section, sectionName string, attribute string, index ...int) error {
	p.lock()
	defer p.unLock()
	setIndex := -1
	if len(index) > 0 && index[0] > -1 {
		setIndex = index[0]
	}
	st, ok := p.Parsers[sectionType]
	if !ok {
		return errors.ErrSectionMissing
	}
	section, ok := st[sectionName]
	if !ok {
		return errors.ErrSectionMissing
	}
	return section.Delete(attribute, setIndex)
}

// Insert put attribute on defined index, in case of single attributes, index is ignored
func (p *configParser) Insert(sectionType Section, sectionName string, attribute string, data common.ParserData, index ...int) error {
	p.lock()
	defer p.unLock()
	setIndex := -1
	if len(index) > 0 && index[0] > -1 {
		setIndex = index[0]
	}
	st, ok := p.Parsers[sectionType]
	if !ok {
		return errors.ErrSectionMissing
	}
	section, ok := st[sectionName]
	if !ok {
		return errors.ErrSectionMissing
	}
	return section.Insert(attribute, data, setIndex)
}

// HasParser checks if we have a parser for attribute
func (p *configParser) HasParser(sectionType Section, attribute string) bool {
	p.lock()
	defer p.unLock()
	st, ok := p.Parsers[sectionType]
	if !ok {
		return false
	}
	sectionName := ""
	for name := range st {
		sectionName = name

		break
	}
	section, ok := st[sectionName]
	if !ok {
		return false
	}
	return section.HasParser(attribute)
}

func (p *configParser) getSortedList(data map[string]*Parsers) []string {
	result := make([]string, len(data))
	index := 0
	for parserSectionName := range data {
		result[index] = parserSectionName
		index++
	}
	sort.Strings(result)
	return result
}

// getSortedListWithFrom returns list of parses sorted,
// since every section can have a from that might depend on,
// we take that into account
func getSortedListWithFrom(data map[string]*Parsers) ([]string, error) {
	sortedSections := make([]sorter.Section, len(data))
	index := 0
	for parserSectionName := range data {
		sortedSections[index] = sorter.Section{
			Name: parserSectionName,
			From: data[parserSectionName].DefaultSectionName,
		}
		index++
	}
	err := sorter.Sort(sortedSections)
	result := make([]string, len(data))
	for index, value := range sortedSections {
		result[index] = value.Name
	}
	return result, err
}
