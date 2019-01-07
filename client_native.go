package client_native

import (
	"log"

	"github.com/haproxytech/client-native/configuration"
	"github.com/haproxytech/client-native/runtime"
)

// LogFunc - default log function is from the stdlib
var LogFunc func(string, ...interface{}) = log.Printf

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
func (c *HAProxyClient) Init(configurationClient *configuration.Client, runtimeClient *runtime.Client) error {
	var err error
	if configurationClient == nil {
		configurationClient, err = configuration.DefaultClient()
		if err != nil {
			return err
		}
	}

	if runtimeClient == nil {
		runtimeClient, err = runtime.DefaultClient()
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
	Configuration *configuration.Client
	Runtime       *runtime.Client
}
