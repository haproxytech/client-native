package runtime

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"

	native_errors "github.com/haproxytech/client-native/v6/errors"
	"github.com/haproxytech/client-native/v6/models"
)

// parseCrlEntries parses one or more CRL entries.
// example:
// Filename: /etc/ssl/crlfile.pem
// Status: Used
// Certificate Revocation List :
// Version 1
// Signature Algorithm: sha256WithRSAEncryption
// Issuer: /C=FR/O=HAProxy Technologies/CN=Intermediate CA2
// Last Update: Apr 1 14:45:39 2023 GMT
// Next Update: Sep 8 14:45:39 2048 GMT
// Revoked Certificates:
//
//	Serial Number: 1008
//	  Revocation Date: Apr 1 14:45:36 2023 GMT
func parseCrlEntries(response string) (*models.SslCrlEntries, error) {
	response = strings.TrimSpace(response)
	if response == "" {
		return nil, native_errors.ErrNotFound
	}

	entries := models.SslCrlEntries{}
	var entry *models.SslCrlEntry
	var filename string
	var status string
	var serialNumber string
	var parseErr error

	strings.SplitSeq(response, "\n")(func(line string) bool {
		line = strings.TrimSpace(line)
		if line == "" {
			return true
		}
		key, val, found := strings.Cut(line, ":")
		if !found {
			if strings.HasPrefix(line, "Version ") {
				_, val, _ = strings.Cut(line, " ")
				if entry != nil {
					entry.Version = val
				}
			}
			return true
		}

		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)

		if strings.HasPrefix(key, "Certificate Revocation List") {
			// Save the current entry and start a new one.
			if entry != nil {
				entries = append(entries, entry)
			}
			entry = &models.SslCrlEntry{
				StorageName: filename,
				Status:      status,
			}
			return true
		}

		switch key {
		// Filename and Status are common to all entries.
		case "Filename":
			filename = val
		case "Status":
			status = val

		case "Issuer":
			entry.Issuer = val
		case "Signature Algorithm":
			entry.SignatureAlgorithm = val
		case "Last Update", "This Update":
			update, err := time.Parse("Jan 2 15:04:05 2006 GMT", val)
			if err != nil {
				parseErr = fmt.Errorf("cannot parse crl last update : %s %w", err.Error(), native_errors.ErrGeneral)
				return false
			}
			entry.LastUpdate = strfmt.Date(update)
		case "Next Update":
			update, err := time.Parse("Jan 2 15:04:05 2006 GMT", val)
			if err != nil {
				parseErr = fmt.Errorf("cannot parse crl next update : %s %w", err.Error(), native_errors.ErrGeneral)
				return false
			}
			entry.NextUpdate = strfmt.Date(update)
		case "Serial Number":
			serialNumber = val
		case "Revocation Date":
			revocationDate, err := time.Parse("Jan 2 15:04:05 2006 GMT", val)
			if err != nil {
				parseErr = fmt.Errorf("cannot parse crl revocation date : %s %w", err.Error(), native_errors.ErrGeneral)
				return false
			}
			rev := &models.RevokedCertificates{
				RevocationDate: strfmt.Date(revocationDate),
				SerialNumber:   serialNumber,
			}
			entry.RevokedCertificates = append(entry.RevokedCertificates, rev)
		}

		return true
	})

	if parseErr != nil {
		return nil, parseErr
	}

	if entry != nil {
		entries = append(entries, entry)
	}

	return &entries, nil
}

// parseCrl parses one line from crl files array and return it structured
func parseCrl(line string) *models.SslCrl {
	line = strings.TrimSpace(line)
	if line == "" || line[0] == '#' || line[0] == '*' {
		return nil
	}
	split := strings.Split(line, "/")
	crl := &models.SslCrl{
		StorageName: strings.TrimSpace(line),
		Description: split[len(split)-1],
	}
	return crl
}

