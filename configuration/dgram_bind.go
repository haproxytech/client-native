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
	parser "github.com/haproxytech/config-parser/v5"
	parser_errors "github.com/haproxytech/config-parser/v5/errors"
	"github.com/haproxytech/config-parser/v5/params"
	"github.com/haproxytech/config-parser/v5/types"

	"github.com/haproxytech/client-native/v6/models"
)

type DgramBind interface {
	GetDgramBinds(logForward string, transactionID string) (int64, models.DgramBinds, error)
	GetDgramBind(name string, logForward string, transactionID string) (int64, *models.DgramBind, error)
	DeleteDgramBind(name string, logForward string, transactionID string, version int64) error
	CreateDgramBind(logForward string, data *models.DgramBind, transactionID string, version int64) error
	EditDgramBind(name string, logForward string, data *models.DgramBind, transactionID string, version int64) error
}

// GetDgramBinds returns configuration version and an array of
// configured binds in the specified logForward. Returns error on fail.
func (c *client) GetDgramBinds(logForward string, transactionID string) (int64, models.DgramBinds, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	binds, err := ParseDgramBinds(logForward, p)
	if err != nil {
		return v, nil, c.HandleError("", "log-forward", logForward, "", false, err)
	}

	return v, binds, nil
}

// GetDgramBind returns configuration version and a requested dgram-bind
// in the specified logForward. Returns error on fail or if dgram-bind does not exist.
func (c *client) GetDgramBind(name string, logForward string, transactionID string) (int64, *models.DgramBind, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	dBind, _ := GetDgramBindByName(name, logForward, p)
	if dBind == nil {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("DgramBind %s does not exist in log-forward %s", name, logForward))
	}

	return v, dBind, nil
}

// DeleteDgramBind deletes a dgram-bind in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteDgramBind(name string, logForward string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	dBind, i := GetDgramBindByName(name, logForward, p)
	if dBind == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("DgramBind %s does not exist in log-forward %s", name, logForward))
		return c.HandleError(name, "log-forward", logForward, t, transactionID == "", e)
	}

	if err := p.Delete(parser.LogForward, logForward, "dgram-bind", i); err != nil {
		return c.HandleError(name, "log-forward", logForward, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateDgramBind creates a dgram-bind in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateDgramBind(logForward string, data *models.DgramBind, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
		validationErr = validateDgramBindParams(data)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	dBind, _ := GetDgramBindByName(data.Name, logForward, p)
	if dBind != nil {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("DgramBind %s already exists in log-forward %s", data.Name, logForward))
		return c.HandleError(data.Name, "log-forward", logForward, t, transactionID == "", e)
	}

	if err := p.Insert(parser.LogForward, logForward, "dgram-bind", SerializeDgramBind(*data), -1); err != nil {
		return c.HandleError(data.Name, "log-forward", logForward, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditDgramBind edits a dgram-bind in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditDgramBind(name string, logForward string, data *models.DgramBind, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
		validationErr = validateDgramBindParams(data)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	dBind, i := GetDgramBindByName(name, logForward, p)
	if dBind == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("DgramBind %v does not exist in log-forward %s", name, logForward))
		return c.HandleError(data.Name, "log-forward", logForward, t, transactionID == "", e)
	}

	if err := p.Set(parser.LogForward, logForward, "dgram-bind", SerializeDgramBind(*data), i); err != nil {
		return c.HandleError(data.Name, "log-forward", logForward, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseDgramBinds(logForward string, p parser.Parser) (models.DgramBinds, error) {
	var dBinds models.DgramBinds

	data, err := p.Get(parser.LogForward, logForward, "dgram-bind", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return dBinds, nil
		}
		return nil, err
	}

	ondiskBinds := data.([]types.DgramBind) //nolint:forcetypeassert
	for _, ondiskBind := range ondiskBinds {
		b := ParseDgramBind(ondiskBind)
		if b != nil {
			dBinds = append(dBinds, b)
		}
	}
	return dBinds, nil
}

func ParseDgramBind(ondiskDgramBind types.DgramBind) *models.DgramBind {
	b := &models.DgramBind{}
	if strings.HasPrefix(ondiskDgramBind.Path, "/") {
		b.Address = ondiskDgramBind.Path
	} else {
		addSlice := strings.Split(ondiskDgramBind.Path, ":")
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
	for _, p := range ondiskDgramBind.Params {
		switch v := p.(type) {
		case *params.BindOptionWord:
			if v.Name == "transparent" {
				b.Transparent = true
			}
		case *params.BindOptionValue:
			switch v.Name {
			case "name":
				b.Name = v.Value
			case "interface":
				b.Interface = v.Value
			case "namespace":
				b.Namespace = v.Value
			}
		}
	}
	if b.Name == "" {
		b.Name = ondiskDgramBind.Path
	}
	return b
}

func SerializeDgramBind(b models.DgramBind) types.DgramBind {
	dBind := types.DgramBind{
		Params: []params.DgramBindOption{},
	}
	if b.Port != nil {
		dBind.Path = b.Address + ":" + strconv.FormatInt(*b.Port, 10)
		if b.PortRangeEnd != nil {
			dBind.Path = dBind.Path + "-" + strconv.FormatInt(*b.PortRangeEnd, 10)
		}
	} else {
		dBind.Path = b.Address
	}

	dBind.Params = []params.DgramBindOption{}
	if b.Name != "" {
		dBind.Params = append(dBind.Params, &params.BindOptionValue{Name: "name", Value: b.Name})
	} else if dBind.Path != "" {
		dBind.Params = append(dBind.Params, &params.BindOptionValue{Name: "name", Value: dBind.Path})
	}

	if b.Transparent {
		dBind.Params = append(dBind.Params, &params.BindOptionWord{Name: "transparent"})
	}

	if b.Interface != "" {
		dBind.Params = append(dBind.Params, &params.BindOptionValue{Name: "interface", Value: b.Interface})
	}

	if b.Namespace != "" {
		dBind.Params = append(dBind.Params, &params.BindOptionValue{Name: "namespace", Value: b.Namespace})
	}

	return dBind
}

func GetDgramBindByName(name string, logForward string, p parser.Parser) (*models.DgramBind, int) {
	dBinds, err := ParseDgramBinds(logForward, p)
	if err != nil {
		return nil, 0
	}

	for i, b := range dBinds {
		if b.Name == name {
			return b, i
		}
	}
	return nil, 0
}

func validateDgramBindParams(data *models.DgramBind) error {
	if data.Port != nil {
		if data.PortRangeEnd != nil && *data.PortRangeEnd <= *data.Port {
			return fmt.Errorf("port upper bound %d less or equal than lower bound %d in dgram-bind %s", *data.PortRangeEnd, *data.Port, data.Name)
		}
		return nil
	}
	if data.Address != "" {
		return nil
	}
	return fmt.Errorf("missing port or address in dgram-bind %s", data.Name)
}
