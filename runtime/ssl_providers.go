package runtime

import (
	"fmt"
	"strings"

	native_errors "github.com/haproxytech/client-native/v6/errors"
	"github.com/haproxytech/client-native/v6/models"
)

func parseSSLProviders(response string) (*models.SslProviders, error) {
	if response == "" {
		return nil, native_errors.ErrNotFound
	}
	sslProviders := &models.SslProviders{}
	parts := strings.Split(response, "\n")
	for _, p := range parts[1:] {
		_, valueString, found := strings.Cut(p, "-")
		if !found {
			continue
		}
		sslProviders.Providers = append(sslProviders.Providers, strings.TrimSpace(valueString))
	}
	return sslProviders, nil
}

// ShowSSLProviders shows the names of the providers loaded by OpenSSL during init
func (s *SingleRuntime) ShowSSLProviders() (*models.SslProviders, error) {
	cmd := "show ssl providers\n"
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}

	return parseSSLProviders(response)
}
