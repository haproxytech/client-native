// Copyright 2019 HAProxy Technologies
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

package configuration

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v5/config-parser"
	parser_errors "github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/params"
	"github.com/haproxytech/client-native/v5/config-parser/types"

	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

type Server interface {
	GetServers(parentType string, parentName string, transactionID string) (int64, models.Servers, error)
	GetServer(name string, parentType string, parentName string, transactionID string) (int64, *models.Server, error)
	DeleteServer(name string, parentType string, parentName string, transactionID string, version int64) error
	CreateServer(parentType string, parentName string, data *models.Server, transactionID string, version int64) error
	EditServer(name string, parentType string, parentName string, data *models.Server, transactionID string, version int64) error
	GetServerSwitchingRules(backend string, transactionID string) (int64, models.ServerSwitchingRules, error)
	GetServerSwitchingRule(id int64, backend string, transactionID string) (int64, *models.ServerSwitchingRule, error)
	DeleteServerSwitchingRule(id int64, backend string, transactionID string, version int64) error
	CreateServerSwitchingRule(backend string, data *models.ServerSwitchingRule, transactionID string, version int64) error
	EditServerSwitchingRule(id int64, backend string, data *models.ServerSwitchingRule, transactionID string, version int64) error
}

// GetServers returns configuration version and an array of
// configured servers in the specified backend. Returns error on fail.
func (c *client) GetServers(parentType string, parentName string, transactionID string) (int64, models.Servers, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	servers, err := ParseServers(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, servers, nil
}

// GetServer returns configuration version and a requested server
// in the specified backend. Returns error on fail or if server does not exist.
func (c *client) GetServer(name string, parentType string, parentName string, transactionID string) (int64, *models.Server, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	server, _ := GetServerByName(name, parentType, parentName, p)
	if server == nil {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("server %s does not exist in %s %s", name, parentName, parentType))
	}

	return v, server, nil
}

