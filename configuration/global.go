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

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v4"
	parser_errors "github.com/haproxytech/config-parser/v4/errors"
	"github.com/haproxytech/config-parser/v4/params"
	"github.com/haproxytech/config-parser/v4/types"

	"github.com/haproxytech/client-native/v3/misc"
	"github.com/haproxytech/client-native/v3/models"
)

type Global interface {
	GetGlobalConfiguration(transactionID string) (int64, *models.Global, error)
	PushGlobalConfiguration(data *models.Global, transactionID string, version int64) error
}

// GetGlobalConfiguration returns configuration version and a
// struct representing Global configuration
func (c *client) GetGlobalConfiguration(transactionID string) (int64, *models.Global, error) {
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
func (c *client) PushGlobalConfiguration(data *models.Global, transactionID string, version int64) error {
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

	if err := SerializeGlobalSection(p, data); err != nil {
		return err
	}
	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func ParseGlobalSection(p parser.Parser) (*models.Global, error) { //nolint:gocognit,gocyclo,cyclop
	var chroot string
	data, err := p.Get(parser.Global, parser.GlobalSectionName, "chroot")
	if err == nil {
		chrootParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("chroot")
		}
		chroot = chrootParser.Value
	}

	var caBase string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ca-base")
	if err == nil {
		caBaseParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("ca-base")
		}
		caBase = caBaseParser.Value
	}

	var crtBase string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "crt-base")
	if err == nil {
		crtBaseParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("crt-base")
		}
		crtBase = crtBaseParser.Value
	}

	var srvStateBase string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "server-state-base")
	if err == nil {
		srvStateBase = data.(*types.StringC).Value
	}

	var srvStateFile string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "server-state-file")
	if err == nil {
		srvStateFile = data.(*types.StringC).Value
	}

	var hardStop *int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "hard-stop-after")
	if err == nil {
		hardStopParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("hard-stop-after")
		}
		hardStop = misc.ParseTimeout(hardStopParser.Value)
	}

	var localPeer string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "localpeer")
	if err == nil {
		userParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("localpeer")
		}
		localPeer = userParser.Value
	}

	var uid int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "uid")
	if err == nil {
		uidParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("uid")
		}
		uid = uidParser.Value
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

	var gid int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "gid")
	if err == nil {
		gidParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("gid")
		}
		gid = gidParser.Value
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
	if !errors.Is(err, parser_errors.ErrFetch) {
		daemon = "enabled"
	}

	var masterWorker bool
	_, err = p.Get(parser.Global, parser.GlobalSectionName, "master-worker")
	if !errors.Is(err, parser_errors.ErrFetch) {
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
			rAPI.BindParams = parseBindParams(s.Params)
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
	if errors.Is(err, parser_errors.ErrFetch) {
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

	var sslLoadExtraDelExt bool
	sslLoadExtraDelExt, err = parseBoolOption(p, "ssl-load-extra-del-ext")
	if err != nil {
		return nil, err
	}

	var sslModeAsync string
	data, _ = p.Get(parser.Global, parser.GlobalSectionName, "ssl-mode-async")
	if _, ok := data.(*types.SslModeAsync); ok {
		sslModeAsync = "enabled"
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "external-check")
	externalCheck := true
	if errors.Is(err, parser_errors.ErrFetch) {
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

	var h1CaseAdjusts []*models.H1CaseAdjust
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "h1-case-adjust")
	if err == nil {
		cases, ok := data.([]types.H1CaseAdjust)
		if !ok {
			return nil, misc.CreateTypeAssertError("h1-case-adjust")
		}
		for _, c := range cases {
			from := c.From
			to := c.To
			h1CaseAdjusts = append(h1CaseAdjusts, &models.H1CaseAdjust{From: &from, To: &to})
		}
	}

	var h1CaseAdjustFile string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "h1-case-adjust-file")
	if err == nil {
		caseFileParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("h1-case-adjust-file")
		}
		h1CaseAdjustFile = caseFileParser.Value
	}

	var busyPolling bool
	_, err = p.Get(parser.Global, parser.GlobalSectionName, "busy-polling")
	if !errors.Is(err, parser_errors.ErrFetch) {
		busyPolling = true
	}

	var maxSpreadChecks int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "max-spread-checks")
	if err == nil {
		maxSpreadChecksParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("max-spread-checks")
		}
		maxSpreadChecks = maxSpreadChecksParser.Value
	}

	var maxconnrate int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "maxconnrate")
	if err == nil {
		maxconnrateParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("maxconnrate")
		}
		maxconnrate = maxconnrateParser.Value
	}

	var maxcomprate int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "maxcomprate")
	if err == nil {
		maxcomprateParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("maxcomprate")
		}
		maxcomprate = maxcomprateParser.Value
	}

	var maxcompcpuusage int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "maxcompcpuusage")
	if err == nil {
		maxcompcpuusageParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("maxcompcpuusage")
		}
		maxcompcpuusage = maxcompcpuusageParser.Value
	}

	var maxpipes int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "maxpipes")
	if err == nil {
		maxpipesParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("maxpipes")
		}
		maxpipes = maxpipesParser.Value
	}

	var maxsessrate int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "maxsessrate")
	if err == nil {
		maxsessrateParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("maxsessrate")
		}
		maxsessrate = maxsessrateParser.Value
	}

	var maxsslconn int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "maxsslconn")
	if err == nil {
		maxsslconnParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("maxsslconn")
		}
		maxsslconn = maxsslconnParser.Value
	}

	var maxsslrate int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "maxsslrate")
	if err == nil {
		maxsslrateParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("maxsslrate")
		}
		maxsslrate = maxsslrateParser.Value
	}

	var maxzlibmem int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "maxzlibmem")
	if err == nil {
		maxzlibmemParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("maxzlibmem")
		}
		maxzlibmem = maxzlibmemParser.Value
	}

	var noepoll bool
	_, err = p.Get(parser.Global, parser.GlobalSectionName, "noepoll")
	if !errors.Is(err, parser_errors.ErrFetch) {
		noepoll = true
	}

	var nokqueue bool
	_, err = p.Get(parser.Global, parser.GlobalSectionName, "nokqueue")
	if !errors.Is(err, parser_errors.ErrFetch) {
		nokqueue = true
	}

	var noevports bool
	_, err = p.Get(parser.Global, parser.GlobalSectionName, "noevports")
	if !errors.Is(err, parser_errors.ErrFetch) {
		noevports = true
	}

	var nopoll bool
	_, err = p.Get(parser.Global, parser.GlobalSectionName, "nopoll")
	if !errors.Is(err, parser_errors.ErrFetch) {
		nopoll = true
	}

	var nosplice bool
	_, err = p.Get(parser.Global, parser.GlobalSectionName, "nosplice")
	if !errors.Is(err, parser_errors.ErrFetch) {
		nosplice = true
	}

	var nogetaddrinfo bool
	_, err = p.Get(parser.Global, parser.GlobalSectionName, "nogetaddrinfo")
	if !errors.Is(err, parser_errors.ErrFetch) {
		nogetaddrinfo = true
	}

	var noreuseport bool
	_, err = p.Get(parser.Global, parser.GlobalSectionName, "noreuseport")
	if !errors.Is(err, parser_errors.ErrFetch) {
		noreuseport = true
	}

	profilingTasks, err := parseAutoOnOffOption(p, "profiling.tasks")
	if err != nil {
		return nil, err
	}

	var spreadChecks int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "spread-checks")
	if err == nil {
		spreadChecksParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("spread-checks")
		}
		spreadChecks = spreadChecksParser.Value
	}

	tuneOptions, err := parseTuneOptions(p)
	if err != nil {
		return nil, err
	}
	// deprecated option
	dhParam := int64(0)
	if tuneOptions != nil {
		dhParam = tuneOptions.SslDefaultDhParam
	}

	global := &models.Global{
		UID:                          uid,
		User:                         user,
		Gid:                          gid,
		Group:                        group,
		Chroot:                       chroot,
		Localpeer:                    localPeer,
		CaBase:                       caBase,
		CrtBase:                      crtBase,
		ServerStateBase:              srvStateBase,
		ServerStateFile:              srvStateFile,
		HardStopAfter:                hardStop,
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
		SslLoadExtraDelExt:           sslLoadExtraDelExt,
		SslModeAsync:                 sslModeAsync,
		TuneOptions:                  tuneOptions,
		TuneSslDefaultDhParam:        dhParam,
		ExternalCheck:                externalCheck,
		LuaLoads:                     luaLoads,
		LuaPrependPath:               luaPrependPath,
		LogSendHostname:              globalLogSendHostName,
		H1CaseAdjusts:                h1CaseAdjusts,
		H1CaseAdjustFile:             h1CaseAdjustFile,
		BusyPolling:                  busyPolling,
		MaxSpreadChecks:              maxSpreadChecks,
		Maxconnrate:                  maxconnrate,
		Maxcomprate:                  maxcomprate,
		Maxcompcpuusage:              maxcompcpuusage,
		Maxpipes:                     maxpipes,
		Maxsessrate:                  maxsessrate,
		Maxsslconn:                   maxsslconn,
		Maxsslrate:                   maxsslrate,
		Maxzlibmem:                   maxzlibmem,
		Noepoll:                      noepoll,
		Nokqueue:                     nokqueue,
		Noevports:                    noevports,
		Nopoll:                       nopoll,
		Nosplice:                     nosplice,
		Nogetaddrinfo:                nogetaddrinfo,
		Noreuseport:                  noreuseport,
		ProfilingTasks:               profilingTasks,
		SpreadChecks:                 spreadChecks,
	}

	return global, nil
}

