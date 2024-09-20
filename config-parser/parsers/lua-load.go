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
	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type LuaLoad struct {
	data        []types.LuaLoad
	preComments []string // comments that appear before the actual line
}

func (l *LuaLoad) parse(line string, parts []string, comment string) (*types.LuaLoad, error) {
	if len(parts) < 2 {
		return nil, &errors.ParseError{Parser: "LuaLoad", Line: line}
	}
	lua := &types.LuaLoad{
		File:    parts[1],
		Comment: comment,
	}
	return lua, nil
}

func (l *LuaLoad) Result() ([]common.ReturnResultLine, error) {
	if len(l.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(l.data))
	for index, data := range l.data {
		result[index] = common.ReturnResultLine{
			Data:    "lua-load " + data.File,
			Comment: data.Comment,
		}
	}
	return result, nil
}
