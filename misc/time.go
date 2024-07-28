// Copyright 2024 HAProxy Technologies
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

// HAProxy's Time suffixes
type TimeSuffix string

const (
	None    TimeSuffix = "none"
	Milli   TimeSuffix = "ms"
	Sec     TimeSuffix = "s"
	Min     TimeSuffix = "m"
	Hour    TimeSuffix = "h"
	Day     TimeSuffix = "d"
	Nearest TimeSuffix = Day

	micro  = -1000
	milli  = 1
	second = 1000
	minute = 60000
	hour   = 3600000
	day    = 86400000
)

// ParseTimeout returns the number of milliseconds in a timeout string.
func ParseTimeout(tOut string) *int64 {
	return parseTimeout(tOut, milli)
}

func ParseTimeoutDefaultSeconds(tOut string) *int64 {
	return parseTimeout(tOut, second)
}

func parseTimeout(tOut string, defaultMultiplier int64) *int64 {
	n := len(tOut)
	if n == 0 {
		return nil
	}

	multiplier := defaultMultiplier
	trim := 1

	switch tOut[n-1] {
	case 's':
		if n < 2 {
			return nil
		}
		switch tOut[n-2] {
		case 'u': // us
			multiplier = micro
			trim = 2
		case 'm': // ms
			multiplier = milli
			trim = 2
		default: // s
			multiplier = second
		}
	case 'm':
		multiplier = minute
	case 'h':
		multiplier = hour
	case 'd':
		multiplier = day
	default:
		trim = 0
	}

	v, err := strconv.ParseInt(tOut[:n-trim], 10, 64)
	if err != nil || v < 0 {
		return nil
	}

	if multiplier > 1 {
		v *= multiplier
	} else if multiplier == micro {
		if v >= 1000 {
			v /= 1000
		} else if v > 0 {
			v = 1
		}
	}

	return &v
}

// Serialize a number of milliseconds as per HAProxy's Time format.
func SerializeTime(ms int64, preferredSuffix string) string {
	switch TimeSuffix(preferredSuffix) {
	case Nearest:
		break
	case Milli:
		goto millisec
	case Sec:
		goto second
	case Min:
		goto minute
	case Hour:
		goto hour
	case None:
		fallthrough
	default:
		return itoa(ms)
	}

	if ms >= day && ms%day == 0 {
		return itoa(ms/day) + string(Day)
	}
hour:
	if ms >= hour && ms%hour == 0 {
		return itoa(ms/hour) + string(Hour)
	}
minute:
	if ms >= minute && ms%minute == 0 {
		return itoa(ms/minute) + string(Min)
	}
second:
	if ms >= second && ms%second == 0 {
		return itoa(ms/second) + string(Sec)
	}
millisec:
	return itoa(ms) + string(Milli)
}

func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}
