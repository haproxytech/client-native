package configuration

import (
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

func (self *LBCTLConfigurationClient) GetListeners(frontend string) (*models.GetListenersOKBody, error) {
	listenersString, err := self.executeLBCTL("l7-listener-dump", "", frontend)
	if err != nil {
		return nil, err
	}

	listeners := self.parseListeners(listenersString)

	v, err := self.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetListenersOKBody{Version: v, Data: listeners}, nil
}

func (self *LBCTLConfigurationClient) GetListener(name string, frontend string) (*models.GetListenerOKBody, error) {
	listenerStr, err := self.executeLBCTL("l7-listener-show", "", frontend, name)
	if err != nil {
		return nil, err
	}
	listener := &models.Listener{Name: name}

	self.parseObject(listenerStr, listener)

	v, err := self.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetListenerOKBody{Version: v, Data: listener}, nil
}

func (self *LBCTLConfigurationClient) DeleteListener(name string, frontend string, transactionID string, version int64) error {
	return self.deleteObject(name, "listener", frontend, "", transactionID, version)
}

func (self *LBCTLConfigurationClient) CreateListener(frontend string, data *models.Listener, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return validationErr
	}
	return self.createObject(data.Name, "listener", frontend, "", data, nil, transactionID, version)
}

func (self *LBCTLConfigurationClient) EditListener(name string, frontend string, data *models.Listener, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return validationErr
	}
	ondiskLst, err := self.GetListener(name, frontend)
	if err != nil {
		return err
	}

	return self.editObject(name, "listener", frontend, "", data, ondiskLst.Data, nil, transactionID, version)
}

func (self *LBCTLConfigurationClient) parseListeners(response string) models.Listeners {
	listeners := make(models.Listeners, 0, 1)
	for _, listenerStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(listenerStr) == "" {
			continue
		}
		name, _ := splitHeaderLine(listenerStr)

		listenerObj := &models.Listener{Name: name}
		self.parseObject(listenerStr, listenerObj)
		listeners = append(listeners, listenerObj)
	}
	return listeners
}
