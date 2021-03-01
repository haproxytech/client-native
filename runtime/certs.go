package runtime

import (
	"fmt"
	"strings"
	"time"

	native_errors "github.com/haproxytech/client-native/v2/errors"
	"github.com/haproxytech/client-native/v2/models"
)

type (
	SslCertEntries []*SslCertEntry
	SslCertEntry   struct {
		StorageName             string
		Status                  string
		Serial                  string
		NotBefore               time.Time
		NotAfter                time.Time
		SubjectAlternativeNames []string
		Algorithm               string
		SHA1FingerPrint         string
		Subject                 string
		Issuer                  string
		ChainSubject            string
		ChainIssuer             string
	}
)

// ShowCerts returns Certs files description from runtime
func (s *SingleRuntime) ShowCerts() (models.SslCertificates, error) {
	cmd := "show ssl cert"
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound) //nolint:errorlint
	}
	return s.parseCerts(response), nil
}

// parseCerts parses output from `show cert` command and return array of certificates
// First line in output represents format and is ignored
// Sample output format:
// /etc/ssl/cert-0.pem
// /etc/ssl/...
//
func (s *SingleRuntime) parseCerts(output string) models.SslCertificates {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil
	}
	certs := models.SslCertificates{}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		c := s.parseCert(line)
		if c != nil {
			certs = append(certs, c)
		}
	}
	return certs
}

// parseCert parses one line from cert files array and return it structured
func (s *SingleRuntime) parseCert(line string) *models.SslCertificate {
	if line == "" || strings.HasPrefix(strings.TrimSpace(line), "# filename") {
		return nil
	}
	split := strings.Split(line, "/")
	cert := &models.SslCertificate{
		StorageName: strings.TrimSpace(line),
		Description: split[len(split)-1],
	}
	return cert
}

// GetCert returns one structured runtime certs
func (s *SingleRuntime) GetCert(storageName string) (*models.SslCertificate, error) {
	if storageName == "" {
		return nil, fmt.Errorf("%s %w", "Argument storageName empty", native_errors.ErrGeneral)
	}
	certs, err := s.ShowCerts()
	if err != nil {
		return nil, err
	}

	for _, c := range certs {
		if c.StorageName == storageName {
			return c, nil
		}
	}
	return nil, fmt.Errorf("%s %w", storageName, native_errors.ErrNotFound)
}

// ShowCertEntry returns one CrtList runtime entries
func (s *SingleRuntime) ShowCertEntry(storageName string) (*SslCertEntry, error) {
	if storageName == "" {
		return nil, fmt.Errorf("%s %w", "Argument storageName empty", native_errors.ErrGeneral)
	}
	cmd := fmt.Sprintf("show ssl cert %s", storageName)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound) //nolint:errorlint
	}
	return parseCertEntry(response)
}

// parseCertEntry parses one entry in one CrtList file/runtime and returns it structured
// example:
// Filename: /etc/ssl/cert-2.pem
// Status: Used
// Serial: 0D933C1B1089BF660AE5253A245BB388
// notBefore: Sep  9 00:00:00 2020 GMT
// notAfter: Sep 14 12:00:00 2021 GMT
// Subject Alternative Name: DNS:*.platform.domain.com, DNS:uaa.platform.domain.com
// Algorithm: RSA4096
// SHA1 FingerPrint: 59242F1838BDEF3E7DAFC83FFE4DD6C03B88805C
// Subject: /C=DE/ST=Baden-Württemberg/L=Walldorf/O=ORG SE/CN=*.platform.domain.com
// Issuer: /C=US/O=DigiCert Inc/CN=DigiCert SHA2 Secure Server CA
// Chain Subject: /C=US/O=DigiCert Inc/CN=DigiCert SHA2 Secure Server CA
// Chain Issuer: /C=US/O=DigiCert Inc/OU=www.digicert.com/CN=DigiCert Global Root CA
func parseCertEntry(response string) (*SslCertEntry, error) {
	if response == "" || strings.HasPrefix(strings.TrimSpace(response), "#") {
		return nil, native_errors.ErrNotFound
	}

	c := &SslCertEntry{}
	parts := strings.Split(response, "\n")
	for _, p := range parts {
		index := strings.Index(p, ":")
		if index == -1 {
			continue
		}
		keyString := strings.TrimSpace(p[0:index])
		valueString := strings.TrimSpace(p[index+1:])

		switch key := keyString; {
		case key == "Filename":
			c.StorageName = valueString
		case key == "Status":
			c.Status = valueString
		case key == "Serial":
			c.Serial = valueString
		case key == "notBefore":
			c.NotBefore, _ = time.Parse("Jan 2 15:04:05 2006 MST", valueString)
		case key == "notAfter":
			c.NotAfter, _ = time.Parse("Jan 2 15:04:05 2006 MST", valueString)
		case key == "Subject Alternative Name":
			c.SubjectAlternativeNames = strings.Split(valueString, ", ")
		case key == "Algorithm":
			c.Algorithm = valueString
		case key == "SHA1 FingerPrint":
			c.SHA1FingerPrint = valueString
		case key == "Subject":
			c.Subject = valueString
		case key == "Issuer":
			c.Issuer = valueString
		case key == "Chain Subject":
			c.ChainSubject = valueString
		case key == "Chain Issuer":
			c.ChainIssuer = valueString
		}
	}

	return c, nil
}

// NewCertEntry adds an entry into the CrtList file
func (s *SingleRuntime) NewCertEntry(storageName string) error {
	if storageName == "" {
		return fmt.Errorf("%s %w", "Argument storageName empty", native_errors.ErrGeneral)
	}
	cmd := fmt.Sprintf("new ssl cert %s", storageName)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral) //nolint:errorlint
	}
	if !strings.Contains(response, "New empty certificate store") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// SetCertEntry adds an entry into the CrtList file
func (s *SingleRuntime) SetCertEntry(storageName string, payload string) error {
	if storageName == "" || payload == "" {
		return fmt.Errorf("%s %w", "Argument storageName or payload empty", native_errors.ErrGeneral)
	}
	cmd := fmt.Sprintf("set ssl cert %s <<\n%s\n", storageName, payload)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral) //nolint:errorlint
	}
	if !strings.Contains(response, "Transaction created for certificate") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// CommitCertEntry adds an entry into the CrtList file
func (s *SingleRuntime) CommitCertEntry(storageName string) error {
	if storageName == "" {
		return fmt.Errorf("%s %w", "Argument storageName empty", native_errors.ErrGeneral)
	}
	cmd := fmt.Sprintf("commit ssl cert %s", storageName)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral) //nolint:errorlint
	}
	if !(strings.Contains(response, "Committing") && strings.Contains(response, "Success!")) {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// AbortCertEntry adds an entry into the CrtList file
func (s *SingleRuntime) AbortCertEntry(storageName string) error {
	if storageName == "" {
		return fmt.Errorf("%s %w", "Argument storageName empty", native_errors.ErrGeneral)
	}
	cmd := fmt.Sprintf("abort ssl cert %s", storageName)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral) //nolint:errorlint
	}
	if !strings.Contains(response, "Transaction aborted for certificate") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// DeleteCertEntry adds an entry into the CrtList file
func (s *SingleRuntime) DeleteCertEntry(storageName string) error {
	if storageName == "" {
		return fmt.Errorf("%s %w", "Argument storageName empty", native_errors.ErrGeneral)
	}
	cmd := fmt.Sprintf("del ssl cert %s", storageName)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral) //nolint:errorlint
	}
	if !(strings.Contains(response, "Certificate") && strings.Contains(response, "deleted!")) {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}