// DeleteServer deletes a server in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteServer(name string, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	server, i := GetServerByName(name, parentType, parentName, p)
	if server == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("server %s does not exist in %s %s", name, parentName, parentType))
		return c.HandleError(name, parentType, parentName, t, transactionID == "", e)
	}

	if err := p.Delete(sectionType(parentType), parentName, "server", i); err != nil {
		return c.HandleError(name, parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateServer creates a server in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateServer(parentType string, parentName string, data *models.Server, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	server, _ := GetServerByName(data.Name, parentType, parentName, p)
	if server != nil {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("server %s already exists in %s %s", data.Name, parentName, parentType))
		return c.HandleError(data.Name, parentType, parentName, t, transactionID == "", e)
	}

	if err := p.Insert(sectionType(parentType), parentName, "server", SerializeServer(*data), -1); err != nil {
		return c.HandleError(data.Name, parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditServer edits a server in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditServer(name string, parentType string, parentName string, data *models.Server, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	server, i := GetServerByName(name, parentType, parentName, p)
	if server == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("server %v does not exist in %s %s", name, parentName, parentType))
		return c.HandleError(data.Name, parentType, parentName, t, transactionID == "", e)
	}

	if err := p.Set(sectionType(parentType), parentName, "server", SerializeServer(*data), i); err != nil {
		return c.HandleError(data.Name, parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseServers(parentType string, parentName string, p parser.Parser) (models.Servers, error) {
	servers := models.Servers{}

	data, err := p.Get(sectionType(parentType), parentName, "server", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return servers, nil
		}
		return nil, err
	}

	ondiskServers, ok := data.([]types.Server)
	if !ok {
		return nil, misc.CreateTypeAssertError("server")
	}
	for _, ondiskServer := range ondiskServers {
		s := ParseServer(ondiskServer)
		if s != nil {
			servers = append(servers, s)
		}
	}
	return servers, nil
}

func parseAddress(address string) (ipOrAddress string, port *int64) {
	if strings.HasPrefix(address, "[") && strings.ContainsRune(address, ']') { // IPv6 with port [2001:0DB8:0000:0000:0000:0000:1428:57ab]:80
		split := strings.Split(address, "]")
		split[0] = strings.TrimPrefix(split[0], "[")
		if len(split) == 2 { // has port
			split[1] = strings.ReplaceAll(split[1], ":", "")
			p, err := strconv.ParseInt(split[1], 10, 64)
			if err == nil {
				port = &p
			}
		}
		return split[0], port
	}

	switch c := strings.Count(address, ":"); {
	case c == 1: // IPv4 with port 127.0.0.1:80
		split := strings.Split(address, ":")
		p, err := strconv.ParseInt(split[1], 10, 64)
		if err == nil {
			port = &p
		}
		return split[0], port
	case c > 1: // IPv6 2001:0DB8:0000:0000:0000:0000:1428:57ab
		// Assume the last element is the port number.
		// This is an imperfect solution, which is why dataplaneapi
		// adds brackets to IPv6 when it can.
		idx := strings.LastIndex(address, ":")
		p, err := strconv.ParseUint(address[idx+1:], 10, 16)
		if err != nil {
			return address, nil
		}
		return address[:idx], misc.Int64P(int(p))
	case c == 0:
		return address, nil // IPv4 or socket address
	default:
		return "", nil
	}
}

func parseServerParams(serverOptions []params.ServerOption, serverParams *models.ServerParams) { //nolint:gocognit,gocyclo,cyclop,cyclop,maintidx
	var sourceArgs []string
	for _, p := range serverOptions {
		switch v := p.(type) {
		case *params.ServerOptionWord:
			switch v.Name {
			case "agent-check":
				serverParams.AgentCheck = "enabled"
			case "no-agent-check":
				serverParams.AgentCheck = "disabled"
			case "allow-0rtt":
				serverParams.Allow0rtt = true
			case "backup":
				serverParams.Backup = "enabled"
			case "no-backup":
				serverParams.Backup = "disabled"
			case "check":
				serverParams.Check = "enabled"
			case "no-check":
				serverParams.Check = "disabled"
			case "check-send-proxy":
				serverParams.CheckSendProxy = "enabled"
			case "check-ssl":
				serverParams.CheckSsl = "enabled"
			case "no-check-ssl":
				serverParams.CheckSsl = "disabled"
			case "check-via-socks4":
				serverParams.CheckViaSocks4 = "enabled"
			case "disabled":
				serverParams.Maintenance = "enabled"
			case "enabled":
				serverParams.Maintenance = "disabled"
			case "force-sslv3":
				serverParams.ForceSslv3 = "enabled"
			case "force-tlsv10":
				serverParams.ForceTlsv10 = "enabled"
			case "no-sslv3":
				serverParams.NoSslv3 = "enabled"
			case "no-tlsv10":
				serverParams.ForceTlsv10 = "disabled"
			case "force-tlsv11":
				serverParams.ForceTlsv11 = "enabled"
			case "no-tlsv11":
				serverParams.ForceTlsv11 = "disabled"
			case "force-tlsv12":
				serverParams.ForceTlsv12 = "enabled"
			case "no-tlsv12":
				serverParams.ForceTlsv12 = "disabled"
			case "force-tlsv13":
				serverParams.ForceTlsv13 = "enabled"
			case "no-tlsv13":
				serverParams.ForceTlsv13 = "disabled"
			case "send-proxy":
				serverParams.SendProxy = "enabled"
			case "no-send-proxy":
				serverParams.SendProxy = "disabled"
			case "send-proxy-v2":
				serverParams.SendProxyV2 = "enabled"
			case "no-send-proxy-v2":
				serverParams.SendProxyV2 = "disabled"
			case "send-proxy-v2-ssl":
				serverParams.SendProxyV2Ssl = "enabled"
			case "no-send-proxy-v2-ssl":
				serverParams.SendProxyV2Ssl = "disabled"
			case "send-proxy-v2-ssl-cn":
				serverParams.SendProxyV2SslCn = "enabled"
			case "no-send-proxy-v2-ssl-cn":
				serverParams.SendProxyV2SslCn = "disabled"
			case "ssl":
				serverParams.Ssl = "enabled"
			case "no-ssl":
				serverParams.Ssl = "disabled"
			case "ssl-reuse":
				serverParams.SslReuse = "enabled"
			case "no-ssl-reuse":
				serverParams.SslReuse = "disabled"
			case "tls-tickets":
				serverParams.TLSTickets = "enabled"
			case "no-tls-tickets":
				serverParams.TLSTickets = "disabled"
			case "tfo":
				serverParams.Tfo = "enabled"
			case "no-tfo":
				serverParams.Tfo = "disabled"
			case "stick":
				serverParams.Stick = "enabled"
			case "non-stick":
				serverParams.Stick = "disabled"
			}
		case *params.ServerOptionValue:
			switch v.Name {
			case "agent-send":
				serverParams.AgentSend = v.Value
			case "agent-inter":
				serverParams.AgentInter = misc.ParseTimeout(v.Value)
			case "agent-addr":
				serverParams.AgentAddr = v.Value
			case "agent-port":
				p, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && p != 0 {
					serverParams.AgentPort = &p
				}
			case "alpn":
				serverParams.Alpn = v.Value
			case "ca-file":
				serverParams.SslCafile = v.Value
			case "check-alpn":
				serverParams.CheckAlpn = v.Value
			case "check-proto":
				serverParams.CheckProto = v.Value
			case "check-sni":
				serverParams.CheckSni = v.Value
			case "ciphers":
				serverParams.Ciphers = v.Value
			case "ciphersuites":
				serverParams.Ciphersuites = v.Value
			case "client-sigalgs":
				serverParams.ClientSigalgs = v.Value
			case "cookie":
				serverParams.Cookie = v.Value
			case "crl-file":
				serverParams.CrlFile = v.Value
			case "crt":
				serverParams.SslCertificate = v.Value
			case "curves":
				serverParams.Curves = v.Value
			case "error-limit":
				c, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					serverParams.ErrorLimit = c
				}
			case "fall":
				serverParams.Fall = misc.ParseTimeout(v.Value)
			case "init-addr":
				serverParams.InitAddr = &v.Value
			case "inter":
				serverParams.Inter = misc.ParseTimeout(v.Value)
			case "fastinter":
				serverParams.Fastinter = misc.ParseTimeout(v.Value)
			case "downinter":
				serverParams.Downinter = misc.ParseTimeout(v.Value)
			case "log-bufsize":
				l, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					serverParams.LogBufsize = &l
				}
			case "log-proto":
				serverParams.LogProto = v.Value
			case "maxconn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					serverParams.Maxconn = &m
				}
			case "maxqueue":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					serverParams.Maxqueue = &m
				}
			case "max-reuse":
				c, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					serverParams.MaxReuse = &c
				}
			case "minconn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					serverParams.Minconn = &m
				}
			case "namespace":
				serverParams.Namespace = v.Value
			case "npn":
				serverParams.Npn = v.Value
			case "observe":
				serverParams.Observe = v.Value
			case "on-error":
				serverParams.OnError = v.Value
			case "on-marked-down":
				serverParams.OnMarkedDown = v.Value
			case "on-marked-up":
				serverParams.OnMarkedUp = v.Value
			case "pool-low-conn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					serverParams.PoolLowConn = &m
				}
			case "pool-max-conn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					serverParams.PoolMaxConn = &m
				}
			case "pool-purge-delay":
				serverParams.PoolPurgeDelay = misc.ParseTimeout(v.Value)
			case "addr":
				serverParams.HealthCheckAddress = v.Value
			case "port":
				p, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					serverParams.HealthCheckPort = &p
				}
			case "proto":
				serverParams.Proto = v.Value
			case "redir":
				serverParams.Redir = v.Value
			case "rise":
				serverParams.Rise = misc.ParseTimeout(v.Value)
			case "resolve-opts":
				serverParams.ResolveOpts = v.Value
			case "resolve-prefer":
				serverParams.ResolvePrefer = v.Value
			case "resolve-net":
				serverParams.ResolveNet = v.Value
			case "resolvers":
				serverParams.Resolvers = v.Value
			case "proxy-v2-options":
				serverParams.ProxyV2Options = strings.Split(v.Value, ",")
			case "shard":
				if v.Value != "" {
					p, err := strconv.ParseInt(v.Value, 10, 64)
					if err == nil {
						serverParams.Shard = p
					}
				}
			case "sigalgs":
				serverParams.Sigalgs = v.Value
			case "slowstart":
				serverParams.Slowstart = misc.ParseTimeout(v.Value)
			case "sni":
				serverParams.Sni = v.Value
			case "source":
				serverParams.Source = v.Value
			case "usesrc", "interface":
				sourceArgs = append(sourceArgs, v.Name+" "+v.Value)
			case "ssl-max-ver":
				serverParams.SslMaxVer = v.Value
			case "ssl-min-ver":
				serverParams.SslMinVer = v.Value
			case "socks4":
				serverParams.Socks4 = v.Value
			case "tcp-ut":
				serverParams.TCPUt = misc.ParseTimeout(v.Value)
			case "track":
				serverParams.Track = v.Value
			case "verify":
				serverParams.Verify = v.Value
			case "verifyhost":
				serverParams.Verifyhost = v.Value
			case "weight":
				w, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					serverParams.Weight = &w
				}
			case "ws":
				serverParams.Ws = v.Value
			}
		case *params.ServerOptionIDValue:
			if v.Name == "set-proxy-v2-tlv-fmt" {
				serverParams.SetProxyV2TlvFmt = &models.ServerParamsSetProxyV2TlvFmt{
					ID:    &v.ID,
					Value: &v.Value,
				}
			}
		}
	}
	// Add corresponding arguments to the source option.
	if serverParams.Source != "" && len(sourceArgs) > 0 {
		serverParams.Source += " " + strings.Join(sourceArgs, " ")
	}
}

