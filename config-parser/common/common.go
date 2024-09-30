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

package common

import (
	"strings"
)

// StringSplitIgnoreEmpty while spliting, removes empty items
func StringSplitIgnoreEmpty(s string, separators ...rune) []string {
	f := func(c rune) bool {
		willSplit := false
		for _, sep := range separators {
			if c == sep {
				willSplit = true
				break
			}
		}
		return willSplit
	}
	return strings.FieldsFunc(s, f)
}

// StringSplitWithCommentIgnoreEmpty while splitting, removes empty items, if we have comment, separate it
func StringSplitWithCommentIgnoreEmpty(s string) (data []string, comment string) { //nolint:gocognit,nonamedreturns
	data = []string{}

	singleQuoteActive := false
	doubleQuoteActive := false
	var buff strings.Builder
	for index, c := range s {
		if !singleQuoteActive && !doubleQuoteActive {
			if (c == '#' && index == 0) || (c == '#' && s[index-1] != '\\') {
				if buff.Len() > 0 {
					data = append(data, buff.String())
					buff.Reset()
				}
				index++
				for ; index < len(s); index++ {
					if s[index] != ' ' {
						break
					}
				}
				comment = s[index:]
				return data, comment
			}
			if (c == ' ' && index == 0) || (c == ' ' && s[index-1] != '\\') || c == '\t' {
				if buff.Len() > 0 {
					data = append(data, buff.String())
					buff.Reset()
				}
				continue
			}
		}
		buff.WriteRune(c)
		if c == '"' {
			if doubleQuoteActive {
				if index == 0 || s[index-1] != '\\' {
					doubleQuoteActive = false
				}
			} else if !singleQuoteActive {
				if index == 0 || s[index-1] != '\\' {
					doubleQuoteActive = true
				}
			}
		}
		if c == '\'' {
			if singleQuoteActive {
				singleQuoteActive = false
			} else if !doubleQuoteActive {
				if index == 0 || s[index-1] != '\\' {
					singleQuoteActive = true
				}
			}
		}
	}
	if buff.Len() > 0 {
		data = append(data, buff.String())
	}
	return data, comment
}

// StringExtractComment checks if comment is added
func StringExtractComment(s string) string {
	p := StringSplitIgnoreEmpty(s, '#')
	if len(p) > 1 {
		return p[len(p)-1]
	}
	return ""
}

// SplitRequest searches for "if" or "unless" and returns result
func SplitRequest(parts []string) (command, condition []string) { //nolint:nonamedreturns
	if len(parts) == 0 {
		return []string{}, []string{}
	}
	index := 0
	found := false
	for index < len(parts) {
		switch parts[index] {
		case "if", "unless":
			found = true
		}
		if found {
			break
		}
		index++
	}
	command = parts[:index]
	condition = parts[index:]
	return command, condition
}

// CutRight slices string around the last occurrence of sep returning the text
// before and after sep. The found result reports whether sep appears in s. If
// sep does not appear in s, cut returns s, "", false.
func CutRight(s, sep string) (before, after string, found bool) { //nolint:nonamedreturns
	pos := strings.LastIndex(s, sep)
	if pos < 0 {
		before = s
		found = false
	} else {
		before = s[:pos]
		after = s[pos+len(sep):]
		found = true
	}
	return before, after, found
}
