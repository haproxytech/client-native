/*
Copyright 2025 HAProxy Technologies

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

import "strings"

type SSLBindOption BindOption

// Options accepted by ssl-f-use (sslbindconf).
var sslfuseOptions = map[string]func() SSLBindOption{ //nolint:gochecknoglobals
	// bind options
	"allow-0rtt":     func() SSLBindOption { return &BindOptionWord{Name: "allow-0rtt"} },
	"alpn":           func() SSLBindOption { return &BindOptionValue{Name: "alpn"} },
	"ca-file":        func() SSLBindOption { return &BindOptionValue{Name: "ca-file"} },
	"ciphers":        func() SSLBindOption { return &BindOptionValue{Name: "ciphers"} },
	"ciphersuites":   func() SSLBindOption { return &BindOptionValue{Name: "ciphersuites"} },
	"client-sigalgs": func() SSLBindOption { return &BindOptionValue{Name: "client-sigalgs"} },
	"crl-file":       func() SSLBindOption { return &BindOptionValue{Name: "crl-file"} },
	"curves":         func() SSLBindOption { return &BindOptionValue{Name: "curves"} },
	"ecdhe":          func() SSLBindOption { return &BindOptionValue{Name: "ecdhe"} },
	"no-alpn":        func() SSLBindOption { return &BindOptionWord{Name: "no-alpn"} },
	"no-ca-names":    func() SSLBindOption { return &BindOptionWord{Name: "no-ca-names"} },
	"npn":            func() SSLBindOption { return &BindOptionValue{Name: "npn"} },
	"sigalgs":        func() SSLBindOption { return &BindOptionValue{Name: "sigalgs"} },
	"ssl-max-ver":    func() SSLBindOption { return &BindOptionValue{Name: "ssl-max-ver"} },
	"ssl-min-ver":    func() SSLBindOption { return &BindOptionValue{Name: "ssl-min-ver"} },
	"verify":         func() SSLBindOption { return &BindOptionValue{Name: "verify"} },
	// crt-store load options
	"crt":         func() SSLBindOption { return &BindOptionValue{Name: "crt"} },
	"key":         func() SSLBindOption { return &BindOptionValue{Name: "key"} },
	"ocsp":        func() SSLBindOption { return &BindOptionValue{Name: "ocsp"} },
	"issuer":      func() SSLBindOption { return &BindOptionValue{Name: "issuer"} },
	"sctl":        func() SSLBindOption { return &BindOptionValue{Name: "sctl"} },
	"ocsp-update": func() SSLBindOption { return &BindOptionOnOff{Name: "ocsp-update"} },
}

func ParseSSLBindOptions(options []string) ([]SSLBindOption, error) {
	result := make([]SSLBindOption, 0, len(options))
	i := 0
	for i < len(options) {
		bindOption := getSSLBindOption(options[i])
		if bindOption == nil {
			return nil, &NotFoundError{Have: options[i]}
		}
		size, err := bindOption.Parse(options, i)
		if err != nil {
			return nil, err
		}
		result = append(result, bindOption)
		i += size
	}
	return result, nil
}

func getSSLBindOption(option string) SSLBindOption {
	if builder, found := sslfuseOptions[option]; found {
		return builder()
	}
	return nil
}

func SSLBindOptionsString(options []SSLBindOption) string {
	var sb strings.Builder
	sb.Grow(64)
	first := true
	for _, parser := range options {
		if parser.Valid() {
			if !first {
				sb.WriteByte(' ')
			} else {
				first = false
			}
			sb.WriteString(parser.String())
		}
	}
	return sb.String()
}

type BindOptionOnOff struct {
	Name  string
	Value string
}

func (b *BindOptionOnOff) Parse(options []string, currentIndex int) (int, error) {
	if currentIndex+1 < len(options) {
		if options[currentIndex] == b.Name {
			b.Value = options[currentIndex+1]
			if !b.Valid() {
				return 0, &NotAllowedValuesError{Have: b.Value, Want: []string{"on", "off"}}
			}
			return 2, nil
		}
		return 0, &NotFoundError{Have: options[currentIndex], Want: b.Name}
	}
	return 0, &NotEnoughParamsError{}
}

func (b BindOptionOnOff) Valid() bool {
	return b.Value == "on" || b.Value == "off"
}

func (b BindOptionOnOff) String() string {
	if b.Name == "" || b.Value == "" {
		return ""
	}
	return b.Name + " " + b.Value
}
