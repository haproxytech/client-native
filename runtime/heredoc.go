// Copyright 2025 HAProxy Technologies
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

package runtime

import "strings"

// terminateHeredocPayload returns payload guaranteed to end with "\n".
//
// HAProxy's `<<` heredoc terminates only on a blank line. The conventional
// command shape used in this package is "<command> <<\n<payload>\n", so for
// the heredoc to terminate, payload itself must end with "\n" (producing
// "\n\n" — a blank line — at the end of the command).
//
// Without that, the outer framing in runtime_single_client.go that appends
// ";quit\n" (master-socket / worker > 0 path) places ";quit" immediately
// after the payload's lone trailing newline — turning what should be the
// blank-line terminator into a non-empty command line — and HAProxy then
// blocks reading the socket until the deadline fires.
func terminateHeredocPayload(payload string) string {
	if strings.HasSuffix(payload, "\n") {
		return payload
	}
	return payload + "\n"
}
