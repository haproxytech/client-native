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
type ServerOption interface { //nolint:iface
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
		if strings.HasPrefix(options[currentIndex], b.Name+"(") {
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

var serverOptionFactoryMethods = map[string]func() ServerOption{ //nolint:gochecknoglobals
	"agent-check":             func() ServerOption { return &ServerOptionWord{Name: "agent-check"} },
	"no-agent-check":          func() ServerOption { return &ServerOptionWord{Name: "no-agent-check"} },
	"allow-0rtt":              func() ServerOption { return &ServerOptionWord{Name: "allow-0rtt"} },
	"backup":                  func() ServerOption { return &ServerOptionWord{Name: "backup"} },
	"no-backup":               func() ServerOption { return &ServerOptionWord{Name: "no-backup"} },
	"check":                   func() ServerOption { return &ServerOptionWord{Name: "check"} },
	"no-check":                func() ServerOption { return &ServerOptionWord{Name: "no-check"} },
	"check-send-proxy":        func() ServerOption { return &ServerOptionWord{Name: "check-send-proxy"} },
	"no-check-send-proxy":     func() ServerOption { return &ServerOptionWord{Name: "no-check-send-proxy"} },
	"check-ssl":               func() ServerOption { return &ServerOptionWord{Name: "check-ssl"} },
	"no-check-ssl":            func() ServerOption { return &ServerOptionWord{Name: "no-check-ssl"} },
	"check-pool-conn-name":    func() ServerOption { return &ServerOptionValue{Name: "check-pool-conn-name"} },
	"check-reuse-pool":        func() ServerOption { return &ServerOptionWord{Name: "check-reuse-pool"} },
	"no-check-reuse-pool":     func() ServerOption { return &ServerOptionWord{Name: "no-check-reuse-pool"} },
	"check-via-socks4":        func() ServerOption { return &ServerOptionWord{Name: "check-via-socks4"} },
	"disabled":                func() ServerOption { return &ServerOptionWord{Name: "disabled"} },
	"enabled":                 func() ServerOption { return &ServerOptionWord{Name: "enabled"} },
	"force-sslv3":             func() ServerOption { return &ServerOptionWord{Name: "force-sslv3"} },
	"no-sslv3":                func() ServerOption { return &ServerOptionWord{Name: "no-sslv3"} },
	"force-tlsv10":            func() ServerOption { return &ServerOptionWord{Name: "force-tlsv10"} },
	"no-tlsv10":               func() ServerOption { return &ServerOptionWord{Name: "no-tlsv10"} },
	"force-tlsv11":            func() ServerOption { return &ServerOptionWord{Name: "force-tlsv11"} },
	"no-tlsv11":               func() ServerOption { return &ServerOptionWord{Name: "no-tlsv11"} },
	"force-tlsv12":            func() ServerOption { return &ServerOptionWord{Name: "force-tlsv12"} },
	"no-tlsv12":               func() ServerOption { return &ServerOptionWord{Name: "no-tlsv12"} },
	"force-tlsv13":            func() ServerOption { return &ServerOptionWord{Name: "force-tlsv13"} },
	"no-tlsv13":               func() ServerOption { return &ServerOptionWord{Name: "no-tlsv13"} },
	"renegotiate":             func() ServerOption { return &ServerOptionWord{Name: "renegotiate"} },
	"no-renegotiate":          func() ServerOption { return &ServerOptionWord{Name: "no-renegotiate"} },
	"send-proxy":              func() ServerOption { return &ServerOptionWord{Name: "send-proxy"} },
	"no-send-proxy":           func() ServerOption { return &ServerOptionWord{Name: "no-send-proxy"} },
	"send-proxy-v2":           func() ServerOption { return &ServerOptionWord{Name: "send-proxy-v2"} },
	"no-send-proxy-v2":        func() ServerOption { return &ServerOptionWord{Name: "no-send-proxy-v2"} },
	"send-proxy-v2-ssl":       func() ServerOption { return &ServerOptionWord{Name: "send-proxy-v2-ssl"} },
	"no-send-proxy-v2-ssl":    func() ServerOption { return &ServerOptionWord{Name: "no-send-proxy-v2-ssl"} },
	"send-proxy-v2-ssl-cn":    func() ServerOption { return &ServerOptionWord{Name: "send-proxy-v2-ssl-cn"} },
	"no-send-proxy-v2-ssl-cn": func() ServerOption { return &ServerOptionWord{Name: "no-send-proxy-v2-ssl-cn"} },
	"ssl":                     func() ServerOption { return &ServerOptionWord{Name: "ssl"} },
	"no-ssl":                  func() ServerOption { return &ServerOptionWord{Name: "no-ssl"} },
	"ssl-reuse":               func() ServerOption { return &ServerOptionWord{Name: "ssl-reuse"} },
	"no-ssl-reuse":            func() ServerOption { return &ServerOptionWord{Name: "no-ssl-reuse"} },
	"stick":                   func() ServerOption { return &ServerOptionWord{Name: "stick"} },
	"non-stick":               func() ServerOption { return &ServerOptionWord{Name: "non-stick"} },
	"tfo":                     func() ServerOption { return &ServerOptionWord{Name: "tfo"} },
	"no-tfo":                  func() ServerOption { return &ServerOptionWord{Name: "no-tfo"} },
	"tls-tickets":             func() ServerOption { return &ServerOptionWord{Name: "tls-tickets"} },
	"no-tls-tickets":          func() ServerOption { return &ServerOptionWord{Name: "no-tls-tickets"} },
	"addr":                    func() ServerOption { return &ServerOptionValue{Name: "addr"} },
	"agent-send":              func() ServerOption { return &ServerOptionValue{Name: "agent-send"} },
	"agent-inter":             func() ServerOption { return &ServerOptionValue{Name: "agent-inter"} },
	"agent-addr":              func() ServerOption { return &ServerOptionValue{Name: "agent-addr"} },
	"agent-port":              func() ServerOption { return &ServerOptionValue{Name: "agent-port"} },
	"alpn":                    func() ServerOption { return &ServerOptionValue{Name: "alpn"} },
	"ca-file":                 func() ServerOption { return &ServerOptionValue{Name: "ca-file"} },
	"check-alpn":              func() ServerOption { return &ServerOptionValue{Name: "check-alpn"} },
	"check-proto":             func() ServerOption { return &ServerOptionValue{Name: "check-proto"} },
	"check-sni":               func() ServerOption { return &ServerOptionValue{Name: "check-sni"} },
	"ciphers":                 func() ServerOption { return &ServerOptionValue{Name: "ciphers"} },
	"ciphersuites":            func() ServerOption { return &ServerOptionValue{Name: "ciphersuites"} },
	"client-sigalgs":          func() ServerOption { return &ServerOptionValue{Name: "client-sigalgs"} },
	"cookie":                  func() ServerOption { return &ServerOptionValue{Name: "cookie"} },
	"crl-file":                func() ServerOption { return &ServerOptionValue{Name: "crl-file"} },
	"crt":                     func() ServerOption { return &ServerOptionValue{Name: "crt"} },
	"curves":                  func() ServerOption { return &ServerOptionValue{Name: "curves"} },
	"error-limit":             func() ServerOption { return &ServerOptionValue{Name: "error-limit"} },
	"fall":                    func() ServerOption { return &ServerOptionValue{Name: "fall"} },
	"id":                      func() ServerOption { return &ServerOptionValue{Name: "id"} },
	"idle-ping":               func() ServerOption { return &ServerOptionValue{Name: "idle-ping"} },
	"init-addr":               func() ServerOption { return &ServerOptionValue{Name: "init-addr"} },
	"init-state":              func() ServerOption { return &ServerOptionValue{Name: "init-state"} },
	"inter":                   func() ServerOption { return &ServerOptionValue{Name: "inter"} },
	"fastinter":               func() ServerOption { return &ServerOptionValue{Name: "fastinter"} },
	"downinter":               func() ServerOption { return &ServerOptionValue{Name: "downinter"} },
	"log-proto":               func() ServerOption { return &ServerOptionValue{Name: "log-proto"} },
	"maxconn":                 func() ServerOption { return &ServerOptionValue{Name: "maxconn"} },
	"maxqueue":                func() ServerOption { return &ServerOptionValue{Name: "maxqueue"} },
	"max-reuse":               func() ServerOption { return &ServerOptionValue{Name: "max-reuse"} },
	"minconn":                 func() ServerOption { return &ServerOptionValue{Name: "minconn"} },
	"namespace":               func() ServerOption { return &ServerOptionValue{Name: "namespace"} },
	"npn":                     func() ServerOption { return &ServerOptionValue{Name: "npn"} },
	"observe":                 func() ServerOption { return &ServerOptionValue{Name: "observe"} },
	"on-error":                func() ServerOption { return &ServerOptionValue{Name: "on-error"} },
	"on-marked-down":          func() ServerOption { return &ServerOptionValue{Name: "on-marked-down"} },
	"on-marked-up":            func() ServerOption { return &ServerOptionValue{Name: "on-marked-up"} },
	"pool-max-conn":           func() ServerOption { return &ServerOptionValue{Name: "pool-max-conn"} },
	"pool-purge-delay":        func() ServerOption { return &ServerOptionValue{Name: "pool-purge-delay"} },
	"port":                    func() ServerOption { return &ServerOptionValue{Name: "port"} },
	"proto":                   func() ServerOption { return &ServerOptionValue{Name: "proto"} },
	"redir":                   func() ServerOption { return &ServerOptionValue{Name: "redir"} },
	"rise":                    func() ServerOption { return &ServerOptionValue{Name: "rise"} },
	"resolve-opts":            func() ServerOption { return &ServerOptionValue{Name: "resolve-opts"} },
	"resolve-prefer":          func() ServerOption { return &ServerOptionValue{Name: "resolve-prefer"} },
	"resolve-net":             func() ServerOption { return &ServerOptionValue{Name: "resolve-net"} },
	"resolvers":               func() ServerOption { return &ServerOptionValue{Name: "resolvers"} },
	"proxy-v2-options":        func() ServerOption { return &ServerOptionValue{Name: "proxy-v2-options"} },
	"shard":                   func() ServerOption { return &ServerOptionValue{Name: "shard"} },
	"sigalgs":                 func() ServerOption { return &ServerOptionValue{Name: "sigalgs"} },
	"slowstart":               func() ServerOption { return &ServerOptionValue{Name: "slowstart"} },
	"sni":                     func() ServerOption { return &ServerOptionValue{Name: "sni"} },
	"sni-auto":                func() ServerOption { return &ServerOptionWord{Name: "sni-auto"} },
	"no-sni-auto":             func() ServerOption { return &ServerOptionWord{Name: "no-sni-auto"} },
	"source":                  func() ServerOption { return &ServerOptionValue{Name: "source"} },
	"strict-maxconn":          func() ServerOption { return &ServerOptionWord{Name: "strict-maxconn"} },
	"usesrc":                  func() ServerOption { return &ServerOptionValue{Name: "usesrc"} },
	"interface":               func() ServerOption { return &ServerOptionValue{Name: "interface"} },
	"ssl-max-ver":             func() ServerOption { return &ServerOptionValue{Name: "ssl-max-ver"} },
	"ssl-min-ver":             func() ServerOption { return &ServerOptionValue{Name: "ssl-min-ver"} },
	"socks4":                  func() ServerOption { return &ServerOptionValue{Name: "socks4"} },
	"tcp-md5sig":              func() ServerOption { return &ServerOptionValue{Name: "tcp-md5sig"} },
	"tcp-ut":                  func() ServerOption { return &ServerOptionValue{Name: "tcp-ut"} },
	"track":                   func() ServerOption { return &ServerOptionValue{Name: "track"} },
	"verify":                  func() ServerOption { return &ServerOptionValue{Name: "verify"} },
	"verifyhost":              func() ServerOption { return &ServerOptionValue{Name: "verifyhost"} },
	"no-verifyhost":           func() ServerOption { return &ServerOptionWord{Name: "no-verifyhost"} },
	"weight":                  func() ServerOption { return &ServerOptionValue{Name: "weight"} },
	"pool-low-conn":           func() ServerOption { return &ServerOptionValue{Name: "pool-low-conn"} },
	"ws":                      func() ServerOption { return &ServerOptionValue{Name: "ws"} },
	"log-bufsize":             func() ServerOption { return &ServerOptionValue{Name: "log-bufsize"} },
	"guid":                    func() ServerOption { return &ServerOptionValue{Name: "guid"} },
	"pool-conn-name":          func() ServerOption { return &ServerOptionValue{Name: "pool-conn-name"} },
	"hash-key":                func() ServerOption { return &ServerOptionValue{Name: "hash-key"} },
}

var serverParamOptionFactoryMethods = map[string]func() ServerOption{ //nolint:gochecknoglobals
	"set-proxy-v2-tlv-fmt": func() ServerOption { return &ServerOptionIDValue{Name: "set-proxy-v2-tlv-fmt"} },
}

func getServerOption(option string) ServerOption {
	if factoryMethod, found := serverOptionFactoryMethods[option]; found {
		return factoryMethod()
	}

	option = strings.Split(option, "(")[0]

	if factoryMethod, found := serverParamOptionFactoryMethods[option]; found {
		return factoryMethod()
	}
	return nil
}

// Parse ...
func ParseServerOptions(options []string) ([]ServerOption, error) {
	result := []ServerOption{}
	currentIndex := 0
	for currentIndex < len(options) {
		serverOption := getServerOption(options[currentIndex])
		if serverOption == nil {
			return nil, &NotFoundError{Have: options[currentIndex]}
		}
		size, err := serverOption.Parse(options, currentIndex)
		if err != nil {
			return nil, err
		}
		result = append(result, serverOption)
		currentIndex += size
	}
	return result, nil
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