func SerializeGlobalSection(p parser.Parser, data *models.Global) error { //nolint:gocognit,gocyclo,cyclop
	pChroot := &types.StringC{
		Value: data.Chroot,
	}
	if data.Chroot == "" {
		pChroot = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "chroot", pChroot); err != nil {
		return err
	}
	pCaBase := &types.StringC{
		Value: data.CaBase,
	}
	if data.CaBase == "" {
		pCaBase = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ca-base", pCaBase); err != nil {
		return err
	}
	pCrtBase := &types.StringC{
		Value: data.CrtBase,
	}
	if data.CrtBase == "" {
		pCrtBase = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "crt-base", pCrtBase); err != nil {
		return err
	}
	pLocalPeer := &types.StringC{
		Value: data.Localpeer,
	}
	if data.Localpeer == "" {
		pLocalPeer = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "localpeer", pLocalPeer); err != nil {
		return err
	}
	pSrvStateBase := &types.StringC{
		Value: data.ServerStateBase,
	}
	if data.ServerStateBase == "" {
		pSrvStateBase = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "server-state-base", pSrvStateBase); err != nil {
		return err
	}
	pSrvStateFile := &types.StringC{
		Value: data.ServerStateFile,
	}
	if data.ServerStateFile == "" {
		pSrvStateFile = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "server-state-file", pSrvStateFile); err != nil {
		return err
	}
	var pHardStop *types.StringC
	if data.HardStopAfter != nil {
		pHardStop = &types.StringC{
			Value: strconv.FormatInt(*data.HardStopAfter, 10),
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "hard-stop-after", pHardStop); err != nil {
		return err
	}
	pUID := &types.Int64C{
		Value: data.UID,
	}
	if data.UID == 0 {
		pUID = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "uid", pUID); err != nil {
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
	pGID := &types.Int64C{
		Value: data.Gid,
	}
	if data.Gid == 0 {
		pGID = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "gid", pGID); err != nil {
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
		socket := types.Socket{
			Path:   *rAPI.Address,
			Params: []params.BindOption{},
		}
		socket.Params = serializeBindParams(rAPI.BindParams, "")
		sockets = append(sockets, socket)
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

	if err := serializeBoolOption(p, "ssl-load-extra-del-ext", data.SslLoadExtraDelExt); err != nil {
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
	if err := p.Set(parser.Global, parser.GlobalSectionName, "external-check", pExternalCheck); err != nil {
		return err
	}

	busyPolling := &types.Enabled{}
	if !data.BusyPolling {
		busyPolling = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "busy-polling", busyPolling); err != nil {
		return err
	}

	maxSpreadChecks := &types.Int64C{
		Value: data.MaxSpreadChecks,
	}
	if data.MaxSpreadChecks == 0 {
		maxSpreadChecks = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "max-spread-checks", maxSpreadChecks); err != nil {
		return err
	}

	maxconnrate := &types.Int64C{
		Value: data.Maxconnrate,
	}
	if data.Maxconnrate == 0 {
		maxconnrate = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "maxconnrate", maxconnrate); err != nil {
		return err
	}

	maxcomprate := &types.Int64C{
		Value: data.Maxcomprate,
	}
	if data.Maxcomprate == 0 {
		maxcomprate = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "maxcomprate", maxcomprate); err != nil {
		return err
	}

	maxcompcpuusage := &types.Int64C{
		Value: data.Maxcompcpuusage,
	}
	if data.Maxcompcpuusage == 0 {
		maxcompcpuusage = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "maxcompcpuusage", maxcompcpuusage); err != nil {
		return err
	}

	maxpipes := &types.Int64C{
		Value: data.Maxpipes,
	}
	if data.Maxpipes == 0 {
		maxpipes = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "maxpipes", maxpipes); err != nil {
		return err
	}

	maxsessrate := &types.Int64C{
		Value: data.Maxsessrate,
	}
	if data.Maxsessrate == 0 {
		maxsessrate = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "maxsessrate", maxsessrate); err != nil {
		return err
	}

	maxsslconn := &types.Int64C{
		Value: data.Maxsslconn,
	}
	if data.Maxsslconn == 0 {
		maxsslconn = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "maxsslconn", maxsslconn); err != nil {
		return err
	}

	maxsslrate := &types.Int64C{
		Value: data.Maxsslrate,
	}
	if data.Maxsslrate == 0 {
		maxsslrate = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "maxsslrate", maxsslrate); err != nil {
		return err
	}

	maxzlibmem := &types.Int64C{
		Value: data.Maxzlibmem,
	}
	if data.Maxzlibmem == 0 {
		maxzlibmem = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "maxzlibmem", maxzlibmem); err != nil {
		return err
	}

	noepoll := &types.Enabled{}
	if !data.Noepoll {
		noepoll = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "noepoll", noepoll); err != nil {
		return err
	}

	nokqueue := &types.Enabled{}
	if !data.Nokqueue {
		nokqueue = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "nokqueue", nokqueue); err != nil {
		return err
	}

	noevports := &types.Enabled{}
	if !data.Noevports {
		noevports = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "noevports", noevports); err != nil {
		return err
	}

	nopoll := &types.Enabled{}
	if !data.Nopoll {
		nopoll = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "nopoll", nopoll); err != nil {
		return err
	}

	nosplice := &types.Enabled{}
	if !data.Nosplice {
		nosplice = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "nosplice", nosplice); err != nil {
		return err
	}

	nogetaddrinfo := &types.Enabled{}
	if !data.Nogetaddrinfo {
		nogetaddrinfo = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "nogetaddrinfo", nogetaddrinfo); err != nil {
		return err
	}

	noreuseport := &types.Enabled{}
	if !data.Noreuseport {
		noreuseport = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "noreuseport", noreuseport); err != nil {
		return err
	}

	if err := serializeAutoOnOffOption(p, "profiling.tasks", data.ProfilingTasks); err != nil {
		return err
	}

	spreadChecks := &types.Int64C{
		Value: data.SpreadChecks,
	}
	if data.SpreadChecks == 0 {
		spreadChecks = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "spread-checks", spreadChecks); err != nil {
		return err
	}

	pH1CaseAdjusts := []types.H1CaseAdjust{}
	if data.H1CaseAdjusts != nil && len(data.H1CaseAdjusts) > 0 {
		for _, caseAdjust := range data.H1CaseAdjusts {
			if caseAdjust != nil && caseAdjust.From != nil && caseAdjust.To != nil {
				ca := types.H1CaseAdjust{From: *caseAdjust.From, To: *caseAdjust.To}
				pH1CaseAdjusts = append(pH1CaseAdjusts, ca)
			}
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "h1-case-adjust", pH1CaseAdjusts); err != nil {
		return err
	}

	pH1CaseAdjustFile := &types.StringC{Value: data.H1CaseAdjustFile}
	if data.H1CaseAdjustFile == "" {
		pH1CaseAdjustFile = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "h1-case-adjust-file", pH1CaseAdjustFile); err != nil {
		return err
	}

	// deprecated option
	if data.TuneSslDefaultDhParam != 0 {
		if data.TuneOptions != nil && data.TuneOptions.SslDefaultDhParam == 0 {
			data.TuneOptions.SslDefaultDhParam = data.TuneSslDefaultDhParam
		}
		if data.TuneOptions == nil {
			data.TuneOptions = &models.GlobalTuneOptions{SslDefaultDhParam: data.TuneSslDefaultDhParam}
		}
	}
	return serializeTuneOptions(p, data.TuneOptions)
}

func serializeTuneOptions(p parser.Parser, options *models.GlobalTuneOptions) error { //nolint:gocognit,gocyclo,cyclop
	if options == nil {
		return nil
	}
	if err := serializeInt64POption(p, "tune.buffers.limit", options.BuffersLimit); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.buffers.reserve", options.BuffersReserve); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.bufsize", options.Bufsize); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.comp.maxlevel", options.CompMaxlevel); err != nil {
		return err
	}
	if err := serializeBoolOption(p, "tune.fail-alloc", options.FailAlloc); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.h2.header-table-size", options.H2HeaderTableSize); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.h2.initial-window-size", options.H2InitialWindowSize); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.h2.max-concurrent-streams", options.H2MaxConcurrentStreams); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.h2.max-frame-size", options.H2MaxFrameSize); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.http.cookielen", options.HTTPCookielen); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.http.logurilen", options.HTTPLogurilen); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.http.maxhdr", options.HTTPMaxhdr); err != nil {
		return err
	}
	if err := serializeOnOffOption(p, "tune.idle-pool.shared", options.IdlePoolShared); err != nil {
		return err
	}
	if err := serializeTimeoutSizeOption(p, "tune.idletimer", options.Idletimer); err != nil {
		return err
	}
	if err := serializeOnOffOption(p, "tune.listener.multi-queue", options.ListenerMultiQueue); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.lua.forced-yield", options.LuaForcedYield); err != nil {
		return err
	}
	if err := serializeBoolOption(p, "tune.lua.maxmem", options.LuaMaxmem); err != nil {
		return err
	}
	if err := serializeTimeoutSizeOption(p, "tune.lua.session-timeout", options.LuaSessionTimeout); err != nil {
		return err
	}
	if err := serializeTimeoutSizeOption(p, "tune.lua.task-timeout", options.LuaTaskTimeout); err != nil {
		return err
	}
	if err := serializeTimeoutSizeOption(p, "tune.lua.service-timeout", options.LuaServiceTimeout); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.maxaccept", options.Maxaccept); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.maxpollevents", options.Maxpollevents); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.maxrewrite", options.Maxrewrite); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.pattern.cache-size", options.PatternCacheSize); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.pipesize", options.Pipesize); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.pool-high-fd-ratio", options.PoolHighFdRatio); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.pool-low-fd-ratio", options.PoolLowFdRatio); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.rcvbuf.client", options.RcvbufClient); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.rcvbuf.server", options.RcvbufServer); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.recv_enough", options.RecvEnough); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.runqueue-depth", options.RunqueueDepth); err != nil {
		return err
	}
	if err := serializeOnOffOption(p, "tune.sched.low-latency", options.SchedLowLatency); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.sndbuf.client", options.SndbufClient); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.sndbuf.server", options.SndbufServer); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.ssl.cachesize", options.SslCachesize); err != nil {
		return err
	}
	if err := serializeBoolOption(p, "tune.ssl.force-private-cache", options.SslForcePrivateCache); err != nil {
		return err
	}
	if err := serializeOnOffOption(p, "tune.ssl.keylog", options.SslKeylog); err != nil {
		return err
	}
	if err := serializeTimeoutSizeOption(p, "tune.ssl.lifetime", options.SslLifetime); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.ssl.maxrecord", options.SslMaxrecord); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.ssl.default-dh-param", options.SslDefaultDhParam); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.ssl.ssl-ctx-cache-size", options.SslCtxCacheSize); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.ssl.capture-buffer-size", options.SslCaptureBufferSize); err != nil {
		return err
	}
	if err := serializeTimeoutSizeOption(p, "tune.vars.global-max-size", options.VarsGlobalMaxSize); err != nil {
		return err
	}
	if err := serializeTimeoutSizeOption(p, "tune.vars.proc-max-size", options.VarsProcMaxSize); err != nil {
		return err
	}
	if err := serializeTimeoutSizeOption(p, "tune.vars.reqres-max-size", options.VarsReqresMaxSize); err != nil {
		return err
	}
	if err := serializeTimeoutSizeOption(p, "tune.vars.sess-max-size", options.VarsSessMaxSize); err != nil {
		return err
	}
	if err := serializeTimeoutSizeOption(p, "tune.vars.txn-max-size", options.VarsTxnMaxSize); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.zlib.memlevel", options.ZlibMemlevel); err != nil {
		return err
	}
	return serializeInt64Option(p, "tune.zlib.windowsize", options.ZlibWindowsize)
}

