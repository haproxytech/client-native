package configuration

import (
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetFrontends returns a struct with configuration version and an array of
// configured frontends. Returns error on fail.
func (c *LBCTLConfigurationClient) GetFrontends(transactionID string) (*models.GetFrontendsOKBody, error) {
	frontendsStr, err := c.executeLBCTL("l7-service-dump", transactionID)
	if err != nil {
		return nil, err
	}
	frontends := c.parseFrontends(frontendsStr)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetFrontendsOKBody{Version: v, Data: frontends}, nil
}

// GetFrontend returns a struct with configuration version and a requested frontend.
// Returns error on fail or if frontend does not exist.
func (c *LBCTLConfigurationClient) GetFrontend(name string, transactionID string) (*models.GetFrontendOKBody, error) {
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

	return &models.GetFrontendOKBody{Version: v, Data: frontend}, nil
}

// DeleteFrontend deletes a frontend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) DeleteFrontend(name string, transactionID string, version int64) error {
	return c.deleteObject(name, "service", "", "", transactionID, version)
}

// EditFrontend edits a frontend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) EditFrontend(name string, data *models.Frontend, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}
	ondiskFrontend, err := c.GetFrontend(name, transactionID)
	if err != nil {
		return err
	}
	return c.editObject(name, "service", "", "", data, ondiskFrontend.Data, nil, transactionID, version)
}

// CreateFrontend creates a frontend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) CreateFrontend(data *models.Frontend, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}
	return c.createObject(data.Name, "service", "", "", data, nil, transactionID, version)
}

func (c *LBCTLConfigurationClient) parseFrontends(response string) models.Frontends {
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
