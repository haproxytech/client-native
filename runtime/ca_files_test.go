package runtime

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v5/models"
)

func TestSingleRuntime_ShowCAFiles(t *testing.T) {
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
		want           models.SslCaFiles
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Simple show ca-files, should return 3 files",
			fields: fields{socketPath: haProxy.Addr().String()},
			want: models.SslCaFiles{
				&models.SslCaFile{
					StorageName: "cafile.pem",
					Count:       "3 certificate(s)",
				},
			},
			socketResponse: map[string]string{
				"show ssl ca-files\n": ` # filename
					cafile.pem - 3 certificate(s)
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
			got, err := s.ShowCAFiles()
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowCAFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("SingleRuntime.ShowCAFiles() = %v, want %v", got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSingleRuntime_GetCAFile(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		caFile string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *models.SslCaFile
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "GetCAFile, should return a ca file",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "cafile.pem",
			},
			want: &models.SslCaFile{
				StorageName: "cafile.pem",
				Count:       "3 certificate(s)",
			},
			socketResponse: map[string]string{
				"show ssl ca-file\n": ` # filename
					cafile.pem - 3 certificate(s)
				`,
			},
		},
		{
			name:   "Get unknown ca files, should not return a ca file",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "/etc/ssl/cafile-36.pem",
			},
			want:    nil,
			wantErr: true,
			socketResponse: map[string]string{
				"show ssl ca-file\n": ` # filename
					cafile.pem - 3 certificate(s)
				`,
			},
		},
		{
			name:   "Get ca files with empty argument, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "",
			},
			want:           nil,
			wantErr:        true,
			socketResponse: map[string]string{},
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
			got, err := s.GetCAFile(tt.args.caFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.GetCAFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SingleRuntime.GetCAFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleRuntime_ShowCAFile(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	notBefore, _ := time.Parse("Jan 2 15:04:05 2006 MST", "Sep  9 00:00:00 2020 GMT")
	notAfter, _ := time.Parse("Jan 2 15:04:05 2006 MST", "Sep 14 12:00:00 2021 GMT")
	index := int64(2)
	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		caFile string
		index  *int64
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *models.SslCertificate
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Simple show ca-files, should return a ca file",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "/etc/ssl/cafile-0.pem",
			},
			want: &models.SslCertificate{
				StorageName:             "/etc/ssl/cafile-0.pem",
				Serial:                  "0D933C1B1089BF660AE5253A245BB388",
				NotBefore:               (*strfmt.DateTime)(&notBefore),
				NotAfter:                (*strfmt.DateTime)(&notAfter),
				SubjectAlternativeNames: "DNS:*.platform.domain.com, DNS:uaa.platform.domain.com",
				Algorithm:               "RSA4096",
				Sha1FingerPrint:         "59242F1838BDEF3E7DAFC83FFE4DD6C03B88805C",
				Subject:                 "/C=FR/ST=Some-State/O=HAProxy Technologies/CN=HAProxy Technologies CA",
				Issuers:                 "/C=FR/ST=Some-State/O=HAProxy Technologies/CN=HAProxy Technologies CA",
			},
			socketResponse: map[string]string{
				"show ssl ca-file /etc/ssl/cafile-0.pem\n": ` Filename: /etc/ssl/cafile-0.pem
					Serial: 0D933C1B1089BF660AE5253A245BB388
					notBefore: Sep  9 00:00:00 2020 GMT
					notAfter: Sep 14 12:00:00 2021 GMT
					Subject Alternative Name: DNS:*.platform.domain.com, DNS:uaa.platform.domain.com
					Algorithm: RSA4096
					SHA1 FingerPrint: 59242F1838BDEF3E7DAFC83FFE4DD6C03B88805C
					Subject: /C=FR/ST=Some-State/O=HAProxy Technologies/CN=HAProxy Technologies CA
					Issuer: /C=FR/ST=Some-State/O=HAProxy Technologies/CN=HAProxy Technologies CA
				`,
			},
		},
		{
			name:   "Simple show ca-files with index, should return a ca file",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "/etc/ssl/cafile-0.pem",
				index:  &index,
			},
			want: &models.SslCertificate{
				StorageName:             "/etc/ssl/cafile-0.pem",
				Serial:                  "0D933C1B1089BF660AE5253A245BB388",
				NotBefore:               (*strfmt.DateTime)(&notBefore),
				NotAfter:                (*strfmt.DateTime)(&notAfter),
				SubjectAlternativeNames: "DNS:*.platform.domain.com, DNS:uaa.platform.domain.com",
				Algorithm:               "RSA4096",
				Sha1FingerPrint:         "59242F1838BDEF3E7DAFC83FFE4DD6C03B88805C",
				Subject:                 "/C=FR/ST=Some-State/O=HAProxy Technologies/CN=HAProxy Technologies CA",
				Issuers:                 "/C=FR/ST=Some-State/O=HAProxy Technologies/CN=HAProxy Technologies CA",
			},
			socketResponse: map[string]string{
				"show ssl ca-file /etc/ssl/cafile-0.pem:2\n": ` Filename: /etc/ssl/cafile-0.pem
					Serial: 0D933C1B1089BF660AE5253A245BB388
					notBefore: Sep  9 00:00:00 2020 GMT
					notAfter: Sep 14 12:00:00 2021 GMT
					Subject Alternative Name: DNS:*.platform.domain.com, DNS:uaa.platform.domain.com
					Algorithm: RSA4096
					SHA1 FingerPrint: 59242F1838BDEF3E7DAFC83FFE4DD6C03B88805C
					Subject: /C=FR/ST=Some-State/O=HAProxy Technologies/CN=HAProxy Technologies CA
					Issuer: /C=FR/ST=Some-State/O=HAProxy Technologies/CN=HAProxy Technologies CA
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
			got, err := s.ShowCAFile(tt.args.caFile, tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowCAFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SingleRuntime.ShowCAFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleRuntime_NewCAFile(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		caFile string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Create new ca-file, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "/etc/ssl/new_cafile.pem",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"new ssl ca-file /etc/ssl/new_cafile.pem\n": ` New CA file created '/etc/ssl/new_cafile.pem'!
				`,
			},
		},
		{
			name:   "Create new ca-file without ca file, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "",
			},
			wantErr:        true,
			socketResponse: nil,
		},
		{
			name:   "Create new ca-file which already exists, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "/etc/ssl/existing_cafile.pem",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"new ssl ca-file /etc/ssl/existing_cafile.pem\n": ` CAFile '/etc/ssl/existing_cafile.pem' already exists!
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
			if err := s.NewCAFile(tt.args.caFile); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.NewCAFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_SetCAFile(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		caFile  string
		payload string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Create new ca-file, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile:  "/etc/ssl/new_cafile.pem",
				payload: "-----BEGIN CERTIFICATE-----<redacted>...",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"set ssl ca-file /etc/ssl/new_cafile.pem <<\n-----BEGIN CERTIFICATE-----<redacted>...\n": ` Transaction created for ca-file /etc/ssl/new_cafile.pem!
				`,
			},
		},
		{
			name:   "Create new ca-file with a wrong payload, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile:  "/etc/ssl/wrong_cafile.pem",
				payload: "-----BEGIN CERTIFICATE-----<redacted_wrong_cert>...",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"set ssl ca-file /etc/ssl/wrong_cafile.pem <<\n-----BEGIN CERTIFICATE-----<redacted_wrong_cafile>...\n": ` unable to load certificate from file 'wrong_cafile.pem'.
					Can't load the payload
					Can't update wrong_cafile.pem!
					`,
			},
		},
		{
			name:   "Create new cafile without payload, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile:  "/etc/ssl/wrong_cafile.pem",
				payload: "",
			},
			wantErr:        true,
			socketResponse: map[string]string{},
		},
		{
			name:   "Create new cafile without ca file, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile:  "",
				payload: "-----BEGIN CERTIFICATE-----<redacted>...",
			},
			wantErr:        true,
			socketResponse: map[string]string{},
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
			if err := s.SetCAFile(tt.args.caFile, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.SetCAFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_CommitCAFile(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		caFile string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Commit updated ca-file, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "/etc/ssl/updated_cafile.pem",
			},
			socketResponse: map[string]string{
				"commit ssl ca-file /etc/ssl/updated_cafile.pem\n": ` Committing /etc/ssl/updated_cafile.pem
				Success!
				`,
			},
		},
		{
			name:   "Commit incomplete ca-file, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "/etc/ssl/incomplete_cafile.pem",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"commit ssl ca-file /etc/ssl/incomplete_cafile.pem\n": ` The transaction must contain at least a ca-file!
				Can't commit /etc/ssl/incomplete_cafile.pem!
				`,
			},
		},
		{
			name:   "Commit without ca-file, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "",
			},
			wantErr:        true,
			socketResponse: map[string]string{},
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
			if err := s.CommitCAFile(tt.args.caFile); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.CommitCAFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_AbortCAFile(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		caFile string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Abort updated ca-file, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "/etc/ssl/updated_cafile.pem",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"abort ssl ca-file /etc/ssl/updated_cafile.pem\n": ` Transaction aborted for certificate '/etc/ssl/updated_cafile.pem'!
				`,
			},
		},
		{
			name:   "Abort ca-file no transaction exists, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "/etc/ssl/no_transaction_cafile.pem",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"abort ssl ca-file /etc/ssl/no_transaction_cafile.pem\n": ` No ongoing transaction!
				`,
			},
		},
		{
			name:   "Abort without ca-file, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "",
			},
			wantErr:        true,
			socketResponse: map[string]string{},
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
			if err := s.AbortCAFile(tt.args.caFile); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.AbortCAFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_AddCAFileEntry(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		caFile  string
		payload string
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
				caFile:  "/etc/haproxy/ca-file",
				payload: "-----BEGIN CERTIFICATE-----<redacted>...",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"add ssl ca-file /etc/haproxy/ca-file <<\n-----BEGIN CERTIFICATE-----<redacted>...\n": ` transaction created for CA '/etc/haproxy/ca-file'.
				Success!
				`,
			},
		},
		{
			name:   "Create new ca-file with a wrong payload, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile:  "/etc/ssl/wrong_cafile.pem",
				payload: "-----BEGIN CERTIFICATE-----<redacted_wrong_cert>...",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"add ssl ca-file /etc/ssl/wrong_cafile.pem <<\n-----BEGIN CERTIFICATE-----<redacted_wrong_cafile>...\n": ` unable to insert certificate into file 'wrong_cafile.pem'.
					Can't load the payload
					Can't update wrong_cafile.pem!
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
			if err := s.AddCAFileEntry(tt.args.caFile, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.AddCAFileEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_DeleteCAFile(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		caFile string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Delete ca-file, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "/etc/ssl/delete_cafile.pem",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"del ssl ca-file /etc/ssl/delete_cafile.pem\n": ` CA file '/etc/ssl/delete_cafile.pem' deleted!
				`,
			},
		},
		{
			name:   "Delete cert without ca-file, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "",
			},
			wantErr:        true,
			socketResponse: nil,
		},
		{
			name:   "Delete ca file which does not exist, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				caFile: "/etc/ssl/not_existing_cafile.pem",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"del ssl ca-file /etc/ssl/not_existing_cafile.pem\n": ` Can't remove the ca-file: ca-file '/etc/ssl/not_existing_cafile.pem' doesn't exist!
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
			if err := s.DeleteCAFile(tt.args.caFile); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.DeleteCAFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
