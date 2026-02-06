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

package params

import (
	"fmt"
	"slices"
	"strings"
)

// BindOption ...
type BindOption interface { //nolint:iface
	Parse(options []string, currentIndex int) (int, error)
	Valid() bool
	String() string
}

// BindOptionWord ...
type BindOptionWord struct {
	Name string
}

// Parse ...
func (b *BindOptionWord) Parse(_ []string, _ int) (int, error) {
	return 1, nil
}

// Valid ...
func (b *BindOptionWord) Valid() bool {
	return b.Name != ""
}

// String ...
func (b *BindOptionWord) String() string {
	return b.Name
}

// BindOptionDoubleWord ...
type BindOptionDoubleWord struct {
	Name  string
	Value string
}

// Parse ...
func (b *BindOptionDoubleWord) Parse(options []string, currentIndex int) (int, error) {
	if currentIndex+1 < len(options) {
		if options[currentIndex] == b.Name && b.Value == options[currentIndex+1] {
			return 2, nil
		}
		return 0, &NotFoundError{
			Have: fmt.Sprintf("%s %s", options[currentIndex], options[currentIndex]),
			Want: fmt.Sprintf("%s %s", b.Name, b.Value),
		}
	}
	return 0, &NotEnoughParamsError{}
}

// Valid ...
func (b *BindOptionDoubleWord) Valid() bool {
	return b.Name != "" && b.Value != ""
}

// String ...
func (b *BindOptionDoubleWord) String() string {
	if b.Name == "" || b.Value == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", b.Name, b.Value)
}

// BindOptionValue ...
type BindOptionValue struct {
	Name  string
	Value string
}

// BindOptionValueValidation ...
type BindOptionValueValidation struct {
	AllowedValues []string
}

//nolint:gochecknoglobals
var bindOptionValuesValidation = map[string]BindOptionValueValidation{
	"quic-socket": {
		AllowedValues: []string{"connection", "listener"},
	},
}

// Parse ...
func (b *BindOptionValue) Parse(options []string, currentIndex int) (int, error) {
	if currentIndex+1 < len(options) {
		if options[currentIndex] == b.Name {
			b.Value = options[currentIndex+1]
			if optionValuesValidation, ok := bindOptionValuesValidation[options[currentIndex]]; ok {
				if !slices.Contains(optionValuesValidation.AllowedValues, b.Value) {
					return 0, &NotAllowedValuesError{
						Have: options[currentIndex+1],
						Want: optionValuesValidation.AllowedValues,
					}
				}
			}
			return 2, nil
		}
		return 0, &NotFoundError{Have: options[currentIndex], Want: b.Name}
	}
	return 0, &NotEnoughParamsError{}
}

// Valid ...
func (b *BindOptionValue) Valid() bool {
	return b.Value != ""
}

// String ...
func (b *BindOptionValue) String() string {
	if b.Name == "" || b.Value == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", b.Name, b.Value)
}

// BindOptionValue ...
type BindOptionParams struct {
	Name   string
	Value  string
	Params []string
}

// BindOptionParamsValidation ...
type BindOptionParamsValidation struct {
	AllowedValues []string
}

//nolint:gochecknoglobals
var bindOptionParamsValidation = map[string]BindOptionParamsValidation{
	"quic-cc-algo": {
		AllowedValues: []string{"cubic", "newreno", "bbr", "nocc"},
	},
}

// Parse ...
func (b *BindOptionParams) Parse(options []string, currentIndex int) (int, error) {
	if currentIndex+1 < len(options) {
		name := options[currentIndex]
		if name == b.Name {
			value, params, err := b.splitParams(options[currentIndex+1])
			if err != nil {
				return 0, err
			}
			b.Value = value
			b.Params = params
			if err = b.validateAllowedValues(); err != nil {
				return 0, err
			}
			return 2, nil
		}
		return 0, &NotFoundError{Have: name, Want: b.Name}
	}
	return 0, &NotEnoughParamsError{}
}

func (b *BindOptionParams) validateAllowedValues() error {
	if optionValuesValidation, ok := bindOptionParamsValidation[b.Name]; ok {
		if !slices.Contains(optionValuesValidation.AllowedValues, b.Value) {
			return &NotAllowedValuesError{
				Have: b.Value,
				Want: optionValuesValidation.AllowedValues,
			}
		}
	}
	return nil
}

func (b *BindOptionParams) splitParams(value string) (string, []string, error) {
	if !strings.HasSuffix(value, ")") {
		return value, nil, nil
	}

	value = strings.TrimSuffix(value, ")")
	parts := strings.Split(value, "(")
	if len(parts) != 2 {
		return "", nil, &NotEnoughParamsError{}
	}

	if len(parts[1]) == 0 {
		return "", nil, &NotEnoughParamsError{}
	}

	params := strings.Split(parts[1], ",")
	return parts[0], params, nil
}

// Valid ...
func (b *BindOptionParams) Valid() bool {
	return b.Value != ""
}

