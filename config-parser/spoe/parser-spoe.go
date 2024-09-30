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

package spoe

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/gofrs/flock"
	"github.com/google/renameio/maybe"
	parser "github.com/haproxytech/client-native/v5/config-parser"
	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/spoe/parsers"
	"github.com/haproxytech/client-native/v5/config-parser/spoe/types"
)

/*type parser.Section string

const (
	//spoe sections
	SPOEAgent   Section = "spoe-agent"
	SPOEGroup   Section = "spoe-group"
	SPOEMessage Section = "spoe-message"
)

const (
	CommentsSectionName = "data"
	GlobalSectionName   = "data"
	DefaultSectionName  = "data"
)*/

// Parser reads and writes configuration on given file
type Parser struct {
	Parsers map[string]map[parser.Section]map[string]*parser.Parsers
	mutex   *sync.Mutex
}

func (p *Parser) lock() {
	p.mutex.Lock()
}

func (p *Parser) unLock() {
	p.mutex.Unlock()
}

// Get get attribute from defaults section
func (p *Parser) Get(scope string, sectionType parser.Section, sectionName string, attribute string, createIfNotExist ...bool) (common.ParserData, error) {
	p.lock()
	defer p.unLock()
	psrs, ok := p.Parsers[scope]
	if !ok {
		return nil, errors.ErrScopeMissing
	}
	st, ok := psrs[sectionType]
	if !ok {
		return nil, errors.ErrSectionMissing
	}
	section, ok := st[sectionName]
	if !ok {
		return nil, errors.ErrSectionMissing
	}
	createNew := false
	if len(createIfNotExist) > 0 && createIfNotExist[0] {
		createNew = true
	}
	return section.Get(attribute, createNew)
}

