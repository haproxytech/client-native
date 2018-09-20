package configuration

import (
	"github.com/haproxytech/models"
)

type ConfigurationClientParams struct {
	configurationFile string
}

// NewConfigurationClientParams creates a new configuration client.
func NewConfigurationClientParams(configurationFile string) *ConfigurationClientParams {
	return &ConfigurationClientParams{configurationFile: configurationFile}
}

// ConfigurationFile changes the configuration file on the client
func (self *ConfigurationClientParams) ConfigurationFile() string {
	return self.configurationFile
}

type Client interface {
	//transaction methods
	GetTransactions(status string) (*models.Transactions, error)
	GetTransaction(id string) (*models.Transaction, error)
	StartTransaction(version int64) (*models.Transaction, error)
	CommitTransaction(id string) error
	//version method
	GetVersion() (int64, error)
	//site methods
	GetSites() (*models.GetSitesOKBody, error)
	GetSite(name string) (*models.GetSiteOKBody, error)
	DeleteSite(name string, transactionID string, version int64) error
	CreateSite(data *models.Site, transactionID string, version int64) error
	EditSite(name string, data *models.Site, transactionID string, version int64) error
	//frontend methods
	GetFrontends() (*models.GetFrontendsOKBody, error)
	GetFrontend(name string) (*models.GetFrontendOKBody, error)
	DeleteFrontend(name string, transactionID string, version int64) error
	CreateFrontend(data *models.Frontend, transactionID string, version int64) error
	EditFrontend(name string, data *models.Frontend, transactionID string, version int64) error
	//backend methods
	GetBackends() (*models.GetBackendsOKBody, error)
	GetBackend(name string) (*models.GetBackendOKBody, error)
	DeleteBackend(name string, transactionID string, version int64) error
	CreateBackend(data *models.Backend, transactionID string, version int64) error
	EditBackend(name string, data *models.Backend, transactionID string, version int64) error
	//server methods
	GetServers(backend string) (*models.GetServersOKBody, error)
	GetServer(name string, backend string) (*models.GetServerOKBody, error)
	DeleteServer(name string, backend string, transactionID string, version int64) error
	CreateServer(backend string, data *models.Server, transactionID string, version int64) error
	EditServer(name string, backend string, data *models.Server, transactionID string, version int64) error
	//listener methods
	GetListeners(frontend string) (*models.GetListenersOKBody, error)
	GetListener(name string, frontend string) (*models.GetListenerOKBody, error)
	DeleteListener(name string, frontend string, transactionID string, version int64) error
	CreateListener(frontend string, data *models.Listener, transactionID string, version int64) error
	EditListener(name string, frontend string, data *models.Listener, transactionID string, version int64) error
}
