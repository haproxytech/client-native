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
	"reflect"
	"testing"

	"github.com/go-openapi/strfmt"

	"github.com/haproxytech/client-native/v6/config-parser/parsers/filters"
	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
)

// TestFilterModelSchema asserts that the generated Frontend and Backend models
// no longer expose a FilterSequenceList field once the standalone
// filter-sequence support has been removed.
func TestFilterModelSchema(t *testing.T) {
	tests := []struct {
		name string
		typ  reflect.Type
	}{
		{"frontend", reflect.TypeOf(models.Frontend{})},
		{"backend", reflect.TypeOf(models.Backend{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, ok := tt.typ.FieldByName("FilterSequenceList"); ok {
				t.Fatalf("%s must not expose FilterSequenceList after filter-sequence removal", tt.name)
			}
		})
	}
}

// TestCompressionFilterConversion covers the comp-req and comp-res filters,
// which share the compression directive family and must keep working through
// model validation, serialization to the typed parser data and parsing back.
func TestCompressionFilterConversion(t *testing.T) { //nolint:gocognit
	metadata := map[string]interface{}{"key": "value"}
	wantComment := `{"key":"value"}`

	tests := []struct {
		name      string
		model     models.Filter
		wantType  reflect.Type
		wantLine  string
		wantModel string
	}{
		{
			name:      "comp-req",
			model:     models.Filter{Type: "comp-req", Metadata: metadata},
			wantType:  reflect.TypeOf(&filters.CompReq{}),
			wantLine:  "filter comp-req",
			wantModel: "comp-req",
		},
		{
			name:      "comp-res",
			model:     models.Filter{Type: "comp-res", Metadata: metadata},
			wantType:  reflect.TypeOf(&filters.CompRes{}),
			wantLine:  "filter comp-res",
			wantModel: "comp-res",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.model.Validate(strfmt.Default); err != nil {
				t.Fatalf("model validation failed: %v", err)
			}

			serialized := SerializeFilter(tt.model, &options.ConfigurationOptions{})
			if reflect.TypeOf(serialized) != tt.wantType {
				t.Fatalf("serialized type = %T, want %v", serialized, tt.wantType)
			}

			var line string
			var comment string
			switch v := serialized.(type) {
			case *filters.CompReq:
				if !v.Enabled {
					t.Fatal("comp-req Enabled must be true")
				}
				line, comment = v.Result().Data, v.Result().Comment
			case *filters.CompRes:
				if !v.Enabled {
					t.Fatal("comp-res Enabled must be true")
				}
				line, comment = v.Result().Data, v.Result().Comment
			default:
				t.Fatalf("unexpected serialized type %T", serialized)
			}

			if line != tt.wantLine {
				t.Fatalf("result directive = %q, want %q", line, tt.wantLine)
			}
			if comment != wantComment {
				t.Fatalf("result metadata comment = %q, want %q", comment, wantComment)
			}

			parsed := ParseFilter(serialized)
			if parsed == nil {
				t.Fatal("ParseFilter returned nil")
			}
			if parsed.Type != tt.wantModel {
				t.Fatalf("parsed type = %q, want %q", parsed.Type, tt.wantModel)
			}
			if !reflect.DeepEqual(parsed.Metadata, metadata) {
				t.Fatalf("parsed metadata = %v, want %v", parsed.Metadata, metadata)
			}
		})
	}
}