// GetOne get attribute from defaults section
func (p *Parser) GetOne(scope string, sectionType parser.Section, sectionName string, attribute string, index ...int) (common.ParserData, error) {
	p.lock()
	defer p.unLock()
	setIndex := -1
	if len(index) > 0 && index[0] > -1 {
		setIndex = index[0]
	}
	psrs, ok := p.Parsers[scope]
	if !ok {
		return nil, errors.ErrScopeMissing
	}
	st, ok := psrs[sectionType]
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
func (p *Parser) SectionsGet(scope string, sectionType parser.Section) ([]string, error) {
	p.lock()
	defer p.unLock()
	psrs, ok := p.Parsers[scope]
	if !ok {
		return nil, errors.ErrScopeMissing
	}
	st, ok := psrs[sectionType]
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

// ScopeDelete deletes one section of sectionType
func (p *Parser) ScopeDelete(scope string) error {
	p.lock()
	defer p.unLock()
	if _, ok := p.Parsers[scope]; !ok {
		return errors.ErrScopeMissing
	}
	delete(p.Parsers, scope)
	return nil
}

// SectionsDelete deletes one section of sectionType
func (p *Parser) SectionsDelete(scope string, sectionType parser.Section, sectionName string) error {
	p.lock()
	defer p.unLock()
	psrs, ok := p.Parsers[scope]
	if !ok {
		return errors.ErrScopeMissing
	}
	_, ok = psrs[sectionType]
	if !ok {
		return errors.ErrSectionMissing
	}
	delete(psrs[sectionType], sectionName)
	return nil
}

// ScopeCreate creates one section of sectionType
func (p *Parser) ScopeCreate(scope string) error {
	p.lock()
	defer p.unLock()

	if _, ok := p.Parsers[scope]; ok {
		return errors.ErrScopeAlreadyExists
	}
	par := map[parser.Section]map[string]*parser.Parsers{}
	p.Parsers[scope] = par
	par[parser.Comments] = map[string]*parser.Parsers{
		parser.CommentsSectionName: getStartParser(),
	}

	par[parser.SPOEAgent] = map[string]*parser.Parsers{}
	par[parser.SPOEGroup] = map[string]*parser.Parsers{}
	par[parser.SPOEMessage] = map[string]*parser.Parsers{}
	return nil
}

func (p *Parser) SectionsCreate(scope string, sectionType parser.Section, sectionName string) error {
	p.lock()
	defer p.unLock()
	psrs, ok := p.Parsers[scope]
	if !ok {
		return errors.ErrScopeMissing
	}
	st, ok := psrs[sectionType]
	if !ok {
		return errors.ErrSectionMissing
	}
	_, ok = st[sectionName]
	if ok {
		return errors.ErrSectionAlreadyExists
	}

	parsers := parser.ConfiguredParsers{
		State:    "",
		Active:   p.Parsers[scope][parser.Comments][parser.CommentsSectionName],
		Comments: p.Parsers[scope][parser.Comments][parser.CommentsSectionName],
	}

	parts := []string{string(sectionType), sectionName}
	comment := ""
	p.ProcessLine(fmt.Sprintf("%s %s", sectionType, sectionName), parts, comment, parsers, scope)
	return nil
}

// Set sets attribute from defaults section, can be nil to disable/remove
func (p *Parser) Set(scope string, sectionType parser.Section, sectionName string, attribute string, data common.ParserData, index ...int) error {
	p.lock()
	defer p.unLock()
	setIndex := -1
	if len(index) > 0 && index[0] > -1 {
		setIndex = index[0]
	}
	psrs, ok := p.Parsers[scope]
	if !ok {
		return errors.ErrScopeMissing
	}
	st, ok := psrs[sectionType]
	if !ok {
		return errors.ErrSectionMissing
	}
	section, ok := st[sectionName]
	if !ok {
		return errors.ErrSectionMissing
	}
	return section.Set(attribute, data, setIndex)
}

// Delete remove attribute on defined index, in case of single attributes, index is ignored
func (p *Parser) Delete(scope string, sectionType parser.Section, sectionName string, attribute string, index ...int) error {
	p.lock()
	defer p.unLock()
	setIndex := -1
	if len(index) > 0 && index[0] > -1 {
		setIndex = index[0]
	}
	psrs, ok := p.Parsers[scope]
	if !ok {
		return errors.ErrScopeMissing
	}
	st, ok := psrs[sectionType]
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
func (p *Parser) Insert(scope string, sectionType parser.Section, sectionName string, attribute string, data common.ParserData, index ...int) error {
	p.lock()
	defer p.unLock()
	setIndex := -1
	if len(index) > 0 && index[0] > -1 {
		setIndex = index[0]
	}
	psrs, ok := p.Parsers[scope]
	if !ok {
		return errors.ErrScopeMissing
	}
	st, ok := psrs[sectionType]
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
func (p *Parser) HasParser(scope string, sectionType parser.Section, attribute string) bool {
	p.lock()
	defer p.unLock()
	psrs, ok := p.Parsers[scope]
	if !ok {
		return false
	}
	st, ok := psrs[sectionType]
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

func (p *Parser) writeParsers(sectionName string, parsersData *parser.Parsers, result io.StringWriter, useIndentation bool) {
	sectionNameWritten := false
	if sectionName == "" {
		sectionNameWritten = true
	}
	for _, parserName := range parsersData.ParserSequence {
		parser := parsersData.Parsers[string(parserName)]
		lines, _, err := parser.ResultAll()
		if err != nil {
			continue
		}
		if !sectionNameWritten {
			_, _ = result.WriteString("\n")
			_, _ = result.WriteString(sectionName)
			_, _ = result.WriteString("\n")
			sectionNameWritten = true
		}
		for _, line := range lines {
			if useIndentation {
				_, _ = result.WriteString("  ")
			}
			_, _ = result.WriteString(line.Data)
			if line.Comment != "" {
				_, _ = result.WriteString(" # ")
				_, _ = result.WriteString(line.Comment)
			}
			_, _ = result.WriteString("\n")
		}
	}
}

func (p *Parser) getSortedList(data map[string]*parser.Parsers) []string {
	result := make([]string, len(data))
	index := 0
	for parserSectionName := range data {
		result[index] = parserSectionName
		index++
	}
	sort.Strings(result)
	return result
}

// String returns configuration in writable form
func (p *Parser) String() string {
	p.lock()
	defer p.unLock()
	var result strings.Builder

	scopes := []string{}
	for scope := range p.Parsers {
		scopes = append(scopes, scope)
	}
	sort.Strings(scopes)
	firstScope := true
	// for scope, data := range p.Parsers {
	for _, scope := range scopes {
		data := p.Parsers[scope]
		if scope != "" {
			if !firstScope {
				result.WriteString("\n")
			} else {
				firstScope = false
			}
			result.WriteString(scope)
			// result.WriteString("\n")
		}

		p.writeParsers("", data[parser.Comments][parser.CommentsSectionName], &result, false)

		sections := []parser.Section{parser.SPOEAgent, parser.SPOEGroup, parser.SPOEMessage}

		for _, section := range sections {
			sortedSections := p.getSortedList(data[section])
			for _, sectionName := range sortedSections {
				p.writeParsers(fmt.Sprintf("%s %s", section, sectionName), data[section][sectionName], &result, true)
			}
		}
	}
	return result.String()
}

func (p *Parser) Save(filename string) error {
	f := flock.New(filename)
	if err := f.Lock(); err != nil {
		return err
	}
	d1 := []byte(p.String())
	err := maybe.WriteFile(filename, d1, 0o644)
	if err != nil {
		f.Unlock() //nolint:errcheck
		return err
	}
	if err := f.Unlock(); err != nil {
		errMsg := err.Error()
		return fmt.Errorf("%w %s", parser.UnlockError{}, errMsg)
	}
	return nil
}

func (p *Parser) IsScope(line string) bool {
	if line == "" {
		return false
	}
	if line[0] == '[' && line[len(line)-1] == ']' {
		return true
	}
	return false
}

// ProcessLine parses line plus determines if we need to change state
func (p *Parser) ProcessLine(line string, parts []string, comment string, config parser.ConfiguredParsers, scope string) (psrs parser.ConfiguredParsers, resultScope string) { //nolint:nonamedreturns
	if p.IsScope(line) {
		scope = line
		_ = p.ScopeCreate(scope)
		return config, scope
	}
	for _, section := range config.Active.ParserSequence {
		prsr := config.Active.Parsers[string(section)]
		if newState, err := prsr.Parse(line, parts, comment); err == nil {
			// should we have an option to remove it when found?
			if newState != "" {
				// log.Printf("change state from %s to %s\n", state, newState)
				config.State = newState
				if config.State == "" {
					config.Active = config.Comments
				}
				if config.State == string(parser.SPOEAgent) {
					parserSectionName := prsr.(*parsers.SPOESection) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.SPOESection) //nolint:forcetypeassert
					config.SPOEAgent = getSPOEAgentParser()
					p.Parsers[scope][parser.SPOEAgent][data.Name] = config.SPOEAgent
					config.Active = config.SPOEAgent
				}
				if config.State == string(parser.SPOEGroup) {
					parserSectionName := prsr.(*parsers.SPOESection) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.SPOESection) //nolint:forcetypeassert
					config.SPOEGroup = getSPOEGroupParser()
					p.Parsers[scope][parser.SPOEGroup][data.Name] = config.SPOEGroup
					config.Active = config.SPOEGroup
				}
				if config.State == string(parser.SPOEMessage) {
					parserSectionName := prsr.(*parsers.SPOESection) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.SPOESection) //nolint:forcetypeassert
					config.SPOEMessage = getSPOEMessageParser()
					p.Parsers[scope][parser.SPOEMessage][data.Name] = config.SPOEMessage
					config.Active = config.SPOEMessage
				}
			}
			break
		}
	}
	return config, scope
}

func (p *Parser) LoadData(filename string) error {
	dat, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return p.ParseData(string(dat))
}

func (p *Parser) ParseData(dat string) error {
	p.mutex = &sync.Mutex{}

	p.Parsers = map[string]map[parser.Section]map[string]*parser.Parsers{}
	par := map[parser.Section]map[string]*parser.Parsers{}
	p.Parsers[""] = par
	par[parser.Comments] = map[string]*parser.Parsers{
		parser.CommentsSectionName: getStartParser(),
	}

	par[parser.SPOEAgent] = map[string]*parser.Parsers{}
	par[parser.SPOEGroup] = map[string]*parser.Parsers{}
	par[parser.SPOEMessage] = map[string]*parser.Parsers{}

	parsers := parser.ConfiguredParsers{
		State:    "",
		Active:   par[parser.Comments][parser.CommentsSectionName],
		Comments: par[parser.Comments][parser.CommentsSectionName],
	}

	lines := common.StringSplitIgnoreEmpty(dat, '\n')
	scope := ""
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts, comment := common.StringSplitWithCommentIgnoreEmpty(line)
		if len(parts) == 0 && comment != "" {
			parts = []string{""}
		}
		if len(parts) == 0 {
			continue
		}
		// this is the difference, no previous line is sent to parsers
		parsers, scope = p.ProcessLine(line, parts, comment, parsers, scope)
	}
	return nil
}
