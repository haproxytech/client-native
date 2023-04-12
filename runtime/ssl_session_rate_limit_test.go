package runtime

import (
	"testing"
)

func TestSingleRuntime_SetRateLimitSSLSessionGlobal(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	type args struct {
		payload uint64
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Set the SSL session global rate limit, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				payload: 2,
			},
			wantErr: false,
			socketResponse: map[string]string{
				"set rate-limit ssl-sessions global 2\n": ``,
			},
		},
		{
			name:   "Set the SSL session global rate limit with invalid value, should return error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				payload: 9999999,
			},
			wantErr: true,
			socketResponse: map[string]string{
				"set rate-limit ssl-sessions global 9999999\n": ` Value out of range.`,
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
			if err := s.SetRateLimitSSLSessionGlobal(tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.SetRateLimitSSLSessionGlobal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
