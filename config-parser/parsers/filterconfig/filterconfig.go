/*
Copyright 2026 HAProxy Technologies

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

// Package filterconfig holds the per-class grammars of the filter-config
// directive. The directive grammar itself is in parsers.FilterConfig.
package filterconfig

import "github.com/haproxytech/client-native/v6/models"

// Class parses the parameters of one filter-config class and serializes them
// back into the directive line.
type Class interface {
	Parse(parts []string, config *models.FilterConfig) error
	Result(config models.FilterConfig) string
}

type entry struct {
	parser Class
	hasID  bool
}

//nolint:gochecknoglobals
var classes = map[string]entry{
	"trace":         {trace{}, true},
	"http-comp-req": {compression{}, false},
	"http-comp-res": {compression{}, false},
	"decomp-req":    {decompression{}, false},
	"decomp-res":    {decompression{}, false},
	"spoe":          {spoe{}, true},
	"cache-store":   {cacheStore{}, true},
	"fcgi":          {fcgi{}, true},
	"bwlim-in":      {bandwidth{}, true},
	"bwlim-out":     {bandwidth{}, true},
	"lua":           {lua{}, true},
}

// Get returns the parser of a filter-config class.
func Get(class string) (Class, bool) {
	e, ok := classes[class]
	return e.parser, ok
}

// HasID reports whether the class requires an instance id.
func HasID(class string) bool {
	return classes[class].hasID
}
