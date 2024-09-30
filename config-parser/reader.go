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
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/parsers/extra"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

func (p *configParser) LoadData(path string) error {
	if p.Options.Log {
		p.Options.Logger.Debugf("%sreading data from %s", p.Options.LogPrefix, path)
	}
	p.Options.Path = path
	dat, err := os.ReadFile(p.Options.Path)
	if err != nil {
		return err
	}
	return p.Process(bytes.NewReader(dat))
}

func (p *configParser) Process(reader io.Reader) error {
	p.Init()

	parsers := ConfiguredParsers{
		State:          "",
		ActiveComments: nil,
		Active:         p.Parsers[Comments][CommentsSectionName],
		Comments:       p.Parsers[Comments][CommentsSectionName],
		Global:         p.Parsers[Global][GlobalSectionName],
	}

	bufferedScanner := bufio.NewScanner(reader)

	var line string

	if p.Options.Log {
		p.Options.Logger.Debugf("%sprocessing of data started", p.Options.LogPrefix)
	}
	for bufferedScanner.Scan() {
		line = bufferedScanner.Text()

		if line == "" {
			if parsers.State == "" {
				parsers.State = "#"
			}
			continue
		}
		parts, comment := common.StringSplitWithCommentIgnoreEmpty(line)
		if len(parts) == 0 && comment != "" {
			switch {
			case strings.HasPrefix(line, "# _version"):
				parts = []string{"# _version"}
			case strings.Contains(line, "config-snippet"):
				parts = []string{"config-snippet"}
			case strings.HasPrefix(line, "# _md5hash"):
				parts = []string{"# _md5hash"}
			default:
				parts = []string{""}
			}
		}
		if len(parts) == 0 {
			continue
		}
		if p.Options.Log {
			p.Options.Logger.Tracef("%sprocessing line: %s", p.Options.LogPrefix, line)
		}
		parsers = p.ProcessLine(line, parts, comment, parsers)
	}
	if parsers.ActiveComments != nil {
		parsers.Active.PostComments = parsers.ActiveComments
	}
	if parsers.ActiveSectionComments != nil {
		parsers.Active.PostComments = append(parsers.Active.PostComments, parsers.ActiveSectionComments...)
	}
	if p.Options.Log {
		p.Options.Logger.Debugf("%sprocessing of data ended", p.Options.LogPrefix)
	}
	return nil
}