// String ...
func (b *BindOptionParams) String() string {
	if b.Name == "" || b.Value == "" {
		return ""
	}

	var result strings.Builder
	result.WriteString(b.Name)
	result.WriteString(" ")
	result.WriteString(b.Value)
	if len(b.Params) > 0 {
		result.WriteString("(")
		result.WriteString(strings.Join(b.Params, ","))
		result.WriteString(")")
	}

	return result.String()
}

func getBindOption(option string) BindOption {
	if factoryMethod, found := bindOptionFactoryMethods[option]; found {
		return factoryMethod()
	}
	return nil
}

var bindOptionFactoryMethods = map[string]func() BindOption{ //nolint:gochecknoglobals
	"accept-proxy":          func() BindOption { return &BindOptionWord{Name: "accept-proxy"} },
	"allow-0rtt":            func() BindOption { return &BindOptionWord{Name: "allow-0rtt"} },
	"defer-accept":          func() BindOption { return &BindOptionWord{Name: "defer-accept"} },
	"force-sslv3":           func() BindOption { return &BindOptionWord{Name: "force-sslv3"} },
	"force-tlsv10":          func() BindOption { return &BindOptionWord{Name: "force-tlsv10"} },
	"force-tlsv11":          func() BindOption { return &BindOptionWord{Name: "force-tlsv11"} },
	"force-tlsv12":          func() BindOption { return &BindOptionWord{Name: "force-tlsv12"} },
	"force-tlsv13":          func() BindOption { return &BindOptionWord{Name: "force-tlsv13"} },
	"generate-certificates": func() BindOption { return &BindOptionWord{Name: "generate-certificates"} },
	"no-alpn":               func() BindOption { return &BindOptionWord{Name: "no-alpn"} },
	"no-ca-names":           func() BindOption { return &BindOptionWord{Name: "no-ca-names"} },
	"no-sslv3":              func() BindOption { return &BindOptionWord{Name: "no-sslv3"} },
	"tls-tickets":           func() BindOption { return &BindOptionWord{Name: "tls-tickets"} },
	"no-tls-tickets":        func() BindOption { return &BindOptionWord{Name: "no-tls-tickets"} },
	"no-tlsv10":             func() BindOption { return &BindOptionWord{Name: "no-tlsv10"} },
	"no-tlsv11":             func() BindOption { return &BindOptionWord{Name: "no-tlsv11"} },
	"no-tlsv12":             func() BindOption { return &BindOptionWord{Name: "no-tlsv12"} },
	"no-tlsv13":             func() BindOption { return &BindOptionWord{Name: "no-tlsv13"} },
	"prefer-client-ciphers": func() BindOption { return &BindOptionWord{Name: "prefer-client-ciphers"} },
	"ssl":                   func() BindOption { return &BindOptionWord{Name: "ssl"} },
	"strict-sni":            func() BindOption { return &BindOptionWord{Name: "strict-sni"} },
	"no-strict-sni":         func() BindOption { return &BindOptionWord{Name: "no-strict-sni"} },
	"tfo":                   func() BindOption { return &BindOptionWord{Name: "tfo"} },
	"transparent":           func() BindOption { return &BindOptionWord{Name: "transparent"} },
	"v4v6":                  func() BindOption { return &BindOptionWord{Name: "v4v6"} },
	"v6only":                func() BindOption { return &BindOptionWord{Name: "v6only"} },
	"quic-force-retry":      func() BindOption { return &BindOptionWord{Name: "quic-force-retry"} },

	"expose-fd": func() BindOption { return &BindOptionDoubleWord{Name: "expose-fd", Value: "listeners"} },

	"accept-netscaler-cip": func() BindOption { return &BindOptionValue{Name: "accept-netscaler-cip"} },
	"alpn":                 func() BindOption { return &BindOptionValue{Name: "alpn"} },
	"backlog":              func() BindOption { return &BindOptionValue{Name: "backlog"} },
	"curves":               func() BindOption { return &BindOptionValue{Name: "curves"} },
	"ecdhe":                func() BindOption { return &BindOptionValue{Name: "ecdhe"} },
	"ca-file":              func() BindOption { return &BindOptionValue{Name: "ca-file"} },
	"ca-ignore-err":        func() BindOption { return &BindOptionValue{Name: "ca-ignore-err"} },
	"ca-sign-file":         func() BindOption { return &BindOptionValue{Name: "ca-sign-file"} },
	"ca-sign-pass":         func() BindOption { return &BindOptionValue{Name: "ca-sign-pass"} },
	"ca-verify-file":       func() BindOption { return &BindOptionValue{Name: "ca-verify-file"} },
	"cc":                   func() BindOption { return &BindOptionValue{Name: "cc"} },
	"ciphers":              func() BindOption { return &BindOptionValue{Name: "ciphers"} },
	"ciphersuites":         func() BindOption { return &BindOptionValue{Name: "ciphersuites"} },
	"client-sigalgs":       func() BindOption { return &BindOptionValue{Name: "client-sigalgs"} },
	"crl-file":             func() BindOption { return &BindOptionValue{Name: "crl-file"} },
	"crt":                  func() BindOption { return &BindOptionValue{Name: "crt"} },
	"crt-ignore-err":       func() BindOption { return &BindOptionValue{Name: "crt-ignore-err"} },
	"crt-list":             func() BindOption { return &BindOptionValue{Name: "crt-list"} },
	"gid":                  func() BindOption { return &BindOptionValue{Name: "gid"} },
	"group":                func() BindOption { return &BindOptionValue{Name: "group"} },
	"id":                   func() BindOption { return &BindOptionValue{Name: "id"} },
	"idle-ping":            func() BindOption { return &BindOptionValue{Name: "idle-ping"} },
	"interface":            func() BindOption { return &BindOptionValue{Name: "interface"} },
	"label":                func() BindOption { return &BindOptionValue{Name: "label"} },
	"level":                func() BindOption { return &BindOptionValue{Name: "level"} },
	"severity-output":      func() BindOption { return &BindOptionValue{Name: "severity-output"} },
	"maxconn":              func() BindOption { return &BindOptionValue{Name: "maxconn"} },
	"mode":                 func() BindOption { return &BindOptionValue{Name: "mode"} },
	"mss":                  func() BindOption { return &BindOptionValue{Name: "mss"} },
	"name":                 func() BindOption { return &BindOptionValue{Name: "name"} },
	"namespace":            func() BindOption { return &BindOptionValue{Name: "namespace"} },
	"nice":                 func() BindOption { return &BindOptionValue{Name: "nice"} },
	"npn":                  func() BindOption { return &BindOptionValue{Name: "npn"} },
	"ocsp-update":          func() BindOption { return &BindOptionValue{Name: "ocsp-update"} },
	"process":              func() BindOption { return &BindOptionValue{Name: "process"} },
	"proto":                func() BindOption { return &BindOptionValue{Name: "proto"} },
	"sigalgs":              func() BindOption { return &BindOptionValue{Name: "sigalgs"} },
	"ssl-max-ver":          func() BindOption { return &BindOptionValue{Name: "ssl-max-ver"} },
	"ssl-min-ver":          func() BindOption { return &BindOptionValue{Name: "ssl-min-ver"} },
	"tcp-md5sig":           func() BindOption { return &BindOptionValue{Name: "tcp-md5sig"} },
	"tcp-ss":               func() BindOption { return &BindOptionValue{Name: "tcp-ss"} },
	"tcp-ut":               func() BindOption { return &BindOptionValue{Name: "tcp-ut"} },
	"thread":               func() BindOption { return &BindOptionValue{Name: "thread"} },
	"tls-ticket-keys":      func() BindOption { return &BindOptionValue{Name: "tls-ticket-keys"} },
	"uid":                  func() BindOption { return &BindOptionValue{Name: "uid"} },
	"user":                 func() BindOption { return &BindOptionValue{Name: "user"} },
	"verify":               func() BindOption { return &BindOptionValue{Name: "verify"} },
	"quic-socket":          func() BindOption { return &BindOptionValue{Name: "quic-socket"} },
	"nbconn":               func() BindOption { return &BindOptionValue{Name: "nbconn"} },
	"guid-prefix":          func() BindOption { return &BindOptionValue{Name: "guid-prefix"} },
	"default-crt":          func() BindOption { return &BindOptionValue{Name: "default-crt"} },

	"quic-cc-algo": func() BindOption { return &BindOptionParams{Name: "quic-cc-algo"} },
	"ktls":         func() BindOption { return &BindOptionOnOff{Name: "ktls"} },
}

