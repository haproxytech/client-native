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

package http

import (
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/parsers/actions"
	httpActions "github.com/haproxytech/client-native/v6/config-parser/parsers/http/actions"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type Responses struct {
	Name        string
	Mode        string
	data        []types.Action
	preComments []string // comments that appear before the actual line
}

func (h *Responses) Init() {
	h.Name = "http-response"
	h.data = []types.Action{}
}

func (h *Responses) ParseHTTPResponse(response types.Action, parts []string, comment string) error {
	err := response.Parse(parts, types.HTTP, comment)
	if err != nil {
		return &errors.ParseError{Parser: "HTTPResponseLines", Line: ""}
	}
	h.data = append(h.data, response)
	return nil
}

func (h *Responses) Parse(line string, parts []string, comment string) (string, error) { //nolint:cyclop,gocyclo
	if len(parts) >= 2 && parts[0] == "http-response" {
		var err error
		switch parts[1] {
		case "add-header":
			err = h.ParseHTTPResponse(&httpActions.AddHeader{}, parts, comment)
		case "allow":
			err = h.ParseHTTPResponse(&httpActions.Allow{}, parts, comment)
		case "cache-store":
			err = h.ParseHTTPResponse(&httpActions.CacheStore{}, parts, comment)
		case "capture":
			err = h.ParseHTTPResponse(&httpActions.Capture{}, parts, comment)
		case "del-header":
			err = h.ParseHTTPResponse(&httpActions.DelHeader{}, parts, comment)
		case "deny":
			err = h.ParseHTTPResponse(&httpActions.Deny{}, parts, comment)
		case "pause":
			err = h.ParseHTTPResponse(&httpActions.Pause{}, parts, comment)
		case "redirect":
			err = h.ParseHTTPResponse(&httpActions.Redirect{}, parts, comment)
		case "replace-header":
			err = h.ParseHTTPResponse(&httpActions.ReplaceHeader{}, parts, comment)
		case "replace-value":
			err = h.ParseHTTPResponse(&httpActions.ReplaceValue{}, parts, comment)
		case "return":
			err = h.ParseHTTPResponse(&httpActions.Return{}, parts, comment)
		case "send-spoe-group":
			err = h.ParseHTTPResponse(&actions.SendSpoeGroup{}, parts, comment)
		case "set-header":
			err = h.ParseHTTPResponse(&httpActions.SetHeader{}, parts, comment)
		case "set-log-level":
			err = h.ParseHTTPResponse(&actions.SetLogLevel{}, parts, comment)
		case "set-mark":
			err = h.ParseHTTPResponse(&actions.SetMark{}, parts, comment)
		case "set-nice":
			err = h.ParseHTTPResponse(&actions.SetNice{}, parts, comment)
		case "set-status":
			err = h.ParseHTTPResponse(&httpActions.SetStatus{}, parts, comment)
		case "set-timeout":
			err = h.ParseHTTPResponse(&httpActions.SetTimeout{}, parts, comment)
		case "set-tos":
			err = h.ParseHTTPResponse(&actions.SetTos{}, parts, comment)
		case "silent-drop":
			err = h.ParseHTTPResponse(&actions.SilentDrop{}, parts, comment)
		case "strict-mode":
			err = h.ParseHTTPResponse(&httpActions.StrictMode{}, parts, comment)
		case "wait-for-body":
			err = h.ParseHTTPResponse(&httpActions.WaitForBody{}, parts, comment)
		case "set-bandwidth-limit":
			err = h.ParseHTTPResponse(&actions.SetBandwidthLimit{}, parts, comment)
		case "set-fc-mark":
			err = h.ParseHTTPResponse(&actions.SetFcMark{}, parts, comment)
		case "set-fc-tos":
			err = h.ParseHTTPResponse(&actions.SetFcTos{}, parts, comment)
		case "do-log":
			err = h.ParseHTTPResponse(&actions.DoLog{}, parts, comment)
		default:
			switch {
			case strings.HasPrefix(parts[1], "track-sc"):
				err = h.ParseHTTPResponse(&actions.TrackSc{}, parts, comment)
			case strings.HasPrefix(parts[1], "add-acl("):
				err = h.ParseHTTPResponse(&httpActions.AddACL{}, parts, comment)
			case strings.HasPrefix(parts[1], "del-acl("):
				err = h.ParseHTTPResponse(&httpActions.DelACL{}, parts, comment)
			case strings.HasPrefix(parts[1], "lua."):
				err = h.ParseHTTPResponse(&actions.Lua{}, parts, comment)
			case strings.HasPrefix(parts[1], "sc-add-gpc("):
				err = h.ParseHTTPResponse(&actions.ScAddGpc{}, parts, comment)
			case strings.HasPrefix(parts[1], "sc-inc-gpc("):
				err = h.ParseHTTPResponse(&actions.ScIncGpc{}, parts, comment)
			case strings.HasPrefix(parts[1], "sc-inc-gpc0("):
				err = h.ParseHTTPResponse(&actions.ScIncGpc0{}, parts, comment)
			case strings.HasPrefix(parts[1], "sc-inc-gpc1("):
				err = h.ParseHTTPResponse(&actions.ScIncGpc1{}, parts, comment)
			case strings.HasPrefix(parts[1], "sc-set-gpt("):
				err = h.ParseHTTPResponse(&actions.ScSetGpt{}, parts, comment)
			case strings.HasPrefix(parts[1], "sc-set-gpt0("):
				err = h.ParseHTTPResponse(&actions.ScSetGpt0{}, parts, comment)
			case strings.HasPrefix(parts[1], "set-map("):
				err = h.ParseHTTPResponse(&httpActions.SetMap{}, parts, comment)
			case strings.HasPrefix(parts[1], "del-map("):
				err = h.ParseHTTPResponse(&httpActions.DelMap{}, parts, comment)
			case strings.HasPrefix(parts[1], "set-var("):
				err = h.ParseHTTPResponse(&actions.SetVar{}, parts, comment)
			case strings.HasPrefix(parts[1], "set-var-fmt("):
				err = h.ParseHTTPResponse(&actions.SetVarFmt{}, parts, comment)
			case strings.HasPrefix(parts[1], "unset-var("):
				err = h.ParseHTTPResponse(&actions.UnsetVar{}, parts, comment)
			default:
				return "", &errors.ParseError{Parser: "HTTPResponseLines", Line: line}
			}
		}
		return "", err
	}
	return "", &errors.ParseError{Parser: "HTTPResponseLines", Line: line}
}

func (h *Responses) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, res := range h.data {
		result[index] = common.ReturnResultLine{
			Data:    "http-response " + res.String(),
			Comment: res.GetComment(),
		}
	}
	return result, nil
}
