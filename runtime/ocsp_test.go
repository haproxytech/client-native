package runtime

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/haproxytech/client-native/v6/models"
)

func TestSingleRuntime_SetOcspResponse(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	type args struct {
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
			name:   "Set the ssl ocsp response, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				payload: "$(base64 resp.der)",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"set ssl ocsp-response <<\n$(base64 resp.der)\n": ` OCSP Response updated!\n
				`,
			},
		},
		{
			name:   "Set the ssl ocsp response with no response, should return error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				payload: "$(base64 no_resp.der)",
			},
			wantErr: true,
			socketResponse: map[string]string{
				"set ssl ocsp-response <<\n$(base64 resp.der)\n": ` Some error
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
			if err := s.SetOcspResponse(tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.SetOcspResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSingleRuntime_ShowOcspResponses(t *testing.T) {
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
		responseValues []*models.SslCertificateID
	}{
		{
			name:    "show the ssl ocsp responses, should return no error",
			fields:  fields{socketPath: haProxy.Addr().String()},
			wantErr: false,
			socketResponse: map[string]string{
				"show ssl ocsp-response\n": ` Certificate ID key : 303b300906052b0e03021a050004148a83e0060faff709ca7e9b95522a2e81635fda0a0414f652b0e435d5ea923851508f0adbe92d85de007a0202100a
  Certificate path : /path_to_cert/foo.pem
    Certificate ID:
      Issuer Name Hash: 8A83E0060FAFF709CA7E9B95522A2E81635FDA0A
      Issuer Key Hash: F652B0E435D5EA923851508F0ADBE92D85DE007A
      Serial Number: 100A
				`,
			},
			responseValues: []*models.SslCertificateID{
				{
					CertificateIDKey: "303b300906052b0e03021a050004148a83e0060faff709ca7e9b95522a2e81635fda0a0414f652b0e435d5ea923851508f0adbe92d85de007a0202100a",
					CertificatePath:  "/path_to_cert/foo.pem",
					CertificateID: &models.CertificateID{
						IssuerNameHash: "8A83E0060FAFF709CA7E9B95522A2E81635FDA0A",
						IssuerKeyHash:  "F652B0E435D5EA923851508F0ADBE92D85DE007A",
						SerialNumber:   "100A",
					},
				},
			},
		},
		{
			name:    "show the ssl ocsp response with some error",
			fields:  fields{socketPath: haProxy.Addr().String()},
			wantErr: true,
			socketResponse: map[string]string{
				"show ssl ocsp-response\n": ``,
			},
			responseValues: nil,
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
			certificateIds, err := s.ShowOcspResponses()
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowOcspResponses() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(certificateIds, tt.responseValues) {
				t.Errorf("SingleRuntime.ShowOcspResponses() = %+v, want %+v", certificateIds[0], tt.responseValues[0])
			}
		})
	}
}

func TestSingleRuntime_ShowOcspResponse(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	producedAt, _ := time.Parse("Jan 2 15:04:05 2006 GMT", "May 27 15:43:38 2021 GMT")
	thisUpdate, _ := time.Parse("Jan 2 15:04:05 2006 GMT", "May 27 15:43:38 2021 GMT")
	nextUpdate, _ := time.Parse("Jan 2 15:04:05 2006 GMT", "Oct 12 15:43:38 2048 GMT")
	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	type args struct {
		idOrPath string
		ofmt     []OcspResponseFmt
	}
	tests := []struct {
		name                string
		fields              fields
		args                args
		wantErr             bool
		socketResponse      map[string]string
		responseValues      *models.SslOcspResponse
		base64ResponseValue string
	}{
		{
			name:    "Show the ssl ocsp response with text, should return no error",
			fields:  fields{socketPath: haProxy.Addr().String()},
			args:    args{idOrPath: "303b300906052b0e03021a050004148a83e0060faff709ca7e9b95522a2e81635fda0a0414f652b0e435d5ea923851508f0adbe92d85de007a0202100a"},
			wantErr: false,
			socketResponse: map[string]string{
				"show ssl ocsp-response 303b300906052b0e03021a050004148a83e0060faff709ca7e9b95522a2e81635fda0a0414f652b0e435d5ea923851508f0adbe92d85de007a0202100a\n": ` OCSP Response Data:
  OCSP Response Status: successful (0x0)
  Response Type: Basic OCSP Response
  Version: 1 (0x0)
  Responder Id: C = FR, O = HAProxy Technologies, CN = ocsp.haproxy.com
  Produced At: May 27 15:43:38 2021 GMT
  Responses:
  Certificate ID:
    Hash Algorithm: sha1
    Issuer Name Hash: 8A83E0060FAFF709CA7E9B95522A2E81635FDA0A
    Issuer Key Hash: F652B0E435D5EA923851508F0ADBE92D85DE007A
    Serial Number: 100A
  Cert Status: good
  This Update: May 27 15:43:38 2021 GMT
  Next Update: Oct 12 15:43:38 2048 GMT
				`,
			},
			responseValues: &models.SslOcspResponse{
				OcspResponseStatus: "successful",
				ResponseType:       "Basic OCSP Response",
				Version:            "1",
				ResponderID: []string{
					"C = FR",
					"O = HAProxy Technologies",
					"CN = ocsp.haproxy.com",
				},
				ProducedAt: strfmt.Date(producedAt),
				Responses: &models.OCSPResponses{
					CertificateID: &models.CertificateID{
						HashAlgorithm:  "sha1",
						IssuerNameHash: "8A83E0060FAFF709CA7E9B95522A2E81635FDA0A",
						IssuerKeyHash:  "F652B0E435D5EA923851508F0ADBE92D85DE007A",
						SerialNumber:   "100A",
					},
					CertStatus: "good",
					ThisUpdate: strfmt.Date(thisUpdate),
					NextUpdate: strfmt.Date(nextUpdate),
				},
			},
		},
		{
			name:    "Show the ssl ocsp response with base64, should return no error",
			fields:  fields{socketPath: haProxy.Addr().String()},
			args:    args{idOrPath: "/path_to_cert/foo.pem", ofmt: []OcspResponseFmt{OcspFmtBase64}},
			wantErr: false,
			socketResponse: map[string]string{
				"show ssl ocsp-response base64 /path_to_cert/foo.pem\n": ` MIIB8woBAKCCAewwggHoBgkrBgEFBQcwAQEEggHZMIIB1TCBvqE
				`,
			},
			responseValues: &models.SslOcspResponse{
				Base64Response: "MIIB8woBAKCCAewwggHoBgkrBgEFBQcwAQEEggHZMIIB1TCBvqE",
			},
		},
		{
			name:    "Show the ssl ocsp response with error",
			fields:  fields{socketPath: haProxy.Addr().String()},
			args:    args{idOrPath: "303b300906052b0e03021a050004148a83e0060faff709ca7e9b95522a2e81635fda0a0414f652b0e435d5ea923851508f0adbe92d85de007a0202100a"},
			wantErr: true,
			socketResponse: map[string]string{
				"show ssl ocsp-response 303b300906052b0e03021a050004148a83e0060faff709ca7e9b95522a2e81635fda0a0414f652b0e435d5ea923851508f0adbe92d85de007a0202100a\n": ``,
			},
			responseValues: nil,
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
			ocspResponse, err := s.ShowOcspResponse(tt.args.idOrPath, tt.args.ofmt...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowOcspResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if ocspResponse != nil && !reflect.DeepEqual(ocspResponse, tt.responseValues) {
				t.Errorf("SingleRuntime.ShowOcspResponse() = %v, want %v", ocspResponse, tt.responseValues)
			}
		})
	}
}

func TestSingleRuntime_UpdateOcspResponse(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	producedAt, _ := time.Parse("Jan 2 15:04:05 2006 GMT", "May 27 15:43:38 2021 GMT")
	thisUpdate, _ := time.Parse("Jan 2 15:04:05 2006 GMT", "May 27 15:43:38 2021 GMT")
	nextUpdate, _ := time.Parse("Jan 2 15:04:05 2006 GMT", "Oct 12 15:43:38 2048 GMT")

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	type args struct {
		certFile string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		socketResponse map[string]string
		responseValues *models.SslOcspResponse
	}{
		{
			name:   "update the ssl ocsp response, should return no error",
			fields: fields{socketPath: haProxy.Addr().String()},
			args: args{
				certFile: "cert.pem",
			},
			wantErr: false,
			socketResponse: map[string]string{
				"update ssl ocsp-response cert.pem\n": ` OCSP Response Data:
  OCSP Response Status: successful (0x0)
  Response Type: Basic OCSP Response
  Version: 1 (0x0)
  Responder Id: C = FR, O = HAProxy Technologies, CN = ocsp.haproxy.com
  Produced At: May 27 15:43:38 2021 GMT
  Responses:
  Certificate ID:
    Hash Algorithm: sha1
    Issuer Name Hash: 8A83E0060FAFF709CA7E9B95522A2E81635FDA0A
    Issuer Key Hash: F652B0E435D5EA923851508F0ADBE92D85DE007A
    Serial Number: 100A
  Cert Status: good
  This Update: May 27 15:43:38 2021 GMT
  Next Update: Oct 12 15:43:38 2048 GMT
				`,
			},
			responseValues: &models.SslOcspResponse{
				OcspResponseStatus: "successful",
				ResponseType:       "Basic OCSP Response",
				Version:            "1",
				ResponderID: []string{
					"C = FR",
					"O = HAProxy Technologies",
					"CN = ocsp.haproxy.com",
				},
				ProducedAt: strfmt.Date(producedAt),
				Responses: &models.OCSPResponses{
					CertificateID: &models.CertificateID{
						HashAlgorithm:  "sha1",
						IssuerNameHash: "8A83E0060FAFF709CA7E9B95522A2E81635FDA0A",
						IssuerKeyHash:  "F652B0E435D5EA923851508F0ADBE92D85DE007A",
						SerialNumber:   "100A",
					},
					CertStatus: "good",
					ThisUpdate: strfmt.Date(thisUpdate),
					NextUpdate: strfmt.Date(nextUpdate),
				},
			},
		},
		{
			name:    "update the ssl ocsp response with some error",
			fields:  fields{socketPath: haProxy.Addr().String()},
			wantErr: true,
			socketResponse: map[string]string{
				"update ssl ocsp-response cert.pem\n": ` Some error
				`,
			},
			responseValues: nil,
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
			ocspResponse, err := s.UpdateOcspResponse(tt.args.certFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.UpdateOcspResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(ocspResponse, tt.responseValues) {
				t.Errorf("SingleRuntime.UpdateOcspResponse() = %v, want %v", ocspResponse, tt.responseValues)
			}
		})
	}
}

func TestSingleRuntime_ShowOcspUpdates(t *testing.T) {
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
		responseValues []*models.SslOcspUpdate
	}{
		{
			name:    "show the ssl ocsp updates, should return no error",
			fields:  fields{socketPath: haProxy.Addr().String()},
			wantErr: false,
			socketResponse: map[string]string{
				"show ssl ocsp-updates\n": ` OCSP Certid | Path | Next Update | Last Update | Successes | Failures | Last Update Status | Last Update Status (str)
      303b300906052b0e03021a050004148a83e0060faff709ca7e9b95522a2e81635fda0a0414f652b0e435d5ea923851508f0adbe92d85de007a02021015 | /path_to_cert/cert.pem | 30/Jan/2023:00:08:09 +0000 | - | 0 | 1 | 2 | HTTP error`,
			},
			responseValues: []*models.SslOcspUpdate{
				{
					CertID:              "303b300906052b0e03021a050004148a83e0060faff709ca7e9b95522a2e81635fda0a0414f652b0e435d5ea923851508f0adbe92d85de007a02021015",
					Path:                "/path_to_cert/cert.pem",
					NextUpdate:          "30/Jan/2023:00:08:09 +0000",
					Successes:           0,
					Failures:            1,
					LastUpdateStatus:    2,
					LastUpdateStatusStr: "HTTP error",
				},
			},
		},
		{
			name:    "show the ssl ocsp updates with some error",
			fields:  fields{socketPath: haProxy.Addr().String()},
			wantErr: true,
			socketResponse: map[string]string{
				"show ssl ocsp-updates\n": ``,
			},
			responseValues: nil,
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
			ocspUpdates, err := s.ShowOcspUpdates()
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowOcspUpdates() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(ocspUpdates, tt.responseValues) {
				t.Errorf("SingleRuntime.ShowOcspUpdates() = %+v, want %+v", ocspUpdates[0], tt.responseValues[0])
			}
		})
	}
}
