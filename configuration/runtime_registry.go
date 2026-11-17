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

package configuration

// registry lists everything that can be applied through the runtime API
// instead of a reload. To make another attribute or object runtime-capable,
// add an entry here; the engine in runtime_apply.go needs no change.
//
// A plain attribute gets an apply function:
//
//	parser.Frontends: {
//		"maxconn": {apply: func(c change) ([]step, bool) { ... }},
//	},
//
// An object list (named lines such as servers) gets an object spec:
//
//	parser.Backends: {
//		"server": {object: &objectSpec{
//			items:  serverItems,
//			fields: map[string]itemApplier{"weight": setServerWeight},
//			add:    addServer, // nil = reload
//			del:    delServer, // nil = reload
//		}},
//	},
var registry = runtimeRegistry{} //nolint:gochecknoglobals
