package runtime

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"

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
		oneLineData := &models.NativeStatsItems{}
		tString := strings.ToLower(line[1])
		if tString == "backend" || tString == "frontend" {
			oneLineData.Name = line[0]
			oneLineData.Type = tString
		} else {
			oneLineData.Name = tString
			oneLineData.Type = "server"
			oneLineData.BackendName = line[0]
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
func (s *SingleRuntime) GetInfo() (models.ProcessInfoHaproxy, error) {
	dataStr, err := s.ExecuteRaw("show info typed")
	data := models.ProcessInfoHaproxy{}
	if err != nil {
		fmt.Println(err.Error())
		return data, err
	}
	return parseInfo(dataStr)
}

func parseInfo(info string) (models.ProcessInfoHaproxy, error) {
	data := models.ProcessInfoHaproxy{}

	for _, line := range strings.Split(info, "\n") {
		fields := strings.Split(line, ":")
		fID := strings.TrimSpace(strings.Split(fields[0], ".")[0])
		switch fID {
		case "1":
			data.Version = fields[3]
		case "2":
			d := strfmt.Date{}
			err := d.Scan(strings.Replace(fields[3], "/", "-", -1))
			if err == nil {
				data.ReleaseDate = d
			}
		case "4":
			nbproc, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.Processes = &nbproc
			}
		case "8":
			uptime, err := strconv.ParseInt(fields[3], 10, 64)
			if err == nil {
				data.Uptime = &uptime
			}
		}
		data.Time = strfmt.DateTime(time.Now())
	}

	return data, nil
}