// parseCrls parses output from `show ssl crl-file` command and return array of crls
// First line in output represents format and is ignored
// Sample output format:
// /etc/ssl/crl-0.pem
// /etc/ssl/...
func parseCrls(output string) models.SslCrls {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil
	}

	lines := strings.Split(output, "\n")
	crls := make(models.SslCrls, 0, len(lines))
	for _, line := range lines {
		c := parseCrl(line)
		if c != nil {
			crls = append(crls, c)
		}
	}
	return crls
}

// ShowCrlFiles returns Crl files description from runtime
func (s *SingleRuntime) ShowCrlFiles() (models.SslCrls, error) {
	cmd := "show ssl crl-file"
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}
	return parseCrls(response), nil
}

// GetCrlFile returns a Crl file description
func (s *SingleRuntime) GetCrlFile(crlFile string) (*models.SslCrl, error) {
	if crlFile == "" {
		return nil, fmt.Errorf("%s %w", "Argument crlFile empty", native_errors.ErrGeneral)
	}
	caFiles, err := s.ShowCrlFiles()
	if err != nil {
		return nil, err
	}

	for _, c := range caFiles {
		if c.StorageName == crlFile {
			return c, nil
		}
	}
	return nil, fmt.Errorf("%s %w", crlFile, native_errors.ErrNotFound)
}

// ShowCrlFile returns one or more entries in a Crl file
func (s *SingleRuntime) ShowCrlFile(crlFile string, index *int64) (*models.SslCrlEntries, error) {
	if crlFile == "" {
		return nil, fmt.Errorf("%s %w", "Argument crlFile empty", native_errors.ErrGeneral)
	}
	if index != nil {
		crlFile = fmt.Sprintf("%s:%d", crlFile, *index)
	}
	response, err := s.ExecuteWithResponse("show ssl crl-file " + crlFile)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}
	return parseCrlEntries(response)
}

// NewCrlFile creates a new empty Crl
func (s *SingleRuntime) NewCrlFile(crlFile string) error {
	if crlFile == "" {
		return fmt.Errorf("%s %w", "Argument crlFile empty", native_errors.ErrGeneral)
	}
	response, err := s.ExecuteWithResponse("new ssl crl-file " + crlFile)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !strings.Contains(response, "New CRL file created") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// SetCrlFile sets a payload to a Crl file
func (s *SingleRuntime) SetCrlFile(crlFile string, payload string) error {
	if crlFile == "" {
		return fmt.Errorf("%s %w", "Argument crlFile empty", native_errors.ErrGeneral)
	}
	if payload == "" {
		return fmt.Errorf("%s %w", "Argument payload empty", native_errors.ErrGeneral)
	}
	cmd := fmt.Sprintf("set ssl crl-file %s <<\n%s\n", crlFile, payload)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !transactionOK(response) {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// CommitCrlFile commits a Crl file
func (s *SingleRuntime) CommitCrlFile(crlFile string) error {
	if crlFile == "" {
		return fmt.Errorf("%s %w", "Argument crlFile empty", native_errors.ErrGeneral)
	}
	response, err := s.ExecuteWithResponse("commit ssl crl-file " + crlFile)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !(strings.Contains(response, "Committing") && strings.Contains(response, "Success!")) {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// AbortCrlFile aborts and destroys a Crl file update transaction
func (s *SingleRuntime) AbortCrlFile(crlFile string) error {
	if crlFile == "" {
		return fmt.Errorf("%s %w", "Argument crlFile empty", native_errors.ErrGeneral)
	}
	response, err := s.ExecuteWithResponse("abort ssl crl-file " + crlFile)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !strings.Contains(response, "Transaction aborted") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// DeleteCrlFile deletes a Crl file
func (s *SingleRuntime) DeleteCrlFile(crlFile string) error {
	if crlFile == "" {
		return fmt.Errorf("%s %w", "Argument crlFile empty", native_errors.ErrGeneral)
	}
	response, err := s.ExecuteWithResponse("del ssl crl-file " + crlFile)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !(strings.Contains(response, "CRL file") && strings.Contains(response, "deleted!")) {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}
