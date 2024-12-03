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

package types

import (
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/parsers/actions"
	tcpActions "github.com/haproxytech/client-native/v6/config-parser/parsers/tcp/actions"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type Content struct {
	Action  types.Action
	Comment string
}

func (f *Content) ParseAction(action types.Action, parts []string) error {
	if action.Parse(parts, types.TCP, "") != nil {
		return &errors.ParseError{Parser: "TCPRequestContent", Line: ""}
	}

	f.Action = action
	return nil
}

func (f *Content) Parse(parts []string, comment string) error {
	if comment != "" {
		f.Comment = comment
	}
	if len(parts) < 3 {
		return stderrors.New("not enough params")
	}
	var err error
	switch parts[2] {
	case "accept":
		err = f.ParseAction(&actions.Accept{}, parts)
	case "reject":
		err = f.ParseAction(&actions.Reject{}, parts)
	case "capture":
		err = f.ParseAction(&tcpActions.Capture{}, parts)
	case "set-priority-class":
		err = f.ParseAction(&actions.SetPriorityClass{}, parts)
	case "set-priority-offset":
		err = f.ParseAction(&actions.SetPriorityOffset{}, parts)
	case "set-dst":
		err = f.ParseAction(&actions.SetDst{}, parts)
	case "set-dst-port":
		err = f.ParseAction(&actions.SetDstPort{}, parts)
	case "silent-drop":
		err = f.ParseAction(&actions.SilentDrop{}, parts)
	case "send-spoe-group":
		err = f.ParseAction(&actions.SendSpoeGroup{}, parts)
	case "use-service":
		err = f.ParseAction(&actions.UseService{}, parts)
	case "set-bandwidth-limit":
		err = f.ParseAction(&actions.SetBandwidthLimit{}, parts)
	case "set-log-level":
		err = f.ParseAction(&actions.SetLogLevel{}, parts)
	case "set-mark":
		err = f.ParseAction(&actions.SetMark{}, parts)
	case "set-nice":
		err = f.ParseAction(&actions.SetNice{}, parts)
	case "set-tos":
		err = f.ParseAction(&actions.SetTos{}, parts)
	case "set-src-port":
		err = f.ParseAction(&actions.SetSrcPort{}, parts)
	case "switch-mode":
		err = f.ParseAction(&tcpActions.SwitchMode{}, parts)
	case "close":
		err = f.ParseAction(&tcpActions.Close{}, parts)
	case "set-bc-mark":
		err = f.ParseAction(&actions.SetBcMark{}, parts)
	case "set-bc-tos":
		err = f.ParseAction(&actions.SetBcTos{}, parts)
	case "set-fc-mark":
		err = f.ParseAction(&actions.SetFcMark{}, parts)
	case "set-fc-tos":
		err = f.ParseAction(&actions.SetFcTos{}, parts)
	case "set-retries":
		err = f.ParseAction(&actions.SetRetries{}, parts)
	case "do-log":
		err = f.ParseAction(&actions.DoLog{}, parts)
	default:
		switch {
		case strings.HasPrefix(parts[2], "track-sc"):
			err = f.ParseAction(&actions.TrackSc{}, parts)
		case strings.HasPrefix(parts[2], "lua."):
			err = f.ParseAction(&actions.Lua{}, parts)
		case strings.HasPrefix(parts[2], "sc-add-gpc("):
			err = f.ParseAction(&actions.ScAddGpc{}, parts)
		case strings.HasPrefix(parts[2], "sc-inc-gpc("):
			err = f.ParseAction(&actions.ScIncGpc{}, parts)
		case strings.HasPrefix(parts[2], "sc-inc-gpc0("):
			err = f.ParseAction(&actions.ScIncGpc0{}, parts)
		case strings.HasPrefix(parts[2], "sc-inc-gpc1("):
			err = f.ParseAction(&actions.ScIncGpc1{}, parts)
		case strings.HasPrefix(parts[2], "sc-set-gpt("):
			err = f.ParseAction(&actions.ScSetGpt{}, parts)
		case strings.HasPrefix(parts[2], "sc-set-gpt0"):
			err = f.ParseAction(&actions.ScSetGpt0{}, parts)
		case strings.HasPrefix(parts[2], "set-var("):
			err = f.ParseAction(&actions.SetVar{}, parts)
		case strings.HasPrefix(parts[2], "set-var-fmt"):
			err = f.ParseAction(&actions.SetVarFmt{}, parts)
		case strings.HasPrefix(parts[2], "unset-var"):
			err = f.ParseAction(&actions.UnsetVar{}, parts)
		case strings.HasPrefix(parts[2], "do-resolve"):
			err = f.ParseAction(&actions.DoResolve{}, parts)
		default:
			return fmt.Errorf("unsupported action %s", parts[2])
		}
	}
	return err
}

func (f *Content) String() string {
	var result strings.Builder

	result.WriteString("content")
	result.WriteString(" ")
	result.WriteString(f.Action.String())

	return result.String()
}

func (f *Content) GetComment() string {
	return f.Comment
}
