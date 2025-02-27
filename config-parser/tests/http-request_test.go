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

package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/haproxytech/client-native/v6/config-parser/parsers/http"
)

func TestHTTPRequestSetPath(t *testing.T) {
	parser := &http.Requests{}

	line := strings.TrimSpace("http-request set-path /%[hdr(host)]%[path]")

	err := ProcessLine(line, parser)
	if err != nil {
		t.Error(err)
	}

	result, err := parser.Result()
	if err != nil {
		t.Error(err)
	}

	var returnLine string

	if result[0].Comment == "" {
		returnLine = result[0].Data
	} else {
		returnLine = fmt.Sprintf("%s # %s", result[0].Data, result[0].Comment)
	}

	if line != returnLine {
		t.Errorf("error: has [%s] expects [%s]", returnLine, line)
	}
}

func TestHTTPRequestSetPathFail(t *testing.T) {
	parser := &http.Requests{}

	line := strings.TrimSpace("--- ---")

	err := ProcessLine(line, parser)

	if err == nil {
		t.Errorf("error: did not throw error for line [%s]", line)
	}

	_, err = parser.Result()
	if err == nil {
		t.Errorf("error: did not throw error on result for line [%s]", line)
	}
}

func TestHTTPRequestConnectionTrackSc0(t *testing.T) {
	parser := &http.Requests{}

	line := strings.TrimSpace("http-request track-sc0 src")

	err := ProcessLine(line, parser)
	if err != nil {
		t.Error(err)
	}

	result, err := parser.Result()
	if err != nil {
		t.Error(err)
	}

	var returnLine string

	if result[0].Comment == "" {
		returnLine = result[0].Data
	} else {
		returnLine = fmt.Sprintf("%s # %s", result[0].Data, result[0].Comment)
	}

	if line != returnLine {
		t.Errorf("error: has [%s] expects [%s]", returnLine, line)
	}
}

func TestHTTPRequestConnectionTrackSc0WithCondition(t *testing.T) {
	parser := &http.Requests{}

	line := strings.TrimSpace("http-request track-sc0 src if some_check")

	err := ProcessLine(line, parser)
	if err != nil {
		t.Error(err)
	}

	result, err := parser.Result()
	if err != nil {
		t.Error(err)
	}

	var returnLine string

	if result[0].Comment == "" {
		returnLine = result[0].Data
	} else {
		returnLine = fmt.Sprintf("%s # %s", result[0].Data, result[0].Comment)
	}

	if line != returnLine {
		t.Errorf("error: has [%s] expects [%s]", returnLine, line)
	}
}
