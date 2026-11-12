/*
Copyright 2026 HAProxy Technologies

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

package filterconfig

import (
	"strings"

	"github.com/haproxytech/client-native/v6/models"
)

// lua passes its arguments to the filter registered as lua.<id>; their meaning
// is defined by the Lua script.
type lua struct{}

func (lua) Parse(parts []string, config *models.FilterConfig) error {
	if len(parts) > 0 {
		config.LuaArgs = append([]string{}, parts...)
	}
	return nil
}

func (lua) Result(config models.FilterConfig) string {
	var result strings.Builder
	for _, arg := range config.LuaArgs {
		result.WriteString(" ")
		result.WriteString(arg)
	}
	return result.String()
}
