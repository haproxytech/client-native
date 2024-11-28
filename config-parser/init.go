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
	"sync"
)

type ConfiguredParsers struct {
	_                     [0]int
	State                 string
	ActiveComments        []string
	ActiveSectionComments []string
	Active                *Parsers
	Previous              *Parsers
	HasDefaultParser      bool
	Comments              *Parsers
	Defaults              *Parsers
	Global                *Parsers
	Frontend              *Parsers
	Backend               *Parsers
	Listen                *Parsers
	Resolver              *Parsers
	Userlist              *Parsers
	Peers                 *Parsers
	Mailers               *Parsers
	Cache                 *Parsers
	Program               *Parsers
	HTTPErrors            *Parsers
	Ring                  *Parsers
	LogForward            *Parsers
	FCGIApp               *Parsers
	CrtStore              *Parsers
	Traces                *Parsers
	// spoe parsers
	SPOEAgent   *Parsers
	SPOEGroup   *Parsers
	SPOEMessage *Parsers
}

func (p *configParser) Init() {
	if p.Options.Log {
		p.Options.Logger.Debugf("%sinit", p.Options.LogPrefix)
	}
	p.initParserMaps()
	for _, sections := range p.Parsers {
		for _, parsers := range sections {
			for _, parser := range parsers.Parsers {
				parser.Init()
			}
		}
	}
}

func (p *configParser) initParserMaps() {
	p.mutex = &sync.Mutex{}
	p.lastDefaultsSectionName = ""

	p.Parsers = map[Section]map[string]*Parsers{}

	p.Parsers[Comments] = map[string]*Parsers{
		CommentsSectionName: p.getStartParser(),
	}

	p.Parsers[Defaults] = map[string]*Parsers{}

	p.Parsers[Global] = map[string]*Parsers{
		GlobalSectionName: p.getGlobalParser(),
	}

	p.Parsers[Frontends] = map[string]*Parsers{}
	p.Parsers[Backends] = map[string]*Parsers{}
	p.Parsers[Listen] = map[string]*Parsers{}
	p.Parsers[Resolvers] = map[string]*Parsers{}
	p.Parsers[UserList] = map[string]*Parsers{}
	p.Parsers[Peers] = map[string]*Parsers{}
	p.Parsers[Mailers] = map[string]*Parsers{}
	p.Parsers[Cache] = map[string]*Parsers{}
	p.Parsers[Program] = map[string]*Parsers{}
	p.Parsers[HTTPErrors] = map[string]*Parsers{}
	p.Parsers[Ring] = map[string]*Parsers{}
	p.Parsers[LogForward] = map[string]*Parsers{}
	p.Parsers[FCGIApp] = map[string]*Parsers{}
	p.Parsers[CrtStore] = map[string]*Parsers{}
	p.Parsers[Traces] = map[string]*Parsers{}
}
