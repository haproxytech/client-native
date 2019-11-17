// This file is generated, don't edit manually, see README.md for details.

package client_native

import (
	"github.com/haproxytech/models"
)

// IRuntimeClient ...
type IRuntimeClient interface {
	//Init must be given path to runtime socket and nbproc that is not 0 when in master worker mode
	//
	//Deprecated: use InitWithSockets or InitWithMasterSocket instead
	Init(socketPath []string, masterSocketPath string, nbproc int) error
	InitWithSockets(socketPath map[int]string) error
	InitWithMasterSocket(masterSocketPath string, nbproc int) error
	//GetStats returns stats from the socket
	GetStats() models.NativeStats
	//GetInfo returns info from the socket
	GetInfo() (models.ProcessInfos, error)
	//SetFrontendMaxConn set maxconn for frontend
	SetFrontendMaxConn(frontend string, maxconn int) error
	//SetServerAddr set ip [port] for server
	SetServerAddr(backend, server string, ip string, port int) error
	//SetServerState set state for server
	SetServerState(backend, server string, state string) error
	//SetServerWeight set weight for server
	SetServerWeight(backend, server string, weight string) error
	//SetServerHealth set health for server
	SetServerHealth(backend, server string, health string) error
	//EnableAgentCheck enable agent check for server
	EnableAgentCheck(backend, server string) error
	//DisableAgentCheck disable agent check for server
	DisableAgentCheck(backend, server string) error
	//EnableServer marks server as UP
	EnableServer(backend, server string) error
	//DisableServer marks server as DOWN for maintenance
	DisableServer(backend, server string) error
	//SetServerAgentAddr set agent-addr for server
	SetServerAgentAddr(backend, server string, addr string) error
	//SetServerAgentSend set agent-send for server
	SetServerAgentSend(backend, server string, send string) error
	//GetServerState returns server runtime state
	GetServersState(backend string) (models.RuntimeServers, error)
	//GetServerState returns server runtime state
	GetServerState(backend, server string) (*models.RuntimeServer, error)
	//SetServerCheckPort set health heck port for server
	SetServerCheckPort(backend, server string, port int) error
	//Show tables show tables from runtime API and return it structured, if process is 0, return for all processes
	ShowTables(process int) (models.StickTables, error)
	//GetTableEntries returns all entries for specified table in the given process with filters and a key
	GetTableEntries(name string, process int, filter []string, key string) (models.StickTableEntries, error)
	//Show table show tables {name} from runtime API associated with process id and return it structured
	ShowTable(name string, process int) (*models.StickTable, error)
	//ExecuteRaw does not procces response, just returns its values for all processes
	ExecuteRaw(command string) ([]string, error)
}
