package configuration

import (
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

// GetListeners returns a struct with configuration version and an array of
// configured listeners in the specified frontend. Returns error on fail.
func (c *LBCTLClient) GetListeners(frontend string, transactionID string) (*models.GetListenersOKBody, error) {
	listenersString, err := c.executeLBCTL("l7-listener-dump", transactionID, frontend)
	if err != nil {
		return nil, err
	}

	listeners := c.parseListeners(listenersString)

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetListenersOKBody{Version: v, Data: listeners}, nil
}

// GetListener returns a struct with configuration version and a requested listener
// in the specified frontend. Returns error on fail or if listener does not exist.
func (c *LBCTLClient) GetListener(name string, frontend string, transactionID string) (*models.GetListenerOKBody, error) {
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

	return &models.GetListenerOKBody{Version: v, Data: listener}, nil
}

// DeleteListener deletes a listener in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) DeleteListener(name string, frontend string, transactionID string, version int64) error {
	return c.deleteObject(name, "listener", frontend, "", transactionID, version)
}

// CreateListener creates a listener in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) CreateListener(frontend string, data *models.Listener, transactionID string, version int64) error {
	if c.UseValidation() {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	return c.createObject(data.Name, "listener", frontend, "", data, nil, transactionID, version)
}

// EditListener edits a listener in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *LBCTLClient) EditListener(name string, frontend string, data *models.Listener, transactionID string, version int64) error {
	if c.UseValidation() {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	ondiskLst, err := c.GetListener(name, frontend, transactionID)
	if err != nil {
		return err
	}

	return c.editObject(name, "listener", frontend, "", data, ondiskLst.Data, nil, transactionID, version)
}

func (c *LBCTLClient) parseListeners(response string) models.Listeners {
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
