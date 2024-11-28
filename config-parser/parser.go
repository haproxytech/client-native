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
	"bytes"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/options"
)

type Section string

const (
	Comments   Section = "#"
	Defaults   Section = "defaults"
	Global     Section = "global"
	Resolvers  Section = "resolvers"
	UserList   Section = "userlist"
	Peers      Section = "peers"
	Mailers    Section = "mailers"
	Frontends  Section = "frontend"
	Backends   Section = "backend"
	Listen     Section = "listen"
	Cache      Section = "cache"
	Program    Section = "program"
	HTTPErrors Section = "http-errors"
	Ring       Section = "ring"
	LogForward Section = "log-forward"
	FCGIApp    Section = "fcgi-app"
	CrtStore   Section = "crt-store"
	Traces     Section = "traces"
	// spoe sections
	SPOEAgent   Section = "spoe-agent"
	SPOEGroup   Section = "spoe-group"
	SPOEMessage Section = "spoe-message"
)

const (
	CommentsSectionName = "data"
	GlobalSectionName   = "data"
	TracesSectionName   = ""
)

var DefaultSectionName = "" //nolint:gochecknoglobals

type Parser interface {
	LoadData(path string) error
	Process(reader io.Reader) error
	String() string
	Save(filename string) error
	StringWithHash() (string, error)
	Get(sectionType Section, sectionName string, attribute string, createIfNotExist ...bool) (common.ParserData, error)
	GetResult(sectionType Section, sectionName string, attribute string) ([]common.ReturnResultLine, error)
	GetPreComments(sectionType Section, sectionName string, attribute string) ([]string, error)
	GetOne(sectionType Section, sectionName string, attribute string, index ...int) (common.ParserData, error)
	SectionsGet(sectionType Section) ([]string, error)
	SectionsDelete(sectionType Section, sectionName string) error
	SectionsCreate(sectionType Section, sectionName string) error
	SectionsDefaultsFromGet(sectionType Section, sectionName string) (string, error)
	SectionsDefaultsFromSet(sectionType Section, sectionName, defaultsSection string) error
	Set(sectionType Section, sectionName string, attribute string, data common.ParserData, index ...int) error
	SetPreComments(sectionType Section, sectionName string, attribute string, preComment []string) error
	Delete(sectionType Section, sectionName string, attribute string, index ...int) error
	Insert(sectionType Section, sectionName string, attribute string, data common.ParserData, index ...int) error
	HasParser(sectionType Section, attribute string) bool
	SetLoggerState(active bool) error
}

type UnlockError struct{}

func (e UnlockError) Error() string {
	return "failed funclock"
}

// configParser reads and writes configuration on given file
type configParser struct {
	Parsers                 map[Section]map[string]*Parsers
	Options                 options.Parser
	lastDefaultsSectionName string
	mutex                   *sync.Mutex
}

func New(opt ...options.ParserOption) (Parser, error) {
	p := &configParser{
		Options:                 options.Parser{},
		lastDefaultsSectionName: "",
	}
	var err error

	for _, option := range opt {
		err = option.Set(&p.Options)
		if err != nil {
			return nil, err
		}
	}
	if p.Options.Path != "" {
		if p.Options.Log {
			p.Options.Logger.Infof("%sreading configuration from file: %s", p.Options.LogPrefix, p.Options.Path)
		}
		dat, err := os.ReadFile(p.Options.Path)
		if err != nil {
			return nil, err
		}
		return p, p.Process(bytes.NewReader(dat))
	}
	if p.Options.Reader != nil {
		if p.Options.Log {
			p.Options.Logger.Infof("%sreading configuration from reader", p.Options.LogPrefix)
		}
		return p, p.Process(p.Options.Reader)
	}
	if p.Options.Log {
		p.Options.Logger.Warningf("%sno configuration source provided", p.Options.LogPrefix)
	}
	return p, p.Process(strings.NewReader(""))
}
