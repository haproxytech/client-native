package runtime

import (
	"reflect"
	"testing"
	"time"

	"github.com/haproxytech/client-native/v2/models"
)

func TestSingleRuntime_ShowCerts(t *testing.T) {
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
		want           models.SslCertificates
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Simple show certs, should return 3 files",
			fields: fields{socketPath: haProxy.Addr().String()},
			want: models.SslCertificates{
				&models.SslCertificate{
					StorageName: "/etc/ssl/cert-0.pem",
					Description: "cert-0.pem",
				},
				&models.SslCertificate{
					StorageName: "/etc/ssl/cert-1.pem",
					Description: "cert-1.pem",
				},
				&models.SslCertificate{
					StorageName: "/etc/ssl/cert-2.pem",
					Description: "cert-2.pem",
				},
			},
			socketResponse: map[string]string{
				"show ssl cert\n": ` # filename
					/etc/ssl/cert-0.pem
					/etc/ssl/cert-1.pem
					/etc/ssl/cert-2.pem
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
			got, err := s.ShowCerts()
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowCerts() error = %v, wantErr %v", err, tt.wantErr)
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

func TestSingleRuntime_GetCert(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		storageName string
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
			name:   "Get certs, should return a cert",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "/etc/ssl/cert-0.pem",
			},
			want: &models.SslCertificate{
				StorageName: "/etc/ssl/cert-0.pem",
				Description: "cert-0.pem",
			},
			socketResponse: map[string]string{
				"show ssl cert\n": ` # filename
					/etc/ssl/cert-0.pem
					/etc/ssl/cert-1.pem
					/etc/ssl/cert-2.pem
				`,
			},
		},
		{
			name:   "Get unknown certs, should return a cert",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "/etc/ssl/cert-36.pem",
			},
			want:    nil,
			wantErr: true,
			socketResponse: map[string]string{
				"show ssl cert\n": ` # filename
					/etc/ssl/cert-0.pem
					/etc/ssl/cert-1.pem
					/etc/ssl/cert-2.pem
				`,
			},
		},
		{
			name:   "Get certs with empty argument, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "",
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
			got, err := s.GetCert(tt.args.storageName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.GetCert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SingleRuntime.GetCert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleRuntime_ShowCertEntry(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	notBefore, _ := time.Parse("Jan 2 15:04:05 2006 MST", "Sep  9 00:00:00 2020 GMT")
	notAfter, _ := time.Parse("Jan 2 15:04:05 2006 MST", "Sep 14 12:00:00 2021 GMT")
	type fields struct {
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		storageName string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *SslCertEntry
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Simple show certs, should return a cert",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "/etc/ssl/cert-0.pem",
			},
			want: &SslCertEntry{
				StorageName: "/etc/ssl/cert-0.pem",
				Status:      "Used",
				Serial:      "0D933C1B1089BF660AE5253A245BB388",
				NotBefore:   notBefore,
				NotAfter:    notAfter,
				SubjectAlternativeNames: []string{
					"DNS:*.platform.domain.com",
					"DNS:uaa.platform.domain.com",
				},
				Algorithm:       "RSA4096",
				SHA1FingerPrint: "59242F1838BDEF3E7DAFC83FFE4DD6C03B88805C",
				Subject:         "/C=DE/ST=Baden-Württemberg/L=Walldorf/O=ORG SE/CN=*.platform.domain.com",
				Issuer:          "/C=US/O=DigiCert Inc/CN=DigiCert SHA2 Secure Server CA",
				ChainSubject:    "/C=US/O=DigiCert Inc/CN=DigiCert SHA2 Secure Server CA",
				ChainIssuer:     "/C=US/O=DigiCert Inc/OU=www.digicert.com/CN=DigiCert Global Root CA",
			},
			socketResponse: map[string]string{
				"show ssl cert /etc/ssl/cert-0.pem\n": ` Filename: /etc/ssl/cert-0.pem
					Status: Used
					Serial: 0D933C1B1089BF660AE5253A245BB388
					notBefore: Sep  9 00:00:00 2020 GMT
					notAfter: Sep 14 12:00:00 2021 GMT
					Subject Alternative Name: DNS:*.platform.domain.com, DNS:uaa.platform.domain.com
					Algorithm: RSA4096
					SHA1 FingerPrint: 59242F1838BDEF3E7DAFC83FFE4DD6C03B88805C
					Subject: /C=DE/ST=Baden-Württemberg/L=Walldorf/O=ORG SE/CN=*.platform.domain.com
					Issuer: /C=US/O=DigiCert Inc/CN=DigiCert SHA2 Secure Server CA
					Chain Subject: /C=US/O=DigiCert Inc/CN=DigiCert SHA2 Secure Server CA
					Chain Issuer: /C=US/O=DigiCert Inc/OU=www.digicert.com/CN=DigiCert Global Root CA
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
			got, err := s.ShowCertEntry(tt.args.storageName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowCertEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SingleRuntime.ShowCertEntries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleRuntime_NewCertEntry(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		jobs       chan Task
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		storageName string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Create new cert, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "/etc/ssl/new_cert.pem",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"new ssl cert /etc/ssl/new_cert.pem\n": ` New empty certificate store '/etc/ssl/new_cert.pem'!
				`,
			},
		},
		{
			name:   "Create new cert without storageName, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "",
			},
			wantErr:        true,
			socketResponse: nil,
		},
		{
			name:   "Create new cert which already exists, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "/etc/ssl/existing_cert.pem",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"new ssl cert /etc/ssl/existing_cert.pem\n": ` Certificate '/etc/ssl/existing_cert.pem' already exists!
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
			if err := s.NewCertEntry(tt.args.storageName); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.NewCertEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_SetCertEntry(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		jobs       chan Task
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		storageName string
		payload     string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Create new cert, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "/etc/ssl/new_cert.pem",
				payload:     "-----BEGIN CERTIFICATE-----<redacted>...",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"set ssl cert /etc/ssl/new_cert.pem <<\n-----BEGIN CERTIFICATE-----<redacted>...\n": ` Transaction created for certificate /etc/ssl/new_cert.pem!
				`,
			},
		},
		{
			name:   "Create new cert, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "/etc/ssl/wrong_cert.pem",
				payload:     "-----BEGIN CERTIFICATE-----<redacted_wrong_cert>...",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"set ssl cert /etc/ssl/wrong_cert.pem <<\n-----BEGIN CERTIFICATE-----<redacted_wrong_cert>...\n": ` unable to load certificate from file 'wrong_cert.pem'.
					Can't load the payload
					Can't update wrong_cert.pem!
					`,
			},
		},
		{
			name:   "Create new cert, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "/etc/ssl/wrong_cert.pem",
				payload:     "",
			},
			wantErr:        true,
			socketResponse: map[string]string{},
		},
		{
			name:   "Create new cert, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "",
				payload:     "-----BEGIN CERTIFICATE-----<redacted>...",
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
			if err := s.SetCertEntry(tt.args.storageName, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.SetCertEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_CommitCertEntry(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		jobs       chan Task
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		storageName string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Commit updated cert, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "/etc/ssl/updated_cert.pem",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"commit ssl cert /etc/ssl/updated_cert.pem\n": ` Committing /etc/ssl/updated_cert.pem"
				Success!
				`,
			},
		},
		{
			name:   "Commit incomplete cert, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "/etc/ssl/incomplete_cert.pem",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"commit ssl cert /etc/ssl/incomplete_cert.pem\n": ` The transaction must contain at least a certificate and a private key!
				Can't commit /etc/ssl/incomplete_cert.pem!
				`,
			},
		},
		{
			name:   "Commit without storageName, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "",
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
			if err := s.CommitCertEntry(tt.args.storageName); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.CommitCertEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_AbortCertEntry(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		jobs       chan Task
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		storageName string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Abort updated cert, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "/etc/ssl/updated_cert.pem",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"abort ssl cert /etc/ssl/updated_cert.pem\n": ` Transaction aborted for certificate '/etc/ssl/updated_cert.pem'!
				`,
			},
		},
		{
			name:   "Abort cert no transaction exists, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "/etc/ssl/no_transaction_cert.pem",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"abort ssl cert /etc/ssl/no_transaction_cert.pem\n": ` No ongoing transaction!
				`,
			},
		},
		{
			name:   "Abort without storageName, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "",
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
			if err := s.AbortCertEntry(tt.args.storageName); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.AbortCertEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_DeleteCertEntry(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		jobs       chan Task
		socketPath string
		worker     int
		process    int
	}
	type args struct {
		storageName string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Delete cert, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "/etc/ssl/delete_cert.pem",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"del ssl cert /etc/ssl/delete_cert.pem\n": ` Certificate '/etc/ssl/delete_cert.pem' deleted!
				`,
			},
		},
		{
			name:   "Delete cert without storageName, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "",
			},
			wantErr:        true,
			socketResponse: nil,
		},
		{
			name:   "Delete cert which not exists, should return an error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				storageName: "/etc/ssl/not_existing_cert.pem",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"del ssl cert /etc/ssl/not_existing_cert.pem\n": ` Can't remove the certificate: certificate '/etc/ssl/not_existing_cert.pem' doesn't exist!
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
			if err := s.DeleteCertEntry(tt.args.storageName); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.DeleteCertEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
