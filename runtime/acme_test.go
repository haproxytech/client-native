package runtime

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v6/models"
)

func TestAcmeStatus(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	dt := func(date string) strfmt.DateTime {
		tm, e := time.Parse(time.RFC3339, date)
		if e != nil {
			t.Error("dt:", e)
		}
		return strfmt.DateTime(tm)
	}

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	tests := []struct {
		name           string
		fields         fields
		want           models.AcmeStatus
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "Acme status example from official docs",
			fields: fields{socketPath: haProxy.Addr().String()},
			want: models.AcmeStatus{
				&models.AcmeCertificateStatus{
					Certificate:      "ecdsa.pem",
					AcmeSection:      "LE",
					State:            "Running",
					ExpiryDate:       dt("2020-01-18T09:31:12Z"),
					ExpiriesIn:       "0d 0h00m00s",
					ScheduledRenewal: dt("2020-01-15T21:31:12Z"),
					RenewalIn:        "0d 0h00m00s",
				},
				&models.AcmeCertificateStatus{
					Certificate:      "foobar.pem.rsa",
					AcmeSection:      "LE",
					State:            "Scheduled",
					ExpiryDate:       dt("2025-08-04T11:50:54Z"),
					ExpiriesIn:       "89d 23h01m13s",
					ScheduledRenewal: dt("2025-07-27T23:50:55Z"),
					RenewalIn:        "82d 11h01m14s",
				},
			},
			socketResponse: map[string]string{
				"acme status\n": `
# certificate	section	state	expiration date (UTC)	expires in	scheduled date (UTC)	scheduled in
ecdsa.pem	LE	Running	2020-01-18T09:31:12Z	0d 0h00m00s	2020-01-15T21:31:12Z	0d 0h00m00s
foobar.pem.rsa	LE	Scheduled	2025-08-04T11:50:54Z	89d 23h01m13s	2025-07-27T23:50:55Z	82d 11h01m14s
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
			got, err := s.AcmeStatus()
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.AcmeStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := range got {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("SingleRuntime.AcmeStatus() = %v, want %v", got[i], tt.want[i])
				}
			}
		})
	}
}
