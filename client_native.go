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
	"context"
	"log"

	"github.com/haproxytech/client-native/v3/configuration"
	"github.com/haproxytech/client-native/v3/runtime"
	runtime_opt "github.com/haproxytech/client-native/v3/runtime/options"
	"github.com/haproxytech/client-native/v3/spoe"
	"github.com/haproxytech/client-native/v3/storage"
)

// LogFunc - default log function is from the stdlib
var LogFunc = log.Printf //nolint:gochecknoglobals

// DefaultClient with sane defaults
func DefaultClient() (*HAProxyClient, error) {
	c := &HAProxyClient{}

	err := c.Init(nil, nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Init HAProxyClient
func (c *HAProxyClient) Init(configurationClient *configuration.Client, runtimeClient runtime.Runtime) error {
	var err error
	if configurationClient == nil {
		configurationClient, err = configuration.DefaultClient()
		if err != nil {
			return err
		}
	}

	if runtimeClient == nil {
		runtimeClient, err = runtime.New(context.Background(), runtime_opt.SocketDefault())
		if err != nil {
			return err
		}
	}

	c.Configuration = configurationClient
	c.Runtime = runtimeClient
	return nil
}

// HAProxyClient Native client for managing configuration and spitting out HAProxy stats
type HAProxyClient struct {
	Configuration  *configuration.Client
	Runtime        runtime.Runtime
	MapStorage     storage.Storage
	SSLCertStorage storage.Storage
	Spoe           spoe.Spoe
}
