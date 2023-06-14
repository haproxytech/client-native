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

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v5"
	"github.com/haproxytech/config-parser/v5/common"
	parser_errors "github.com/haproxytech/config-parser/v5/errors"
	"github.com/haproxytech/config-parser/v5/params"
	"github.com/haproxytech/config-parser/v5/types"

	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
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

func ParseGlobalSection(p parser.Parser) (*models.Global, error) { //nolint:gocognit,gocyclo,cyclop,maintidx
	var anonkey *int64
	data, err := p.Get(parser.Global, parser.GlobalSectionName, "anonkey")
	if err == nil {
		anonkeyParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("anonkey")
		}
		anonkey = &anonkeyParser.Value
	}

	var clusterSecret string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "cluster-secret")
	if err == nil {
		csParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("cluster-secret")
		}
		clusterSecret = csParser.Value
	}

	var chroot string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "chroot")
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

	var node string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "node")
	if err == nil {
		nodeParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("node")
		}
		node = nodeParser.Value
	}

	var description string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "description")
	if err == nil {
		descriptionParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("description")
		}
		description = descriptionParser.Value
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "expose-experimental-directives")
	exposeExperimentalDirectives := true
	if errors.Is(err, parser_errors.ErrFetch) {
		exposeExperimentalDirectives = false
	}

	var grace *int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "grace")
	if err == nil {
		graceParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("grace")
		}
		grace = misc.ParseTimeout(graceParser.Value)
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "insecure-fork-wanted")
	insecureForkWanted := true
	if errors.Is(err, parser_errors.ErrFetch) {
		insecureForkWanted = false
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "insecure-setuid-wanted")
	insecureSetuidWanted := true
	if errors.Is(err, parser_errors.ErrFetch) {
		insecureSetuidWanted = false
	}

	var issuersChainPath string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "issuers-chain-path")
	if err == nil {
		issuersChainPathParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("issuers-chain-path")
		}
		issuersChainPath = issuersChainPathParser.Value
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "h2-workaround-bogus-websocket-clients")
	h2WorkaroundBogusWebsocketClients := true
	if errors.Is(err, parser_errors.ErrFetch) {
		h2WorkaroundBogusWebsocketClients = false
	}

	var luaLoadPerThread string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "lua-load-per-thread")
	if err == nil {
		luaLoadPerThreadParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("lua-load-per-thread")
		}
		luaLoadPerThread = luaLoadPerThreadParser.Value
	}

	var mworkerMaxReloads *int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "mworker-max-reloads")
	if errors.Is(err, parser_errors.ErrFetch) {
		mworkerMaxReloads = nil
	} else {
		mworkerMaxReloadsParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("mworker-max-reloads")
		}
		mworkerMaxReloads = &mworkerMaxReloadsParser.Value
	}

	var numaCPUMapping string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "numa-cpu-mapping")
	if err == nil {
		numaCPUMappingParser, ok := data.(*types.NumaCPUMapping)
		if !ok {
			return nil, misc.CreateTypeAssertError("numa-cpu-mapping")
		}
		if numaCPUMappingParser.NoOption {
			numaCPUMapping = "disabled"
		} else {
			numaCPUMapping = "enabled"
		}
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "pp2-never-send-local")
	pp2NeverSendLocal := true
	if errors.Is(err, parser_errors.ErrFetch) {
		pp2NeverSendLocal = false
	}

	var ulimitn int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ulimit-n")
	if err == nil {
		ulimitnParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("ulimit-n")
		}
		ulimitn = ulimitnParser.Value
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "set-dumpable")
	setDumpable := true
	if errors.Is(err, parser_errors.ErrFetch) {
		setDumpable = false
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "strict-limits")
	strictLimits := true
	if errors.Is(err, parser_errors.ErrFetch) {
		strictLimits = false
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
			ondiskMap := m
			cpuMap := &models.CPUMap{
				Process: &ondiskMap.Process,
				CPUSet:  &ondiskMap.CPUSet,
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

	httpClientResolversDisabled, err := parseOnOffOption(p, "httpclient.resolvers.disabled")
	if err != nil {
		return nil, err
	}

	httpClientResolversID, err := parseStringOption(p, "httpclient.resolvers.id")
	if err != nil {
		return nil, err
	}

	var httpClientResolversPrefer string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "httpclient.resolvers.prefer")
	if err == nil {
		httpClientResolversPreferParser, ok := data.(*types.HTTPClientResolversPrefer)
		if !ok {
			return nil, misc.CreateTypeAssertError("httpclient.resolvers.prefer")
		}
		httpClientResolversPrefer = httpClientResolversPreferParser.Type
	}

	httpClientSSLCaFile, err := parseStringOption(p, "httpclient.ssl.ca-file")
	if err != nil {
		return nil, err
	}

	var httpClientSSLVerify *string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "httpclient.ssl.verify")
	if err == nil {
		httpClientSSLVerifyParser, ok := data.(*types.HTTPClientSSLVerify)
		if !ok {
			return nil, misc.CreateTypeAssertError("httpclient.ssl.verify")
		}
		httpClientSSLVerify = &httpClientSSLVerifyParser.Type
	}

	preallocFD, err := parseBoolOption(p, "prealloc-fd")
	if err != nil {
		return nil, err
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

	var sslDefaultBindCurves string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-bind-curves")
	if err == nil {
		sslDefaultBindCurvesParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("ssl-default-bind-curves")
		}
		sslDefaultBindCurves = sslDefaultBindCurvesParser.Value
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-skip-self-issued-ca")
	sslSkipSelfIssuedCa := true
	if errors.Is(err, parser_errors.ErrFetch) {
		sslSkipSelfIssuedCa = false
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

	var sslBindSigalgs string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-bind-sigalgs")
	if err == nil {
		sslBindSigalgsParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("ssl-default-bind-sigalgs")
		}
		sslBindSigalgs = sslBindSigalgsParser.Value
	}

	var sslBindClientSigalgs string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-bind-client-sigalgs")
	if err == nil {
		sslBindClientSigalgsParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("ssl-default-bind-client-sigalgs")
		}
		sslBindClientSigalgs = sslBindClientSigalgsParser.Value
	}

	var sslDefaultServerCiphers string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-default-server-ciphers")
	if err == nil {
		sslDefaultServerCiphersParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("ssl-default-server-ciphers")
		}
		sslDefaultServerCiphers = sslDefaultServerCiphersParser.Value
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

	var sslModeAsync string
	data, _ = p.Get(parser.Global, parser.GlobalSectionName, "ssl-mode-async")
	if _, ok := data.(*types.SslModeAsync); ok {
		sslModeAsync = "enabled"
	}

	var sslDhParamFile string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-dh-param-file")
	if err == nil {
		sslDhParamFileParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("ssl-dh-param-file")
		}
		sslDhParamFile = sslDhParamFileParser.Value
	}

	var sslServerVerify string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-server-verify")
	if err == nil {
		sslServerVerifyParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("ssl-server-verify")
		}
		sslServerVerify = sslServerVerifyParser.Value
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

	maxSpreadChecksStr, err := parseStringOption(p, "max-spread-checks")
	if err != nil {
		return nil, err
	}
	maxSpreadChecks := misc.ParseTimeout(maxSpreadChecksStr)

	closeSpreadTimeStr, err := parseStringOption(p, "close-spread-time")
	if err != nil {
		return nil, err
	}
	closeSpreadTime := misc.ParseTimeout(closeSpreadTimeStr)

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

	var noQuic bool
	_, err = p.Get(parser.Global, parser.GlobalSectionName, "no-quic")
	if !errors.Is(err, parser_errors.ErrFetch) {
		noQuic = true
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

	var threadGroups int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "thread-groups")
	if err == nil {
		threadGroupsParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("thread-groups")
		}
		threadGroups = threadGroupsParser.Value
	}

	var statsMaxconn *int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "stats maxconn")
	if errors.Is(err, parser_errors.ErrFetch) {
		statsMaxconn = nil
	} else {
		statsMaxconnParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("stats maxconn")
		}
		statsMaxconn = &statsMaxconnParser.Value
	}

	var SSLLoadExtraFiles string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-load-extra-files")
	if err == nil {
		SSLLoadExtraFilesParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("ssl-load-extra-files")
		}
		SSLLoadExtraFiles = SSLLoadExtraFilesParser.Value
	}

	var threadGroupLines []*models.ThreadGroup
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "thread-group")
	if err == nil {
		items, ok := data.([]types.ThreadGroup)
		if !ok {
			return nil, misc.CreateTypeAssertError("thread-group")
		}
		for _, item := range items {
			g := item.Group
			nor := item.NumOrRange
			threadGroupLines = append(threadGroupLines, &models.ThreadGroup{
				Group:      &g,
				NumOrRange: &nor,
			})
		}
	}

	var sslEngines []*models.SslEngine
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "ssl-engine")
	if err == nil {
		items, ok := data.([]types.SslEngine)
		if !ok {
			return nil, misc.CreateTypeAssertError("ssl-engine")
		}
		for _, item := range items {
			name := item.Name
			algo := strings.Join(item.Algorithms, ",")
			sslEngines = append(sslEngines, &models.SslEngine{
				Name:       &name,
				Algorithms: &algo,
			})
		}
	}

	wurflOptions := models.GlobalWurflOptions{}

	var wurflDataFile string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "wurfl-data-file")
	if err == nil {
		wurflDataFileParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("wurfl-data-file")
		}
		wurflDataFile = wurflDataFileParser.Value
		wurflOptions.DataFile = wurflDataFile
	}

	var wurflInformationList string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "wurfl-information-list")
	if err == nil {
		wurflInformationListParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("wurfl-information-list")
		}
		wurflInformationList = wurflInformationListParser.Value
		wurflOptions.InformationList = wurflInformationList
	}

	var wurflInformationListSeparator string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "wurfl-information-list-separator")
	if err == nil {
		wurflInformationListSeparatorParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("wurfl-information-list-separator")
		}
		wurflInformationListSeparator = wurflInformationListSeparatorParser.Value
		wurflOptions.InformationListSeparator = wurflInformationListSeparator
	}

	var wurflPatchFile string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "wurfl-patch-file")
	if err == nil {
		wurflPatchFileParser, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError("wurfl-patch-file")
		}
		wurflPatchFile = wurflPatchFileParser.Value
		wurflOptions.PatchFile = wurflPatchFile
	}

	var wurflCacheSize int64
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "wurfl-cache-size")
	if err == nil {
		wurflCacheSizeParser, ok := data.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("wurfl-cache-size")
		}
		wurflCacheSize = wurflCacheSizeParser.Value
		wurflOptions.CacheSize = wurflCacheSize
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "quiet")
	quiet := true
	if errors.Is(err, parser_errors.ErrFetch) {
		quiet = false
	}

	_, err = p.Get(parser.Global, parser.GlobalSectionName, "zero-warning")
	zeroWarning := true
	if errors.Is(err, parser_errors.ErrFetch) {
		zeroWarning = false
	}

	var setVars []*models.SetVar
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "set-var")
	if err == nil {
		vars, ok := data.([]types.SetVar)
		if !ok {
			return nil, misc.CreateTypeAssertError("set-var")
		}
		for _, v := range vars {
			name := v.Name
			expr := v.Expr.String()
			setVar := &models.SetVar{
				Name: &name,
				Expr: &expr,
			}
			setVars = append(setVars, setVar)
		}
	}

	var setVarFormats []*models.SetVarFmt
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "set-var-fmt")
	if err == nil {
		formats, ok := data.([]types.SetVarFmt)
		if !ok {
			return nil, misc.CreateTypeAssertError("set-var-fmt")
		}
		for _, f := range formats {
			name := f.Name
			format := f.Format
			setVarFmt := &models.SetVarFmt{
				Name:   &name,
				Format: &format,
			}
			setVarFormats = append(setVarFormats, setVarFmt)
		}
	}

	var presetEnvs []*models.PresetEnv
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "presetenv")
	if err == nil {
		envs, ok := data.([]types.StringKeyValueC)
		if !ok {
			return nil, misc.CreateTypeAssertError("presetenv")
		}
		for _, e := range envs {
			env := &models.PresetEnv{
				Name:  &e.Key,
				Value: &e.Value,
			}
			presetEnvs = append(presetEnvs, env)
		}
	}

	var setEnvs []*models.SetEnv
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "setenv")
	if err == nil {
		envs, ok := data.([]types.StringKeyValueC)
		if !ok {
			return nil, misc.CreateTypeAssertError("setenv")
		}
		for _, e := range envs {
			ondiskEnv := e
			env := &models.SetEnv{
				Name:  &ondiskEnv.Key,
				Value: &ondiskEnv.Value,
			}
			setEnvs = append(setEnvs, env)
		}
	}

	var resetEnv string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "resetenv")
	if err == nil {
		resetEnvParser, ok := data.(*types.StringSliceC)
		if !ok {
			return nil, misc.CreateTypeAssertError("resetenv")
		}
		if len(resetEnvParser.Value) > 0 {
			resetEnv = strings.Join(resetEnvParser.Value, " ")
		}
	}

	var unsetEnv string
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "unsetenv")
	if err == nil {
		unsetEnvParser, ok := data.(*types.StringSliceC)
		if !ok {
			return nil, misc.CreateTypeAssertError("unsetenv")
		}
		if len(unsetEnvParser.Value) > 0 {
			unsetEnv = strings.Join(unsetEnvParser.Value, " ")
		}
	}

	tuneOptions, err := parseTuneOptions(p)
	if err != nil {
		return nil, err
	}

	deviceAtlasOptions, err := parseDeviceAtlasOptions(p)
	if err != nil {
		return nil, err
	}

	fiftyOneDegreesOptions, err := parseFiftyOneDegreesOptions(p)
	if err != nil {
		return nil, err
	}

	var defaultPath *models.GlobalDefaultPath
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "default-path")
	if err == nil {
		defaultPathParser, ok := data.(*types.DefaultPath)
		if !ok {
			return nil, misc.CreateTypeAssertError("default-path")
		}
		defaultPath = &models.GlobalDefaultPath{
			Type: defaultPathParser.Type,
			Path: defaultPathParser.Path,
		}
	}

	// deprecated option
	dhParam := int64(0)
	if tuneOptions != nil {
		dhParam = tuneOptions.SslDefaultDhParam
	}

	global := &models.Global{
		Anonkey:                           anonkey,
		PresetEnvs:                        presetEnvs,
		SetEnvs:                           setEnvs,
		Resetenv:                          resetEnv,
		Unsetenv:                          unsetEnv,
		UID:                               uid,
		User:                              user,
		Gid:                               gid,
		Group:                             group,
		ClusterSecret:                     clusterSecret,
		Chroot:                            chroot,
		Localpeer:                         localPeer,
		CaBase:                            caBase,
		CrtBase:                           crtBase,
		ServerStateBase:                   srvStateBase,
		ServerStateFile:                   srvStateFile,
		HardStopAfter:                     hardStop,
		Daemon:                            daemon,
		DefaultPath:                       defaultPath,
		MasterWorker:                      masterWorker,
		Maxconn:                           mConn,
		Nbproc:                            nbproc,
		Nbthread:                          nbthread,
		Pidfile:                           pidfile,
		RuntimeAPIs:                       rAPIs,
		StatsTimeout:                      statsTimeout,
		CPUMaps:                           cpuMaps,
		HttpclientResolversDisabled:       httpClientResolversDisabled,
		HttpclientResolversID:             httpClientResolversID,
		HttpclientResolversPrefer:         httpClientResolversPrefer,
		HttpclientSslCaFile:               httpClientSSLCaFile,
		HttpclientSslVerify:               httpClientSSLVerify,
		PreallocFd:                        preallocFD,
		SslDefaultBindCiphers:             sslBindCiphers,
		SslDefaultBindCiphersuites:        sslBindCiphersuites,
		SslDefaultBindCurves:              sslDefaultBindCurves,
		SslDefaultBindOptions:             sslBindOptions,
		SslDefaultServerCiphers:           sslDefaultServerCiphers,
		SslDefaultServerCiphersuites:      sslServerCiphersuites,
		SslDefaultServerOptions:           sslServerOptions,
		SslModeAsync:                      sslModeAsync,
		SslSkipSelfIssuedCa:               sslSkipSelfIssuedCa,
		TuneOptions:                       tuneOptions,
		TuneSslDefaultDhParam:             dhParam,
		ExternalCheck:                     externalCheck,
		LuaLoads:                          luaLoads,
		LuaPrependPath:                    luaPrependPath,
		LogSendHostname:                   globalLogSendHostName,
		H1CaseAdjusts:                     h1CaseAdjusts,
		H1CaseAdjustFile:                  h1CaseAdjustFile,
		BusyPolling:                       busyPolling,
		MaxSpreadChecks:                   maxSpreadChecks,
		CloseSpreadTime:                   closeSpreadTime,
		Maxconnrate:                       maxconnrate,
		Maxcomprate:                       maxcomprate,
		Maxcompcpuusage:                   maxcompcpuusage,
		Maxpipes:                          maxpipes,
		Maxsessrate:                       maxsessrate,
		Maxsslconn:                        maxsslconn,
		Maxsslrate:                        maxsslrate,
		Maxzlibmem:                        maxzlibmem,
		NoQuic:                            noQuic,
		Noepoll:                           noepoll,
		Nokqueue:                          nokqueue,
		Noevports:                         noevports,
		Nopoll:                            nopoll,
		Nosplice:                          nosplice,
		Nogetaddrinfo:                     nogetaddrinfo,
		Noreuseport:                       noreuseport,
		ProfilingTasks:                    profilingTasks,
		SpreadChecks:                      spreadChecks,
		ThreadGroups:                      threadGroups,
		StatsMaxconn:                      statsMaxconn,
		SslLoadExtraFiles:                 SSLLoadExtraFiles,
		ThreadGroupLines:                  threadGroupLines,
		Node:                              node,
		Description:                       description,
		ExposeExperimentalDirectives:      exposeExperimentalDirectives,
		Grace:                             grace,
		InsecureForkWanted:                insecureForkWanted,
		InsecureSetuidWanted:              insecureSetuidWanted,
		IssuersChainPath:                  issuersChainPath,
		H2WorkaroundBogusWebsocketClients: h2WorkaroundBogusWebsocketClients,
		LuaLoadPerThread:                  luaLoadPerThread,
		MworkerMaxReloads:                 mworkerMaxReloads,
		NumaCPUMapping:                    numaCPUMapping,
		Pp2NeverSendLocal:                 pp2NeverSendLocal,
		Ulimitn:                           ulimitn,
		SetDumpable:                       setDumpable,
		StrictLimits:                      strictLimits,
		WurflOptions:                      &wurflOptions,
		DeviceAtlasOptions:                deviceAtlasOptions,
		FiftyOneDegreesOptions:            fiftyOneDegreesOptions,
		Quiet:                             quiet,
		ZeroWarning:                       zeroWarning,
		SslEngines:                        sslEngines,
		SslDhParamFile:                    sslDhParamFile,
		SslServerVerify:                   sslServerVerify,
		SetVars:                           setVars,
		SetVarFmts:                        setVarFormats,
		SslDefaultBindSigalgs:             sslBindSigalgs,
		SslDefaultBindClientSigalgs:       sslBindClientSigalgs,
	}

	return global, nil
}

