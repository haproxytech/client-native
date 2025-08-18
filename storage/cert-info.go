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
	"crypto/sha1" //nolint:gosec
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Information about stored certificates to be returned by the API.
type CertificatesInfo struct {
	NotAfter, NotBefore     *time.Time
	DNS, IPs, Issuers       string
	AuthorityKeyID          string
	SubjectKeyID            string
	Serial                  string
	Algorithm               string
	Sha1FingerPrint         string
	Sha256FingerPrint       string
	Subject                 string
	SubjectAlternativeNames string
}

// Private struct to store unique info about multiple certificates.
type certsInfo struct {
	NotAfter, NotBefore     time.Time
	DNS, IPs, Issuers       map[string]struct{}
	AuthorityKeyID          string
	SubjectKeyID            string
	Serial                  string
	Algorithm               string
	Sha1FingerPrint         string
	Sha256FingerPrint       string
	Subject                 string
	SubjectAlternativeNames []string
	Certs                   []*x509.Certificate
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

	crt, err := findLeafCertificate(ci.Certs)
	if err == nil {
		// Format the keys as OpenSSL does: hex digits in uppercase, colon-separated
		ci.AuthorityKeyID = formatFingerprint(crt.AuthorityKeyId)
		ci.SubjectKeyID = formatFingerprint(crt.SubjectKeyId)
		ci.Serial = crt.SerialNumber.String()
		ci.Algorithm = crt.SignatureAlgorithm.String()
		// Format the fingerprint as OpenSSL does: hex digits in uppercase, colon-separated
		fingerPrint := sha1.Sum(crt.Raw) //nolint:gosec
		ci.Sha1FingerPrint = formatFingerprint(fingerPrint[:])
		fingerPrint256 := sha256.Sum256(crt.Raw)
		ci.Sha256FingerPrint = formatFingerprint(fingerPrint256[:])
		ci.Subject = crt.Subject.CommonName
		ci.SubjectAlternativeNames = crt.DNSNames
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

	ci.Certs = append(ci.Certs, crt)

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
		ci.DNS[crt.Subject.CommonName] = struct{}{}
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

// formatFingerprint formats a byte array as: hex digits in uppercase, colon-separated
func formatFingerprint(fingerprint []byte) string {
	parts := make([]string, len(fingerprint))
	for i, b := range fingerprint {
		parts[i] = fmt.Sprintf("%02X", b)
	}
	return strings.Join(parts, ":")
}

// findLeafCertificate returns the first leaf certificate in the chain.
func findLeafCertificate(certs []*x509.Certificate) (*x509.Certificate, error) {
	if len(certs) == 0 {
		return nil, errors.New("empty certificate chain")
	}
	if len(certs) == 1 {
		return certs[0], nil
	}

	// Create a map to check if a certificate is someone else's issuer
	isIssuer := make(map[string]bool)
	for _, cert := range certs {
		isIssuer[cert.Issuer.String()] = true
	}

	// Find the starting certificate (a certificate whose issuer is not in the list)
	for _, cert := range certs {
		if !cert.IsCA && (cert.Subject.CommonName != "" || len(cert.DNSNames) != 0) && !isIssuer[cert.Subject.String()] {
			return cert, nil
		}
	}

	return nil, errors.New("no leaf certificate found")
}

// Transform a dirty *certsInfo into a clean *CertificatesInfo for the API.
func (ci *certsInfo) toCertificatesInfo() *CertificatesInfo {
	csi := &CertificatesInfo{
		DNS:                     strings.Join(mapKeys(ci.DNS), ", "),
		IPs:                     strings.Join(mapKeys(ci.IPs), ", "),
		Issuers:                 strings.Join(mapKeys(ci.Issuers), ", "),
		AuthorityKeyID:          ci.AuthorityKeyID,
		SubjectKeyID:            ci.SubjectKeyID,
		Serial:                  ci.Serial,
		Algorithm:               ci.Algorithm,
		Sha1FingerPrint:         ci.Sha1FingerPrint,
		Sha256FingerPrint:       ci.Sha256FingerPrint,
		Subject:                 ci.Subject,
		SubjectAlternativeNames: strings.Join(ci.SubjectAlternativeNames, ", "),
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
