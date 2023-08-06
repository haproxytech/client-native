// Code generated with struct_equal_generator; DO NOT EDIT.

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

package models

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/brianvoe/gofakeit/v6"

	jsoniter "github.com/json-iterator/go"
)

func TestStickTableEntryEqual(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b StickTableEntry
	}{}
	for i := 0; i < 100; i++ {
		var sample StickTableEntry
		var result StickTableEntry
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err = enc.Encode(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b StickTableEntry
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if !result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected StickTableEntry to be equal, but it is not %s %s", a, b)
		}
	}
}

func TestStickTableEntryEqualFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b StickTableEntry
	}{}
	for i := 0; i < 100; i++ {
		var sample StickTableEntry
		var result StickTableEntry
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.BytesInCnt = Ptr(*sample.BytesInCnt + 1)
		result.BytesInRate = Ptr(*sample.BytesInRate + 1)
		result.BytesOutCnt = Ptr(*sample.BytesOutCnt + 1)
		result.BytesOutRate = Ptr(*sample.BytesOutRate + 1)
		result.ConnCnt = Ptr(*sample.ConnCnt + 1)
		result.ConnCur = Ptr(*sample.ConnCur + 1)
		result.ConnRate = Ptr(*sample.ConnRate + 1)
		result.Exp = Ptr(*sample.Exp + 1)
		result.Gpc0 = Ptr(*sample.Gpc0 + 1)
		result.Gpc0Rate = Ptr(*sample.Gpc0Rate + 1)
		result.Gpc1 = Ptr(*sample.Gpc1 + 1)
		result.Gpc1Rate = Ptr(*sample.Gpc1Rate + 1)
		result.Gpt0 = Ptr(*sample.Gpt0 + 1)
		result.HTTPErrCnt = Ptr(*sample.HTTPErrCnt + 1)
		result.HTTPErrRate = Ptr(*sample.HTTPErrRate + 1)
		result.HTTPReqCnt = Ptr(*sample.HTTPReqCnt + 1)
		result.HTTPReqRate = Ptr(*sample.HTTPReqRate + 1)
		result.ServerID = Ptr(*sample.ServerID + 1)
		result.SessCnt = Ptr(*sample.SessCnt + 1)
		result.SessRate = Ptr(*sample.SessRate + 1)
		result.Use = !sample.Use
		samples = append(samples, struct {
			a, b StickTableEntry
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Equal(sample.b)
		if result {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected StickTableEntry to be different, but it is not %s %s", a, b)
		}
	}
}

func TestStickTableEntryDiff(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b StickTableEntry
	}{}
	for i := 0; i < 100; i++ {
		var sample StickTableEntry
		var result StickTableEntry
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err = enc.Encode(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&result)
		if err != nil {
			t.Errorf(err.Error())
		}

		samples = append(samples, struct {
			a, b StickTableEntry
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 0 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected StickTableEntry to be equal, but it is not %s %s, %v", a, b, result)
		}
	}
}

func TestStickTableEntryDiffFalse(t *testing.T) {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	samples := []struct {
		a, b StickTableEntry
	}{}
	for i := 0; i < 100; i++ {
		var sample StickTableEntry
		var result StickTableEntry
		err := gofakeit.Struct(&sample)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = gofakeit.Struct(&result)
		if err != nil {
			t.Errorf(err.Error())
		}
		result.BytesInCnt = Ptr(*sample.BytesInCnt + 1)
		result.BytesInRate = Ptr(*sample.BytesInRate + 1)
		result.BytesOutCnt = Ptr(*sample.BytesOutCnt + 1)
		result.BytesOutRate = Ptr(*sample.BytesOutRate + 1)
		result.ConnCnt = Ptr(*sample.ConnCnt + 1)
		result.ConnCur = Ptr(*sample.ConnCur + 1)
		result.ConnRate = Ptr(*sample.ConnRate + 1)
		result.Exp = Ptr(*sample.Exp + 1)
		result.Gpc0 = Ptr(*sample.Gpc0 + 1)
		result.Gpc0Rate = Ptr(*sample.Gpc0Rate + 1)
		result.Gpc1 = Ptr(*sample.Gpc1 + 1)
		result.Gpc1Rate = Ptr(*sample.Gpc1Rate + 1)
		result.Gpt0 = Ptr(*sample.Gpt0 + 1)
		result.HTTPErrCnt = Ptr(*sample.HTTPErrCnt + 1)
		result.HTTPErrRate = Ptr(*sample.HTTPErrRate + 1)
		result.HTTPReqCnt = Ptr(*sample.HTTPReqCnt + 1)
		result.HTTPReqRate = Ptr(*sample.HTTPReqRate + 1)
		result.ServerID = Ptr(*sample.ServerID + 1)
		result.SessCnt = Ptr(*sample.SessCnt + 1)
		result.SessRate = Ptr(*sample.SessRate + 1)
		result.Use = !sample.Use
		samples = append(samples, struct {
			a, b StickTableEntry
		}{sample, result})
	}

	for _, sample := range samples {
		result := sample.a.Diff(sample.b)
		if len(result) != 23 {
			json := jsoniter.ConfigCompatibleWithStandardLibrary
			a, err := json.Marshal(&sample.a)
			if err != nil {
				t.Errorf(err.Error())
			}
			b, err := json.Marshal(&sample.b)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Errorf("Expected StickTableEntry to be different in 23 cases, but it is not (%d) %s %s", len(result), a, b)
		}
	}
}
