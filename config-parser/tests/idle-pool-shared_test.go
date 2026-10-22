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

package tests

import (
	"errors"
	"testing"

	configparser "github.com/haproxytech/client-native/v6/config-parser"
	parsererrors "github.com/haproxytech/client-native/v6/config-parser/errors"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

func globalParserWithLine(t *testing.T, line string) configparser.Parser {
	t.Helper()

	p, err := configparser.New(parser_options.String("global\n  " + line + "\n"))
	if err != nil {
		t.Fatalf("parser.New: %v", err)
	}
	return p
}

// tune.idle-pool.shared is the only directive accepting the full value.
func TestIdlePoolSharedAcceptsOnOffFull(t *testing.T) {
	for _, value := range []string{"on", "off", "full"} {
		t.Run(value, func(t *testing.T) {
			p := globalParserWithLine(t, "tune.idle-pool.shared "+value)
			data, err := p.Get(configparser.Global, configparser.GlobalSectionName, "tune.idle-pool.shared")
			if err != nil {
				t.Fatalf("Get: %v", err)
			}
			got, ok := data.(*types.StringC)
			if !ok || got.Value != value {
				t.Fatalf("got %#v, want %q", data, value)
			}
		})
	}
}

func TestIdlePoolSharedRejectsInvalid(t *testing.T) {
	for _, value := range []string{"enabled", "full extra", ""} {
		t.Run(value, func(t *testing.T) {
			p := globalParserWithLine(t, "tune.idle-pool.shared "+value)
			if _, err := p.Get(configparser.Global, configparser.GlobalSectionName, "tune.idle-pool.shared"); !errors.Is(err, parsererrors.ErrFetch) {
				t.Fatalf("got %v, want ErrFetch", err)
			}
		})
	}
}

// Other on/off directives must keep rejecting full.
func TestOtherOnOffRejectsFull(t *testing.T) {
	p := globalParserWithLine(t, "tune.listener.multi-queue full")
	if _, err := p.Get(configparser.Global, configparser.GlobalSectionName, "tune.listener.multi-queue"); !errors.Is(err, parsererrors.ErrFetch) {
		t.Fatalf("got %v, want ErrFetch", err)
	}
}
