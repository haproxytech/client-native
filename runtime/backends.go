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

package runtime

import (
	stderrors "errors"
	"fmt"
	"time"
)

// AddBackend creates an unpublished backend that takes all its settings from the named defaults section
func (s *SingleRuntime) AddBackend(name, defaults, mode, guid string) error {
	if name == "" || defaults == "" {
		return stderrors.New("bad request")
	}
	cmd := fmt.Sprintf("add backend %s from %s", name, defaults)
	if mode != "" {
		cmd += " mode " + mode
	}
	if guid != "" {
		cmd += " guid " + guid
	}
	return s.Execute(cmd)
}

// DeleteBackend removes a backend; it must be unpublished, without servers and not referenced by the configuration
func (s *SingleRuntime) DeleteBackend(name string) error {
	if name == "" {
		return stderrors.New("bad request")
	}
	return s.Execute("del backend " + name)
}

// PublishBackend makes the backend available to content switching rules
func (s *SingleRuntime) PublishBackend(name string) error {
	if name == "" {
		return stderrors.New("bad request")
	}
	return s.Execute("publish backend " + name)
}

// UnpublishBackend makes the backend unavailable to content switching rules
func (s *SingleRuntime) UnpublishBackend(name string) error {
	if name == "" {
		return stderrors.New("bad request")
	}
	return s.Execute("unpublish backend " + name)
}

// WaitBackendRemovable waits up to timeout for the backend to become removable by DeleteBackend
func (s *SingleRuntime) WaitBackendRemovable(name string, timeout time.Duration) error {
	if name == "" || timeout < time.Millisecond {
		return stderrors.New("bad request")
	}
	return s.Execute(fmt.Sprintf("wait %dms be-removable %s", timeout.Milliseconds(), name))
}
