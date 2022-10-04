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
	parser "github.com/haproxytech/config-parser/v4"
	parser_errors "github.com/haproxytech/config-parser/v4/errors"
	"github.com/haproxytech/config-parser/v4/params"
	"github.com/haproxytech/config-parser/v4/types"

	"github.com/haproxytech/client-native/v4/misc"
	"github.com/haproxytech/client-native/v4/models"
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

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
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

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
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

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
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

func parseServerParams(serverOptions []params.ServerOption) (s models.ServerParams) { //nolint:gocognit,gocyclo,cyclop,cyclop,maintidx
	for _, p := range serverOptions {
		switch v := p.(type) {
		case *params.ServerOptionWord:
			switch v.Name {
			case "agent-check":
				s.AgentCheck = "enabled"
			case "no-agent-check":
				s.AgentCheck = "disabled"
			case "allow-0rtt":
				s.Allow0rtt = true
			case "backup":
				s.Backup = "enabled"
			case "no-backup":
				s.Backup = "disabled"
			case "check":
				s.Check = "enabled"
			case "no-check":
				s.Check = "disabled"
			case "check-send-proxy":
				s.CheckSendProxy = "enabled"
			case "check-ssl":
				s.CheckSsl = "enabled"
			case "no-check-ssl":
				s.CheckSsl = "disabled"
			case "check-via-socks4":
				s.CheckViaSocks4 = "enabled"
			case "disabled":
				s.Maintenance = "enabled"
			case "enabled":
				s.Maintenance = "disabled"
			case "force-sslv3":
				s.ForceSslv3 = "enabled"
			case "force-tlsv10":
				s.ForceTlsv10 = "enabled"
			case "no-tlsv10":
				s.ForceTlsv10 = "disabled"
			case "force-tlsv11":
				s.ForceTlsv11 = "enabled"
			case "no-tlsv11":
				s.ForceTlsv11 = "disabled"
			case "force-tlsv12":
				s.ForceTlsv12 = "enabled"
			case "no-tlsv12":
				s.ForceTlsv12 = "disabled"
			case "force-tlsv13":
				s.ForceTlsv13 = "enabled"
			case "no-tlsv13":
				s.ForceTlsv13 = "disabled"
			case "send-proxy":
				s.SendProxy = "enabled"
			case "no-send-proxy":
				s.SendProxy = "disabled"
			case "send-proxy-v2":
				s.SendProxyV2 = "enabled"
			case "no-send-proxy-v2":
				s.SendProxyV2 = "disabled"
			case "send-proxy-v2-ssl":
				s.SendProxyV2Ssl = "enabled"
			case "send-proxy-v2-ssl-cn":
				s.SendProxyV2SslCn = "enabled"
			case "ssl":
				s.Ssl = "enabled"
			case "no-ssl":
				s.Ssl = "disabled"
			case "ssl-reuse":
				s.SslReuse = "enabled"
			case "no-ssl-reuse":
				s.SslReuse = "disabled"
			case "tls-tickets":
				s.TLSTickets = "enabled"
			case "no-tls-tickets":
				s.TLSTickets = "disabled"
			case "tfo":
				s.Tfo = "enabled"
			case "no-tfo":
				s.Tfo = "disabled"
			case "stick":
				s.Stick = "enabled"
			case "no-stick":
				s.Stick = "disabled"
			}
		case *params.ServerOptionValue:
			switch v.Name {
			case "agent-send":
				s.AgentSend = v.Value
			case "agent-inter":
				s.AgentInter = misc.ParseTimeout(v.Value)
			case "agent-addr":
				s.AgentAddr = v.Value
			case "agent-port":
				p, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && p != 0 {
					s.AgentPort = &p
				}
			case "alpn":
				s.Alpn = v.Value
			case "ca-file":
				s.SslCafile = v.Value
			case "check-alpn":
				s.CheckAlpn = v.Value
			case "check-proto":
				s.CheckProto = v.Value
			case "check-sni":
				s.CheckSni = v.Value
			case "ciphers":
				s.Ciphers = v.Value
			case "ciphersuites":
				s.Ciphersuites = v.Value
			case "cookie":
				s.Cookie = v.Value
			case "crl-file":
				s.CrlFile = v.Value
			case "crt":
				s.SslCertificate = v.Value
			case "error-limit":
				c, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					s.ErrorLimit = c
				}
			case "fall":
				c, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					s.Fall = &c
				}
			case "init-addr":
				s.InitAddr = &v.Value
			case "inter":
				s.Inter = misc.ParseTimeout(v.Value)
			case "fastinter":
				s.Fastinter = misc.ParseTimeout(v.Value)
			case "downinter":
				s.Downinter = misc.ParseTimeout(v.Value)
			case "log-proto":
				s.LogProto = v.Value
			case "maxconn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					s.Maxconn = &m
				}
			case "maxqueue":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					s.Maxqueue = &m
				}
			case "max-reuse":
				c, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					s.MaxReuse = &c
				}
			case "minconn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					s.Minconn = &m
				}
			case "namespace":
				s.Namespace = v.Value
			case "npn":
				s.Npn = v.Value
			case "observe":
				s.Observe = v.Value
			case "on-error":
				s.OnError = v.Value
			case "on-marked-down":
				s.OnMarkedDown = v.Value
			case "on-marked-up":
				s.OnMarkedUp = v.Value
			case "pool-low-conn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					s.PoolLowConn = &m
				}
			case "pool-max-conn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					s.PoolMaxConn = &m
				}
			case "pool-purge-delay":
				d, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					s.PoolPurgeDelay = &d
				}
			case "port":
				p, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					s.HealthCheckPort = &p
				}
			case "proto":
				s.Proto = v.Value
			case "redir":
				s.Redir = v.Value
			case "rise":
				s.Rise = misc.ParseTimeout(v.Value)
			case "resolve-opts":
				s.ResolveOpts = v.Value
			case "resolve-prefer":
				s.ResolvePrefer = v.Value
			case "resolve-net":
				s.ResolveNet = v.Value
			case "resolvers":
				s.Resolvers = v.Value
			case "proxy-v2-options":
				s.ProxyV2Options = strings.Split(v.Value, ",")
			case "slowstart":
				s.Slowstart = misc.ParseTimeout(v.Value)
			case "sni":
				s.Sni = v.Value
			case "source":
				s.Source = v.Value
			case "ssl-max-ver":
				s.SslMaxVer = v.Value
			case "ssl-min-ver":
				s.SslMinVer = v.Value
			case "socks4":
				s.Socks4 = v.Value
			case "tcp-ut":
				d, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					s.TCPUt = d
				}
			case "track":
				s.Track = v.Value
			case "verify":
				s.Verify = v.Value
			case "verifyhost":
				s.Verifyhost = v.Value
			case "weight":
				w, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					s.Weight = &w
				}
			case "ws":
				s.Ws = v.Value
			}
		}
	}
	return s
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
	s.ServerParams = parseServerParams(ondiskServer.Params)
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
	if s.ForceSslv3 == "disabled" {
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
		options = append(options, &params.ServerOptionWord{Name: "no-stick"})
	}
	if s.Tfo == "enabled" {
		options = append(options, &params.ServerOptionWord{Name: "tfo"})
	}
	if s.Tfo == "disabled" {
		options = append(options, &params.ServerOptionWord{Name: "no-tfo"})
	}
	// ServerOptionValue
	/*
		if s.ID != nil {
			options = append(options, &params.ServerOptionValue{Name: "id", Value: strconv.FormatInt(*s.ID, 10)})
		}*/
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
	if s.Cookie != "" {
		options = append(options, &params.ServerOptionValue{Name: "cookie", Value: s.Cookie})
	}
	if s.CrlFile != "" {
		options = append(options, &params.ServerOptionValue{Name: "crl-file", Value: s.CrlFile})
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
		options = append(options, &params.ServerOptionValue{Name: "non", Value: s.Npn})
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
	if s.TCPUt != 0 {
		options = append(options, &params.ServerOptionValue{Name: "tcp-ut", Value: strconv.FormatInt(s.TCPUt, 10)})
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
		if misc.IsIPv6(s.Address) {
			server.Address = fmt.Sprintf("[%s]:%d", s.Address, *s.Port)
		} else {
			server.Address = fmt.Sprintf("%s:%d", s.Address, *s.Port)
		}
	} else {
		server.Address = s.Address
	}
	server.Params = serializeServerParams(s.ServerParams)
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
