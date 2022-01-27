// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package configuration

import (
	"context"
	"fmt"
	"os/exec"

	parser "github.com/haproxytech/config-parser/v4"
	parser_options "github.com/haproxytech/config-parser/v4/options"

	"github.com/haproxytech/client-native/v3/configuration/options"
)

type Configuration interface {
	Parser
	ACL
	Backend
	Bind
	Cache
	Capture
	Defaults
	Frontend
	Filter
	Global
	HTTPCheck
	HTTPRequestRule
	HTTPResponseRule
	LogTarget
	Nameserver
	PeerEntry
	PeerSection
	Raw
	Resolver
	Server
	ServerTemplate
	Site
	StickRule
	ServiceI
	TCPCheck
	TCPRequestRule
	TCPResponseRule
	Transactions
	TransactionHandling
	Version
	Userlist
	User
	Group
}

func New(ctx context.Context, opt ...options.ConfigurationOption) (Configuration, error) {
	c := &client{}
	var err error

	optionsWrapper := options.ConfigurationOptions{}

	for _, option := range opt {
		err = option.Set(&optionsWrapper)
		if err != nil {
			return nil, err
		}
	}

	if optionsWrapper.TransactionDir == "" {
		optionsWrapper.TransactionDir = options.DefaultTransactionDir
	}

	if optionsWrapper.ConfigurationFile == "" {
		optionsWrapper.ConfigurationFile = options.DefaultConfigurationFile
	}

	if optionsWrapper.Haproxy == "" {
		optionsWrapper.Haproxy = options.DefaultHaproxy
	}

	// #nosec G204
	if err1 := exec.Command(optionsWrapper.Haproxy, "-v").Run(); err1 != nil {
		return nil, NewConfError(ErrCannotFindHAProxy, fmt.Sprintf("path to HAProxy binary not valid: %s", optionsWrapper.Haproxy))
	}

	c.TransactionClient = c
	c.ConfigurationOptions = optionsWrapper

	c.parsers = make(map[string]parser.Parser)
	c.services = make(map[string]*Service)
	if err2 := c.InitTransactionParsers(); err2 != nil {
		return nil, err2
	}

	parserOptions := []parser_options.ParserOption{}
	if c.ConfigurationOptions.UseMd5Hash {
		parserOptions = append(parserOptions, parser_options.UseMd5Hash)
	}
	parserOptions = append(parserOptions, parser_options.Path(optionsWrapper.ConfigurationFile))
	p, err := parser.New(parserOptions...)
	if err != nil {
		return nil, NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read %s", c.ConfigurationFile))
	}

	c.parser = p

	return c, nil
}

func (c *client) Parser() parser.Parser {
	return c.parser
}