func serializeTimeoutSizeOption(p parser.Parser, option string, data *int64) error {
	var value *types.StringC
	if data == nil {
		value = nil
	} else {
		value = &types.StringC{Value: strconv.FormatInt(*data, 10)}
	}
	return p.Set(parser.Global, parser.GlobalSectionName, option, value)
}

func serializeBoolOption(p parser.Parser, option string, data bool) error {
	value := &types.Enabled{}
	if !data {
		value = nil
	}
	return p.Set(parser.Global, parser.GlobalSectionName, option, value)
}

func serializeOnOffOption(p parser.Parser, option, data string) error {
	var value *types.StringC
	switch data {
	case "enabled":
		value = &types.StringC{Value: "on"}
	case "disabled":
		value = &types.StringC{Value: "off"}
	default:
		value = nil
	}
	return p.Set(parser.Global, parser.GlobalSectionName, option, value)
}

func serializeAutoOnOffOption(p parser.Parser, option, data string) error {
	var value *types.StringC
	switch data {
	case "auto":
		value = &types.StringC{Value: "auto"}
	case "enabled":
		value = &types.StringC{Value: "on"}
	case "disabled":
		value = &types.StringC{Value: "off"}
	default:
		value = nil
	}
	return p.Set(parser.Global, parser.GlobalSectionName, option, value)
}

