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

package configuration

import (
	"errors"
	"fmt"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parsererrors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type FCGIApp interface {
	GetFCGIApplications(transactionID string) (int64, models.FCGIApps, error)
	GetFCGIApplication(name string, transactionID string) (int64, *models.FCGIApp, error)
	DeleteFCGIApplication(name string, transactionID string, version int64) error
	EditFCGIApplication(name string, data *models.FCGIApp, transactionID string, version int64) error
	CreateFCGIApplication(data *models.FCGIApp, transactionID string, version int64) error
}

func maxReqs(p parser.Parser, name string) (*types.OptionMaxReqs, error) {
	data, err := p.Get(parser.FCGIApp, name, "option max-reqs")
	if err != nil {
		return nil, err
	}

	d := data.(*types.OptionMaxReqs)

	return d, nil
}

func genericString(p parser.Parser, name, attribute string) (*types.StringC, error) {
	data, err := p.Get(parser.FCGIApp, name, attribute)
	if err != nil {
		if errors.Is(err, parsererrors.ErrFetch) {
			return &types.StringC{}, nil
		}

		return nil, err
	}

	return data.(*types.StringC), nil
}

func logStderr(p parser.Parser, name string) ([]*models.FCGILogStderr, error) {
	data, err := p.Get(parser.FCGIApp, name, "log-stderr")
	if err != nil {
		if errors.Is(err, parsererrors.ErrFetch) {
			return []*models.FCGILogStderr{}, nil
		}

		return nil, err
	}

	d := data.([]types.LogStdErr)

	logstderr := make([]*models.FCGILogStderr, 0, len(d))

	for _, log := range d {
		if log.Global {
			logstderr = append(logstderr, &models.FCGILogStderr{
				Global: true,
			})

			continue
		}

		l := models.FCGILogStderr{
			Address:  log.Address,
			Facility: log.Facility,
			Format:   log.Format,
			Global:   log.Global,
			Len:      log.Length,
			Level:    log.Level,
			Minlevel: log.MinLevel,
		}

		if len(log.SampleRange) > 0 && log.SampleSize > 0 {
			l.Sample = &models.FCGILogStderrSample{
				Ranges: misc.StringP(log.SampleRange),
				Size:   misc.Int64P(int(log.SampleSize)),
			}
		}

		logstderr = append(logstderr, &l)
	}

	return logstderr, nil
}

func setParam(p parser.Parser, name string) ([]*models.FCGISetParam, error) {
	data, err := p.Get(parser.FCGIApp, name, "set-param")
	if err != nil {
		if errors.Is(err, parsererrors.ErrFetch) {
			return []*models.FCGISetParam{}, nil
		}

		return nil, err
	}

	d := data.([]types.SetParam)

	sp := make([]*models.FCGISetParam, 0, len(d))

	for _, i := range d {
		sp = append(sp, &models.FCGISetParam{
			Cond:     i.Criterion,
			CondTest: i.Value,
			Format:   i.Format,
			Name:     i.Name,
		})
	}

	return sp, nil
}

func passHeader(p parser.Parser, name string) ([]*models.FCGIPassHeader, error) {
	data, err := p.Get(parser.FCGIApp, name, "pass-header")
	if err != nil {
		if errors.Is(err, parsererrors.ErrFetch) {
			return []*models.FCGIPassHeader{}, nil
		}

		return nil, err
	}

	d := data.([]types.PassHeader)

	ph := make([]*models.FCGIPassHeader, 0, len(d))

	for _, i := range d {
		ph = append(ph, &models.FCGIPassHeader{
			Cond:     i.Criterion,
			CondTest: i.Value,
			Name:     i.Name,
		})
	}

	return ph, nil
}

func ParseFCGIApp(p parser.Parser, name string) (*models.FCGIApp, error) {
	app := &models.FCGIApp{
		FCGIAppBase: models.FCGIAppBase{Name: name},
	}

	docRoot, err := genericString(p, name, "docroot")
	if err != nil {
		docRoot = &types.StringC{}
	}
	app.Docroot = &docRoot.Value

	index, err := genericString(p, name, "index")
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return nil, err
	}
	app.Index = index.Value

	app.LogStderrs, err = logStderr(p, name)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return nil, err
	}

	maxreqs, err := maxReqs(p, name)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return nil, err
	}
	if maxreqs != nil {
		app.MaxReqs = maxreqs.Reqs
	}

	app.PassHeaders, err = passHeader(p, name)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return nil, err
	}

	pathInfo, err := genericString(p, name, "path-info")
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return nil, err
	}
	app.PathInfo = pathInfo.Value

	app.SetParams, err = setParam(p, name)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return nil, err
	}

	_, mpxsConns := NewParseSection(parser.FCGIApp, name, p).checkOptions("mpxs-conns")
	if mpxsConns != nil {
		app.MpxsConns = mpxsConns.(string)
	}

	_, keepConn := NewParseSection(parser.FCGIApp, name, p).checkOptions("keep-conn")
	if keepConn != nil {
		app.KeepConn = keepConn.(string)
	}

	_, getValues := NewParseSection(parser.FCGIApp, name, p).checkOptions("get-values")
	if getValues != nil {
		app.GetValues = getValues.(string)
	}

	return app, nil
}

func (c *client) GetFCGIApplications(transactionID string) (int64, models.FCGIApps, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	fNames, err := p.SectionsGet(parser.FCGIApp)
	if err != nil {
		return v, nil, err
	}

	var apps models.FCGIApps

	for _, name := range fNames {
		app, parseErr := ParseFCGIApp(p, name)
		if parseErr != nil {
			continue
		}

		apps = append(apps, app)
	}

	return v, apps, nil
}

