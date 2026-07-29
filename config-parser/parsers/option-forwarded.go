/*
Copyright 2026 HAProxy Technologies

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

package parsers

import (
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type OptionForwardedAttributeName string

const (
	optionForwardedHost      OptionForwardedAttributeName = "host"
	optionForwardedHostExpr  OptionForwardedAttributeName = "host-expr"
	optionForwardedBy        OptionForwardedAttributeName = "by"
	optionForwardedByExpr    OptionForwardedAttributeName = "by-expr"
	optionForwardedByPort    OptionForwardedAttributeName = "by_port"
	optionForwardedByPortExp OptionForwardedAttributeName = "by_port-expr"
	optionForwardedFor       OptionForwardedAttributeName = "for"
	optionForwardedForExpr   OptionForwardedAttributeName = "for-expr"
	optionForwardedForPort   OptionForwardedAttributeName = "for_port"
	optionForwardedForPortEx OptionForwardedAttributeName = "for_port-expr"
	optionForwardedProto     OptionForwardedAttributeName = "proto"
)

type OptionForwarded struct {
	data        *types.OptionForwarded
	preComments []string // comments that appear before the actual line
}

type optionForwardedAttribute struct {
	family     OptionForwardedAttributeName
	needsValue bool
	set        func(*types.OptionForwarded, string)
}

func optionForwardedAttributeByName(name OptionForwardedAttributeName) (optionForwardedAttribute, bool) {
	switch name {
	case optionForwardedProto:
		return optionForwardedAttribute{
			family: optionForwardedProto,
			set: func(data *types.OptionForwarded, _ string) {
				data.Proto = true
			},
		}, true
	case optionForwardedHost:
		return optionForwardedAttribute{
			family: optionForwardedHost,
			set: func(data *types.OptionForwarded, _ string) {
				data.Host = true
			},
		}, true
	case optionForwardedHostExpr:
		return optionForwardedAttribute{
			family:     optionForwardedHost,
			needsValue: true,
			set:        func(data *types.OptionForwarded, value string) { data.HostExpr = value },
		}, true
	case optionForwardedBy:
		return optionForwardedAttribute{
			family: optionForwardedBy,
			set: func(data *types.OptionForwarded, _ string) {
				data.By = true
			},
		}, true
	case optionForwardedByExpr:
		return optionForwardedAttribute{
			family:     optionForwardedBy,
			needsValue: true,
			set:        func(data *types.OptionForwarded, value string) { data.ByExpr = value },
		}, true
	case optionForwardedByPort:
		return optionForwardedAttribute{
			family: optionForwardedByPort,
			set: func(data *types.OptionForwarded, _ string) {
				data.ByPort = true
			},
		}, true
	case optionForwardedByPortExp:
		return optionForwardedAttribute{
			family:     optionForwardedByPort,
			needsValue: true,
			set:        func(data *types.OptionForwarded, value string) { data.ByPortExpr = value },
		}, true
	case optionForwardedFor:
		return optionForwardedAttribute{
			family: optionForwardedFor,
			set: func(data *types.OptionForwarded, _ string) {
				data.For = true
			},
		}, true
	case optionForwardedForExpr:
		return optionForwardedAttribute{
			family:     optionForwardedFor,
			needsValue: true,
			set:        func(data *types.OptionForwarded, value string) { data.ForExpr = value },
		}, true
	case optionForwardedForPort:
		return optionForwardedAttribute{
			family: optionForwardedForPort,
			set: func(data *types.OptionForwarded, _ string) {
				data.ForPort = true
			},
		}, true
	case optionForwardedForPortEx:
		return optionForwardedAttribute{
			family:     optionForwardedForPort,
			needsValue: true,
			set:        func(data *types.OptionForwarded, value string) { data.ForPortExpr = value },
		}, true
	default:
		return optionForwardedAttribute{}, false
	}
}

/*
option forwarded [ proto ] [ host | host-expr <host_expr> ] [ by | by-expr <by_expr> ] [ by_port | by_port-expr <by_port_expr>] [ for | for-expr <for_expr> ] [ for_port | for_port-expr <for_port_expr>]
*/
func (s *OptionForwarded) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 2 && parts[0] == "no" && parts[1] == "option" && parts[2] == "forwarded" {
		if len(parts) != 3 {
			return "", errors.ErrInvalidData
		}
		s.data = &types.OptionForwarded{
			NoOption: true,
			Comment:  comment,
		}
		return "", nil
	}
	if len(parts) > 1 && parts[0] == "option" && parts[1] == "forwarded" {
		data := &types.OptionForwarded{
			Comment: comment,
		}
		if err := parseOptionForwardedAttributes(data, parts[2:]); err != nil {
			return "", err
		}
		if err := validateOptionForwarded(data); err != nil {
			return "", err
		}
		s.data = data
		return "", nil
	}
	return "", &errors.ParseError{Parser: "option forwarded", Line: line}
}

