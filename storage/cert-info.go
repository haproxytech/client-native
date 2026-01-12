// Copyright 2023 HAProxy Technologies
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

package storage

import (
	"crypto/x509"
	"encoding/pem"
	"strings"
	"time"
)

// Information about stored certificates to be returned by the API.
type CertificatesInfo struct {
	NotAfter, NotBefore *time.Time
	DNS, IPs, Issuers   string
}

// Private struct to store unique info about multiple certificates.
type certsInfo struct {
	NotAfter, NotBefore time.Time
	DNS, IPs, Issuers   map[string]struct{}
}

func newCertsInfo() *certsInfo {
	return &certsInfo{
		DNS:     make(map[string]struct{}),
		IPs:     make(map[string]struct{}),
		Issuers: make(map[string]struct{}),
	}
}

func ParseCertificatesInfo(bundle []byte) (*CertificatesInfo, error) {
	ci := newCertsInfo()

	for {
		block, rest := pem.Decode(bundle)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			err := ci.parseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
		}
		bundle = rest
	}

	return ci.toCertificatesInfo(), nil
}

// This function is called for each certificate found in a PEM file.
// Populates *certsInfo with unique info about each certificate.
func (ci *certsInfo) parseCertificate(der []byte) error {
	crt, err := x509.ParseCertificate(der)
	if err != nil {
		return err
	}

	// Only keep the earliest expiration date.
	if ci.NotAfter.IsZero() || crt.NotAfter.Before(ci.NotAfter) {
		ci.NotAfter = crt.NotAfter
	}
	// Only keep the youngest NotBefore date.
	if crt.NotBefore.After(ci.NotBefore) {
		ci.NotBefore = crt.NotBefore
	}
	ci.Issuers[crt.Issuer.CommonName] = struct{}{}

	if !crt.IsCA {
		if crt.Subject.CommonName != "" {
			ci.DNS[crt.Subject.CommonName] = struct{}{}
		}
		// Alternate Subject Names
		for _, name := range crt.DNSNames {
			ci.DNS[name] = struct{}{}
		}
		// Certificates can accepts IP addresses/ranges (Fusion does this).
		for _, ip := range crt.IPAddresses {
			ci.IPs[ip.String()] = struct{}{}
		}
	}

	return nil
}

// Transform a dirty *certsInfo into a clean *CertificatesInfo for the API.
func (ci *certsInfo) toCertificatesInfo() *CertificatesInfo {
	csi := &CertificatesInfo{
		DNS:     strings.Join(mapKeys(ci.DNS), ", "),
		IPs:     strings.Join(mapKeys(ci.IPs), ", "),
		Issuers: strings.Join(mapKeys(ci.Issuers), ", "),
	}
	if !ci.NotAfter.IsZero() {
		csi.NotAfter = &ci.NotAfter
	}
	if !ci.NotBefore.IsZero() {
		csi.NotBefore = &ci.NotBefore
	}
	return csi
}

func mapKeys(m map[string]struct{}) []string {
	list := make([]string, len(m))
	i := 0
	for key := range m {
		list[i] = key
		i++
	}
	return list
}
