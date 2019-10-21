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
	"strconv"

	"github.com/haproxytech/client-native/misc"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser"
	"github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/params"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetGlobalConfiguration returns configuration version and a
// struct representing Global configuration
func (c *Client) GetGlobalConfiguration(transactionID string) (int64, *models.Global, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "daemon")
	daemon := "enabled"
	if err == errors.ErrFetch {
		daemon = "disabled"
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "master-worker")
	masterWorker := true
	if err == errors.ErrFetch {
		masterWorker = false
	}

	data, err := p.Get(parser.Global, parser.GlobalSectionName, "maxconn")
	mConn := int64(0)
	if err == nil {
		maxConn := data.(*types.Int64C)
		mConn = maxConn.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "nbproc")
	nbproc := int64(0)
	if err == nil {
		nbProcParser := data.(*types.Int64C)
		nbproc = nbProcParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "nbthread")
	nbthread := int64(0)
	if err == nil {
		nbthreadParser := data.(*types.Int64C)
		nbthread = nbthreadParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "pidfile")
	pidfile := ""
	if err == nil {
		pidfileParser := data.(*types.StringC)
		pidfile = pidfileParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "stats socket")
	rAPIs := []*models.RuntimeAPI{}
	if err == nil {
		sockets := data.([]types.Socket)
		for _, s := range sockets {
			p := s.Path
			rAPI := &models.RuntimeAPI{Address: &p}
			for _, p := range s.Params {
				switch v := p.(type) {
				case *params.BindOptionDoubleWord:
					if v.Name == "expose-fd" && v.Value == "listener" {
						rAPI.ExposeFdListeners = true
					}
				case *params.BindOptionValue:
					if v.Name == "level" {
						rAPI.Level = v.Value
					} else if v.Name == "mode" {
						rAPI.Mode = v.Value
					} else if v.Name == "process" {
						rAPI.Process = v.Value
					}
				}
			}
			rAPIs = append(rAPIs, rAPI)
		}
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "cpu-map")
	cpuMaps := []*models.CPUMap{}
	if err == nil {
		cMaps := data.([]types.CPUMap)
		for _, m := range cMaps {
			cpuMap := &models.CPUMap{
				Process: &m.Process,
				CPUSet:  &m.CPUSet,
			}
			cpuMaps = append(cpuMaps, cpuMap)
		}
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "stats timeout")
	var statsTimeout *int64
	if err == errors.ErrFetch {
		statsTimeout = nil
	} else {
		statsTimeoutParser := data.(*types.StringC)
		statsTimeout = misc.ParseTimeout(statsTimeoutParser.Value)
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-bind-ciphers")
	sslCiphers := ""
	if err == nil {
		sslCiphersParser := data.(*types.StringC)
		sslCiphers = sslCiphersParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-bind-options")
	sslOptions := ""
	if err == nil {
		sslOptionsParser := data.(*types.StringC)
		sslOptions = sslOptionsParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "tune.ssl.default-dh-param")
	dhParam := int64(0)
	if err == nil {
		dhParamsParser := data.(*types.Int64C)
		dhParam = dhParamsParser.Value
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "external-check")
	externalCheck := true
	if err == errors.ErrFetch {
		externalCheck = false
	}

	g := &models.Global{
		Daemon:                daemon,
		MasterWorker:          masterWorker,
		Maxconn:               mConn,
		Nbproc:                nbproc,
		Nbthread:              nbthread,
		Pidfile:               pidfile,
		RuntimeApis:           rAPIs,
		StatsTimeout:          statsTimeout,
		CPUMaps:               cpuMaps,
		SslDefaultBindCiphers: sslCiphers,
		SslDefaultBindOptions: sslOptions,
		TuneSslDefaultDhParam: dhParam,
		ExternalCheck:         externalCheck,
	}

	return v, g, nil
}

// PushGlobalConfiguration pushes a Global config struct to global
// config gile
func (c *Client) PushGlobalConfiguration(data *models.Global, transactionID string, version int64) error {
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

	pDaemon := &types.Enabled{}
	if data.Daemon != "enabled" {
		pDaemon = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "daemon", pDaemon)

	pMasterWorker := &types.Enabled{}
	if data.MasterWorker == false {
		pMasterWorker = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "master-worker", pMasterWorker)

	pMaxConn := &types.Int64C{
		Value: data.Maxconn,
	}
	if data.Maxconn == 0 {
		pMaxConn = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "maxconn", pMaxConn)

	pNbProc := &types.Int64C{
		Value: data.Nbproc,
	}
	if data.Nbproc == 0 {
		pNbProc = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "nbproc", pNbProc)

	pNbthread := &types.Int64C{
		Value: data.Nbthread,
	}
	if data.Nbthread == 0 {
		pNbProc = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "nbthread", pNbthread)

	pPidfile := &types.StringC{
		Value: data.Pidfile,
	}
	if data.Pidfile == "" {
		pPidfile = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "pidfile", pPidfile)

	sockets := []types.Socket{}
	for _, rAPI := range data.RuntimeApis {
		s := types.Socket{
			Path:   *rAPI.Address,
			Params: []params.BindOption{},
		}
		if rAPI.ExposeFdListeners {
			p := &params.BindOptionDoubleWord{Name: "expose-fd", Value: "listeners"}
			s.Params = append(s.Params, p)
		}
		if rAPI.Level != "" {
			p := &params.BindOptionValue{Name: "level", Value: rAPI.Level}
			s.Params = append(s.Params, p)
		}
		if rAPI.Mode != "" {
			p := &params.BindOptionValue{Name: "mode", Value: rAPI.Mode}
			s.Params = append(s.Params, p)
		}
		if rAPI.Process != "" {
			p := &params.BindOptionValue{Name: "process", Value: rAPI.Process}
			s.Params = append(s.Params, p)
		}
		sockets = append(sockets, s)
	}
	p.Set(parser.Global, parser.GlobalSectionName, "stats socket", sockets)

	var statsTimeout *types.StringC
	if data.StatsTimeout != nil {
		statsTimeout = &types.StringC{Value: strconv.FormatInt(*data.StatsTimeout, 10)}
	} else {
		statsTimeout = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "stats timeout", statsTimeout)

	cpuMaps := []types.CPUMap{}
	for _, cpuMap := range data.CPUMaps {
		cm := types.CPUMap{
			Process: *cpuMap.Process,
			CPUSet:  *cpuMap.CPUSet,
		}
		cpuMaps = append(cpuMaps, cm)
	}
	p.Set(parser.Global, parser.GlobalSectionName, "cpu-map", cpuMaps)

	pSSLCiphers := &types.StringC{
		Value: data.SslDefaultBindCiphers,
	}
	if data.SslDefaultBindCiphers == "" {
		pSSLCiphers = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "ssl-default-bind-ciphers", pSSLCiphers)

	pSSLOptions := &types.StringC{
		Value: data.SslDefaultBindOptions,
	}
	if data.SslDefaultBindCiphers == "" {
		pSSLOptions = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "ssl-default-bind-options", pSSLOptions)

	pDhParams := &types.Int64C{
		Value: data.TuneSslDefaultDhParam,
	}
	if data.TuneSslDefaultDhParam == 0 {
		pDhParams = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "tune.ssl.default-dh-param", pDhParams)

	pExternalCheck := &types.Enabled{}
	if data.ExternalCheck == false {
		pExternalCheck = nil
	}
	p.Set(parser.Global, parser.GlobalSectionName, "external-check", pExternalCheck)
	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}
