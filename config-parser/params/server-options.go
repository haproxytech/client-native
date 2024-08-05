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
	"strings"
)

// ServerOption ...
type ServerOption interface {
	Parse(options []string, currentIndex int) (int, error)
	Valid() bool
	String() string
}

// ServerOptionWord ...
type ServerOptionWord struct {
	Name string
}

// Parse ...
func (b *ServerOptionWord) Parse(options []string, currentIndex int) (int, error) {
	if currentIndex < len(options) {
		if options[currentIndex] == b.Name {
			return 1, nil
		}
		return 0, &NotFoundError{Have: options[currentIndex], Want: b.Name}
	}
	return 0, &NotEnoughParamsError{}
}

// Valid ...
func (b *ServerOptionWord) Valid() bool {
	return b.Name != ""
}

// String ...
func (b *ServerOptionWord) String() string {
	return b.Name
}

// ServerOptionDoubleWord ...
type ServerOptionDoubleWord struct {
	Name  string
	Value string
}

// Parse ...
func (b *ServerOptionDoubleWord) Parse(options []string, currentIndex int) (int, error) {
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
func (b *ServerOptionDoubleWord) Valid() bool {
	return b.Name != "" && b.Value != ""
}

// String ...
func (b *ServerOptionDoubleWord) String() string {
	if b.Name == "" || b.Value == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", b.Name, b.Value)
}

// ServerOptionValue ...
type ServerOptionValue struct {
	Name  string
	Value string
}

// Parse ...
func (b *ServerOptionValue) Parse(options []string, currentIndex int) (int, error) {
	if currentIndex+1 < len(options) {
		if options[currentIndex] == b.Name {
			b.Value = options[currentIndex+1]
			return 2, nil
		}
		return 0, &NotFoundError{Have: options[currentIndex], Want: b.Name}
	}
	return 0, &NotEnoughParamsError{}
}

// Valid ...
func (b *ServerOptionValue) Valid() bool {
	return b.Value != ""
}

// String ...
func (b *ServerOptionValue) String() string {
	if b.Name == "" || b.Value == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", b.Name, b.Value)
}

type ServerOptionIDValue struct {
	ID    string
	Name  string
	Value string
}

// Parse ...
func (b *ServerOptionIDValue) Parse(options []string, currentIndex int) (int, error) {
	if currentIndex+1 < len(options) {
		if strings.HasPrefix(options[currentIndex], fmt.Sprintf("%s(", b.Name)) {
			words := strings.Split(options[currentIndex], "(")
			if len(words) != 2 {
				return 0, &NotEnoughParamsError{}
			}
			if !strings.HasSuffix(words[1], ")") {
				return 0, &NotEnoughParamsError{}
			}
			b.ID = words[1][:len(words[1])-1]
			b.Value = options[currentIndex+1]
			return 2, nil
		}
		return 0, &NotFoundError{Have: options[currentIndex], Want: b.Name}
	}
	return 0, &NotEnoughParamsError{}
}

// Valid ...
func (b *ServerOptionIDValue) Valid() bool {
	return b.Value != "" && b.ID != ""
}

// String ...
func (b *ServerOptionIDValue) String() string {
	if b.Name == "" || b.ID == "" || b.Value == "" {
		return ""
	}
	return fmt.Sprintf("%s(%s) %s", b.Name, b.ID, b.Value)
}

func getServerOptions() []ServerOption {
	return []ServerOption{
		&ServerOptionWord{Name: "agent-check"},
		&ServerOptionWord{Name: "allow-0rtt"},
		&ServerOptionWord{Name: "backup"},
		&ServerOptionWord{Name: "check"},
		&ServerOptionWord{Name: "check-send-proxy"},
		&ServerOptionWord{Name: "check-ssl"},
		&ServerOptionWord{Name: "check-via-socks4"},
		&ServerOptionWord{Name: "disabled"},
		&ServerOptionWord{Name: "enabled"},
		&ServerOptionWord{Name: "force-sslv3"},
		&ServerOptionWord{Name: "force-tlsv10"},
		&ServerOptionWord{Name: "force-tlsv11"},
		&ServerOptionWord{Name: "force-tlsv12"},
		&ServerOptionWord{Name: "force-tlsv13"},
		&ServerOptionWord{Name: "no-agent-check"},
		&ServerOptionWord{Name: "no-backup"},
		&ServerOptionWord{Name: "no-check"},
		&ServerOptionWord{Name: "no-check-ssl"},
		&ServerOptionWord{Name: "no-send-proxy"},
		&ServerOptionWord{Name: "no-send-proxy-v2"},
		&ServerOptionWord{Name: "no-send-proxy-v2-ssl"},
		&ServerOptionWord{Name: "no-send-proxy-v2-ssl-cn"},
		&ServerOptionWord{Name: "no-ssl"},
		&ServerOptionWord{Name: "no-ssl-reuse"},
		&ServerOptionWord{Name: "no-sslv3"},
		&ServerOptionWord{Name: "no-tls-tickets"},
		&ServerOptionWord{Name: "no-tlsv10"},
		&ServerOptionWord{Name: "no-tlsv11"},
		&ServerOptionWord{Name: "no-tlsv12"},
		&ServerOptionWord{Name: "no-tlsv13"},
		&ServerOptionWord{Name: "no-verifyhost"},
		&ServerOptionWord{Name: "no-tfo"},
		&ServerOptionWord{Name: "non-stick"},
		&ServerOptionWord{Name: "send-proxy"},
		&ServerOptionWord{Name: "send-proxy-v2"},
		&ServerOptionWord{Name: "send-proxy-v2-ssl"},
		&ServerOptionWord{Name: "send-proxy-v2-ssl-cn"},
		&ServerOptionWord{Name: "ssl"},
		&ServerOptionWord{Name: "ssl-reuse"},
		&ServerOptionWord{Name: "stick"},
		&ServerOptionWord{Name: "tfo"},
		&ServerOptionWord{Name: "tls-tickets"},
		&ServerOptionValue{Name: "addr"},
		&ServerOptionValue{Name: "agent-send"},
		&ServerOptionValue{Name: "agent-inter"},
		&ServerOptionValue{Name: "agent-addr"},
		&ServerOptionValue{Name: "agent-port"},
		&ServerOptionValue{Name: "alpn"},
		&ServerOptionValue{Name: "ca-file"},
		&ServerOptionValue{Name: "check-alpn"},
		&ServerOptionValue{Name: "check-proto"},
		&ServerOptionValue{Name: "check-sni"},
		&ServerOptionValue{Name: "ciphers"},
		&ServerOptionValue{Name: "ciphersuites"},
		&ServerOptionValue{Name: "client-sigalgs"},
		&ServerOptionValue{Name: "cookie"},
		&ServerOptionValue{Name: "crl-file"},
		&ServerOptionValue{Name: "crt"},
		&ServerOptionValue{Name: "curves"},
		&ServerOptionValue{Name: "error-limit"},
		&ServerOptionValue{Name: "fall"},
		&ServerOptionValue{Name: "id"},
		&ServerOptionValue{Name: "init-addr"},
		&ServerOptionValue{Name: "inter"},
		&ServerOptionValue{Name: "fastinter"},
		&ServerOptionValue{Name: "downinter"},
		&ServerOptionValue{Name: "log-proto"},
		&ServerOptionValue{Name: "maxconn"},
		&ServerOptionValue{Name: "maxqueue"},
		&ServerOptionValue{Name: "max-reuse"},
		&ServerOptionValue{Name: "minconn"},
		&ServerOptionValue{Name: "namespace"},
		&ServerOptionValue{Name: "npn"},
		&ServerOptionValue{Name: "observe"},
		&ServerOptionValue{Name: "on-error"},
		&ServerOptionValue{Name: "on-marked-down"},
		&ServerOptionValue{Name: "on-marked-up"},
		&ServerOptionValue{Name: "pool-max-conn"},
		&ServerOptionValue{Name: "pool-purge-delay"},
		&ServerOptionValue{Name: "port"},
		&ServerOptionValue{Name: "proto"},
		&ServerOptionValue{Name: "redir"},
		&ServerOptionValue{Name: "rise"},
		&ServerOptionValue{Name: "resolve-opts"},
		&ServerOptionValue{Name: "resolve-prefer"},
		&ServerOptionValue{Name: "resolve-net"},
		&ServerOptionValue{Name: "resolvers"},
		&ServerOptionValue{Name: "proxy-v2-options"},
		&ServerOptionValue{Name: "shard"},
		&ServerOptionValue{Name: "sigalgs"},
		&ServerOptionValue{Name: "slowstart"},
		&ServerOptionValue{Name: "sni"},
		&ServerOptionValue{Name: "source"},
		&ServerOptionValue{Name: "usesrc"},
		&ServerOptionValue{Name: "interface"},
		&ServerOptionValue{Name: "ssl-max-ver"},
		&ServerOptionValue{Name: "ssl-min-ver"},
		&ServerOptionValue{Name: "socks4"},
		&ServerOptionValue{Name: "tcp-ut"},
		&ServerOptionValue{Name: "track"},
		&ServerOptionValue{Name: "verify"},
		&ServerOptionValue{Name: "verifyhost"},
		&ServerOptionValue{Name: "weight"},
		&ServerOptionValue{Name: "pool-low-conn"},
		&ServerOptionValue{Name: "ws"},
		&ServerOptionValue{Name: "log-bufsize"},
		&ServerOptionValue{Name: "guid"},
		&ServerOptionIDValue{Name: "set-proxy-v2-tlv-fmt"},
		&ServerOptionValue{Name: "pool-conn-name"},
		&ServerOptionValue{Name: "hash-key"},
	}
}

// Parse ...
func ParseServerOptions(options []string) []ServerOption {
	result := []ServerOption{}
	currentIndex := 0
	for currentIndex < len(options) {
		found := false
		for _, parser := range getServerOptions() {
			if size, err := parser.Parse(options, currentIndex); err == nil {
				result = append(result, parser)
				found = true
				currentIndex += size
			}
		}
		if !found {
			currentIndex++
		}
	}
	return result
}

func ServerOptionsString(options []ServerOption) string {
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

// CreateServerOptionWord creates valid one word value
func CreateServerOptionWord(name string) (ServerOptionWord, ErrParseServerOption) {
	b := ServerOptionWord{
		Name: name,
	}
	_, err := b.Parse([]string{name}, 0)
	return b, err
}

// CreateServerOptionDoubleWord creates valid two word value
func CreateServerOptionDoubleWord(name1 string, name2 string) (ServerOptionDoubleWord, ErrParseServerOption) {
	b := ServerOptionDoubleWord{
		Name:  name1,
		Value: name2,
	}
	_, err := b.Parse([]string{name1, name2}, 0)
	return b, err
}

// CreateServerOptionValue creates valid option with value
func CreateServerOptionValue(name string, value string) (ServerOptionValue, ErrParseServerOption) {
	b := ServerOptionValue{
		Name:  name,
		Value: value,
	}
	_, err := b.Parse([]string{name, value}, 0)
	return b, err
}
