// Copyright 2022 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package runtime

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/client-native/v6/runtime/options"
)

type Maps interface {
	GetMapsDir() (string, error)
	// GetMapsPath returns runtime map file path or map id
	GetMapsPath(name string) (string, error)
	// ShowMaps returns structured unique map files
	ShowMaps() (models.Maps, error)
	// CreateMap creates a new map file with its entries
	CreateMap(file io.Reader, header multipart.FileHeader) (*models.Map, error)
	// GetMap returns one structured runtime map file
	GetMap(name string) (*models.Map, error)
	// ClearMap removes all map entries from the map file. If forceDelete is true, deletes file from disk
	ClearMap(name string, forceDelete bool) error
	// ShowMapEntries list all map entries by map file name
	ShowMapEntries(name string) (models.MapEntries, error)
	// AddMapPayload adds multiple entries to the map file
	AddMapPayload(name, payload string) error
	// AddMapEntry adds an entry into the map file
	AddMapEntry(name, key, value string) error
	// GetMapEntry returns one map runtime setting
	GetMapEntry(name, id string) (*models.MapEntry, error)
	// SetMapEntry replace the value corresponding to each id in a map
	SetMapEntry(name, id, value string) error
	// DeleteMapEntry deletes all the map entries from the map by its id
	DeleteMapEntry(name, id string) error
	ParseMapEntries(output string) models.MapEntries
	// ParseMapEntriesFromFile reads entries from file
	ParseMapEntriesFromFile(inputFile io.Reader, hasID bool) models.MapEntries
	AddMapPayloadVersioned(name string, entries models.MapEntries) error
	// PrepareMap allocates a new map version
	PrepareMap(name string) (version string, err error)
	// CommitMap commits all changes made to a map version
	CommitMap(version, name string) error
}

type Servers interface {
	// AddServer adds a new server to a backend
	AddServer(backend, name, attributes string) error
	// DeleteServer removes a server from a backend
	DeleteServer(backend, name string) error
	// SetServerAddr set ip [port] for server
	SetServerAddr(backend, server string, ip string, port int) error
	// SetServerState set state for server
	SetServerState(backend, server string, state string) error
	// SetServerWeight set weight for server
	SetServerWeight(backend, server string, weight string) error
	// SetServerHealth set health for server
	SetServerHealth(backend, server string, health string) error
	// EnableServerHealth
	EnableServerHealth(backend, server string) error
	// EnableAgentCheck enable agent check for server
	EnableAgentCheck(backend, server string) error
	// DisableAgentCheck disable agent check for server
	DisableAgentCheck(backend, server string) error
	// EnableServer marks server as UP
	EnableServer(backend, server string) error
	// DisableServer marks server as DOWN for maintenance
	DisableServer(backend, server string) error
	// SetServerAgentAddr set agent-addr for server
	SetServerAgentAddr(backend, server string, addr string) error
	// SetServerAgentSend set agent-send for server
	SetServerAgentSend(backend, server string, send string) error
	// GetServerState returns server runtime state
	GetServersState(backend string) (models.RuntimeServers, error)
	// GetServerState returns server runtime state
	GetServerState(backend, server string) (*models.RuntimeServer, error)
	// SetServerCheckPort set health heck port for server
	SetServerCheckPort(backend, server string, port int) error
}

type ACLs interface {
	GetACLFiles() (files models.ACLFiles, err error)
	GetACLFile(id string) (files *models.ACLFile, err error)
	GetACLFilesEntries(id string) (files models.ACLFilesEntries, err error)
	GetACLFileEntry(id, value string) (fileEntry *models.ACLFileEntry, err error)
	AddACLFileEntry(id, value string) error
	DeleteACLFileEntry(id, value string) error
	AddACLAtomic(aclID string, entries models.ACLFilesEntries) error
}

type Tables interface {
	// SetTableEntry create or update a stick-table entry in the table.
	SetTableEntry(table, key string, dataType models.StickTableEntry) error
	// Show tables show tables from runtime API and return it structured
	ShowTables() (models.StickTables, error)
	// GetTableEntries returns all entries for specified table with filters and a key
	GetTableEntries(name string, filter []string, key string) (models.StickTableEntries, error)
	// Show table show tables {name} from runtime API and return it structured
	ShowTable(name string) (*models.StickTable, error)
}

