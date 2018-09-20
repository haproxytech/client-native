package configuration

import (
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

func (self *LBCTLConfigurationClient) GetFrontends() (*models.GetFrontendsOKBody, error) {
	frontendsStr, err := self.executeLBCTL("l7-service-dump", "")
	if err != nil {
		return nil, err
	}
	frontends := self.parseFrontends(frontendsStr)

	v, err := self.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetFrontendsOKBody{Version: v, Data: frontends}, nil
}

func (self *LBCTLConfigurationClient) GetFrontend(name string) (*models.GetFrontendOKBody, error) {
	frontendStr, err := self.executeLBCTL("l7-service-show", "", name)
	if err != nil {
		return nil, err
	}
	frontend := &models.Frontend{Name: name}

	self.parseObject(frontendStr, frontend)

	v, err := self.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetFrontendOKBody{Version: v, Data: frontend}, nil
}

func (self *LBCTLConfigurationClient) DeleteFrontend(name string, transactionID string, version int64) error {
	return self.deleteObject(name, "service", "", "", transactionID, version)
}

func (self *LBCTLConfigurationClient) EditFrontend(name string, data *models.Frontend, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return validationErr
	}
	ondiskFrontend, err := self.GetFrontend(name)
	if err != nil {
		return err
	}
	return self.editObject(name, "service", "", "", data, ondiskFrontend.Data, nil, transactionID, version)
}

func (self *LBCTLConfigurationClient) CreateFrontend(data *models.Frontend, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return validationErr
	}
	return self.createObject(data.Name, "service", "", "", data, nil, transactionID, version)
}

func (self *LBCTLConfigurationClient) parseFrontends(response string) models.Frontends {
	frontends := make(models.Frontends, 0, 1)
	for _, frontendStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(frontendStr) == "" {
			continue
		}
		name := strings.TrimSpace(frontendStr[strings.Index(frontendStr, ".service ")+9 : strings.Index(frontendStr, "\n")])

		frontendObj := &models.Frontend{Name: name}
		self.parseObject(frontendStr, frontendObj)
		frontends = append(frontends, frontendObj)
	}
	return frontends
}