func (c *client) GetFCGIApplication(name string, transactionID string) (int64, *models.FCGIApp, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.FCGIApp, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("FCGI application %s does not exist", name))
	}

	app, parseErr := ParseFCGIApp(p, name)
	if parseErr != nil {
		return 0, nil, parseErr
	}

	return v, app, nil
}

func (c *client) DeleteFCGIApplication(name string, transactionID string, version int64) error {
	return c.deleteSection(parser.FCGIApp, name, transactionID, version)
}

func (c *client) EditFCGIApplication(name string, data *models.FCGIApp, transactionID string, version int64) error { //nolint:revive
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

	if !c.checkSectionExists(parser.FCGIApp, data.Name, p) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s does not exists", parser.FCGIApp, data.Name))

		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = SerializeFCGIAppSection(p, data); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) CreateFCGIApplication(data *models.FCGIApp, transactionID string, version int64) error {
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

	if c.checkSectionExists(parser.FCGIApp, data.Name, p) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.FCGIApp, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(parser.FCGIApp, data.Name); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if err = SerializeFCGIAppSection(p, data); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

//nolint:gocognit
func SerializeFCGIAppSection(p parser.Parser, data *models.FCGIApp) (err error) {
	if data == nil {
		return fmt.Errorf("empty FCGI app")
	}

	if data.Docroot == nil && len(*data.Docroot) == 0 {
		return fmt.Errorf("missing required docroot")
	} else if err = p.Set(parser.FCGIApp, data.Name, "docroot", types.StringC{Value: *data.Docroot}); err != nil {
		return err
	}

	if err = p.Set(parser.FCGIApp, data.Name, "option get-values", serializeSimpleOption(data.GetValues)); err != nil {
		return err
	}

	if len(data.Index) > 0 {
		if err = p.Set(parser.FCGIApp, data.Name, "index", types.StringC{Value: data.Index}); err != nil {
			return err
		}
	} else {
		if err = p.Set(parser.FCGIApp, data.Name, "index", nil); err != nil {
			return err
		}
	}

	if err = p.Set(parser.FCGIApp, data.Name, "option keep-conn", serializeSimpleOption(data.KeepConn)); err != nil {
		return err
	}

	if err = p.Set(parser.FCGIApp, data.Name, "log-stderr", nil); err != nil {
		return err
	}
	for _, i := range serializeLogStderr(data.LogStderrs) {
		if err = p.Set(parser.FCGIApp, data.Name, "log-stderr", i); err != nil {
			return err
		}
	}

	if data.MaxReqs > 0 {
		if err = p.Set(parser.FCGIApp, data.Name, "option max-reqs", types.OptionMaxReqs{Reqs: data.MaxReqs}); err != nil {
			return err
		}
	} else {
		if err = p.Set(parser.FCGIApp, data.Name, "option max-reqs", nil); err != nil {
			return err
		}
	}

	if err = p.Set(parser.FCGIApp, data.Name, "option mpxs-conns", serializeSimpleOption(data.GetValues)); err != nil {
		return err
	}

	if err = p.Set(parser.FCGIApp, data.Name, "pass-header", nil); err != nil {
		return err
	}
	for _, i := range serializePassHeader(data.PassHeaders) {
		if err = p.Set(parser.FCGIApp, data.Name, "pass-header", i); err != nil {
			return err
		}
	}

	if len(data.PathInfo) > 0 {
		if err = p.Set(parser.FCGIApp, data.Name, "path-info", types.StringC{Value: data.PathInfo}); err != nil {
			return err
		}
	} else {
		if err = p.Set(parser.FCGIApp, data.Name, "path-info", nil); err != nil {
			return err
		}
	}

	if err = p.Set(parser.FCGIApp, data.Name, "set-param", nil); err != nil {
		return err
	}
	for _, i := range serializeSetParam(data.SetParams) {
		if err = p.Set(parser.FCGIApp, data.Name, "set-param", i); err != nil {
			return err
		}
	}

	return nil
}

func serializeSetParam(setParams []*models.FCGISetParam) []types.SetParam {
	out := make([]types.SetParam, 0, len(setParams))

	for _, param := range setParams {
		p := types.SetParam{
			Name:      param.Name,
			Format:    param.Format,
			Criterion: param.Cond,
			Value:     param.CondTest,
		}

		out = append(out, p)
	}

	return out
}

func serializePassHeader(headers []*models.FCGIPassHeader) []types.PassHeader {
	out := make([]types.PassHeader, 0, len(headers))

	for _, header := range headers {
		h := types.PassHeader{
			Name:      header.Name,
			Criterion: header.Cond,
			Value:     header.CondTest,
		}

		out = append(out, h)
	}

	return out
}

func serializeSimpleOption(value string) *types.SimpleOption {
	switch value {
	case "enabled":
		return &types.SimpleOption{}
	case "disabled":
		return &types.SimpleOption{NoOption: true}
	default:
		return nil
	}
}

func serializeLogStderr(logs []*models.FCGILogStderr) []types.LogStdErr {
	out := make([]types.LogStdErr, 0, len(logs))

	for _, log := range logs {
		t := types.LogStdErr{
			Global:   log.Global,
			Address:  log.Address,
			Length:   log.Len,
			Format:   log.Format,
			Facility: log.Facility,
			Level:    log.Level,
			MinLevel: log.Minlevel,
		}

		if sample := log.Sample; sample != nil {
			t.SampleRange = *sample.Ranges
			t.SampleSize = *sample.Size
		}

		out = append(out, t)
	}

	return out
}
