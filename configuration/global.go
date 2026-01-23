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
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/parsers"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/common"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/params"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/configuration/options"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
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

// ValidateGlobalSection performs validation of the Global section that cannot be done by swagger 2.0
func ValidateGlobalSection(data *models.Global) error {
	for i, cpuSet := range data.CPUSets {
		if cpuSet != nil {
			if cpuSet.Directive == nil {
				return fmt.Errorf("cpu_set.%d.directive must not be empty", i)
			}
			if *cpuSet.Directive == parsers.CPUSetResetDirective && len(cpuSet.Set) > 0 {
				return fmt.Errorf("cpu_set.%d.set must be empty when directive is %s", i, parsers.CPUSetResetDirective)
			}
			if *cpuSet.Directive != parsers.CPUSetResetDirective && len(cpuSet.Set) == 0 {
				return fmt.Errorf("cpu_set.%d.set must not be empty when directive is %s", i, *cpuSet.Directive)
			}
		}
	}
	return nil
}

// PushGlobalConfiguration pushes a Global config struct to global
// config file
func (c *client) PushGlobalConfiguration(data *models.Global, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}

		if err := ValidateGlobalSection(data); err != nil {
			return NewConfError(ErrValidationError, err.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if err := SerializeGlobalSection(p, data, &c.ConfigurationOptions); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func parseCPUMaps(p parser.Parser) ([]*models.CPUMap, error) {
	var cpuMaps []*models.CPUMap
	data, err := p.Get(parser.Global, parser.GlobalSectionName, "cpu-map")
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
	return cpuMaps, nil
}

func parseCPUSets(p parser.Parser) ([]*models.CPUSet, error) {
	var cpuSet []*models.CPUSet
	d, err := p.Get(parser.Global, parser.GlobalSectionName, "cpu-set")
	if err == nil {
		cpuSets, ok := d.([]types.CPUSet)
		if !ok {
			return nil, misc.CreateTypeAssertError("cpu-set")
		}
		for _, c := range cpuSets {
			directive := c.Directive
			cpuSet = append(cpuSet, &models.CPUSet{
				Directive: &directive,
				Set:       c.Set,
			})
		}
	}
	return cpuSet, nil
}

func parseH1CaseAdjusts(p parser.Parser) ([]*models.H1CaseAdjust, error) {
	var h1CaseAdjusts []*models.H1CaseAdjust
	data, err := p.Get(parser.Global, parser.GlobalSectionName, "h1-case-adjust")
	if err == nil {
		cases, ok := data.([]types.H1CaseAdjust)
		if !ok {
			return h1CaseAdjusts, misc.CreateTypeAssertError("h1-case-adjust")
		}
		for _, c := range cases {
			from := c.From
			to := c.To
			h1CaseAdjusts = append(h1CaseAdjusts, &models.H1CaseAdjust{From: &from, To: &to})
		}
	}
	return h1CaseAdjusts, nil
}

func parseRuntimeAPIs(p parser.Parser) ([]*models.RuntimeAPI, error) {
	var rAPIs []*models.RuntimeAPI
	data, err := p.Get(parser.Global, parser.GlobalSectionName, "stats socket")
	if err == nil {
		sockets, ok := data.([]types.Socket)
		if !ok {
			return nil, misc.CreateTypeAssertError("stats socket")
		}
		for _, s := range sockets {
			p := s.Path
			rAPI := &models.RuntimeAPI{Address: &p}
			rAPI.BindParams, rAPI.Name = parseBindParams(s.Params)
			rAPIs = append(rAPIs, rAPI)
		}
	}
	return rAPIs, nil
}

func parseSetVars(p parser.Parser) ([]*models.SetVar, error) {
	var setVars []*models.SetVar
	data, err := p.Get(parser.Global, parser.GlobalSectionName, "set-var")
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
	return setVars, nil
}

func parseSetVarFmts(p parser.Parser) ([]*models.SetVarFmt, error) {
	var setVarFormats []*models.SetVarFmt
	data, err := p.Get(parser.Global, parser.GlobalSectionName, "set-var-fmt")
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

	return setVarFormats, nil
}

func parseThreadGroupLines(p parser.Parser) ([]*models.ThreadGroup, error) {
	var threadGroupLines []*models.ThreadGroup
	data, err := p.Get(parser.Global, parser.GlobalSectionName, "thread-group")
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
	return threadGroupLines, nil
}

func parseDefaultPath(p parser.Parser) (*models.GlobalDefaultPath, error) {
	var defaultPath *models.GlobalDefaultPath
	data, err := p.Get(parser.Global, parser.GlobalSectionName, "default-path")
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
	return defaultPath, nil
}

func parseLuaOptions(p parser.Parser) (*models.LuaOptions, error) {
	options := &models.LuaOptions{}
	isEmpty := true

	luaLoadPerThread, err := parseStringOption(p, "lua-load-per-thread")
	if err != nil {
		return nil, err
	}
	if luaLoadPerThread != "" {
		isEmpty = false
		options.LoadPerThread = luaLoadPerThread
	}

	var luaPrependPath []*models.LuaPrependPath
	cpLuaPrependPath, err := p.Get(parser.Global, parser.GlobalSectionName, "lua-prepend-path")
	if err == nil {
		lpp, ok := cpLuaPrependPath.([]types.LuaPrependPath)
		if !ok {
			return nil, misc.CreateTypeAssertError("lua-prepend-path")
		}
		for _, l := range lpp {
			path := l.Path
			luaPrependPath = append(luaPrependPath, &models.LuaPrependPath{Path: &path, Type: l.Type})
		}
	}
	if luaPrependPath != nil {
		options.PrependPath = luaPrependPath
		isEmpty = false
	}

	var luaLoads []*models.LuaLoad
	cpLuaLoads, err := p.Get(parser.Global, parser.GlobalSectionName, "lua-load")
	if err == nil {
		luas, ok := cpLuaLoads.([]types.LuaLoad)
		if !ok {
			return nil, misc.CreateTypeAssertError("lua-load")
		}
		for _, l := range luas {
			file := l.File
			luaLoads = append(luaLoads, &models.LuaLoad{File: &file})
		}
	}
	if luaLoads != nil {
		options.Loads = luaLoads
		isEmpty = false
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func parseOcspUpdateOptions(p parser.Parser) (*models.OcspUpdateOptions, error) {
	options := &models.OcspUpdateOptions{}
	isEmpty := true
	ocspUpdateDisable, err := parseOnOffOption(p, "ocsp-update.disable")
	if err != nil {
		return nil, err
	}
	switch ocspUpdateDisable {
	case "disabled":
		isEmpty = false
		options.Disable = misc.BoolP(false)
	case "enabled":
		isEmpty = false
		options.Disable = misc.BoolP(true)
	default:
		options.Disable = nil
	}

	minDelayP, err := parseInt64POption(p, "ocsp-update.mindelay")
	if err != nil {
		return nil, err
	}
	if minDelayP != nil {
		isEmpty = false
		options.Mindelay = minDelayP
	}

	maxDelayP, err := parseInt64POption(p, "ocsp-update.maxdelay")
	if err != nil {
		return nil, err
	}
	if maxDelayP != nil {
		isEmpty = false
		options.Maxdelay = maxDelayP
	}

	addressPort, err := parseStringOption(p, "ocsp-update.httpproxy")
	if err != nil {
		return nil, err
	}
	address, port := ParseAddress(addressPort)
	if address != "" {
		isEmpty = false
		options.Httpproxy = &models.OcspUpdateOptionsHttpproxy{}
		options.Httpproxy.Address = address
		options.Httpproxy.Port = port
	}

	mode, err := parseOnOffOption(p, "ocsp-update.mode")
	if err != nil {
		return nil, err
	}
	if mode != "" {
		isEmpty = false
		options.Mode = mode
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}

	return options, nil
}

func parsePerformanceOptions(p parser.Parser) (*models.PerformanceOptions, error) { //nolint:gocognit,gocyclo,cyclop,maintidx
	options := &models.PerformanceOptions{}
	isEmpty := true
	busyPolling, err := parseBoolOption(p, "busy-polling")
	if err != nil {
		return nil, err
	}
	if busyPolling {
		isEmpty = false
		options.BusyPolling = busyPolling
	}

	noktls, err := parseBoolOption(p, "noktls")
	if err != nil {
		return nil, err
	}
	if noktls {
		isEmpty = false
		options.Noktls = noktls
	}

	maxSpreadChecks, err := parseTimeoutOption(p, "max-spread-checks")
	if err != nil {
		return nil, err
	}
	if maxSpreadChecks != nil {
		isEmpty = false
		options.MaxSpreadChecks = maxSpreadChecks
	}

	maxcompcpuusage, err := parseInt64Option(p, "maxcompcpuusage")
	if err != nil {
		return nil, err
	}
	if maxcompcpuusage != 0 {
		isEmpty = false
		options.Maxcompcpuusage = maxcompcpuusage
	}

	maxconnrate, err := parseInt64Option(p, "maxconnrate")
	if err != nil {
		return nil, err
	}
	if maxconnrate != 0 {
		isEmpty = false
		options.Maxconnrate = maxconnrate
	}

	maxcomprate, err := parseInt64Option(p, "maxcomprate")
	if err != nil {
		return nil, err
	}
	if maxcomprate != 0 {
		isEmpty = false
		options.Maxcomprate = maxcomprate
	}

	maxpipes, err := parseInt64Option(p, "maxpipes")
	if err != nil {
		return nil, err
	}
	if maxpipes != 0 {
		isEmpty = false
		options.Maxpipes = maxpipes
	}

	maxsessrate, err := parseInt64Option(p, "maxsessrate")
	if err != nil {
		return nil, err
	}
	if maxsessrate != 0 {
		isEmpty = false
		options.Maxsessrate = maxsessrate
	}

	maxconn, err := parseInt64Option(p, "maxconn")
	if err != nil {
		return nil, err
	}
	if maxconn != 0 {
		isEmpty = false
		options.Maxconn = maxconn
	}

	maxzlibmem, err := parseInt64Option(p, "maxzlibmem")
	if err != nil {
		return nil, err
	}
	if maxzlibmem != 0 {
		isEmpty = false
		options.Maxzlibmem = maxzlibmem
	}

	noepoll, err := parseBoolOption(p, "noepoll")
	if err != nil {
		return nil, err
	}
	if noepoll {
		isEmpty = false
		options.Noepoll = noepoll
	}

	nokqueue, err := parseBoolOption(p, "nokqueue")
	if err != nil {
		return nil, err
	}
	if nokqueue {
		isEmpty = false
		options.Nokqueue = nokqueue
	}

	noevports, err := parseBoolOption(p, "noevports")
	if err != nil {
		return nil, err
	}
	if noevports {
		isEmpty = false
		options.Noevports = noevports
	}

	nopoll, err := parseBoolOption(p, "nopoll")
	if err != nil {
		return nil, err
	}
	if nopoll {
		isEmpty = false
		options.Nopoll = nopoll
	}

	nosplice, err := parseBoolOption(p, "nosplice")
	if err != nil {
		return nil, err
	}
	if nosplice {
		isEmpty = false
		options.Nosplice = nosplice
	}

	nogetaddrinfo, err := parseBoolOption(p, "nogetaddrinfo")
	if err != nil {
		return nil, err
	}
	if nogetaddrinfo {
		isEmpty = false
		options.Nogetaddrinfo = nogetaddrinfo
	}

	noreuseport, err := parseBoolOption(p, "noreuseport")
	if err != nil {
		return nil, err
	}
	if noreuseport {
		isEmpty = false
		options.Noreuseport = noreuseport
	}

	profilingTasks, err := parseAutoOnOffOption(p, "profiling.tasks")
	if err != nil {
		return nil, err
	}
	if profilingTasks != "" {
		isEmpty = false
		options.ProfilingTasks = profilingTasks
	}

	profilingMemory, err := parseOnOffOption(p, "profiling.memory")
	if err != nil {
		return nil, err
	}
	if profilingMemory != "" {
		isEmpty = false
		options.ProfilingMemory = profilingMemory
	}

	srvStateBase, err := parseStringOption(p, "server-state-base")
	if err != nil {
		return nil, err
	}
	if srvStateBase != "" {
		isEmpty = false
		options.ServerStateBase = srvStateBase
	}

	srvStateFile, err := parseStringOption(p, "server-state-file")
	if err != nil {
		return nil, err
	}
	if srvStateFile != "" {
		isEmpty = false
		options.ServerStateFile = srvStateFile
	}

	spreadChecks, err := parseInt64Option(p, "spread-checks")
	if err != nil {
		return nil, err
	}
	if spreadChecks != 0 {
		isEmpty = false
		options.SpreadChecks = spreadChecks
	}

	threadHardLimit, err := parseInt64POption(p, "thread-hard-limit")
	if err != nil {
		return nil, err
	}
	if threadHardLimit != nil {
		isEmpty = false
		options.ThreadHardLimit = threadHardLimit
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}

	return options, nil
}

func parseSSLOptions(p parser.Parser) (*models.SslOptions, error) { //nolint:gocognit,gocyclo,cyclop,maintidx
	options := &models.SslOptions{}
	isEmpty := true
	var sslEngines []*models.SslEngine
	data, err := p.Get(parser.Global, parser.GlobalSectionName, "ssl-engine")
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
	if sslEngines != nil {
		isEmpty = false
		options.SslEngines = sslEngines
	}

	sched, err := parseStringOption(p, "acme.scheduler")
	if err != nil {
		return nil, err
	}
	if sched != "" {
		isEmpty = false
		options.AcmeScheduler = sched
	}

	caBase, err := parseStringOption(p, "ca-base")
	if err != nil {
		return nil, err
	}
	if caBase != "" {
		isEmpty = false
		options.CaBase = caBase
	}

	crtBase, err := parseStringOption(p, "crt-base")
	if err != nil {
		return nil, err
	}
	if crtBase != "" {
		isEmpty = false
		options.CrtBase = crtBase
	}

	sslBindCiphers, err := parseStringOption(p, "ssl-default-bind-ciphers")
	if err != nil {
		return nil, err
	}
	if sslBindCiphers != "" {
		isEmpty = false
		options.DefaultBindCiphers = sslBindCiphers
	}

	sslBindCiphersuites, err := parseStringOption(p, "ssl-default-bind-ciphersuites")
	if err != nil {
		return nil, err
	}
	if sslBindCiphersuites != "" {
		isEmpty = false
		options.DefaultBindCiphersuites = sslBindCiphersuites
	}

	sslDefaultBindCurves, err := parseStringOption(p, "ssl-default-bind-curves")
	if err != nil {
		return nil, err
	}
	if sslDefaultBindCurves != "" {
		isEmpty = false
		options.DefaultBindCurves = sslDefaultBindCurves
	}

	sslDefaultServerCurves, err := parseStringOption(p, "ssl-default-server-curves")
	if err != nil {
		return nil, err
	}
	if sslDefaultServerCurves != "" {
		isEmpty = false
		options.DefaultServerCurves = sslDefaultServerCurves
	}

	sslSkipSelfIssuedCa, err := parseBoolOption(p, "ssl-skip-self-issued-ca")
	if err != nil {
		return nil, err
	}
	if sslSkipSelfIssuedCa {
		isEmpty = false
		options.SkipSelfIssuedCa = sslSkipSelfIssuedCa
	}

	sslBindOptions, err := parseStringOption(p, "ssl-default-bind-options")
	if err != nil {
		return nil, err
	}
	if sslBindOptions != "" {
		isEmpty = false
		options.DefaultBindOptions = sslBindOptions
	}

	sslBindSigalgs, err := parseStringOption(p, "ssl-default-bind-sigalgs")
	if err != nil {
		return nil, err
	}
	if sslBindSigalgs != "" {
		isEmpty = false
		options.DefaultBindSigalgs = sslBindSigalgs
	}

	sslBindClientSigalgs, err := parseStringOption(p, "ssl-default-bind-client-sigalgs")
	if err != nil {
		return nil, err
	}
	if sslBindClientSigalgs != "" {
		isEmpty = false
		options.DefaultBindClientSigalgs = sslBindClientSigalgs
	}

	sslServerSigalgs, err := parseStringOption(p, "ssl-default-server-sigalgs")
	if err != nil {
		return nil, err
	}
	if sslServerSigalgs != "" {
		isEmpty = false
		options.DefaultServerSigalgs = sslServerSigalgs
	}

	sslServerClientSigalgs, err := parseStringOption(p, "ssl-default-server-client-sigalgs")
	if err != nil {
		return nil, err
	}
	if sslServerClientSigalgs != "" {
		isEmpty = false
		options.DefaultServerClientSigalgs = sslServerClientSigalgs
	}

	sslDefaultServerCiphers, err := parseStringOption(p, "ssl-default-server-ciphers")
	if err != nil {
		return nil, err
	}
	if sslDefaultServerCiphers != "" {
		isEmpty = false
		options.DefaultServerCiphers = sslDefaultServerCiphers
	}

	sslServerCiphersuites, err := parseStringOption(p, "ssl-default-server-ciphersuites")
	if err != nil {
		return nil, err
	}
	if sslServerCiphersuites != "" {
		isEmpty = false
		options.DefaultServerCiphersuites = sslServerCiphersuites
	}

	sslServerOptions, err := parseStringOption(p, "ssl-default-server-options")
	if err != nil {
		return nil, err
	}
	if sslServerOptions != "" {
		isEmpty = false
		options.DefaultServerOptions = sslServerOptions
	}

	data, _ = p.Get(parser.Global, parser.GlobalSectionName, "ssl-mode-async")
	if _, ok := data.(*types.SslModeAsync); ok {
		isEmpty = false
		options.ModeAsync = "enabled"
	}

	sslDhParamFile, err := parseStringOption(p, "ssl-dh-param-file")
	if err != nil {
		return nil, err
	}
	if sslDhParamFile != "" {
		isEmpty = false
		options.DhParamFile = sslDhParamFile
	}

	sslPropquery, err := parseStringOption(p, "ssl-propquery")
	if err != nil {
		return nil, err
	}
	if sslPropquery != "" {
		isEmpty = false
		options.Propquery = sslPropquery
	}

	sslProvider, err := parseStringOption(p, "ssl-provider")
	if err != nil {
		return nil, err
	}
	if sslProvider != "" {
		isEmpty = false
		options.Provider = sslProvider
	}

	sslProviderPath, err := parseStringOption(p, "ssl-provider-path")
	if err != nil {
		return nil, err
	}
	if sslProviderPath != "" {
		isEmpty = false
		options.ProviderPath = sslProviderPath
	}

	sslServerVerify, err := parseStringOption(p, "ssl-server-verify")
	if err != nil {
		return nil, err
	}
	if sslServerVerify != "" {
		isEmpty = false
		options.ServerVerify = sslServerVerify
	}

	issuersChainPath, err := parseStringOption(p, "issuers-chain-path")
	if err != nil {
		return nil, err
	}
	if issuersChainPath != "" {
		isEmpty = false
		options.IssuersChainPath = issuersChainPath
	}

	sslLoadExtraFiles, err := parseStringOption(p, "ssl-load-extra-files")
	if err != nil {
		return nil, err
	}
	if sslLoadExtraFiles != "" {
		isEmpty = false
		options.LoadExtraFiles = sslLoadExtraFiles
	}

	sslPassphraseCmd, err := parseStringOption(p, "ssl-passphrase-cmd")
	if err != nil {
		return nil, err
	}
	if sslPassphraseCmd != "" {
		isEmpty = false
		options.PassphraseCmd = sslPassphraseCmd
	}

	maxsslconn, err := parseInt64Option(p, "maxsslconn")
	if err != nil {
		return nil, err
	}
	if maxsslconn != 0 {
		isEmpty = false
		options.Maxsslconn = maxsslconn
	}

	maxsslrate, err := parseInt64Option(p, "maxsslrate")
	if err != nil {
		return nil, err
	}
	if maxsslrate != 0 {
		isEmpty = false
		options.Maxsslrate = maxsslrate
	}

	sslSecurityLevel, err := parseInt64POption(p, "ssl-security-level")
	if err != nil {
		return nil, err
	}
	if sslSecurityLevel != nil {
		isEmpty = false
		options.SecurityLevel = sslSecurityLevel
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func parseDebugOptions(p parser.Parser) (*models.DebugOptions, error) {
	options := &models.DebugOptions{}
	isEmpty := true

	anonkey, err := parseInt64POption(p, "anonkey")
	if err != nil {
		return nil, err
	}
	if anonkey != nil {
		isEmpty = false
		options.Anonkey = anonkey
	}

	stress, err := parseInt64POption(p, "stress-level")
	if err != nil {
		return nil, err
	}
	if stress != nil {
		isEmpty = false
		options.StressLevel = stress
	}

	quiet, err := parseBoolOption(p, "quiet")
	if err != nil {
		return nil, err
	}
	if quiet {
		isEmpty = false
		options.Quiet = quiet
	}

	zeroWarning, err := parseBoolOption(p, "zero-warning")
	if err != nil {
		return nil, err
	}
	if zeroWarning {
		isEmpty = false
		options.ZeroWarning = zeroWarning
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func parseEnvironmentOptions(p parser.Parser) (*models.EnvironmentOptions, error) {
	options := &models.EnvironmentOptions{}
	isEmpty := true
	var presetEnvs []*models.PresetEnv
	data, err := p.Get(parser.Global, parser.GlobalSectionName, "presetenv")
	if err == nil {
		envs, ok := data.([]types.StringKeyValueC)
		if !ok {
			return nil, misc.CreateTypeAssertError("presetenv")
		}
		for _, e := range envs {
			env := &models.PresetEnv{
				Name:  misc.Ptr(e.Key),
				Value: misc.Ptr(e.Value),
			}
			presetEnvs = append(presetEnvs, env)
		}
	}
	if presetEnvs != nil {
		isEmpty = false
		options.PresetEnvs = presetEnvs
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
	if setEnvs != nil {
		isEmpty = false
		options.SetEnvs = setEnvs
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
	if resetEnv != "" {
		isEmpty = false
		options.Resetenv = resetEnv
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
	if unsetEnv != "" {
		isEmpty = false
		options.Unsetenv = unsetEnv
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func parseDeviceAtlasOptions(p parser.Parser) (*models.DeviceAtlasOptions, error) {
	options := &models.DeviceAtlasOptions{}
	isEmpty := true
	var option string
	var err error

	option, err = parseStringOption(p, "deviceatlas-json-file")
	if err != nil {
		return nil, err
	}
	if option != "" {
		isEmpty = false
		options.JSONFile = option
	}

	option, err = parseStringOption(p, "deviceatlas-log-level")
	if err != nil {
		return nil, err
	}
	if option != "" {
		isEmpty = false
		options.LogLevel = option
	}

	option, err = parseStringOption(p, "deviceatlas-separator")
	if err != nil {
		return nil, err
	}
	if option != "" {
		isEmpty = false
		options.Separator = option
	}

	option, err = parseStringOption(p, "deviceatlas-properties-cookie")
	if err != nil {
		return nil, err
	}
	if option != "" {
		isEmpty = false
		options.PropertiesCookie = option
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func parseFiftyOneDegreesOptions(p parser.Parser) (*models.FiftyOneDegreesOptions, error) {
	options := &models.FiftyOneDegreesOptions{}
	isEmpty := true
	var option string
	var optionInt int64
	var err error

	option, err = parseStringOption(p, "51degrees-data-file")
	if err != nil {
		return nil, err
	}
	if option != "" {
		isEmpty = false
		options.DataFile = option
	}

	option, err = parseStringOption(p, "51degrees-property-name-list")
	if err != nil {
		return nil, err
	}
	if option != "" {
		isEmpty = false
		options.PropertyNameList = option
	}

	option, err = parseStringOption(p, "51degrees-property-separator")
	if err != nil {
		return nil, err
	}
	if option != "" {
		isEmpty = false
		options.PropertySeparator = option
	}

	optionInt, err = parseInt64Option(p, "51degrees-cache-size")
	if err != nil {
		return nil, err
	}
	if optionInt != 0 {
		isEmpty = false
		options.CacheSize = optionInt
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func parseHardenOptions(p parser.Parser) (*models.GlobalHarden, error) {
	options := &models.GlobalHarden{}
	hardenRejectPrivilgedPortQuic, err := parseOnOffOption(p, "harden.reject-privileged-ports.quic")
	if err != nil {
		return nil, err
	}
	hardenRejectPrivilgedPortTCP, err := parseOnOffOption(p, "harden.reject-privileged-ports.tcp")
	if err != nil {
		return nil, err
	}
	if hardenRejectPrivilgedPortQuic != "" || hardenRejectPrivilgedPortTCP != "" {
		options.RejectPrivilegedPorts = &models.GlobalHardenRejectPrivilegedPorts{}
		options.RejectPrivilegedPorts.Quic = hardenRejectPrivilgedPortQuic
		options.RejectPrivilegedPorts.TCP = hardenRejectPrivilgedPortTCP
	} else {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func parseHTTPClientOptions(p parser.Parser) (*models.HTTPClientOptions, error) {
	options := &models.HTTPClientOptions{}
	isEmpty := true

	httpClientResolversDisabled, err := parseOnOffOption(p, "httpclient.resolvers.disabled")
	if err != nil {
		return nil, err
	}
	if httpClientResolversDisabled != "" {
		isEmpty = false
		options.ResolversDisabled = httpClientResolversDisabled
	}

	httpClientResolversID, err := parseStringOption(p, "httpclient.resolvers.id")
	if err != nil {
		return nil, err
	}
	if httpClientResolversID != "" {
		isEmpty = false
		options.ResolversID = httpClientResolversID
	}

	data, err := p.Get(parser.Global, parser.GlobalSectionName, "httpclient.resolvers.prefer")
	if err == nil {
		resolverPreferParser, ok := data.(*types.HTTPClientResolversPrefer)
		if !ok {
			return nil, misc.CreateTypeAssertError("httpclient.resolvers.prefer")
		}
		if resolverPreferParser != nil {
			isEmpty = false
			options.ResolversPrefer = resolverPreferParser.Type
		}
	}

	httpClientRetries, err := parseInt64Option(p, "httpclient.retries")
	if err != nil {
		return nil, err
	}
	if httpClientRetries != 0 {
		isEmpty = false
		options.Retries = httpClientRetries
	}

	httpClientSSLCaFile, err := parseStringOption(p, "httpclient.ssl.ca-file")
	if err != nil {
		return nil, err
	}
	if httpClientSSLCaFile != "" {
		isEmpty = false
		options.SslCaFile = httpClientSSLCaFile
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
	if httpClientSSLVerify != nil {
		isEmpty = false
		options.SslVerify = httpClientSSLVerify
	}

	httpClientTimeoutConnect, err := parseTimeoutOption(p, "httpclient.timeout.connect")
	if err != nil {
		return nil, err
	}
	if httpClientTimeoutConnect != nil {
		isEmpty = false
		options.TimeoutConnect = httpClientTimeoutConnect
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func parseWurflOptions(p parser.Parser) (*models.WurflOptions, error) {
	options := &models.WurflOptions{}
	isEmpty := true
	wurflDataFile, err := parseStringOption(p, "wurfl-data-file")
	if err != nil {
		return nil, err
	}
	if wurflDataFile != "" {
		isEmpty = false
		options.DataFile = wurflDataFile
	}

	wurflInformationList, err := parseStringOption(p, "wurfl-information-list")
	if err != nil {
		return nil, err
	}
	if wurflInformationList != "" {
		isEmpty = false
		options.InformationList = wurflInformationList
	}

	wurflInformationListSeparator, err := parseStringOption(p, "wurfl-information-list-separator")
	if err != nil {
		return nil, err
	}
	if wurflInformationListSeparator != "" {
		isEmpty = false
		options.InformationListSeparator = wurflInformationListSeparator
	}

	wurflPatchFile, err := parseStringOption(p, "wurfl-patch-file")
	if err != nil {
		return nil, err
	}
	if wurflPatchFile != "" {
		isEmpty = false
		options.PatchFile = wurflPatchFile
	}

	wurflCacheSize, err := parseInt64Option(p, "wurfl-cache-size")
	if err != nil {
		return nil, err
	}
	if wurflCacheSize != 0 {
		isEmpty = false
		options.CacheSize = wurflCacheSize
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func parseTuneOptions(p parser.Parser) (*models.TuneOptions, error) { //nolint:gocognit, gocyclo, cyclop,maintidx
	options := &models.TuneOptions{}
	isEmpty := true
	var intOption int64
	var intPOption *int64
	var boolOption bool
	var strOption string
	var err error

	strOption, err = parseOnOffOption(p, "tune.applet.zero-copy-forwarding")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		isEmpty = false
		options.AppletZeroCopyForwarding = strOption
	}

	intOption, err = parseInt64Option(p, "tune.comp.maxlevel")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.CompMaxlevel = intOption
	}

	boolOption, err = parseBoolOption(p, "tune.disable-fast-forward")
	if err != nil {
		return nil, err
	}
	if boolOption {
		isEmpty = false
		options.DisableFastForward = boolOption
	}

	boolOption, err = parseBoolOption(p, "tune.disable-zero-copy-forwarding")
	if err != nil {
		return nil, err
	}
	if boolOption {
		isEmpty = false
		options.DisableZeroCopyForwarding = boolOption
	}

	strOption, err = parseStringOption(p, "tune.epoll.mask-events")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		if words := strings.Split(strOption, ","); len(words) > 0 && words[0] != "" {
			isEmpty = false
			options.EpollMaskEvents = words
		}
	}

	intOption, err = parseInt64Option(p, "tune.events.max-events-at-once")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.EventsMaxEventsAtOnce = intOption
	}

	boolOption, err = parseBoolOption(p, "tune.fail-alloc")
	if err != nil {
		return nil, err
	}
	if boolOption {
		isEmpty = false
		options.FailAlloc = boolOption
	}

	intPOption, err = parseInt64POption(p, "tune.glitches.kill.cpu-usage")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.GlitchesKillCPUUsage = intPOption
	}

	intOption, err = parseInt64Option(p, "tune.h2.header-table-size")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.H2HeaderTableSize = intOption
	}

	intPOption, err = parseInt64POption(p, "tune.h2.initial-window-size")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.H2InitialWindowSize = intPOption
	}

	intOption, err = parseInt64Option(p, "tune.h2.max-concurrent-streams")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.H2MaxConcurrentStreams = intOption
	}

	intOption, err = parseInt64Option(p, "tune.h2.max-frame-size")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.H2MaxFrameSize = intOption
	}

	intOption, err = parseInt64Option(p, "tune.http.cookielen")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.HTTPCookielen = intOption
	}

	intOption, err = parseInt64Option(p, "tune.http.logurilen")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.HTTPLogurilen = intOption
	}

	intOption, err = parseInt64Option(p, "tune.http.maxhdr")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.HTTPMaxhdr = intOption
	}

	strOption, err = parseOnOffOption(p, "tune.idle-pool.shared")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		isEmpty = false
		options.IdlePoolShared = strOption
	}

	intPOption, err = parseTimeoutOption(p, "tune.idletimer")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.Idletimer = intPOption
	}

	strOption, err = parseListenerDefaultShards(p, "tune.listener.default-shards")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		isEmpty = false
		options.ListenerDefaultShards = strOption
	}

	strOption, err = parseOnOffOption(p, "tune.listener.multi-queue")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		isEmpty = false
		options.ListenerMultiQueue = strOption
	}

	intPOption, err = parseInt64POption(p, "tune.max-checks-per-thread")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.MaxChecksPerThread = intPOption
	}

	intPOption, err = parseInt64POption(p, "tune.max-rules-at-once")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.MaxRulesAtOnce = intPOption
	}

	intOption, err = parseInt64Option(p, "tune.maxaccept")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.Maxaccept = intOption
	}

	intOption, err = parseInt64Option(p, "tune.maxpollevents")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.Maxpollevents = intOption
	}

	intOption, err = parseInt64Option(p, "tune.maxrewrite")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.Maxrewrite = intOption
	}

	intPOption, err = parseInt64POption(p, "tune.memory.hot-size")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.MemoryHotSize = intPOption
	}

	intPOption, err = parseInt64POption(p, "tune.pattern.cache-size")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.PatternCacheSize = intPOption
	}

	intOption, err = parseInt64Option(p, "tune.peers.max-updates-at-once")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.PeersMaxUpdatesAtOnce = intOption
	}

	intOption, err = parseInt64Option(p, "tune.pool-high-fd-ratio")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.PoolHighFdRatio = intOption
	}

	intOption, err = parseInt64Option(p, "tune.pool-low-fd-ratio")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.PoolLowFdRatio = intOption
	}

	intPOption, err = parseInt64POption(p, "tune.ring.queues")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.RingQueues = intPOption
	}

	intOption, err = parseInt64Option(p, "tune.runqueue-depth")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.RunqueueDepth = intOption
	}

	strOption, err = parseOnOffOption(p, "tune.sched.low-latency")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		isEmpty = false
		options.SchedLowLatency = strOption
	}

	intPOption, err = parseInt64POption(p, "tune.stick-counters")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.StickCounters = intPOption
	}

	strOption, err = parseStringOption(p, "tune.takeover-other-tg-connections")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		isEmpty = false
		options.TakeoverOtherTgConnections = strOption
	}

	strOption, err = parseOnOffOption(p, "tune.fd.edge-triggered")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		isEmpty = false
		options.FdEdgeTriggered = strOption
	}

	strOption, err = parseOnOffOption(p, "tune.h1.zero-copy-fwd-recv")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		isEmpty = false
		options.H1ZeroCopyFwdRecv = strOption
	}

	strOption, err = parseOnOffOption(p, "tune.h1.zero-copy-fwd-send")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		isEmpty = false
		options.H1ZeroCopyFwdSend = strOption
	}

	intPOption, err = parseInt64POption(p, "tune.h2.be.glitches-threshold")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.H2BeGlitchesThreshold = intPOption
	}

	intOption, err = parseInt64Option(p, "tune.h2.be.initial-window-size")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.H2BeInitialWindowSize = intOption
	}

	intOption, err = parseInt64Option(p, "tune.h2.be.max-concurrent-streams")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.H2BeMaxConcurrentStreams = intOption
	}

	intPOption, err = parseSizeOption(p, "tune.h2.be.rxbuf")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.H2BeRxbuf = intPOption
	}

	intPOption, err = parseInt64POption(p, "tune.h2.fe.glitches-threshold")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.H2FeGlitchesThreshold = intPOption
	}

	intOption, err = parseInt64Option(p, "tune.h2.fe.initial-window-size")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.H2FeInitialWindowSize = intOption
	}

	intOption, err = parseInt64Option(p, "tune.h2.fe.max-concurrent-streams")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.H2FeMaxConcurrentStreams = intOption
	}

	intPOption, err = parseInt64POption(p, "tune.h2.fe.max-total-streams")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.H2FeMaxTotalStreams = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.h2.fe.rxbuf")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.H2FeRxbuf = intPOption
	}

	strOption, err = parseOnOffOption(p, "tune.h2.zero-copy-fwd-send")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		isEmpty = false
		options.H2ZeroCopyFwdSend = strOption
	}

	intPOption, err = parseSizeOption(p, "tune.notsent-lowat.client")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.NotsentLowatClient = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.notsent-lowat.server")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.NotsentLowatServer = intPOption
	}

	strOption, err = parseOnOffOption(p, "tune.pt.zero-copy-forwarding")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		isEmpty = false
		options.PtZeroCopyForwarding = strOption
	}

	intPOption, err = parseInt64POption(p, "tune.renice.runtime")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.ReniceRuntime = intPOption
	}

	intPOption, err = parseInt64POption(p, "tune.renice.startup")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.ReniceStartup = intPOption
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func parseTuneBufferOptions(p parser.Parser) (*models.TuneBufferOptions, error) { //nolint:gocognit
	options := &models.TuneBufferOptions{}
	isEmpty := true

	intPOption, err := parseInt64POption(p, "tune.buffers.limit")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.BuffersLimit = intPOption
	}

	intOption, err := parseInt64Option(p, "tune.buffers.reserve")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.BuffersReserve = intOption
	}

	intPOption, err = parseSizeOption(p, "tune.bufsize")
	if err != nil {
		return nil, err
	}
	if intPOption != nil && *intPOption != 0 {
		isEmpty = false
		options.Bufsize = *intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.bufsize.small")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.BufsizeSmall = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.pipesize")
	if err != nil {
		return nil, err
	}
	if intPOption != nil && *intPOption != 0 {
		isEmpty = false
		options.Pipesize = *intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.rcvbuf.backend")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.RcvbufBackend = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.rcvbuf.client")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.RcvbufClient = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.rcvbuf.frontend")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.RcvbufFrontend = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.rcvbuf.server")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.RcvbufServer = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.recv_enough")
	if err != nil {
		return nil, err
	}
	if intPOption != nil && *intPOption != 0 {
		isEmpty = false
		options.RecvEnough = *intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.sndbuf.backend")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.SndbufBackend = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.sndbuf.client")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.SndbufClient = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.sndbuf.frontend")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.SndbufFrontend = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.sndbuf.server")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.SndbufServer = intPOption
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func parseTuneLuaOptions(p parser.Parser) (*models.TuneLuaOptions, error) {
	var intOption int64
	var intPOption *int64
	var strOption string
	var err error
	options := &models.TuneLuaOptions{}
	isEmpty := true

	strOption, err = parseStringOption(p, "tune.lua.bool-sample-conversion")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		isEmpty = false
		options.BoolSampleConversion = strOption
	}

	intOption, err = parseInt64Option(p, "tune.lua.forced-yield")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		options.ForcedYield = intOption
		isEmpty = false
	}

	intPOption, err = parseInt64POption(p, "tune.lua.maxmem")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		options.Maxmem = intPOption
		isEmpty = false
	}

	strOption, err = parseOnOffOption(p, "tune.lua.log.loggers")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		options.LogLoggers = strOption
		isEmpty = false
	}

	strOption, err = parseAutoOnOffOption(p, "tune.lua.log.stderr")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		options.LogStderr = strOption
		isEmpty = false
	}

	intPOption, err = parseTimeoutOption(p, "tune.lua.session-timeout")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		options.SessionTimeout = intPOption
		isEmpty = false
	}

	intPOption, err = parseTimeoutOption(p, "tune.lua.burst-timeout")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		options.BurstTimeout = intPOption
		isEmpty = false
	}

	intPOption, err = parseTimeoutOption(p, "tune.lua.task-timeout")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		options.TaskTimeout = intPOption
		isEmpty = false
	}

	intPOption, err = parseTimeoutOption(p, "tune.lua.service-timeout")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		options.ServiceTimeout = intPOption
		isEmpty = false
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}

	return options, nil
}

