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

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parsererrors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"

	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type LogForward interface {
	GetLogForwards(transactionID string) (int64, models.LogForwards, error)
	GetLogForward(name string, transactionID string) (int64, *models.LogForward, error)
	DeleteLogForward(name string, transactionID string, version int64) error
	CreateLogForward(data *models.LogForward, transactionID string, version int64) error
	EditLogForward(name string, data *models.LogForward, transactionID string, version int64) error
}

// GetLogForwards returns configuration version and an array of
// configured log forwards. Returns error on fail.
func (c *client) GetLogForwards(transactionID string) (int64, models.LogForwards, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	lNames, err := p.SectionsGet(parser.LogForward)
	if err != nil {
		return v, nil, err
	}

	ls := []*models.LogForward{}
	for _, name := range lNames {
		_, l, err := c.GetLogForward(name, transactionID)
		if err == nil {
			ls = append(ls, l)
		}
	}

	return v, ls, nil
}

// GetLogForward returns configuration version and a requested log forward.
// Returns error on fail or if log forward does not exist.
func (c *client) GetLogForward(name string, transactionID string) (int64, *models.LogForward, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.LogForward, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("log forward %s does not exist", name))
	}

	lf := &models.LogForward{LogForwardBase: models.LogForwardBase{Name: name}}
	if err = ParseLogForward(p, lf); err != nil {
		return 0, nil, err
	}
	return v, lf, nil
}

func ParseLogForward(p parser.Parser, lf *models.LogForward) error {
	backlog, err := p.Get(parser.LogForward, lf.Name, "backlog", false)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		bl, ok := backlog.(*types.Int64C)
		if !ok {
			return misc.CreateTypeAssertError("backlog")
		}
		if bl.Value > 0 {
			lf.Backlog = &bl.Value
		}
	}

	maxconn, err := p.Get(parser.LogForward, lf.Name, "maxconn", false)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		mc, ok := maxconn.(*types.Int64C)
		if !ok {
			return misc.CreateTypeAssertError("maxconn")
		}
		if mc.Value > 0 {
			lf.Maxconn = &mc.Value
		}
	}

	// timeouts
	tConnect, err := p.Get(parser.LogForward, lf.Name, "timeout client", false)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		tc, ok := tConnect.(*types.SimpleTimeout)
		if !ok {
			return misc.CreateTypeAssertError("timeout client")
		}
		if tc.Value != "" {
			lf.TimeoutClient = misc.ParseTimeout(tc.Value)
		}
	}
	return nil
}

// DeleteLogForward deletes a log forward in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteLogForward(name string, transactionID string, version int64) error {
	return c.deleteSection(parser.LogForward, name, transactionID, version)
}

// CreateLogForward creates a log forward in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateLogForward(data *models.LogForward, transactionID string, version int64) error {
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

	if c.checkSectionExists(parser.LogForward, data.Name, p) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.LogForward, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(parser.LogForward, data.Name); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if err = SerializeLogForwardSection(p, data, &c.ConfigurationOptions); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditLogForward edits a log forward in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditLogForward(name string, data *models.LogForward, transactionID string, version int64) error { //nolint:revive
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

	if !c.checkSectionExists(parser.LogForward, data.Name, p) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s does not exists", parser.LogForward, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = SerializeLogForwardSection(p, data, &c.ConfigurationOptions); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

func SerializeLogForwardSection(p parser.Parser, data *models.LogForward, opt *options.ConfigurationOptions) error {
	if data == nil {
		return errors.New("empty log forward")
	}

	var err error

	if data.Backlog == nil {
		if err = p.Set(parser.LogForward, data.Name, "backlog", nil); err != nil {
			return err
		}
	} else {
		d := types.Int64C{Value: *data.Backlog}
		if err = p.Set(parser.LogForward, data.Name, "backlog", d); err != nil {
			return err
		}
	}

	if data.Maxconn == nil {
		if err = p.Set(parser.LogForward, data.Name, "maxconn", nil); err != nil {
			return err
		}
	} else {
		d := types.Int64C{Value: *data.Maxconn}
		if err = p.Set(parser.LogForward, data.Name, "maxconn", d); err != nil {
			return err
		}
	}

	if data.TimeoutClient == nil {
		if err = p.Set(parser.LogForward, data.Name, "timeout client", nil); err != nil {
			return err
		}
	} else {
		tc := types.SimpleTimeout{Value: misc.SerializeTime(*data.TimeoutClient, opt.PreferredTimeSuffix)}
		if err = p.Set(parser.LogForward, data.Name, "timeout client", tc); err != nil {
			return err
		}
	}

	return err
}
