package configuration

import (
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetBackends returns a struct with configuration version and an array of
// configured backends. Returns error on fail.
func (c *LBCTLConfigurationClient) GetBackends() (*models.GetBackendsOKBody, error) {
	backendsString, err := c.executeLBCTL("l7-farm-dump", "")
	if err != nil {
		return nil, err
	}

	backends := c.parseBackends(backendsString)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetBackendsOKBody{Version: v, Data: backends}, nil
}

// GetBackend returns a struct with configuration version and a requested backend.
// Returns error on fail or if backend does not exist.
func (c *LBCTLConfigurationClient) GetBackend(name string) (*models.GetBackendOKBody, error) {
	backendStr, err := c.executeLBCTL("l7-farm-show", "", name)
	if err != nil {
		return nil, err
	}
	backend := &models.Backend{Name: name}

	c.parseObject(backendStr, backend)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetBackendOKBody{Version: v, Data: backend}, nil
}

// DeleteBackend deletes a backend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) DeleteBackend(name string, transactionID string, version int64) error {
	return c.deleteObject(name, "farm", "", "", transactionID, version)
}

// CreateBackend creates a backend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) CreateBackend(data *models.Backend, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}
	return c.createObject(data.Name, "farm", "", "", data, nil, transactionID, version)
}

// EditBackend edits a backend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) EditBackend(name string, data *models.Backend, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}
	ondiskBck, err := c.GetBackend(name)
	if err != nil {
		return err
	}
	return c.editObject(name, "farm", "", "", data, ondiskBck.Data, nil, transactionID, version)
}

func (c *LBCTLConfigurationClient) parseBackends(response string) models.Backends {
	backends := make(models.Backends, 0, 1)
	for _, backendStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(backendStr) == "" {
			continue
		}
		name := strings.TrimSpace(backendStr[strings.Index(backendStr, ".farm ")+6 : strings.Index(backendStr, "\n")])

		backendObj := &models.Backend{Name: name}
		c.parseObject(backendStr, backendObj)
		backends = append(backends, backendObj)
	}
	return backends
}
