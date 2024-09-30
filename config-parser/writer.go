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
	//nolint:gosec
	"crypto/md5" // G501: Blocklisted import crypto/md5: weak cryptographic primitive
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/gofrs/flock"
	"github.com/google/renameio/maybe"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

// String returns configuration in writable form
func (p *configParser) String() string {
	if p.Options.Log {
		p.Options.Logger.Debugf("%screating string representation", p.Options.LogPrefix)
	}
	p.lock()
	defer p.unLock()
	var result strings.Builder

	p.writeParsers("", p.Parsers[Comments][CommentsSectionName], &result, false)
	p.writeParsers("global", p.Parsers[Global][GlobalSectionName], &result, true)

	sections := []Section{Defaults, UserList, Peers, Mailers, Resolvers, Cache, Ring, LogForward, HTTPErrors, Frontends, Backends, Listen, Program, FCGIApp}

	for _, section := range sections {
		var sortedSections []string
		if section == Defaults {
			var err error
			sortedSections, err = getSortedListWithFrom(p.Parsers[section])
			if err != nil && p.Options.Log {
				p.Options.Logger.Errorf("%s", err.Error())
			}
		} else {
			sortedSections = p.getSortedList(p.Parsers[section])
		}

		for _, sectionName := range sortedSections {
			var sName string
			if sectionName != "" {
				sName = fmt.Sprintf("%s %s", section, sectionName)
			} else {
				sName = string(section)
			}
			p.writeParsers(sName, p.Parsers[section][sectionName], &result, true)
		}
	}
	return result.String()
}

func (p *configParser) Save(filename string) error {
	if p.Options.Log {
		p.Options.Logger.Debugf("%ssaving configuration to file %s", p.Options.LogPrefix, filename)
	}
	if p.Options.UseMd5Hash {
		data, err := p.StringWithHash()
		if err != nil {
			return err
		}
		return p.save([]byte(data), filename)
	}
	return p.save([]byte(p.String()), filename)
}

func (p *configParser) save(data []byte, filename string) error {
	f := flock.New(filename)
	if err := f.Lock(); err != nil {
		return err
	}
	err := maybe.WriteFile(filename, data, 0o644)
	if err != nil {
		f.Unlock() //nolint:errcheck
		return err
	}
	if err := f.Unlock(); err != nil {
		errMsg := err.Error()
		return fmt.Errorf("%w %s", UnlockError{}, errMsg)
	}
	return nil
}

func (p *configParser) StringWithHash() (string, error) {
	var result strings.Builder
	content := p.String()
	//nolint:gosec
	hash := md5.Sum([]byte(content))
	result.WriteString(fmt.Sprintf("# _md5hash=%x\n", hash))
	result.WriteString(content)
	if err := p.Set(Comments, CommentsSectionName, "# _md5hash", &types.ConfigHash{Value: hex.EncodeToString(hash[:])}); err != nil {
		return "", err
	}

	return result.String(), nil
}

func (p *configParser) writeSection(sectionName string, comments []string, defaultsSection string, result io.StringWriter) {
	_, _ = result.WriteString("\n")
	for _, line := range comments {
		_, _ = result.WriteString("# ")
		_, _ = result.WriteString(line)
		_, _ = result.WriteString("\n")
	}
	_, _ = result.WriteString(sectionName)
	if defaultsSection != "" {
		_, _ = result.WriteString(" from ")
		_, _ = result.WriteString(defaultsSection)
	}
	_, _ = result.WriteString("\n")
}

func (p *configParser) writeParsers(sectionName string, parsersData *Parsers, result io.StringWriter, useIndentation bool) {
	sectionNameWritten := false
	switch sectionName {
	case "":
		sectionNameWritten = true
	case "global":
		break
	default:
		p.writeSection(sectionName, parsersData.PreComments, parsersData.DefaultSectionName, result)
		sectionNameWritten = true
	}
	for _, parserName := range parsersData.ParserSequence {
		parser := parsersData.Parsers[string(parserName)]
		lines, comments, err := parser.ResultAll()
		if err != nil {
			continue
		}
		if !sectionNameWritten {
			p.writeSection(sectionName, parsersData.PreComments, parsersData.DefaultSectionName, result)
			sectionNameWritten = true
		}
		for _, line := range comments {
			if useIndentation {
				_, _ = result.WriteString("  ")
			}
			_, _ = result.WriteString("# ")
			_, _ = result.WriteString(line)
			_, _ = result.WriteString("\n")
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
	for _, line := range parsersData.PostComments {
		if useIndentation {
			_, _ = result.WriteString("  ")
		}
		_, _ = result.WriteString("# ")
		_, _ = result.WriteString(line)
		_, _ = result.WriteString("\n")
	}
}
