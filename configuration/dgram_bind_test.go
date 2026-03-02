package configuration

import (
	"testing"

	"github.com/haproxytech/client-native/v6/config-parser/params"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/stretchr/testify/require"
)

func TestParseDgramBindIPv6(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		params          []params.DgramBindOption
		expectedAddress string
		expectedPort    *int64
		expectedName    string
	}{
		{
			name:            "bracketed IPv6 with port",
			path:            "[fd66:c3ec:c7fc::7c]:123",
			params:          []params.DgramBindOption{&params.BindOptionValue{Name: "name", Value: "ntp6"}},
			expectedAddress: "fd66:c3ec:c7fc::7c",
			expectedPort:    int64P(123),
			expectedName:    "ntp6",
		},
		{
			name:            "bracketed IPv6 with port and no name param",
			path:            "[2a04:4b07:21dc:218::1]:53",
			expectedAddress: "2a04:4b07:21dc:218::1",
			expectedPort:    int64P(53),
			expectedName:    "[2a04:4b07:21dc:218::1]:53",
		},
		{
			name:            "double colon with port",
			path:            ":::443",
			expectedAddress: "::",
			expectedPort:    int64P(443),
			expectedName:    ":::443",
		},
		{
			name:            "IPv4 with port",
			path:            "10.0.0.1:8080",
			params:          []params.DgramBindOption{&params.BindOptionValue{Name: "name", Value: "test"}},
			expectedAddress: "10.0.0.1",
			expectedPort:    int64P(8080),
			expectedName:    "test",
		},
		{
			name:            "wildcard with port",
			path:            ":443",
			expectedAddress: "",
			expectedPort:    int64P(443),
			expectedName:    ":443",
		},
		{
			name:            "unix socket path",
			path:            "/var/run/haproxy.sock",
			expectedAddress: "/var/run/haproxy.sock",
			expectedPort:    nil,
			expectedName:    "/var/run/haproxy.sock",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ondisk := types.DgramBind{
				Path:   tc.path,
				Params: tc.params,
			}
			result := ParseDgramBind(ondisk)
			require.NotNil(t, result, "ParseDgramBind returned nil for path %q", tc.path)
			require.Equal(t, tc.expectedAddress, result.Address, "address mismatch")
			if tc.expectedPort != nil {
				require.NotNil(t, result.Port, "expected port to be set")
				require.Equal(t, *tc.expectedPort, *result.Port, "port mismatch")
			} else {
				require.Nil(t, result.Port, "expected port to be nil")
			}
			require.Equal(t, tc.expectedName, result.Name, "name mismatch")
		})
	}
}

func TestSerializeDgramBindIPv6(t *testing.T) {
	tests := []struct {
		name         string
		dgramBind    models.DgramBind
		expectedPath string
	}{
		{
			name: "IPv6 address with port",
			dgramBind: models.DgramBind{
				Address: "fd66:c3ec:c7fc::7c",
				Port:    int64P(123),
				Name:    "ntp6",
			},
			expectedPath: "[fd66:c3ec:c7fc::7c]:123",
		},
		{
			name: "IPv6 double colon with port",
			dgramBind: models.DgramBind{
				Address: "::",
				Port:    int64P(443),
				Name:    "test",
			},
			expectedPath: "[::]:443",
		},
		{
			name: "IPv4 address with port",
			dgramBind: models.DgramBind{
				Address: "10.0.0.1",
				Port:    int64P(8080),
				Name:    "test",
			},
			expectedPath: "10.0.0.1:8080",
		},
		{
			name: "address only, no port",
			dgramBind: models.DgramBind{
				Address: "/var/run/haproxy.sock",
				Name:    "test",
			},
			expectedPath: "/var/run/haproxy.sock",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := SerializeDgramBind(tc.dgramBind)
			require.Equal(t, tc.expectedPath, result.Path, "serialized path mismatch")
		})
	}
}

func TestParseDgramBindRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"IPv6 bracketed", "[fd66:c3ec:c7fc::7c]:123"},
		{"IPv6 full", "[2a04:4b07:21dc:218::1]:53"},
		{"IPv6 double colon", "[::]:443"},
		{"IPv4", "10.0.0.1:8080"},
		{"wildcard", ":443"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ondisk := types.DgramBind{
				Path:   tc.path,
				Params: []params.DgramBindOption{&params.BindOptionValue{Name: "name", Value: "roundtrip"}},
			}
			parsed := ParseDgramBind(ondisk)
			require.NotNil(t, parsed, "ParseDgramBind returned nil for path %q", tc.path)

			serialized := SerializeDgramBind(*parsed)
			require.Equal(t, tc.path, serialized.Path, "round-trip path mismatch")
		})
	}
}

func int64P(v int64) *int64 {
	return &v
}
