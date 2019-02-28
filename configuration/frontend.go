package configuration

import (
	"fmt"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	"github.com/haproxytech/models"
)

// GetFrontends returns a struct with configuration version and an array of
// configured frontends. Returns error on fail.
func (c *Client) GetFrontends(transactionID string) (*models.GetFrontendsOKBody, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return nil, err
	}

	fNames, err := p.SectionsGet(parser.Frontends)
	if err != nil {
		return nil, err
	}

	frontends := []*models.Frontend{}
	for _, name := range fNames {
		f := &models.Frontend{Name: name}
		if err := c.parseSection(f, parser.Frontends, name, p); err != nil {
			continue
		}
		frontends = append(frontends, f)
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	return &models.GetFrontendsOKBody{Version: v, Data: frontends}, nil
}

// GetFrontend returns a struct with configuration version and a requested frontend.
// Returns error on fail or if frontend does not exist.
func (c *Client) GetFrontend(name string, transactionID string) (*models.GetFrontendOKBody, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return nil, err
	}

	if !c.checkSectionExists(parser.Frontends, name, p) {
		return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Frontend %s does not exist", name))
	}

	frontend := &models.Frontend{Name: name}
	if err := c.parseSection(frontend, parser.Frontends, name, p); err != nil {
		return nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	return &models.GetFrontendOKBody{Version: v, Data: frontend}, nil
}

// DeleteFrontend deletes a frontend in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteFrontend(name string, transactionID string, version int64) error {
	if err := c.deleteSection(parser.Frontends, name, transactionID, version); err != nil {
		return err
	}
	return nil
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

	if err := c.editSection(parser.Frontends, name, data, transactionID, version); err != nil {
		return err
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

	if err := c.createSection(parser.Frontends, data.Name, data, transactionID, version); err != nil {
		return err
	}

	return nil
}
