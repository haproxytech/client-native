package configuration

import (
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/haproxytech/models"
)

func (self *LBCTLConfigurationClient) GetServers(backend string) (*models.GetServersOKBody, error) {
	serversString, err := self.executeLBCTL("l7-server-dump", "", backend)
	if err != nil {
		return nil, err
	}

	servers := self.parseServers(serversString)

	v, err := self.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetServersOKBody{Version: v, Data: servers}, nil
}

func (self *LBCTLConfigurationClient) GetServer(name string, backend string) (*models.GetServerOKBody, error) {
	serverStr, err := self.executeLBCTL("l7-server-show", "", backend, name)
	if err != nil {
		return nil, err
	}
	server := &models.Server{Name: name}

	self.parseObject(serverStr, server)

	v, err := self.GetVersion()
	if err != nil {
		return nil, err
	}

	return &models.GetServerOKBody{Version: v, Data: server}, nil
}

func (self *LBCTLConfigurationClient) DeleteServer(name string, backend string, transactionID string, version int64) error {
	return self.deleteObject(name, "server", backend, "", transactionID, version)
}

func (self *LBCTLConfigurationClient) CreateServer(backend string, data *models.Server, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return validationErr
	}
	return self.createObject(data.Name, "server", backend, "", data, nil, transactionID, version)
}

func (self *LBCTLConfigurationClient) EditServer(name string, backend string, data *models.Server, transactionID string, version int64) error {
	validationErr := data.Validate(strfmt.Default)
	if validationErr != nil {
		return validationErr
	}
	ondiskSrv, err := self.GetServer(name, backend)
	if err != nil {
		return err
	}

	return self.editObject(name, "server", backend, "", data, ondiskSrv.Data, nil, transactionID, version)
}

func (self *LBCTLConfigurationClient) parseServers(response string) models.Servers {
	servers := make(models.Servers, 0, 1)
	for _, serverStr := range strings.Split(response, "\n\n") {
		if strings.TrimSpace(serverStr) == "" {
			continue
		}
		name, _ := splitHeaderLine(serverStr)

		serverObj := &models.Server{Name: name}
		self.parseObject(serverStr, serverObj)
		servers = append(servers, serverObj)
	}
	return servers
}
