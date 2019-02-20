package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/models"
)

func TestGetFilters(t *testing.T) {
	filters, err := client.GetFilters("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(filters.Data) != 3 {
		t.Errorf("%v filters returned, expected 3", len(filters.Data))
	}

	if filters.Version != version {
		t.Errorf("Version %v returned, expected %v", filters.Version, version)
	}

	for _, f := range filters.Data {
		if *f.ID == 0 {
			if f.Type != "trace" {
				t.Errorf("%v: Type not trace: %v", *f.ID, f.Type)
			}
			if f.TraceName != "BEFORE-HTTP-COMP" {
				t.Errorf("%v: TraceName not BEFORE-HTTP-COMP: %v", *f.ID, f.TraceName)
			}
			if f.TraceRndParsing != true {
				t.Errorf("%v: TraceRndParsing not true: %v", *f.ID, f.TraceRndParsing)
			}
			if f.TraceHexdump != true {
				t.Errorf("%v: TraceHexdump not true: %v", *f.ID, f.TraceHexdump)
			}
		} else if *f.ID == 1 {
			if f.Type != "compression" {
				t.Errorf("%v: Type not compression: %v", *f.ID, f.Type)
			}
		} else if *f.ID == 2 {
			if f.Type != "trace" {
				t.Errorf("%v: Type not trace: %v", *f.ID, f.Type)
			}
			if f.TraceName != "AFTER-HTTP-COMP" {
				t.Errorf("%v: TraceName not AFTER-HTTP-COMP: %v", *f.ID, f.TraceName)
			}
			if f.TraceRndForwarding != true {
				t.Errorf("%v: TraceRndForwarding not true: %v", *f.ID, f.TraceRndForwarding)
			}
		} else {
			t.Errorf("Expext only filter 1, 2 or 3, %v found", *f.ID)
		}
	}

	fJSON, err := filters.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	filters, err = client.GetFilters("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(filters.Data) > 0 {
		t.Errorf("%v filters returned, expected 0", len(filters.Data))
	}

	if !t.Failed() {
		fmt.Println("GetFilters succesful\nResponse: \n" + string(fJSON) + "\n")
	}
}

func TestGetFilter(t *testing.T) {
	filter, err := client.GetFilter(0, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	f := filter.Data

	if *f.ID != 0 {
		t.Errorf("Filter ID 0, %v found", *f.ID)
	}
	if f.Type != "trace" {
		t.Errorf("%v: Type not trace: %v", *f.ID, f.Type)
	}
	if f.TraceName != "BEFORE-HTTP-COMP" {
		t.Errorf("%v: TraceName not BEFORE-HTTP-COMP: %v", *f.ID, f.TraceName)
	}
	if f.TraceRndParsing != true {
		t.Errorf("%v: TraceRndParsing not true: %v", *f.ID, f.TraceRndParsing)
	}
	if f.TraceHexdump != true {
		t.Errorf("%v: TraceHexdump not true: %v", *f.ID, f.TraceHexdump)
	}

	fJSON, err := f.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, err = client.GetFilter(3, "backend", "test2", "")
	if err == nil {
		t.Error("Should throw error, non existant filter")
	}

	if !t.Failed() {
		fmt.Println("GetFilter succesful\nResponse: \n" + string(fJSON) + "\n")
	}
}

func TestCreateEditDeleteFilter(t *testing.T) {
	// TestCreateFilter
	id := int64(1)
	f := &models.Filter{
		ID:         &id,
		Type:       "spoe",
		SpoeEngine: "test",
		SpoeConfig: "test.cfg",
	}

	err := client.CreateFilter("frontend", "test", f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ondiskF, err := client.GetFilter(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskF.Data, f) {
		fmt.Printf("Created filter: %v\n", ondiskF.Data)
		fmt.Printf("Given filter: %v\n", f)
		t.Error("Created filter not equal to given filter")
	}

	if ondiskF.Version != version {
		t.Errorf("Version %v returned, expected %v", ondiskF.Version, version)
	}

	if !t.Failed() {
		fmt.Println("CreateFilter successful")
	}

	// TestEditFilter
	f = &models.Filter{
		ID:         &id,
		Type:       "spoe",
		SpoeConfig: "bla.cfg",
		SpoeEngine: "bla",
	}

	err = client.EditFilter(1, "frontend", "test", f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	ondiskF, err = client.GetFilter(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskF.Data, f) {
		fmt.Printf("Edited filter: %v\n", ondiskF.Data)
		fmt.Printf("Given filter: %v\n", f)
		t.Error("Edited filter not equal to given filter")
	}

	if ondiskF.Version != version {
		t.Errorf("Version %v returned, expected %v", ondiskF.Version, version)
	}

	if !t.Failed() {
		fmt.Println("EditFilter successful")
	}

	// TestDeleteFilter
	err = client.DeleteFilter(3, "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, err = client.GetFilter(3, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteFilter failed, filter 3 still exists")
	}

	err = client.DeleteFilter(1, "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant filter")
		version++
	}

	if !t.Failed() {
		fmt.Println("DeleteFilter successful")
	}
}