func parseTuneQuicOptions(p parser.Parser) (*models.TuneQuicOptions, error) {
	options := &models.TuneQuicOptions{}
	isEmpty := true
	intPOption, err := parseInt64POption(p, "tune.quic.frontend.conn-tx-buffers.limit")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.FrontendConnTxBuffersLimit = intPOption
	}

	intPOption, err = parseTimeoutOption(p, "tune.quic.frontend.max-idle-timeout")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.FrontendMaxIdleTimeout = intPOption
	}

	intPOption, err = parseInt64POption(p, "tune.quic.frontend.max-streams-bidi")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.FrontendMaxStreamsBidi = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.quic.frontend.max-tx-mem")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.FrontendMaxTxMemory = intPOption
	}

	intPOption, err = parseInt64POption(p, "tune.quic.max-frame-loss")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.MaxFrameLoss = intPOption
	}

	intPOption, err = parseInt64POption(p, "tune.quic.reorder-ratio")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.ReorderRatio = intPOption
	}

	intPOption, err = parseInt64POption(p, "tune.quic.retry-threshold")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.RetryThreshold = intPOption
	}

	so, err := p.Get(parser.Global, parser.GlobalSectionName, "tune.quic.socket-owner")
	if err == nil {
		value, ok := so.(*types.QuicSocketOwner)
		if !ok {
			return nil, misc.CreateTypeAssertError("tune.quic.socket-owner")
		}
		if value != nil && value.Owner != "" {
			isEmpty = false
			options.SocketOwner = value.Owner
		}
	}

	strOption, err := parseOnOffOption(p, "tune.quic.zero-copy-fwd-send")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		isEmpty = false
		options.ZeroCopyFwdSend = strOption
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func parseTuneSSLOptions(p parser.Parser) (*models.TuneSslOptions, error) {
	options := &models.TuneSslOptions{}
	isEmpty := true
	intPOption, err := parseInt64POption(p, "tune.ssl.cachesize")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.Cachesize = intPOption
	}

	boolOption, err := parseBoolOption(p, "tune.ssl.force-private-cache")
	if err != nil {
		return nil, err
	}
	if boolOption {
		isEmpty = false
		options.ForcePrivateCache = boolOption
	}

	strOption, err := parseOnOffOption(p, "tune.ssl.keylog")
	if err != nil {
		return nil, err
	}
	if strOption != "" {
		isEmpty = false
		options.Keylog = strOption
	}

	intPOption, err = parseTimeoutOption(p, "tune.ssl.lifetime")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.Lifetime = intPOption
	}

	intPOption, err = parseInt64POption(p, "tune.ssl.maxrecord")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.Maxrecord = intPOption
	}

	intOption, err := parseInt64Option(p, "tune.ssl.default-dh-param")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.DefaultDhParam = intOption
	}

	intOption, err = parseInt64Option(p, "tune.ssl.ssl-ctx-cache-size")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.CtxCacheSize = intOption
	}

	intPOption, err = parseInt64POption(p, "tune.ssl.capture-buffer-size")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.CaptureBufferSize = intPOption
	}

	intPOption, err = parseInt64POption(p, "tune.ssl.ocsp-update.maxdelay")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.OcspUpdateMaxDelay = intPOption
	}

	intPOption, err = parseInt64POption(p, "tune.ssl.ocsp-update.mindelay")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.OcspUpdateMinDelay = intPOption
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func parseTuneVarsOptions(p parser.Parser) (*models.TuneVarsOptions, error) {
	options := &models.TuneVarsOptions{}
	isEmpty := true
	intPOption, err := parseSizeOption(p, "tune.vars.global-max-size")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.GlobalMaxSize = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.vars.proc-max-size")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.ProcMaxSize = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.vars.reqres-max-size")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.ReqresMaxSize = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.vars.sess-max-size")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.SessMaxSize = intPOption
	}

	intPOption, err = parseSizeOption(p, "tune.vars.txn-max-size")
	if err != nil {
		return nil, err
	}
	if intPOption != nil {
		isEmpty = false
		options.TxnMaxSize = intPOption
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func parseTuneZlibOptions(p parser.Parser) (*models.TuneZlibOptions, error) {
	options := &models.TuneZlibOptions{}
	isEmpty := true
	intOption, err := parseInt64Option(p, "tune.zlib.memlevel")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.Memlevel = intOption
	}

	intOption, err = parseInt64Option(p, "tune.zlib.windowsize")
	if err != nil {
		return nil, err
	}
	if intOption != 0 {
		isEmpty = false
		options.Windowsize = intOption
	}

	if isEmpty {
		return nil, nil //nolint:nilnil
	}
	return options, nil
}

