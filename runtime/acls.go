package runtime

import (
	"fmt"
	"regexp"
	"strings"

	native_errors "github.com/haproxytech/client-native/v2/errors"
	"github.com/haproxytech/client-native/v2/models"
)

// ShowACLS returns Acl files description from runtime
func (s *SingleRuntime) ShowACLS() (models.ACLFiles, error) {
	response, err := s.ExecuteWithResponse("show acl")
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound) //nolint:errorlint
	}
	return s.parseACLS(response), nil
}

// parseACLS parses output from `show acl` command and return array of acl files
// First line in output represents format and is ignored
// Sample output format:
// # id (file) description
// -0 (/etc/acl/blocklist.txt) pattern loaded from file '/etc/acl/blocklist.txt' used by acl at file '/usr/local/etc/haproxy/haproxy.cfg' line 59
// -1 () acl 'src' file '/usr/local/etc/haproxy/haproxy.cfg' line 59
func (s *SingleRuntime) parseACLS(output string) models.ACLFiles {
	if output == "" {
		return nil
	}
	acls := models.ACLFiles{}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		m := s.parseACL(line)
		if m != nil {
			acls = append(acls, m)
		}
	}
	return acls
}

// parseACL parses one line from ACL files array and return it structured
func (s *SingleRuntime) parseACL(line string) *models.ACLFile {
	if line == "" || strings.HasPrefix(strings.TrimSpace(line), "# id") {
		return nil
	}

	parts := strings.Fields(line)
	if len(parts) < 3 {
		return nil
	}

	m := &models.ACLFile{
		ID:          parts[0],
		StorageName: findStorageName(parts[1], line),
		Description: strings.Join(parts[2:], " "),
	}

	return m
}

// findStorageName checks if acl name exists and extracts it
func findStorageName(name, line string) string {
	name = strings.TrimSuffix(strings.TrimPrefix(name, "("), ")")
	if name == "" {
		re := regexp.MustCompile(`acl\s'(.*)'\sfile`)
		matches := re.FindStringSubmatch(line)
		if matches != nil {
			name = matches[1]
		}
	}
	return name
}

// GetACL returns one structured runtime Acl file
func (s *SingleRuntime) GetACL(storageName string) (*models.ACLFile, error) {
	if storageName == "" {
		return nil, fmt.Errorf("%s %w", "Argument nameOrFile empty", native_errors.ErrGeneral)
	}
	acls, err := s.ShowACLS()
	if err != nil {
		return nil, err
	}

	for _, m := range acls {
		if m.StorageName == storageName {
			return m, nil
		}
	}
	return nil, fmt.Errorf("%s %w", storageName, native_errors.ErrNotFound)
}

// ShowACLFileEntries returns one acl runtime entries
func (s *SingleRuntime) ShowACLFileEntries(storageName string) (models.ACLFilesEntries, error) {
	if storageName == "" {
		return nil, fmt.Errorf("%s %w", "Argument file empty", native_errors.ErrGeneral)
	}
	cmd := fmt.Sprintf("show acl %s", storageName)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound) //nolint:errorlint
	}
	return ParseACLFileEntries(response, true)
}

// ParseACLFileEntries parses array of entries in one Acl file
// One line sample entry:
// ID			  Value
// 0x560f3f9e8600 10.178.160.0
func ParseACLFileEntries(output string, hasID bool) (models.ACLFilesEntries, error) {
	if output == "" || strings.HasPrefix(strings.TrimSpace(output), "Unknown ACL identifier.") {
		return nil, native_errors.ErrNotFound
	}
	me := models.ACLFilesEntries{}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		e := parseACLFileEntry(line, hasID)
		if e != nil {
			me = append(me, e)
		}
	}
	return me, nil
}

// parseACLFileEntry parses one entry in one Acl file/runtime and returns it structured
func parseACLFileEntry(line string, hasID bool) *models.ACLFileEntry {
	if line == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
		return nil
	}

	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil
	}

	m := &models.ACLFileEntry{}
	if hasID {
		m.ID = parts[0] // acl entries from runtime have ID
		m.Value = parts[1]
	} else {
		m.Value = parts[0]
	}
	return m
}

// AddACLFileEntry adds an entry into the Acl file
func (s *SingleRuntime) AddACLFileEntry(aclID, value string) error {
	if aclID == "" || value == "" {
		return fmt.Errorf("%s %w", "One or more Arguments empty", native_errors.ErrGeneral)
	}
	m, _ := s.GetACLFileEntry(aclID, value)
	if m != nil {
		return fmt.Errorf("%w", native_errors.ErrAlreadyExists)
	}
	cmd := fmt.Sprintf("add acl #%s %s", aclID, value)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrGeneral) //nolint:errorlint
	}
	if strings.Contains(response, "not") && strings.Contains(response, "valid") {
		return fmt.Errorf("%s %w", strings.TrimSpace(response), native_errors.ErrGeneral)
	}
	return nil
}

// GetACLFileEntry returns one Acl runtime setting
func (s *SingleRuntime) GetACLFileEntry(aclID, value string) (*models.ACLFileEntry, error) {
	cmd := fmt.Sprintf("get acl #%s %s", aclID, value)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound) //nolint:errorlint
	}

	matched := false
	m := &models.ACLFileEntry{}
	parts := strings.Split(response, ",")
	for _, p := range parts {
		kv := strings.Split(p, "=")
		switch key := strings.TrimSpace(kv[0]); {
		case key == "pattern":
			m.Value = strings.Trim(strings.TrimSpace(kv[1]), "\"")
		case key == "match":
			matched = true
		}
	}

	if m.Value == "" || !matched {
		return nil, fmt.Errorf("%s %w", value, native_errors.ErrNotFound)
	}
	return m, nil
}

// DeleteACLFileEntry deletes all the Acl entries from the Acl by its value
func (s *SingleRuntime) DeleteACLFileEntry(aclID, value string) error {
	if aclID == "" || value == "" {
		return fmt.Errorf("%s %w", "One or more Arguments empty", native_errors.ErrGeneral)
	}
	cmd := fmt.Sprintf("del acl #%s %s", aclID, value)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound) //nolint:errorlint
	}
	if strings.Contains(response, "not") && strings.Contains(response, "found") {
		return fmt.Errorf("%s %w", strings.TrimSpace(response), native_errors.ErrGeneral)
	}
	return nil
}
