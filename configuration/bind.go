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

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v3"
	parser_errors "github.com/haproxytech/config-parser/v3/errors"
	"github.com/haproxytech/config-parser/v3/params"
	"github.com/haproxytech/config-parser/v3/types"

	"github.com/haproxytech/client-native/v2/models"
)

// GetBinds returns configuration version and an array of
// configured binds in the specified frontend. Returns error on fail.
func (c *Client) GetBinds(frontend string, transactionID string) (int64, models.Binds, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	binds, err := ParseBinds(frontend, p)
	if err != nil {
		return v, nil, c.HandleError("", "frontend", frontend, "", false, err)
	}

	return v, binds, nil
}

// GetBind returns configuration version and a requested bind
// in the specified frontend. Returns error on fail or if bind does not exist.
func (c *Client) GetBind(name string, frontend string, transactionID string) (int64, *models.Bind, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	bind, _ := GetBindByName(name, frontend, p)
	if bind == nil {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Bind %s does not exist in frontend %s", name, frontend))
	}

	return v, bind, nil
}

// DeleteBind deletes a bind in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteBind(name string, frontend string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	bind, i := GetBindByName(name, frontend, p)
	if bind == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Bind %s does not exist in frontend %s", name, frontend))
		return c.HandleError(name, "frontend", frontend, t, transactionID == "", e)
	}

	if err := p.Delete(parser.Frontends, frontend, "bind", i); err != nil {
		return c.HandleError(name, "frontend", frontend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// CreateBind creates a bind in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateBind(frontend string, data *models.Bind, transactionID string, version int64) error {
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

	if data.PortRangeEnd != nil && *data.Port >= *data.PortRangeEnd {
		e := NewConfError(ErrGeneralError, fmt.Sprintf("Bind port range end %d has to be greater start %d", *data.PortRangeEnd, *data.Port))
		return c.HandleError(data.Name, "frontend", frontend, t, transactionID == "", e)
	}

	bind, _ := GetBindByName(data.Name, frontend, p)
	if bind != nil {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("Bind %s already exists in frontend %s", data.Name, frontend))
		return c.HandleError(data.Name, "frontend", frontend, t, transactionID == "", e)
	}

	if err := p.Insert(parser.Frontends, frontend, "bind", SerializeBind(*data), -1); err != nil {
		return c.HandleError(data.Name, "frontend", frontend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// EditBind edits a bind in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditBind(name string, frontend string, data *models.Bind, transactionID string, version int64) error {
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

	bind, i := GetBindByName(name, frontend, p)
	if bind == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Bind %v does not exist in frontend %s", name, frontend))
		return c.HandleError(data.Name, "frontend", frontend, t, transactionID == "", e)
	}

	if err := p.Set(parser.Frontends, frontend, "bind", SerializeBind(*data), i); err != nil {
		return c.HandleError(data.Name, "frontend", frontend, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func ParseBinds(frontend string, p *parser.Parser) (models.Binds, error) {
	binds := models.Binds{}

	data, err := p.Get(parser.Frontends, frontend, "bind", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return binds, nil
		}
		return nil, err
	}

	ondiskBinds := data.([]types.Bind)
	for _, ondiskBind := range ondiskBinds {
		b := ParseBind(ondiskBind)
		if b != nil {
			binds = append(binds, b)
		}
	}
	return binds, nil
}

func ParseBind(ondiskBind types.Bind) *models.Bind { //nolint:gocognit,gocyclo
	b := &models.Bind{
		Name: ondiskBind.Path,
	}
	if strings.HasPrefix(ondiskBind.Path, "/") {
		b.Address = ondiskBind.Path
	} else {
		addSlice := strings.Split(ondiskBind.Path, ":")
		switch n := len(addSlice); {
		case n == 0:
			return nil
		case n == 4: // :::443
			b.Address = "::"
			if addSlice[3] != "" {
				p, err := strconv.ParseInt(addSlice[3], 10, 64)
				if err == nil {
					b.Port = &p
				}
			}
		case n > 1:
			b.Address = addSlice[0]
			ports := strings.Split(addSlice[1], "-")

			// *:<port>
			if ports[0] != "" {
				port, err := strconv.ParseInt(ports[0], 10, 64)
				if err == nil {
					b.Port = &port
				}
			}
			// *:<port-first>-<port-last>
			if b.Port != nil && len(ports) == 2 {
				portRangeEnd, err := strconv.ParseInt(ports[1], 10, 64)
				// Deny inverted interval.
				if err == nil && (*b.Port < portRangeEnd) {
					b.PortRangeEnd = &portRangeEnd
				}
			}
		case n > 0:
			b.Address = addSlice[0]

		}
	}
	for _, p := range ondiskBind.Params {
		switch v := p.(type) {
		case *params.BindOptionDoubleWord:
			if v.Name == "expose-fd" && v.Value == "listener" {
				b.ExposeFdListeners = true
			}
		case *params.BindOptionWord:
			switch v.Name {
			case "ssl":
				b.Ssl = true
			case "transparent":
				b.Transparent = true
			case "accept-proxy":
				b.AcceptProxy = true
			case "v4v6":
				b.V4v6 = true
			case "allow-0rtt":
				b.Allow0rtt = true
			case "defer-accept":
				b.DeferAccept = true
			case "force-sslv3":
				b.ForceSslv3 = true
			case "force-tlsv10":
				b.ForceTlsv10 = true
			case "force-tlsv11":
				b.ForceTlsv11 = true
			case "force-tlsv12":
				b.ForceTlsv12 = true
			case "force-tlsv13":
				b.ForceTlsv13 = true
			case "generate-certificates":
				b.GenerateCertificates = true
			case "no-ca-names":
				b.NoCaNames = true
			case "no-sslv3":
				b.NoSslv3 = true
			case "no-tls-tickets":
				b.NoTLSTickets = true
			case "no-tlsv10":
				b.NoTlsv10 = true
			case "no-tlsv11":
				b.NoTlsv11 = true
			case "no-tlsv12":
				b.NoTlsv12 = true
			case "no-tlsv13":
				b.NoTlsv13 = true
			case "prefer-client-ciphers":
				b.PreferClientCiphers = true
			case "strict-sni":
				b.StrictSni = true
			case "tfo":
				b.Tfo = true
			case "v6only":
				b.V6only = true
			}
		case *params.BindOptionValue:
			switch v.Name {
			case "name":
				b.Name = v.Value
			case "process":
				b.Process = v.Value
			case "tcp-ut":
				t, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && t != 0 {
					b.TCPUserTimeout = &t
				}
			case "crt":
				b.SslCertificate = v.Value
			case "ca-file":
				b.SslCafile = v.Value
			case "verify":
				b.Verify = v.Value
			case "alpn":
				b.Alpn = v.Value
			case "accept-netscaler-cip":
				mn, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && mn != 0 {
					b.AcceptNetscalerCip = mn
				}
			case "backlog":
				b.Backlog = v.Value
			case "curves":
				b.Curves = v.Value
			case "ecdhe":
				b.Ecdhe = v.Value
			case "ca-ignore-err":
				b.CaIgnoreErr = v.Value
			case "ca-sign-file":
				b.CaSignFile = v.Value
			case "ca-sign-pass":
				b.CaSignPass = v.Value
			case "ciphers":
				b.Ciphers = v.Value
			case "ciphersuites":
				b.Ciphersuites = v.Value
			case "crl-file":
				b.CrlFile = v.Value
			case "crt-ignore-err":
				b.CrtIgnoreErr = v.Value
			case "crt-list":
				b.CrtList = v.Value
			case "gid":
				gid, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && gid != 0 {
					b.Gid = gid
				}
			case "group":
				b.Group = v.Value
			case "id":
				b.ID = v.Value
			case "interface":
				b.Interface = v.Value
			case "level":
				b.Level = v.Value
			case "severity-output":
				b.SeverityOutput = v.Value
			case "maxconn":
				m, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && m != 0 {
					b.Maxconn = m
				}
			case "mode":
				b.Mode = v.Value
			case "mss":
				b.Mss = v.Value
			case "namespace":
				b.Namespace = v.Value
			case "nice":
				n, err := strconv.ParseInt(v.Value, 10, 64)
				if err == nil && n != 0 {
					b.Nice = n
				}
			case "npn":
				b.Npn = v.Value
			case "proto":
				b.Proto = v.Value
			case "ssl-max-ver":
				b.SslMaxVer = v.Value
			case "ssl-min-ver":
				b.SslMinVer = v.Value
			case "tls-ticket-keys":
				b.TLSTicketKeys = v.Value
			case "uid":
				b.UID = v.Value
			case "user":
				b.User = v.Value
			}
		}
	}
	return b
}

func SerializeBind(b models.Bind) types.Bind { //nolint:gocognit,gocyclo
	bind := types.Bind{
		Params: []params.BindOption{},
	}
	if b.Port != nil {
		bind.Path = b.Address + ":" + strconv.FormatInt(*b.Port, 10)
		if b.PortRangeEnd != nil {
			bind.Path = bind.Path + "-" + strconv.FormatInt(*b.PortRangeEnd, 10)
		}
	} else {
		bind.Path = b.Address
	}
	if b.Name != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "name", Value: b.Name})
	} else {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "name", Value: bind.Path})
	}
	if b.Process != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "process", Value: b.Process})
	}
	if b.SslCertificate != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "crt", Value: b.SslCertificate})
	}
	if b.SslCafile != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "ca-file", Value: b.SslCafile})
	}
	if b.TCPUserTimeout != nil {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "tcp-ut", Value: strconv.FormatInt(*b.TCPUserTimeout, 10)})
	}
	if b.Ssl {
		bind.Params = append(bind.Params, &params.BindOptionWord{Name: "ssl"})
	}
	if b.V4v6 {
		bind.Params = append(bind.Params, &params.BindOptionWord{Name: "v4v6"})
	}
	if b.Verify != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "verify", Value: b.Verify})
	}
	if b.Alpn != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "alpn", Value: b.Alpn})
	}
	if b.Transparent {
		bind.Params = append(bind.Params, &params.BindOptionWord{Name: "transparent"})
	}
	if b.AcceptProxy {
		bind.Params = append(bind.Params, &params.BindOptionWord{Name: "accept-proxy"})
	}
	if b.Allow0rtt {
		bind.Params = append(bind.Params, &params.BindOptionWord{Name: "allow-0rtt"})
	}
	if b.AcceptNetscalerCip != 0 {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "accept-netscaler-cip", Value: strconv.FormatInt(b.AcceptNetscalerCip, 10)})
	}
	if b.Backlog != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "backlog", Value: b.Backlog})
	}
	if b.Curves != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "curves", Value: b.Curves})
	}
	if b.Ecdhe != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "ecdhe", Value: b.Ecdhe})
	}
	if b.CaIgnoreErr != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "ca-ignore-err", Value: b.CaIgnoreErr})
	}
	if b.CaSignFile != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "ca-sign-file", Value: b.CaSignFile})
	}
	if b.CaSignPass != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "ca-sign-pass", Value: b.CaSignPass})
	}
	if b.Ciphers != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "ciphers ", Value: b.Ciphers})
	}
	if b.Ciphersuites != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "ciphersuites ", Value: b.Ciphersuites})
	}
	if b.CrlFile != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "crl-file ", Value: b.CrlFile})
	}
	if b.CrtIgnoreErr != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "crt-ignore-err", Value: b.CrtIgnoreErr})
	}
	if b.CrtList != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "crt-list", Value: b.CrtList})
	}
	if b.DeferAccept {
		bind.Params = append(bind.Params, &params.BindOptionWord{Name: "defer-accept"})
	}
	if b.ExposeFdListeners {
		bind.Params = append(bind.Params, &params.ServerOptionDoubleWord{Name: "expose-fd", Value: "listeners"})
	}
	if b.ForceSslv3 {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "force-sslv3"})
	}
	if b.ForceTlsv10 {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "force-tlsv10"})
	}
	if b.ForceTlsv11 {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "force-tlsv11"})
	}
	if b.ForceTlsv12 {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "force-tlsv12"})
	}
	if b.ForceTlsv13 {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "force-tlsv13"})
	}
	if b.GenerateCertificates {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "generate-certificates"})
	}
	if b.Gid != 0 {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "gid", Value: strconv.FormatInt(b.Gid, 10)})
	}
	if b.Group != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "group", Value: b.Group})
	}
	if b.ID != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "id", Value: b.ID})
	}
	if b.Interface != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "interface", Value: b.Interface})
	}
	if b.Level != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "level", Value: b.Level})
	}
	if b.SeverityOutput != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "severity-output", Value: b.SeverityOutput})
	}
	if b.Maxconn != 0 {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "maxconn", Value: strconv.FormatInt(b.Maxconn, 10)})
	}
	if b.Mode != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "mode", Value: b.Mode})
	}
	if b.Mss != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "mss", Value: b.Mss})
	}
	if b.Namespace != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "namespace", Value: b.Namespace})
	}
	if b.NoCaNames {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "no-ca-names"})
	}
	if b.NoSslv3 {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "no-sslv3"})
	}
	if b.NoTLSTickets {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "no-tls-tickets"})
	}
	if b.NoTlsv10 {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "no-tlsv10"})
	}
	if b.NoTlsv11 {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "no-tlsv11"})
	}
	if b.NoTlsv12 {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "no-tlsv12"})
	}
	if b.NoTlsv13 {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "no-tlsv13"})
	}
	if b.Npn != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "npn", Value: b.Npn})
	}
	if b.PreferClientCiphers {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "prefer-client-ciphers"})
	}
	if b.Proto != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "proto", Value: b.Proto})
	}
	if b.SslMaxVer != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "ssl-max-ver", Value: b.SslMaxVer})
	}
	if b.SslMinVer != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "ssl-min-ver", Value: b.SslMinVer})
	}
	if b.StrictSni {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "strict-sni"})
	}
	if b.Tfo {
		bind.Params = append(bind.Params, &params.ServerOptionWord{Name: "tfo"})
	}
	if b.TLSTicketKeys != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "tls-ticket-keys", Value: b.TLSTicketKeys})
	}
	if b.V6only {
		bind.Params = append(bind.Params, &params.BindOptionWord{Name: "v6only"})
	}
	if b.UID != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "uid", Value: b.UID})
	}
	if b.User != "" {
		bind.Params = append(bind.Params, &params.BindOptionValue{Name: "user", Value: b.User})
	}

	return bind
}

func GetBindByName(name string, frontend string, p *parser.Parser) (*models.Bind, int) {
	binds, err := ParseBinds(frontend, p)
	if err != nil {
		return nil, 0
	}

	for i, b := range binds {
		if b.Name == name {
			return b, i
		}
	}
	return nil, 0
}
