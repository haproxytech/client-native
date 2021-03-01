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
	goerrors "errors"
	"strconv"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v3"
	"github.com/haproxytech/config-parser/v3/errors"
	"github.com/haproxytech/config-parser/v3/params"
	"github.com/haproxytech/config-parser/v3/types"

	"github.com/haproxytech/client-native/v2/misc"
	"github.com/haproxytech/client-native/v2/models"
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

	g, err := ParseGlobalSection(p)
	if err != nil {
		return 0, nil, err
	}

	return v, g, nil
}

// PushGlobalConfiguration pushes a Global config struct to global
// config file
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

	if err := SerializeGlobalSection(p, data); err != nil {
		return err
	}
	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func ParseGlobalSection(p *parser.Parser) (*models.Global, error) { //nolint:gocognit,gocyclo
	data, err := p.Get(parser.Global, parser.GlobalSectionName, "chroot")
	chroot := ""
	if err == nil {
		chrootParser := data.(*types.StringC)
		chroot = chrootParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "user")
	user := ""
	if err == nil {
		userParser := data.(*types.StringC)
		user = userParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "group")
	group := ""
	if err == nil {
		groupParser := data.(*types.StringC)
		group = groupParser.Value
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "daemon")
	daemon := "enabled"
	if goerrors.Is(err, errors.ErrFetch) {
		daemon = "disabled"
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "master-worker")
	masterWorker := true
	if goerrors.Is(err, errors.ErrFetch) {
		masterWorker = false
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "maxconn")
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
					switch v.Name {
					case "level":
						rAPI.Level = v.Value
					case "mode":
						rAPI.Mode = v.Value
					case "process":
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
	if goerrors.Is(err, errors.ErrFetch) {
		statsTimeout = nil
	} else {
		statsTimeoutParser := data.(*types.StringC)
		statsTimeout = misc.ParseTimeout(statsTimeoutParser.Value)
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-bind-ciphers")
	sslBindCiphers := ""
	if err == nil {
		sslBindCiphersParser := data.(*types.StringC)
		sslBindCiphers = sslBindCiphersParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-bind-ciphersuites")
	sslBindCiphersuites := ""
	if err == nil {
		sslBindCiphersuitesParser := data.(*types.StringC)
		sslBindCiphersuites = sslBindCiphersuitesParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-bind-options")
	sslBindOptions := ""
	if err == nil {
		sslBindOptionsParser := data.(*types.StringC)
		sslBindOptions = sslBindOptionsParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-server-ciphers")
	sslServerCiphers := ""
	if err == nil {
		sslServerCiphersParser := data.(*types.StringC)
		sslServerCiphers = sslServerCiphersParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-server-ciphersuites")
	sslServerCiphersuites := ""
	if err == nil {
		sslServerCiphersuitesParser := data.(*types.StringC)
		sslServerCiphersuites = sslServerCiphersuitesParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-server-options")
	sslServerOptions := ""
	if err == nil {
		sslServerOptionsParser := data.(*types.StringC)
		sslServerOptions = sslServerOptionsParser.Value
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "tune.ssl.default-dh-param")
	dhParam := int64(0)
	if err == nil {
		dhParamsParser := data.(*types.Int64C)
		dhParam = dhParamsParser.Value
	}

	data, _ = p.Get(parser.Global, parser.GlobalSectionName, "ssl-mode-async")
	sslModeAsync := "disabled"
	if _, ok := data.(*types.SslModeAsync); ok {
		sslModeAsync = "enabled"
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "external-check")
	externalCheck := true
	if goerrors.Is(err, errors.ErrFetch) {
		externalCheck = false
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "lua-load")
	luaLoads := []*models.LuaLoad{}
	if err == nil {
		luas := data.([]types.LuaLoad)
		for _, l := range luas {
			file := l.File
			luaLoads = append(luaLoads, &models.LuaLoad{File: &file})
		}
	}

	data, err = p.Get(parser.Global, parser.GlobalSectionName, "log-send-hostname")
	logSendHostName := "disabled"
	globalLogSendHostName := &models.GlobalLogSendHostname{Enabled: &logSendHostName}
	if err == nil {
		logSendHostNameParser := data.(*types.StringC)
		globalLogSendHostName.Param = logSendHostNameParser.Value
		logSendHostName = "enabled"
		globalLogSendHostName.Enabled = &logSendHostName
	}

	g := &models.Global{
		User:                         user,
		Group:                        group,
		Chroot:                       chroot,
		Daemon:                       daemon,
		MasterWorker:                 masterWorker,
		Maxconn:                      mConn,
		Nbproc:                       nbproc,
		Nbthread:                     nbthread,
		Pidfile:                      pidfile,
		RuntimeAPIs:                  rAPIs,
		StatsTimeout:                 statsTimeout,
		CPUMaps:                      cpuMaps,
		SslDefaultBindCiphers:        sslBindCiphers,
		SslDefaultBindCiphersuites:   sslBindCiphersuites,
		SslDefaultBindOptions:        sslBindOptions,
		SslDefaultServerCiphers:      sslServerCiphers,
		SslDefaultServerCiphersuites: sslServerCiphersuites,
		SslDefaultServerOptions:      sslServerOptions,
		SslModeAsync:                 sslModeAsync,
		TuneSslDefaultDhParam:        dhParam,
		ExternalCheck:                externalCheck,
		LuaLoads:                     luaLoads,
		LogSendHostname:              globalLogSendHostName,
	}

	return g, nil
}

func SerializeGlobalSection(p *parser.Parser, data *models.Global) error { //nolint:gocognit,gocyclo
	pChroot := &types.StringC{
		Value: data.Chroot,
	}
	if data.Chroot == "" {
		pChroot = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "chroot", pChroot); err != nil {
		return err
	}
	pUser := &types.StringC{
		Value: data.User,
	}
	if data.User == "" {
		pUser = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "user", pUser); err != nil {
		return err
	}
	pGroup := &types.StringC{
		Value: data.Group,
	}
	if data.Group == "" {
		pGroup = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "group", pGroup); err != nil {
		return err
	}
	pDaemon := &types.Enabled{}
	if data.Daemon != "enabled" {
		pDaemon = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "daemon", pDaemon); err != nil {
		return err
	}
	pMasterWorker := &types.Enabled{}
	if !data.MasterWorker {
		pMasterWorker = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "master-worker", pMasterWorker); err != nil {
		return err
	}
	pMaxConn := &types.Int64C{
		Value: data.Maxconn,
	}
	if data.Maxconn == 0 {
		pMaxConn = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "maxconn", pMaxConn); err != nil {
		return err
	}
	pNbProc := &types.Int64C{
		Value: data.Nbproc,
	}
	if data.Nbproc == 0 {
		pNbProc = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "nbproc", pNbProc); err != nil {
		return err
	}
	pNbthread := &types.Int64C{
		Value: data.Nbthread,
	}
	if data.Nbthread == 0 {
		pNbthread = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "nbthread", pNbthread); err != nil {
		return err
	}
	pPidfile := &types.StringC{
		Value: data.Pidfile,
	}
	if data.Pidfile == "" {
		pPidfile = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "pidfile", pPidfile); err != nil {
		return err
	}
	sockets := []types.Socket{}
	for _, rAPI := range data.RuntimeAPIs {
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
	if err := p.Set(parser.Global, parser.GlobalSectionName, "stats socket", sockets); err != nil {
		return err
	}
	var statsTimeout *types.StringC
	if data.StatsTimeout != nil {
		statsTimeout = &types.StringC{Value: strconv.FormatInt(*data.StatsTimeout, 10)}
	} else {
		statsTimeout = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "stats timeout", statsTimeout); err != nil {
		return err
	}
	cpuMaps := []types.CPUMap{}
	for _, cpuMap := range data.CPUMaps {
		cm := types.CPUMap{
			Process: *cpuMap.Process,
			CPUSet:  *cpuMap.CPUSet,
		}
		cpuMaps = append(cpuMaps, cm)
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "cpu-map", cpuMaps); err != nil {
		return err
	}
	pSSLBindCiphers := &types.StringC{
		Value: data.SslDefaultBindCiphers,
	}
	if data.SslDefaultBindCiphers == "" {
		pSSLBindCiphers = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-default-bind-ciphers", pSSLBindCiphers); err != nil {
		return err
	}
	pSSLBindCiphersuites := &types.StringC{
		Value: data.SslDefaultBindCiphersuites,
	}
	if data.SslDefaultBindCiphersuites == "" {
		pSSLBindCiphersuites = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-default-bind-ciphersuites", pSSLBindCiphersuites); err != nil {
		return err
	}
	pSSLBindOptions := &types.StringC{
		Value: data.SslDefaultBindOptions,
	}
	if data.SslDefaultBindOptions == "" {
		pSSLBindOptions = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-default-bind-options", pSSLBindOptions); err != nil {
		return err
	}
	pSSLServerCiphers := &types.StringC{
		Value: data.SslDefaultServerCiphers,
	}
	if data.SslDefaultServerCiphers == "" {
		pSSLServerCiphers = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-default-server-ciphers", pSSLServerCiphers); err != nil {
		return err
	}
	pSSLServerCiphersuites := &types.StringC{
		Value: data.SslDefaultServerCiphersuites,
	}
	if data.SslDefaultServerCiphersuites == "" {
		pSSLServerCiphersuites = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-default-server-ciphersuites", pSSLServerCiphersuites); err != nil {
		return err
	}
	pSSLServerOptions := &types.StringC{
		Value: data.SslDefaultServerOptions,
	}
	if data.SslDefaultServerOptions == "" {
		pSSLServerOptions = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-default-server-options", pSSLServerOptions); err != nil {
		return err
	}
	pDhParams := &types.Int64C{
		Value: data.TuneSslDefaultDhParam,
	}
	if data.TuneSslDefaultDhParam == 0 {
		pDhParams = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "tune.ssl.default-dh-param", pDhParams); err != nil {
		return err
	}
	sslModeAsync := &types.SslModeAsync{}
	if data.SslModeAsync != "enabled" {
		sslModeAsync = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-mode-async", sslModeAsync); err != nil {
		return err
	}

	luaLoads := []types.LuaLoad{}
	for _, lua := range data.LuaLoads {
		ll := types.LuaLoad{
			File: *lua.File,
		}
		luaLoads = append(luaLoads, ll)
	}

	if err := p.Set(parser.Global, parser.GlobalSectionName, "lua-load", luaLoads); err != nil {
		return err
	}

	logSendHostName := &types.StringC{}
	if data.LogSendHostname == nil || *data.LogSendHostname.Enabled == "disabled" {
		logSendHostName = nil
	} else {
		logSendHostName.Value = data.LogSendHostname.Param
	}

	if err := p.Set(parser.Global, parser.GlobalSectionName, "log-send-hostname", logSendHostName); err != nil {
		return err
	}

	pExternalCheck := &types.Enabled{}
	if !data.ExternalCheck {
		pExternalCheck = nil
	}

	return p.Set(parser.Global, parser.GlobalSectionName, "external-check", pExternalCheck)
}
