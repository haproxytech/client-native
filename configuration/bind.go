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
	"fmt"
	"strconv"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/params"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
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

	binds, err := c.parseBinds(frontend, p)
	if err != nil {
		return v, nil, c.handleError("", "frontend", frontend, "", false, err)
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

	bind, _ := c.getBindByName(name, frontend, p)
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

	bind, i := c.getBindByName(name, frontend, p)
	if bind == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Bind %s does not exist in frontend %s", name, frontend))
		return c.handleError(name, "frontend", frontend, t, transactionID == "", e)
	}

	if err := p.Delete(parser.Frontends, frontend, "bind", i); err != nil {
		return c.handleError(name, "frontend", frontend, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
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

	bind, _ := c.getBindByName(data.Name, frontend, p)
	if bind != nil {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("Bind %s already exists in frontend %s", data.Name, frontend))
		return c.handleError(data.Name, "frontend", frontend, t, transactionID == "", e)
	}

	if err := p.Insert(parser.Frontends, frontend, "bind", serializeBind(*data), -1); err != nil {
		return c.handleError(data.Name, "frontend", frontend, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
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

	bind, i := c.getBindByName(name, frontend, p)
	if bind == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Bind %v does not exist in frontend %s", name, frontend))
		return c.handleError(data.Name, "frontend", frontend, t, transactionID == "", e)
	}

	if err := p.Set(parser.Frontends, frontend, "bind", serializeBind(*data), i); err != nil {
		return c.handleError(data.Name, "frontend", frontend, t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func (c *Client) parseBinds(frontend string, p *parser.Parser) (models.Binds, error) {
	binds := models.Binds{}

	data, err := p.Get(parser.Frontends, frontend, "bind", false)
	if err != nil {
		if err == parser_errors.ErrFetch {
			return binds, nil
		}
		return nil, err
	}

	ondiskBinds := data.([]types.Bind)
	for _, ondiskBind := range ondiskBinds {
		b := parseBind(ondiskBind)
		if b != nil {
			binds = append(binds, b)
		}
	}
	return binds, nil
}

func parseBind(ondiskBind types.Bind) *models.Bind {
	b := &models.Bind{
		Name: ondiskBind.Path,
	}
	if strings.HasPrefix(ondiskBind.Path, "/") {
		b.Address = ondiskBind.Path
	} else {
		addSlice := strings.Split(ondiskBind.Path, ":")
		if len(addSlice) == 0 {
			return nil
		} else if len(addSlice) == 4 { // :::443
			b.Address = "::"
			if addSlice[3] != "" {
				p, err := strconv.ParseInt(addSlice[3], 10, 64)
				if err == nil {
					b.Port = &p
				}
			}
		} else if len(addSlice) > 1 {
			b.Address = addSlice[0]
			if addSlice[1] != "" {
				p, err := strconv.ParseInt(addSlice[1], 10, 64)
				if err == nil {
					b.Port = &p
				}
			}
		} else if len(addSlice) > 0 {
			b.Address = addSlice[0]
		}
	}
	for _, p := range ondiskBind.Params {
		switch v := p.(type) {
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
			}
		}
	}
	return b
}

func serializeBind(b models.Bind) types.Bind {
	bind := types.Bind{
		Params: []params.BindOption{},
	}
	if b.Port != nil {
		bind.Path = b.Address + ":" + strconv.FormatInt(*b.Port, 10)
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

	return bind
}

func (c *Client) getBindByName(name string, frontend string, p *parser.Parser) (*models.Bind, int) {
	binds, err := c.parseBinds(frontend, p)
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
