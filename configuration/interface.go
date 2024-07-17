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
	"path/filepath"
	"strconv"
	"strings"

	parser "github.com/haproxytech/config-parser/v5"
	parser_options "github.com/haproxytech/config-parser/v5/options"

	"github.com/haproxytech/client-native/v6/configuration/options"
)

type Configuration interface {
	Parser
	ACL
	Backend
	Bind
	Cache
	Capture
	CrtLoad
	CrtStore
	Defaults
	DgramBind
	FCGIApp
	Filter
	Frontend
	Global
	HTTPCheck
	HTTPRequestRule
	HTTPResponseRule
	HTTPAfterResponseRule
	HTTPErrorRule
	HTTPErrorsSection
	LogTarget
	LogForward
	MailerEntry
	MailersSection
	Nameserver
	PeerEntry
	PeerSection
	Program
	Raw
	Resolver
	Ring
	Server
	ServerTemplate
	Site
	StickRule
	ServiceI
	Table
	TCPCheck
	TCPRequestRule
	TCPResponseRule
	Transactions
	TransactionHandling
	Version
	Userlist
	User
	Group
	Structured
}

func New(ctx context.Context, opt ...options.ConfigurationOption) (Configuration, error) { //nolint:revive
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

	if optionsWrapper.BackupsDir == "" {
		optionsWrapper.BackupsDir = filepath.Dir(optionsWrapper.ConfigurationFile)
	}

	if optionsWrapper.Haproxy == "" {
		optionsWrapper.Haproxy = options.DefaultHaproxy
	}

	if optionsWrapper.PreferredTimeSuffix == "" {
		optionsWrapper.PreferredTimeSuffix = options.DefaultTimeSuffix
	}

	versionString, err := c.fetchVersion(optionsWrapper.Haproxy)
	if err != nil {
		return nil, NewConfError(ErrCannotFindHAProxy, fmt.Sprintf("path to HAProxy binary not valid: %s, err: %s", optionsWrapper.Haproxy, err.Error()))
	}
	noNamedDefaultsFrom := noNamedDefaultsFrom(versionString)
	c.noNamedDefaultsFrom = true

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

	// HAProxy lower than 2.4 doesn't support from keyword to inherit defaults section, so if it's lower than that don't set it
	// if it isn't set then add from to all frontend/backend section that have it unset to the proper defaults section
	if noNamedDefaultsFrom {
		parserOptions = append(parserOptions, parser_options.NoNamedDefaultsFrom)
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

func (c *client) fetchVersion(haproxy string) (string, error) {
	versionString, err := exec.Command(haproxy, "-v").Output()
	if err != nil {
		return "", err
	}
	return string(versionString), nil
}

func getVersionNumbers(version string) (int64, int64, int64, error) {
	if !strings.HasPrefix(version, "HAProxy version ") {
		return 0, 0, 0, fmt.Errorf("not a haproxy version string")
	}
	version = version[strings.Index(version, "HAProxy version ")+len("HAProxy version "):]
	versionSlice := strings.SplitN(version, "-", 2)
	if len(versionSlice) != 2 {
		return 0, 0, 0, fmt.Errorf("not a haproxy version string")
	}

	versionInts := strings.SplitN(versionSlice[0], ".", 3)
	if len(versionInts) != 3 {
		return 0, 0, 0, fmt.Errorf("not a haproxy version string")
	}

	major, err := strconv.ParseInt(versionInts[0], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("not a haproxy version string")
	}
	minor, err := strconv.ParseInt(versionInts[1], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("not a haproxy version string")
	}
	patch, err := strconv.ParseInt(versionInts[2], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("not a haproxy version string")
	}
	return major, minor, patch, nil
}

func noNamedDefaultsFrom(version string) bool {
	major, minor, _, err := getVersionNumbers(version)
	if err != nil {
		return true
	}

	return major < 2 || (major == 2 && minor < 4)
}
