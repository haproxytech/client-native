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

import "github.com/haproxytech/client-native/v6/config-parser/errors"

// ExtractComment extracts the comment from the metadata map.
func ExtractComment(metadata map[string]any) string {
	var comment string
	if metadata == nil {
		return comment
	}
	if cmt, ok := metadata["comment"]; ok && cmt != "" {
		comment, _ = cmt.(string)
	}
	return comment
}

// ExtractCommentWithErr extracts the comment from the metadata map and returns an error if the comment is not a string.
func ExtractCommentWithErr(metadata map[string]any) (string, error) {
	var comment string
	if metadata == nil {
		return comment, nil
	}
	if cmt, ok := metadata["comment"]; ok && cmt != "" {
		comment, ok = cmt.(string)
		if !ok {
			return "", errors.ErrInvalidData
		}
	}
	return comment, nil
}
