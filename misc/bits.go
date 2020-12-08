// Copyright 2019 HAProxy Technologies
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

package misc

import (
	"strconv"
)

type Bits uint64

const (
	B0 Bits = 1 << iota
	B1
	B2
	B3
	B4
	B5
	B6
	B7
)

func (b *Bits) Any(flags ...Bits) bool {
	for _, flag := range flags {
		if *b&flag != 0 {
			return true
		}
	}
	return false
}

// GetServerAdminState parses srv_admin_state, srv_admin_state is a mask
func GetServerAdminState(state string) (string, error) {
	// 23	0100011 - mask for maint
	// 18   0011000 - mask for drain
	val, err := strconv.ParseUint(state, 16, 64)
	if err != nil {
		return "", err
	}
	mask := Bits(val)
	if mask.Any(B0, B1, B5) {
		return "maint", nil
	}
	if mask.Any(B3, B4) {
		return "drain", nil
	}
	return "ready", nil
}
