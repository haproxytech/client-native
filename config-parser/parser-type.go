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
	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
)

type ParserInterface interface { //nolint:revive
	Init()
	Parse(line string, parts []string, comment string) (changeState string, err error)
	PreParse(line string, parts []string, preComments []string, comment string) (changeState string, err error)
	GetParserName() string
	Get(createIfNotExist bool) (common.ParserData, error)
	GetPreComments() ([]string, error)
	GetOne(index int) (common.ParserData, error)
	Delete(index int) error
	Insert(data common.ParserData, index int) error
	Set(data common.ParserData, index int) error
	SetPreComments(preComment []string)
	ResultAll() ([]common.ReturnResultLine, []string, error)
}

type Parsers struct {
	Parsers            map[string]ParserInterface
	ParserSequence     []Section
	PreComments        []string
	PostComments       []string
	DefaultSectionName string
}

func (p *Parsers) Get(attribute string, createIfNotExist ...bool) (common.ParserData, error) {
	createNew := false
	if len(createIfNotExist) > 0 && createIfNotExist[0] {
		createNew = true
	}
	if parser, ok := p.Parsers[attribute]; ok {
		return parser.Get(createNew)
	}

	return nil, errors.ErrParserMissing
}

func (p *Parsers) GetResult(attribute string) ([]common.ReturnResultLine, error) {
	if parser, ok := p.Parsers[attribute]; ok {
		lines, _, err := parser.ResultAll()
		if err != nil {
			return nil, errors.ErrParserMissing
		}
		return lines, nil
	}
	return nil, errors.ErrParserMissing
}

func (p *Parsers) GetPreComments(attribute string) ([]string, error) {
	if parser, ok := p.Parsers[attribute]; ok {
		return parser.GetPreComments()
	}

	return nil, errors.ErrParserMissing
}

func (p *Parsers) GetOne(attribute string, index ...int) (common.ParserData, error) {
	setIndex := -1
	if len(index) > 0 && index[0] > -1 {
		setIndex = index[0]
	}
	if parser, ok := p.Parsers[attribute]; ok {
		return parser.GetOne(setIndex)
	}

	return nil, errors.ErrParserMissing
}

// HasParser checks if we have a parser for attribute
func (p *Parsers) HasParser(attribute string) bool {
	_, hasParser := p.Parsers[attribute]
	return hasParser
}

// Set sets data in parser, if you can have multiple items, index is a must
func (p *Parsers) Set(attribute string, data common.ParserData, index ...int) error {
	setIndex := -1
	if len(index) > 0 && index[0] > -1 {
		setIndex = index[0]
	}
	if parser, ok := p.Parsers[attribute]; ok {
		return parser.Set(data, setIndex)
	}

	return errors.ErrAttributeNotFound
}

// SetPreComments sets comment lines before parser
func (p *Parsers) SetPreComments(attribute string, preComment []string) error {
	if parser, ok := p.Parsers[attribute]; ok {
		parser.SetPreComments(preComment)
		return nil
	}

	return errors.ErrAttributeNotFound
}

func (p *Parsers) Insert(attribute string, data common.ParserData, index ...int) error {
	setIndex := -1
	if len(index) > 0 && index[0] > -1 {
		setIndex = index[0]
	}
	if parser, ok := p.Parsers[attribute]; ok {
		return parser.Insert(data, setIndex)
	}

	return errors.ErrAttributeNotFound
}

func (p *Parsers) Delete(attribute string, index ...int) error {
	setIndex := -1
	if len(index) > 0 && index[0] > -1 {
		setIndex = index[0]
	}
	if parser, ok := p.Parsers[attribute]; ok {
		return parser.Delete(setIndex)
	}

	return errors.ErrAttributeNotFound
}
