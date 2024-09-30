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
package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/haproxytech/client-native/v5/config-parser/parsers/http"
)

func TestRedirect(t *testing.T) {
	data := [][2]string{
		{
			`http-request redirect code 301 location https://site.d1.randomsite.com/healthcheck.html if { path_beg -m beg -i /health-check.html AND hdr_beg(host) -i www }`,
			`http-request redirect location https://site.d1.randomsite.com/healthcheck.html code 301 if { path_beg -m beg -i /health-check.html AND hdr_beg(host) -i www }`,
		},
		{
			`http-request redirect scheme https code 301 if http u_login`,
			`http-request redirect scheme https code 301 if http u_login`,
		},
		{
			`http-request redirect code 301 prefix https://%[req.hdr(Host)] set-cookie SEEN=1 if !cookie_set`,
			`http-request redirect prefix https://%[req.hdr(Host)] code 301 set-cookie SEEN=1 if !cookie_set`,
		},
	}
	parser := &http.Requests{}

	for _, d := range data {

		line := strings.TrimSpace(d[0])
		expected := d[1]

		parser.Init()
		err := ProcessLine(line, parser)
		if err != nil {
			t.Errorf(err.Error())
		}

		result, err := parser.Result()
		if err != nil {
			t.Errorf(err.Error())
		}

		var actual string

		if result[0].Comment == "" {
			actual = result[0].Data
		} else {
			actual = fmt.Sprintf("%s # %s", result[0].Data, result[0].Comment)
		}

		if expected != actual {
			t.Errorf(fmt.Sprintf("error: has [%s] expects [%s]", actual, expected))
		}
	}
}
