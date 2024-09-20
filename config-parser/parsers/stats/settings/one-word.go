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

package stats

import "errors"

type OneWord struct {
	Name    string
	Comment string
}

func (o *OneWord) Parse(parts []string, comment string) error {
	if len(parts) < 2 {
		return errors.New("not enough params")
	}

	o.Name = parts[1]
	o.Comment = comment
	return nil
}

func (o *OneWord) String() string {
	return o.Name
}

func (o *OneWord) GetComment() string {
	return o.Comment
}
