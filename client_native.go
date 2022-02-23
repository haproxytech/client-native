// Copyright 2019 HAProxy Technologies
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

package clientnative

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/haproxytech/client-native/v3/configuration"
	"github.com/haproxytech/client-native/v3/runtime"
	"github.com/haproxytech/client-native/v3/spoe"
	"github.com/haproxytech/client-native/v3/storage"
)

// LogFunc - default log function is from the stdlib
var (
	LogFunc               = log.Printf //nolint:gochecknoglobals
	ErrOptionNotAvailable = errors.New("option is not available")
)

// HAProxyClient Native client for managing configuration and spitting out HAProxy stats
type haProxyClient struct {
	configuration   configuration.Configuration
	runtime         runtime.Runtime
	mapStorage      storage.Storage
	sslCertStorage  storage.Storage
	generalStorage  storage.Storage
	spoe            spoe.Spoe
	configurationMu sync.RWMutex
	runtimeMu       sync.RWMutex
}

func (c *haProxyClient) Configuration() (configuration.Configuration, error) {
	c.configurationMu.RLock()
	defer c.configurationMu.RUnlock()
	if c.configuration == nil {
		return nil, fmt.Errorf("configuration: %w", ErrOptionNotAvailable)
	}
	return c.configuration, nil
}

func (c *haProxyClient) ReplaceConfiguration(configurationClient configuration.Configuration) {
	c.configurationMu.Lock()
	defer c.configurationMu.Unlock()
	c.configuration = configurationClient
}

func (c *haProxyClient) ReplaceRuntime(runtime runtime.Runtime) {
	c.runtimeMu.Lock()
	defer c.runtimeMu.Unlock()
	c.runtime = runtime
}

func (c *haProxyClient) Runtime() (runtime.Runtime, error) {
	c.runtimeMu.RLock()
	defer c.runtimeMu.RUnlock()
	if c.runtime == nil {
		return nil, fmt.Errorf("runtime: %w", ErrOptionNotAvailable)
	}
	return c.runtime, nil
}

func (c *haProxyClient) MapStorage() (storage.Storage, error) {
	if c.mapStorage == nil {
		return nil, fmt.Errorf("map storage: %w", ErrOptionNotAvailable)
	}
	return c.mapStorage, nil
}

func (c *haProxyClient) SSLCertStorage() (storage.Storage, error) {
	if c.sslCertStorage == nil {
		return nil, fmt.Errorf("ssl cert storage: %w", ErrOptionNotAvailable)
	}
	return c.sslCertStorage, nil
}

func (c *haProxyClient) GeneralStorage() (storage.Storage, error) {
	if c.generalStorage == nil {
		return nil, fmt.Errorf("general files storage: %w", ErrOptionNotAvailable)
	}
	return c.generalStorage, nil
}

func (c *haProxyClient) Spoe() (spoe.Spoe, error) {
	if c.spoe == nil {
		return nil, fmt.Errorf("spoe: %w", ErrOptionNotAvailable)
	}
	return c.spoe, nil
}
