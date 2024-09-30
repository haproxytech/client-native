/*
Copyright 2021 HAProxy Technologies

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

package options

type useMd5Hash struct{}

func (u useMd5Hash) Set(p *Parser) error {
	p.UseMd5Hash = true
	return nil
}

// UseMd5Hash sets flag to use md5 hash
var UseMd5Hash = useMd5Hash{} //nolint:gochecknoglobals