func ParseGlobalSection(p parser.Parser) (*models.Global, error) { //nolint:gocognit,gocyclo,cyclop,maintidx
	global := &models.Global{}
	cpuMaps, err := parseCPUMaps(p)
	if err != nil {
		return nil, err
	}
	global.CPUMaps = cpuMaps

	cpuSets, err := parseCPUSets(p)
	if err != nil {
		return nil, err
	}
	global.CPUSets = cpuSets

	h1CaseAdjusts, err := parseH1CaseAdjusts(p)
	if err != nil {
		return nil, err
	}
	global.H1CaseAdjusts = h1CaseAdjusts

	runtimeAPIs, err := parseRuntimeAPIs(p)
	if err != nil {
		return nil, err
	}
	global.RuntimeAPIs = runtimeAPIs

	setVars, err := parseSetVars(p)
	if err != nil {
		return nil, err
	}
	global.SetVars = setVars

	setVarFmts, err := parseSetVarFmts(p)
	if err != nil {
		return nil, err
	}
	global.SetVarFmts = setVarFmts

	threadGroupLines, err := parseThreadGroupLines(p)
	if err != nil {
		return nil, err
	}
	global.ThreadGroupLines = threadGroupLines

	chroot, err := parseStringOption(p, "chroot")
	if err != nil {
		return nil, err
	}
	global.Chroot = chroot

	closeSpreadTime, err := parseTimeoutOption(p, "close-spread-time")
	if err != nil {
		return nil, err
	}
	global.CloseSpreadTime = closeSpreadTime

	clusterSecret, err := parseStringOption(p, "cluster-secret")
	if err != nil {
		return nil, err
	}
	global.ClusterSecret = clusterSecret

	cpuPolicy, err := parseStringOption(p, "cpu-policy")
	if err != nil {
		return nil, err
	}
	global.CPUPolicy = cpuPolicy

	daemon, err := parseBoolOption(p, "daemon")
	if err != nil {
		return nil, err
	}
	global.Daemon = daemon

	debugOptions, err := parseDebugOptions(p)
	if err != nil {
		return nil, err
	}
	global.DebugOptions = debugOptions

	defaultPath, err := parseDefaultPath(p)
	if err != nil {
		return nil, err
	}
	global.DefaultPath = defaultPath

	description, err := parseStringOption(p, "description")
	if err != nil {
		return nil, err
	}
	global.Description = description

	deviceAtlasOptions, err := parseDeviceAtlasOptions(p)
	if err != nil {
		return nil, err
	}
	global.DeviceAtlasOptions = deviceAtlasOptions

	dnsAcceptFamily, err := parseStringOption(p, "dns-accept-family")
	if err != nil {
		return nil, err
	}
	global.DNSAcceptFamily = dnsAcceptFamily

	shmStatsFile, err := parseStringOption(p, "shm-stats-file")
	if err != nil {
		return nil, err
	}
	global.ShmStatsFile = shmStatsFile

	shmStatSileMaxObjects, err := parseInt64POption(p, "shm-stats-file-max-objects")
	if err != nil {
		return nil, err
	}
	global.ShmStatsFileMaxObjects = shmStatSileMaxObjects

	envOptions, err := parseEnvironmentOptions(p)
	if err != nil {
		return nil, err
	}
	global.EnvironmentOptions = envOptions

	exposeDeprecatedDirectives, err := parseBoolOption(p, "expose-deprecated-directives")
	if err != nil {
		return nil, err
	}
	global.ExposeDeprecatedDirectives = exposeDeprecatedDirectives

	exposeExperimentalDirectives, err := parseBoolOption(p, "expose-experimental-directives")
	if err != nil {
		return nil, err
	}
	global.ExposeExperimentalDirectives = exposeExperimentalDirectives

	externalCheck, err := parseBoolOption(p, "external-check")
	if err != nil {
		return nil, err
	}
	global.ExternalCheck = externalCheck

	fiftyOneDegreesOptions, err := parseFiftyOneDegreesOptions(p)
	if err != nil {
		return nil, err
	}
	global.FiftyOneDegreesOptions = fiftyOneDegreesOptions

	global.ForceCfgParserPause, err = parseTimeoutOption(p, "force-cfg-parser-pause")
	if err != nil {
		return nil, err
	}

	gid, err := parseInt64Option(p, "gid")
	if err != nil {
		return nil, err
	}
	global.Gid = gid

	grace, err := parseTimeoutOption(p, "grace")
	if err != nil {
		return nil, err
	}
	global.Grace = grace

	group, err := parseStringOption(p, "group")
	if err != nil {
		return nil, err
	}
	global.Group = group

	h1CaseAdjustFile, err := parseStringOption(p, "h1-case-adjust-file")
	if err != nil {
		return nil, err
	}
	global.H1CaseAdjustFile = h1CaseAdjustFile

	h1AcceptPayloadWithAnyMethod, err := parseBoolOption(p, "h1-accept-payload-with-any-method")
	if err != nil {
		return nil, err
	}
	global.H1AcceptPayloadWithAnyMethod = h1AcceptPayloadWithAnyMethod

	h1DoNotCloseOnInsecureTransferEncoding, err := parseBoolOption(p, "h1-do-not-close-on-insecure-transfer-encoding")
	if err != nil {
		return nil, err
	}
	global.H1DoNotCloseOnInsecureTransferEncoding = h1DoNotCloseOnInsecureTransferEncoding

	h2WorkaroundBogusWebsocketClients, err := parseBoolOption(p, "h2-workaround-bogus-websocket-clients")
	if err != nil {
		return nil, err
	}
	global.H2WorkaroundBogusWebsocketClients = h2WorkaroundBogusWebsocketClients

	hardStop, err := parseTimeoutOption(p, "hard-stop-after")
	if err != nil {
		return nil, err
	}
	global.HardStopAfter = hardStop

	harden, err := parseHardenOptions(p)
	if err != nil {
		return nil, err
	}
	global.Harden = harden

	httpClientOptions, err := parseHTTPClientOptions(p)
	if err != nil {
		return nil, err
	}
	global.HTTPClientOptions = httpClientOptions

	var errCodes []*models.HTTPCodes
	data, err := p.Get(parser.Global, parser.GlobalSectionName, "http-err-codes")
	if err == nil {
		errCodesParser, ok := data.([]types.HTTPErrCodes)
		if !ok {
			return nil, misc.CreateTypeAssertError("http-err-codes")
		}
		for _, e := range errCodesParser {
			errCode := &models.HTTPCodes{
				Value: misc.Ptr(e.Value),
			}
			errCodes = append(errCodes, errCode)
		}
	}
	global.HTTPErrCodes = errCodes

	var failCodes []*models.HTTPCodes
	data, err = p.Get(parser.Global, parser.GlobalSectionName, "http-fail-codes")
	if err == nil {
		failCodesParser, ok := data.([]types.HTTPFailCodes)
		if !ok {
			return nil, misc.CreateTypeAssertError("http-fail-codes")
		}
		for _, f := range failCodesParser {
			failCode := &models.HTTPCodes{
				Value: misc.Ptr(f.Value),
			}
			failCodes = append(failCodes, failCode)
		}
	}
	global.HTTPFailCodes = failCodes

	insecureForkWanted, err := parseBoolOption(p, "insecure-fork-wanted")
	if err != nil {
		return nil, err
	}
	global.InsecureForkWanted = insecureForkWanted

	insecureSetUIDWanted, err := parseBoolOption(p, "insecure-setuid-wanted")
	if err != nil {
		return nil, err
	}
	global.InsecureSetuidWanted = insecureSetUIDWanted

	limitedQuic, err := parseBoolOption(p, "limited-quic")
	if err != nil {
		return nil, err
	}
	global.LimitedQuic = limitedQuic

	localPeer, err := parseStringOption(p, "localpeer")
	if err != nil {
		return nil, err
	}
	global.Localpeer = localPeer

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
	global.LogSendHostname = globalLogSendHostName

	luaOptions, err := parseLuaOptions(p)
	if err != nil {
		return nil, err
	}
	global.LuaOptions = luaOptions

	masterWorker, err := parseBoolOption(p, "master-worker")
	if err != nil {
		return nil, err
	}
	global.MasterWorker = masterWorker

	mworkerMaxReloads, err := parseInt64POption(p, "mworker-max-reloads")
	if err != nil {
		return nil, err
	}
	global.MworkerMaxReloads = mworkerMaxReloads

	nbthread, err := parseInt64Option(p, "nbthread")
	if err != nil {
		return nil, err
	}
	global.Nbthread = nbthread

	noQuic, err := parseBoolOption(p, "no-quic")
	if err != nil {
		return nil, err
	}
	global.NoQuic = noQuic

	node, err := parseStringOption(p, "node")
	if err != nil {
		return nil, err
	}
	global.Node = node

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
	global.NumaCPUMapping = numaCPUMapping

	ocspUpdate, err := parseOcspUpdateOptions(p)
	if err != nil {
		return nil, err
	}
	global.OcspUpdateOptions = ocspUpdate

	performanceOptions, err := parsePerformanceOptions(p)
	if err != nil {
		return nil, err
	}
	global.PerformanceOptions = performanceOptions

	pidfile, err := parseStringOption(p, "pidfile")
	if err != nil {
		return nil, err
	}
	global.Pidfile = pidfile

	pp2NeverSendLocal, err := parseBoolOption(p, "pp2-never-send-local")
	if err != nil {
		return nil, err
	}
	global.Pp2NeverSendLocal = pp2NeverSendLocal

	preallocFD, err := parseBoolOption(p, "prealloc-fd")
	if err != nil {
		return nil, err
	}
	global.PreallocFd = preallocFD

	setDumpable, err := parseBoolOption(p, "set-dumpable")
	if err != nil {
		return nil, err
	}
	global.SetDumpable = setDumpable

	setcap, err := parseStringOption(p, "setcap")
	if err != nil {
		return nil, err
	}
	global.Setcap = setcap

	sslOptions, err := parseSSLOptions(p)
	if err != nil {
		return nil, err
	}
	global.SslOptions = sslOptions

	statsFile, err := parseStringOption(p, "stats-file")
	if err != nil {
		return nil, err
	}
	global.StatsFile = statsFile

	statsMaxconn, err := parseInt64POption(p, "stats maxconn")
	if err != nil {
		return nil, err
	}
	global.StatsMaxconn = statsMaxconn

	statsTimeout, err := parseTimeoutOption(p, "stats timeout")
	if err != nil {
		return nil, err
	}
	global.StatsTimeout = statsTimeout

	strictLimits, err := parseBoolOption(p, "strict-limits")
	if err != nil {
		return nil, err
	}
	global.StrictLimits = strictLimits

	threadGroups, err := parseInt64Option(p, "thread-groups")
	if err != nil {
		return nil, err
	}
	global.ThreadGroups = threadGroups

	tuneBufferOptions, err := parseTuneBufferOptions(p)
	if err != nil {
		return nil, err
	}
	global.TuneBufferOptions = tuneBufferOptions

	tuneLuaOptions, err := parseTuneLuaOptions(p)
	if err != nil {
		return nil, err
	}
	global.TuneLuaOptions = tuneLuaOptions

	tuneOptions, err := parseTuneOptions(p)
	if err != nil {
		return nil, err
	}
	global.TuneOptions = tuneOptions

	tuneQuicOptions, err := parseTuneQuicOptions(p)
	if err != nil {
		return nil, err
	}
	global.TuneQuicOptions = tuneQuicOptions

	tuneSSLOptions, err := parseTuneSSLOptions(p)
	if err != nil {
		return nil, err
	}
	global.TuneSslOptions = tuneSSLOptions

	tuneVarsOptions, err := parseTuneVarsOptions(p)
	if err != nil {
		return nil, err
	}
	global.TuneVarsOptions = tuneVarsOptions

	tuneZlibOptions, err := parseTuneZlibOptions(p)
	if err != nil {
		return nil, err
	}
	global.TuneZlibOptions = tuneZlibOptions

	uid, err := parseInt64Option(p, "uid")
	if err != nil {
		return nil, err
	}
	global.UID = uid

	user, err := parseStringOption(p, "user")
	if err != nil {
		return nil, err
	}
	global.User = user

	ulimitn, err := parseInt64Option(p, "ulimit-n")
	if err != nil {
		return nil, err
	}
	global.Ulimitn = ulimitn

	global.WarnBlockedTrafficAfter, err = parseTimeoutOption(p, "warn-blocked-traffic-after")
	if err != nil {
		return nil, err
	}

	wurflOptions, err := parseWurflOptions(p)
	if err != nil {
		return nil, err
	}
	global.WurflOptions = wurflOptions

	return global, nil
}

