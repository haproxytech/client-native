package runtime

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/errors"
	"github.com/haproxytech/client-native/v6/models"
)

// SetTableEntry create or update a stick-table entry in the table.
func (s *SingleRuntime) SetTableEntry(table, key string, dataType models.StickTableEntry) error {
	b, err := json.Marshal(dataType)
	if err != nil {
		return err
	}
	var marshalDataType map[string]any
	if err = json.Unmarshal(b, &marshalDataType); err != nil {
		return err
	}

	const setTableEntryCommand = "set table %s key %s data.%s %v"
	for k, v := range marshalDataType {
		if k == "id" || k == "key" || k == "use" {
			continue
		}
		command := fmt.Sprintf(setTableEntryCommand, table, key, k, v)
		_, err := s.ExecuteWithResponse(command)
		if err != nil {
			return err
		}
	}
	return nil
}

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

	lines := strings.SplitSeq(response, "\n")
	for line := range lines {
		stkT := s.parseStickTable(line)
		if stkT == nil || stkT.Name != name {
			continue
		}
		return stkT, nil
	}
	return nil, fmt.Errorf("no data for table %s: %w", name, errors.ErrNotFound)
}

// GetTableEntries returns Stick Tables entries
func (s *SingleRuntime) GetTableEntries(name string, filter []string, key string) (models.StickTableEntries, error) {
	cmd := "show table " + name

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
	stkTable := &models.StickTable{}

	stkStrings := strings.SplitSeq(output, ",")

	for stkT := range stkStrings {
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

func parseStickTableEntry(output string) *models.StickTableEntry { //nolint:gocognit,gocyclo,cyclop, maintidx
	idData := strings.SplitN(output, ":", 2)
	if len(idData) != 2 {
		return nil
	}
	entry := &models.StickTableEntry{ID: idData[0]}
	data := parseStickTableEntryLine(strings.TrimSpace(idData[1]))
	for k, v := range data {
		switch key := k; {
		case key == "server_id":
			sID, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.ServerID = &sID
			}
		case strings.HasPrefix(key, "gpc("):
			gpc, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.Gpc = &models.StickTableEntryGpc{
					Value: &gpc,
				}
				if index, ok := parseIndex(key); ok {
					entry.Gpc.Idx = index
				}
			}
		case strings.HasPrefix(key, "gpc_rate(") && strings.Contains(key, ","):
			gpcRate, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.GpcRate = &models.StickTableEntryGpcRate{
					Value: &gpcRate,
				}
				if index, ok := parseIndexRate(key); ok {
					entry.GpcRate.Idx = index
				}
			}

		case key == "gpc0":
			gpc0, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.Gpc0 = &gpc0
			}
		case strings.HasPrefix(key, "gpc0_rate("):
			gpc0Rate, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.Gpc0Rate = &gpc0Rate
			}
		case key == "gpc1":
			gpc1, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.Gpc1 = &gpc1
			}
		case key == "gpt0":
			gpt0, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.Gpt0 = &gpt0
			}
		case strings.HasPrefix(key, "gpt("):
			gpt, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.Gpt = &models.StickTableEntryGpt{
					Value: &gpt,
				}
				if index, ok := parseIndex(key); ok {
					entry.Gpt.Idx = index
				}
			}
		case strings.HasPrefix(key, "gpc1_rate("):
			gpc1Rate, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.Gpc1Rate = &gpc1Rate
			}
		case key == "conn_cnt":
			connCnt, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.ConnCnt = &connCnt
			}
		case key == "conn_cur":
			connCur, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.ConnCur = &connCur
			}
		case strings.HasPrefix(key, "conn_rate("):
			connRate, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.ConnRate = &connRate
			}
		case key == "sess_cnt":
			sessCnt, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.SessCnt = &sessCnt
			}
		case strings.HasPrefix(key, "sess_rate("):
			sessRate, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.SessRate = &sessRate
			}
		case key == "http_fail_cnt":
			httpFailCnt, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.HTTPFailCnt = &httpFailCnt
			}
		case strings.HasPrefix(key, "http_fail_rate("):
			httpFailRate, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.HTTPFailRate = &httpFailRate
			}

		case key == "http_req_cnt":
			httpReqCnt, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.HTTPReqCnt = &httpReqCnt
			}
		case strings.HasPrefix(key, "http_req_rate("):
			httpReqRate, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.HTTPReqRate = &httpReqRate
			}
		case key == "http_err_cnt":
			httpErrCnt, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.HTTPErrCnt = &httpErrCnt
			}
		case strings.HasPrefix(key, "http_err_rate("):
			httpErrRate, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.HTTPErrRate = &httpErrRate
			}
		case key == "bytes_in_cnt":
			bytesInCnt, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.BytesInCnt = &bytesInCnt
			}
		case strings.HasPrefix(key, "bytes_in_rate("):
			bytesInRate, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.BytesInRate = &bytesInRate
			}
		case key == "bytes_out_cnt":
			bytesOutCnt, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.BytesOutCnt = &bytesOutCnt
			}
		case strings.HasPrefix(key, "bytes_out_rate("):
			bytesOutRate, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.BytesOutRate = &bytesOutRate
			}
		case key == "use":
			entry.Use = strings.TrimSpace(v) == "1"
		case key == "exp":
			exp, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err == nil {
				entry.Exp = &exp
			}
		case key == "key":
			entry.Key = strings.TrimSpace(v)
		}
	}
	return entry
}

func parseStickTableEntryLine(data string) map[string]string {
	words := strings.Split(data, " ")
	retData := make(map[string]string)

	currentKey := ""
	for _, word := range words {
		if currentKey != "" {
			retData[currentKey] = fmt.Sprintf("%s %s", retData[currentKey], word)
			if !strings.HasSuffix(word, "\\") {
				currentKey = ""
			}
		} else {
			kv := strings.Split(word, "=")
			if len(kv) == 2 {
				retData[kv[0]] = kv[1]
				if strings.HasPrefix(kv[1], "\"") && strings.HasSuffix(kv[1], "\\") {
					currentKey = kv[0]
				}
			}
		}
	}
	return retData
}

func parseIndex(key string) (int64, bool) {
	// Extract content between parentheses
	indexStr := strings.TrimPrefix(key, "gpc(")
	indexStr = strings.TrimSuffix(indexStr, ")")

	index, err := strconv.ParseInt(strings.TrimSpace(indexStr), 10, 64)
	if err != nil {
		return 0, false
	}

	return index, true
}

func parseIndexRate(key string) (int64, bool) {
	// Extract content between parentheses
	params := strings.TrimPrefix(key, "gpc_rate(")
	params = strings.TrimSuffix(params, ")")

	// Split by comma
	parts := strings.Split(params, ",")
	if len(parts) != 2 {
		return 0, false
	}

	index, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
	if err != nil {
		return 0, false
	}

	return index, true
}