func ParseServer(ondiskServer types.Server) *models.Server {
	s := &models.Server{
		Name: ondiskServer.Name,
	}
	address, port := parseAddress(ondiskServer.Address)
	if address == "" {
		return nil
	}
	s.Address = address
	s.Port = port
	for _, p := range ondiskServer.Params {
		if v, ok := p.(*params.ServerOptionValue); ok {
			if v.Name == "id" {
				p, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					s.ID = &p
				}
			}
		}
	}
	parseServerParams(ondiskServer.Params, &s.ServerParams)
	return s
}

func serializeServerParams(s models.ServerParams) (options []params.ServerOption) { //nolint:gocognit,gocyclo,cyclop,cyclop,maintidx
	// ServerOptionWord
	if s.AgentCheck == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "agent-check"})
	}
	if s.AgentCheck == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-agent-check"})
	}
	if s.Allow0rtt {
		options = append(options, &params.ServerOptionWord{Name: "allow-0rtt"})
	}
	if s.Backup == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "backup"})
	}
	if s.Backup == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-backup"})
	}
	if s.Check == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "check"})
	}
	if s.Check == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-check"})
	}
	if s.CheckSendProxy == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "check-send-proxy"})
	}
	if s.CheckSendProxy == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-check-send-proxy"})
	}
	if s.CheckSsl == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "check-ssl"})
	}
	if s.CheckSsl == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-check-ssl"})
	}
	if s.CheckViaSocks4 == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "check-via-socks4"})
	}
	if s.ForceSslv3 == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "force-sslv3"})
	}
	if s.NoSslv3 == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-sslv3"})
	}
	if s.ForceTlsv10 == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "force-tlsv10"})
	}
	if s.ForceTlsv10 == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-tlsv10"})
	}
	if s.ForceTlsv11 == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "force-tlsv11"})
	}
	if s.ForceTlsv11 == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-tlsv11"})
	}
	if s.ForceTlsv12 == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "force-tlsv12"})
	}
	if s.ForceTlsv12 == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-tlsv12"})
	}
	if s.ForceTlsv13 == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "force-tlsv13"})
	}
	if s.ForceTlsv13 == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-tlsv13"})
	}
	if s.Maintenance == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "disabled"})
	}
	if s.Maintenance == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "enabled"})
	}
	if s.SendProxy == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "send-proxy"})
	}
	if s.SendProxy == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-send-proxy"})
	}
	if s.SendProxyV2 == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "send-proxy-v2"})
	}
	if s.SendProxyV2 == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-send-proxy-v2"})
	}
	if s.SendProxyV2Ssl == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "send-proxy-v2-ssl"})
	}
	if s.SendProxyV2Ssl == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-send-proxy-v2-ssl"})
	}
	if s.SendProxyV2SslCn == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "send-proxy-v2-ssl-cn"})
	}
	if s.SendProxyV2SslCn == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-send-proxy-v2-ssl-cn"})
	}
	if s.Ssl == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "ssl"})
	}
	if s.Ssl == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-ssl"})
	}
	if s.SslReuse == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "ssl-reuse"})
	}
	if s.SslReuse == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-ssl-reuse"})
	}
	if s.TLSTickets == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "tls-tickets"})
	}
	if s.TLSTickets == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-tls-tickets"})
	}
	if s.Stick == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "stick"})
	}
	if s.Stick == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "non-stick"})
	}
	if s.Tfo == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "tfo"})
	}
	if s.Tfo == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-tfo"})
	}
	if s.AgentSend != "" {
		options = append(options, &params.ServerOptionValue{Name: "agent-send", Value: s.AgentSend})
	}
	if s.AgentInter != nil {
		options = append(options, &params.ServerOptionValue{Name: "agent-inter", Value: strconv.FormatInt(*s.AgentInter, 10)})
	}
	if s.AgentAddr != "" {
		options = append(options, &params.ServerOptionValue{Name: "agent-addr", Value: s.AgentAddr})
	}
	if s.AgentPort != nil {
		options = append(options, &params.ServerOptionValue{Name: "agent-port", Value: strconv.FormatInt(*s.AgentPort, 10)})
	}
	if s.Alpn != "" {
		options = append(options, &params.ServerOptionValue{Name: "alpn", Value: s.Alpn})
	}
	if s.SslCafile != "" { // ca-file
		options = append(options, &params.ServerOptionValue{Name: "ca-file", Value: s.SslCafile})
	}
	if s.CheckAlpn != "" {
		options = append(options, &params.ServerOptionValue{Name: "check-alpn", Value: s.CheckAlpn})
	}
	if s.CheckProto != "" {
		options = append(options, &params.ServerOptionValue{Name: "check-proto", Value: s.CheckProto})
	}
	if s.CheckSni != "" {
		options = append(options, &params.ServerOptionValue{Name: "check-sni", Value: s.CheckSni})
	}
	if s.Ciphers != "" {
		options = append(options, &params.ServerOptionValue{Name: "ciphers", Value: s.Ciphers})
	}
	if s.Ciphersuites != "" {
		options = append(options, &params.ServerOptionValue{Name: "ciphersuites", Value: s.Ciphersuites})
	}
	if s.ClientSigalgs != "" {
		options = append(options, &params.ServerOptionValue{Name: "client-sigalgs", Value: s.ClientSigalgs})
	}
	if s.Cookie != "" {
		options = append(options, &params.ServerOptionValue{Name: "cookie", Value: s.Cookie})
	}
	if s.CrlFile != "" {
		options = append(options, &params.ServerOptionValue{Name: "crl-file", Value: s.CrlFile})
	}
	if s.Curves != "" {
		options = append(options, &params.ServerOptionValue{Name: "curves", Value: s.Curves})
	}
	if s.SslCertificate != "" {
		options = append(options, &params.ServerOptionValue{Name: "crt", Value: s.SslCertificate})
	}
	if s.ErrorLimit != 0 {
		options = append(options, &params.ServerOptionValue{Name: "error-limit", Value: strconv.FormatInt(s.ErrorLimit, 10)})
	}
	if s.Fall != nil {
		options = append(options, &params.ServerOptionValue{Name: "fall", Value: strconv.FormatInt(*s.Fall, 10)})
	}
	if s.InitAddr != nil {
		options = append(options, &params.ServerOptionValue{Name: "init-addr", Value: *s.InitAddr})
	}
	if s.Inter != nil {
		options = append(options, &params.ServerOptionValue{Name: "inter", Value: strconv.FormatInt(*s.Inter, 10)})
	}
	if s.Fastinter != nil {
		options = append(options, &params.ServerOptionValue{Name: "fastinter", Value: strconv.FormatInt(*s.Fastinter, 10)})
	}
	if s.Downinter != nil {
		options = append(options, &params.ServerOptionValue{Name: "downinter", Value: strconv.FormatInt(*s.Downinter, 10)})
	}
	if s.LogBufsize != nil {
		options = append(options, &params.ServerOptionValue{Name: "log-bufsize", Value: strconv.FormatInt(*s.LogBufsize, 10)})
	}
	if s.LogProto != "" {
		options = append(options, &params.ServerOptionValue{Name: "log-proto", Value: s.LogProto})
	}
	if s.Maxconn != nil {
		options = append(options, &params.ServerOptionValue{Name: "maxconn", Value: strconv.FormatInt(*s.Maxconn, 10)})
	}
	if s.Maxqueue != nil {
		options = append(options, &params.ServerOptionValue{Name: "maxqueue", Value: strconv.FormatInt(*s.Maxqueue, 10)})
	}
	if s.MaxReuse != nil {
		options = append(options, &params.ServerOptionValue{Name: "max-reuse", Value: strconv.FormatInt(*s.MaxReuse, 10)})
	}
	if s.Minconn != nil {
		options = append(options, &params.ServerOptionValue{Name: "minconn", Value: strconv.FormatInt(*s.Minconn, 10)})
	}
	if s.Namespace != "" {
		options = append(options, &params.ServerOptionValue{Name: "namespace", Value: s.Namespace})
	}
	if s.Npn != "" {
		options = append(options, &params.ServerOptionValue{Name: "npn", Value: s.Npn})
	}
	if s.Observe != "" {
		options = append(options, &params.ServerOptionValue{Name: "observe", Value: s.Observe})
	}
	if s.OnError != "" {
		options = append(options, &params.ServerOptionValue{Name: "on-error", Value: s.OnError})
	}
	if s.OnMarkedDown != "" {
		options = append(options, &params.ServerOptionValue{Name: "on-marked-down", Value: s.OnMarkedDown})
	}
	if s.OnMarkedUp != "" {
		options = append(options, &params.ServerOptionValue{Name: "on-marked-up", Value: s.OnMarkedUp})
	}
	if s.PoolLowConn != nil {
		options = append(options, &params.ServerOptionValue{Name: "pool-low-conn", Value: strconv.FormatInt(*s.PoolLowConn, 10)})
	}
	if s.PoolMaxConn != nil {
		options = append(options, &params.ServerOptionValue{Name: "pool-max-conn", Value: strconv.FormatInt(*s.PoolMaxConn, 10)})
	}
	if s.PoolPurgeDelay != nil {
		options = append(options, &params.ServerOptionValue{Name: "pool-purge-delay", Value: strconv.FormatInt(*s.PoolPurgeDelay, 10)})
	}
	if s.HealthCheckAddress != "" {
		options = append(options, &params.ServerOptionValue{Name: "addr", Value: s.HealthCheckAddress})
	}
	if s.HealthCheckPort != nil {
		options = append(options, &params.ServerOptionValue{Name: "port", Value: strconv.FormatInt(*s.HealthCheckPort, 10)})
	}
	if s.Proto != "" {
		options = append(options, &params.ServerOptionValue{Name: "proto", Value: s.Proto})
	}
	if s.Redir != "" {
		options = append(options, &params.ServerOptionValue{Name: "redir", Value: s.Redir})
	}
	if s.Rise != nil {
		options = append(options, &params.ServerOptionValue{Name: "rise", Value: strconv.FormatInt(*s.Rise, 10)})
	}
	if s.ResolveOpts != "" {
		options = append(options, &params.ServerOptionValue{Name: "resolve-opts", Value: s.ResolveOpts})
	}
	if s.ResolvePrefer != "" {
		options = append(options, &params.ServerOptionValue{Name: "resolve-prefer", Value: s.ResolvePrefer})
	}
	if s.ResolveNet != "" {
		options = append(options, &params.ServerOptionValue{Name: "resolve-net", Value: s.ResolveNet})
	}
	if s.Resolvers != "" {
		options = append(options, &params.ServerOptionValue{Name: "resolvers", Value: s.Resolvers})
	}
	if len(s.ProxyV2Options) > 0 {
		options = append(options, &params.ServerOptionValue{Name: "proxy-v2-options", Value: strings.Join(s.ProxyV2Options, ",")})
	}
	if s.Shard != 0 {
		options = append(options, &params.ServerOptionValue{Name: "shard", Value: strconv.FormatInt(s.Shard, 10)})
	}
	if s.Sigalgs != "" {
		options = append(options, &params.ServerOptionValue{Name: "sigalgs", Value: s.Sigalgs})
	}
	if s.Slowstart != nil {
		options = append(options, &params.ServerOptionValue{Name: "slowstart", Value: strconv.FormatInt(*s.Slowstart, 10)})
	}
	if s.Sni != "" {
		options = append(options, &params.ServerOptionValue{Name: "sni", Value: s.Sni})
	}
	if s.Source != "" {
		options = append(options, &params.ServerOptionValue{Name: "source", Value: s.Source})
	}
	if s.SslMaxVer != "" {
		options = append(options, &params.ServerOptionValue{Name: "ssl-max-ver", Value: s.SslMaxVer})
	}
	if s.SslMinVer != "" {
		options = append(options, &params.ServerOptionValue{Name: "ssl-min-ver", Value: s.SslMinVer})
	}
	if s.Socks4 != "" {
		options = append(options, &params.ServerOptionValue{Name: "socks4", Value: s.Socks4})
	}
	if s.SetProxyV2TlvFmt != nil && s.SetProxyV2TlvFmt.ID != nil && s.SetProxyV2TlvFmt.Value != nil {
		options = append(options, &params.ServerOptionIDValue{Name: "set-proxy-v2-tlv-fmt", ID: *s.SetProxyV2TlvFmt.ID, Value: *s.SetProxyV2TlvFmt.Value})
	}
	if s.TCPUt != nil {
		options = append(options, &params.ServerOptionValue{Name: "tcp-ut", Value: strconv.FormatInt(*s.TCPUt, 10)})
	}
	if s.Track != "" {
		options = append(options, &params.ServerOptionValue{Name: "track", Value: s.Track})
	}
	if s.Verify != "" {
		options = append(options, &params.ServerOptionValue{Name: "verify", Value: s.Verify})
	}
	if s.Verifyhost != "" {
		options = append(options, &params.ServerOptionValue{Name: "verifyhost", Value: s.Verifyhost})
	}
	if s.Weight != nil {
		options = append(options, &params.ServerOptionValue{Name: "weight", Value: strconv.FormatInt(*s.Weight, 10)})
	}
	if s.Ws != "" {
		options = append(options, &params.ServerOptionValue{Name: "ws", Value: s.Ws})
	}
	return options
}

func SerializeServer(s models.Server) types.Server {
	server := types.Server{
		Name:   s.Name,
		Params: []params.ServerOption{},
	}
	if s.Port != nil {
		server.Address = fmt.Sprintf("%s:%d", misc.SanitizeIPv6Address(s.Address), *s.Port)
	} else {
		server.Address = misc.SanitizeIPv6Address(s.Address)
	}
	server.Params = serializeServerParams(s.ServerParams)
	if s.ID != nil {
		server.Params = append(server.Params, &params.ServerOptionValue{Name: "id", Value: strconv.FormatInt(*s.ID, 10)})
	}
	return server
}

func GetServerByName(name string, parentType string, parentName string, p parser.Parser) (*models.Server, int) {
	servers, err := ParseServers(parentType, parentName, p)
	if err != nil {
		return nil, 0
	}

	for i, s := range servers {
		if s.Name == name {
			return s, i
		}
	}
	return nil, 0
}

func sectionType(parentType string) parser.Section {
	var sectionType parser.Section
	switch parentType {
	case "backend":
		sectionType = parser.Backends
	case "ring":
		sectionType = parser.Ring
	case "peers":
		sectionType = parser.Peers
	}
	return sectionType
}
