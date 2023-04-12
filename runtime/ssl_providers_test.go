package runtime

import (
	"reflect"
	"testing"

	"github.com/haproxytech/client-native/v6/models"
)

func TestSingleRuntime_ShowSSLProviders(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	tests := []struct {
		name           string
		fields         fields
		wantErr        bool
		socketResponse map[string]string
		responseValues *models.SslProviders
	}{
		{
			name:    "Show SSL Providers, should return no error",
			fields:  fields{socketPath: haProxy.Addr().String()},
			wantErr: false,
			socketResponse: map[string]string{
				"show ssl providers\n": `Loaded providers :
    - fips
    - base`,
			},
			responseValues: &models.SslProviders{
				Providers: []string{"fips", "base"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.masterWorkerMode)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			sslProviders, err := s.ShowSSLProviders()
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowSSLProviders() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(sslProviders, tt.responseValues) {
				t.Errorf("SingleRuntime.ShowSSLProviders() = %v, want %v", sslProviders, tt.responseValues)
			}
		})
	}
}
