package runtime

import (
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/haproxytech/models"
)

//GetStats fetches HAProxy stats from runtime API
func (s *SingleRuntime) GetStats() (models.NativeStats, error) {
	rawdata, err := s.ExecuteRaw("show stat")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(rawdata[2:], "\n")
	result := models.NativeStats{}
	keys := strings.Split(lines[0], ",")
	//data := []map[string]string{}
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
		oneLineData := &models.NativeStatsItems{
			Name: line[0],
			Type: strings.ToLower(line[1]),
		}
		var st models.NativeStatsItemsStats
		err := mapstructure.WeakDecode(data, &st)
		if err != nil {
			continue
		}
		oneLineData.Stats = &st
		result = append(result, oneLineData)
	}
	return result, nil
}

//GetInfo fetches HAProxy info from runtime API
func (s *SingleRuntime) GetInfo() (string, error) {
	data, err := s.ExecuteRaw("show stat")
	if err != nil {
		return "", err
	}
	return data, nil
}
