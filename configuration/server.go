package configuration

import (
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetServers returns a struct with configuration version and an array of
// configured servers in the specified backend. Returns error on fail.
func (c *Client) GetServers(backend string, transactionID string) (*models.GetServersOKBody, error) {
	if c.Cache.Enabled() {
		servers, found := c.Cache.Servers.Get(backend, transactionID)
		if found {
			return &models.GetServersOKBody{Version: c.Cache.Version.Get(transactionID), Data: servers}, nil
		}
	}
	serversString, err := c.executeLBCTL("l7-server-dump", transactionID, backend)
	if err != nil {
		return nil, err
	}

	servers := c.parseServers(serversString)

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.Servers.SetAll(backend, transactionID, servers)
	}
	return &models.GetServersOKBody{Version: v, Data: servers}, nil
}

// GetServer returns a struct with configuration version and a requested server
// in the specified backend. Returns error on fail or if server does not exist.
func (c *Client) GetServer(name string, backend string, transactionID string) (*models.GetServerOKBody, error) {
	if c.Cache.Enabled() {
		server, found := c.Cache.Servers.GetOne(name, backend, transactionID)
		if found {
			return &models.GetServerOKBody{Version: c.Cache.Version.Get(transactionID), Data: server}, nil
		}
	}
	serverStr, err := c.executeLBCTL("l7-server-show", transactionID, backend, name)
	if err != nil {
		return nil, err
	}
	server := &models.Server{Name: name}

	c.parseObject(serverStr, server)

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.Servers.Set(name, backend, transactionID, server)
	}
	return &models.GetServerOKBody{Version: v, Data: server}, nil
}

// DeleteServer deletes a server in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteServer(name string, backend string, transactionID string, version int64) error {
	err := c.deleteObject(name, "server", backend, "", transactionID, version)
	if err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.Servers.Delete(name, backend, transactionID)
	}
	return nil
}

// CreateServer creates a server in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateServer(backend string, data *models.Server, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	err := c.createObject(data.Name, "server", backend, "", data, nil, transactionID, version)
	if err != nil {
		return err
	}

	if c.Cache.Enabled() {
		c.Cache.Servers.Set(data.Name, backend, transactionID, data)
	}
	return nil
}

// EditServer edits a server in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditServer(name string, backend string, data *models.Server, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	ondiskSrv, err := c.GetServer(name, backend, transactionID)
	if err != nil {
		return err
	}

	err = c.editObject(name, "server", backend, "", data, ondiskSrv.Data, nil, transactionID, version)
	if err != nil {
		return err
	}

	if c.Cache.Enabled() {
		c.Cache.Servers.Set(name, backend, transactionID, data)
	}
	return nil
}

func (c *Client) parseServers(response string) models.Servers {
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
