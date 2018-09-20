package configuration

import (
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

func (self *LBCTLConfigurationClient) GetBackends() (*models.GetBackendsOKBody, error) {
	backendsString, err := self.executeLBCTL("l7-farm-dump", "")
	if err != nil {
		return nil, err
	}

	backends := self.parseBackends(backendsString)

	v, err := self.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetBackendsOKBody{Version: v, Data: backends}, nil
}

func (self *LBCTLConfigurationClient) GetBackend(name string) (*models.GetBackendOKBody, error) {
	backendStr, err := self.executeLBCTL("l7-farm-show", "", name)
	if err != nil {
		return nil, err
	}
	backend := &models.Backend{Name: name}

	self.parseObject(backendStr, backend)

	v, err := self.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetBackendOKBody{Version: v, Data: backend}, nil
}

func (self *LBCTLConfigurationClient) DeleteBackend(name string, transactionID string, version int64) error {
	return self.deleteObject(name, "farm", "", "", transactionID, version)
}

func (self *LBCTLConfigurationClient) CreateBackend(data *models.Backend, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return validationErr
	}
	return self.createObject(data.Name, "farm", "", "", data, nil, transactionID, version)
}

func (self *LBCTLConfigurationClient) EditBackend(name string, data *models.Backend, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return validationErr
	}
	ondiskBck, err := self.GetBackend(name)
	if err != nil {
		return err
	}
	return self.editObject(name, "farm", "", "", data, ondiskBck.Data, nil, transactionID, version)
}

func (self *LBCTLConfigurationClient) parseBackends(response string) models.Backends {
	backends := make(models.Backends, 0, 1)
	for _, backendStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(backendStr) == "" {
			continue
		}
		name := strings.TrimSpace(backendStr[strings.Index(backendStr, ".farm ")+6 : strings.Index(backendStr, "\n")])

		backendObj := &models.Backend{Name: name}
		self.parseObject(backendStr, backendObj)
		backends = append(backends, backendObj)
	}
	return backends
}
