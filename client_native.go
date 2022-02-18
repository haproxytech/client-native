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
	configuration  configuration.Configuration
	runtime        runtime.Runtime
	mapStorage     storage.Storage
	sslCertStorage storage.Storage
	spoe           spoe.Spoe
}

func (c *haProxyClient) Configuration() (configuration.Configuration, error) {
	if c.configuration == nil {
		return nil, fmt.Errorf("configuration: %w", ErrOptionNotAvailable)
	}
	return c.configuration, nil
}

func (c *haProxyClient) ReplaceConfiguration(configurationClient configuration.Configuration) {
	c.configuration = configurationClient
}

func (c *haProxyClient) ReplaceRuntime(runtime runtime.Runtime) {
	c.runtime = runtime
}

func (c *haProxyClient) Runtime() (runtime.Runtime, error) {
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

func (c *haProxyClient) Spoe() (spoe.Spoe, error) {
	if c.spoe == nil {
		return nil, fmt.Errorf("spoe: %w", ErrOptionNotAvailable)
	}
	return c.spoe, nil
}
