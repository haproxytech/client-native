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
package actions

import "io"

func CheckParsePair(parts []string, i *int, str *string) {
	if (*i + 1) < len(parts) {
		*str = parts[*i+1]
		*i++
	}
}

func CheckWritePair(sb io.StringWriter, key string, value string) {
	if value != "" {
		_, _ = sb.WriteString(" ")
		if key != "" {
			_, _ = sb.WriteString(key)
			_, _ = sb.WriteString(" ")
		}
		_, _ = sb.WriteString(value)
	}
}