func serializeDebugOptions(p parser.Parser, options *models.DebugOptions) error {
	if options == nil {
		options = &models.DebugOptions{}
	}

	if err := serializeInt64POption(p, "anonkey", options.Anonkey); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "stress-level", options.StressLevel); err != nil {
		return err
	}
	if err := serializeBoolOption(p, "quiet", options.Quiet); err != nil {
		return err
	}
	return serializeBoolOption(p, "zero-warning", options.ZeroWarning)
}

func serializeDeviceAtlasOptions(p parser.Parser, options *models.DeviceAtlasOptions) error {
	if options == nil {
		options = &models.DeviceAtlasOptions{}
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
	return serializeStringOption(p, "deviceatlas-properties-cookie", options.PropertiesCookie)
}

func serializeEnvironmentOptions(p parser.Parser, options *models.EnvironmentOptions) error {
	if options == nil {
		options = &models.EnvironmentOptions{}
	}
	presetEnvs := []types.StringKeyValueC{}
	if len(options.PresetEnvs) > 0 {
		for _, presetEnv := range options.PresetEnvs {
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
	if len(options.SetEnvs) > 0 {
		for _, presetEnv := range options.SetEnvs {
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
	if options.Resetenv == "" {
		resetenv = nil
	} else {
		resetenv.Value = strings.Split(options.Resetenv, " ")
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "resetenv", resetenv); err != nil {
		return err
	}

	unsetenv := &types.StringSliceC{}
	if options.Unsetenv == "" {
		unsetenv = nil
	} else {
		unsetenv.Value = strings.Split(options.Unsetenv, " ")
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "unsetenv", unsetenv); err != nil {
		return err
	}
	return nil
}

func serializeHardenOptions(p parser.Parser, options *models.GlobalHarden) error {
	rppQuic := ""
	rppTCP := ""
	if options != nil {
		if options.RejectPrivilegedPorts != nil {
			rppQuic = options.RejectPrivilegedPorts.Quic
			rppTCP = options.RejectPrivilegedPorts.TCP
		}
	}

	if err := serializeOnOffOption(p, "harden.reject-privileged-ports.quic", rppQuic); err != nil {
		return err
	}
	return serializeOnOffOption(p, "harden.reject-privileged-ports.tcp", rppTCP)
}

func serializeHTTPClientOptions(p parser.Parser, options *models.HTTPClientOptions, configOptions *options.ConfigurationOptions) error {
	if options == nil {
		options = &models.HTTPClientOptions{}
	}

	if err := serializeOnOffOption(p, "httpclient.resolvers.disabled", options.ResolversDisabled); err != nil {
		return err
	}

	if err := serializeStringOption(p, "httpclient.resolvers.id", options.ResolversID); err != nil {
		return err
	}

	pHTTPClientResolversPrefer := &types.HTTPClientResolversPrefer{
		Type: options.ResolversPrefer,
	}
	if options.ResolversPrefer == "" {
		pHTTPClientResolversPrefer = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "httpclient.resolvers.prefer", pHTTPClientResolversPrefer); err != nil {
		return err
	}

	if err := serializeStringOption(p, "httpclient.ssl.ca-file", options.SslCaFile); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "httpclient.retries", options.Retries); err != nil {
		return err
	}

	var pHTTPClientSSLCaFile *types.HTTPClientSSLVerify
	if options.SslVerify != nil {
		pHTTPClientSSLCaFile = &types.HTTPClientSSLVerify{
			Type: *options.SslVerify,
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "httpclient.ssl.verify", pHTTPClientSSLCaFile); err != nil {
		return err
	}

	return serializeTimeoutOption(p, "httpclient.timeout.connect", options.TimeoutConnect, configOptions)
}

func serializeSSLOptions(p parser.Parser, options *models.SslOptions) error { //nolint:gocognit
	if options == nil {
		options = &models.SslOptions{}
	}

	sslEngines := []types.SslEngine{}
	if len(options.SslEngines) > 0 {
		for _, sslEngine := range options.SslEngines {
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

	if err := serializeStringOption(p, "acme.scheduler", options.AcmeScheduler); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ca-base", options.CaBase); err != nil {
		return err
	}

	if err := serializeStringOption(p, "crt-base", options.CrtBase); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-default-bind-ciphers", options.DefaultBindCiphers); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-default-bind-ciphersuites", options.DefaultBindCiphersuites); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-default-bind-client-sigalgs", options.DefaultBindClientSigalgs); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-default-bind-curves", options.DefaultBindCurves); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-default-bind-options", options.DefaultBindOptions); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-default-bind-sigalgs", options.DefaultBindSigalgs); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-default-server-ciphers", options.DefaultServerCiphers); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-default-server-ciphersuites", options.DefaultServerCiphersuites); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-default-server-client-sigalgs", options.DefaultServerClientSigalgs); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-default-server-curves", options.DefaultServerCurves); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-default-server-options", options.DefaultServerOptions); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-default-server-sigalgs", options.DefaultServerSigalgs); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-dh-param-file", options.DhParamFile); err != nil {
		return err
	}

	if err := serializeStringOption(p, "issuers-chain-path", options.IssuersChainPath); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-load-extra-files", options.LoadExtraFiles); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-passphrase-cmd", options.PassphraseCmd); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "maxsslconn", options.Maxsslconn); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "maxsslrate", options.Maxsslrate); err != nil {
		return err
	}

	sslModeAsync := &types.SslModeAsync{}
	if options.ModeAsync != "enabled" {
		sslModeAsync = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "ssl-mode-async", sslModeAsync); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-propquery", options.Propquery); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-provider", options.Provider); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-provider-path", options.ProviderPath); err != nil {
		return err
	}

	if err := serializeStringOption(p, "ssl-server-verify", options.ServerVerify); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "ssl-skip-self-issued-ca", options.SkipSelfIssuedCa); err != nil {
		return err
	}

	return serializeInt64POption(p, "ssl-security-level", options.SecurityLevel)
}