func parseOptionForwardedAttributes(data *types.OptionForwarded, parts []string) error {
	families := map[OptionForwardedAttributeName]bool{}
	for index := 0; index < len(parts); index++ {
		attribute, ok := optionForwardedAttributeByName(OptionForwardedAttributeName(parts[index]))
		if !ok {
			return errors.ErrInvalidData
		}
		if alreadyParsedOptionForwardedFamily(families, attribute.family) {
			return errors.ErrInvalidData
		}
		value := ""
		if attribute.needsValue {
			index++
			if index == len(parts) {
				return errors.ErrInvalidData
			}
			value = parts[index]
		}
		attribute.set(data, value)
	}
	return nil
}

func alreadyParsedOptionForwardedFamily(families map[OptionForwardedAttributeName]bool, family OptionForwardedAttributeName) bool {
	if families[family] {
		return true
	}
	families[family] = true
	return false
}

func validateOptionForwarded(data *types.OptionForwarded) error {
	if data == nil {
		return errors.ErrInvalidData
	}
	if data.NoOption {
		if data.Proto || data.Host || data.HostExpr != "" || data.By || data.ByExpr != "" ||
			data.ByPort || data.ByPortExpr != "" || data.For || data.ForExpr != "" ||
			data.ForPort || data.ForPortExpr != "" {
			return errors.ErrInvalidData
		}
		return nil
	}
	if data.Host && data.HostExpr != "" {
		return errors.ErrInvalidData
	}
	if data.By && data.ByExpr != "" {
		return errors.ErrInvalidData
	}
	if data.ByPort && data.ByPortExpr != "" {
		return errors.ErrInvalidData
	}
	if data.For && data.ForExpr != "" {
		return errors.ErrInvalidData
	}
	if data.ForPort && data.ForPortExpr != "" {
		return errors.ErrInvalidData
	}
	return nil
}

func (s *OptionForwarded) Result() ([]common.ReturnResultLine, error) {
	if s.data == nil {
		return nil, errors.ErrFetch
	}
	if err := validateOptionForwarded(s.data); err != nil {
		return nil, err
	}
	var sb strings.Builder
	if s.data.NoOption {
		sb.WriteString("no ")
	}
	sb.WriteString("option forwarded")
	if !s.data.NoOption {
		writeOptionForwardedAttributes(&sb, s.data)
	}
	return []common.ReturnResultLine{
		{
			Data:    sb.String(),
			Comment: s.data.Comment,
		},
	}, nil
}

func writeOptionForwardedAttributes(sb *strings.Builder, data *types.OptionForwarded) {
	if data.Proto {
		writeOptionForwardedAttribute(sb, optionForwardedProto)
	}
	if data.Host {
		writeOptionForwardedAttribute(sb, optionForwardedHost)
	}
	if data.HostExpr != "" {
		writeOptionForwardedAttributeValue(sb, optionForwardedHostExpr, data.HostExpr)
	}
	if data.By {
		writeOptionForwardedAttribute(sb, optionForwardedBy)
	}
	if data.ByExpr != "" {
		writeOptionForwardedAttributeValue(sb, optionForwardedByExpr, data.ByExpr)
	}
	if data.ByPort {
		writeOptionForwardedAttribute(sb, optionForwardedByPort)
	}
	if data.ByPortExpr != "" {
		writeOptionForwardedAttributeValue(sb, optionForwardedByPortExp, data.ByPortExpr)
	}
	if data.For {
		writeOptionForwardedAttribute(sb, optionForwardedFor)
	}
	if data.ForExpr != "" {
		writeOptionForwardedAttributeValue(sb, optionForwardedForExpr, data.ForExpr)
	}
	if data.ForPort {
		writeOptionForwardedAttribute(sb, optionForwardedForPort)
	}
	if data.ForPortExpr != "" {
		writeOptionForwardedAttributeValue(sb, optionForwardedForPortEx, data.ForPortExpr)
	}
}

func writeOptionForwardedAttribute(sb *strings.Builder, name OptionForwardedAttributeName) {
	sb.WriteString(" ")
	sb.WriteString(string(name))
}

func writeOptionForwardedAttributeValue(sb *strings.Builder, name OptionForwardedAttributeName, value string) {
	writeOptionForwardedAttribute(sb, name)
	sb.WriteString(" ")
	sb.WriteString(value)
}
