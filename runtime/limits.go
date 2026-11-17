// Copyright 2026 HAProxy Technologies
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
	stderrors "errors"
	"fmt"
)

// SetGlobalMaxConn sets the process-wide maxconn
func (s *SingleRuntime) SetGlobalMaxConn(maxconn int) error {
	if maxconn < 0 {
		return stderrors.New("bad request")
	}
	return s.Execute(fmt.Sprintf("set maxconn global %d", maxconn))
}

// SetRateLimitConnectionsGlobal sets the process-wide connections per second limit
func (s *SingleRuntime) SetRateLimitConnectionsGlobal(value uint64) error {
	return s.Execute(fmt.Sprintf("set rate-limit connections global %d", value))
}

// SetRateLimitSessionsGlobal sets the process-wide sessions per second limit
func (s *SingleRuntime) SetRateLimitSessionsGlobal(value uint64) error {
	return s.Execute(fmt.Sprintf("set rate-limit sessions global %d", value))
}

// SetRateLimitHTTPCompressionGlobal sets the process-wide HTTP compression rate limit in kB/s
func (s *SingleRuntime) SetRateLimitHTTPCompressionGlobal(value uint64) error {
	return s.Execute(fmt.Sprintf("set rate-limit http-compression global %d", value))
}
