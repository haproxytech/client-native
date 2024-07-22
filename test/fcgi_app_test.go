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

package test

import (
	_ "embed"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

func fcgiAppExpectation() map[string]models.FCGIApps {
	initStructuredExpected()
	res := StructuredToFCGIAppMap()
	// Add individual entries
	for _, vs := range res {
		for _, v := range vs {
			key := v.Name
			res[key] = models.FCGIApps{v}
		}
	}
	return res
}

func TestGetFCGIApps(t *testing.T) {
	m := make(map[string]models.FCGIApps)
	v, apps, err := clientTest.GetFCGIApplications("")
	m[""] = apps

	require.NoError(t, err)
	require.Len(t, apps, 2)
	require.Equal(t, version, v)
	checkFCGIApp(t, m)
}

func TestGetFCGIApp(t *testing.T) {
	mapps := make(map[string]models.FCGIApps)
	v, app, err := clientTest.GetFCGIApplication("test", "")

	require.NoError(t, err)
	require.NotNil(t, app)
	mapps["test"] = models.FCGIApps{app}

	require.Equal(t, version, v)

	checkFCGIApp(t, mapps)
}

func checkFCGIApp(t *testing.T, got map[string]models.FCGIApps) {
	exp := fcgiAppExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if g.Name == w.Name {
					require.True(t, g.FCGIAppBase.Equal(w.FCGIAppBase), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
	}
}

func TestCreateEditDeleteFCGIApp(t *testing.T) {
	// Creating an application
	app := &models.FCGIApp{
		FCGIAppBase: models.FCGIAppBase{
			Docroot:   misc.StringP("/path/to/my/chroot/app"),
			GetValues: models.FCGIAppBaseGetValuesEnabled,
			Index:     "index.php",
			KeepConn:  models.FCGIAppBaseKeepConnEnabled,
			LogStderrs: []*models.FCGILogStderr{
				{
					Address:  "127.0.0.1:514",
					Facility: "debug",
				},
				{
					Address:  "127.0.0.1:415",
					Facility: "error",
					Format:   "fmt",
					Len:      1024,
					Level:    "info",
					Minlevel: "debug",
					Sample: &models.FCGILogStderrSample{
						Ranges: misc.StringP("1,2-9"),
						Size:   misc.Int64P(10),
					},
				},
			},
			MaxReqs:   100,
			MpxsConns: models.FCGIAppBaseMpxsConnsEnabled,
			Name:      "created",
			PassHeaders: []*models.FCGIPassHeader{
				{
					Cond:     models.FCGIPassHeaderCondIf,
					CondTest: "end",
					Name:     "x-end",
				},
				{
					Cond:     models.FCGIPassHeaderCondUnless,
					CondTest: "start",
					Name:     "x-start",
				},
			},
			PathInfo: `^(/.+\.php)(/.*)?$`,
			SetParams: []*models.FCGISetParam{
				{
					Cond:     models.FCGISetParamCondIf,
					CondTest: "start",
					Format:   "fmt",
					Name:     "start",
				},
				{
					Cond:     models.FCGISetParamCondUnless,
					CondTest: "end",
					Format:   "fmt",
					Name:     "end",
				},
			},
		},
	}
	require.NoError(t, clientTest.CreateFCGIApplication(app, "", version))
	version++
	// Ensuring data is correct
	configVersion, found, err := clientTest.GetFCGIApplication(app.Name, "")
	assert.Nil(t, err)
	assert.Equal(t, configVersion, version)
	assert.EqualValues(t, app, found)
	// Expecting failure due to duplicated entry
	assert.Error(t, clientTest.CreateFCGIApplication(app, "", version))
	// Updating and ensuring data is properly reflected
	app.Docroot = misc.StringP("/no/more/chroot")
	app.GetValues = models.FCGIAppBaseGetValuesDisabled
	app.Index = "old_index.php"
	app.KeepConn = models.FCGIAppBaseGetValuesDisabled
	app.LogStderrs[0].Level = "critical"
	app.LogStderrs = []*models.FCGILogStderr{
		app.LogStderrs[0],
	}
	app.MaxReqs = 10
	app.MpxsConns = models.FCGIAppBaseMpxsConnsDisabled
	app.PassHeaders[0].Name = "x-middle"
	app.PassHeaders[0].CondTest = "middle"
	app.PassHeaders = []*models.FCGIPassHeader{
		app.PassHeaders[0],
	}
	app.PathInfo = `^(/.+\.pl)(/.*)?$`
	app.SetParams[0].Name = "middle"
	app.PassHeaders[0].CondTest = "middle"
	app.PassHeaders = []*models.FCGIPassHeader{
		app.PassHeaders[0],
	}

	require.NoError(t, clientTest.EditFCGIApplication(app.Name, app, "", version))
	version++

	configVersion, found, err = clientTest.GetFCGIApplication(app.Name, "")
	assert.Nil(t, err)
	assert.Equal(t, configVersion, version)
	assert.EqualValues(t, app, found)
	// Deleting and ensuring triggering not found error
	assert.NoError(t, clientTest.DeleteFCGIApplication(app.Name, "", version))
	version++

	_, _, err = clientTest.GetFCGIApplication(app.Name, "")
	assert.Error(t, err)
}
