package runtime

import (
	"reflect"
	"testing"

	"github.com/haproxytech/client-native/v5/misc"
)

func TestSingleRuntime_ShowCrtLists(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	tests := []struct {
		name           string
		fields         fields
		want           CrtLists
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Simple show crt-list files, should return a file",
			fields: fields{socketPath: haProxy.Addr().String()},
			want: CrtLists{
				&CrtList{
					File: "/etc/haproxy/crt-list",
				},
			},
			socketResponse: map[string]string{
				"show ssl crt-list\n": ` /etc/haproxy/crt-list
				`,
			},
		},
		{
			name:   "Simple show crt-list files, should return a nothing",
			fields: fields{socketPath: haProxy.Addr().String()},
			want:   nil,
			socketResponse: map[string]string{
				"show ssl crt-list\n": `
				`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.process, tt.fields.worker)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			got, err := s.ShowCrtLists()
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowCrtLists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("SingleRuntime.ShowCrtLists() = %v, want %v", got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSingleRuntime_GetCrtList(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		file string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *CrtList
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Get specific crt-list files, should return a file",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				file: "/etc/haproxy/crt-list",
			},
			want: &CrtList{
				File: "/etc/haproxy/crt-list",
			},
			socketResponse: map[string]string{
				"show ssl crt-list\n": ` /etc/haproxy/crt-list
				`,
			},
		},
		{
			name:   "Get a not known crt-list files, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				file: "/etc/haproxy/not-known-list",
			},
			want:    nil,
			wantErr: true,
			socketResponse: map[string]string{
				"show ssl crt-list\n": ` /etc/haproxy/crt-list
				`,
			},
		},
		{
			name:   "Get a no crt-list files, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				file: "/etc/haproxy/crt-list",
			},
			want:    nil,
			wantErr: true,
			socketResponse: map[string]string{
				"show ssl crt-list\n": `
				`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.process, tt.fields.worker)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			got, err := s.GetCrtList(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.GetCrtList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SingleRuntime.GetCrtList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleRuntime_ShowCrtListEntries(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		file string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           CrtListEntries
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Get crt-list entries of crt-list file, should return 3 entries",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				file: "/etc/haproxy/crt-list",
			},
			want: CrtListEntries{
				&CrtListEntry{
					LineNumber: 1,
					File:       "/etc/ssl/cert-0.pem",
					SNIFilter: []string{
						"!*.crt-test.platform.domain.com",
						"!connectivitynotification.platform.domain.com",
						"!connectivitytunnel.platform.domain.com",
						"!authentication.cert.another.domain.com",
						"!*.authentication.cert.another.domain.com",
					},
				},
				&CrtListEntry{
					LineNumber:    2,
					File:          "/etc/ssl/cert-1.pem",
					SSLBindConfig: "verify optional ca-file /etc/ssl/ca-file-1.pem",
					SNIFilter: []string{
						"*.crt-test.platform.domain.com",
						"!connectivitynotification.platform.domain.com",
					},
				},
				&CrtListEntry{
					LineNumber:    4,
					File:          "/etc/ssl/cert-2.pem",
					SSLBindConfig: "verify required ca-file /etc/ssl/ca-file-2.pem",
					SNIFilter:     []string{},
				},
			},
			socketResponse: map[string]string{
				"show ssl crt-list -n /etc/haproxy/crt-list\n": ` # /etc/ssl/crt-list
					/etc/ssl/cert-0.pem:1 !*.crt-test.platform.domain.com !connectivitynotification.platform.domain.com !connectivitytunnel.platform.domain.com !authentication.cert.another.domain.com !*.authentication.cert.another.domain.com
					/etc/ssl/cert-1.pem:2 [verify optional ca-file /etc/ssl/ca-file-1.pem] *.crt-test.platform.domain.com !connectivitynotification.platform.domain.com
					/etc/ssl/cert-2.pem:4 [verify required ca-file /etc/ssl/ca-file-2.pem]
				`,
			},
		},
		{
			name:   "Get crt-list entries of crt-list file, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				file: "/etc/haproxy/not_known_list",
			},
			want:    nil,
			wantErr: true,
			socketResponse: map[string]string{
				"show ssl crt-list -n /etc/haproxy/not_known_list\n": ` didn't find the specified filename
				`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.process, tt.fields.worker)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			got, err := s.ShowCrtListEntries(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowCrtListEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("SingleRuntime.ShowCrtListEntries() = %v, want %v", got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSingleRuntime_AddCrtListEntry(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		crtList string
		entry   CrtListEntry
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "add crt-list entries to crt-list file, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crtList: "/etc/haproxy/crt-list",
				entry: CrtListEntry{
					File:          "/etc/ssl/cert-0.pem",
					SSLBindConfig: "alpn h2",
					SNIFilter: []string{
						"test.domain.com",
					},
				},
			},
			wantErr: false,
			socketResponse: map[string]string{
				"add ssl crt-list /etc/haproxy/crt-list <<\n/etc/ssl/cert-0.pem [alpn h2] test.domain.com\n": ` Inserting certificate '/etc/ssl/cert-0.pem' in crt-list '/etc/ssl/crt-list'.
				Success!
				`,
			},
		},
		{
			name:   "add crt-list entries to crt-list file without SSLBindConfig, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crtList: "/etc/haproxy/crt-list",
				entry: CrtListEntry{
					File: "/etc/ssl/cert-0.pem",
					SNIFilter: []string{
						"test.domain.com",
					},
				},
			},
			wantErr: false,
			socketResponse: map[string]string{
				"add ssl crt-list /etc/haproxy/crt-list <<\n/etc/ssl/cert-0.pem test.domain.com\n": ` Inserting certificate '/etc/ssl/cert-0.pem' in crt-list '/etc/ssl/crt-list'.
				Success!
				`,
			},
		},
		{
			name:   "add crt-list entries to crt-list file with a not known pem, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crtList: "/etc/haproxy/crt-list",
				entry: CrtListEntry{
					File:      "/etc/ssl/not_known.pem",
					SNIFilter: []string{},
				},
			},
			wantErr: true,
			socketResponse: map[string]string{
				"add ssl crt-list /etc/haproxy/crt-list <<\n/etc/ssl/not_known.pem\n": ` Can't edit the crt-list: certificate '/etc/ssl/cert-26.pem' does not exist!
				`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.process, tt.fields.worker)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			if err := s.AddCrtListEntry(tt.args.crtList, tt.args.entry); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.AddCrtListEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_DeleteCrtListEntry(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		crtList    string
		certFile   string
		lineNumber int
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "delete crt-list entries of crt-list, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crtList:    "/etc/haproxy/crt-list",
				certFile:   "/etc/ssl/cert-1.pem",
				lineNumber: 5,
			},
			wantErr: false,
			socketResponse: map[string]string{
				"del ssl crt-list /etc/haproxy/crt-list /etc/ssl/cert-1.pem:5\n": ` Entry '/etc/ssl/cert-1.pem' deleted in crtlist '/etc/ssl/crt-list'!
				`,
			},
		},
		{
			name:   "delete crt-list entries of crt-list, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crtList:    "/etc/haproxy/crt-list",
				certFile:   "/etc/ssl/not_known.pem",
				lineNumber: 10,
			},
			wantErr: true,
			socketResponse: map[string]string{
				"del ssl crt-list /etc/haproxy/crt-list /etc/ssl/not_known.pem:10\n": ` Can't edit the crt-list: certificate '/etc/ssl/not_known.pem' does not exist!
				`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.process, tt.fields.worker)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			if err := s.DeleteCrtListEntry(tt.args.crtList, tt.args.certFile, misc.Int64P(tt.args.lineNumber)); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.DeleteCrtListEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
