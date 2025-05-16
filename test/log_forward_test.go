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

package test

import (
	_ "embed"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func logForwardExpectation() map[string]models.LogForwards {
	initStructuredExpected()
	res := StructuredToLogForwardMap()
	// Add individual entries
	for _, vs := range res {
		for _, v := range vs {
			key := v.Name
			res[key] = models.LogForwards{v}
		}
	}
	return res
}

func TestLogForwards(t *testing.T) {
	m := make(map[string]models.LogForwards)
	v, logForwards, err := clientTest.GetLogForwards("")
	if err != nil {
		t.Error(err.Error())
	}
	if len(logForwards) != 1 {
		t.Errorf("%v logForwards returned, expected 1", len(logForwards))
	}
	for _, v := range logForwards {
		m[v.Name] = models.LogForwards{v}
	}

	if v != version {
		t.Errorf("version %v returned, expected %v", v, version)
	}

	checkLogForward(t, m)
}

func checkLogForward(t *testing.T, got map[string]models.LogForwards) {
	exp := logForwardExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if g.Name == w.Name {
					require.True(t, g.LogForwardBase.Equal(w.LogForwardBase), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
	}
}

func TestGetLogForward(t *testing.T) {
	v, r, err := clientTest.GetLogForward("sylog-loadb", "")
	if err != nil {
		t.Error(err.Error())
	}
	m := make(map[string]models.LogForwards)
	m["sylog-loadb"] = models.LogForwards{r}

	if v != version {
		t.Errorf("version %v returned, expected %v", v, version)
	}

	checkLogForward(t, m)

	_, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetLogForward("doesnotexist", "")
	if err == nil {
		t.Error("should throw error, non existent log forwards section")
	}
}

func TestCreateEditDeleteLogForward(t *testing.T) {
	backlog := int64(50)
	maxconn := int64(2000)
	TimeoutClient := int64(5)

	lf := &models.LogForward{
		LogForwardBase: models.LogForwardBase{
			Name:             "created_log_forward",
			AssumeRfc6587Ntf: true,
			DontParseLog:     true,
			Backlog:          &backlog,
			Maxconn:          &maxconn,
			TimeoutClient:    &TimeoutClient,
		},
	}
	err := clientTest.CreateLogForward(lf, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, logForward, err := clientTest.GetLogForward("created_log_forward", "")
	require.NoError(t, err, "failed to get log forward")
	require.Equal(t, lf, logForward)
	require.Equal(t, version, v, "version not incremented")

	err = clientTest.CreateLogForward(lf, "", version)
	if err == nil {
		t.Error("should throw error log forward already exists")
		version++
	}

	backlog++
	maxconn++
	TimeoutClient++
	lf.AssumeRfc6587Ntf = false
	lf.DontParseLog = false
	lf.Backlog = &backlog
	lf.Maxconn = &maxconn
	lf.TimeoutClient = &TimeoutClient
	err = clientTest.EditLogForward("created_log_forward", lf, "", version)
	require.NoError(t, err, "failed to edit log forward")
	version++
	if v, _ := clientTest.GetVersion(""); v != version {
		t.Fatalf("version not incremented")
	}
	v, logForward, err = clientTest.GetLogForward("created_log_forward", "")
	require.NoError(t, err, "failed to get log forward")
	require.Equal(t, lf, logForward, "log forward has not been updated")

	err = clientTest.DeleteLogForward("created_log_forward", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("version not incremented")
	}

	err = clientTest.DeleteLogForward("created_log_forward", "", 999999)
	if err != nil {
		if confErr, ok := err.(*configuration.ConfError); ok {
			if !confErr.Is(configuration.ErrVersionMismatch) {
				t.Error("should throw configuration.ErrVersionMismatch error")
			}
		} else {
			t.Error("should throw configuration.ErrVersionMismatch error")
		}
	}
	_, _, err = clientTest.GetLogForward("created_log_forward", "")
	if err == nil {
		t.Error("deleteLogForward failed, log forward created_log_forward still exists")
	}

	err = clientTest.DeleteLogForward("doesnotexist", "", version)
	if err == nil {
		t.Error("should throw error, non existent log forward")
		version++
	}
}
