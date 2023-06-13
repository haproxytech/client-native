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
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/client-native/v5/models"
)

func TestGetFilters(t *testing.T) { //nolint:gocognit
	v, filters, err := clientTest.GetFilters("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(filters) != 6 {
		t.Errorf("%v filters returned, expected 3", len(filters))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, f := range filters {
		switch *f.Index {
		case 0:
			if f.Type != "trace" {
				t.Errorf("%v: Type not trace: %v", *f.Index, f.Type)
			}
			if f.TraceName != "BEFORE-HTTP-COMP" {
				t.Errorf("%v: TraceName not BEFORE-HTTP-COMP: %v", *f.Index, f.TraceName)
			}
			if f.TraceRndParsing != true {
				t.Errorf("%v: TraceRndParsing not true: %v", *f.Index, f.TraceRndParsing)
			}
			if f.TraceHexdump != true {
				t.Errorf("%v: TraceHexdump not true: %v", *f.Index, f.TraceHexdump)
			}
		case 1:
			if f.Type != "compression" {
				t.Errorf("%v: Type not compression: %v", *f.Index, f.Type)
			}
		case 2:
			if f.Type != "trace" {
				t.Errorf("%v: Type not trace: %v", *f.Index, f.Type)
			}
			if f.TraceName != "AFTER-HTTP-COMP" {
				t.Errorf("%v: TraceName not AFTER-HTTP-COMP: %v", *f.Index, f.TraceName)
			}
			if f.TraceRndForwarding != true {
				t.Errorf("%v: TraceRndForwarding not true: %v", *f.Index, f.TraceRndForwarding)
			}
		case 3:
			if f.Type != "fcgi-app" {
				t.Errorf("%v: Type not fcgi-app: %v", *f.Index, f.Type)
			}
			if f.AppName != "my-app" {
				t.Errorf("%v: AppName not my-app: %v", *f.Index, f.AppName)
			}
		case 4:
			if f.Type != "bwlim-in" {
				t.Errorf("%v: Type not bwlim-in: %v", *f.Index, f.Type)
			}
			if f.BandwidthLimitName != "in" {
				t.Errorf("%v: BandwidthLimitName not in: %v", *f.Index, f.BandwidthLimitName)
			}
			if f.DefaultLimit != 1024 {
				t.Errorf("%v: DefaultLimit not 1024: %v", *f.Index, f.DefaultLimit)
			}
			if f.DefaultPeriod != 1000 {
				t.Errorf("%v: DefaultPeriod not 1000: %v", *f.Index, f.DefaultPeriod)
			}
			if f.MinSize != 1048576 {
				t.Errorf("%v: MinSize not 1048576: %v", *f.Index, f.MinSize)
			}
		case 5:
			if f.Type != "bwlim-out" {
				t.Errorf("%v: Type not bwlim-out: %v", *f.Index, f.Type)
			}
			if f.BandwidthLimitName != "out" {
				t.Errorf("%v: BandwidthLimitName not out: %v", *f.Index, f.BandwidthLimitName)
			}
			if f.Limit != 1024 {
				t.Errorf("%v: Limit not 1024: %v", *f.Index, f.Limit)
			}
			if f.Key != "name(arg1)" {
				t.Errorf("%v: Limit not name(arg1): %v", *f.Index, f.Key)
			}
			if f.Table != "st_src_global" {
				t.Errorf("%v: Limit not st_src_global: %v", *f.Index, f.Table)
			}
			if f.MinSize != 32 {
				t.Errorf("%v: MinSize not 32: %v", *f.Index, f.MinSize)
			}
		default:
			t.Errorf("Expext only filter 1, 2 or 3, %v found", *f.Index)
		}
	}

	_, filters, err = clientTest.GetFilters("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(filters) > 0 {
		t.Errorf("%v filters returned, expected 0", len(filters))
	}
}

func TestGetFilter(t *testing.T) {
	v, f, err := clientTest.GetFilter(0, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if *f.Index != 0 {
		t.Errorf("Filter ID 0, %v found", *f.Index)
	}
	if f.Type != "trace" {
		t.Errorf("%v: Type not trace: %v", *f.Index, f.Type)
	}
	if f.TraceName != "BEFORE-HTTP-COMP" {
		t.Errorf("%v: TraceName not BEFORE-HTTP-COMP: %v", *f.Index, f.TraceName)
	}
	if f.TraceRndParsing != true {
		t.Errorf("%v: TraceRndParsing not true: %v", *f.Index, f.TraceRndParsing)
	}
	if f.TraceHexdump != true {
		t.Errorf("%v: TraceHexdump not true: %v", *f.Index, f.TraceHexdump)
	}

	_, err = f.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetFilter(3, "backend", "test2", "")
	if err == nil {
		t.Error("Should throw error, non existant filter")
	}
}

func TestCreateEditDeleteFilter(t *testing.T) {
	// TestCreateFilter
	id := int64(1)
	f := &models.Filter{
		Index:      &id,
		Type:       "spoe",
		SpoeEngine: "test",
		SpoeConfig: "test.cfg",
	}

	err := clientTest.CreateFilter("frontend", "test", f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskF, err := clientTest.GetFilter(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskF, f) {
		fmt.Printf("Created filter: %v\n", ondiskF)
		fmt.Printf("Given filter: %v\n", f)
		t.Error("Created filter not equal to given filter")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditFilter
	f = &models.Filter{
		Index:      &id,
		Type:       "spoe",
		SpoeConfig: "bla.cfg",
		SpoeEngine: "bla",
	}

	err = clientTest.EditFilter(1, "frontend", "test", f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskF, err = clientTest.GetFilter(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskF, f) {
		fmt.Printf("Edited filter: %v\n", ondiskF)
		fmt.Printf("Given filter: %v\n", f)
		t.Error("Edited filter not equal to given filter")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteFilter
	err = clientTest.DeleteFilter(3, "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, filters, _ := clientTest.GetFilters("frontend", "test", "")
	_ = filters

	_, _, err = clientTest.GetFilter(6, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteFilter failed, filter 5 still exists")
	}

	err = clientTest.DeleteFilter(1, "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant filter")
		version++
	}
}
