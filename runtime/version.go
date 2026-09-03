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
	parts := strings.Split(version, "-")
	sv, err := semver.NewVersion(version)
	if err != nil {
		// Only use the first part.
		sv, err = semver.NewVersion(parts[0])
		if err != nil {
			return err
		}
	}
	v.Version = sv

	// Look for a git commit hash.
	if len(parts) > 1 {
		for _, part := range parts[1:] {
			if isGitTag(part) {
				v.Commit = part
				break
			}
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

// Matches git tags between 6 and 8 characters long.
func isGitTag(s string) bool {
	n := len(s)
	if n < 6 || n > 8 {
		return false
	}
	for _, c := range s {
		match := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')
		if !match {
			return false
		}
	}
	return true
}
