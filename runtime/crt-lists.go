package runtime

import (
	"fmt"
	"strconv"
	"strings"

	native_errors "github.com/haproxytech/client-native/v5/errors"
	"github.com/haproxytech/client-native/v5/models"
)

// These type aliases are provided for backward compatibility.
type CrtListEntry = models.SslCrtListEntry //nolint:gofumpt
type CrtListEntries = models.SslCrtListEntries
type CrtList = models.SslCrtList
type CrtLists = models.SslCrtLists

// ShowCrtLists returns CrtList files description from runtime
func (s *SingleRuntime) ShowCrtLists() (models.SslCrtLists, error) {
	response, err := s.ExecuteWithResponse("show ssl crt-list")
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}
	return s.parseCrtLists(response), nil
}

// parseCrtLists parses output from `show crt-list` command and return array of crt-list files
// First line in output represents format and is ignored
// Sample output format:
// /etc/ssl/crt-list
// /etc/ssl/...
func (s *SingleRuntime) parseCrtLists(output string) models.SslCrtLists {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil
	}

	lines := strings.Split(output, "\n")
	list := make(models.SslCrtLists, 0, len(lines))
	for _, line := range lines {
		if line != "" {
			list = append(list, &models.SslCrtList{File: line})
		}
	}

	return list
}

// GetCrtList returns one structured runtime CrtList file
func (s *SingleRuntime) GetCrtList(file string) (*models.SslCrtList, error) {
	crtLists, err := s.ShowCrtLists()
	if err != nil {
		return nil, err
	}

	for _, m := range crtLists {
		if m.File == file {
			return m, nil
		}
	}
	return nil, fmt.Errorf("%s %w", file, native_errors.ErrNotFound)
}

// ShowCrtListEntries returns one CrtList runtime entries
func (s *SingleRuntime) ShowCrtListEntries(file string) (models.SslCrtListEntries, error) {
	cmd := "show ssl crt-list -n " + file
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}
	return ParseCrtListEntries(response)
}

// ParseCrtListEntries parses array of entries in one CrtList file
// One line sample entry:
// /etc/ssl/cert-0.pem !*.crt-test.platform.domain.com !connectivitynotification.platform.domain.com !connectivitytunnel.platform.domain.com !authentication.cert.another.domain.com !*.authentication.cert.another.domain.com
// /etc/ssl/cert-1.pem [verify optional ca-file /etc/ssl/ca-file-1.pem] *.crt-test.platform.domain.com !connectivitynotification.platform.domain.com !connectivitytunnel.platform.domain.com !authentication.cert.another.domain.com !*.authentication.cert.another.domain.com
// /etc/ssl/cert-2.pem [verify required ca-file /etc/ssl/ca-file-2.pem]
func ParseCrtListEntries(output string) (models.SslCrtListEntries, error) {
	output = strings.TrimSpace(output)
	if output == "" || strings.HasPrefix(output, "didn't find the specified filename") {
		return nil, native_errors.ErrNotFound
	}

	lines := strings.Split(output, "\n")
	ce := make(models.SslCrtListEntries, 0, len(lines))
	for _, line := range lines {
		entry := parseCrtListEntry(line)
		if entry != nil {
			ce = append(ce, entry)
		}
	}
	return ce, nil
}

// parseCrtListEntry parses one entry in one CrtList file/runtime and returns it structured
// example:
// cert1.pem
// cert2.pem [alpn h2,http/1.1]
// certW.pem                   *.domain.tld !secure.domain.tld
// certS.pem [curves X25519:P-256 ciphers ECDHE-ECDSA-AES256-GCM-SHA384] secure.domain.tld
func parseCrtListEntry(line string) *models.SslCrtListEntry {
	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] == '#' {
		return nil
	}

	// The first word is the certificate with its line number.
	parts := strings.SplitN(line, " ", 2)
	if len(parts) == 0 {
		return nil
	}
	entry := models.SslCrtListEntry{}
	certLine := strings.Split(parts[0], ":")
	if len(certLine) > 2 || len(certLine) == 0 {
		return nil
	}
	entry.File = certLine[0]
	if len(certLine) == 2 {
		lineno, err := strconv.ParseInt(certLine[1], 10, 64)
		if err != nil || lineno < 1 {
			return nil
		}
		entry.LineNumber = lineno
	}

	if len(parts) == 1 {
		return &entry
	}

	rest := parts[1]

	// The next optional part is the SSL bind config.
	if len(rest) > 0 && rest[0] == '[' {
		end := strings.IndexByte(rest, ']')
		if end == -1 {
			return nil // should never happen
		}
		if end > 1 {
			entry.SSLBindConfig = rest[1:end]
		}
		if len(rest) > end {
			rest = rest[end+1:]
		} else {
			return &entry
		}
	}

	// The last optional part are the SNI filters.
	entry.SNIFilter = strings.Split(strings.TrimSpace(rest), " ")
	if len(entry.SNIFilter) == 1 && entry.SNIFilter[0] == "" {
		entry.SNIFilter = []string{}
	}

	return &entry
}

// AddCrtListEntry adds an entry into the CrtList file
func (s *SingleRuntime) AddCrtListEntry(crtList string, entry models.SslCrtListEntry) error {
	if crtList == "" {
		return fmt.Errorf("%s %w", "Argument crtList empty", native_errors.ErrGeneral)
	}
	if entry.File == "" {
		return fmt.Errorf("%s %w", "Filename empty", native_errors.ErrGeneral)
	}

	// The syntax of the command changes if any of those are set.
	extended := entry.SSLBindConfig != "" || len(entry.SNIFilter) > 0

	var sb strings.Builder
	sb.Grow(64)
	sb.WriteString("add ssl crt-list ")
	sb.WriteString(crtList)
	sb.WriteByte(' ')
	if extended {
		sb.WriteString("<<\n")
	}
	sb.WriteString(entry.File)
	if entry.SSLBindConfig != "" {
		sb.WriteString(" [")
		sb.WriteString(entry.SSLBindConfig)
		sb.WriteByte(']')
	}
	if len(entry.SNIFilter) > 0 {
		sb.WriteByte(' ')
		sb.WriteString(strings.Join(entry.SNIFilter, " "))
	}
	sb.WriteByte('\n')

	response, err := s.ExecuteWithResponse(sb.String())
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !strings.Contains(response, "Success") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// DeleteCrtListEntry deletes all the CrtList entries from the CrtList by its id
func (s *SingleRuntime) DeleteCrtListEntry(crtList, certFile string, lineNumber *int64) error {
	lineno := ""
	if lineNumber != nil {
		lineno = fmt.Sprintf(":%d", *lineNumber)
	}
	cmd := fmt.Sprintf("del ssl crt-list %s %s%s", crtList, certFile, lineno)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}
	if !strings.Contains(response, "deleted in crtlist") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}