func serializePerformanceOptions(p parser.Parser, options *models.PerformanceOptions, configOptions *options.ConfigurationOptions) error {
	if options == nil {
		options = &models.PerformanceOptions{}
	}

	if err := serializeBoolOption(p, "busy-polling", options.BusyPolling); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "noktls", options.Noktls); err != nil {
		return err
	}

	if err := serializeTimeoutOption(p, "max-spread-checks", options.MaxSpreadChecks, configOptions); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "maxcomprate", options.Maxcomprate); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "maxcompcpuusage", options.Maxcompcpuusage); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "maxconn", options.Maxconn); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "maxconnrate", options.Maxconnrate); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "maxpipes", options.Maxpipes); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "maxsessrate", options.Maxsessrate); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "maxconn", options.Maxconn); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "maxzlibmem", options.Maxzlibmem); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "noepoll", options.Noepoll); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "noevports", options.Noevports); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "nogetaddrinfo", options.Nogetaddrinfo); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "nopoll", options.Nopoll); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "noreuseport", options.Noreuseport); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "nosplice", options.Nosplice); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "nokqueue", options.Nokqueue); err != nil {
		return err
	}

	if err := serializeAutoOnOffOption(p, "profiling.tasks", options.ProfilingTasks); err != nil {
		return err
	}

	if err := serializeOnOffOption(p, "profiling.memory", options.ProfilingMemory); err != nil {
		return err
	}

	if err := serializeStringOption(p, "server-state-file", options.ServerStateFile); err != nil {
		return err
	}

	if err := serializeStringOption(p, "server-state-base", options.ServerStateBase); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "spread-checks", options.SpreadChecks); err != nil {
		return err
	}

	return serializeInt64POption(p, "thread-hard-limit", options.ThreadHardLimit)
}

