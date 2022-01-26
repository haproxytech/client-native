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

// GetServerTemplatess returns configuration version and an array of
// configured server templates in the specified backend. Returns error on fail.
func (c *Client) GetServerTemplates(backend string, transactionID string) (int64, models.ServerTemplates, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	templates, err := ParseServerTemplates(backend, p)
	if err != nil {
		return v, nil, c.HandleError("", "backend", backend, "", false, err)
	}

	return v, templates, nil
}

// GetServerTemplate returns configuration version and a requested server template
// in the specified backend. Returns error on fail or if server template does not exist.
func (c *Client) GetServerTemplate(prefix string, backend string, transactionID string) (int64, *models.ServerTemplate, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	template, _ := GetServerTemplateByPrefix(prefix, backend, p)
	if template == nil {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Server template %s does not exist in backend %s", prefix, backend))
	}

	return v, template, nil
}

// DeleteServerTemplate deletes a server template in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteServerTemplate(prefix string, backend string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	template, i := GetServerTemplateByPrefix(prefix, backend, p)
	if template == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Server template %s does not exist in backend %s", prefix, backend))
		return c.HandleError(prefix, "backend", backend, t, transactionID == "", e)
	}

	if err := p.Delete(parser.Backends, backend, "server-template", i); err != nil {
		return c.HandleError(prefix, "backend", backend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// CreateServerTemplate creates a server template in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateServerTemplate(backend string, data *models.ServerTemplate, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	template, _ := GetServerTemplateByPrefix(data.Prefix, backend, p)
	if template != nil {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("Server template %s already exists in backend %s", data.Prefix, backend))
		return c.HandleError(data.Prefix, "backend", backend, t, transactionID == "", e)
	}

	if err := p.Insert(parser.Backends, backend, "server-template", SerializeServerTemplate(*data), -1); err != nil {
		return c.HandleError(data.Prefix, "backend", backend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// EditServerTemplate edits a server template in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditServerTemplate(prefix string, backend string, data *models.ServerTemplate, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	template, i := GetServerTemplateByPrefix(prefix, backend, p)
	if template == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Server template %v does not exist in backend %s", prefix, backend))
		return c.HandleError(data.Prefix, "backend", backend, t, transactionID == "", e)
	}

	if err := p.Set(parser.Backends, backend, "server-template", SerializeServerTemplate(*data), i); err != nil {
		return c.HandleError(data.Prefix, "backend", backend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func ParseServerTemplates(backend string, p parser.Parser) (models.ServerTemplates, error) {
	templates := models.ServerTemplates{}

	data, err := p.Get(parser.Backends, backend, "server-template", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return templates, nil
		}
		return nil, err
	}

	ondiskServerTemplates := data.([]types.ServerTemplate)
	for _, ondiskServerTemplate := range ondiskServerTemplates {
		template := ParseServerTemplate(ondiskServerTemplate)
		if template != nil {
			templates = append(templates, template)
		}
	}
	return templates, nil
}

func ParseServerTemplate(ondiskServerTemplate types.ServerTemplate) *models.ServerTemplate { //nolint:gocognit,gocyclo,dupl,cyclop
	st := &models.ServerTemplate{
		Prefix:     ondiskServerTemplate.Prefix,
		NumOrRange: ondiskServerTemplate.NumOrRange,
		Fqdn:       ondiskServerTemplate.Fqdn,
		Port:       &ondiskServerTemplate.Port,
	}
	for _, p := range ondiskServerTemplate.Params { //nolint:gocognit,gocyclo,dupl,cyclop
		switch v := p.(type) {
		case *params.ServerOptionWord:
			switch v.Name {
			case "agent-check":
				st.AgentCheck = "enabled"
			case "no-agent-check":
				st.AgentCheck = "disabled"
			case "allow-0rtt":
				st.Allow0rtt = true
			case "backup":
				st.Backup = "enabled"
			case "no-backup":
				st.Backup = "disabled"
			case "check":
				st.Check = "enabled"
			case "no-check":
				st.Check = "disabled"
			case "check-send-proxy":
				st.CheckSendProxy = "enabled"
			case "check-ssl":
				st.CheckSsl = "enabled"
			case "no-check-ssl":
				st.CheckSsl = "disabled"
			case "check-via-socks4":
				st.CheckViaSocks4 = "enabled"
			case "disabled":
				st.Maintenance = "enabled"
			case "enabled":
				st.Maintenance = "disabled"
			case "force-sslv3":
				st.ForceSslv3 = "enabled"
			case "force-tlsv10":
				st.ForceTlsv10 = "enabled"
			case "no-tlsv10":
				st.ForceTlsv10 = "disabled"
			case "force-tlsv11":
				st.ForceTlsv11 = "enabled"
			case "no-tlsv11":
				st.ForceTlsv11 = "disabled"
			case "force-tlsv12":
				st.ForceTlsv12 = "enabled"
			case "no-tlsv12":
				st.ForceTlsv12 = "disabled"
			case "force-tlsv13":
				st.ForceTlsv13 = "enabled"
			case "no-tlsv13":
				st.ForceTlsv13 = "disabled"
			case "send-proxy":
				st.SendProxy = "enabled"
			case "no-send-proxy":
				st.SendProxy = "disabled"
			case "send-proxy-v2":
				st.SendProxyV2 = "enabled"
			case "no-send-proxy-v2":
				st.SendProxyV2 = "disabled"
			case "send-proxy-v2-ssl":
				st.SendProxyV2Ssl = "enabled"
			case "send-proxy-v2-ssl-cn":
				st.SendProxyV2SslCn = "enabled"
			case "ssl":
				st.Ssl = "enabled"
			case "no-ssl":
				st.Ssl = "disabled"
			case "ssl-reuse":
				st.SslReuse = "enabled"
			case "no-ssl-reuse":
				st.SslReuse = "disabled"
			case "tls-tickets":
				st.TLSTickets = "enabled"
			case "no-tls-tickets":
				st.TLSTickets = "disabled"
			case "tfo":
				st.Tfo = "enabled"
			case "no-tfo":
				st.Tfo = "disabled"
			case "stick":
				st.Stick = "enabled"
			case "no-stick":
				st.Stick = "disabled"

			}
		case *params.ServerOptionValue: //nolint:gocognit,gocyclo,dupl,cyclop
			switch v.Name {
			case "agent-send":
				st.AgentSend = v.Value
			case "agent-inter":
				st.AgentInter = misc.ParseTimeout(v.Value)
			case "agent-addr":
				st.AgentAddr = v.Value
			case "agent-port":
				p, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && p != 0 {
					st.AgentPort = &p
				}
			case "alpn":
				st.Alpn = v.Value
			case "ca-file":
				st.SslCafile = v.Value
			case "check-alpn":
				st.CheckAlpn = v.Value
			case "check-proto":
				st.CheckProto = v.Value
			case "check-sni":
				st.CheckSni = v.Value
			case "ciphers":
				st.Ciphers = v.Value
			case "ciphersuites":
				st.Ciphersuites = v.Value
			case "cookie":
				st.Cookie = v.Value
			case "crl-file":
				st.CrlFile = v.Value
			case "crt":
				st.SslCertificate = v.Value
			case "error-limit":
				c, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && c != 0 {
					st.ErrorLimit = c
				}
			case "fall":
				c, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && c != 0 {
					st.Fall = &c
				}
			case "init-addr":
				st.InitAddr = &v.Value
			case "inter":
				st.Inter = misc.ParseTimeout(v.Value)
			case "fastinter":
				st.Fastinter = misc.ParseTimeout(v.Value)
			case "downinter":
				st.Downinter = misc.ParseTimeout(v.Value)
			case "log-proto":
				st.LogProto = v.Value
			case "maxconn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && m != 0 {
					st.Maxconn = &m
				}
			case "maxqueue":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && m != 0 {
					st.Maxqueue = &m
				}
			case "max-reuse":
				c, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && c != 0 {
					st.MaxReuse = &c
				}
			case "minconn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && m != 0 {
					st.Minconn = &m
				}
			case "namespace":
				st.Namespace = v.Value
			case "npn":
				st.Npn = v.Value
			case "observe":
				st.Observe = v.Value
			case "on-error":
				st.OnError = v.Value
			case "on-marked-down":
				st.OnMarkedDown = v.Value
			case "on-marked-up":
				st.OnMarkedUp = v.Value
			case "pool-low-conn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && m != 0 {
					st.PoolLowConn = &m
				}
			case "pool-max-conn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && m != 0 {
					st.PoolMaxConn = &m
				}
			case "pool-purge-delay":
				d, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && d != 0 {
					st.PoolPurgeDelay = &d
				}
			case "port":
				p, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil {
					st.HealthCheckPort = &p
				}
			case "proto":
				st.Proto = v.Value
			case "redir":
				st.Redir = v.Value
			case "rise":
				st.Rise = misc.ParseTimeout(v.Value)
			case "resolve-opts":
				st.ResolveOpts = v.Value
			case "resolve-prefer":
				st.ResolvePrefer = v.Value
			case "resolve-net":
				st.ResolveNet = v.Value
			case "resolvers":
				st.Resolvers = v.Value
			case "proxy-v2-options":
				st.ProxyV2Options = strings.Split(v.Value, ",")
			case "slowstart":
				st.Slowstart = misc.ParseTimeout(v.Value)
			case "sni":
				st.Sni = v.Value
			case "source":
				st.Source = v.Value
			case "ssl-max-ver":
				st.SslMaxVer = v.Value
			case "ssl-min-ver":
				st.SslMinVer = v.Value
			case "socks4":
				st.Socks4 = v.Value
			case "tcp-ut":
				d, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && d != 0 {
					st.TCPUt = d
				}
			case "track":
				st.Track = v.Value
			case "verify":
				st.Verify = v.Value
			case "verifyhost":
				st.Verifyhost = v.Value
			case "weight":
				w, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && w != 0 {
					st.Weight = &w
				}
			}
		}
	}
	return st
}

func SerializeServerTemplate(s models.ServerTemplate) types.ServerTemplate { //nolint:gocognit,gocyclo,dupl,cyclop
	srv := types.ServerTemplate{
		Prefix:     s.Prefix,
		NumOrRange: s.NumOrRange,
		Fqdn:       s.Fqdn,
		Port:       *s.Port,
		Params:     []params.ServerOption{},
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

func GetServerTemplateByPrefix(prefix string, backend string, p parser.Parser) (*models.ServerTemplate, int) {
	templates, err := ParseServerTemplates(backend, p)
	if err != nil {
		return nil, 0
	}
	for i, template := range templates {
		if template.Prefix == prefix {
			return template, i
		}
	}
	return nil, 0
}
