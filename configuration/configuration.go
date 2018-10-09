package configuration

import (
	"github.com/haproxytech/models"
)

// ClientParams is just a placeholder for all client options
type ClientParams struct {
	configurationFile string
}

// NewConfigurationClientParams creates a new configuration client.
func NewConfigurationClientParams(configurationFile string) *ClientParams {
	return &ClientParams{configurationFile: configurationFile}
}

// ConfigurationFile changes the configuration file on the client
func (c *ClientParams) ConfigurationFile() string {
	return c.configurationFile
}

// Client interface is the interface used for managing HAProxy configuration
// file.
type Client interface {
	//transaction methods
	GetTransactions(status string) (*models.Transactions, error)
	GetTransaction(id string) (*models.Transaction, error)
	StartTransaction(version int64) (*models.Transaction, error)
	CommitTransaction(id string) error
	DeleteTransaction(id string) error
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
	//backend switching rule methods
	GetBackendSwitchingRules(frontend string) (*models.GetBackendSwitchingRulesOKBody, error)
	GetBackendSwitchingRule(id int64, frontend string) (*models.GetBackendSwitchingRuleOKBody, error)
	DeleteBackendSwitchingRule(id int64, frontend string, transactionID string, version int64) error
	CreateBackendSwitchingRule(frontend string, data *models.BackendSwitchingRule, transactionID string, version int64) error
	EditBackendSwitchingRule(id int64, frontend string, data *models.BackendSwitchingRule, transactionID string, version int64) error
	//server switching rule methods
	GetServerSwitchingRules(backend string) (*models.GetServerSwitchingRulesOKBody, error)
	GetServerSwitchingRule(id int64, backend string) (*models.GetServerSwitchingRuleOKBody, error)
	DeleteServerSwitchingRule(id int64, backend string, transactionID string, version int64) error
	CreateServerSwitchingRule(backend string, data *models.ServerSwitchingRule, transactionID string, version int64) error
	EditServerSwitchingRule(id int64, backend string, data *models.ServerSwitchingRule, transactionID string, version int64) error
	//filter methods
	GetFilters(parentType, parentName string) (*models.GetFiltersOKBody, error)
	GetFilter(id int64, parentType, parentName string) (*models.GetFilterOKBody, error)
	DeleteFilter(id int64, parentType, parentName string, transactionID string, version int64) error
	CreateFilter(parentType, parentName string, data *models.Filter, transactionID string, version int64) error
	EditFilter(id int64, parentType, parentName string, data *models.Filter, transactionID string, version int64) error
	//HTTP Request Rules methods
	GetHTTPRequestRules(parentType, parentName string) (*models.GetHTTPRequestRulesOKBody, error)
	GetHTTPRequestRule(id int64, parentType, parentName string) (*models.GetHTTPRequestRuleOKBody, error)
	DeleteHTTPRequestRule(id int64, parentType, parentName string, transactionID string, version int64) error
	CreateHTTPRequestRule(parentType, parentName string, data *models.HTTPRequestRule, transactionID string, version int64) error
	EditHTTPRequestRule(id int64, parentType, parentName string, data *models.HTTPRequestRule, transactionID string, version int64) error
	//HTTP Response Rules methods
	GetHTTPResponseRules(parentType, parentName string) (*models.GetHTTPResponseRulesOKBody, error)
	GetHTTPResponseRule(id int64, parentType, parentName string) (*models.GetHTTPResponseRuleOKBody, error)
	DeleteHTTPResponseRule(id int64, parentType, parentName string, transactionID string, version int64) error
	CreateHTTPResponseRule(parentType, parentName string, data *models.HTTPResponseRule, transactionID string, version int64) error
	EditHTTPResponseRule(id int64, parentType, parentName string, data *models.HTTPResponseRule, transactionID string, version int64) error
}
