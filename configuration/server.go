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

	"github.com/haproxytech/client-native/v3/misc"
	"github.com/haproxytech/client-native/v3/models"
)

type Server interface {
	GetServers(backend string, transactionID string) (int64, models.Servers, error)
	GetServer(name string, backend string, transactionID string) (int64, *models.Server, error)
	DeleteServer(name string, backend string, transactionID string, version int64) error
	CreateServer(backend string, data *models.Server, transactionID string, version int64) error
	EditServer(name string, backend string, data *models.Server, transactionID string, version int64) error
	GetServerSwitchingRules(backend string, transactionID string) (int64, models.ServerSwitchingRules, error)
	GetServerSwitchingRule(id int64, backend string, transactionID string) (int64, *models.ServerSwitchingRule, error)
	DeleteServerSwitchingRule(id int64, backend string, transactionID string, version int64) error
	CreateServerSwitchingRule(backend string, data *models.ServerSwitchingRule, transactionID string, version int64) error
	EditServerSwitchingRule(id int64, backend string, data *models.ServerSwitchingRule, transactionID string, version int64) error
}

// GetServers returns configuration version and an array of
// configured servers in the specified backend. Returns error on fail.
func (c *client) GetServers(backend string, transactionID string) (int64, models.Servers, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	servers, err := ParseServers(backend, p)
	if err != nil {
		return v, nil, c.HandleError("", "backend", backend, "", false, err)
	}

	return v, servers, nil
}

// GetServer returns configuration version and a requested server
// in the specified backend. Returns error on fail or if server does not exist.
func (c *client) GetServer(name string, backend string, transactionID string) (int64, *models.Server, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	server, _ := GetServerByName(name, backend, p)
	if server == nil {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Server %s does not exist in backend %s", name, backend))
	}

	return v, server, nil
}

