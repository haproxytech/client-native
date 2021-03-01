package runtime

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v2/models"
)

// ShowTables returns Stick Tables descriptions from runtime
func (s *SingleRuntime) ShowTables() (models.StickTables, error) {
	response, err := s.ExecuteWithResponse("show table")
	if err != nil {
		return nil, err
	}
	return s.parseStickTables(response), nil
}

// ShowTables returns one Stick Table descriptions from runtime
func (s *SingleRuntime) ShowTable(name string) (*models.StickTable, error) {
	response, err := s.ExecuteWithResponse("show table")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(response, "\n")
	for _, line := range lines {
		stkT := s.parseStickTable(line)
		if stkT == nil || stkT.Name != name {
			continue
		} else {
			return stkT, nil
		}
	}
	return nil, nil
}

// GetTableEntries returns Stick Tables entries
func (s *SingleRuntime) GetTableEntries(name string, filter []string, key string) (models.StickTableEntries, error) {
	cmd := fmt.Sprintf("show table %s", name)

	// use only first filter here
	if len(filter) > 0 {
		cmd = fmt.Sprintf("%s data.%s", cmd, filter[0])
	}

	if key != "" {
		cmd = fmt.Sprintf("%s key %s", cmd, key)
	}

	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(response, "\n")
	entries := models.StickTableEntries{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		entry := parseStickTableEntry(line)
		if entry != nil {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

func (s *SingleRuntime) parseStickTables(output string) models.StickTables {
	lines := strings.Split(output, "\n")

	stkTables := models.StickTables{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || !strings.HasPrefix(strings.TrimSpace(line), "# table:") {
			continue
		}
		stkTable := s.parseStickTable(line)

		if stkTable == nil {
			continue
		}

		stkTables = append(stkTables, stkTable)
	}
	return stkTables
}

func (s *SingleRuntime) parseStickTable(output string) *models.StickTable {
	if !strings.HasPrefix(output, "# table:") {
		return nil
	}
	proc := int64(s.process)
	stkTable := &models.StickTable{Process: &proc}

	stkStrings := strings.Split(output, ",")

	for _, stkT := range stkStrings {
		switch {
		case strings.HasPrefix(stkT, "# table:"):
			stkTable.Name = strings.TrimSpace(stkT[len("# table:"):])
		case strings.HasPrefix(stkT, " type:"):
			stkTable.Type = strings.TrimSpace(stkT[len(" type:"):])
		case strings.HasPrefix(stkT, " size:"):
			s, _ := strconv.ParseInt(strings.TrimSpace(stkT[len(" size:"):]), 10, 64)
			stkTable.Size = &s
		case strings.HasPrefix(stkT, " used:"):
			u, _ := strconv.ParseInt(strings.TrimSpace(stkT[len(" used:"):]), 10, 64)
			stkTable.Used = &u
		}
	}

	return stkTable
}

func parseStickTableEntry(output string) *models.StickTableEntry { //nolint:gocognit,gocyclo
	idData := strings.SplitN(output, ":", 2)
	if len(idData) != 2 {
		return nil
	}
	entry := &models.StickTableEntry{ID: idData[0]}
	data := strings.Split(strings.TrimSpace(idData[1]), " ")
	for _, d := range data {
		kv := strings.Split(d, "=")
		switch key := kv[0]; {
		case key == "server_id":
			sID, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.ServerID = &sID
			}
		case key == "gpc0":
			gpc0, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.Gpc0 = &gpc0
			}
		case strings.HasPrefix(key, "gpc0_rate("):
			gpc0Rate, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.Gpc0Rate = &gpc0Rate
			}
		case key == "gpc1":
			gpc1, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.Gpc1 = &gpc1
			}
		case strings.HasPrefix(key, "gpc1_rate("):
			gpc1Rate, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.Gpc1Rate = &gpc1Rate
			}
		case key == "conn_cnt":
			connCnt, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.ConnCnt = &connCnt
			}
		case key == "conn_cur":
			connCur, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.ConnCur = &connCur
			}
		case strings.HasPrefix(key, "conn_rate("):
			connRate, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.ConnRate = &connRate
			}
		case key == "sess_cnt":
			sessCnt, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.SessCnt = &sessCnt
			}
		case strings.HasPrefix(key, "sess_rate("):
			sessRate, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.SessRate = &sessRate
			}
		case key == "http_req_cnt":
			httpReqCnt, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.HTTPReqCnt = &httpReqCnt
			}
		case strings.HasPrefix(key, "http_req_rate("):
			httpReqRate, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.HTTPReqRate = &httpReqRate
			}
		case key == "http_err_cnt":
			httpErrCnt, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.HTTPErrCnt = &httpErrCnt
			}
		case strings.HasPrefix(key, "http_err_rate("):
			httpErrRate, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.HTTPErrRate = &httpErrRate
			}
		case key == "bytes_in_cnt":
			bytesInCnt, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.BytesInCnt = &bytesInCnt
			}
		case strings.HasPrefix(key, "bytes_in_rate("):
			bytesInRate, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.BytesInRate = &bytesInRate
			}
		case key == "bytes_out_cnt":
			bytesOutCnt, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.BytesOutCnt = &bytesOutCnt
			}
		case strings.HasPrefix(key, "bytes_out_rate("):
			bytesOutRate, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.BytesOutRate = &bytesOutRate
			}
		case key == "use":
			entry.Use = strings.TrimSpace(kv[1]) == "1"
		case key == "exp":
			exp, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
			if err == nil {
				entry.Exp = &exp
			}
		case key == "key":
			entry.Key = strings.TrimSpace(kv[1])
		}
	}
	return entry
}
