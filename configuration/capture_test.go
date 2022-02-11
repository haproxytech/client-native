// Copyright 2019 HAProxy Technologies
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
	"io/ioutil"
	"testing"

	"github.com/haproxytech/client-native/v3/models"
)

func generateDeclareCaptureConfig(config string) (string, error) {
	f, err := ioutil.TempFile("/tmp", "capture")
	if err != nil {
		return "", err
	}
	err = prepareTestFile(config, f.Name())
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func TestDeclareCapture(t *testing.T) {
	config := `# _version=1
global
	daemon

defaults
	maxconn 2000

frontend test
	declare capture request len 1
	declare capture response len 2

frontend test_second
	declare capture request len 111
	declare capture response len 222

frontend test_replace
	declare capture request len 1
	declare capture response len 2

frontend test_add

frontend test_delete
	declare capture request len 1
	declare capture response len 2
	`
	configFile, err := generateConfig(config)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = deleteTestFile(configFile)
	}()

	tests := []struct {
		name              string
		configurationFile string
		want              int64
		wantErr           bool
	}{
		{
			name:              "declare captures",
			configurationFile: configFile,
			want:              1,
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := prepareClient(tt.configurationFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareClient error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// fetch tests, first frontend
			_, declareCaptures, err := c.GetDeclareCaptures("test", "")
			if err != nil {
				t.Error(err.Error())
			}

			if declareCaptures == nil {
				t.Errorf("No declare capture lines found, expected 2")
			}

			counter := Counter{0}

			_, declareCapture, err := c.GetDeclareCapture(0, "test", "")
			if declareCapture.Type != "request" {
				t.Errorf("Declare capture type %v returned, expected %v", declareCapture.Type, "request")
			}
			if declareCapture.Length != 1 {
				t.Errorf("Declare capture length %v returned, expected %v", declareCapture.Length, 1)
			}

			_, declareCapture, err = c.GetDeclareCapture(counter.increment(), "test", "")
			if declareCapture.Type != "response" {
				t.Errorf("Declare capture type %v returned, expected %v", declareCapture.Type, "response")
			}
			if declareCapture.Length != 2 {
				t.Errorf("Declare capture length %v returned, expected %v", declareCapture.Length, 2)
			}

			// fetch tests, second frontend
			_, declareCaptures, err = c.GetDeclareCaptures("test_second", "")
			if err != nil {
				t.Error(err.Error())
			}

			if declareCaptures == nil {
				t.Errorf("No declare capture lines found, expected 2")
			}

			counter = Counter{0}

			_, declareCapture, err = c.GetDeclareCapture(0, "test_second", "")
			if declareCapture.Type != "request" {
				t.Errorf("Declare capture type %v returned, expected %v", declareCapture.Type, "request")
			}
			if declareCapture.Length != 111 {
				t.Errorf("Declare capture length %v returned, expected %v", declareCapture.Length, 111)
			}

			_, declareCapture, err = c.GetDeclareCapture(counter.increment(), "test_second", "")
			if declareCapture.Type != "response" {
				t.Errorf("Declare capture type %v returned, expected %v", declareCapture.Type, "response")
			}
			if declareCapture.Length != 222 {
				t.Errorf("Declare capture length %v returned, expected %v", declareCapture.Length, 222)
			}

			// replace tests
			index := int64(0)
			edited := models.Capture{
				Index:  &index,
				Type:   "request",
				Length: 12345,
			}
			if c.EditDeclareCapture(0, "test_replace", &edited, "", 1) != nil {
				t.Errorf("Edit of an existing declare capture failed")
			}

			index = int64(1)
			edited = models.Capture{
				Index:  &index,
				Type:   "response",
				Length: 12345,
			}
			if c.EditDeclareCapture(1, "test_replace", &edited, "", 2) != nil {
				t.Errorf("Edit of an existing declare capture failed")
			}

			// add tests
			index = int64(0)
			add := models.Capture{
				Index:  &index,
				Type:   "request",
				Length: 1,
			}
			if c.CreateDeclareCapture("test_add", &add, "", 3) != nil {
				t.Errorf("Adding a new declare capture request failed")
			}

			add = models.Capture{
				Index:  &index,
				Type:   "response",
				Length: 2,
			}
			if c.CreateDeclareCapture("test_add", &add, "", 4) != nil {
				t.Errorf("Adding a new declare capture response failed")
			}
			// delete tests
			if c.DeleteDeclareCapture(1, "test_delete", "", 5) != nil {
				t.Errorf("Deleting an existing declare capture failed")
			}
			if c.DeleteDeclareCapture(0, "test_delete", "", 6) != nil {
				t.Errorf("Deleting an existing declare capture failed")
			}
		})
	}
}