func serializeLuaOptions(p parser.Parser, options *models.LuaOptions) error {
	if options == nil {
		options = &models.LuaOptions{}
	}

	if err := serializeStringOption(p, "lua-load-per-thread", options.LoadPerThread); err != nil {
		return err
	}

	luaLoads := []types.LuaLoad{}
	for _, lua := range options.Loads {
		ll := types.LuaLoad{
			File: *lua.File,
		}
		luaLoads = append(luaLoads, ll)
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "lua-load", luaLoads); err != nil {
		return err
	}

	luaPrependPath := []types.LuaPrependPath{}
	for _, l := range options.PrependPath {
		lpp := types.LuaPrependPath{
			Path: *l.Path,
			Type: l.Type,
		}
		luaPrependPath = append(luaPrependPath, lpp)
	}
	return p.Set(parser.Global, parser.GlobalSectionName, "lua-prepend-path", luaPrependPath)
}

func SerializeGlobalSection(p parser.Parser, data *models.Global, opt *options.ConfigurationOptions) error { //nolint:gocognit,gocyclo,cyclop,maintidx
	cpuMaps := make([]types.CPUMap, len(data.CPUMaps))
	for i, cpuMap := range data.CPUMaps {
		cm := types.CPUMap{
			Process: *cpuMap.Process,
			CPUSet:  *cpuMap.CPUSet,
		}
		cpuMaps[i] = cm
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "cpu-map", cpuMaps); err != nil {
		return err
	}

	cpuSets := []types.CPUSet{}
	for _, cpuSet := range data.CPUSets {
		cs := types.CPUSet{
			Directive: *cpuSet.Directive,
			Set:       cpuSet.Set,
		}
		cpuSets = append(cpuSets, cs)
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "cpu-set", cpuSets); err != nil {
		return err
	}

	pH1CaseAdjusts := []types.H1CaseAdjust{}
	if len(data.H1CaseAdjusts) > 0 {
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

	sockets := []types.Socket{}
	for _, rAPI := range data.RuntimeAPIs {
		socket := types.Socket{
			Path:   *rAPI.Address,
			Params: []params.BindOption{},
		}
		socket.Params = serializeBindParams(rAPI.BindParams, rAPI.Name, "", opt)
		sockets = append(sockets, socket)
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "stats socket", sockets); err != nil {
		return err
	}

	setVars := []types.SetVar{}
	if len(data.SetVars) > 0 {
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
	if len(data.SetVarFmts) > 0 {
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

	threadGroupLines := []types.ThreadGroup{}
	if len(data.ThreadGroupLines) > 0 {
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

	if err := serializeStringOption(p, "chroot", data.Chroot); err != nil {
		return err
	}

	if err := serializeTimeoutOption(p, "close-spread-time", data.CloseSpreadTime, opt); err != nil {
		return err
	}

	if err := serializeStringOption(p, "cluster-secret", data.ClusterSecret); err != nil {
		return err
	}

	if err := serializeStringOption(p, "cpu-policy", data.CPUPolicy); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "daemon", data.Daemon); err != nil {
		return err
	}

	if err := serializeDebugOptions(p, data.DebugOptions); err != nil {
		return err
	}

	if err := serializeDefaultPath(p, data.DefaultPath); err != nil {
		return err
	}

	if err := serializeStringOption(p, "dns-accept-family", data.DNSAcceptFamily); err != nil {
		return err
	}

	if err := serializeStringOption(p, "shm-stats-file", data.ShmStatsFile); err != nil {
		return err
	}

	if err := serializeInt64POption(p, "shm-stats-file-max-objects", data.ShmStatsFileMaxObjects); err != nil {
		return err
	}

	if err := serializeStringOption(p, "description", data.Description); err != nil {
		return err
	}

	if err := serializeDeviceAtlasOptions(p, data.DeviceAtlasOptions); err != nil {
		return err
	}

	if err := serializeEnvironmentOptions(p, data.EnvironmentOptions); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "expose-deprecated-directives", data.ExposeDeprecatedDirectives); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "expose-experimental-directives", data.ExposeExperimentalDirectives); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "external-check", data.ExternalCheck); err != nil {
		return err
	}

	if err := serializeFiftyOneDegreesOptions(p, data.FiftyOneDegreesOptions); err != nil {
		return err
	}

	if err := serializeTimeoutOption(p, "force-cfg-parser-pause", data.ForceCfgParserPause, opt); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "gid", data.Gid); err != nil {
		return err
	}

	if err := serializeTimeoutOption(p, "grace", data.Grace, opt); err != nil {
		return err
	}

	if err := serializeStringOption(p, "group", data.Group); err != nil {
		return err
	}

	if err := serializeStringOption(p, "h1-case-adjust-file", data.H1CaseAdjustFile); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "h2-workaround-bogus-websocket-clients", data.H2WorkaroundBogusWebsocketClients); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "h1-accept-payload-with-any-method", data.H1AcceptPayloadWithAnyMethod); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "h1-do-not-close-on-insecure-transfer-encoding", data.H1DoNotCloseOnInsecureTransferEncoding); err != nil {
		return err
	}

	if err := serializeTimeoutOption(p, "hard-stop-after", data.HardStopAfter, opt); err != nil {
		return err
	}

	if err := serializeHardenOptions(p, data.Harden); err != nil {
		return err
	}

	if err := serializeHTTPClientOptions(p, data.HTTPClientOptions, opt); err != nil {
		return err
	}

	httpErrCodes := []types.HTTPErrCodes{}
	if len(data.HTTPErrCodes) > 0 {
		for _, errCodes := range data.HTTPErrCodes {
			if errCodes != nil {
				errCode := types.HTTPErrCodes{
					StringC: types.StringC{
						Value: *errCodes.Value,
					},
				}
				httpErrCodes = append(httpErrCodes, errCode)
			}
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "http-err-codes", httpErrCodes); err != nil {
		return err
	}

	httpFailCodes := []types.HTTPFailCodes{}
	if len(data.HTTPFailCodes) > 0 {
		for _, failCodes := range data.HTTPFailCodes {
			if failCodes != nil {
				failCode := types.HTTPFailCodes{
					StringC: types.StringC{
						Value: *failCodes.Value,
					},
				}
				httpFailCodes = append(httpFailCodes, failCode)
			}
		}
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "http-fail-codes", httpFailCodes); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "insecure-fork-wanted", data.InsecureForkWanted); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "insecure-setuid-wanted", data.InsecureSetuidWanted); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "limited-quic", data.LimitedQuic); err != nil {
		return err
	}

	if err := serializeStringOption(p, "localpeer", data.Localpeer); err != nil {
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

	if err := serializeLuaOptions(p, data.LuaOptions); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "master-worker", data.MasterWorker); err != nil {
		return err
	}

	if err := serializeInt64POption(p, "mworker-max-reloads", data.MworkerMaxReloads); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "nbthread", data.Nbthread); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "no-quic", data.NoQuic); err != nil {
		return err
	}

	if err := serializeStringOption(p, "node", data.Node); err != nil {
		return err
	}

	numaCPUMapping := &types.NumaCPUMapping{}
	switch data.NumaCPUMapping {
	case "":
		numaCPUMapping = nil
	case "disabled":
		numaCPUMapping.NoOption = true
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "numa-cpu-mapping", numaCPUMapping); err != nil {
		return err
	}

	if err := serializeOcspUpdateOptions(p, data.OcspUpdateOptions); err != nil {
		return err
	}

	if err := serializePerformanceOptions(p, data.PerformanceOptions, opt); err != nil {
		return err
	}

	if err := serializeStringOption(p, "pidfile", data.Pidfile); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "pp2-never-send-local", data.Pp2NeverSendLocal); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "prealloc-fd", data.PreallocFd); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "set-dumpable", data.SetDumpable); err != nil {
		return err
	}

	if err := serializeStringOption(p, "setcap", data.Setcap); err != nil {
		return err
	}

	if err := serializeSSLOptions(p, data.SslOptions); err != nil {
		return err
	}

	if err := serializeStringOption(p, "stats-file", data.StatsFile); err != nil {
		return err
	}

	if err := serializeInt64POption(p, "stats maxconn", data.StatsMaxconn); err != nil {
		return err
	}

	if err := serializeTimeoutOption(p, "stats timeout", data.StatsTimeout, opt); err != nil {
		return err
	}

	if err := serializeBoolOption(p, "strict-limits", data.StrictLimits); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "thread-groups", data.ThreadGroups); err != nil {
		return err
	}

	if err := serializeTuneBufferOptions(p, data.TuneBufferOptions); err != nil {
		return err
	}

	if err := serializeTuneLuaOptions(p, data.TuneLuaOptions, opt); err != nil {
		return err
	}

	if err := serializeTuneOptions(p, data.TuneOptions, opt); err != nil {
		return err
	}

	if err := serializeTuneQuicOptions(p, data.TuneQuicOptions, opt); err != nil {
		return err
	}

	if err := serializeTuneSSLOptions(p, data.TuneSslOptions, opt); err != nil {
		return err
	}

	if err := serializeTuneVarsOptions(p, data.TuneVarsOptions); err != nil {
		return err
	}

	if err := serializeTuneZlibOptions(p, data.TuneZlibOptions); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "uid", data.UID); err != nil {
		return err
	}

	if err := serializeInt64Option(p, "ulimit-n", data.Ulimitn); err != nil {
		return err
	}

	if err := serializeStringOption(p, "user", data.User); err != nil {
		return err
	}

	if err := serializeTimeoutOption(p, "warn-blocked-traffic-after", data.WarnBlockedTrafficAfter, opt); err != nil {
		return err
	}

	return serializeWurflOptions(p, data.WurflOptions)
}

