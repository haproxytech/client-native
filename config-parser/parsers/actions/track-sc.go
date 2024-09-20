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

package actions

import (
	stderrors "errors"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type TrackScT string

const (
	TrackScType TrackScT = "track-sc"
)

type TrackSc struct {
	Type         TrackScT
	StickCounter int64
	Key          string
	Table        string
	Comment      string
	Cond         string
	CondTest     string
}

func (f *TrackSc) Parse(parts []string, parserType types.ParserType, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	var data string
	var command []string
	switch parserType {
	case types.HTTP:
		if len(parts) < 3 {
			return stderrors.New("not enough params")
		}
		data = parts[1]
		command = parts[2:]
	case types.TCP:
		if len(parts) < 4 {
			return stderrors.New("not enough params")
		}
		data = parts[2]
		command = parts[3:]
	}

	f.Type = TrackScType
	counterS := strings.TrimPrefix(data, string(TrackScType))
	counter, err := strconv.ParseInt(counterS, 10, 64)
	if err != nil {
		return stderrors.New("failed to parse stick-counter")
	}
	f.StickCounter = counter

	return f.parseCommand(command)
}

func (f *TrackSc) parseCommand(command []string) error {
	// command contains only <key>
	if len(command) == 1 {
		f.Key = command[0]
		return nil
	}

	if len(command) < 3 {
		return stderrors.New("not enough params")
	}
	command, condition := common.SplitRequest(command)
	if len(command) > 1 && command[1] == "table" {
		if len(command) < 3 {
			return stderrors.New("not enough params")
		}
		f.Key = command[0]
		f.Table = command[2]
	}
	if len(command) == 1 {
		f.Key = command[0]
	}
	if len(condition) > 1 {
		f.Cond = condition[0]
		f.CondTest = strings.Join(condition[1:], " ")
	}
	return nil
}

func (f *TrackSc) String() string {
	var result strings.Builder
	result.WriteString(string(f.Type))
	result.WriteString(strconv.FormatInt(f.StickCounter, 10))
	result.WriteString(" ")
	result.WriteString(f.Key)
	if f.Table != "" {
		result.WriteString(" table ")
		result.WriteString(f.Table)
	}
	if f.Cond != "" {
		result.WriteString(" ")
		result.WriteString(f.Cond)
		result.WriteString(" ")
		result.WriteString(f.CondTest)
	}
	return result.String()
}

func (f *TrackSc) GetComment() string {
	return f.Comment
}
