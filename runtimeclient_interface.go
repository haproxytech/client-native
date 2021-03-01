// This file is generated, don't edit manually, see README.md for details.

package clientnative

import (
	"io"
	"mime/multipart"

	"github.com/haproxytech/client-native/v2/models"
)

// IRuntimeClient ...
type IRuntimeClient interface {
	// Init must be given path to runtime socket and nbproc that is not 0 when in master worker mode
	//
	// Deprecated: use InitWithSockets or InitWithMasterSocket instead
	Init(socketPath []string, masterSocketPath string, nbproc int) error
	// GetMapsPath returns runtime map file path or map id
	GetMapsPath(name string) (string, error)
	InitWithSockets(socketPath map[int]string) error
	InitWithMasterSocket(masterSocketPath string, nbproc int) error
	// GetStats returns stats from the socket
	GetStats() models.NativeStats
	// GetInfo returns info from the socket
	GetInfo() (models.ProcessInfos, error)
	// SetFrontendMaxConn set maxconn for frontend
	SetFrontendMaxConn(frontend string, maxconn int) error
	// SetServerAddr set ip [port] for server
	SetServerAddr(backend, server string, ip string, port int) error
	// SetServerState set state for server
	SetServerState(backend, server string, state string) error
	// SetServerWeight set weight for server
	SetServerWeight(backend, server string, weight string) error
	// SetServerHealth set health for server
	SetServerHealth(backend, server string, health string) error
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
	// Show tables show tables from runtime API and return it structured, if process is 0, return for all processes
	ShowTables(process int) (models.StickTables, error)
	// GetTableEntries returns all entries for specified table in the given process with filters and a key
	GetTableEntries(name string, process int, filter []string, key string) (models.StickTableEntries, error)
	// Show table show tables {name} from runtime API associated with process id and return it structured
	ShowTable(name string, process int) (*models.StickTable, error)
	// ExecuteRaw does not procces response, just returns its values for all processes
	ExecuteRaw(command string) ([]string, error)
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
}
