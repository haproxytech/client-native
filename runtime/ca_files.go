package runtime

import (
	"fmt"
	"strings"

	native_errors "github.com/haproxytech/client-native/v6/errors"
	"github.com/haproxytech/client-native/v6/models"
)

// parseCAFiles parses output from `show ca files` command and return array of certificates
// First line in output represents format and is ignored
// Sample output format:
// cafile.pem - 3 certificates(3)
// @system - 1 certificates(1)
func (s *SingleRuntime) parseCAFiles(output string) models.SslCaFiles {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil
	}
	certs := models.SslCaFiles{}

	strings.SplitSeq(output, "\n")(func(line string) bool {
		if c := s.parseCAFile(line); c != nil {
			certs = append(certs, c)
		}
		return true
	})
	return certs
}

// parseCAFile parses one line from ca files array and return it structured
func (s *SingleRuntime) parseCAFile(line string) *models.SslCaFile {
	line = strings.TrimSpace(line)
	if line == "" || line[0] == '#' || line[0] == '*' {
		return nil
	}
	name, count, found := strings.Cut(line, " - ")
	if !found {
		return nil
	}
	return &models.SslCaFile{
		StorageName: strings.TrimSpace(name),
		Count:       strings.TrimSpace(count),
	}
}

// ShowCAFiles returns CA files description from runtime
func (s *SingleRuntime) ShowCAFiles() (models.SslCaFiles, error) {
	cmd := "show ssl ca-file"
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}
	return s.parseCAFiles(response), nil
}

// GetCAFile returns one CA file description
func (s *SingleRuntime) GetCAFile(caFile string) (*models.SslCaFile, error) {
	if caFile == "" {
		return nil, fmt.Errorf("%s %w", "Argument caFile empty", native_errors.ErrGeneral)
	}
	caFiles, err := s.ShowCAFiles()
	if err != nil {
		return nil, err
	}

	for _, c := range caFiles {
		if c.StorageName == caFile {
			return c, nil
		}
	}
	return nil, fmt.Errorf("%s %w", caFile, native_errors.ErrNotFound)
}

// ShowCAFile returns one CA file
func (s *SingleRuntime) ShowCAFile(caFile string, index *int64) (*models.SslCertificate, error) {
	if caFile == "" {
		return nil, fmt.Errorf("%s %w", "Argument caFile empty", native_errors.ErrGeneral)
	}
	if index != nil {
		caFile = fmt.Sprintf("%s:%d", caFile, *index)
	}
	response, err := s.ExecuteWithResponse("show ssl ca-file " + caFile)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}
	return parseCertEntry(response)
}

// NewCAFile creates a new empty CA file
func (s *SingleRuntime) NewCAFile(caFile string) error {
	if caFile == "" {
		return fmt.Errorf("%s %w", "Argument caFile empty", native_errors.ErrGeneral)
	}
	response, err := s.ExecuteWithResponse("new ssl ca-file " + caFile)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !strings.Contains(response, "New CA file created") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// SetCAFile sets a certificate payload to a CA file
func (s *SingleRuntime) SetCAFile(caFile, payload string) error {
	if caFile == "" {
		return fmt.Errorf("%s %w", "Argument caFile empty", native_errors.ErrGeneral)
	}
	if payload == "" {
		return fmt.Errorf("%s %w", "Argument payload empty", native_errors.ErrGeneral)
	}
	cmd := fmt.Sprintf("set ssl ca-file %s <<\n%s", caFile, payload)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !transactionOK(response) {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// CommitCAFile commits a CA file
func (s *SingleRuntime) CommitCAFile(caFile string) error {
	if caFile == "" {
		return fmt.Errorf("%s %w", "Argument caFile empty", native_errors.ErrGeneral)
	}
	response, err := s.ExecuteWithResponse("commit ssl ca-file " + caFile)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !strings.Contains(response, "Committing") || !strings.Contains(response, "Success!") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// AbortCAFile aborts and destroys a CA file update transaction
func (s *SingleRuntime) AbortCAFile(caFile string) error {
	if caFile == "" {
		return fmt.Errorf("%s %w", "Argument caFile empty", native_errors.ErrGeneral)
	}
	response, err := s.ExecuteWithResponse("abort ssl ca-file " + caFile)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !strings.Contains(response, "Transaction aborted for certificate") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// DeleteCAFile deletes a CA file
func (s *SingleRuntime) DeleteCAFile(caFile string) error {
	if caFile == "" {
		return fmt.Errorf("%s %w", "Argument caFile empty", native_errors.ErrGeneral)
	}
	response, err := s.ExecuteWithResponse("del ssl ca-file " + caFile)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if strings.Contains(response, "in use") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	if !strings.Contains(response, "CA file") || !strings.Contains(response, "deleted!") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// AddCAFileEntry adds an entry into the CA file
func (s *SingleRuntime) AddCAFileEntry(caFile, payload string) error {
	cmd := fmt.Sprintf("add ssl ca-file %s <<\n%s\n", caFile, payload)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !transactionOK(response) {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

func transactionOK(response string) bool {
	// The case is not always consistent in HAProxy.
	response = strings.ToLower(response)
	return strings.Contains(response, "transaction created") ||
		strings.Contains(response, "transaction updated")
}
