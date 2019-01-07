package configuration

import (
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetFrontends returns a struct with configuration version and an array of
// configured frontends. Returns error on fail.
func (c *Client) GetFrontends(transactionID string) (*models.GetFrontendsOKBody, error) {
	if c.Cache.Enabled() {
		frontends, found := c.Cache.Frontends.Get(transactionID)
		if found {
			return &models.GetFrontendsOKBody{Version: c.Cache.Version.Get(), Data: frontends}, nil
		}
	}
	frontendsStr, err := c.executeLBCTL("l7-service-dump", transactionID)
	if err != nil {
		return nil, err
	}
	frontends := c.parseFrontends(frontendsStr)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.Frontends.SetAll(transactionID, frontends)
	}
	return &models.GetFrontendsOKBody{Version: v, Data: frontends}, nil
}

// GetFrontend returns a struct with configuration version and a requested frontend.
// Returns error on fail or if frontend does not exist.
func (c *Client) GetFrontend(name string, transactionID string) (*models.GetFrontendOKBody, error) {
	if c.Cache.Enabled() {
		frontend, found := c.Cache.Frontends.GetOne(name, transactionID)
		if found {
			return &models.GetFrontendOKBody{Version: c.Cache.Version.Get(), Data: frontend}, nil
		}
	}
	frontendStr, err := c.executeLBCTL("l7-service-show", transactionID, name)
	if err != nil {
		return nil, err
	}
	frontend := &models.Frontend{Name: name}

	c.parseObject(frontendStr, frontend)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.Frontends.Set(frontend.Name, transactionID, frontend)
	}
	return &models.GetFrontendOKBody{Version: v, Data: frontend}, nil
}

// DeleteFrontend deletes a frontend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteFrontend(name string, transactionID string, version int64) error {
	err := c.deleteObject(name, "service", "", "", transactionID, version)
	if err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.DeleteFrontendCache(name, transactionID)
	}
	return err
}

// EditFrontend edits a frontend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditFrontend(name string, data *models.Frontend, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	ondiskFrontend, err := c.GetFrontend(name, transactionID)
	if err != nil {
		return err
	}
	err = c.editObject(name, "service", "", "", data, ondiskFrontend.Data, nil, transactionID, version)
	if err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.Frontends.Set(name, transactionID, data)
	}
	return nil
}

// CreateFrontend creates a frontend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateFrontend(data *models.Frontend, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	err := c.createObject(data.Name, "service", "", "", data, nil, transactionID, version)
	if err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.Frontends.Set(data.Name, transactionID, data)
	}
	return nil
}

func (c *Client) parseFrontends(response string) models.Frontends {
	frontends := make(models.Frontends, 0, 1)
	for _, frontendStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(frontendStr) == "" {
			continue
		}
		name := strings.TrimSpace(frontendStr[strings.Index(frontendStr, ".service ")+9 : strings.Index(frontendStr, "\n")])

		frontendObj := &models.Frontend{Name: name}
		c.parseObject(frontendStr, frontendObj)
		frontends = append(frontends, frontendObj)
	}
	return frontends
}