func serializeDefaultPath(p parser.Parser, data *models.GlobalDefaultPath) error {
	dp := &types.DefaultPath{}
	if data != nil {
		dp.Type = data.Type
		dp.Path = data.Path
	} else {
		dp = nil
	}
	return p.Set(parser.Global, parser.GlobalSectionName, "default-path", dp)
}

func serializeOcspUpdateOptions(p parser.Parser, options *models.OcspUpdateOptions) error {
	if options == nil {
		options = &models.OcspUpdateOptions{}
	}
	disable := ""
	if options.Disable != nil {
		switch *options.Disable {
		case true:
			disable = "enabled"
		case false:
			disable = "disabled"
		}
	}
	if err := serializeOnOffOption(p, "ocsp-update.disable", disable); err != nil {
		return err
	}

	if options.Maxdelay != nil && options.Mindelay != nil && *options.Maxdelay < *options.Mindelay {
		return errors.New("ocsp-update.maxdelay must be greater than ocsp-update.mindelay")
	}

	if err := serializeInt64POption(p, "ocsp-update.mindelay", options.Mindelay); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "ocsp-update.maxdelay", options.Maxdelay); err != nil {
		return err
	}

	addr := ""
	if options.Httpproxy != nil {
		addr = options.Httpproxy.Address
		if options.Httpproxy.Port != nil {
			addr = fmt.Sprintf("%s:%d", addr, *options.Httpproxy.Port)
		}
	}
	if err := serializeStringOption(p, "ocsp-update.httpproxy", addr); err != nil {
		return err
	}

	return serializeOnOffOption(p, "ocsp-update.mode", options.Mode)
}

func serializeWurflOptions(p parser.Parser, options *models.WurflOptions) error {
	if options == nil {
		options = &models.WurflOptions{}
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
	return serializeInt64Option(p, "wurfl-cache-size", options.CacheSize)
}

func serializeFiftyOneDegreesOptions(p parser.Parser, options *models.FiftyOneDegreesOptions) error {
	if options == nil {
		options = &models.FiftyOneDegreesOptions{}
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
	return serializeInt64Option(p, "51degrees-cache-size", options.CacheSize)
}

func serializeTuneBufferOptions(p parser.Parser, options *models.TuneBufferOptions) error {
	if options == nil {
		options = &models.TuneBufferOptions{}
	}
	if err := serializeInt64POption(p, "tune.buffers.limit", options.BuffersLimit); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.buffers.reserve", options.BuffersReserve); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.bufsize", &options.Bufsize); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.bufsize.small", options.BufsizeSmall); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.pipesize", &options.Pipesize); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.rcvbuf.backend", options.RcvbufBackend); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.rcvbuf.client", options.RcvbufClient); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.rcvbuf.frontend", options.RcvbufFrontend); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.rcvbuf.server", options.RcvbufServer); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.recv_enough", &options.RecvEnough); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.sndbuf.backend", options.SndbufBackend); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.sndbuf.client", options.SndbufClient); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.sndbuf.frontend", options.SndbufFrontend); err != nil {
		return err
	}
	return serializeSizeOption(p, "tune.sndbuf.server", options.SndbufServer)
}

