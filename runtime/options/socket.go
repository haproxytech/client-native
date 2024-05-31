/*
Copyright 2022 HAProxy Technologies

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

type socket struct {
	path string
}

func Socket(path string) RuntimeOption {
	return socket{path}
}

// SocketDefault uses /var/run/haproxy.sock as socket path
func SocketDefault() RuntimeOption {
	return Socket("/var/run/haproxy.sock")
}

func (u socket) Set(p *RuntimeOptions) error {
	p.Socket = u.path
	return nil
}
