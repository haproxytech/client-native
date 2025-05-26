package runtime

import (
	"testing"
	"time"
)

func TestEventListener(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		address string
		sink    string
		flags   []string
	}
	tests := []struct {
		name           string
		fields         fields
		want           string
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "acme newcert event",
			fields: fields{address: haProxy.Addr().String(), sink: "dpapi", flags: []string{"-w", "-0"}},
			want:   "acme newcert foobar.pem.rsa",
			socketResponse: map[string]string{
				"show events dpapi -w -0\n": "<0>2025-05-19T15:56:23.059755+02:00 acme newcert foobar.pem.rsa\n\x00",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			l, err := NewEventListener("unix", tt.fields.address, tt.fields.sink, time.Second, tt.fields.flags...)
			if err != nil {
				t.Errorf("NewEventListener() error = %v", err)
				return
			}
			got, err := l.Listen(t.Context())
			if (err != nil) != tt.wantErr {
				t.Errorf("EventListener.Listen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Message != tt.want {
				t.Errorf("EventListener.Listen() = %v, want %v", got, tt.want)
			}
			// Already closed anyway.
			if err = l.Close(); err != nil {
				t.Errorf("EventListener.Close() error = %v", err)
			}
		})
	}
}
