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
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestStructuredLogForwards(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.LogForwards)
	v, logForwards, err := clientTest.GetStructuredLogForwards("")
	require.NoError(t, err)
	require.Equal(t, 1, len(logForwards), "%v logForwards returned, expected 1", len(logForwards))

	for _, v := range logForwards {
		m[v.Name] = models.LogForwards{v}
	}

	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	checkStructuredLogForward(t, m)
}

func checkStructuredLogForward(t *testing.T, got map[string]models.LogForwards) {
	exp := logForwardExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if g.Name == w.Name {
					require.True(t, g.Equal(*w), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
	}
}

func TestGetStructuredLogForward(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	v, r, err := clientTest.GetStructuredLogForward("sylog-loadb", "")
	require.NoError(t, err)
	m := make(map[string]models.LogForwards)
	m["sylog-loadb"] = models.LogForwards{r}
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	checkStructuredLogForward(t, m)

	_, err = r.MarshalBinary()
	require.NoError(t, err)

	_, _, err = clientTest.GetStructuredLogForward("doesnotexist", "")
	require.Error(t, err, "should throw error, non existent log forwards section")
}

func TestCreateEditDeleteStructuredLogForward(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	backlog := int64(50)
	maxconn := int64(2000)
	TimeoutClient := int64(5)

	lf := &models.LogForward{
		LogForwardBase: models.LogForwardBase{
			Name:          "created_log_forward",
			Backlog:       &backlog,
			Maxconn:       &maxconn,
			TimeoutClient: &TimeoutClient,
		},
		Binds: map[string]models.Bind{
			"192.168.1.1:9200": {
				BindParams: models.BindParams{Name: "192.168.1.1:9200"},
				Address:    "192.168.1.1",
				Port:       misc.Ptr[int64](9200),
			},
		},
		DgramBinds: map[string]models.DgramBind{
			"bind1": {
				Name:    "bind1",
				Address: "192.168.1.1",
				Port:    misc.Ptr[int64](820),
			},
		},
		LogTargetList: models.LogTargets{
			&models.LogTarget{
				Address:  "192.169.0.1",
				Facility: "mail",
				Global:   true,
			},
		},
	}
	err = clientTest.CreateStructuredLogForward(lf, "", version)
	require.NoError(t, err)
	version++

	v, logForward, err := clientTest.GetStructuredLogForward("created_log_forward", "")
	require.NoError(t, err)

	require.True(t, logForward.Equal(*lf), "log_forward=%s - diff %v", logForward.Name, cmp.Diff(*logForward, *lf))
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.CreateStructuredLogForward(lf, "", version)
	require.Error(t, err, "should throw error log forward already exists")

	err = clientTest.DeleteLogForward("created_log_forward", "", version)
	require.NoError(t, err)
	version++

	v, _ = clientTest.GetVersion("")
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.DeleteLogForward("created_log_forward", "", 999999)
	require.Error(t, err, "Should throw error, non existent frontend")
	require.ErrorIs(t, err, configuration.ErrVersionMismatch, "Should throw configuration.ErrVersionMismatch error")

	_, _, err = clientTest.GetStructuredLogForward("created_log_forward", "")
	require.Error(t, err, "deleteLogForward failed, log forward created_log_forward still exists")

	err = clientTest.DeleteLogForward("doesnotexist", "", version)
	require.Error(t, err, "should throw error, non existent log forward")
}
