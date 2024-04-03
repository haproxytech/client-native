package runtime

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	native_errors "github.com/haproxytech/client-native/v5/errors"
)

type CrtLists []*CrtList

type CrtList struct {
	File string
}

type CrtListEntries []*CrtListEntry

type CrtListEntry struct {
	File          string
	SSLBindConfig string
	SNIFilter     []string
	LineNumber    int
}

// ShowCrtLists returns CrtList files description from runtime
func (s *SingleRuntime) ShowCrtLists() (CrtLists, error) {
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
func (s *SingleRuntime) parseCrtLists(output string) CrtLists {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil
	}
	crtLists := CrtLists{}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		c := s.parseCrtList(line)
		if c != nil {
			crtLists = append(crtLists, c)
		}
	}
	return crtLists
}

// parseCrtList parses one line from CrtList files array and return it structured
func (s *SingleRuntime) parseCrtList(line string) *CrtList {
	if line == "" {
		return nil
	}
	crtList := &CrtList{
		File: line,
	}
	return crtList
}

// GetCrtList returns one structured runtime CrtList file
func (s *SingleRuntime) GetCrtList(file string) (*CrtList, error) {
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
func (s *SingleRuntime) ShowCrtListEntries(file string) (CrtListEntries, error) {
	cmd := fmt.Sprintf("show ssl crt-list -n %s", file)
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
func ParseCrtListEntries(output string) (CrtListEntries, error) {
	output = strings.TrimSpace(output)
	if output == "" || strings.HasPrefix(output, "didn't find the specified filename") {
		return nil, native_errors.ErrNotFound
	}
	ce := CrtListEntries{}

	lines := strings.Split(strings.TrimSpace(output), "\n")
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
func parseCrtListEntry(line string) *CrtListEntry {
	if line == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
		return nil
	}

	c := &CrtListEntry{}
	re := regexp.MustCompile(`(\S+)(?:\s\[(.*)\])?(?:\s(.*))?`)
	matches := re.FindStringSubmatch(line)
	if matches != nil {
		split := strings.Split(matches[1], ":")
		linenumber, _ := strconv.ParseInt(split[1], 0, 32)
		c.LineNumber = int(linenumber)
		c.File = split[0]
		c.SSLBindConfig = matches[2]
		c.SNIFilter = strings.Fields(matches[3])
	}

	return c
}

// AddCrtListEntry adds an entry into the CrtList file
func (s *SingleRuntime) AddCrtListEntry(crtList string, entry CrtListEntry) error {
	cmd := fmt.Sprintf("add ssl crt-list %s <<\n%s", crtList, entry.File)
	if entry.SSLBindConfig != "" {
		cmd = fmt.Sprintf("%s [%s]", cmd, entry.SSLBindConfig)
	}
	for _, sni := range entry.SNIFilter {
		cmd = fmt.Sprintf("%s %s", cmd, sni)
	}
	cmd += "\n"
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral)
	}
	if !strings.Contains(response, "Success") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}

// DeleteCrtListEntry deletes all the CrtList entries from the CrtList by its id
func (s *SingleRuntime) DeleteCrtListEntry(crtList, certFile string, lineNumber int) error {
	cmd := fmt.Sprintf("del ssl crt-list %s %s:%v", crtList, certFile, lineNumber)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}
	if !strings.Contains(response, "deleted in crtlist") {
		return fmt.Errorf("%s %w", response, native_errors.ErrGeneral)
	}
	return nil
}
