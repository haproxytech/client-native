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

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/parsers/actions"
	httpActions "github.com/haproxytech/client-native/v5/config-parser/parsers/http/actions"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type Requests struct {
	Name        string
	Mode        string
	data        []types.Action
	preComments []string // comments that appear before the actual line
}

func (h *Requests) Init() {
	h.Name = "http-request"
	h.data = []types.Action{}
}

func (h *Requests) ParseHTTPRequest(request types.Action, parts []string, comment string) error {
	err := request.Parse(parts, types.HTTP, comment)
	if err != nil {
		return &errors.ParseError{Parser: "HTTPRequestLines", Line: ""}
	}
	h.data = append(h.data, request)
	return nil
}

func (h *Requests) Parse(line string, parts []string, comment string) (string, error) { //nolint:gocyclo,cyclop,maintidx
	if len(parts) >= 2 && parts[0] == "http-request" {
		var err error
		switch parts[1] {
		case "add-header":
			err = h.ParseHTTPRequest(&httpActions.AddHeader{}, parts, comment)
		case "allow":
			err = h.ParseHTTPRequest(&httpActions.Allow{}, parts, comment)
		case "auth":
			err = h.ParseHTTPRequest(&httpActions.Auth{}, parts, comment)
		case "cache-use":
			err = h.ParseHTTPRequest(&httpActions.CacheUse{}, parts, comment)
		case "capture":
			err = h.ParseHTTPRequest(&httpActions.Capture{}, parts, comment)
		case "del-header":
			err = h.ParseHTTPRequest(&httpActions.DelHeader{}, parts, comment)
		case "deny":
			err = h.ParseHTTPRequest(&httpActions.Deny{}, parts, comment)
		case "disable-l7-retry":
			err = h.ParseHTTPRequest(&httpActions.DisableL7Retry{}, parts, comment)
		case "early-hint":
			err = h.ParseHTTPRequest(&httpActions.EarlyHint{}, parts, comment)
		case "normalize-uri":
			err = h.ParseHTTPRequest(&httpActions.NormalizeURI{}, parts, comment)
		case "redirect":
			err = h.ParseHTTPRequest(&httpActions.Redirect{}, parts, comment)
		case "reject":
			err = h.ParseHTTPRequest(&actions.Reject{}, parts, comment)
		case "replace-header":
			err = h.ParseHTTPRequest(&httpActions.ReplaceHeader{}, parts, comment)
		case "replace-path":
			err = h.ParseHTTPRequest(&httpActions.ReplacePath{}, parts, comment)
		case "replace-pathq":
			err = h.ParseHTTPRequest(&httpActions.ReplacePathQ{}, parts, comment)
		case "replace-uri":
			err = h.ParseHTTPRequest(&httpActions.ReplaceURI{}, parts, comment)
		case "replace-value":
			err = h.ParseHTTPRequest(&httpActions.ReplaceValue{}, parts, comment)
		case "return":
			err = h.ParseHTTPRequest(&httpActions.Return{}, parts, comment)
		case "send-spoe-group":
			err = h.ParseHTTPRequest(&actions.SendSpoeGroup{}, parts, comment)
		case "set-dst":
			err = h.ParseHTTPRequest(&actions.SetDst{}, parts, comment)
		case "set-dst-port":
			err = h.ParseHTTPRequest(&actions.SetDstPort{}, parts, comment)
		case "set-header":
			err = h.ParseHTTPRequest(&httpActions.SetHeader{}, parts, comment)
		case "set-log-level":
			err = h.ParseHTTPRequest(&actions.SetLogLevel{}, parts, comment)
		case "set-mark":
			err = h.ParseHTTPRequest(&actions.SetMark{}, parts, comment)
		case "set-method":
			err = h.ParseHTTPRequest(&httpActions.SetMethod{}, parts, comment)
		case "set-nice":
			err = h.ParseHTTPRequest(&actions.SetNice{}, parts, comment)
		case "set-path":
			err = h.ParseHTTPRequest(&httpActions.SetPath{}, parts, comment)
		case "set-pathq":
			err = h.ParseHTTPRequest(&httpActions.SetPathQ{}, parts, comment)
		case "set-priority-class":
			err = h.ParseHTTPRequest(&actions.SetPriorityClass{}, parts, comment)
		case "set-priority-offset":
			err = h.ParseHTTPRequest(&actions.SetPriorityOffset{}, parts, comment)
		case "set-query":
			err = h.ParseHTTPRequest(&httpActions.SetQuery{}, parts, comment)
		case "set-src":
			err = h.ParseHTTPRequest(&httpActions.SetSrc{}, parts, comment)
		case "set-src-port":
			err = h.ParseHTTPRequest(&actions.SetSrcPort{}, parts, comment)
		case "set-timeout":
			err = h.ParseHTTPRequest(&httpActions.SetTimeout{}, parts, comment)
		case "set-tos":
			err = h.ParseHTTPRequest(&actions.SetTos{}, parts, comment)
		case "set-uri":
			err = h.ParseHTTPRequest(&httpActions.SetURI{}, parts, comment)
		case "silent-drop":
			err = h.ParseHTTPRequest(&actions.SilentDrop{}, parts, comment)
		case "strict-mode":
			err = h.ParseHTTPRequest(&httpActions.StrictMode{}, parts, comment)
		case "tarpit":
			err = h.ParseHTTPRequest(&httpActions.Tarpit{}, parts, comment)
		case "use-service":
			err = h.ParseHTTPRequest(&actions.UseService{}, parts, comment)
		case "wait-for-body":
			err = h.ParseHTTPRequest(&httpActions.WaitForBody{}, parts, comment)
		case "wait-for-handshake":
			err = h.ParseHTTPRequest(&httpActions.WaitForHandshake{}, parts, comment)
		case "set-bandwidth-limit":
			err = h.ParseHTTPRequest(&actions.SetBandwidthLimit{}, parts, comment)
		default:
			switch {
			case strings.HasPrefix(parts[1], "track-sc"):
				err = h.ParseHTTPRequest(&actions.TrackSc{}, parts, comment)
			case strings.HasPrefix(parts[1], "add-acl("):
				err = h.ParseHTTPRequest(&httpActions.AddACL{}, parts, comment)
			case strings.HasPrefix(parts[1], "del-acl("):
				err = h.ParseHTTPRequest(&httpActions.DelACL{}, parts, comment)
			case strings.HasPrefix(parts[1], "set-map("):
				err = h.ParseHTTPRequest(&httpActions.SetMap{}, parts, comment)
			case strings.HasPrefix(parts[1], "del-map("):
				err = h.ParseHTTPRequest(&httpActions.DelMap{}, parts, comment)
			case strings.HasPrefix(parts[1], "do-resolve("):
				err = h.ParseHTTPRequest(&actions.DoResolve{}, parts, comment)
			case strings.HasPrefix(parts[1], "lua."):
				err = h.ParseHTTPRequest(&actions.Lua{}, parts, comment)
			case strings.HasPrefix(parts[1], "sc-add-gpc("):
				err = h.ParseHTTPRequest(&actions.ScAddGpc{}, parts, comment)
			case strings.HasPrefix(parts[1], "sc-inc-gpc("):
				err = h.ParseHTTPRequest(&actions.ScIncGpc{}, parts, comment)
			case strings.HasPrefix(parts[1], "sc-inc-gpc0("):
				err = h.ParseHTTPRequest(&actions.ScIncGpc0{}, parts, comment)
			case strings.HasPrefix(parts[1], "sc-inc-gpc1("):
				err = h.ParseHTTPRequest(&actions.ScIncGpc1{}, parts, comment)
			case strings.HasPrefix(parts[1], "sc-set-gpt("):
				err = h.ParseHTTPRequest(&actions.ScSetGpt{}, parts, comment)
			case strings.HasPrefix(parts[1], "sc-set-gpt0("):
				err = h.ParseHTTPRequest(&actions.ScSetGpt0{}, parts, comment)
			case strings.HasPrefix(parts[1], "set-var("):
				err = h.ParseHTTPRequest(&actions.SetVar{}, parts, comment)
			case strings.HasPrefix(parts[1], "set-var-fmt("):
				err = h.ParseHTTPRequest(&actions.SetVarFmt{}, parts, comment)
			case strings.HasPrefix(parts[1], "unset-var("):
				err = h.ParseHTTPRequest(&actions.UnsetVar{}, parts, comment)
			default:
				return "", &errors.ParseError{Parser: "HTTPRequestLines", Line: line}
			}
		}
		return "", err
	}
	return "", &errors.ParseError{Parser: "HTTPRequestLines", Line: line}
}

func (h *Requests) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, req := range h.data {
		result[index] = common.ReturnResultLine{
			Data:    "http-request " + req.String(),
			Comment: req.GetComment(),
		}
	}
	return result, nil
}
