// This file is generated, don't edit manually, see README.md for details.

package clientnative

import (
	parser "github.com/haproxytech/config-parser/v3"

	"github.com/haproxytech/client-native/v2/configuration"
	"github.com/haproxytech/client-native/v2/models"
)

// IConfigurationClient ...
type IConfigurationClient interface {
	// GetACLs returns configuration version and an array of
	// configured ACL lines in the specified parent. Returns error on fail.
	GetACLs(parentType, parentName string, transactionID string) (int64, models.Acls, error)
	// GetACL returns configuration version and a requested ACL line
	// in the specified parent. Returns error on fail or if ACL line does not exist.
	GetACL(id int64, parentType, parentName string, transactionID string) (int64, *models.ACL, error)
	// DeleteACL deletes a ACL line in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteACL(id int64, parentType string, parentName string, transactionID string, version int64) error
	// CreateACL creates a ACL line in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateACL(parentType string, parentName string, data *models.ACL, transactionID string, version int64) error
	// EditACL edits a ACL line in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditACL(id int64, parentType string, parentName string, data *models.ACL, transactionID string, version int64) error
	// GetBackends returns configuration version and an array of
	// configured backends. Returns error on fail.
	GetBackends(transactionID string) (int64, models.Backends, error)
	// GetBackend returns configuration version and a requested backend.
	// Returns error on fail or if backend does not exist.
	GetBackend(name string, transactionID string) (int64, *models.Backend, error)
	// DeleteBackend deletes a backend in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteBackend(name string, transactionID string, version int64) error
	// CreateBackend creates a backend in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateBackend(data *models.Backend, transactionID string, version int64) error
	// EditBackend edits a backend in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditBackend(name string, data *models.Backend, transactionID string, version int64) error
	// GetBackendSwitchingRules returns configuration version and an array of
	// configured backend switching rules in the specified frontend. Returns error on fail.
	GetBackendSwitchingRules(frontend string, transactionID string) (int64, models.BackendSwitchingRules, error)
	// GetBackendSwitchingRule returns configuration version and a requested backend switching rule
	// in the specified frontend. Returns error on fail or if backend switching rule does not exist.
	GetBackendSwitchingRule(id int64, frontend string, transactionID string) (int64, *models.BackendSwitchingRule, error)
	// DeleteBackendSwitchingRule deletes a backend switching rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteBackendSwitchingRule(id int64, frontend string, transactionID string, version int64) error
	// CreateBackendSwitchingRule creates a backend switching rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateBackendSwitchingRule(frontend string, data *models.BackendSwitchingRule, transactionID string, version int64) error
	// EditBackendSwitchingRule edits a backend switching rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditBackendSwitchingRule(id int64, frontend string, data *models.BackendSwitchingRule, transactionID string, version int64) error
	// GetBinds returns configuration version and an array of
	// configured binds in the specified frontend. Returns error on fail.
	GetBinds(frontend string, transactionID string) (int64, models.Binds, error)
	// GetBind returns configuration version and a requested bind
	// in the specified frontend. Returns error on fail or if bind does not exist.
	GetBind(name string, frontend string, transactionID string) (int64, *models.Bind, error)
	// DeleteBind deletes a bind in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteBind(name string, frontend string, transactionID string, version int64) error
	// CreateBind creates a bind in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateBind(frontend string, data *models.Bind, transactionID string, version int64) error
	// EditBind edits a bind in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditBind(name string, frontend string, data *models.Bind, transactionID string, version int64) error
	// Init initializes a Client
	Init(options configuration.ClientParams) error
	// HasParser checks whether transaction exists in parser
	HasParser(transactionID string) bool
	// GetParserTransactions returns parser transactions
	GetParserTransactions() models.Transactions
	// GetParser returns a parser for given transactionID, if transactionID is "", it returns "master" parser
	GetParser(transactionID string) (*parser.Parser, error)
	// AddParser adds parser to parser map
	AddParser(transactionID string) error
	// DeleteParser deletes parser from parsers map
	DeleteParser(transactionID string) error
	// CommitParser commits transaction parser, deletes it from parsers map, and replaces master Parser
	CommitParser(transactionID string) error
	// InitTransactionParsers checks transactions and initializes parsers map with transactions in_progress
	InitTransactionParsers() error
	// GetVersion returns configuration file version
	GetVersion(transactionID string) (int64, error)
	IncrementVersion() error
	IncrementTransactionVersion(transactionID string) error
	LoadData(filename string) error
	Save(transactionFile, transactionID string) error
	GetFailedParserTransactionVersion(transactionID string) (int64, error)
	// GetDefaultsConfiguration returns configuration version and a
	// struct representing Defaults configuration
	GetDefaultsConfiguration(transactionID string) (int64, *models.Defaults, error)
	// PushDefaultsConfiguration pushes a Defaults config struct to global
	// config file
	PushDefaultsConfiguration(data *models.Defaults, transactionID string, version int64) error
	// GetFilters returns configuration version and an array of
	// configured filters in the specified parent. Returns error on fail.
	GetFilters(parentType, parentName string, transactionID string) (int64, models.Filters, error)
	// GetFilter returns configuration version and a requested filter
	// in the specified parent. Returns error on fail or if filter does not exist.
	GetFilter(id int64, parentType, parentName string, transactionID string) (int64, *models.Filter, error)
	// DeleteFilter deletes a filter in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteFilter(id int64, parentType string, parentName string, transactionID string, version int64) error
	// CreateFilter creates a filter in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateFilter(parentType string, parentName string, data *models.Filter, transactionID string, version int64) error
	// EditFilter edits a filter in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditFilter(id int64, parentType string, parentName string, data *models.Filter, transactionID string, version int64) error
	// GetFrontends returns configuration version and an array of
	// configured frontends. Returns error on fail.
	GetFrontends(transactionID string) (int64, models.Frontends, error)
	// GetFrontend returns configuration version and a requested frontend.
	// Returns error on fail or if frontend does not exist.
	GetFrontend(name string, transactionID string) (int64, *models.Frontend, error)
	// DeleteFrontend deletes a frontend in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteFrontend(name string, transactionID string, version int64) error
	// EditFrontend edits a frontend in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditFrontend(name string, data *models.Frontend, transactionID string, version int64) error
	// CreateFrontend creates a frontend in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateFrontend(data *models.Frontend, transactionID string, version int64) error
	// GetGlobalConfiguration returns configuration version and a
	// struct representing Global configuration
	GetGlobalConfiguration(transactionID string) (int64, *models.Global, error)
	// PushGlobalConfiguration pushes a Global config struct to global
	// config file
	PushGlobalConfiguration(data *models.Global, transactionID string, version int64) error
	// GetHTTPRequestRules returns configuration version and an array of
	// configured http request rules in the specified parent. Returns error on fail.
	GetHTTPRequestRules(parentType, parentName string, transactionID string) (int64, models.HTTPRequestRules, error)
	// GetHTTPRequestRule returns configuration version and a requested http request rule
	// in the specified parent. Returns error on fail or if http request rule does not exist.
	GetHTTPRequestRule(id int64, parentType, parentName string, transactionID string) (int64, *models.HTTPRequestRule, error)
	// DeleteHTTPRequestRule deletes a http request rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteHTTPRequestRule(id int64, parentType string, parentName string, transactionID string, version int64) error
	// CreateHTTPRequestRule creates a http request rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateHTTPRequestRule(parentType string, parentName string, data *models.HTTPRequestRule, transactionID string, version int64) error
	// EditHTTPRequestRule edits a http request rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditHTTPRequestRule(id int64, parentType string, parentName string, data *models.HTTPRequestRule, transactionID string, version int64) error
	// GetHTTPResponseRules returns configuration version and an array of
	// configured http response rules in the specified parent. Returns error on fail.
	GetHTTPResponseRules(parentType, parentName string, transactionID string) (int64, models.HTTPResponseRules, error)
	// GetHTTPResponseRule returns configuration version and a responseed http response rule
	// in the specified parent. Returns error on fail or if http response rule does not exist.
	GetHTTPResponseRule(id int64, parentType, parentName string, transactionID string) (int64, *models.HTTPResponseRule, error)
	// DeleteHTTPResponseRule deletes a http response rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteHTTPResponseRule(id int64, parentType string, parentName string, transactionID string, version int64) error
	// CreateHTTPResponseRule creates a http response rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateHTTPResponseRule(parentType string, parentName string, data *models.HTTPResponseRule, transactionID string, version int64) error
	// EditHTTPResponseRule edits a http response rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditHTTPResponseRule(id int64, parentType string, parentName string, data *models.HTTPResponseRule, transactionID string, version int64) error
	// GetLogTargets returns configuration version and an array of
	// configured log targets in the specified parent. Returns error on fail.
	GetLogTargets(parentType, parentName string, transactionID string) (int64, models.LogTargets, error)
	// GetLogTarget returns configuration version and a requested log target
	// in the specified parent. Returns error on fail or if log target does not exist.
	GetLogTarget(id int64, parentType, parentName string, transactionID string) (int64, *models.LogTarget, error)
	// DeleteLogTarget deletes a log target in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteLogTarget(id int64, parentType string, parentName string, transactionID string, version int64) error
	// CreateLogTarget creates a log target in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateLogTarget(parentType string, parentName string, data *models.LogTarget, transactionID string, version int64) error
	// EditLogTarget edits a log target in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditLogTarget(id int64, parentType string, parentName string, data *models.LogTarget, transactionID string, version int64) error
	// GetNameservers returns configuration version and an array of
	// configured namservers in the specified resolvers section. Returns error on fail.
	GetNameservers(resolverSection string, transactionID string) (int64, models.Nameservers, error)
	// GetNameserver returns configuration version and a requested nameserver
	// in the specified resolvers section. Returns error on fail or if nameserver does not exist.
	GetNameserver(name string, resolverSection string, transactionID string) (int64, *models.Nameserver, error)
	// DeleteNameserver deletes an nameserver in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteNameserver(name string, resolverSection string, transactionID string, version int64) error
	// CreateNameserver creates a nameserver in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateNameserver(resolverSection string, data *models.Nameserver, transactionID string, version int64) error
	// EditNameserver edits a nameserver in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditNameserver(name string, resolverSection string, data *models.Nameserver, transactionID string, version int64) error
	// GetPeerEntries returns configuration version and an array of
	// configured binds in the specified peers section. Returns error on fail.
	GetPeerEntries(peerSection string, transactionID string) (int64, models.PeerEntries, error)
	// GetPeerEntry returns configuration version and a requested peer entry
	// in the specified peer section. Returns error on fail or if bind does not exist.
	GetPeerEntry(name string, peerSection string, transactionID string) (int64, *models.PeerEntry, error)
	// DeletePeerEntry deletes an peer entry in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeletePeerEntry(name string, peerSection string, transactionID string, version int64) error
	// CreatePeerEntry creates a peer entry in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreatePeerEntry(peerSection string, data *models.PeerEntry, transactionID string, version int64) error
	// EditPeerEntry edits a peer entry in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditPeerEntry(name string, peerSection string, data *models.PeerEntry, transactionID string, version int64) error
	// GetPeerSections returns configuration version and an array of
	// configured peer sections. Returns error on fail.
	GetPeerSections(transactionID string) (int64, models.PeerSections, error)
	// GetPeerSection returns configuration version and a requested peer section.
	// Returns error on fail or if peer section does not exist.
	GetPeerSection(name string, transactionID string) (int64, *models.PeerSection, error)
	// DeletePeerSection deletes a peerSection in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeletePeerSection(name string, transactionID string, version int64) error
	// CreatePeerSection creates a peerSection in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreatePeerSection(data *models.PeerSection, transactionID string, version int64) error
	// GetRawConfiguration returns configuration version and a
	// string containing raw config file
	GetRawConfiguration(transactionID string, version int64) (int64, string, error)
	// PostRawConfiguration pushes given string to the config file if the version
	// matches
	PostRawConfiguration(config *string, version int64, skipVersionCheck bool, onlyValidate ...bool) error
	// GetResolvers returns configuration version and an array of
	// configured resolvers. Returns error on fail.
	GetResolvers(transactionID string) (int64, models.Resolvers, error)
	// GetResolver returns configuration version and a requested resolver.
	// Returns error on fail or if resolver does not exist.
	GetResolver(name string, transactionID string) (int64, *models.Resolver, error)
	// DeleteResolver deletes a resolver in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteResolver(name string, transactionID string, version int64) error
	// EditResolver edits a resolver in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditResolver(name string, data *models.Resolver, transactionID string, version int64) error
	// CreateResolver creates a resolver in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateResolver(data *models.Resolver, transactionID string, version int64) error
	// GetServers returns configuration version and an array of
	// configured servers in the specified backend. Returns error on fail.
	GetServers(backend string, transactionID string) (int64, models.Servers, error)
	// GetServer returns configuration version and a requested server
	// in the specified backend. Returns error on fail or if server does not exist.
	GetServer(name string, backend string, transactionID string) (int64, *models.Server, error)
	// DeleteServer deletes a server in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteServer(name string, backend string, transactionID string, version int64) error
	// CreateServer creates a server in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateServer(backend string, data *models.Server, transactionID string, version int64) error
	// EditServer edits a server in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditServer(name string, backend string, data *models.Server, transactionID string, version int64) error
	// GetServerSwitchingRules returns configuration version and an array of
	// configured server switching rules in the specified backend. Returns error on fail.
	GetServerSwitchingRules(backend string, transactionID string) (int64, models.ServerSwitchingRules, error)
	// GetServerSwitchingRule returns configuration version and a requested server switching rule
	// in the specified backend. Returns error on fail or if server switching rule does not exist.
	GetServerSwitchingRule(id int64, backend string, transactionID string) (int64, *models.ServerSwitchingRule, error)
	// DeleteServerSwitchingRule deletes a server switching rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteServerSwitchingRule(id int64, backend string, transactionID string, version int64) error
	// CreateServerSwitchingRule creates a server switching rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateServerSwitchingRule(backend string, data *models.ServerSwitchingRule, transactionID string, version int64) error
	// EditServerSwitchingRule edits a server switching rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditServerSwitchingRule(id int64, backend string, data *models.ServerSwitchingRule, transactionID string, version int64) error
	// NewService creates and returns a new Service instance.
	// name indicates the name of the service and only one Service instance with the given name can be created.
	NewService(name string, scaling configuration.ScalingParams) (*configuration.Service, error)
	// DeleteService removes the Service instance specified by name from the client.
	DeleteService(name string)
	// GetSites returns configuration version and an array of
	// configured sites. Returns error on fail.
	GetSites(transactionID string) (int64, models.Sites, error)
	// GetSite returns configuration version and a requested site.
	// Returns error on fail or if backend does not exist.
	GetSite(name string, transactionID string) (int64, *models.Site, error)
	// CreateSite creates a site in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateSite(data *models.Site, transactionID string, version int64) error
	// EditSite edits a site in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditSite(name string, data *models.Site, transactionID string, version int64) error
	// DeleteSite deletes a site in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteSite(name string, transactionID string, version int64) error
	// GetStickRules returns configuration version and an array of
	// configured stick rules in the specified backend. Returns error on fail.
	GetStickRules(backend string, transactionID string) (int64, models.StickRules, error)
	// GetStickRule returns configuration version and a requested stick rule
	// in the specified backend. Returns error on fail or if stick rule does not exist.
	GetStickRule(id int64, backend string, transactionID string) (int64, *models.StickRule, error)
	// DeleteStickRule deletes a stick rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteStickRule(id int64, backend string, transactionID string, version int64) error
	// CreateStickRule creates a stick rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateStickRule(backend string, data *models.StickRule, transactionID string, version int64) error
	// EditStickRule edits a stick rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditStickRule(id int64, backend string, data *models.StickRule, transactionID string, version int64) error
	// GetTCPRequestRules returns configuration version and an array of
	// configured TCP request rules in the specified parent. Returns error on fail.
	GetTCPRequestRules(parentType, parentName string, transactionID string) (int64, models.TCPRequestRules, error)
	// GetTCPRequestRule returns configuration version and a requested tcp request rule
	// in the specified parent. Returns error on fail or if http request rule does not exist.
	GetTCPRequestRule(id int64, parentType, parentName string, transactionID string) (int64, *models.TCPRequestRule, error)
	// DeleteTCPRequestRule deletes a tcp request rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteTCPRequestRule(id int64, parentType string, parentName string, transactionID string, version int64) error
	// CreateTCPRequestRule creates a tcp request rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateTCPRequestRule(parentType string, parentName string, data *models.TCPRequestRule, transactionID string, version int64) error
	// EditTCPRequestRule edits a tcp request rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditTCPRequestRule(id int64, parentType string, parentName string, data *models.TCPRequestRule, transactionID string, version int64) error
	// GetTCPResponseRules returns configuration version and an array of
	// configured tcp response rules in the specified backend. Returns error on fail.
	GetTCPResponseRules(backend string, transactionID string) (int64, models.TCPResponseRules, error)
	// GetTCPResponseRule returns configuration version and a requested tcp response rule
	// in the specified backend. Returns error on fail or if tcp response rule does not exist.
	GetTCPResponseRule(id int64, backend string, transactionID string) (int64, *models.TCPResponseRule, error)
	// DeleteTCPResponseRule deletes a tcp response rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	DeleteTCPResponseRule(id int64, backend string, transactionID string, version int64) error
	// CreateTCPResponseRule creates a tcp response rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	CreateTCPResponseRule(backend string, data *models.TCPResponseRule, transactionID string, version int64) error
	// EditTCPResponseRule edits a tcp response rule in configuration. One of version or transactionID is
	// mandatory. Returns error on fail, nil on success.
	EditTCPResponseRule(id int64, backend string, data *models.TCPResponseRule, transactionID string, version int64) error
	// GetConfigurationVersion returns configuration version
	GetConfigurationVersion(transactionID string) (int64, error)
}
