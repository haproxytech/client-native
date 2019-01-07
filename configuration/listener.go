package configuration

import (
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetListeners returns a struct with configuration version and an array of
// configured listeners in the specified frontend. Returns error on fail.
func (c *Client) GetListeners(frontend string, transactionID string) (*models.GetListenersOKBody, error) {
	if c.Cache.Enabled() {
		listeners, found := c.Cache.Listeners.Get(frontend, transactionID)
		if found {
			return &models.GetListenersOKBody{Version: c.Cache.Version.Get(), Data: listeners}, nil
		}
	}
	listenersString, err := c.executeLBCTL("l7-listener-dump", transactionID, frontend)
	if err != nil {
		return nil, err
	}

	listeners := c.parseListeners(listenersString)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.Listeners.SetAll(frontend, transactionID, listeners)
	}
	return &models.GetListenersOKBody{Version: v, Data: listeners}, nil
}

// GetListener returns a struct with configuration version and a requested listener
// in the specified frontend. Returns error on fail or if listener does not exist.
func (c *Client) GetListener(name string, frontend string, transactionID string) (*models.GetListenerOKBody, error) {
	if c.Cache.Enabled() {
		listener, found := c.Cache.Listeners.GetOne(name, frontend, transactionID)
		if found {
			return &models.GetListenerOKBody{Version: c.Cache.Version.Get(), Data: listener}, nil
		}
	}
	listenerStr, err := c.executeLBCTL("l7-listener-show", transactionID, frontend, name)
	if err != nil {
		return nil, err
	}
	listener := &models.Listener{Name: name}

	c.parseObject(listenerStr, listener)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	if c.Cache.Enabled() {
		c.Cache.Listeners.Set(name, frontend, transactionID, listener)
	}
	return &models.GetListenerOKBody{Version: v, Data: listener}, nil
}

// DeleteListener deletes a listener in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteListener(name string, frontend string, transactionID string, version int64) error {
	err := c.deleteObject(name, "listener", frontend, "", transactionID, version)
	if err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.Listeners.Delete(name, frontend, transactionID)
	}
	return nil
}

// CreateListener creates a listener in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateListener(frontend string, data *models.Listener, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	err := c.createObject(data.Name, "listener", frontend, "", data, nil, transactionID, version)
	if err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.Listeners.Set(data.Name, frontend, transactionID, data)
	}
	return nil
}

// EditListener edits a listener in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditListener(name string, frontend string, data *models.Listener, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	ondiskLst, err := c.GetListener(name, frontend, transactionID)
	if err != nil {
		return err
	}

	err = c.editObject(name, "listener", frontend, "", data, ondiskLst.Data, nil, transactionID, version)
	if err != nil {
		return err
	}

	if c.Cache.Enabled() {
		c.Cache.Listeners.Set(name, frontend, transactionID, data)
	}
	return nil
}

func (c *Client) parseListeners(response string) models.Listeners {
	listeners := make(models.Listeners, 0, 1)
	for _, listenerStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(listenerStr) == "" {
			continue
		}
		name, _ := splitHeaderLine(listenerStr)

		listenerObj := &models.Listener{Name: name}
		c.parseObject(listenerStr, listenerObj)
		listeners = append(listeners, listenerObj)
	}
	return listeners
}
