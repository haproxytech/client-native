// Copyright 2024 HAProxy Technologies
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

package test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func tracesExpcectation() *models.Traces {
	initStructuredExpected()
	res := StructuredToTracesMap()
	return &res
}

func TestGetTraces(t *testing.T) {
	v, traces, err := clientTest.GetTraces("")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	checkTraces(t, traces)
}

func checkTraces(t *testing.T, traces *models.Traces) {
	want := tracesExpcectation()
	require.True(t, traces.Equal(*want), "diff %v", cmp.Diff(traces, want))
}

func TestCreateEditDeleteTraces(t *testing.T) {

	// We cannot start by creating a new traces section, since there is
	// already one present in $testConf, and this is a unique section.
	// So let's delete first.
	err := clientTest.DeleteTraces("", version)
	if err != nil {
		t.Error(err)
	} else {
		version++
	}

	traces := &models.Traces{
		Entries: models.TraceEntries{
			&models.TraceEntry{Trace: "test trace entry 1"},
			&models.TraceEntry{Trace: "test trace entry 2"},
			&models.TraceEntry{Trace: "test trace entry 3"},
		},
	}

	err = clientTest.CreateTraces(traces, "", version)
	if err != nil {
		t.Fatal(err)
	} else {
		version++
	}

	v, created, err := clientTest.GetTraces("")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(created, traces) {
		t.Log(cmp.Diff(created, traces))
		t.Fatal("Created traces section not equal to the given one")
	}

	if v != version {
		t.Errorf("found version %d, expected %d", v, version)
	}

	err = clientTest.CreateTraces(traces, "", version)
	if err == nil {
		t.Fatal("Should throw error traces section already exists")
		version++
	}

	// Modify the section.
	traces.Entries[1].Trace = "test trace entry 2 bis"
	traces.Entries = append(traces.Entries, &models.TraceEntry{Trace: "test trace entry 4"})
	err = clientTest.EditTraces(traces, "", version)
	if err != nil {
		t.Errorf("EditTraces: %v", err)
	} else {
		version++
	}

	// Check if the modification was effective.
	v, created, err = clientTest.GetTraces("")
	if err != nil {
		t.Fatal(err)
	}
	if created == nil {
		t.Fatal("got a nil Traces")
	}
	if v != version {
		t.Errorf("found version %d, expected %d", v, version)
	}
	if !reflect.DeepEqual(created, traces) {
		t.Log(cmp.Diff(created, traces))
		t.Fatal("Created traces section not equal to the given one")
	}

	// Delete the section.
	err = clientTest.DeleteTraces("", version)
	if err != nil {
		t.Error(err)
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	err = clientTest.DeleteTraces("", 999999)
	if err != nil {
		if confErr, ok := err.(*configuration.ConfError); ok {
			if !confErr.Is(configuration.ErrVersionMismatch) {
				t.Error("Should throw configuration.ErrVersionMismatch error")
			}
		} else {
			t.Error("Should throw configuration.ErrVersionMismatch error")
		}
	}

	_, _, err = clientTest.GetTraces("")
	if err == nil {
		t.Error("DeleteTraces failed: traces section still exists")
	}

	err = clientTest.DeleteTraces("", version)
	if err == nil {
		t.Error("Should throw error, non existent Traces")
		version++
	}

	// Re-create the traces section for the following tests.
	traces = &models.Traces{
		Entries: models.TraceEntries{
			&models.TraceEntry{Trace: "h1 sink buf1 level developer verbosity complete start now"},
			&models.TraceEntry{Trace: "h2 sink buf2 level developer verbosity complete start now"},
		},
	}
	err = clientTest.CreateTraces(traces, "", version)
	if err != nil {
		t.Fatal(err)
	} else {
		version++
	}
}

func TestAddRemoveTraceEntries(t *testing.T) {
	v, traces, err := clientTest.GetTraces("")
	if err != nil {
		t.Fatal(err)
	}
	if len(traces.Entries) == 0 {
		t.Fatal("zero log entries found in testConf")
	}

	initial := len(traces.Entries)

	entry := &models.TraceEntry{Trace: "h2 oops buf2"}

	if err = clientTest.CreateTraceEntry(entry, "", v); err != nil {
		t.Fatal("CreateTraceEntry: ", err)
	}
	version++

	v, traces, err = clientTest.GetTraces("")
	if err != nil {
		t.Fatal(err)
	}

	if len(traces.Entries) != initial+1 {
		for i, e := range traces.Entries {
			t.Logf("trace entry %d: %v", i, e.Trace)
		}
		t.Fatal("wrong number of entries")
	}

	if err = clientTest.DeleteTraceEntry(entry, "", v); err != nil {
		t.Fatal("DeleteTraceEntry: ", err)
	}
	version++

	v, traces, err = clientTest.GetTraces("")
	if err != nil {
		t.Fatal(err)
	}

	if len(traces.Entries) != initial {
		for i, e := range traces.Entries {
			t.Logf("trace entry %d: %v", i, e.Trace)
		}
		t.Fatal("wrong number of entries")
	}
}
