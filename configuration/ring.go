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

type Ring interface {
	GetRings(transactionID string) (int64, models.Rings, error)
	GetRing(name string, transactionID string) (int64, *models.Ring, error)
	DeleteRing(name string, transactionID string, version int64) error
	CreateRing(data *models.Ring, transactionID string, version int64) error
	EditRing(name string, data *models.Ring, transactionID string, version int64) error
}

// GetRings returns configuration version and an array of
// configured rings. Returns error on fail.
func (c *client) GetRings(transactionID string) (int64, models.Rings, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	rNames, err := p.SectionsGet(parser.Ring)
	if err != nil {
		return v, nil, err
	}

	rings := []*models.Ring{}
	for _, name := range rNames {
		_, a, err := c.GetRing(name, transactionID)
		if err == nil {
			rings = append(rings, a)
		}
	}

	return v, rings, nil
}

// GetRing returns configuration version and a requested ring.
// Returns error on fail or if ring does not exist.
func (c *client) GetRing(name string, transactionID string) (int64, *models.Ring, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !p.SectionExists(parser.Ring, name) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("ring %s does not exist", name))
	}

	ring := &models.Ring{RingBase: models.RingBase{Name: name}}
	if err = ParseRingSection(p, ring); err != nil {
		return 0, nil, err
	}
	return v, ring, nil
}

func ParseRingSection(p parser.Parser, ring *models.Ring) error { //nolint:gocognit
	if data, err := p.SectionGet(parser.Ring, ring.Name); err == nil {
		d, ok := data.(types.Section)
		if ok {
			ring.Metadata = misc.ParseMetadata(d.Comment)
		}
	}
	description, err := p.Get(parser.Ring, ring.Name, "description", false)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		desc, ok := description.(*types.StringC)
		if !ok {
			return misc.CreateTypeAssertError("description")
		}

		if desc.Value != "" {
			ring.Description = desc.Value
		}
	}

	format, err := p.Get(parser.Ring, ring.Name, "format", true)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		f, ok := format.(*types.StringC)
		if !ok {
			return misc.CreateTypeAssertError("format")
		}
		if f.Value != "" {
			ring.Format = f.Value
		}
	}

	maxlen, err := p.Get(parser.Ring, ring.Name, "maxlen", true)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		mx, ok := maxlen.(*types.Int64C)
		if !ok {
			return misc.CreateTypeAssertError("maxlen")
		}
		if mx.Value > 0 {
			ring.Maxlen = &mx.Value
		}
	}

	size, err := p.Get(parser.Ring, ring.Name, "size", true)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		sz, ok := size.(*types.StringC)
		if !ok {
			return misc.CreateTypeAssertError("size")
		}
		if sz.Value != "" {
			ring.Size = misc.ParseSize(sz.Value)
		}
	}

	// timeouts
	tConnect, err := p.Get(parser.Ring, ring.Name, "timeout connect", true)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		tc, ok := tConnect.(*types.SimpleTimeout)
		if !ok {
			return misc.CreateTypeAssertError("timeout connect")
		}
		if tc.Value != "" {
			ring.TimeoutConnect = misc.ParseTimeout(tc.Value)
		}
	}

	tServer, err := p.Get(parser.Ring, ring.Name, "timeout server", true)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		ts, ok := tServer.(*types.SimpleTimeout)
		if !ok {
			return misc.CreateTypeAssertError("timeout server")
		}
		if ts.Value != "" {
			ring.TimeoutServer = misc.ParseTimeout(ts.Value)
		}
	}
	return nil
}

// DeleteRing deletes a ring in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteRing(name string, transactionID string, version int64) error {
	return c.deleteSection(parser.Ring, name, transactionID, version)
}

// CreateRing creates a ring in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateRing(data *models.Ring, transactionID string, version int64) error {
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

	if p.SectionExists(parser.Ring, data.Name) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.Ring, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(parser.Ring, data.Name); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if err = SerializeRingSection(p, data, &c.ConfigurationOptions); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditRing edits a ring in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditRing(name string, data *models.Ring, transactionID string, version int64) error { //nolint:revive
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

	if !p.SectionExists(parser.Ring, data.Name) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s does not exists", parser.Ring, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = SerializeRingSection(p, data, &c.ConfigurationOptions); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

func SerializeRingSection(p parser.Parser, data *models.Ring, opt *options.ConfigurationOptions) error { //nolint:gocognit
	if data == nil {
		return errors.New("empty ring")
	}
	if data.Metadata != nil {
		comment, err := misc.SerializeMetadata(data.Metadata)
		if err != nil {
			return err
		}
		if err := p.SectionCommentSet(parser.Ring, data.Name, comment); err != nil {
			return err
		}
	}
	var err error
	if data.Description == "" {
		if err = p.Set(parser.Ring, data.Name, "description", nil); err != nil {
			return err
		}
	} else {
		d := types.StringC{Value: data.Description}
		if err = p.Set(parser.Ring, data.Name, "description", d); err != nil {
			return err
		}
	}

	if data.Format == "" {
		if err = p.Set(parser.Ring, data.Name, "format", nil); err != nil {
			return err
		}
	} else {
		d := types.StringC{Value: data.Format}
		if err = p.Set(parser.Ring, data.Name, "format", d); err != nil {
			return err
		}
	}

	if data.Maxlen == nil {
		if err = p.Set(parser.Ring, data.Name, "maxlen", nil); err != nil {
			return err
		}
	} else {
		d := types.Int64C{Value: *data.Maxlen}
		if err = p.Set(parser.Ring, data.Name, "maxlen", d); err != nil {
			return err
		}
	}

	if data.Size == nil {
		if err = p.Set(parser.Ring, data.Name, "size", nil); err != nil {
			return err
		}
	} else {
		d := types.StringC{Value: misc.SerializeSize(*data.Size)}
		if err = p.Set(parser.Ring, data.Name, "size", d); err != nil {
			return err
		}
	}

	if data.TimeoutConnect == nil {
		if err = p.Set(parser.Ring, data.Name, "timeout connect", nil); err != nil {
			return err
		}
	} else {
		tc := types.SimpleTimeout{Value: misc.SerializeTime(*data.TimeoutConnect, opt.PreferredTimeSuffix)}
		if err = p.Set(parser.Ring, data.Name, "timeout connect", tc); err != nil {
			return err
		}
	}

	if data.TimeoutServer == nil {
		if err = p.Set(parser.Ring, data.Name, "timeout server", nil); err != nil {
			return err
		}
	} else {
		ts := types.SimpleTimeout{Value: misc.SerializeTime(*data.TimeoutServer, opt.PreferredTimeSuffix)}
		if err = p.Set(parser.Ring, data.Name, "timeout server", ts); err != nil {
			return err
		}
	}

	return err
}
