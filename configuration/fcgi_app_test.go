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

package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/haproxytech/client-native/v4/misc"
	"github.com/haproxytech/client-native/v4/models"
)

func TestGetFCGIApps(t *testing.T) {
	v, apps, err := clientTest.GetFCGIApplications("")

	require.NoError(t, err)
	require.Len(t, apps, 2)
	require.Equal(t, version, v)
}

func TestGetFCGIApp(t *testing.T) {
	v, app, err := clientTest.GetFCGIApplication("test", "")

	require.NoError(t, err)
	require.NotNil(t, app)
	require.Equal(t, version, v)

	assert.Equal(t, "/path/to/chroot", *app.Docroot)
	assert.Equal(t, "index.php", app.Index)
	assert.Equal(t, `^(/.+\.php)(/.*)?$`, app.PathInfo)
	assert.Equal(t, models.FCGIAppGetValuesEnabled, app.GetValues)
	assert.Equal(t, models.FCGIAppGetValuesDisabled, app.KeepConn)
	assert.Equal(t, models.FCGIAppGetValuesDisabled, app.MpxsConns)
	assert.Equal(t, int64(1024), app.MaxReqs)
	// set-param
	require.Len(t, app.SetParams, 3)

	assert.Equal(t, "name", app.SetParams[0].Name)
	assert.Equal(t, "fmt", app.SetParams[0].Format)
	assert.Equal(t, "if", app.SetParams[0].Cond)
	assert.Equal(t, "acl", app.SetParams[0].CondTest)

	assert.Equal(t, "name", app.SetParams[1].Name)
	assert.Equal(t, "fmt", app.SetParams[1].Format)
	assert.Equal(t, "unless", app.SetParams[1].Cond)
	assert.Equal(t, "acl", app.SetParams[1].CondTest)

	assert.Equal(t, "name", app.SetParams[2].Name)
	assert.Equal(t, "fmt", app.SetParams[2].Format)
	assert.Empty(t, app.SetParams[2].Cond)
	assert.Empty(t, app.SetParams[2].CondTest)
	// pass-header
	require.Len(t, app.PassHeaders, 3)

	assert.Equal(t, "x-header", app.PassHeaders[0].Name)
	assert.Equal(t, models.FCGIPassHeaderCondUnless, app.PassHeaders[0].Cond)
	assert.Equal(t, "acl", app.PassHeaders[0].CondTest)

	assert.Equal(t, "x-header", app.PassHeaders[1].Name)
	assert.Equal(t, models.FCGIPassHeaderCondIf, app.PassHeaders[1].Cond)
	assert.Equal(t, "acl", app.PassHeaders[1].CondTest)

	assert.Equal(t, "x-header", app.PassHeaders[1].Name)
	assert.Empty(t, app.PassHeaders[2].Cond)
	assert.Empty(t, app.PassHeaders[2].CondTest)
	// log-stderr

	require.Len(t, app.LogStderrs, 4)

	assert.False(t, app.LogStderrs[0].Global)
	assert.Equal(t, "127.0.0.1:1515", app.LogStderrs[0].Address)
	assert.Equal(t, int64(8192), app.LogStderrs[0].Len)
	assert.Equal(t, "rfc5424", app.LogStderrs[0].Format)
	require.NotNil(t, app.LogStderrs[0].Sample)
	assert.Equal(t, "1,2-5", *app.LogStderrs[0].Sample.Ranges)
	assert.Equal(t, int64(6), *app.LogStderrs[0].Sample.Size)
	assert.Equal(t, "local2", app.LogStderrs[0].Facility)
	assert.Equal(t, "info", app.LogStderrs[0].Level)
	assert.Equal(t, "debug", app.LogStderrs[0].Minlevel)

	assert.False(t, app.LogStderrs[1].Global)
	assert.Equal(t, "127.0.0.1:1515", app.LogStderrs[1].Address)
	assert.Equal(t, int64(8192), app.LogStderrs[1].Len)
	assert.Equal(t, "rfc5424", app.LogStderrs[1].Format)
	require.NotNil(t, app.LogStderrs[1].Sample)
	assert.Equal(t, "1,2-5", *app.LogStderrs[1].Sample.Ranges)
	assert.Equal(t, int64(6), *app.LogStderrs[1].Sample.Size)
	assert.Equal(t, "local2", app.LogStderrs[1].Facility)
	assert.Equal(t, "info", app.LogStderrs[1].Level)
	assert.Empty(t, app.LogStderrs[1].Minlevel)

	assert.False(t, app.LogStderrs[2].Global)
	assert.Equal(t, "127.0.0.1:1515", app.LogStderrs[2].Address)
	assert.Empty(t, app.LogStderrs[2].Len)
	assert.Empty(t, app.LogStderrs[2].Format)
	assert.Nil(t, app.LogStderrs[2].Sample)
	assert.Equal(t, "local2", app.LogStderrs[2].Facility)
	assert.Empty(t, app.LogStderrs[2].Level)
	assert.Empty(t, app.LogStderrs[2].Minlevel)

	assert.True(t, app.LogStderrs[3].Global)
	assert.Empty(t, app.LogStderrs[3].Address)
	assert.Empty(t, app.LogStderrs[3].Len)
	assert.Empty(t, app.LogStderrs[3].Format)
	assert.Nil(t, app.LogStderrs[3].Sample)
	assert.Empty(t, app.LogStderrs[3].Facility)
	assert.Empty(t, app.LogStderrs[3].Level)
	assert.Empty(t, app.LogStderrs[3].Minlevel)
}

func TestCreateEditDeleteFCGIApp(t *testing.T) {
	// Creating an application
	app := &models.FCGIApp{
		Docroot:   misc.StringP("/path/to/my/chroot/app"),
		GetValues: models.FCGIAppGetValuesEnabled,
		Index:     "index.php",
		KeepConn:  models.FCGIAppKeepConnEnabled,
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
		MpxsConns: models.FCGIAppMpxsConnsEnabled,
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
	app.GetValues = models.FCGIAppGetValuesDisabled
	app.Index = "old_index.php"
	app.KeepConn = models.FCGIAppGetValuesDisabled
	app.LogStderrs[0].Level = "critical"
	app.LogStderrs = []*models.FCGILogStderr{
		app.LogStderrs[0],
	}
	app.MaxReqs = 10
	app.MpxsConns = models.FCGIAppMpxsConnsDisabled
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
