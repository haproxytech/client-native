package client_native

import (
	"log"

	"github.com/haproxytech/client-native/configuration"
	"github.com/haproxytech/client-native/stats"
)

// default log function is from the stdlib
var LogFunc func(string, ...interface{}) = log.Printf
var Default = New(nil, nil)

func New(configurationClient configuration.Client, statsClient stats.Client) *HAProxyClient {
	client := new(HAProxyClient)

	if configurationClient == nil {
		configurationClient = configuration.DefaultLBCTLClient()
	}

	if statsClient == nil {
		statsClient = stats.DefaultStatsClient()
	}

	client.Configuration = configurationClient
	client.Stats = statsClient

	return client
}

type HAProxyClient struct {
	Configuration configuration.Client
	Stats         stats.Client
}