// ProcessLine parses line plus determines if we need to change state
func (p *configParser) ProcessLine(line string, parts []string, comment string, config ConfiguredParsers) ConfiguredParsers { //nolint:gocognit,gocyclo,cyclop,maintidx
	if config.State != "" {
		if parts[0] == "" && comment != "" && comment != "##_config-snippet_### BEGIN" && comment != "##_config-snippet_### END" {
			if line[0] == ' ' {
				if config.State != "snippet_beg" {
					config.ActiveComments = append(config.ActiveComments, comment)
					return config
				}
			} else {
				config.ActiveSectionComments = append(config.ActiveSectionComments, comment)
				return config
			}
		}
	}
	parsers := make([]ParserInterface, 0, 2)
	if !p.Options.DisableUnProcessed {
		parsers = append(parsers, config.Active.Parsers[""])
	}

	if config.HasDefaultParser {
		// Default parser name is given in position 0 of ParserSequence
		parsers = append(parsers, config.Active.Parsers[string(config.Active.ParserSequence[0])])
	}
	// We add iteratively the different parts to form a potential parser name
	for i := 1; i <= len(parts) && !config.HasDefaultParser; i++ {
		parserName := strings.Join(parts[:i], " ")
		if parserName == "" {
			continue
		}
		if parserFound, ok := config.Active.Parsers[parserName]; ok {
			parsers = append(parsers, parserFound)
			break
		}
	}
	if len(parsers) < 2 && len(parts) == 1 && parts[0] == "" {
		if parserFound, ok := config.Active.Parsers["#"]; ok {
			parsers = append(parsers, parserFound)
		}
	}
	if (len(parsers) < 2) && len(parts) > 0 && parts[0] == "no" {
		for i := 2; i <= len(parts) && !config.HasDefaultParser; i++ {
			if parserFound, ok := config.Active.Parsers[strings.Join(parts[1:i], " ")]; ok {
				parsers = append(parsers, parserFound)
				break
			}
		}
	}

	for i := len(parsers) - 1; i >= 0; i-- {
		parser := parsers[i]
		if p.Options.Log {
			p.Options.Logger.Tracef("%susing parser [%s]", p.Options.LogPrefix, parser.GetParserName())
		}
		if newState, err := parser.PreParse(line, parts, config.ActiveComments, comment); err == nil {
			if newState != "" {
				if p.Options.Log {
					p.Options.Logger.Debugf("%schange active section to %s\n", p.Options.LogPrefix, newState)
				}
				if config.ActiveComments != nil {
					config.Active.PostComments = config.ActiveComments
				}
				config.State = newState
				switch config.State {
				case "":
					config.Active = config.Comments
				case "defaults":
					parserSectionName := parser.(*extra.Section) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.Section) //nolint:forcetypeassert
					if data.Name == "" {
						var name string
						for i := 1; ; i++ {
							name = "unnamed_defaults_" + strconv.Itoa(i)
							_, exists := p.Parsers[Defaults][name]
							if !exists {
								break
							}
						}
						data.Name = name
					}
					config.Defaults = p.getDefaultParser()
					if data.FromDefaults != "" {
						config.Defaults.DefaultSectionName = data.FromDefaults
					}
					p.Parsers[Defaults][data.Name] = config.Defaults
					config.Active = config.Defaults
					if p.Options.Log {
						p.Options.Logger.Tracef("%defaults section %s active", p.Options.LogPrefix, data.Name)
					}
					if !p.Options.NoNamedDefaultsFrom {
						p.lastDefaultsSectionName = data.Name
					}
					DefaultSectionName = data.Name
				case "global":
					config.Active = config.Global
				case "frontend":
					parserSectionName := parser.(*extra.Section) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.Section) //nolint:forcetypeassert
					config.Frontend = p.getFrontendParser()
					if data.FromDefaults != "" {
						config.Frontend.DefaultSectionName = data.FromDefaults
					} else {
						config.Frontend.DefaultSectionName = p.lastDefaultsSectionName
					}
					p.Parsers[Frontends][data.Name] = config.Frontend
					config.Active = config.Frontend
					if p.Options.Log {
						p.Options.Logger.Tracef("%sfrontend section %s active", p.Options.LogPrefix, data.Name)
					}
				case "backend":
					parserSectionName := parser.(*extra.Section) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.Section) //nolint:forcetypeassert
					config.Backend = p.getBackendParser()
					if data.FromDefaults != "" {
						config.Backend.DefaultSectionName = data.FromDefaults
					} else {
						config.Backend.DefaultSectionName = p.lastDefaultsSectionName
					}
					p.Parsers[Backends][data.Name] = config.Backend
					config.Active = config.Backend
					if p.Options.Log {
						p.Options.Logger.Tracef("%sbackend section %s active", p.Options.LogPrefix, data.Name)
					}
				case "listen":
					parserSectionName := parser.(*extra.Section) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.Section) //nolint:forcetypeassert
					config.Listen = p.getListenParser()
					if data.FromDefaults != "" {
						config.Listen.DefaultSectionName = data.FromDefaults
					} else {
						config.Listen.DefaultSectionName = p.lastDefaultsSectionName
					}
					p.Parsers[Listen][data.Name] = config.Listen
					config.Active = config.Listen
					if p.Options.Log {
						p.Options.Logger.Tracef("%slisten section %s active", p.Options.LogPrefix, data.Name)
					}
				case "resolvers":
					parserSectionName := parser.(*extra.Section) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.Section) //nolint:forcetypeassert
					config.Resolver = p.getResolverParser()
					p.Parsers[Resolvers][data.Name] = config.Resolver
					config.Active = config.Resolver
					if p.Options.Log {
						p.Options.Logger.Tracef("%sresolvers section %s active", p.Options.LogPrefix, data.Name)
					}
				case "userlist":
					parserSectionName := parser.(*extra.Section) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.Section) //nolint:forcetypeassert
					config.Userlist = p.getUserlistParser()
					p.Parsers[UserList][data.Name] = config.Userlist
					config.Active = config.Userlist
					if p.Options.Log {
						p.Options.Logger.Tracef("%suserlist section %s active", p.Options.LogPrefix, data.Name)
					}
				case "fcgi-app":
					parserSectionName := parser.(*extra.Section) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.Section) //nolint:forcetypeassert
					config.FCGIApp = p.getFcgiAppParser()
					p.Parsers[FCGIApp][data.Name] = config.FCGIApp
					config.Active = config.FCGIApp
					if p.Options.Log {
						p.Options.Logger.Tracef("%sfcgi-app section %s active", p.Options.LogPrefix, data.Name)
					}
				case "peers":
					parserSectionName := parser.(*extra.Section) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.Section) //nolint:forcetypeassert
					config.Peers = p.getPeersParser()
					p.Parsers[Peers][data.Name] = config.Peers
					config.Active = config.Peers
					if p.Options.Log {
						p.Options.Logger.Tracef("%speers section %s active", p.Options.LogPrefix, data.Name)
					}
				case "mailers":
					parserSectionName := parser.(*extra.Section) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.Section) //nolint:forcetypeassert
					config.Mailers = p.getMailersParser()
					p.Parsers[Mailers][data.Name] = config.Mailers
					config.Active = config.Mailers
					if p.Options.Log {
						p.Options.Logger.Tracef("%smailers section %s active", p.Options.LogPrefix, data.Name)
					}
				case "cache":
					parserSectionName := parser.(*extra.Section) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.Section) //nolint:forcetypeassert
					config.Cache = p.getCacheParser()
					p.Parsers[Cache][data.Name] = config.Cache
					config.Active = config.Cache
					if p.Options.Log {
						p.Options.Logger.Tracef("%scache section %s active", p.Options.LogPrefix, data.Name)
					}
				case "program":
					parserSectionName := parser.(*extra.Section) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.Section) //nolint:forcetypeassert
					config.Program = p.getProgramParser()
					p.Parsers[Program][data.Name] = config.Program
					config.Active = config.Program
					if p.Options.Log {
						p.Options.Logger.Tracef("%sprogram section %s active", p.Options.LogPrefix, data.Name)
					}
				case "http-errors":
					parserSectionName := parser.(*extra.Section) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.Section) //nolint:forcetypeassert
					config.HTTPErrors = p.getHTTPErrorsParser()
					p.Parsers[HTTPErrors][data.Name] = config.HTTPErrors
					config.Active = config.HTTPErrors
					if p.Options.Log {
						p.Options.Logger.Tracef("%shttp-errors section %s active", p.Options.LogPrefix, data.Name)
					}
				case "ring":
					parserSectionName := parser.(*extra.Section) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.Section) //nolint:forcetypeassert
					config.Ring = p.getRingParser()
					p.Parsers[Ring][data.Name] = config.Ring
					config.Active = config.Ring
					if p.Options.Log {
						p.Options.Logger.Tracef("%sring section %s active", p.Options.LogPrefix, data.Name)
					}
				case "log-forward":
					parserSectionName := parser.(*extra.Section) //nolint:forcetypeassert
					rawData, _ := parserSectionName.Get(false)
					data := rawData.(*types.Section) //nolint:forcetypeassert
					config.LogForward = p.getLogForwardParser()
					p.Parsers[LogForward][data.Name] = config.LogForward
					config.Active = config.LogForward
					if p.Options.Log {
						p.Options.Logger.Tracef("%log-forward section %s active", p.Options.LogPrefix, data.Name)
					}
				case "snippet_beg":
					config.Previous = config.Active
					config.Active = &Parsers{
						Parsers:        map[string]ParserInterface{"config-snippet": parser},
						ParserSequence: []Section{"config-snippet"},
					}
					config.HasDefaultParser = true
				case "snippet_end":
					config.Active = config.Previous
					config.HasDefaultParser = false
				}
				if config.ActiveSectionComments != nil {
					config.Active.PreComments = config.ActiveSectionComments
					config.ActiveSectionComments = nil
				}
			}
			config.ActiveComments = nil
			break
		}
	}

	return config
}
