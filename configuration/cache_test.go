// Copyright 2026 HAProxy Technologies
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
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

func TestCacheEarlyHints(t *testing.T) {
	tests := []struct {
		line  string
		mode  string
		ratio *int64
	}{
		{"early-hints on", "on", nil},
		{"early-hints on ratio 25", "on", misc.Int64P(25)},
		{"early-hints off", "off", nil},
		{"early-hints off ratio 25", "off", misc.Int64P(25)},
		{"early-hints only", "only", nil},
		{"early-hints only ratio 25", "only", misc.Int64P(25)},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			p := newCacheParser(t, "cache test\n  max-age 60\n  "+tt.line+"\n")

			cache := &models.Cache{Name: misc.StringP("test")}
			if err := ParseCacheSection(p, cache); err != nil {
				t.Fatalf("ParseCacheSection: %v", err)
			}
			if cache.EarlyHints == nil {
				t.Fatal("EarlyHints is nil")
			}
			if *cache.EarlyHints.Mode != tt.mode {
				t.Errorf("mode = %q, want %q", *cache.EarlyHints.Mode, tt.mode)
			}
			if (cache.EarlyHints.Ratio == nil) != (tt.ratio == nil) {
				t.Fatalf("ratio = %v, want %v", cache.EarlyHints.Ratio, tt.ratio)
			}
			if tt.ratio != nil && *cache.EarlyHints.Ratio != *tt.ratio {
				t.Errorf("ratio = %d, want %d", *cache.EarlyHints.Ratio, *tt.ratio)
			}

			// Regression check: model must not alias parser storage.
			if cache.EarlyHints.Ratio != nil {
				*cache.EarlyHints.Ratio = 99
				assertConfigLineCount(t, p.String(), tt.line, 1)
			}

			out := newCacheParser(t, "cache test\n")
			if err := SerializeCacheSection(out, &models.Cache{
				Name:       misc.StringP("test"),
				EarlyHints: &models.CacheEarlyHints{Mode: misc.StringP(tt.mode), Ratio: tt.ratio},
				MaxAge:     60,
			}); err != nil {
				t.Fatalf("SerializeCacheSection: %v", err)
			}
			assertConfigLineCount(t, out.String(), tt.line, 1)
		})
	}

	t.Run("absent", func(t *testing.T) {
		p := newCacheParser(t, "cache test\n  max-age 60\n")

		cache := &models.Cache{Name: misc.StringP("test")}
		if err := ParseCacheSection(p, cache); err != nil {
			t.Fatalf("ParseCacheSection: %v", err)
		}
		if cache.EarlyHints != nil {
			t.Errorf("EarlyHints = %v, want nil", cache.EarlyHints)
		}
		if cache.MaxAge != 60 {
			t.Errorf("MaxAge = %d, want 60", cache.MaxAge)
		}
	})

	t.Run("removal", func(t *testing.T) {
		p := newCacheParser(t, "cache test\n  early-hints on\n  max-age 60\n")

		if err := SerializeCacheSection(p, &models.Cache{Name: misc.StringP("test"), MaxAge: 60}); err != nil {
			t.Fatalf("SerializeCacheSection: %v", err)
		}
		config := p.String()
		assertConfigLineCount(t, config, "early-hints on", 0)
		assertConfigLineCount(t, config, "max-age 60", 1)
	})

	t.Run("missing mode", func(t *testing.T) {
		p := newCacheParser(t, "cache test\n")

		err := SerializeCacheSection(p, &models.Cache{
			Name:       misc.StringP("test"),
			EarlyHints: &models.CacheEarlyHints{},
		})
		if err == nil {
			t.Fatal("SerializeCacheSection returned nil error, want validation error for missing mode")
		}
	})
}

func newCacheParser(t *testing.T, config string) parser.Parser {
	t.Helper()

	p, err := parser.New(parser_options.String(config))
	if err != nil {
		t.Fatalf("parser.New: %v", err)
	}
	return p
}
