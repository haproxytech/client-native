package tests

import (
	"bytes"
	"strings"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/config-parser/parsers"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

func TestOptionForwardedReviewGrammar(t *testing.T) {
	tests := map[string]bool{
		"option forwarded proto host by by_port for for_port":                                       true,
		"option forwarded proto host-expr %[req.hdr(host)] by-expr %[src] by_port-expr %[src_port]": true,
		"option forwarded for-expr %[src] for_port-expr %[src_port]":                                true,
		"option forwarded proto-expr %[ssl_fc]":                                                     false,
		"option forwarded host host-expr %[req.hdr(host)]":                                          false,
		"option forwarded by by-expr %[src]":                                                        false,
		"option forwarded by_port by_port-expr %[src_port]":                                         false,
		"option forwarded for for-expr %[src]":                                                      false,
		"option forwarded for_port for_port-expr %[src_port]":                                       false,
		"option forwarded proto proto":                                                              false,
		"option forwarded host host":                                                                false,
		"option forwarded by_port by_port":                                                          false,
		"option forwarded for_port-expr":                                                            false,
		"option forwarded unknown":                                                                  false,
		"no option forwarded trailing":                                                              false,
	}
	parser := &parsers.OptionForwarded{}
	for command, shouldPass := range tests {
		t.Run(command, func(t *testing.T) {
			parser.Init()
			err := ProcessLine(strings.TrimSpace(command), parser)
			if shouldPass && err != nil {
				t.Fatalf("expected parse success, got %v", err)
			}
			if !shouldPass && err == nil {
				t.Fatal("expected parse failure")
			}
		})
	}
}

func TestOptionForwardedSetNoOption(t *testing.T) {
	parser := &parsers.OptionForwarded{}
	if err := parser.Set(types.OptionForwarded{NoOption: true}, -1); err != nil {
		t.Fatalf("set no option: %v", err)
	}
	result, err := parser.Result()
	if err != nil {
		t.Fatalf("result: %v", err)
	}
	if result[0].Data != "no option forwarded" {
		t.Fatalf("expected no option forwarded, got %q", result[0].Data)
	}
}

func TestOptionForwardedSetRejectsConflictingFamily(t *testing.T) {
	parser := &parsers.OptionForwarded{}
	if err := parser.Set(types.OptionForwarded{Proto: true}, -1); err != nil {
		t.Fatalf("set initial valid option forwarded: %v", err)
	}
	err := parser.Set(types.OptionForwarded{Host: true, HostExpr: "%[req.hdr(host)]"}, -1)
	if err == nil {
		t.Fatal("expected conflicting host family to be rejected")
	}
	result, resultErr := parser.Result()
	if resultErr != nil {
		t.Fatalf("result after rejected set: %v", resultErr)
	}
	if result[0].Data != "option forwarded proto" {
		t.Fatalf("invalid set replaced existing data with %q", result[0].Data)
	}
}

func TestOptionForwardedSetRejectsNoOptionWithAttributes(t *testing.T) {
	parser := &parsers.OptionForwarded{}
	err := parser.Set(types.OptionForwarded{NoOption: true, Proto: true}, -1)
	if err == nil {
		t.Fatal("expected no option with extra attributes to be rejected")
	}
	if _, resultErr := parser.Result(); resultErr == nil {
		t.Fatal("invalid no option value was stored")
	}
}

func TestOptionForwardedListenSection(t *testing.T) {
	config := "\nlisten test\n  option forwarded proto\n"
	var buffer bytes.Buffer
	buffer.WriteString(config)
	p, err := parser.New(options.UseListenSectionParsers, options.Reader(&buffer))
	if err != nil {
		t.Fatalf("parse listen config: %v", err)
	}
	result := p.String()
	if result != config {
		t.Fatalf("expected %q, got %q", config, result)
	}
}
