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

package types

//name:spoe-section
//no:sections
//no:init
//no:name
type SPOESection struct {
	Name    string
	Comment string
}

//name:event
//no:sections
//test:ok:event on-client-session
//test:ok:event on-client-session if ! { src -f /etc/haproxy/whitelist.lst }
//test:fail:event
type Event struct {
	Name     string
	Cond     string
	CondTest string
	Comment  string
}
