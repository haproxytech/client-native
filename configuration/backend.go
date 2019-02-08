package configuration

import (
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetBackends returns a struct with configuration version and an array of
// configured backends. Returns error on fail.
func (c *Client) GetBackends(transactionID string) (*models.GetBackendsOKBody, error) {
	if c.Cache.Enabled() {
		backends, found := c.Cache.Backends.Get(transactionID)
		if found {
			return &models.GetBackendsOKBody{Version: c.Cache.Version.Get(transactionID), Data: backends}, nil
		}
	}

	backendsString, err := c.executeLBCTL("l7-farm-dump", transactionID)
	if err != nil {
		return nil, err
	}

	backends := c.parseBackends(backendsString)

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.Backends.SetAll(transactionID, backends)
	}
	return &models.GetBackendsOKBody{Version: v, Data: backends}, nil
}

// GetBackend returns a struct with configuration version and a requested backend.
// Returns error on fail or if backend does not exist.
func (c *Client) GetBackend(name string, transactionID string) (*models.GetBackendOKBody, error) {
	if c.Cache.Enabled() {
		backend, found := c.Cache.Backends.GetOne(name, transactionID)
		if found {
			return &models.GetBackendOKBody{Version: c.Cache.Version.Get(transactionID), Data: backend}, nil
		}
	}

	backendStr, err := c.executeLBCTL("l7-farm-show", transactionID, name)
	if err != nil {
		return nil, err
	}
	backend := &models.Backend{Name: name}

	c.parseObject(backendStr, backend)

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.Backends.Set(name, transactionID, backend)
	}
	return &models.GetBackendOKBody{Version: v, Data: backend}, nil
}

// DeleteBackend deletes a backend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteBackend(name string, transactionID string, version int64) error {
	err := c.deleteObject(name, "farm", "", "", transactionID, version)
	if err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.DeleteBackendCache(name, transactionID)
	}
	return nil
}

// CreateBackend creates a backend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateBackend(data *models.Backend, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	err := c.createObject(data.Name, "farm", "", "", data, nil, transactionID, version)
	if err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.Backends.Set(data.Name, transactionID, data)
	}
	return nil
}

// EditBackend edits a backend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditBackend(name string, data *models.Backend, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	ondiskBck, err := c.GetBackend(name, transactionID)
	if err != nil {
		return err
	}
	err = c.editObject(name, "farm", "", "", data, ondiskBck.Data, nil, transactionID, version)
	if err != nil {
		return err
	}

	if c.Cache.Enabled() {
		c.Cache.Backends.Set(name, transactionID, data)
	}
	return nil
}

func (c *Client) parseBackends(response string) models.Backends {
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
