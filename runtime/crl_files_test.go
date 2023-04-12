package runtime

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/models"
)

func TestSingleRuntime_ShowCrlFiles(t *testing.T) {
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
		want           models.SslCrls
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Simple show crl-files, should return 3 files",
			fields: fields{socketPath: haProxy.Addr().String()},
			want: models.SslCrls{
				&models.SslCrl{
					StorageName: "/etc/ssl/crl-0.pem",
					Description: "crl-0.pem",
				},
				&models.SslCrl{
					StorageName: "/etc/ssl/crl-1.pem",
					Description: "crl-1.pem",
				},
				&models.SslCrl{
					StorageName: "/etc/ssl/crl-2.pem",
					Description: "crl-2.pem",
				},
			},
			socketResponse: map[string]string{
				"show ssl crl-file\n": ` # filename
					/etc/ssl/crl-0.pem
					/etc/ssl/crl-1.pem
					/etc/ssl/crl-2.pem
				`,
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
			got, err := s.ShowCrlFiles()
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowCrlFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("SingleRuntime.ShowCrlFiles() = %v, want %v", got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSingleRuntime_GetCrlFile(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	type args struct {
		crlFile string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *models.SslCrl
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "GetCrlFile, should return a crl file",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/crl-0.pem",
			},
			want: &models.SslCrl{
				StorageName: "/etc/ssl/crl-0.pem",
				Description: "crl-0.pem",
			},
			socketResponse: map[string]string{
				"show ssl crl-file\n": ` # filename
					/etc/ssl/crl-0.pem
					/etc/ssl/crl-1.pem
					/etc/ssl/crl-2.pem
				`,
			},
		},
		{
			name:   "Get unknown crl files, should not return a crl file",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/crl-36.pem",
			},
			want:    nil,
			wantErr: true,
			socketResponse: map[string]string{
				"show ssl crl-file\n": ` # filename
					/etc/ssl/crl-0.pem
					/etc/ssl/crl-1.pem
					/etc/ssl/crl-2.pem
				`,
			},
		},
		{
			name:   "Get crl files with empty argument, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "",
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
			err := s.Init(tt.fields.socketPath, tt.fields.masterWorkerMode)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			got, err := s.GetCrlFile(tt.args.crlFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.GetCrlFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SingleRuntime.GetCrlFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleRuntime_ShowCrlFile(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	thisUpdate, _ := time.Parse("Jan 2 15:04:05 2006 GMT", "Apr 1 14:45:39 2023 GMT")
	nextUpdate, _ := time.Parse("Jan 2 15:04:05 2006 GMT", "Sep 8 14:45:39 2048 GMT")
	revocationDate, _ := time.Parse("Jan 2 15:04:05 2006 GMT", "Apr 1 14:45:36 2023 GMT")
	index := int64(2)
	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	type args struct {
		crlFile string
		index   *int64
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *models.SslCrlEntries
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Simple show crl files, should return a crl file",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/crlfile-0.pem",
			},
			want: &models.SslCrlEntries{
				&models.SslCrlEntry{
					StorageName: "/etc/ssl/crlfile-0.pem",
					Status:      "Used",
					Version:     "2",
					LastUpdate:  strfmt.Date(thisUpdate),
					NextUpdate:  strfmt.Date(nextUpdate),
					RevokedCertificates: []*models.RevokedCertificates{
						{
							SerialNumber:   "1008",
							RevocationDate: strfmt.Date(revocationDate),
						},
					},
					SignatureAlgorithm: "sha256WithRSAEncryption",
					Issuer:             "/C=CA/ST=Some-State/O=HAProxy Technologies/CN=HAProxy Technologies CA",
				},
			},
			socketResponse: map[string]string{
				"show ssl crl-file /etc/ssl/crlfile-0.pem\n": ` Filename: /etc/ssl/crlfile-0.pem
Status: Used
Certificate Revocation List (CRL):
		Version 2
		Signature Algorithm: sha256WithRSAEncryption
		Issuer: /C=CA/ST=Some-State/O=HAProxy Technologies/CN=HAProxy Technologies CA
		This Update: Apr 1 14:45:39 2023 GMT
		Next Update: Sep 8 14:45:39 2048 GMT
		Revoked Certificates:
		Serial Number: 1008
		Revocation Date: Apr 1 14:45:36 2023 GMT
`,
			},
		},
		{
			name:   "Simple show crl files with index, should return a crl file",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/crlfile.pem",
				index:   &index,
			},
			want: &models.SslCrlEntries{
				&models.SslCrlEntry{
					StorageName: "/etc/ssl/crlfile.pem",
					Status:      "Used",
					Version:     "2",
					LastUpdate:  strfmt.Date(thisUpdate),
					NextUpdate:  strfmt.Date(nextUpdate),
					RevokedCertificates: []*models.RevokedCertificates{
						{
							SerialNumber:   "1008",
							RevocationDate: strfmt.Date(revocationDate),
						},
					},
					SignatureAlgorithm: "sha256WithRSAEncryption",
					Issuer:             "/C=CA/ST=Some-State/O=HAProxy Technologies/CN=HAProxy Technologies CA",
				},
			},
			socketResponse: map[string]string{
				"show ssl crl-file /etc/ssl/crlfile.pem:2\n": ` Filename: /etc/ssl/crlfile.pem
Status: Used
		Certificate Revocation List (CRL):
		Version 2
		Signature Algorithm: sha256WithRSAEncryption
		Issuer: /C=CA/ST=Some-State/O=HAProxy Technologies/CN=HAProxy Technologies CA
		This Update: Apr 1 14:45:39 2023 GMT
		Next Update: Sep 8 14:45:39 2048 GMT
		Revoked Certificates:
		Serial Number: 1008
		Revocation Date: Apr 1 14:45:36 2023 GMT
				`,
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
			got, err := s.ShowCrlFile(tt.args.crlFile, tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowCrlFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Error("SingleRuntime.ShowCrlFile(): ", cmp.Diff(got, tt.want))
				//t.Errorf("SingleRuntime.ShowCrlFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleRuntime_NewCrlFile(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	type args struct {
		crlFile string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Create new crl file, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/new_crlfile.pem",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"new ssl crl-file /etc/ssl/new_crlfile.pem\n": ` New CRL file created '/etc/ssl/new_crlfile.pem'!
				`,
			},
		},
		{
			name:   "Create new crl-file without crl file, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "",
			},
			wantErr:        true,
			socketResponse: nil,
		},
		{
			name:   "Create new crl file which already exists, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/existing_crlfile.pem",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"new ssl ca-file /etc/ssl/existing_crlfile.pem\n": ` CrlFile '/etc/ssl/existing_crlfile.pem' already exists!
				`,
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
			if err := s.NewCrlFile(tt.args.crlFile); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.NewCrlFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_SetCrlFile(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	type args struct {
		crlFile string
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
			name:   "Create new crl file, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/new_crlfile.pem",
				payload: "-----BEGIN X509 CRL-----<redacted>...-----END X509 CRL-----",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"set ssl crl-file /etc/ssl/new_crlfile.pem <<\n-----BEGIN X509 CRL-----<redacted>...-----END X509 CRL-----\n": ` transaction created for CRL /etc/ssl/new_crlfile.pem!
				`,
			},
		},
		{
			name:   "Create new crl file with a wrong payload, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/wrong_crlfile.pem",
				payload: "-----BEGIN X509 CRL-----<redacted_wrong_crl>...-----END X509 CRL-----",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"set ssl crl-file /etc/ssl/wrong_crlfile.pem <<\n-----BEGIN X509 CRL-----<redacted_wrong_crl>...-----END X509 CRL-----\n": ` unable to load crl from file 'wrong_crlfile.pem'.
					Can't load the payload
					Can't update wrong_crlfile.pem!
					`,
			},
		},
		{
			name:   "Create new crl file without payload, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/wrong_crlfile.pem",
				payload: "",
			},
			wantErr:        true,
			socketResponse: map[string]string{},
		},
		{
			name:   "Create new crl file without crl file, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "",
				payload: "-----BEGIN X509 CRL-----<redacted>...-----END X509 CRL-----",
			},
			wantErr:        true,
			socketResponse: map[string]string{},
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
			if err := s.SetCrlFile(tt.args.crlFile, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.SetCrlFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_CommitCrlFile(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	type args struct {
		crlFile string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Commit updated crl file, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/updated_crlfile.pem",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"commit ssl crl-file /etc/ssl/updated_crlfile.pem\n": ` Committing /etc/ssl/updated_crlfile.pem"
				Success!
				`,
			},
		},
		{
			name:   "Commit incomplete crl file, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/incomplete_crlfile.pem",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"commit ssl crl-file /etc/ssl/incomplete_crlfile.pem\n": ` The transaction must contain at least a crl file!
				Can't commit /etc/ssl/incomplete_crlfile.pem!
				`,
			},
		},
		{
			name:   "Commit without crl file, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "",
			},
			wantErr:        true,
			socketResponse: map[string]string{},
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
			if err := s.CommitCrlFile(tt.args.crlFile); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.CommitCrlFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_AbortCrlFile(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	type args struct {
		crlFile string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Abort updated crl file, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/updated_crlfile.pem",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"abort ssl crl-file /etc/ssl/updated_crlfile.pem\n": ` Transaction aborted for crl '/etc/ssl/updated_crlfile.pem'!
				`,
			},
		},
		{
			name:   "Abort crl file no transaction exists, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/no_transaction_crlfile.pem",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"abort ssl crl-file /etc/ssl/no_transaction_crlfile.pem\n": ` No ongoing transaction!
				`,
			},
		},
		{
			name:   "Abort without crl file, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "",
			},
			wantErr:        true,
			socketResponse: map[string]string{},
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
			if err := s.AbortCrlFile(tt.args.crlFile); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.AbortCrlFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_DeleteCrlFile(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	type args struct {
		crlFile string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Delete crl file, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/delete_crl.pem",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"del ssl crl-file /etc/ssl/delete_crl.pem\n": ` CRL file '/etc/ssl/delete_crl.pem' deleted!
				`,
			},
		},
		{
			name:   "Delete crl without crl file, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "",
			},
			wantErr:        true,
			socketResponse: nil,
		},
		{
			name:   "Delete crl file which does not exist, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				crlFile: "/etc/ssl/not_existing_crlfile.pem",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"del ssl crl-file /etc/ssl/not_existing_crlfile.pem\n": ` Can't remove the crl-file: crl-file '/etc/ssl/not_existing_crlfile.pem' doesn't exist!
				`,
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
			if err := s.DeleteCrlFile(tt.args.crlFile); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.DeleteCrlFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