func serializeInt64Option(p parser.Parser, option string, data int64) error {
	value := &types.Int64C{
		Value: data,
	}
	if data == 0 {
		value = nil
	}
	return p.Set(parser.Global, parser.GlobalSectionName, option, value)
}

func serializeInt64POption(p parser.Parser, option string, data *int64) error {
	var value *types.Int64C
	if data == nil {
		value = nil
	} else {
		value = &types.Int64C{
			Value: *data,
		}
	}
	return p.Set(parser.Global, parser.GlobalSectionName, option, value)
}

func parseTuneOptions(p parser.Parser) (*models.GlobalTuneOptions, error) { //nolint:gocognit, gocyclo, cyclop
	options := &models.GlobalTuneOptions{}
	var intOption int64
	var intPOption *int64
	var boolOption bool
	var strOption string
	var err error

	intPOption, err = parseInt64POption(p, "tune.buffers.limit")
	if err != nil {
		return nil, err
	}
	options.BuffersLimit = intPOption

	intOption, err = parseInt64Option(p, "tune.buffers.reserve")
	if err != nil {
		return nil, err
	}
	options.BuffersReserve = intOption

	intOption, err = parseInt64Option(p, "tune.bufsize")
	if err != nil {
		return nil, err
	}
	options.Bufsize = intOption

	intOption, err = parseInt64Option(p, "tune.comp.maxlevel")
	if err != nil {
		return nil, err
	}
	options.CompMaxlevel = intOption

	boolOption, err = parseBoolOption(p, "tune.fail-alloc")
	if err != nil {
		return nil, err
	}
	options.FailAlloc = boolOption

	intOption, err = parseInt64Option(p, "tune.h2.header-table-size")
	if err != nil {
		return nil, err
	}
	options.H2HeaderTableSize = intOption

	intPOption, err = parseInt64POption(p, "tune.h2.initial-window-size")
	if err != nil {
		return nil, err
	}
	options.H2InitialWindowSize = intPOption

	intOption, err = parseInt64Option(p, "tune.h2.max-concurrent-streams")
	if err != nil {
		return nil, err
	}
	options.H2MaxConcurrentStreams = intOption

	intOption, err = parseInt64Option(p, "tune.h2.max-frame-size")
	if err != nil {
		return nil, err
	}
	options.H2MaxFrameSize = intOption

	intOption, err = parseInt64Option(p, "tune.http.cookielen")
	if err != nil {
		return nil, err
	}
	options.HTTPCookielen = intOption

	intOption, err = parseInt64Option(p, "tune.http.logurilen")
	if err != nil {
		return nil, err
	}
	options.HTTPLogurilen = intOption

	intOption, err = parseInt64Option(p, "tune.http.maxhdr")
	if err != nil {
		return nil, err
	}
	options.HTTPMaxhdr = intOption

	strOption, err = parseOnOffOption(p, "tune.idle-pool.shared")
	if err != nil {
		return nil, err
	}
	options.IdlePoolShared = strOption

	strOption, err = parseStringOption(p, "tune.idletimer")
	if err != nil {
		return nil, err
	}
	options.Idletimer = misc.ParseTimeout(strOption)

	strOption, err = parseOnOffOption(p, "tune.listener.multi-queue")
	if err != nil {
		return nil, err
	}
	options.ListenerMultiQueue = strOption

	intOption, err = parseInt64Option(p, "tune.lua.forced-yield")
	if err != nil {
		return nil, err
	}
	options.LuaForcedYield = intOption

	boolOption, err = parseBoolOption(p, "tune.lua.maxmem")
	if err != nil {
		return nil, err
	}
	options.LuaMaxmem = boolOption

	strOption, err = parseStringOption(p, "tune.lua.session-timeout")
	if err != nil {
		return nil, err
	}
	options.LuaSessionTimeout = misc.ParseTimeout(strOption)

	strOption, err = parseStringOption(p, "tune.lua.task-timeout")
	if err != nil {
		return nil, err
	}
	options.LuaTaskTimeout = misc.ParseTimeout(strOption)

	strOption, err = parseStringOption(p, "tune.lua.service-timeout")
	if err != nil {
		return nil, err
	}
	options.LuaServiceTimeout = misc.ParseTimeout(strOption)

	intOption, err = parseInt64Option(p, "tune.maxaccept")
	if err != nil {
		return nil, err
	}
	options.Maxaccept = intOption

	intOption, err = parseInt64Option(p, "tune.maxpollevents")
	if err != nil {
		return nil, err
	}
	options.Maxpollevents = intOption

	intOption, err = parseInt64Option(p, "tune.maxrewrite")
	if err != nil {
		return nil, err
	}
	options.Maxrewrite = intOption

	intPOption, err = parseInt64POption(p, "tune.pattern.cache-size")
	if err != nil {
		return nil, err
	}
	options.PatternCacheSize = intPOption

	intOption, err = parseInt64Option(p, "tune.pipesize")
	if err != nil {
		return nil, err
	}
	options.Pipesize = intOption

	intOption, err = parseInt64Option(p, "tune.pool-high-fd-ratio")
	if err != nil {
		return nil, err
	}
	options.PoolHighFdRatio = intOption

	intOption, err = parseInt64Option(p, "tune.pool-low-fd-ratio")
	if err != nil {
		return nil, err
	}
	options.PoolLowFdRatio = intOption

	intPOption, err = parseInt64POption(p, "tune.rcvbuf.client")
	if err != nil {
		return nil, err
	}
	options.RcvbufClient = intPOption

	intPOption, err = parseInt64POption(p, "tune.rcvbuf.server")
	if err != nil {
		return nil, err
	}
	options.RcvbufServer = intPOption

	intOption, err = parseInt64Option(p, "tune.recv_enough")
	if err != nil {
		return nil, err
	}
	options.RecvEnough = intOption

	intOption, err = parseInt64Option(p, "tune.runqueue-depth")
	if err != nil {
		return nil, err
	}
	options.RunqueueDepth = intOption

	strOption, err = parseOnOffOption(p, "tune.sched.low-latency")
	if err != nil {
		return nil, err
	}
	options.SchedLowLatency = strOption

	intPOption, err = parseInt64POption(p, "tune.sndbuf.client")
	if err != nil {
		return nil, err
	}
	options.SndbufClient = intPOption

	intPOption, err = parseInt64POption(p, "tune.sndbuf.server")
	if err != nil {
		return nil, err
	}
	options.SndbufServer = intPOption

	intPOption, err = parseInt64POption(p, "tune.ssl.cachesize")
	if err != nil {
		return nil, err
	}
	options.SslCachesize = intPOption

	boolOption, err = parseBoolOption(p, "tune.ssl.force-private-cache")
	if err != nil {
		return nil, err
	}
	options.SslForcePrivateCache = boolOption

	strOption, err = parseOnOffOption(p, "tune.ssl.keylog")
	if err != nil {
		return nil, err
	}
	options.SslKeylog = strOption

	strOption, err = parseStringOption(p, "tune.ssl.lifetime")
	if err != nil {
		return nil, err
	}
	options.SslLifetime = misc.ParseTimeout(strOption)

	intPOption, err = parseInt64POption(p, "tune.ssl.maxrecord")
	if err != nil {
		return nil, err
	}
	options.SslMaxrecord = intPOption

	intOption, err = parseInt64Option(p, "tune.ssl.default-dh-param")
	if err != nil {
		return nil, err
	}
	options.SslDefaultDhParam = intOption

	intOption, err = parseInt64Option(p, "tune.ssl.ssl-ctx-cache-size")
	if err != nil {
		return nil, err
	}
	options.SslCtxCacheSize = intOption

	intPOption, err = parseInt64POption(p, "tune.ssl.capture-buffer-size")
	if err != nil {
		return nil, err
	}
	options.SslCaptureBufferSize = intPOption

	strOption, err = parseStringOption(p, "tune.vars.global-max-size")
	if err != nil {
		return nil, err
	}
	options.VarsGlobalMaxSize = misc.ParseSize(strOption)

	strOption, err = parseStringOption(p, "tune.vars.proc-max-size")
	if err != nil {
		return nil, err
	}
	options.VarsProcMaxSize = misc.ParseSize(strOption)

	strOption, err = parseStringOption(p, "tune.vars.reqres-max-size")
	if err != nil {
		return nil, err
	}
	options.VarsReqresMaxSize = misc.ParseSize(strOption)

	strOption, err = parseStringOption(p, "tune.vars.sess-max-size")
	if err != nil {
		return nil, err
	}
	options.VarsSessMaxSize = misc.ParseSize(strOption)

	strOption, err = parseStringOption(p, "tune.vars.txn-max-size")
	if err != nil {
		return nil, err
	}
	options.VarsTxnMaxSize = misc.ParseSize(strOption)

	intOption, err = parseInt64Option(p, "tune.zlib.memlevel")
	if err != nil {
		return nil, err
	}
	options.ZlibMemlevel = intOption

	intOption, err = parseInt64Option(p, "tune.zlib.windowsize")
	if err != nil {
		return nil, err
	}
	options.ZlibWindowsize = intOption

	return options, nil
}

