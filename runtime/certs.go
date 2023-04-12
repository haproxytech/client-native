package runtime

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"

	native_errors "github.com/haproxytech/client-native/v6/errors"
	"github.com/haproxytech/client-native/v6/models"
)

// parseCerts parses output from `show cert` command and return array of certificates
// First line in output represents format and is ignored
// Sample output format:
// /etc/ssl/cert-0.pem
// /etc/ssl/...
func (s *SingleRuntime) parseCerts(output string) models.SslCertificates {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil
	}
	certs := models.SslCertificates{}

	// ignore lines about open transactions
	ignore := false

	strings.SplitSeq(output, "\n")(func(line string) bool {
		if strings.HasPrefix(line, "# transaction") {
			ignore = true
			return true
		}
		if strings.HasPrefix(line, "# filename") {
			ignore = false
			return true
		}
		if !ignore {
			if c := s.parseCert(line); c != nil {
				certs = append(certs, c)
			}
		}
		return true
	})

	return certs
}

// parseCert parses one line from cert files array and return it structured
func (s *SingleRuntime) parseCert(line string) *models.SslCertificate {
	line = strings.TrimSpace(line)
	if line == "" || line[0] == '#' {
		return nil
	}
	if strings.HasPrefix(line, "\\*") {
		// Remove the leading antislash added by HAProxy.
		line = line[1:]
	}
	split := strings.Split(line, "/")
	return &models.SslCertificate{
		StorageName: line,
		Description: split[len(split)-1],
	}
}

// parseCertEntry parses one entry in one CrtList file and returns it structured
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
func parseCertEntry(response string) (*models.SslCertificate, error) {
	response = strings.TrimSpace(response)
	if response == "" {
		return nil, native_errors.ErrNotFound
	}

	c := &models.SslCertificate{}
	strings.SplitSeq(response, "\n")(func(line string) bool {
		key, val, found := strings.Cut(line, ": ")
		if !found {
			return true
		}
		switch strings.TrimSpace(key) {
		case "Filename":
			c.StorageName = val
		case "Status":
			c.Status = val
		case "Serial":
			c.Serial = val
		case "notBefore":
			notBefore, _ := time.Parse("Jan 2 15:04:05 2006 MST", val)
			c.NotBefore = (*strfmt.DateTime)(&notBefore)
		case "notAfter":
			notAfter, _ := time.Parse("Jan 2 15:04:05 2006 MST", val)
			c.NotAfter = (*strfmt.DateTime)(&notAfter)
		case "Subject Alternative Name":
			c.SubjectAlternativeNames = val
		case "Algorithm":
			c.Algorithm = val
		case "SHA1 FingerPrint":
			c.Sha1FingerPrint = val
		case "Subject":
			c.Subject = val
		case "Issuer":
			c.Issuers = val
		case "Chain Subject":
			c.ChainSubject = val
		case "Chain Issuer":
			c.ChainIssuer = val
		}
		return true
	})

	return c, nil
}

// ShowCerts returns cert files description from runtime
func (s *SingleRuntime) ShowCerts() (models.SslCertificates, error) {
	cmd := "show ssl cert"
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}
	return s.parseCerts(response), nil
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

func (s *SingleRuntime) ShowCertificate(storageName string) (*models.SslCertificate, error) {
	if storageName == "" {
		return nil, fmt.Errorf("%s %w", "Argument storageName empty", native_errors.ErrGeneral)
	}
	response, err := s.ExecuteWithResponse("show ssl cert " + storageName)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}
	return parseCertEntry(response)
}

func (s *SingleRuntime) NewCertificate(storageName string) error {
	if storageName == "" {
		return fmt.Errorf("%s %w", "Argument storageName empty", native_errors.ErrGeneral)
	}
	cmd := "new ssl cert " + storageName
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !strings.Contains(response, "New empty certificate store") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

func (s *SingleRuntime) SetCertificate(storageName string, payload string) error {
	if storageName == "" {
		return fmt.Errorf("%s %w", "Argument storageName empty", native_errors.ErrGeneral)
	}
	if payload == "" {
		return fmt.Errorf("%s %w", "Argument payload empty", native_errors.ErrGeneral)
	}
	cmd := fmt.Sprintf("set ssl cert %s <<\n%s\n", storageName, payload)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !transactionOK(response) {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

func (s *SingleRuntime) CommitCertificate(storageName string) error {
	if storageName == "" {
		return fmt.Errorf("%s %w", "Argument storageName empty", native_errors.ErrGeneral)
	}
	cmd := "commit ssl cert " + storageName
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !(strings.Contains(response, "Committing") && strings.Contains(response, "Success!")) {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

func (s *SingleRuntime) AbortCertificate(storageName string) error {
	if storageName == "" {
		return fmt.Errorf("%s %w", "Argument storageName empty", native_errors.ErrGeneral)
	}
	cmd := "abort ssl cert " + storageName
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !strings.Contains(response, "Transaction aborted for certificate") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

func (s *SingleRuntime) DeleteCertificate(storageName string) error {
	if storageName == "" {
		return fmt.Errorf("%s %w", "Argument storageName empty", native_errors.ErrGeneral)
	}
	cmd := "del ssl cert " + storageName
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if strings.Contains(response, "in use") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	if !(strings.Contains(response, "Certificate") && strings.Contains(response, "deleted!")) {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}
