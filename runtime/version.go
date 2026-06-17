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
	"strings"

	"github.com/Masterminds/semver/v3"
)

type HAProxyVersion struct {
	*semver.Version

	Commit string
}

func (v *HAProxyVersion) ParseHAProxyVersion(version string) error {
	sv, err := semver.NewVersion(version)
	if err != nil {
		return err
	}
	v.Version = sv
	// Commit lives in the prerelease tail: "dev6-2f6f36-13" → "2f6f36",
	// "a42e6c-11" → "a42e6c", "34b2b10" → "34b2b10".
	if pre := sv.Prerelease(); pre != "" {
		parts := strings.Split(pre, "-")
		if len(parts) < 2 {
			v.Commit = parts[0]
		} else {
			v.Commit = parts[len(parts)-2]
		}
	}
	return nil
}

func IsBiggerOrEqual(minimum, current *HAProxyVersion) bool {
	if current == nil || current.Version == nil {
		return false
	}
	if minimum == nil || minimum.Version == nil {
		return true
	}
	return current.GreaterThanEqual(minimum.Version)
}
