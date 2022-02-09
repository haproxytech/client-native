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
	parser "github.com/haproxytech/config-parser/v4"
	"github.com/haproxytech/config-parser/v4/errors"
	"github.com/haproxytech/config-parser/v4/params"
	"github.com/haproxytech/config-parser/v4/types"

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

func ParseGlobalSection(p parser.Parser) (*models.Global, error) { //nolint:gocognit,gocyclo
	var chroot string
	data, err := p.Get(parser.Global, parser.GlobalSectionName, "chroot")
	if err == nil {
		chrootParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("chroot")
		}
		chroot = chrootParser.Value
	}

	var user string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "user")
	if err == nil {
		userParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("user")
		}
		user = userParser.Value
	}

	var group string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "group")
	if err == nil {
		groupParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("group")
		}
		group = groupParser.Value
	}

	var daemon string
	_, err = p.Get(parser.Global, parser.GlobalSectionName, "daemon")
	if !goerrors.Is(err, errors.ErrFetch) {
		daemon = "enabled"
	}

	var masterWorker bool
	_, err = p.Get(parser.Global, parser.GlobalSectionName, "master-worker")
	if !goerrors.Is(err, errors.ErrFetch) {
		masterWorker = true
	}

	var mConn int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "maxconn")
	if err == nil {
		maxConn, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("maxconn")
		}
		mConn = maxConn.Value
	}

	var nbproc int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "nbproc")
	if err == nil {
		nbProcParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("nbproc")
		}
		nbproc = nbProcParser.Value
	}

	var nbthread int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "nbthread")
	if err == nil {
		nbthreadParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("nbthread")
		}
		nbthread = nbthreadParser.Value
	}

	var pidfile string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "pidfile")
	if err == nil {
		pidfileParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("pidfile")
		}
		pidfile = pidfileParser.Value
	}

	var rAPIs []*models.RuntimeAPI
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "stats socket")
	if err == nil {
		sockets, ok := data.([]types.Socket)
		if !ok {
			return nil, misc.CreateTypeAssertError("stats socket")
		}
		for _, s := range sockets {
			p := s.Path
			rAPI := &models.RuntimeAPI{Address: &p}
			for _, p := range s.Params {
				switch v := p.(type) {
				case *params.BindOptionDoubleWord:
					if v.Name == "expose-fd" && v.Value == "listeners" {
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

	var cpuMaps []*models.CPUMap
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "cpu-map")
	if err == nil {
		cMaps, ok := data.([]types.CPUMap)
		if !ok {
			return nil, misc.CreateTypeAssertError("cpu-map")
		}
		for _, m := range cMaps {
			cpuMap := &models.CPUMap{
				Process: &m.Process,
				CPUSet:  &m.CPUSet,
			}
			cpuMaps = append(cpuMaps, cpuMap)
		}
	}

	var statsTimeout *int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "stats timeout")
	if goerrors.Is(err, errors.ErrFetch) {
		statsTimeout = nil
	} else {
		statsTimeoutParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("stats timeout")
		}
		statsTimeout = misc.ParseTimeout(statsTimeoutParser.Value)
	}

	var sslBindCiphers string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-bind-ciphers")
	if err == nil {
		sslBindCiphersParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("ssl-default-bind-ciphers")
		}
		sslBindCiphers = sslBindCiphersParser.Value
	}

	var sslBindCiphersuites string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-bind-ciphersuites")
	if err == nil {
		sslBindCiphersuitesParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("ssl-default-bind-ciphersuites")
		}
		sslBindCiphersuites = sslBindCiphersuitesParser.Value
	}

	var sslBindOptions string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-bind-options")
	if err == nil {
		sslBindOptionsParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("ssl-default-bind-options")
		}
		sslBindOptions = sslBindOptionsParser.Value
	}

	var sslServerCiphers string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-server-ciphers")
	if err == nil {
		sslServerCiphersParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("ssl-default-server-ciphers")
		}
		sslServerCiphers = sslServerCiphersParser.Value
	}

	var sslServerCiphersuites string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-server-ciphersuites")
	if err == nil {
		sslServerCiphersuitesParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("ssl-default-server-ciphersuites")
		}
		sslServerCiphersuites = sslServerCiphersuitesParser.Value
	}

	var sslServerOptions string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-server-options")
	if err == nil {
		sslServerOptionsParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("ssl-default-server-options")
		}
		sslServerOptions = sslServerOptionsParser.Value
	}

	var dhParam int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "tune.ssl.default-dh-param")
	if err == nil {
		dhParamsParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("tune.ssl.default-dh-param")
		}
		dhParam = dhParamsParser.Value
	}

	var sslModeAsync string
	data, _ = p.Get(parser.Global, parser.GlobalSectionName, "ssl-mode-async")
	if _, ok := data.(*types.SslModeAsync); ok {
		sslModeAsync = "enabled"
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "external-check")
	externalCheck := true
	if goerrors.Is(err, errors.ErrFetch) {
		externalCheck = false
	}

	var luaPrependPath []*models.LuaPrependPath
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "lua-prepend-path")
	if err == nil {
		lpp, ok := data.([]types.LuaPrependPath)
		if !ok {
			return nil, misc.CreateTypeAssertError("lua-prepend-path")
		}
		for _, l := range lpp {
			path := l.Path
			luaPrependPath = append(luaPrependPath, &models.LuaPrependPath{Path: &path, Type: l.Type})
		}
	}

	var luaLoads []*models.LuaLoad
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "lua-load")
	if err == nil {
		luas, ok := data.([]types.LuaLoad)
		if !ok {
			return nil, misc.CreateTypeAssertError("lua-load")
		}
		for _, l := range luas {
			file := l.File
			luaLoads = append(luaLoads, &models.LuaLoad{File: &file})
		}
	}

	var globalLogSendHostName *models.GlobalLogSendHostname
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "log-send-hostname")
	if err == nil {
		logSendHostName := "enabled"
		logSendHostNameParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("log-send-hostname")
		}
		globalLogSendHostName = &models.GlobalLogSendHostname{
			Enabled: &logSendHostName,
			Param:   logSendHostNameParser.Value,
		}
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
		LuaPrependPath:               luaPrependPath,
		LogSendHostname:              globalLogSendHostName,
	}

	return g, nil
}

func SerializeGlobalSection(p parser.Parser, data *models.Global) error { //nolint:gocognit,gocyclo
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

	luaPrependPath := []types.LuaPrependPath{}
	for _, l := range data.LuaPrependPath {
		lpp := types.LuaPrependPath{
			Path: *l.Path,
			Type: l.Type,
		}
		luaPrependPath = append(luaPrependPath, lpp)
	}

	if err := p.Set(parser.Global, parser.GlobalSectionName, "lua-prepend-path", luaPrependPath); err != nil {
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
