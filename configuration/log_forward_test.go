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
	"testing"

	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

func TestLogForwards(t *testing.T) {
	v, logForwards, err := clientTest.GetLogForwards("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(logForwards) != 1 {
		t.Errorf("%v logForwards returned, expected 1", len(logForwards))
	}

	if v != version {
		t.Errorf("version %v returned, expected %v", v, version)
	}

	if logForwards[0].Name != "sylog-loadb" {
		t.Errorf("expected only test, %v found", logForwards[0].Name)
	}
}

func TestGetLogForward(t *testing.T) {
	v, r, err := clientTest.GetLogForward("sylog-loadb", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("version %v returned, expected %v", v, version)
	}

	if r.Name != "sylog-loadb" {
		t.Errorf("expected sylog-loadb log forward, %v found", r.Name)
	}

	if *r.Backlog != *misc.Int64P(10) {
		t.Errorf("expected backlog 10, %v found", r.Backlog)
	}

	_, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetLogForward("doesnotexist", "")
	if err == nil {
		t.Error("should throw error, non existant log forwards section")
	}
}

func TestCreateEditDeleteLogForward(t *testing.T) {
	backlog := int64(50)
	maxconn := int64(2000)
	TimeoutClient := int64(5)

	lf := &models.LogForward{
		Name:          "created_log_forward",
		Backlog:       &backlog,
		Maxconn:       &maxconn,
		TimeoutClient: &TimeoutClient,
	}
	err := clientTest.CreateLogForward(lf, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, logForward, err := clientTest.GetLogForward("created_log_forward", "")
	if err != nil {
		t.Error(err.Error())
	}

	if *logForward.Backlog != *lf.Backlog {
		t.Errorf("backlog expected %v got %v", *logForward.Backlog, *lf.Backlog)
	}

	if *logForward.Maxconn != *lf.Maxconn {
		t.Errorf("maxconn expected %v got %v", *logForward.Maxconn, *lf.Maxconn)
	}

	if *logForward.TimeoutClient != *lf.TimeoutClient {
		t.Errorf("timeout connect expected %v got %v", *logForward.TimeoutClient, *lf.TimeoutClient)
	}

	if v != version {
		t.Errorf("version expected %v got %v", v, version)
	}

	err = clientTest.CreateLogForward(lf, "", version)
	if err == nil {
		t.Error("should throw error log forward already exists")
		version++
	}

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
		if confErr, ok := err.(*ConfError); ok {
			if confErr.Code() != ErrVersionMismatch {
				t.Error("should throw ErrVersionMismatch error")
			}
		} else {
			t.Error("should throw ErrVersionMismatch error")
		}
	}
	_, _, err = clientTest.GetLogForward("created_log_forward", "")
	if err == nil {
		t.Error("deleteLogForward failed, log forward created_log_forward still exists")
	}

	err = clientTest.DeleteLogForward("doesnotexist", "", version)
	if err == nil {
		t.Error("should throw error, non existant log forward")
		version++
	}
}