func SerializeGlobalSection(p parser.Parser, data *models.Global) error { //nolint:gocognit,gocyclo,cyclop,maintidx
	var pAnonkey *types.Int64C
	if data.Anonkey == nil {
		pAnonkey = nil
	} else {
		pAnonkey = &types.Int64C{
			Value: *data.Anonkey,
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "anonkey", pAnonkey); err != nil {
		return err
	}

	pClusterSecret := &types.StringC{
		Value: data.ClusterSecret,
	}
	if data.ClusterSecret == "" {
		pClusterSecret = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "cluster-secret", pClusterSecret); err != nil {
		return err
	}

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

	if err := serializeOnOffOption(p, "httpclient.resolvers.disabled", data.HttpclientResolversDisabled); err != nil {
		return err
	}

	if err := serializeStringOption(p, "httpclient.resolvers.id", data.HttpclientResolversID); err != nil {
		return err
	}

	pHTTPClientResolversPrefer := &types.HTTPClientResolversPrefer{
		Type: data.HttpclientResolversPrefer,
	}
	if data.HttpclientResolversPrefer == "" {
		pHTTPClientResolversPrefer = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "httpclient.resolvers.prefer", pHTTPClientResolversPrefer); err != nil {
		return err
	}

	if err := serializeStringOption(p, "httpclient.ssl.ca-file", data.HttpclientSslCaFile); err != nil {
		return err
	}

	var pHTTPClientSSLCaFile *types.HTTPClientSSLVerify
	if data.HttpclientSslVerify != nil {
		pHTTPClientSSLCaFile = &types.HTTPClientSSLVerify{
			Type: *data.HttpclientSslVerify,
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "httpclient.ssl.verify", pHTTPClientSSLCaFile); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "prealloc-fd", data.PreallocFd); err != nil {
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

	pSSLDefaultBindCurves := &types.StringC{
		Value: data.SslDefaultBindCurves,
	}
	if data.SslDefaultBindCurves == "" {
		pSSLDefaultBindCurves = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-default-bind-curves", pSSLDefaultBindCurves); err != nil {
		return err
	}

	pSslSkipSelfIssuedCa := &types.Enabled{}
	if !data.SslSkipSelfIssuedCa {
		pSslSkipSelfIssuedCa = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-skip-self-issued-ca", pSslSkipSelfIssuedCa); err != nil {
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

	pSSLBindSigalgs := &types.StringC{
		Value: data.SslDefaultBindSigalgs,
	}
	if data.SslDefaultBindSigalgs == "" {
		pSSLBindSigalgs = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-default-bind-sigalgs", pSSLBindSigalgs); err != nil {
		return err
	}

	pSSLBindClientSigalgs := &types.StringC{
		Value: data.SslDefaultBindClientSigalgs,
	}
	if data.SslDefaultBindClientSigalgs == "" {
		pSSLBindClientSigalgs = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-default-bind-client-sigalgs", pSSLBindClientSigalgs); err != nil {
		return err
	}

	pSSLDefaultServerCiphers := &types.StringC{
		Value: data.SslDefaultServerCiphers,
	}
	if data.SslDefaultServerCiphers == "" {
		pSSLDefaultServerCiphers = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-default-server-ciphers", pSSLDefaultServerCiphers); err != nil {
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

	pSSLDhParamFile := &types.StringC{
		Value: data.SslDhParamFile,
	}
	if data.SslDhParamFile == "" {
		pSSLDhParamFile = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-dh-param-file", pSSLDhParamFile); err != nil {
		return err
	}

	pSSLServerVerify := &types.StringC{
		Value: data.SslServerVerify,
	}
	if data.SslServerVerify == "" {
		pSSLServerVerify = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-server-verify", pSSLServerVerify); err != nil {
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

	if err := serializeTimeoutSizeOption(p, "max-spread-checks", data.MaxSpreadChecks); err != nil {
		return err
	}

	if err := serializeTimeoutSizeOption(p, "close-spread-time", data.CloseSpreadTime); err != nil {
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

	noQuic := &types.Enabled{}
	if !data.NoQuic {
		noQuic = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "no-quic", noQuic); err != nil {
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

	node := &types.StringC{Value: data.Node}
	if data.Node == "" {
		node = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "node", node); err != nil {
		return err
	}

	description := &types.StringC{Value: data.Description}
	if data.Description == "" {
		description = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "description", description); err != nil {
		return err
	}

	exposeExperimentalDirectives := &types.Enabled{}
	if !data.ExposeExperimentalDirectives {
		exposeExperimentalDirectives = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "expose-experimental-directives", exposeExperimentalDirectives); err != nil {
		return err
	}

	var grace *types.StringC
	if data.Grace != nil {
		grace = &types.StringC{Value: strconv.FormatInt(*data.Grace, 10)}
	} else {
		grace = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "grace", grace); err != nil {
		return err
	}

	insecureForkWanted := &types.Enabled{}
	if !data.InsecureForkWanted {
		insecureForkWanted = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "insecure-fork-wanted", insecureForkWanted); err != nil {
		return err
	}

	insecureSetuidWanted := &types.Enabled{}
	if !data.InsecureForkWanted {
		insecureSetuidWanted = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "insecure-setuid-wanted", insecureSetuidWanted); err != nil {
		return err
	}

	issuersChainPath := &types.StringC{Value: data.IssuersChainPath}
	if data.IssuersChainPath == "" {
		issuersChainPath = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "issuers-chain-path", issuersChainPath); err != nil {
		return err
	}

	workaround := &types.Enabled{}
	if !data.H2WorkaroundBogusWebsocketClients {
		workaround = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "h2-workaround-bogus-websocket-clients", workaround); err != nil {
		return err
	}

	luaLoadPerThread := &types.StringC{Value: data.LuaLoadPerThread}
	if data.LuaLoadPerThread == "" {
		luaLoadPerThread = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "lua-load-per-thread", luaLoadPerThread); err != nil {
		return err
	}

	var mworkerMaxReloads *types.Int64C
	if data.MworkerMaxReloads == nil {
		mworkerMaxReloads = nil
	} else {
		mworkerMaxReloads = &types.Int64C{
			Value: *data.MworkerMaxReloads,
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "mworker-max-reloads", mworkerMaxReloads); err != nil {
		return err
	}

	numaCPUMapping := &types.NumaCPUMapping{}
	if data.NumaCPUMapping == "" {
		numaCPUMapping = nil
	} else if data.NumaCPUMapping == "disabled" {
		numaCPUMapping.NoOption = true
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "numa-cpu-mapping", numaCPUMapping); err != nil {
		return err
	}

	neverSendLocal := &types.Enabled{}
	if !data.Pp2NeverSendLocal {
		neverSendLocal = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "pp2-never-send-local", neverSendLocal); err != nil {
		return err
	}

	ulimitN := &types.Int64C{
		Value: data.Ulimitn,
	}
	if data.Ulimitn == 0 {
		ulimitN = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ulimit-n", ulimitN); err != nil {
		return err
	}

	setDumpable := &types.Enabled{}
	if !data.SetDumpable {
		setDumpable = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "set-dumpable", setDumpable); err != nil {
		return err
	}

	strictLimits := &types.Enabled{}
	if !data.StrictLimits {
		strictLimits = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "strict-limits", strictLimits); err != nil {
		return err
	}

	threadGroups := &types.Int64C{
		Value: data.ThreadGroups,
	}
	if data.ThreadGroups == 0 {
		threadGroups = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "thread-groups", threadGroups); err != nil {
		return err
	}

	var statsMaxconn *types.Int64C
	if data.StatsMaxconn == nil {
		statsMaxconn = nil
	} else {
		statsMaxconn = &types.Int64C{
			Value: *data.StatsMaxconn,
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "stats maxconn", statsMaxconn); err != nil {
		return err
	}

	SSLLoadExtraFiles := &types.StringC{Value: data.SslLoadExtraFiles}
	if data.SslLoadExtraFiles == "" {
		SSLLoadExtraFiles = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-load-extra-files", SSLLoadExtraFiles); err != nil {
		return err
	}

	quiet := &types.Enabled{}
	if !data.Quiet {
		quiet = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "quiet", quiet); err != nil {
		return err
	}

	zeroWarning := &types.Enabled{}
	if !data.ZeroWarning {
		zeroWarning = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "zero-warning", zeroWarning); err != nil {
		return err
	}

	threadGroupLines := []types.ThreadGroup{}
	if data.ThreadGroupLines != nil && len(data.ThreadGroupLines) > 0 {
		for _, threadGroupLine := range data.ThreadGroupLines {
			if threadGroupLine != nil {
				tgl := types.ThreadGroup{
					Group:      *threadGroupLine.Group,
					NumOrRange: *threadGroupLine.NumOrRange,
				}
				threadGroupLines = append(threadGroupLines, tgl)
			}
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "thread-group", threadGroupLines); err != nil {
		return err
	}

	sslEngines := []types.SslEngine{}
	if data.SslEngines != nil && len(data.SslEngines) > 0 {
		for _, sslEngine := range data.SslEngines {
			if sslEngine != nil {
				se := types.SslEngine{
					Name:       *sslEngine.Name,
					Algorithms: strings.Split(*sslEngine.Algorithms, ","),
				}
				sslEngines = append(sslEngines, se)
			}
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-engine", sslEngines); err != nil {
		return err
	}

	setVars := []types.SetVar{}
	if data.SetVars != nil && len(data.SetVars) > 0 {
		for _, setVar := range data.SetVars {
			if setVar != nil {
				sv := types.SetVar{
					Name: *setVar.Name,
					Expr: common.Expression{Expr: strings.Split(*setVar.Expr, " ")},
				}
				setVars = append(setVars, sv)
			}
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "set-var", setVars); err != nil {
		return err
	}

	setVarFmts := []types.SetVarFmt{}
	if data.SetVarFmts != nil && len(data.SetVarFmts) > 0 {
		for _, setVarFmt := range data.SetVarFmts {
			if setVarFmt != nil {
				svf := types.SetVarFmt{
					Name:   *setVarFmt.Name,
					Format: *setVarFmt.Format,
				}
				setVarFmts = append(setVarFmts, svf)
			}
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "set-var-fmt", setVarFmts); err != nil {
		return err
	}

	presetEnvs := []types.StringKeyValueC{}
	if data.PresetEnvs != nil && len(data.PresetEnvs) > 0 {
		for _, presetEnv := range data.PresetEnvs {
			if presetEnv != nil {
				env := types.StringKeyValueC{
					Key:   *presetEnv.Name,
					Value: *presetEnv.Value,
				}
				presetEnvs = append(presetEnvs, env)
			}
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "presetenv", presetEnvs); err != nil {
		return err
	}

	setEnvs := []types.StringKeyValueC{}
	if data.SetEnvs != nil && len(data.SetEnvs) > 0 {
		for _, presetEnv := range data.SetEnvs {
			if presetEnv != nil {
				env := types.StringKeyValueC{
					Key:   *presetEnv.Name,
					Value: *presetEnv.Value,
				}
				setEnvs = append(setEnvs, env)
			}
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "setenv", setEnvs); err != nil {
		return err
	}

	resetenv := &types.StringSliceC{}
	if data.Resetenv == "" {
		resetenv = nil
	} else {
		resetenv.Value = strings.Split(data.Resetenv, " ")
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "resetenv", resetenv); err != nil {
		return err
	}

	unsetenv := &types.StringSliceC{}
	if data.Unsetenv == "" {
		unsetenv = nil
	} else {
		unsetenv.Value = strings.Split(data.Unsetenv, " ")
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "unsetenv", unsetenv); err != nil {
		return err
	}

	if data.WurflOptions == nil {
		data.WurflOptions = &models.GlobalWurflOptions{}
	}

	if err := serializeWurflOptions(p, data.WurflOptions); err != nil {
		return err
	}

	if err := serializeDeviceAtlasOptions(p, data.DeviceAtlasOptions); err != nil {
		return err
	}

	if err := serializeFiftyOneDegreesOptions(p, data.FiftyOneDegreesOptions); err != nil {
		return err
	}

	defaultPath := &types.DefaultPath{}
	if data.DefaultPath != nil {
		defaultPath.Type = data.DefaultPath.Type
		defaultPath.Path = data.DefaultPath.Path
	} else {
		defaultPath = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "default-path", defaultPath); err != nil {
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

func serializeWurflOptions(p parser.Parser, options *models.GlobalWurflOptions) error {
	if options == nil {
		return nil
	}
	if err := serializeStringOption(p, "wurfl-data-file", options.DataFile); err != nil {
		return err
	}
	if err := serializeStringOption(p, "wurfl-information-list", options.InformationList); err != nil {
		return err
	}
	if err := serializeStringOption(p, "wurfl-information-list-separator", options.InformationListSeparator); err != nil {
		return err
	}
	if err := serializeStringOption(p, "wurfl-patch-file", options.PatchFile); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "wurfl-cache-size", options.CacheSize); err != nil {
		return err
	}
	return nil
}

func serializeDeviceAtlasOptions(p parser.Parser, options *models.GlobalDeviceAtlasOptions) error {
	if options == nil {
		return nil
	}
	if err := serializeStringOption(p, "deviceatlas-json-file", options.JSONFile); err != nil {
		return err
	}
	if err := serializeStringOption(p, "deviceatlas-log-level", options.LogLevel); err != nil {
		return err
	}
	if err := serializeStringOption(p, "deviceatlas-separator", options.Separator); err != nil {
		return err
	}
	if err := serializeStringOption(p, "deviceatlas-properties-cookie", options.PropertiesCookie); err != nil {
		return err
	}
	return nil
}

func serializeFiftyOneDegreesOptions(p parser.Parser, options *models.GlobalFiftyOneDegreesOptions) error {
	if options == nil {
		return nil
	}
	if err := serializeStringOption(p, "51degrees-data-file", options.DataFile); err != nil {
		return err
	}
	if err := serializeStringOption(p, "51degrees-property-name-list", options.PropertyNameList); err != nil {
		return err
	}
	if err := serializeStringOption(p, "51degrees-property-separator", options.PropertySeparator); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "51degrees-cache-size", options.CacheSize); err != nil {
		return err
	}
	return nil
}

func serializeTuneOptions(p parser.Parser, options *models.GlobalTuneOptions) error { //nolint:gocognit,gocyclo,cyclop,maintidx
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
	if err := serializeListenerDefaultShards(p, "tune.listener.default-shards", options.ListenerDefaultShards); err != nil {
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
	if err := serializeTimeoutSizeOption(p, "tune.lua.burst-timeout", options.LuaBurstTimeout); err != nil {
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
	if err := serializeInt64POption(p, "tune.memory.hot-size", options.MemoryHotSize); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.pattern.cache-size", options.PatternCacheSize); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.peers.max-updates-at-once", options.PeersMaxUpdatesAtOnce); err != nil {
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
	if err := serializeInt64POption(p, "tune.quic.frontend.conn-tx-buffers.limit", options.QuicFrontendConnTcBuffersLimit); err != nil {
		return err
	}
	if err := serializeTimeoutSizeOption(p, "tune.quic.frontend.max-idle-timeout", options.QuicFrontendMaxIdleTimeout); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.quic.frontend.max-streams-bidi", options.QuicFrontendMaxStreamsBidi); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.quic.max-frame-loss", options.QuicMaxFrameLoss); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.quic.retry-threshold", options.QuicRetryThreshold); err != nil {
		return err
	}
	value := &types.QuicSocketOwner{Owner: options.QuicSocketOwner}
	if options.QuicSocketOwner == "" {
		value = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "tune.quic.socket-owner", value); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.zlib.memlevel", options.ZlibMemlevel); err != nil {
		return err
	}
	if err := serializeOnOffOption(p, "tune.fd.edge-triggered", options.FdEdgeTriggered); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.h2.be.initial-window-size", options.H2BeInitialWindowSize); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.h2.be.max-concurrent-streams", options.H2BeMaxConcurrentStreams); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.h2.fe.initial-window-size", options.H2FeInitialWindowSize); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.h2.fe.max-concurrent-streams", options.H2FeMaxConcurrentStreams); err != nil {
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

func serializeListenerDefaultShards(p parser.Parser, option, data string) error {
	var value *types.StringC
	switch data {
	case models.GlobalTuneOptionsListenerDefaultShardsByDashProcess:
		value = &types.StringC{Value: models.GlobalTuneOptionsListenerDefaultShardsByDashProcess}
	case models.GlobalTuneOptionsListenerDefaultShardsByDashThread:
		value = &types.StringC{Value: models.GlobalTuneOptionsListenerDefaultShardsByDashThread}
	case models.GlobalTuneOptionsListenerDefaultShardsByDashGroup:
		value = &types.StringC{Value: models.GlobalTuneOptionsListenerDefaultShardsByDashGroup}
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

func serializeStringOption(p parser.Parser, option string, data string) error {
	value := &types.StringC{Value: data}
	if data == "" {
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

func parseDeviceAtlasOptions(p parser.Parser) (*models.GlobalDeviceAtlasOptions, error) {
	options := &models.GlobalDeviceAtlasOptions{}
	var option string
	var err error

	option, err = parseStringOption(p, "deviceatlas-json-file")
	if err != nil {
		return nil, err
	}
	options.JSONFile = option

	option, err = parseStringOption(p, "deviceatlas-log-level")
	if err != nil {
		return nil, err
	}
	options.LogLevel = option

	option, err = parseStringOption(p, "deviceatlas-separator")
	if err != nil {
		return nil, err
	}
	options.Separator = option

	option, err = parseStringOption(p, "deviceatlas-properties-cookie")
	if err != nil {
		return nil, err
	}
	options.PropertiesCookie = option

	return options, nil
}

func parseFiftyOneDegreesOptions(p parser.Parser) (*models.GlobalFiftyOneDegreesOptions, error) {
	options := &models.GlobalFiftyOneDegreesOptions{}
	var option string
	var optionInt int64
	var err error

	option, err = parseStringOption(p, "51degrees-data-file")
	if err != nil {
		return nil, err
	}
	options.DataFile = option

	option, err = parseStringOption(p, "51degrees-property-name-list")
	if err != nil {
		return nil, err
	}
	options.PropertyNameList = option

	option, err = parseStringOption(p, "51degrees-property-separator")
	if err != nil {
		return nil, err
	}
	options.PropertySeparator = option

	optionInt, err = parseInt64Option(p, "51degrees-cache-size")
	if err != nil {
		return nil, err
	}
	options.CacheSize = optionInt

	return options, nil
}

func parseTuneOptions(p parser.Parser) (*models.GlobalTuneOptions, error) { //nolint:gocognit, gocyclo, cyclop,maintidx
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

	strOption, err = parseListenerDefaultShards(p, "tune.listener.default-shards")
	if err != nil {
		return nil, err
	}
	options.ListenerDefaultShards = strOption

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

	strOption, err = parseStringOption(p, "tune.lua.burst-timeout")
	if err != nil {
		return nil, err
	}
	options.LuaBurstTimeout = misc.ParseTimeout(strOption)

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

	intPOption, err = parseInt64POption(p, "tune.memory.hot-size")
	if err != nil {
		return nil, err
	}
	options.MemoryHotSize = intPOption

	intPOption, err = parseInt64POption(p, "tune.pattern.cache-size")
	if err != nil {
		return nil, err
	}
	options.PatternCacheSize = intPOption

	intOption, err = parseInt64Option(p, "tune.peers.max-updates-at-once")
	if err != nil {
		return nil, err
	}
	options.PeersMaxUpdatesAtOnce = intOption

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

	intPOption, err = parseInt64POption(p, "tune.ssl.ocsp-update.maxdelay")
	if err != nil {
		return nil, err
	}
	options.SslOcspUpdateMaxDelay = intPOption

	intPOption, err = parseInt64POption(p, "tune.ssl.ocsp-update.mindelay")
	if err != nil {
		return nil, err
	}
	options.SslOcspUpdateMinDelay = intPOption

	intPOption, err = parseInt64POption(p, "tune.stick-counters")
	if err != nil {
		return nil, err
	}
	options.StickCounters = intPOption

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

	intPOption, err = parseInt64POption(p, "tune.quic.frontend.conn-tx-buffers.limit")
	if err != nil {
		return nil, err
	}
	options.QuicFrontendConnTcBuffersLimit = intPOption

	strOption, err = parseStringOption(p, "tune.quic.frontend.max-idle-timeout")
	if err != nil {
		return nil, err
	}
	options.QuicFrontendMaxIdleTimeout = misc.ParseTimeout(strOption)

	intPOption, err = parseInt64POption(p, "tune.quic.frontend.max-streams-bidi")
	if err != nil {
		return nil, err
	}
	options.QuicFrontendMaxStreamsBidi = intPOption

	intPOption, err = parseInt64POption(p, "tune.quic.max-frame-loss")
	if err != nil {
		return nil, err
	}
	options.QuicMaxFrameLoss = intPOption

	intPOption, err = parseInt64POption(p, "tune.quic.retry-threshold")
	if err != nil {
		return nil, err
	}
	options.QuicRetryThreshold = intPOption

	so, err := p.Get(parser.Global, parser.GlobalSectionName, "tune.quic.socket-owner")
	if err == nil {
		value, ok := so.(*types.QuicSocketOwner)
		if !ok {
			return nil, misc.CreateTypeAssertError("tune.quic.socket-owner")
		}
		options.QuicSocketOwner = value.Owner
	}

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

	strOption, err = parseOnOffOption(p, "tune.fd.edge-triggered")
	if err != nil {
		return nil, err
	}
	options.FdEdgeTriggered = strOption

	intOption, err = parseInt64Option(p, "tune.h2.be.initial-window-size")
	if err != nil {
		return nil, err
	}
	options.H2BeInitialWindowSize = intOption

	intOption, err = parseInt64Option(p, "tune.h2.be.max-concurrent-streams")
	if err != nil {
		return nil, err
	}
	options.H2BeMaxConcurrentStreams = intOption

	intOption, err = parseInt64Option(p, "tune.h2.fe.initial-window-size")
	if err != nil {
		return nil, err
	}
	options.H2FeInitialWindowSize = intOption

	intOption, err = parseInt64Option(p, "tune.h2.fe.max-concurrent-streams")
	if err != nil {
		return nil, err
	}
	options.H2FeMaxConcurrentStreams = intOption

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

func parseListenerDefaultShards(p parser.Parser, option string) (string, error) {
	data, err := p.Get(parser.Global, parser.GlobalSectionName, option)
	if err == nil {
		value, ok := data.(*types.StringC)
		if !ok {
			return "", misc.CreateTypeAssertError(option)
		}
		switch value.Value {
		case models.GlobalTuneOptionsListenerDefaultShardsByDashProcess:
			return models.GlobalTuneOptionsListenerDefaultShardsByDashProcess, nil
		case models.GlobalTuneOptionsListenerDefaultShardsByDashThread:
			return models.GlobalTuneOptionsListenerDefaultShardsByDashThread, nil
		case models.GlobalTuneOptionsListenerDefaultShardsByDashGroup:
			return models.GlobalTuneOptionsListenerDefaultShardsByDashGroup, nil
		default:
			return "", fmt.Errorf("unsupported value for %s: %s", option, value.Value)
		}
	}
	if errors.Is(err, parser_errors.ErrFetch) {
		return "", nil
	}
	return "", err
}
