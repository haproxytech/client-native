// Copyright 2025 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package runtime

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v6/models"
)

func (s *SingleRuntime) AcmeRenew(certificate string) error {
	return s.Execute("acme renew " + certificate)
}

func (s *SingleRuntime) AcmeStatus() (models.AcmeStatus, error) {
	resp, err := s.ExecuteWithResponse("acme status")
	if err != nil {
		return nil, err
	}
	return parseAcmeStatus(resp)
}

func parseAcmeStatus(resp string) (models.AcmeStatus, error) {
	// Example:
	// # certificate   section  state      expiration date (UTC)  expires in        scheduled date (UTC)  scheduled in
	// ecdsa.pem       LE       Running    2020-01-18T09:31:12Z   0d 0h00m00s       2020-01-15T21:31:12Z  0d 0h00m00s
	// foobar.pem.rsa  LE       Scheduled  2025-08-04T11:50:54Z   89d 23h01m13s     2025-07-27T23:50:55Z  82d 11h01m14s

	resp = strings.TrimSpace(resp)
	if resp == "" {
		return models.AcmeStatus{}, nil
	}

	status := make(models.AcmeStatus, 0, 8)
	var err error

	strings.SplitSeq(resp, "\n")(func(line string) bool {
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			// skip the header line
			return true
		}

		const expectedFieldsNo = 7

		parts := strings.SplitN(line, "\t", expectedFieldsNo)
		if len(parts) < expectedFieldsNo {
			err = fmt.Errorf("failed to parse response: not enough fields in line '%s'", line)
			return false
		}

		exp, e := parseAcmeDate(parts[3])
		if e != nil {
			err = fmt.Errorf("failed to parse date in response: '%s': %w", parts[3], e)
			return false
		}
		sched, e := parseAcmeDate(parts[5])
		if e != nil {
			err = fmt.Errorf("failed to parse date in response: '%s': %w", parts[5], e)
			return false
		}

		s := &models.AcmeCertificateStatus{
			Certificate:      parts[0],
			AcmeSection:      parts[1],
			State:            parts[2],
			ExpiryDate:       strfmt.DateTime(exp),
			ExpiriesIn:       parts[4],
			ScheduledRenewal: strfmt.DateTime(sched),
			RenewalIn:        parts[6],
		}

		status = append(status, s)
		return true
	})

	return status, err
}

func parseAcmeDate(s string) (time.Time, error) {
	if s == "-" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, s)
}
