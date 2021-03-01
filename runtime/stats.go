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

package runtime

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/haproxytech/client-native/v2/models"
)

// GetStats fetches HAProxy stats from runtime API
func (s *SingleRuntime) GetStats() *models.NativeStatsCollection {
	rAPI := ""
	if s.worker != 0 {
		rAPI = fmt.Sprintf("%s@%v", s.socketPath, s.worker)
	} else {
		rAPI = s.socketPath
	}
	result := &models.NativeStatsCollection{RuntimeAPI: rAPI}
	rawdata, err := s.ExecuteRaw("show stat")
	if err != nil {
		result.Error = err.Error()
		return result
	}
	lines := strings.Split(rawdata[2:], "\n")
	stats := []*models.NativeStat{}
	keys := strings.Split(lines[0], ",")
	for i := 1; i < len(lines); i++ {
		data := map[string]string{}
		line := strings.Split(lines[i], ",")
		if len(line) < len(keys) {
			continue
		}
		for index, key := range keys {
			if len(line[index]) > 0 {
				data[key] = line[index]
			}
		}
		oneLineData := &models.NativeStat{}
		tString := strings.ToLower(line[1])
		if tString == "backend" || tString == "frontend" {
			oneLineData.Name = line[0]
			oneLineData.Type = tString
		} else {
			oneLineData.Name = tString
			oneLineData.Type = "server"
			oneLineData.BackendName = line[0]
		}

		var st models.NativeStatStats
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result:           &st,
			WeaklyTypedInput: true,
			TagName:          "json",
		})
		if err != nil {
			continue
		}

		err = decoder.Decode(data)
		if err != nil {
			continue
		}
		oneLineData.Stats = &st

		stats = append(stats, oneLineData)
	}
	result.Stats = stats
	return result
}
