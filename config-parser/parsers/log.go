/*
Copyright 2019 HAProxy Technologies

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package parsers

import (
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

//nolint:gochecknoglobals
var logAllowedFacilities = map[string]struct{}{
	"kern": {}, "user": {}, "mail": {}, "daemon": {},
	"auth": {}, "syslog": {}, "lpr": {}, "news": {},
	"uucp": {}, "cron": {}, "auth2": {}, "ftp": {},
	"ntp": {}, "audit": {}, "alert": {}, "cron2": {},
	"local0": {}, "local1": {}, "local2": {}, "local3": {},
	"local4": {}, "local5": {}, "local6": {}, "local7": {},
}

//nolint:gochecknoglobals
var logAllowedLevels = map[string]struct{}{
	"emerg":   {},
	"alert":   {},
	"crit":    {},
	"err":     {},
	"warning": {},
	"notice":  {},
	"info":    {},
	"debug":   {},
}

type Log struct {
	data        []types.Log
	preComments []string // comments that appear before the actual line
}

func (l *Log) Init() {
	l.data = []types.Log{}
	l.preComments = []string{}
}

//nolint:gocognit
func (l *Log) parse(line string, parts []string, comment string) (*types.Log, error) {
	if len(parts) > 1 && parts[1] == "global" {
		return &types.Log{
			Global:  true,
			Comment: comment,
		}, nil
	}
	if len(parts) < 3 {
		return nil, &errors.ParseError{Parser: "Log", Line: line}
	}
	log := &types.Log{
		Address: parts[1],
		Comment: comment,
	}
	// see if we have length
	currIndex := 2
	if currIndex >= len(parts) {
		return log, &errors.ParseError{Parser: "Log", Line: line}
	}
	if parts[currIndex] == "len" {
		currIndex++
		if currIndex >= len(parts) {
			return log, &errors.ParseError{Parser: "Log", Line: line}
		}
		if num, err := strconv.ParseInt(parts[currIndex], 10, 64); err == nil {
			log.Length = num
			currIndex++
		} else {
			return log, &errors.ParseError{Parser: "Log", Line: line}
		}
	}
	if currIndex < len(parts) && parts[currIndex] == "format" {
		currIndex++
		if currIndex >= len(parts) {
			return log, &errors.ParseError{Parser: "Log", Line: line}
		}
		log.Format = parts[currIndex]
		currIndex++
	}
	if currIndex < len(parts) && parts[currIndex] == "sample" {
		currIndex++
		if currIndex >= len(parts) {
			return log, &errors.ParseError{Parser: "Log", Line: line}
		}
		sampleData := strings.Split(parts[currIndex], ":")
		if len(sampleData) != 2 || sampleData[0] == "" {
			return log, &errors.ParseError{Parser: "Log", Line: line}
		}
		log.SampleRange = sampleData[0]
		if num, err := strconv.ParseInt(sampleData[1], 10, 64); err == nil {
			log.SampleSize = num
		} else {
			return log, &errors.ParseError{Parser: "Log", Line: line}
		}
		currIndex++
	}
	if currIndex < len(parts) && parts[currIndex] == "profile" {
		currIndex++
		if currIndex >= len(parts) {
			return log, &errors.ParseError{Parser: "Log", Line: line}
		}
		log.Profile = parts[currIndex]
		currIndex++
	}
	// we must have facility
	if currIndex >= len(parts) {
		return log, &errors.ParseError{Parser: "Log", Line: line}
	}
	facility := parts[currIndex]
	if _, ok := logAllowedFacilities[facility]; !ok {
		return log, &errors.ParseError{Parser: "Log", Line: line}
	}
	log.Facility = facility
	currIndex++
	// level is optional
	if currIndex >= len(parts) {
		return log, nil
	}
	level := parts[currIndex]
	if _, ok := logAllowedLevels[level]; !ok {
		return log, nil
	}
	log.Level = level
	currIndex++
	// min level is optional
	if currIndex >= len(parts) {
		return log, nil
	}
	level = parts[currIndex]
	if _, ok := logAllowedLevels[level]; !ok {
		return log, nil
	}
	log.MinLevel = level
	return log, nil
}

func (l *Log) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] == "log" {
		log, err := l.parse(line, parts, comment)
		if err != nil {
			return "", &errors.ParseError{Parser: "Log", Line: line}
		}
		l.data = append(l.data, *log)
		return "", nil
	}
	if parts[0] == "no" && parts[1] == "log" {
		l.data = append(l.data, types.Log{
			NoLog:   true,
			Comment: comment,
		})
		return "", nil
	}
	return "", &errors.ParseError{Parser: "Log", Line: line}
}

func (l *Log) Result() ([]common.ReturnResultLine, error) {
	if len(l.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(l.data))
	for index, log := range l.data {
		if log.Global {
			result[index] = common.ReturnResultLine{
				Data:    "log global",
				Comment: log.Comment,
			}
			continue
		}
		if log.NoLog {
			result[index] = common.ReturnResultLine{
				Data:    "no log",
				Comment: log.Comment,
			}
			continue
		}
		var sb strings.Builder
		sb.WriteString("log ")
		sb.WriteString(log.Address)
		if log.Length > 0 {
			sb.WriteString(" len ")
			sb.WriteString(strconv.FormatInt(log.Length, 10))
		}
		if log.Format != "" {
			sb.WriteString(" format ")
			sb.WriteString(log.Format)
		}
		if log.SampleRange != "" && log.SampleSize != 0 {
			sb.WriteString(" sample ")
			sb.WriteString(log.SampleRange)
			sb.WriteString(":")
			sb.WriteString(strconv.FormatInt(log.SampleSize, 10))
		}
		if log.Profile != "" {
			sb.WriteString(" profile ")
			sb.WriteString(log.Profile)
		}
		sb.WriteString(" ")
		sb.WriteString(log.Facility)
		if log.Level != "" {
			sb.WriteString(" ")
			sb.WriteString(log.Level)
			if log.MinLevel != "" {
				sb.WriteString(" ")
				sb.WriteString(log.MinLevel)
			}
		}
		result[index] = common.ReturnResultLine{
			Data:    sb.String(),
			Comment: log.Comment,
		}
	}
	return result, nil
}