// DeleteServer deletes a server in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteServer(name string, backend string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	server, i := GetServerByName(name, backend, p)
	if server == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Server %s does not exist in backend %s", name, backend))
		return c.HandleError(name, "backend", backend, t, transactionID == "", e)
	}

	if err := p.Delete(parser.Backends, backend, "server", i); err != nil {
		return c.HandleError(name, "backend", backend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// CreateServer creates a server in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateServer(backend string, data *models.Server, transactionID string, version int64) error {
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

	server, _ := GetServerByName(data.Name, backend, p)
	if server != nil {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("Server %s already exists in backend %s", data.Name, backend))
		return c.HandleError(data.Name, "backend", backend, t, transactionID == "", e)
	}

	if err := p.Insert(parser.Backends, backend, "server", SerializeServer(*data), -1); err != nil {
		return c.HandleError(data.Name, "backend", backend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// EditServer edits a server in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditServer(name string, backend string, data *models.Server, transactionID string, version int64) error {
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

	server, i := GetServerByName(name, backend, p)
	if server == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Server %v does not exist in backend %s", name, backend))
		return c.HandleError(data.Name, "backend", backend, t, transactionID == "", e)
	}

	if err := p.Set(parser.Backends, backend, "server", SerializeServer(*data), i); err != nil {
		return c.HandleError(data.Name, "backend", backend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func ParseServers(backend string, p parser.Parser) (models.Servers, error) {
	servers := models.Servers{}

	data, err := p.Get(parser.Backends, backend, "server", false)
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
		return address, nil
	case c == 0:
		return address, nil // IPv4 or socket address
	default:
		return "", nil
	}
}

func ParseServer(ondiskServer types.Server) *models.Server { //nolint:gocyclo,cyclop,gocognit
	s := &models.Server{
		Name: ondiskServer.Name,
	}
	address, port := parseAddress(ondiskServer.Address)
	if address == "" {
		return nil
	}
	s.Address = address
	s.Port = port
	for _, p := range ondiskServer.Params { //nolint:dupl
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
				if err == nil && c != 0 {
					s.ErrorLimit = c
				}
			case "fall":
				c, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && c != 0 {
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
				if err == nil && m != 0 {
					s.Maxconn = &m
				}
			case "maxqueue":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && m != 0 {
					s.Maxqueue = &m
				}
			case "max-reuse":
				c, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && c != 0 {
					s.MaxReuse = &c
				}
			case "minconn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && m != 0 {
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
				if err == nil && m != 0 {
					s.PoolLowConn = &m
				}
			case "pool-max-conn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && m != 0 {
					s.PoolMaxConn = &m
				}
			case "pool-purge-delay":
				d, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && d != 0 {
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
				if err == nil && d != 0 {
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
				if err == nil && w != 0 {
					s.Weight = &w
				}
			}
		}
	}
	return s
}

func SerializeServer(s models.Server) types.Server { //nolint:gocognit,gocyclo,cyclop,cyclop
	srv := types.Server{
		Name:   s.Name,
		Params: []params.ServerOption{},
	}
	// ServerOptionWord
	if s.AgentCheck == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "agent-check"})
	}
	if s.AgentCheck == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-agent-check"})
	}
	if s.Allow0rtt {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "allow-0rtt"})
	}
	if s.Backup == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "backup"})
	}
	if s.Backup == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-backup"})
	}
	if s.Check == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "check"})
	}
	if s.Check == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-check"})
	}
	if s.CheckSendProxy == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "check-send-proxy"})
	}
	if s.CheckSendProxy == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "check-send-proxy"})
	}
	if s.CheckSsl == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "check-ssl"})
	}
	if s.CheckSsl == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-check-ssl"})
	}
	if s.CheckViaSocks4 == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "check-via-socks4"})
	}
	if s.ForceSslv3 == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "force-sslv3"})
	}
	if s.ForceSslv3 == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-sslv3"})
	}
	if s.ForceTlsv10 == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "force-tlsv10"})
	}
	if s.ForceTlsv10 == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-tlsv10"})
	}
	if s.ForceTlsv11 == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "force-tlsv11"})
	}
	if s.ForceTlsv11 == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-tlsv11"})
	}
	if s.ForceTlsv12 == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "force-tlsv12"})
	}
	if s.ForceTlsv12 == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-tlsv12"})
	}
	if s.ForceTlsv13 == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "force-tlsv13"})
	}
	if s.ForceTlsv13 == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-tlsv13"})
	}
	if s.Maintenance == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "disabled"})
	}
	if s.Maintenance == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "enabled"})
	}
	if s.SendProxy == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "send-proxy"})
	}
	if s.SendProxy == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-send-proxy"})
	}
	if s.SendProxyV2 == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "send-proxy-v2"})
	}
	if s.SendProxyV2 == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-send-proxy-v2"})
	}
	if s.SendProxyV2Ssl == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "send-proxy-v2-ssl"})
	}
	if s.SendProxyV2Ssl == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-send-proxy-v2-ssl"})
	}
	if s.SendProxyV2SslCn == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "send-proxy-v2-ssl-cn"})
	}
	if s.SendProxyV2SslCn == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-send-proxy-v2-ssl-cn"})
	}
	if s.Ssl == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "ssl"})
	}
	if s.Ssl == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-ssl"})
	}
	if s.SslReuse == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "ssl-reuse"})
	}
	if s.SslReuse == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-ssl-reuse"})
	}
	if s.TLSTickets == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "tls-tickets"})
	}
	if s.TLSTickets == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-tls-tickets"})
	}
	if s.Stick == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "stick"})
	}
	if s.Stick == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-stick"})
	}
	if s.Tfo == "enabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "tfo"})
	}
	if s.Tfo == "disabled" {
		srv.Params = append(srv.Params, &params.ServerOptionWord{Name: "no-tfo"})
	}
	// ServerOptionValue
	if s.AgentSend != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "agent-send", Value: s.AgentSend})
	}
	if s.AgentInter != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "agent-inter", Value: strconv.FormatInt(*s.AgentInter, 10)})
	}
	if s.AgentAddr != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "agent-addr", Value: s.AgentAddr})
	}
	if s.AgentPort != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "agent-port", Value: strconv.FormatInt(*s.AgentPort, 10)})
	}
	if s.Alpn != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "alpn", Value: s.Alpn})
	}
	if s.SslCafile != "" { // ca-file
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "ca-file", Value: s.SslCafile})
	}
	if s.CheckAlpn != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "check-alpn", Value: s.CheckAlpn})
	}
	if s.CheckProto != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "check-proto", Value: s.CheckProto})
	}
	if s.CheckSni != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "check-sni", Value: s.CheckSni})
	}
	if s.Ciphers != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "ciphers", Value: s.Ciphers})
	}
	if s.Ciphersuites != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "ciphersuites", Value: s.Ciphersuites})
	}
	if s.Cookie != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "cookie", Value: s.Cookie})
	}
	if s.CrlFile != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "crl-file", Value: s.CrlFile})
	}
	if s.SslCertificate != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "crt", Value: s.SslCertificate})
	}
	if s.ErrorLimit != 0 {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "error-limit", Value: strconv.FormatInt(s.ErrorLimit, 10)})
	}
	if s.Fall != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "fall", Value: strconv.FormatInt(*s.Fall, 10)})
	}
	if s.InitAddr != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "init-addr", Value: *s.InitAddr})
	}
	if s.Inter != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "inter", Value: strconv.FormatInt(*s.Inter, 10)})
	}
	if s.Fastinter != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "fastinter", Value: strconv.FormatInt(*s.Fastinter, 10)})
	}
	if s.Downinter != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "downinter", Value: strconv.FormatInt(*s.Downinter, 10)})
	}
	if s.LogProto != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "log-proto", Value: s.LogProto})
	}
	if s.Maxconn != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "maxconn", Value: strconv.FormatInt(*s.Maxconn, 10)})
	}
	if s.Maxqueue != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "maxqueue", Value: strconv.FormatInt(*s.Maxqueue, 10)})
	}
	if s.MaxReuse != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "max-reuse", Value: strconv.FormatInt(*s.MaxReuse, 10)})
	}
	if s.Minconn != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "minconn", Value: strconv.FormatInt(*s.Minconn, 10)})
	}
	if s.Namespace != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "namespace", Value: s.Namespace})
	}
	if s.Npn != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "non", Value: s.Npn})
	}
	if s.Observe != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "observe", Value: s.Observe})
	}
	if s.OnError != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "on-error", Value: s.OnError})
	}
	if s.OnMarkedDown != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "on-marked-down", Value: s.OnMarkedDown})
	}
	if s.OnMarkedUp != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "on-marked-up", Value: s.OnMarkedUp})
	}
	if s.PoolLowConn != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "pool-low-conn", Value: strconv.FormatInt(*s.PoolLowConn, 10)})
	}
	if s.PoolMaxConn != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "pool-max-conn", Value: strconv.FormatInt(*s.PoolMaxConn, 10)})
	}
	if s.PoolPurgeDelay != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "pool-purge-delay", Value: strconv.FormatInt(*s.PoolPurgeDelay, 10)})
	}
	if s.Port != nil {
		srv.Address = s.Address + ":" + strconv.FormatInt(*s.Port, 10)
	} else {
		srv.Address = s.Address
	}
	if s.HealthCheckPort != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "port", Value: strconv.FormatInt(*s.HealthCheckPort, 10)})
	}
	if s.Proto != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "proto", Value: s.Proto})
	}
	if s.Redir != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "redir", Value: s.Redir})
	}
	if s.Rise != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "rise", Value: strconv.FormatInt(*s.Rise, 10)})
	}
	if s.ResolveOpts != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "resolve-opts", Value: s.ResolveOpts})
	}
	if s.ResolvePrefer != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "resolve-prefer", Value: s.ResolvePrefer})
	}
	if s.ResolveNet != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "resolve-net", Value: s.ResolveNet})
	}
	if s.Resolvers != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "resolvers", Value: s.Resolvers})
	}
	if len(s.ProxyV2Options) > 0 {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "proxy-v2-options", Value: strings.Join(s.ProxyV2Options, ",")})
	}
	if s.Slowstart != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "slowstart", Value: strconv.FormatInt(*s.Slowstart, 10)})
	}
	if s.Sni != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "sni", Value: s.Sni})
	}
	if s.Source != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "source", Value: s.Source})
	}
	if s.SslMaxVer != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "ssl-max-ver", Value: s.SslMaxVer})
	}
	if s.SslMinVer != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "ssl-min-ver", Value: s.SslMinVer})
	}
	if s.Socks4 != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "socks4", Value: s.Socks4})
	}
	if s.TCPUt != 0 {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "tcp-ut", Value: strconv.FormatInt(s.TCPUt, 10)})
	}
	if s.Track != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "track", Value: s.Track})
	}
	if s.Verify != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "verify", Value: s.Verify})
	}
	if s.Verifyhost != "" {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "verifyhost", Value: s.Verifyhost})
	}
	if s.Weight != nil {
		srv.Params = append(srv.Params, &params.ServerOptionValue{Name: "weight", Value: strconv.FormatInt(*s.Weight, 10)})
	}
	return srv
}

func GetServerByName(name string, backend string, p parser.Parser) (*models.Server, int) {
	servers, err := ParseServers(backend, p)
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