// Parse ...
func ParseBindOptions(options []string) ([]BindOption, error) {
	result := []BindOption{}
	currentIndex := 0
	for currentIndex < len(options) {
		bindOption := getBindOption(options[currentIndex])
		if bindOption == nil {
			return nil, &NotFoundError{Have: options[currentIndex]}
		}
		size, err := bindOption.Parse(options, currentIndex)
		if err != nil {
			return nil, err
		}
		result = append(result, bindOption)
		currentIndex += size
	}
	return result, nil
}

func BindOptionsString(options []BindOption) string {
	var sb strings.Builder
	first := true
	for _, parser := range options {
		if parser.Valid() {
			if !first {
				sb.WriteString(" ")
			} else {
				first = false
			}
			sb.WriteString(parser.String())
		}
	}
	return sb.String()
}

// CreateBindOptionWord creates valid one word value
func CreateBindOptionWord(name string) (BindOptionWord, ErrParseBindOption) {
	b := BindOptionWord{
		Name: name,
	}
	_, err := b.Parse([]string{name}, 0)
	return b, err
}

// CreateBindOptionDoubleWord creates valid two word value
func CreateBindOptionDoubleWord(name1 string, name2 string) (BindOptionDoubleWord, ErrParseBindOption) {
	b := BindOptionDoubleWord{
		Name:  name1,
		Value: name2,
	}
	_, err := b.Parse([]string{name1, name2}, 0)
	return b, err
}

// CreateBindOptionValue creates valid option with value
func CreateBindOptionValue(name string, value string) (BindOptionValue, ErrParseBindOption) {
	b := BindOptionValue{
		Name:  name,
		Value: value,
	}
	_, err := b.Parse([]string{name, value}, 0)
	return b, err
}
