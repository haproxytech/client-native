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
package configs //nolint:testpackage

import (
	"errors"
	"slices"
	"strings"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parserErrors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

func newHTTPHardeningParser(t *testing.T, config string) parser.Parser {
	t.Helper()
	p, err := parser.New()
	if err != nil {
		t.Fatal(err)
	}
	if err = p.Process(strings.NewReader(config)); err != nil {
		t.Fatal(err)
	}
	return p
}

func httpHardeningOption(t *testing.T, p parser.Parser, name, keyword string) string {
	t.Helper()
	data, err := p.Get(parser.HTTPHardening, name, keyword, false)
	if err != nil {
		t.Fatal(err)
	}
	option, ok := data.(*types.StringC)
	if !ok {
		t.Fatalf("unexpected data type %T", data)
	}
	return option.Value
}

func TestHTTPHardeningOptionsRoundTrip(t *testing.T) {
	keywords := []string{
		"all",
		"all-reject",
		"all-delete",
		"delete-all-trailer-fields",
		"delete-underscore-in-header-name",
		"delete-underscore-in-trailer-name",
		"delete-unregistered-header-fields",
		"delete-unregistered-trailer-fields",
		"reject-content-length-leading-zero",
		"reject-content-length-with-te",
		"reject-ctl-in-header-trailer",
		"reject-duplicate-host-header",
		"reject-header-trailer-folding",
		"reject-lowercase-method",
		"reject-underscore-in-header-name",
		"reject-underscore-in-trailer-name",
	}
	for _, keyword := range keywords {
		for _, value := range []string{"on", "off"} {
			config := "http-hardening p\n  " + keyword + " " + value + "\n"
			p := newHTTPHardeningParser(t, config)
			if got := httpHardeningOption(t, p, "p", keyword); got != value {
				t.Fatalf("%s: expected %q, got %q", keyword, value, got)
			}
			if out := p.String(); !strings.Contains(out, keyword+" "+value) {
				t.Fatalf("%s: missing %q in output:\n%s", keyword, keyword+" "+value, out)
			}
		}
	}
}

func TestHTTPHardeningProfileOrder(t *testing.T) {
	config := `http-hardening z-base
  all on
  reject-lowercase-method off

http-hardening a-child from z-base
  all-delete off
  reject-lowercase-method on
`
	p := newHTTPHardeningParser(t, config)

	names, err := p.SectionsGet(parser.HTTPHardening)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(names, "z-base") || !slices.Contains(names, "a-child") {
		t.Fatalf("missing profiles: %v", names)
	}

	out := p.String()
	base := strings.Index(out, "http-hardening z-base")
	child := strings.Index(out, "http-hardening a-child from z-base")
	if base < 0 || child < 0 || base > child {
		t.Fatalf("parent was written after child:\n%s", out)
	}

	reparsed := newHTTPHardeningParser(t, out)
	reparsedOut := reparsed.String()
	reparsedBase := strings.Index(reparsedOut, "http-hardening z-base")
	reparsedChild := strings.Index(reparsedOut, "http-hardening a-child from z-base")
	if reparsedBase < 0 || reparsedChild < 0 || reparsedBase > reparsedChild {
		t.Fatalf("parent order changed after round trip:\n%s", reparsedOut)
	}
}

func TestHTTPHardeningLastWins(t *testing.T) {
	config := `http-hardening p
  all on
  reject-lowercase-method off
  all-delete off
  reject-lowercase-method on
`
	p := newHTTPHardeningParser(t, config)

	for keyword, expected := range map[string]string{
		"all":                     "on",
		"all-delete":              "off",
		"reject-lowercase-method": "on",
	} {
		if value := httpHardeningOption(t, p, "p", keyword); value != expected {
			t.Fatalf("%s: expected %q, got %q", keyword, expected, value)
		}
	}
	for _, keyword := range []string{"delete-all-trailer-fields", "reject-duplicate-host-header"} {
		if _, err := p.Get(parser.HTTPHardening, "p", keyword, false); !errors.Is(err, parserErrors.ErrFetch) {
			t.Fatalf("%s: expected ErrFetch for unset option, got %v", keyword, err)
		}
	}

	out := p.String()
	for _, line := range []string{"all on", "all-delete off", "reject-lowercase-method on"} {
		if !strings.Contains(out, line) {
			t.Fatalf("missing %q in output:\n%s", line, out)
		}
	}
	if strings.Count(out, "reject-lowercase-method") != 1 {
		t.Fatalf("repeated keyword was not collapsed:\n%s", out)
	}
}

func TestHTTPHardeningEmptyFieldsSkipped(t *testing.T) {
	config := `http-hardening p
  all on
`
	p := newHTTPHardeningParser(t, config)
	out := p.String()
	for _, line := range []string{"all-reject", "all-delete", "delete-", "reject-"} {
		if strings.Contains(out, line) {
			t.Fatalf("unset field %q was written:\n%s", line, out)
		}
	}
}

func TestHTTPHardeningBuiltinInheritanceOrder(t *testing.T) {
	config := `http-hardening a-child from trusted
  all on

http-hardening trusted
  all off
`
	p := newHTTPHardeningParser(t, config)

	names, err := p.SectionsGet(parser.HTTPHardening)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(names, "a-child") || !slices.Contains(names, "trusted") {
		t.Fatalf("missing profiles: %v", names)
	}

	data, err := p.SectionGet(parser.HTTPHardening, "a-child")
	if err != nil {
		t.Fatal(err)
	}
	section := data.(types.Section)
	if section.From != "trusted" {
		t.Fatalf("expected parent trusted, got %q", section.From)
	}
}

func TestHTTPHardeningDuplicateHeader(t *testing.T) {
	config := `http-hardening p
  all on

http-hardening p
  all off
`
	p := newHTTPHardeningParser(t, config)

	names, err := p.SectionsGet(parser.HTTPHardening)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(names, []string{"p"}) {
		t.Fatalf("unexpected profiles: %v", names)
	}

	if value := httpHardeningOption(t, p, "p", "all"); value != "off" {
		t.Fatalf("last profile did not win: %q", value)
	}
	out := p.String()
	if strings.Count(out, "http-hardening p") != 1 || !strings.Contains(out, "all off") {
		t.Fatalf("duplicate section not collapsed:\n%s", out)
	}
}

func TestHTTPHardeningSectionsCreateAndDelete(t *testing.T) {
	p := newHTTPHardeningParser(t, "http-hardening a\n")

	if err := p.SectionsCreate(parser.HTTPHardening, "b"); err != nil {
		t.Fatal(err)
	}
	names, err := p.SectionsGet(parser.HTTPHardening)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(names, "a") || !slices.Contains(names, "b") {
		t.Fatalf("create did not add profile: %v", names)
	}

	if err := p.SectionsCreate(parser.HTTPHardening, "a"); !errors.Is(err, parserErrors.ErrSectionAlreadyExists) {
		t.Fatalf("expected ErrSectionAlreadyExists, got %v", err)
	}

	if err := p.SectionsDelete(parser.HTTPHardening, "a"); err != nil {
		t.Fatal(err)
	}
	names, err = p.SectionsGet(parser.HTTPHardening)
	if err != nil {
		t.Fatal(err)
	}
	if slices.Contains(names, "a") || !slices.Contains(names, "b") {
		t.Fatalf("delete did not remove profile: %v", names)
	}
}

func TestHTTPHardeningDefaultsFromSet(t *testing.T) {
	p := newHTTPHardeningParser(t, "http-hardening base\n\nhttp-hardening child\n")

	if err := p.SectionsDefaultsFromSet(parser.HTTPHardening, "child", "base"); err != nil {
		t.Fatal(err)
	}
	parent, err := p.SectionsDefaultsFromGet(parser.HTTPHardening, "child")
	if err != nil {
		t.Fatal(err)
	}
	if parent != "base" {
		t.Fatalf("expected base, got %q", parent)
	}

	if err := p.SectionsDefaultsFromSet(parser.HTTPHardening, "child", "trusted"); err != nil {
		t.Fatal(err)
	}
	if err := p.SectionsDefaultsFromSet(parser.HTTPHardening, "child", "auto"); !errors.Is(err, parserErrors.ErrFromDefaultsSectionMissing) {
		t.Fatalf("expected ErrFromDefaultsSectionMissing, got %v", err)
	}
	if err := p.SectionsDefaultsFromSet(parser.HTTPHardening, "child", "child"); !errors.Is(err, parserErrors.ErrCircularDependency) {
		t.Fatalf("expected ErrCircularDependency, got %v", err)
	}
	if err := p.SectionsDefaultsFromSet(parser.HTTPHardening, "child", "missing"); !errors.Is(err, parserErrors.ErrFromDefaultsSectionMissing) {
		t.Fatalf("expected ErrFromDefaultsSectionMissing, got %v", err)
	}
	if err := p.SectionsDefaultsFromSet(parser.HTTPHardening, "child", ""); err != nil {
		t.Fatal(err)
	}
	parent, err = p.SectionsDefaultsFromGet(parser.HTTPHardening, "child")
	if err != nil {
		t.Fatal(err)
	}
	if parent != "" {
		t.Fatalf("expected cleared parent, got %q", parent)
	}
}
