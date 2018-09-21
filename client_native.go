package client_native

import (
	"log"

	"github.com/haproxytech/client-native/configuration"
	"github.com/haproxytech/client-native/stats"
)

// LogFunc - default log function is from the stdlib
var LogFunc func(string, ...interface{}) = log.Printf

// Default HAProxyClient using sane defaults
var Default = New(nil, nil)

// New HAProxyClient constructor
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

// HAProxyClient Native client for managing configuration and spitting out HAProxy stats
type HAProxyClient struct {
	Configuration configuration.Client
	Stats         stats.Client
}
