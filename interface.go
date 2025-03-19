package clientnative

import (
	"context"

	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/options"
	"github.com/haproxytech/client-native/v6/runtime"
	"github.com/haproxytech/client-native/v6/spoe"
	"github.com/haproxytech/client-native/v6/storage"
)

type HAProxyClient interface {
	Configuration() (configuration.Configuration, error)
	Runtime() (runtime.Runtime, error)
	ReplaceConfiguration(configurationClient configuration.Configuration)
	ReplaceRuntime(runtimeClient runtime.Runtime)
	MapStorage() (storage.Storage, error)
	SSLCertStorage() (storage.Storage, error)
	CrtListStorage() (storage.Storage, error)
	GeneralStorage() (storage.Storage, error)
	Spoe() (spoe.Spoe, error)
}

func New(ctx context.Context, opt ...options.Option) (HAProxyClient, error) { //nolint:revive
	o := options.Options{}
	var err error

	for _, option := range opt {
		err = option.Set(&o)
		if err != nil {
			return nil, err
		}
	}

	c := &haProxyClient{
		configuration:  o.Configuration,
		runtime:        o.Runtime,
		mapStorage:     o.MapStorage,
		sslCertStorage: o.SSLCertStorage,
		crtListStorage: o.CrtListStorage,
		generalStorage: o.GeneralStorage,
		spoe:           o.Spoe,
	}

	return c, nil
}