func serializeTuneLuaOptions(p parser.Parser, options *models.TuneLuaOptions, configOptions *options.ConfigurationOptions) error {
	if options == nil {
		options = &models.TuneLuaOptions{}
	}
	if err := serializeStringOption(p, "tune.lua.bool-sample-conversion", options.BoolSampleConversion); err != nil {
		return err
	}
	if err := serializeTimeoutOption(p, "tune.lua.burst-timeout", options.BurstTimeout, configOptions); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.lua.forced-yield", options.ForcedYield); err != nil {
		return err
	}
	if err := serializeOnOffOption(p, "tune.lua.log.loggers", options.LogLoggers); err != nil {
		return err
	}
	if err := serializeAutoOnOffOption(p, "tune.lua.log.stderr", options.LogStderr); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.lua.maxmem", options.Maxmem); err != nil {
		return err
	}
	if err := serializeTimeoutOption(p, "tune.lua.service-timeout", options.ServiceTimeout, configOptions); err != nil {
		return err
	}
	if err := serializeTimeoutOption(p, "tune.lua.session-timeout", options.SessionTimeout, configOptions); err != nil {
		return err
	}
	return serializeTimeoutOption(p, "tune.lua.task-timeout", options.TaskTimeout, configOptions)
}

func serializeTuneQuicOptions(p parser.Parser, options *models.TuneQuicOptions, configOptions *options.ConfigurationOptions) error {
	if options == nil {
		options = &models.TuneQuicOptions{}
	}
	if err := serializeInt64POption(p, "tune.quic.frontend.conn-tx-buffers.limit", options.FrontendConnTxBuffersLimit); err != nil {
		return err
	}
	if err := serializeTimeoutOption(p, "tune.quic.frontend.max-idle-timeout", options.FrontendMaxIdleTimeout, configOptions); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.quic.frontend.max-streams-bidi", options.FrontendMaxStreamsBidi); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.quic.frontend.max-tx-mem", options.FrontendMaxTxMemory); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.quic.max-frame-loss", options.MaxFrameLoss); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.quic.reorder-ratio", options.ReorderRatio); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.quic.retry-threshold", options.RetryThreshold); err != nil {
		return err
	}
	value := &types.QuicSocketOwner{Owner: options.SocketOwner}
	if options.SocketOwner == "" {
		value = nil
	}
	if err := p.Set(parser.Global, parser.GlobalSectionName, "tune.quic.socket-owner", value); err != nil {
		return err
	}
	return serializeOnOffOption(p, "tune.quic.zero-copy-fwd-send", options.ZeroCopyFwdSend)
}

func serializeTuneSSLOptions(p parser.Parser, options *models.TuneSslOptions, configOptions *options.ConfigurationOptions) error {
	if options == nil {
		options = &models.TuneSslOptions{}
	}
	if err := serializeInt64POption(p, "tune.ssl.cachesize", options.Cachesize); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.ssl.capture-buffer-size", options.CaptureBufferSize); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.ssl.ssl-ctx-cache-size", options.CtxCacheSize); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.ssl.default-dh-param", options.DefaultDhParam); err != nil {
		return err
	}
	if err := serializeBoolOption(p, "tune.ssl.force-private-cache", options.ForcePrivateCache); err != nil {
		return err
	}
	if err := serializeOnOffOption(p, "tune.ssl.keylog", options.Keylog); err != nil {
		return err
	}
	if err := serializeTimeoutOption(p, "tune.ssl.lifetime", options.Lifetime, configOptions); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.ssl.maxrecord", options.Maxrecord); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.ssl.ocsp-update.maxdelay", options.OcspUpdateMaxDelay); err != nil {
		return err
	}
	return serializeInt64POption(p, "tune.ssl.ocsp-update.mindelay", options.OcspUpdateMinDelay)
}

func serializeTuneVarsOptions(p parser.Parser, options *models.TuneVarsOptions) error {
	if options == nil {
		options = &models.TuneVarsOptions{}
	}
	if err := serializeSizeOption(p, "tune.vars.global-max-size", options.GlobalMaxSize); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.vars.proc-max-size", options.ProcMaxSize); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.vars.reqres-max-size", options.ReqresMaxSize); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.vars.sess-max-size", options.SessMaxSize); err != nil {
		return err
	}
	return serializeSizeOption(p, "tune.vars.txn-max-size", options.TxnMaxSize)
}

func serializeTuneZlibOptions(p parser.Parser, options *models.TuneZlibOptions) error {
	if options == nil {
		options = &models.TuneZlibOptions{}
	}
	if err := serializeInt64Option(p, "tune.zlib.memlevel", options.Memlevel); err != nil {
		return err
	}
	return serializeInt64Option(p, "tune.zlib.windowsize", options.Windowsize)
}

func serializeTuneOptions(p parser.Parser, options *models.TuneOptions, configOptions *options.ConfigurationOptions) error { //nolint:gocognit,gocyclo,cyclop,maintidx
	if options == nil {
		options = &models.TuneOptions{}
	}
	if err := serializeOnOffOption(p, "tune.applet.zero-copy-forwarding", options.AppletZeroCopyForwarding); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.comp.maxlevel", options.CompMaxlevel); err != nil {
		return err
	}
	if err := serializeBoolOption(p, "tune.disable-fast-forward", options.DisableFastForward); err != nil {
		return err
	}
	if err := serializeBoolOption(p, "tune.disable-zero-copy-forwarding", options.DisableZeroCopyForwarding); err != nil {
		return err
	}
	if err := serializeStringOption(p, "tune.epoll.mask-events", strings.Join(options.EpollMaskEvents, ",")); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.events.max-events-at-once", options.EventsMaxEventsAtOnce); err != nil {
		return err
	}
	if err := serializeBoolOption(p, "tune.fail-alloc", options.FailAlloc); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.glitches.kill.cpu-usage", options.GlitchesKillCPUUsage); err != nil {
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
	if err := serializeTimeoutOption(p, "tune.idletimer", options.Idletimer, configOptions); err != nil {
		return err
	}
	if err := serializeListenerDefaultShards(p, "tune.listener.default-shards", options.ListenerDefaultShards); err != nil {
		return err
	}
	if err := serializeOnOffOption(p, "tune.listener.multi-queue", options.ListenerMultiQueue); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.max-checks-per-thread", options.MaxChecksPerThread); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.max-rules-at-once", options.MaxRulesAtOnce); err != nil {
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
	if err := serializeInt64Option(p, "tune.pool-high-fd-ratio", options.PoolHighFdRatio); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.pool-low-fd-ratio", options.PoolLowFdRatio); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.ring.queues", options.RingQueues); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.runqueue-depth", options.RunqueueDepth); err != nil {
		return err
	}
	if err := serializeOnOffOption(p, "tune.sched.low-latency", options.SchedLowLatency); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.stick-counters", options.StickCounters); err != nil {
		return err
	}
	if err := serializeStringOption(p, "tune.takeover-other-tg-connections", options.TakeoverOtherTgConnections); err != nil {
		return err
	}
	if err := serializeOnOffOption(p, "tune.fd.edge-triggered", options.FdEdgeTriggered); err != nil {
		return err
	}
	if err := serializeOnOffOption(p, "tune.h1.zero-copy-fwd-recv", options.H1ZeroCopyFwdRecv); err != nil {
		return err
	}
	if err := serializeOnOffOption(p, "tune.h1.zero-copy-fwd-send", options.H1ZeroCopyFwdSend); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.h2.be.glitches-threshold", options.H2BeGlitchesThreshold); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.h2.be.initial-window-size", options.H2BeInitialWindowSize); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.h2.be.max-concurrent-streams", options.H2BeMaxConcurrentStreams); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.h2.be.rxbuf", options.H2BeRxbuf); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.h2.fe.glitches-threshold", options.H2FeGlitchesThreshold); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.h2.fe.initial-window-size", options.H2FeInitialWindowSize); err != nil {
		return err
	}
	if err := serializeInt64Option(p, "tune.h2.fe.max-concurrent-streams", options.H2FeMaxConcurrentStreams); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.h2.fe.max-total-streams", options.H2FeMaxTotalStreams); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.h2.fe.rxbuf", options.H2FeRxbuf); err != nil {
		return err
	}
	if err := serializeOnOffOption(p, "tune.h2.zero-copy-fwd-send", options.H2ZeroCopyFwdSend); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.notsent-lowat.client", options.NotsentLowatClient); err != nil {
		return err
	}
	if err := serializeSizeOption(p, "tune.notsent-lowat.server", options.NotsentLowatServer); err != nil {
		return err
	}
	if err := serializeOnOffOption(p, "tune.pt.zero-copy-forwarding", options.PtZeroCopyForwarding); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.renice.runtime", options.ReniceRuntime); err != nil {
		return err
	}
	if err := serializeInt64POption(p, "tune.renice.startup", options.ReniceStartup); err != nil {
		return err
	}
	return nil
}

func serializeTimeoutOption(p parser.Parser, option string, data *int64, opt *options.ConfigurationOptions) error {
	var value *types.StringC
	if data == nil {
		value = nil
	} else {
		value = &types.StringC{Value: misc.SerializeTime(*data, opt.PreferredTimeSuffix)}
	}
	return p.Set(parser.Global, parser.GlobalSectionName, option, value)
}

func serializeSizeOption(p parser.Parser, option string, data *int64) error {
	var value *types.StringC
	if data == nil || *data == 0 {
		value = nil
	} else {
		value = &types.StringC{Value: misc.SerializeSize(*data)}
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
	case models.TuneOptionsListenerDefaultShardsByDashProcess:
		value = &types.StringC{Value: models.TuneOptionsListenerDefaultShardsByDashProcess}
	case models.TuneOptionsListenerDefaultShardsByDashThread:
		value = &types.StringC{Value: models.TuneOptionsListenerDefaultShardsByDashThread}
	case models.TuneOptionsListenerDefaultShardsByDashGroup:
		value = &types.StringC{Value: models.TuneOptionsListenerDefaultShardsByDashGroup}
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

func parseTimeoutOption(p parser.Parser, option string) (*int64, error) {
	data, err := p.Get(parser.Global, parser.GlobalSectionName, option)
	if err == nil {
		value, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError(option)
		}
		return misc.ParseTimeout(value.Value), nil
	}
	if errors.Is(err, parser_errors.ErrFetch) {
		return nil, nil //nolint:nilnil
	}
	return nil, err
}

func parseSizeOption(p parser.Parser, option string) (*int64, error) {
	data, err := p.Get(parser.Global, parser.GlobalSectionName, option)
	if err == nil {
		value, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError(option)
		}
		return misc.ParseSize(value.Value), nil
	}
	if errors.Is(err, parser_errors.ErrFetch) {
		return nil, nil //nolint:nilnil
	}
	return nil, err
}

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
		return nil, nil //nolint:nilnil
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
		case models.TuneOptionsListenerDefaultShardsByDashProcess:
			return models.TuneOptionsListenerDefaultShardsByDashProcess, nil
		case models.TuneOptionsListenerDefaultShardsByDashThread:
			return models.TuneOptionsListenerDefaultShardsByDashThread, nil
		case models.TuneOptionsListenerDefaultShardsByDashGroup:
			return models.TuneOptionsListenerDefaultShardsByDashGroup, nil
		default:
			return "", fmt.Errorf("unsupported value for %s: %s", option, value.Value)
		}
	}
	if errors.Is(err, parser_errors.ErrFetch) {
		return "", nil
	}
	return "", err
}