type Frontend interface {
	// SetFrontendMaxConn set maxconn for frontend
	SetFrontendMaxConn(frontend string, maxconn int) error
}

type Info interface {
	// GetStats returns stats from the socket
	GetStats() models.NativeStats
	// GetInfo returns info from the socket
	GetInfo() (models.ProcessInfo, error)
	// GetVersion() returns running HAProxy version
	GetVersion() (HAProxyVersion, error)
}

type Manage interface {
	// Reloads HAProxy's configuration file. Similar to SIGUSR2. Returns the startup logs.
	Reload() (string, error)
	// Clears max counters
	ClearCounters() error
	// Clears counters
	ClearAllCounters() error
}

type Raw interface {
	// ExecuteRaw does not process response, just returns its value
	ExecuteRaw(command string) (string, error)
}

type SSL interface {
	ShowCAFiles() (models.SslCaFiles, error)
	GetCAFile(name string) (*models.SslCaFile, error)
	ShowCAFile(name string, index *int64) (*models.SslCertificate, error)
	NewCAFile(name string) error
	SetCAFile(name, payload string) error
	CommitCAFile(name string) error
	AbortCAFile(name string) error
	AddCAFileEntry(name, payload string) error
	DeleteCAFile(name string) error
	ShowCerts() (models.SslCertificates, error)
	GetCert(name string) (*models.SslCertificate, error)
	ShowCertificate(name string) (*models.SslCertificate, error)
	DumpCertificate(name string) (string, error)
	NewCertEntry(name string) error
	SetCertEntry(name, payload string) error
	CommitCertEntry(name string) error
	AbortCertEntry(name string) error
	DeleteCertEntry(name string) error
	ShowCrlFiles() (models.SslCrls, error)
	GetCrlFile(name string) (*models.SslCrl, error)
	ShowCrlFile(name string, index *int64) (*models.SslCrlEntries, error)
	NewCrlFile(name string) error
	SetCrlFile(name string, payload string) error
	CommitCrlFile(name string) error
	AbortCrlFile(name string) error
	DeleteCrlFile(name string) error
	ShowCrtLists() (models.SslCrtLists, error)
	GetCrtList(name string) (*models.SslCrtList, error)
	ShowCrtListEntries(name string) (models.SslCrtListEntries, error)
	AddCrtListEntry(crtList string, entry models.SslCrtListEntry) error
	DeleteCrtListEntry(crtList, filename string, lineNumber *int64) error
	ShowOcspResponses() ([]*models.SslCertificateID, error)
	ShowOcspResponse(idOrPath string, ofmt ...OcspResponseFmt) (*models.SslOcspResponse, error)
	ShowOcspUpdates() ([]*models.SslOcspUpdate, error)
	SetOcspResponse(payload string) error
	UpdateOcspResponse(name string) (*models.SslOcspResponse, error)
	ShowSSLProviders() (*models.SslProviders, error)
	SetRateLimitSSLSessionGlobal(value uint64) error
}

type Acme interface {
	AcmeRenew(certificate string) error
	AcmeStatus() (models.AcmeStatus, error)
}

type Runtime interface {
	Info
	Frontend
	Manage
	Maps
	Servers
	ACLs
	Tables
	Raw
	SSL
	Acme
	SocketPath() string
	IsStatsSocket() bool
}

func New(_ context.Context, opt ...options.RuntimeOption) (Runtime, error) {
	c := &client{
		options: options.RuntimeOptions{},
	}
	var err error

	for _, option := range opt {
		err = option.Set(&c.options)
		if err != nil {
			return nil, err
		}
	}

	if c.options.MasterSocketData != nil {
		err = c.initWithMasterSocket(c.options)
	} else {
		err = c.initWithSockets(c.options)
	}
	if err != nil {
		return nil, err
	}

	return c, nil
}
