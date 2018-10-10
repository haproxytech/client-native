package configuration

import (
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetServers returns a struct with configuration version and an array of
// configured servers in the specified backend. Returns error on fail.
func (c *LBCTLConfigurationClient) GetServers(backend string, transactionID string) (*models.GetServersOKBody, error) {
	serversString, err := c.executeLBCTL("l7-server-dump", transactionID, backend)
	if err != nil {
		return nil, err
	}

	servers := c.parseServers(serversString)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetServersOKBody{Version: v, Data: servers}, nil
}

// GetServer returns a struct with configuration version and a requested server
// in the specified backend. Returns error on fail or if server does not exist.
func (c *LBCTLConfigurationClient) GetServer(name string, backend string, transactionID string) (*models.GetServerOKBody, error) {
	serverStr, err := c.executeLBCTL("l7-server-show", transactionID, backend, name)
	if err != nil {
		return nil, err
	}
	server := &models.Server{Name: name}

	c.parseObject(serverStr, server)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetServerOKBody{Version: v, Data: server}, nil
}

// DeleteServer deletes a server in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) DeleteServer(name string, backend string, transactionID string, version int64) error {
	return c.deleteObject(name, "server", backend, "", transactionID, version)
}

// CreateServer creates a server in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) CreateServer(backend string, data *models.Server, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}
	return c.createObject(data.Name, "server", backend, "", data, nil, transactionID, version)
}

// EditServer edits a server in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLConfigurationClient) EditServer(name string, backend string, data *models.Server, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return NewConfError(ErrValidationError, validationErr.Error())
	}
	ondiskSrv, err := c.GetServer(name, backend, transactionID)
	if err != nil {
		return err
	}

	return c.editObject(name, "server", backend, "", data, ondiskSrv.Data, nil, transactionID, version)
}

func (c *LBCTLConfigurationClient) parseServers(response string) models.Servers {
	servers := make(models.Servers, 0, 1)
	for _, serverStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(serverStr) == "" {
			continue
		}
		name, _ := splitHeaderLine(serverStr)

		serverObj := &models.Server{Name: name}
		c.parseObject(serverStr, serverObj)
		servers = append(servers, serverObj)
	}
	return servers
}
