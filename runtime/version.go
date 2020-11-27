// Copyright 2020 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package runtime

import (
	"fmt"
	"strconv"
	"strings"
)

type HAProxyVersion struct {
	Major   int
	Minor   int
	Patch   int
	Commit  string
	Version string
}

func (v *HAProxyVersion) ParseHAProxyVersion(version string) error {
	v.Version = version

	parts := strings.SplitN(version, "-", 2)
	data := strings.SplitN(parts[0], ".", 3)
	major, err := strconv.Atoi(data[0])
	if err == nil {
		v.Major = major
	}
	if len(data) > 1 {
		minor, err := strconv.Atoi(data[1])
		if err == nil {
			v.Minor = minor
		}
	}
	if len(data) > 2 {
		patch, err := strconv.Atoi(data[2])
		if err == nil {
			v.Patch = patch
		}
	}
	if len(parts) < 2 {
		return fmt.Errorf("version is not in correct format [%s]", version)
	}
	data = strings.Split(parts[1], "-")
	if len(parts) < 2 {
		v.Commit = data[0]
	} else {
		v.Commit = data[len(data)-2]
	}
	return nil
}