func parseStringOption(p parser.Parser, option string) (string, error) {
	data, err := p.Get(parser.Global, parser.GlobalSectionName, option)
	if err == nil {
		value, ok := data.(*types.StringC)
		if !ok {
			return "", misc.CreateTypeAssertError(option)
		}
		return value.Value, nil
	}
	if errors.Is(err, parser_errors.ErrFetch) {
		return "", nil
	}
	return "", err
}

func parseInt64Option(p parser.Parser, option string) (int64, error) {
	data, err := p.Get(parser.Global, parser.GlobalSectionName, option)
	if err == nil {
		value, ok := data.(*types.Int64C)
		if !ok {
			return 0, misc.CreateTypeAssertError(option)
		}
		return value.Value, nil
	}
	if errors.Is(err, parser_errors.ErrFetch) {
		return 0, nil
	}
	return 0, err
}

//nolint:nilnil
func parseInt64POption(p parser.Parser, option string) (*int64, error) {
	data, err := p.Get(parser.Global, parser.GlobalSectionName, option)
	if err == nil {
		value, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError(option)
		}
		return misc.Int64P(int(value.Value)), nil
	}
	if errors.Is(err, parser_errors.ErrFetch) {
		return nil, nil
	}
	return nil, err
}

func parseBoolOption(p parser.Parser, option string) (bool, error) {
	_, err := p.Get(parser.Global, parser.GlobalSectionName, option)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, parser_errors.ErrFetch) {
		return false, nil
	}
	return false, err
}

func parseOnOffOption(p parser.Parser, option string) (string, error) {
	data, err := p.Get(parser.Global, parser.GlobalSectionName, option)
	if err == nil {
		value, ok := data.(*types.StringC)
		if !ok {
			return "", misc.CreateTypeAssertError(option)
		}
		switch value.Value {
		case "on":
			return "enabled", nil
		case "off":
			return "disabled", nil
		default:
			return "", fmt.Errorf("unsupported value for %s: %s", option, value.Value)
		}
	}
	if errors.Is(err, parser_errors.ErrFetch) {
		return "", nil
	}
	return "", err
}

func parseAutoOnOffOption(p parser.Parser, option string) (string, error) {
	data, err := p.Get(parser.Global, parser.GlobalSectionName, option)
	if err == nil {
		value, ok := data.(*types.StringC)
		if !ok {
			return "", misc.CreateTypeAssertError(option)
		}
		switch value.Value {
		case "auto":
			return "auto", nil
		case "on":
			return "enabled", nil
		case "off":
			return "disabled", nil
		default:
			return "", fmt.Errorf("unsupported value for %s: %s", option, value.Value)
		}
	}
	if errors.Is(err, parser_errors.ErrFetch) {
		return "", nil
	}
	return "", err
}
