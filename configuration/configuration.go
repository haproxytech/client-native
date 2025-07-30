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
	"reflect"
	"strconv"
	"strings"
	"sync"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/common"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/config-parser/params"
	"github.com/haproxytech/client-native/v6/config-parser/parsers"
	stats "github.com/haproxytech/client-native/v6/config-parser/parsers/stats/settings"
	"github.com/haproxytech/client-native/v6/config-parser/types"

	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

const (
	BackendParentName    = "backend"
	FrontendParentName   = "frontend"
	DefaultsParentName   = "defaults"
	LogForwardParentName = "log_forward"
	PeersParentName      = "peers"
	RingParentName       = "ring"
	GlobalParentName     = "global"
	FCGIAppParentName    = "fcgi-app"
	ResolverParentName   = "resolvers"
	CrtStoreParentName   = "crt-store"
	TracesParentName     = "traces"
	LogProfileParentName = "log-profile"
)

// ClientParams is just a placeholder for all client options
type ClientParams struct {
	ConfigurationFile string
	Haproxy           string
	TransactionDir    string

	// ValidateCmd allows specifying a custom script to validate the transaction file.
	// The injected environment variable DATAPLANEAPI_TRANSACTION_FILE must be used to get the location of the file.
	ValidateCmd               string
	ValidateConfigFilesBefore []string
	ValidateConfigFilesAfter  []string
	PreferredTimeSuffix       string
	BackupsNumber             int
	UseValidation             bool
	PersistentTransactions    bool
	ValidateConfigurationFile bool
	MasterWorker              bool
	SkipFailedTransactions    bool
	UseMd5Hash                bool
}

// Client configuration client
// Parser is the config parser instance that loads "master" configuration file on Init
// and when transaction is committed it gets replaced with the parser from parsers map.
// parsers map contains a config parser for each transaction, which loads data from
// transaction files on StartTransaction, and deletes on CommitTransaction. We save
// data to file on every change for persistence.
type client struct {
	parser   parser.Parser
	parsers  map[string]parser.Parser
	services map[string]*Service
	Transaction
	clientMu sync.Mutex
}

// SetValidateConfigFiles set before and after validation files
func (c *client) SetValidateConfigFiles(before, after []string) {
	c.clientMu.Lock()
	defer c.clientMu.Unlock()
	c.Transaction.ValidateConfigFilesBefore = before
	c.Transaction.ValidateConfigFilesAfter = after
}

// HasParser checks whether transaction exists in parser
func (c *client) HasParser(transactionID string) bool {
	c.clientMu.Lock()
	defer c.clientMu.Unlock()
	_, ok := c.parsers[transactionID]
	return ok
}

// GetParser returns a parser for given transactionID, if transactionID is "", it returns "master" parser
func (c *client) GetParser(transactionID string) (parser.Parser, error) {
	if transactionID == "" {
		return c.parser, nil
	}
	c.clientMu.Lock()
	p, ok := c.parsers[transactionID]
	c.clientMu.Unlock()
	if !ok {
		return nil, NewConfError(ErrTransactionDoesNotExist, transactionID)
	}
	return p, nil
}

// AddParser adds parser to parser map
func (c *client) AddParser(transactionID string) error {
	if transactionID == "" {
		return NewConfError(ErrValidationError, "Not a valid transaction")
	}
	c.clientMu.Lock()
	_, ok := c.parsers[transactionID]
	c.clientMu.Unlock()
	if ok {
		return NewConfError(ErrTransactionAlreadyExists, transactionID)
	}

	parserOptions := []parser_options.ParserOption{}
	if c.ConfigurationOptions.UseMd5Hash {
		parserOptions = append(parserOptions, parser_options.UseMd5Hash)
	}
	if c.noNamedDefaultsFrom {
		parserOptions = append(parserOptions, parser_options.NoNamedDefaultsFrom)
	}

	var tFile string
	var err error
	if c.PersistentTransactions {
		tFile, err = c.GetTransactionFile(transactionID)
		if err != nil {
			return err
		}
	} else {
		tFile = c.ConfigurationFile
	}
	parserOptions = append(parserOptions, parser_options.Path(tFile))
	p, err := parser.New(parserOptions...)
	if err != nil {
		return NewConfError(ErrCannotReadConfFile, "Cannot read "+tFile)
	}
	c.clientMu.Lock()
	c.parsers[transactionID] = p
	c.clientMu.Unlock()
	return nil
}

// DeleteParser deletes parser from parsers map
func (c *client) DeleteParser(transactionID string) error {
	if transactionID == "" {
		return NewConfError(ErrValidationError, "Not a valid transaction")
	}
	c.clientMu.Lock()
	_, ok := c.parsers[transactionID]
	c.clientMu.Unlock()
	if !ok {
		return NewConfError(ErrTransactionDoesNotExist, transactionID)
	}
	c.clientMu.Lock()
	delete(c.parsers, transactionID)
	c.clientMu.Unlock()
	return nil
}

// CommitParser commits transaction parser, deletes it from parsers map, and replaces master Parser
func (c *client) CommitParser(transactionID string) error {
	if transactionID == "" {
		return NewConfError(ErrValidationError, "Not a valid transaction")
	}
	c.clientMu.Lock()
	p, ok := c.parsers[transactionID]
	c.clientMu.Unlock()
	if !ok {
		return NewConfError(ErrTransactionDoesNotExist, transactionID)
	}
	c.parser = p
	delete(c.parsers, transactionID)
	return nil
}

// GetVersion returns configuration file version
func (c *client) GetVersion(transactionID string) (int64, error) {
	return c.getVersion(transactionID)
}

func (c *client) getVersion(transactionID string) (int64, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, NewConfError(ErrCannotReadVersion, err.Error())
	}

	data, _ := p.Get(parser.Comments, parser.CommentsSectionName, "# _version", true)
	ver, _ := data.(*types.ConfigVersion)
	return ver.Value, nil
}

func (c *client) IncrementVersion() error {
	data, _ := c.parser.Get(parser.Comments, parser.CommentsSectionName, "# _version", true)
	ver, _ := data.(*types.ConfigVersion)
	ver.Value++

	if err := c.parser.Save(c.ConfigurationFile); err != nil {
		return NewConfError(ErrCannotSetVersion, err.Error())
	}
	return nil
}

func (c *client) LoadData(filename string) error {
	err := c.parser.LoadData(filename)
	if err != nil {
		return NewConfError(ErrCannotReadConfFile, "cannot read "+filename)
	}
	return nil
}

func (c *client) Save(transactionFile, transactionID string) error {
	if transactionID == "" {
		return c.parser.Save(transactionFile)
	}
	p, err := c.GetParser(transactionID)
	if err != nil {
		return err
	}
	return p.Save(transactionFile)
}

// ParseSection sets the fields of the section based on the provided parser
func ParseSection(object interface{}, section parser.Section, pName string, p parser.Parser) error {
	sp := &SectionParser{
		Object:  object,
		Section: section,
		Name:    pName,
		Parser:  p,
	}
	return sp.Parse()
}

func NewParseSection(section parser.Section, pName string, p parser.Parser) *SectionParser {
	return &SectionParser{
		Section: section,
		Name:    pName,
		Parser:  p,
	}
}

// SectionParser is used set fields of a section based on the provided parser
type SectionParser struct {
	Object  interface{}
	Parser  parser.Parser
	Section parser.Section
	Name    string
}

// Parse parses the sections fields and sets their values with the data from the parser
func (s *SectionParser) Parse() error {
	objValue := reflect.ValueOf(s.Object).Elem()
	for i := range objValue.NumField() {
		typeField := objValue.Type().Field(i)
		field := objValue.FieldByName(typeField.Name)
		val := s.parseField(typeField.Name)
		if val != nil {
			if field.Kind() == reflect.Bool {
				if reflect.ValueOf(val).Kind() == reflect.String {
					if val == "enabled" {
						field.Set(reflect.ValueOf(true))
					} else {
						field.Set(reflect.ValueOf(false))
					}
				} else if reflect.ValueOf(val).Kind() == reflect.Bool {
					field.Set(reflect.ValueOf(val))
				}
			} else {
				field.Set(reflect.ValueOf(val))
			}
		}
	}

	return nil
}

func (s *SectionParser) parseField(fieldName string) interface{} {
	if match, data := s.checkSpecialFields(fieldName); match {
		return data
	}
	if match, data := s.checkTimeouts(fieldName); match {
		return data
	}
	if match, data := s.checkOptions(fieldName); match {
		return data
	}
	if match, data := s.checkSingleLine(fieldName); match {
		return data
	}
	return nil
}

func (s *SectionParser) checkSpecialFields(fieldName string) (bool, interface{}) { //nolint:gocyclo,cyclop
	switch fieldName {
	case "Shards":
		return true, s.shards()
	case "From":
		return true, s.from()
	case "MonitorFail":
		return true, s.monitorFail()
	case "MonitorURI":
		return true, s.monitorURI()
	case "StatsOptions":
		return true, s.statsOptions()
	case "Forwardfor":
		return true, s.forwardfor()
	case "Redispatch":
		return true, s.redispatch()
	case "Balance":
		return true, s.balance()
	case "PersistRule":
		return true, s.persistRule()
	case "Cookie":
		return true, s.cookie()
	case "HashType":
		return true, s.hashType()
	case "ErrorFiles":
		return true, s.errorFiles()
	case "ErrorFilesFromHTTPErrors":
		return true, s.errorfilesFromHTTPErrors()
	case "DefaultServer":
		return true, s.defaultServer()
	case "LoadServerStateFromFile":
		return true, s.loadServerStateFromFile()
	case "StickTable":
		return true, s.stickTable()
	case "AdvCheck":
		return true, s.advCheck()
	case "Logasap":
		return true, s.logasap()
	case "Allbackups":
		return true, s.allbackups()
	case "ExternalCheck":
		return true, s.externalCheck()
	case "ExternalCheckPath":
		return true, s.externalCheckPath()
	case "ExternalCheckCommand":
		return true, s.externalCheckCommand()
	case "DefaultBackend":
		return true, s.defaultBackend()
	case "Clflog":
		return true, s.clflog()
	case "Httplog":
		return true, s.httplog()
	case "HTTPReuse":
		return true, s.httpReuse()
	case "UniqueIDFormat":
		return true, s.uniqueIDFormat()
	case "UniqueIDHeader":
		return true, s.uniqueIDHeader()
	case "HTTPConnectionMode":
		return true, s.httpConnectionMode()
	case "Compression":
		return true, s.compression()
	case "ClitcpkaIdle":
		return true, s.clitcpkaIdle()
	case "ClitcpkaIntvl":
		return true, s.clitcpkaIntvl()
	case "SrvtcpkaIdle":
		return true, s.srvtcpkaIdle()
	case "SrvtcpkaIntvl":
		return true, s.srvtcpkaIntvl()
	case "EmailAlert":
		return true, s.emailAlert()
	case "ServerStateFileName":
		return true, s.serverStateFileName()
	case "UseFCGIApp":
		return true, s.useFcgiApp()
	case "Description":
		return true, s.description()
	case "Errorloc302":
		return true, s.errorloc302()
	case "Errorloc303":
		return true, s.errorloc303()
	case "HTTPRestrictReqHdrNames":
		return true, s.httpRestirctReqHdrNames()
	case "DefaultBind":
		return true, s.defaultBind()
	case "HTTPSendNameHeader":
		return true, s.httpSendNameHeader()
	case "ForcePersistList":
		return true, s.forcePersistList()
	case "IgnorePersistList":
		return true, s.ignorePersistList()
	case "Source":
		return true, s.source()
	case "Originalto":
		return true, s.originalto()
	case "LogSteps":
		return true, s.logSteps()
	default:
		return false, nil
	}
}

func (s *SectionParser) checkTimeouts(fieldName string) (bool, interface{}) {
	if strings.HasSuffix(fieldName, "Timeout") {
		if pName := translateTimeout(fieldName); s.Parser.HasParser(s.Section, pName) {
			data, err := s.get(pName, false)
			if err != nil {
				return true, nil
			}
			timeout := data.(*types.SimpleTimeout) //nolint:forcetypeassert
			return true, misc.ParseTimeout(timeout.Value)
		}
	}
	return false, nil
}

func (s *SectionParser) checkSingleLine(fieldName string) (bool, interface{}) {
	if pName := misc.DashCase(fieldName); s.Parser.HasParser(s.Section, pName) {
		data, err := s.get(pName, false)
		if err != nil {
			return true, nil
		}
		return true, parseOption(data)
	}
	return false, nil
}

func (s *SectionParser) checkOptions(fieldName string) (bool, interface{}) {
	if pName := "option " + misc.DashCase(fieldName); s.Parser.HasParser(s.Section, pName) {
		data, err := s.get(pName, false)
		if err != nil {
			return true, nil
		}
		return true, parseOption(data)
	}
	return false, nil
}

func (s *SectionParser) get(attribute string, createIfNotExists ...bool) (common.ParserData, error) {
	return s.Parser.Get(s.Section, s.Name, attribute, createIfNotExists...)
}

func (s *SectionParser) from() interface{} {
	from, err := s.Parser.SectionsDefaultsFromGet(s.Section, s.Name)
	if err != nil {
		return ""
	}
	return from
}

func (s *SectionParser) httpConnectionMode() interface{} {
	data, err := s.get("option http-tunnel", false)
	if err == nil {
		d := data.(*types.SimpleOption) //nolint:forcetypeassert
		if !d.NoOption {
			return "http-tunnel"
		}
	}

	data, err = s.get("option httpclose", false)
	if err == nil {
		d, ok := data.(*types.SimpleOption)
		if !ok {
			return misc.CreateTypeAssertError("option httpclose")
		}
		if !d.NoOption {
			return "httpclose"
		}
	}
	// deprecated option, alias for httpclose
	data, err = s.get("option forceclose", false)
	if err == nil {
		d := data.(*types.SimpleOption) //nolint:forcetypeassert
		if !d.NoOption {
			return "httpclose"
		}
	}

	data, err = s.get("option http-server-close", false)
	if err == nil {
		d := data.(*types.SimpleOption) //nolint:forcetypeassert
		if !d.NoOption {
			return "http-server-close"
		}
	}

	data, err = s.get("option http-keep-alive", false)
	if err == nil {
		d := data.(*types.SimpleOption) //nolint:forcetypeassert
		if !d.NoOption {
			return "http-keep-alive"
		}
	}
	return nil
}

func (s *SectionParser) uniqueIDHeader() interface{} {
	_, e := s.get("unique-id-format")
	if e != nil {
		return nil
	}
	data, err := s.get("unique-id-header")
	if err == nil {
		d := data.(*types.UniqueIDHeader) //nolint:forcetypeassert
		return d.Name
	}
	return nil
}

func (s *SectionParser) uniqueIDFormat() interface{} {
	data, err := s.get("unique-id-format")
	if err == nil {
		d := data.(*types.UniqueIDFormat)
		return d.LogFormat
	}
	return nil
}

func (s *SectionParser) httpReuse() interface{} {
	data, err := s.get("http-reuse", false)
	if err == nil {
		d := data.(*types.HTTPReuse)
		return d.ShareType
	}
	return nil
}

func (s *SectionParser) httplog() interface{} {
	data, err := s.get("option httplog", false)
	if err == nil {
		d := data.(*types.OptionHTTPLog)
		if !d.NoOption {
			return !d.Clf
		}
	}
	return nil
}

func (s *SectionParser) clflog() interface{} {
	data, err := s.get("option httplog", false)
	if err == nil {
		d := data.(*types.OptionHTTPLog)
		if !d.NoOption {
			return d.Clf
		}
	}
	return nil
}

func (s *SectionParser) defaultBackend() interface{} {
	data, err := s.get("default_backend", false)
	if err != nil {
		return nil
	}
	bck := data.(*types.StringC)
	return bck.Value
}

func (s *SectionParser) externalCheckCommand() interface{} {
	data, err := s.get("external-check command", false)
	if err != nil {
		return nil
	}
	d := data.(*types.ExternalCheckCommand)
	return d.Command
}

func (s *SectionParser) externalCheckPath() interface{} {
	data, err := s.get("external-check path", false)
	if err != nil {
		return nil
	}
	d := data.(*types.ExternalCheckPath)
	return d.Path
}

func (s *SectionParser) externalCheck() interface{} {
	data, err := s.get("option external-check", false)
	if err != nil {
		return nil
	}
	if data.(*types.SimpleOption).NoOption {
		return "disabled"
	}
	return "enabled"
}

func (s *SectionParser) allbackups() interface{} {
	data, err := s.get("option allbackups", false)
	if err != nil {
		return nil
	}
	if data.(*types.SimpleOption).NoOption {
		return "disabled"
	}
	return "enabled"
}

func (s *SectionParser) logasap() interface{} {
	data, err := s.get("option logasap", false)
	if err != nil {
		return nil
	}
	if data.(*types.SimpleOption).NoOption {
		return "disabled"
	}
	return "enabled"
}

func (s *SectionParser) useFcgiApp() interface{} {
	_, e := s.get("use-fcgi-app")
	if e != nil {
		return nil
	}
	data, err := s.get("use-fcgi-app")
	if err == nil {
		d := data.(*types.UseFcgiApp) //nolint:forcetypeassert
		return d.Name
	}
	return nil
}

func (s *SectionParser) advCheck() interface{} {
	if found, data := s.getSslChkData(); found {
		return data
	}

	if found, data := s.getSMTPChkData(); found {
		return data
	}

	if found, data := s.getLdapCheckData(); found {
		return data
	}

	if found, data := s.getMysqlCheckData(); found {
		return data
	}

	if found, data := s.getPgsqlCheckData(); found {
		return data
	}

	if found, data := s.getTCPCheckData(); found {
		return data
	}

	if found, data := s.getRedisCheckData(); found {
		return data
	}

	if found, data := s.getHttpchkData(); found {
		return data
	}

	return nil
}

func (s *SectionParser) getSslChkData() (bool, interface{}) {
	data, err := s.get("option ssl-hello-chk", false)
	if err == nil {
		d := data.(*types.SimpleOption)
		if !d.NoOption {
			return true, "ssl-hello-chk"
		}
	}
	return false, nil
}

func (s *SectionParser) getSMTPChkData() (bool, interface{}) {
	data, err := s.get("option smtpchk", false)
	if err == nil {
		d := data.(*types.OptionSmtpchk)
		if !d.NoOption {
			s.setField("SmtpchkParams", &models.SmtpchkParams{
				Hello:  d.Hello,
				Domain: d.Domain,
			})
			return true, "smtpchk"
		}
	}
	return false, nil
}

func (s *SectionParser) getLdapCheckData() (bool, interface{}) {
	data, err := s.get("option ldap-check", false)
	if err == nil {
		d := data.(*types.SimpleOption)
		if !d.NoOption {
			return true, "ldap-check"
		}
	}
	return false, nil
}

func (s *SectionParser) getMysqlCheckData() (bool, interface{}) {
	data, err := s.get("option mysql-check", false)
	if err == nil {
		d := data.(*types.OptionMysqlCheck)
		if !d.NoOption {
			s.setField("MysqlCheckParams", &models.MysqlCheckParams{
				ClientVersion: d.ClientVersion,
				Username:      d.User,
			})
			return true, "mysql-check"
		}
	}
	return false, nil
}

func (s *SectionParser) getPgsqlCheckData() (bool, interface{}) {
	data, err := s.get("option pgsql-check", false)
	if err == nil {
		d := data.(*types.OptionPgsqlCheck)
		if !d.NoOption {
			s.setField("PgsqlCheckParams", &models.PgsqlCheckParams{
				Username: d.User,
			})
			return true, "pgsql-check"
		}
	}
	return false, nil
}

func (s *SectionParser) getTCPCheckData() (bool, interface{}) {
	data, err := s.get("option tcp-check", false)
	if err == nil {
		d := data.(*types.SimpleOption)
		if !d.NoOption {
			return true, "tcp-check"
		}
	}
	return false, nil
}

func (s *SectionParser) getRedisCheckData() (bool, interface{}) {
	data, err := s.get("option redis-check", false)
	if err == nil {
		d := data.(*types.SimpleOption)
		if !d.NoOption {
			return true, "redis-check"
		}
	}
	return false, nil
}

func (s *SectionParser) getHttpchkData() (bool, interface{}) {
	data, err := s.get("option httpchk", false)
	if err == nil {
		d := data.(*types.OptionHttpchk)
		if !d.NoOption {
			s.setField("HttpchkParams", &models.HttpchkParams{
				Method:  d.Method,
				URI:     d.URI,
				Version: d.Version,
				Host:    d.Host,
			})
			return true, "httpchk"
		}
	}
	return false, nil
}

func (s *SectionParser) setField(fieldName string, data interface{}) {
	objValue := reflect.ValueOf(s.Object).Elem()
	field := objValue.FieldByName(fieldName)
	field.Set(reflect.ValueOf(data))
}

func (s *SectionParser) stickTable() interface{} {
	data, err := s.get("stick-table", false)
	if err != nil {
		return nil
	}
	d := data.(*types.StickTable)
	bst := &models.ConfigStickTable{}

	if d == nil {
		return nil
	}
	bst.Type = d.Type
	bst.Size = misc.ParseSize(d.Size)
	bst.Store = d.Store
	bst.Expire = misc.ParseTimeout(d.Expire)
	bst.Peers = d.Peers

	k, err := strconv.ParseInt(d.Length, 10, 64)
	if err == nil {
		bst.Keylen = &k
	}
	if d.NoPurge {
		bst.Nopurge = true
	}
	if d.SrvKey != "" {
		bst.Srvkey = &d.SrvKey
	}
	if d.WriteTo != "" {
		bst.WriteTo = &d.WriteTo
	}
	return bst
}

func (s *SectionParser) defaultServer() interface{} {
	data, err := s.get("default-server", false)
	if err != nil {
		return nil
	}
	d := data.([]types.DefaultServer)
	dServer := &models.DefaultServer{}
	for _, ds := range d {
		parseServerParams(ds.Params, &dServer.ServerParams)
	}
	return dServer
}

func (s *SectionParser) loadServerStateFromFile() interface{} {
	data, err := s.get("load-server-state-from-file", false)
	if err == nil {
		d := data.(*types.LoadServerStateFromFile)
		return d.Argument
	}
	return nil
}

func (s *SectionParser) errorFiles() interface{} {
	data, err := s.get("errorfile", false)
	if err != nil {
		return nil
	}
	d := data.([]types.ErrorFile)

	dEFiles := []*models.Errorfile{}
	for _, ef := range d {
		dEFile := &models.Errorfile{}
		code, err := strconv.ParseInt(ef.Code, 10, 64)
		if err != nil {
			continue
		}
		dEFile.Code = code
		dEFile.File = ef.File
		dEFiles = append(dEFiles, dEFile)
	}
	if len(dEFiles) == 0 {
		return nil
	}
	return dEFiles
}

func (s *SectionParser) errorfilesFromHTTPErrors() interface{} {
	data, err := s.get("errorfiles", false)
	if err != nil {
		return nil
	}
	d := data.([]types.ErrorFiles)

	dEFiles := make([]*models.Errorfiles, len(d))
	for i, ef := range d {
		dEFile := &models.Errorfiles{}
		dEFile.Codes = ef.Codes
		dEFile.Name = ef.Name
		dEFiles[i] = dEFile
	}
	if len(dEFiles) == 0 {
		return nil
	}
	return dEFiles
}

func (s *SectionParser) hashType() interface{} {
	data, err := s.get("hash-type", false)
	if err != nil {
		return nil
	}
	d := data.(*types.HashType)
	return &models.HashType{
		Method:   d.Method,
		Function: d.Function,
		Modifier: d.Modifier,
	}
}

func (s *SectionParser) cookie() interface{} {
	data, err := s.get("cookie", false)
	if err != nil {
		return nil
	}
	d := data.(*types.Cookie)
	domains := make([]*models.Domain, len(d.Domain))
	for i, domain := range d.Domain {
		domains[i] = &models.Domain{Value: domain}
	}
	if len(d.Domain) == 0 {
		domains = nil
	}
	attrs := make([]*models.Attr, len(d.Attr))
	for i, attr := range d.Attr {
		attrs[i] = &models.Attr{Value: attr}
	}
	if len(d.Attr) == 0 {
		attrs = nil
	}
	return &models.Cookie{
		Attrs:    attrs,
		Domains:  domains,
		Dynamic:  d.Dynamic,
		Httponly: d.Httponly,
		Indirect: d.Indirect,
		Maxidle:  d.Maxidle,
		Maxlife:  d.Maxlife,
		Name:     &d.Name,
		Nocache:  d.Nocache,
		Postonly: d.Postonly,
		Preserve: d.Preserve,
		Type:     d.Type,
		Secure:   d.Secure,
	}
}

func (s *SectionParser) persistRule() interface{} {
	data, err := s.get("persist", false)
	if err != nil {
		return nil
	}
	d := data.(*types.Persist)
	p := &models.PersistRule{
		Type: &d.Type,
	}
	switch prm := d.Params.(type) {
	case *params.PersistRdpCookie:
		p.RdpCookieName = prm.Name
	default:
	}
	return p
}

func (s *SectionParser) balance() interface{} {
	data, err := s.get("balance", false)
	if err != nil {
		return nil
	}
	d := data.(*types.Balance)
	b := &models.Balance{
		Algorithm: &d.Algorithm,
	}
	switch prm := d.Params.(type) {
	case *params.BalanceHdr:
		b.HdrName = prm.Name
		b.HdrUseDomainOnly = prm.UseDomainOnly
	case *params.BalanceRandom:
		b.RandomDraws = prm.Draws
	case *params.BalanceRdpCookie:
		b.RdpCookieName = prm.Name
	case *params.BalanceURI:
		b.URIDepth = prm.Depth
		b.URILen = prm.Len
		b.URIWhole = prm.Whole
		b.URIPathOnly = prm.PathOnly
	case *params.BalanceURLParam:
		b.URLParam = prm.Param
		b.URLParamCheckPost = prm.CheckPost
		b.URLParamMaxWait = prm.MaxWait
	case *params.BalanceHash:
		b.HashExpression = prm.Expression
	}
	return b
}

func (s *SectionParser) redispatch() interface{} {
	data, err := s.get("option redispatch", false)
	if err != nil {
		return nil
	}
	d := data.(*types.OptionRedispatch)
	br := &models.Redispatch{
		Interval: d.Interval,
	}
	if d.NoOption {
		d := "disabled"
		br.Enabled = &d
	} else {
		e := "enabled"
		br.Enabled = &e
	}
	return br
}

func (s *SectionParser) forwardfor() interface{} {
	data, err := s.get("option forwardfor", false)
	if err != nil {
		return nil
	}
	d := data.(*types.OptionForwardFor)
	enabled := "enabled"
	bff := &models.Forwardfor{
		Except:  d.Except,
		Header:  d.Header,
		Ifnone:  d.IfNone,
		Enabled: &enabled,
	}
	return bff
}

func (s *SectionParser) emailAlert() interface{} {
	data, err := s.get("email-alert", false)
	if err != nil {
		return nil
	}
	tokens := data.([]types.EmailAlert)
	ea := &models.EmailAlert{}
	for _, tok := range tokens {
		t := tok
		switch t.Attribute {
		case "from":
			ea.From = &t.Value
		case "to":
			ea.To = &t.Value
		case "level":
			ea.Level = t.Value
		case "myhostname":
			ea.Myhostname = t.Value
		case "mailers":
			ea.Mailers = &t.Value
		}
	}
	return ea
}

func (s *SectionParser) statsOptions() interface{} { //nolint:gocognit
	data, err := s.get("stats", false)
	if err != nil {
		return nil
	}
	ss := data.([]types.StatsSettings)
	opt := &models.StatsOptions{}
	for _, stat := range ss {
		switch v := stat.(type) {
		case *stats.OneWord:
			if v.Name == "enable" {
				opt.StatsEnable = true
			}
			if v.Name == "hide-version" {
				opt.StatsHideVersion = true
			}
			if v.Name == "show-legends" {
				opt.StatsShowLegends = true
			}
			if v.Name == "show-modules" {
				opt.StatsShowModules = true
			}
		case *stats.ShowDesc:
			if v.Desc != "" {
				opt.StatsShowDesc = misc.StringP(v.Desc)
			}
		case *stats.MaxConn:
			d, err := v.Maxconn.Get(false)
			if err != nil {
				return nil
			}
			mc := d.(*types.Int64C)
			opt.StatsMaxconn = mc.Value
		case *stats.Refresh:
			if v.Delay != "" {
				opt.StatsRefreshDelay = misc.ParseTimeout(v.Delay)
			}
		case *stats.ShowNode:
			opt.StatsShowNodeName = misc.StringP(v.Name)
		case *stats.URI:
			if v.Prefix != "" {
				opt.StatsURIPrefix = v.Prefix
			}
		case *stats.Admin:
			if v != nil {
				opt.StatsAdmin = true
				if v.Cond != "" {
					opt.StatsAdminCond = v.Cond
					opt.StatsAdminCondTest = v.CondTest
				}
			}
		case *stats.Realm:
			if v != nil {
				opt.StatsRealm = true
				opt.StatsRealmRealm = misc.StringP(v.Realm)
			}
		case *stats.Auth:
			if v != nil {
				opt.StatsAuths = append(opt.StatsAuths, &models.StatsAuth{
					User:   misc.StringP(v.User),
					Passwd: misc.StringP(v.Password),
				})
			}
		case *stats.HTTPRequest:
			if v != nil && s.Section == parser.Backends {
				parts := strings.Split(v.Type, " ")
				httpRequest := &models.StatsHTTPRequest{
					Type: misc.StringP(parts[0]),
				}
				if len(parts) > 2 && parts[0] == "auth" && parts[1] == "realm" {
					httpRequest.Realm = strings.Join(parts[2:], " ")
				}
				if v.Cond != "" {
					httpRequest.Cond = v.Cond
					httpRequest.CondTest = v.CondTest
				}
				opt.StatsHTTPRequests = append(opt.StatsHTTPRequests, httpRequest)
			}
		}
	}
	return opt
}

func (s *SectionParser) monitorURI() interface{} {
	data, err := s.get("monitor-uri", false)
	if err != nil {
		return nil
	}
	d := data.(*types.MonitorURI)
	return models.MonitorURI(d.URI)
}

func (s *SectionParser) monitorFail() interface{} {
	if s.Section == parser.Frontends {
		data, err := s.get("monitor fail", false)
		if err != nil {
			return nil
		}
		d := data.(*types.MonitorFail)
		return &models.MonitorFail{
			Cond:     &d.Condition,
			CondTest: misc.StringP(strings.Join(d.ACLList, " ")),
		}
	}
	return nil
}

func (s *SectionParser) compression() interface{} { //nolint:gocognit
	compressionFound := false
	compression := &models.Compression{}

	data, err := s.get("compression algo", false)

	if err == nil {
		d, ok := data.(*types.StringSliceC)
		if ok && d != nil && len(d.Value) > 0 {
			compressionFound = true
			compression.Algorithms = d.Value
		}
	}

	data, err = s.get("compression algo-req", false)
	if err == nil {
		d, ok := data.(*types.StringC)
		if ok && d != nil {
			compressionFound = true
			compression.AlgoReq = d.Value
		}
	}

	data, err = s.get("compression algo-res", false)
	if err == nil {
		d, ok := data.(*types.StringSliceC)
		if ok && d != nil && len(d.Value) > 0 {
			compressionFound = true
			compression.AlgosRes = d.Value
		}
	}

	data, err = s.get("compression type", false)
	if err == nil {
		d, ok := data.(*types.StringSliceC)
		if ok && d != nil && len(d.Value) > 0 {
			compressionFound = true
			compression.Types = d.Value
		}
	}

	data, err = s.get("compression type-req", false)
	if err == nil {
		d, ok := data.(*types.StringSliceC)
		if ok && d != nil && len(d.Value) > 0 {
			compressionFound = true
			compression.TypesReq = d.Value
		}
	}

	data, err = s.get("compression type-res", false)
	if err == nil {
		d, ok := data.(*types.StringSliceC)
		if ok && d != nil && len(d.Value) > 0 {
			compressionFound = true
			compression.TypesRes = d.Value
		}
	}

	data, err = s.get("compression offload", false)
	if err == nil {
		d, ok := data.(*types.Enabled)
		if ok && d != nil {
			compressionFound = true
			compression.Offload = true
		}
	}

	data, err = s.get("compression direction", false)
	if err == nil {
		d, ok := data.(*types.StringC)
		if ok && d != nil {
			compressionFound = true
			compression.Direction = d.Value
		}
	}

	if compressionFound {
		return compression
	}
	return nil
}

func (s *SectionParser) clitcpkaIdle() interface{} {
	data, err := s.get("clitcpka-idle", false)
	if err != nil {
		return nil
	}
	d := data.(*types.StringC)
	return misc.ParseTimeoutDefaultSeconds(d.Value)
}

func (s *SectionParser) clitcpkaIntvl() interface{} {
	data, err := s.get("clitcpka-intvl", false)
	if err != nil {
		return nil
	}
	d := data.(*types.StringC)
	return misc.ParseTimeoutDefaultSeconds(d.Value)
}

func (s *SectionParser) srvtcpkaIdle() interface{} {
	data, err := s.get("srvtcpka-idle", false)
	if err != nil {
		return nil
	}
	d := data.(*types.StringC)
	return misc.ParseTimeoutDefaultSeconds(d.Value)
}

func (s *SectionParser) srvtcpkaIntvl() interface{} {
	data, err := s.get("srvtcpka-intvl", false)
	if err != nil {
		return nil
	}
	d := data.(*types.StringC)
	return misc.ParseTimeoutDefaultSeconds(d.Value)
}

func (s *SectionParser) serverStateFileName() interface{} {
	data, err := s.get("server-state-file-name", false)
	if err != nil {
		return nil
	}
	d := data.(*types.StringC)
	return d.Value
}

func (s *SectionParser) description() interface{} {
	data, err := s.get("description", false)
	if err != nil {
		return nil
	}
	d := data.(*types.StringC)
	return d.Value
}

func (s *SectionParser) errorloc302() interface{} {
	data, err := s.get("errorloc302", false)
	if err != nil {
		return nil
	}
	d := data.(*types.ErrorLoc302)
	if d == nil {
		return nil
	}
	intCode, err := strconv.ParseInt(d.Code, 10, 64)
	if err != nil {
		return nil
	}
	value := &models.Errorloc{
		Code: &intCode,
		URL:  &d.URL,
	}
	return value
}

func (s *SectionParser) errorloc303() interface{} {
	data, err := s.get("errorloc303", false)
	if err != nil {
		return nil
	}
	d := data.(*types.ErrorLoc303)
	if d == nil {
		return nil
	}
	intCode, err := strconv.ParseInt(d.Code, 10, 64)
	if err != nil {
		return nil
	}
	value := &models.Errorloc{
		Code: &intCode,
		URL:  &d.URL,
	}
	return value
}

func (s *SectionParser) httpRestirctReqHdrNames() interface{} {
	data, err := s.get("option http-restrict-req-hdr-names", false)
	if err != nil {
		return nil
	}
	d := data.(*types.OptionHTTPRestrictReqHdrNames)
	if d == nil {
		return nil
	}
	return d.Policy
}

func (s *SectionParser) defaultBind() interface{} {
	data, err := s.get("default-bind", false)
	if err != nil {
		return nil
	}

	d := data.(*types.DefaultBind)
	return &models.DefaultBind{
		BindParams: parseBindParams(d.Params),
	}
}

func (s *SectionParser) httpSendNameHeader() interface{} {
	if s.Section == parser.Defaults || s.Section == parser.Backends {
		data, err := s.get("http-send-name-header", false)
		if err != nil {
			return nil
		}
		d := data.(*types.HTTPSendNameHeader)
		if d == nil {
			return nil
		}
		return &d.Name
	}
	return nil
}

func (s *SectionParser) forcePersistList() interface{} {
	if s.Section != parser.Backends {
		return nil
	}
	data, err := s.get("force-persist", false)
	if err != nil {
		return nil
	}
	d := data.([]types.ForcePersist)
	if len(d) == 0 {
		return nil
	}

	items := make([]*models.ForcePersist, len(d))
	for i := range d {
		items[i] = &models.ForcePersist{
			Cond:     &d[i].Cond,
			CondTest: &d[i].CondTest,
		}
	}
	if len(d) == 0 {
		items = nil
	}
	return items
}

func (s *SectionParser) ignorePersistList() interface{} {
	if s.Section != parser.Backends {
		return nil
	}
	data, err := s.get("ignore-persist", false)
	if err != nil {
		return nil
	}
	d := data.([]types.IgnorePersist)
	if len(d) == 0 {
		return nil
	}

	items := make([]*models.IgnorePersist, len(d))
	for i := range d {
		items[i] = &models.IgnorePersist{
			Cond:     &d[i].Cond,
			CondTest: &d[i].CondTest,
		}
	}
	if len(d) == 0 {
		items = nil
	}
	return items
}

func (s *SectionParser) source() interface{} {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		data, err := s.get("source", false)
		if err != nil {
			return nil
		}
		d := data.(*types.Source)
		source := &models.Source{
			Address:       &d.Address,
			AddressSecond: d.AddressSecond,
			Hdr:           d.Hdr,
			Occ:           d.Occ,
			Interface:     d.Interface,
		}
		if d.Port != 0 {
			source.Port = d.Port
		}
		if d.PortSecond != 0 {
			source.PortSecond = d.PortSecond
		}
		switch {
		case d.Client:
			source.Usesrc = models.SourceUsesrcClient
		case d.ClientIP:
			source.Usesrc = models.SourceUsesrcClientip
		case d.HdrIP:
			source.Usesrc = models.SourceUsesrcHdrIP
		case d.AddressSecond != "":
			source.Usesrc = models.SourceUsesrcAddress
		}
		return source
	}
	return nil
}

func (s *SectionParser) shards() interface{} {
	if s.Section == parser.Peers {
		data, err := s.get("shards", false)
		if err != nil {
			return nil
		}

		d := data.(*types.Int64C)

		return d.Value
	}

	return nil
}

func (s *SectionParser) originalto() interface{} {
	data, err := s.get("option originalto", false)
	if err != nil {
		return nil
	}
	d := data.(*types.OptionOriginalTo)
	enabled := "enabled"
	originalto := &models.Originalto{
		Except:  d.Except,
		Header:  d.Header,
		Enabled: &enabled,
	}
	return originalto
}

func (s *SectionParser) logSteps() interface{} {
	if s.Section == parser.Frontends || s.Section == parser.Defaults {
		data, err := s.get("log-steps", false)
		if err != nil {
			return nil
		}
		d := data.(*types.StringC)
		return strings.Split(d.Value, ",")
	}
	return nil
}

// SectionObject represents a configuration section
type SectionObject struct {
	Object  interface{}
	Parser  parser.Parser
	Section parser.Section
	Name    string
	Options *options.ConfigurationOptions
}

// CreateEditSection creates or updates a section in the parser based on the provided object
func CreateEditSection(object interface{}, section parser.Section, pName string, p parser.Parser, opt *options.ConfigurationOptions) error {
	so := SectionObject{
		Object:  object,
		Section: section,
		Name:    pName,
		Parser:  p,
		Options: opt,
	}
	return so.CreateEditSection()
}

// CreateEditSection creates or updates a section in the parser based on the provided object
func (s *SectionObject) CreateEditSection() error {
	objValue := reflect.ValueOf(s.Object)
	if objValue.Kind() == reflect.Ptr {
		objValue = reflect.ValueOf(s.Object).Elem()
	}
	for i := range objValue.NumField() {
		typeField := objValue.Type().Field(i)
		field := objValue.FieldByName(typeField.Name)
		if typeField.Name != "Name" && typeField.Name != "ID" {
			if err := s.setFieldValue(typeField.Name, field); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *SectionObject) setFieldValue(fieldName string, field reflect.Value) error {
	if match := s.checkParams(fieldName); match {
		return nil
	}

	if match, err := s.checkSpecialFields(fieldName, field); match {
		return err
	}

	if match, err := s.checkTimeouts(fieldName, field); match {
		return err
	}

	if match, err := s.checkOptions(fieldName, field); match {
		return err
	}

	if match, err := s.checkSingleLine(fieldName, field); match {
		return err
	}

	return fmt.Errorf("cannot parse option for %s %s: %s", s.Section, s.Name, fieldName)
}

func (s *SectionObject) checkParams(fieldName string) bool {
	return s.Section != parser.FCGIApp && strings.HasSuffix(fieldName, "Params")
}

func (s *SectionObject) checkSpecialFields(fieldName string, field reflect.Value) (bool, error) { //nolint:gocyclo,cyclop
	switch fieldName {
	case "Shard":
		return true, s.shard(field)
	case "From":
		return true, s.from(field)
	case "MonitorURI":
		return true, s.monitorURI(field)
	case "MonitorFail":
		return true, s.monitorFail(field)
	case "StatsOptions":
		return true, s.statsOptions(field)
	case "Forwardfor":
		return true, s.forwardfor(field)
	case "Redispatch":
		return true, s.redispatch(field)
	case "Balance":
		return true, s.balance(field)
	case "PersistRule":
		return true, s.persistRule(field)
	case "Cookie":
		return true, s.cookie(field)
	case "HashType":
		return true, s.hashType(field)
	case "ErrorFiles":
		return true, s.errorFiles(field)
	case "ErrorFilesFromHTTPErrors":
		return true, s.errorFilesFromHTTPErrors(field)
	case "DefaultServer":
		return true, s.defaultServer(field)
	case "LoadServerStateFromFile":
		return true, s.loadServerStateFromFile(field)
	case "StickTable":
		return true, s.stickTable(field)
	case "AdvCheck":
		return true, s.advCheck(field)
	case "Logasap":
		return true, s.logasap(field)
	case "Allbackups":
		return true, s.allbackups(field)
	case "ExternalCheck":
		return true, s.externalCheck(field)
	case "ExternalCheckPath":
		return true, s.externalCheckPath(field)
	case "ExternalCheckCommand":
		return true, s.externalCheckCommand(field)
	case "DefaultBackend":
		return true, s.defaultBackend(field)
	case "HTTPConnectionMode":
		return true, s.httpConnectionMode(field)
	case "HTTPReuse":
		return true, s.httpReuse(field)
	case "UniqueIDFormat":
		return true, s.uniqueIDFormat(field)
	case "UniqueIDHeader":
		return true, s.uniqueIDHeader(field)
	case "Clflog":
		return true, s.clflog(field)
	case "Httplog":
		return true, s.httplog(field)
	case "Compression":
		return true, s.compression(field)
	case "ClitcpkaIdle":
		return true, s.clitcpkaIdle(field)
	case "ClitcpkaIntvl":
		return true, s.clitcpkaIntvl(field)
	case "SrvtcpkaIdle":
		return true, s.srvtcpkaIdle(field)
	case "SrvtcpkaIntvl":
		return true, s.srvtcpkaIntvl(field)
	case "EmailAlert":
		return true, s.emailAlert(field)
	case "ServerStateFileName":
		return true, s.serverStateFileName(field)
	case "Description":
		return true, s.description(field)
	case "Errorloc302":
		return true, s.errorloc302(field)
	case "Errorloc303":
		return true, s.errorloc303(field)
	case "HTTPRestrictReqHdrNames":
		return true, s.httpRestrictReqHdrNames(field)
	case "DefaultBind":
		return true, s.defaultBind(field)
	case "HTTPSendNameHeader":
		return true, s.httpSendNameHeader(field)
	case "ForcePersistList":
		return true, s.forcePersistList(field)
	case "ForcePersist":
		// "ForcePersist" field (force_persist) is deprecated in favour of "ForcePersistList" (force_persist_list).
		// Backward compatibility during the sunset period is handled by callers of this library that perform payload
		// transformation as necessary and remove the deprecated field.
		// "ForcePersist" is explicitly caught and ignored here as a safeguard against a runtime panic that can occur
		// if callers behave unexpectedly. It should be removed at the end of the sunset period along with the field.
		return true, nil
	case "IgnorePersistList":
		return true, s.ignorePersistList(field)
	case "IgnorePersist":
		// "IgnorePersist" field (ignore_persist) is deprecated in favour of "IgnorePersistList" (ignore_persist_list).
		// Backward compatibility during the sunset period is handled by callers of this library that perform payload
		// transformation as necessary and remove the deprecated field.
		// "IgnorePersist" is explicitly caught and ignored here as a safeguard against a runtime panic that can occur
		// if callers behave unexpectedly. It should be removed at the end of the sunset period along with the field.
		return true, nil
	case "Source":
		return true, s.source(field)
	case "Originalto":
		return true, s.originalto(field)
	case "LogSteps":
		return true, s.logSteps(field)
	default:
		return false, nil
	}
}

func (s *SectionObject) checkTimeouts(fieldName string, field reflect.Value) (bool, error) {
	if strings.HasSuffix(fieldName, "Timeout") {
		if pName := translateTimeout(fieldName); s.Parser.HasParser(s.Section, pName) {
			if valueIsNil(field) {
				if err := s.set(pName, nil); err != nil {
					return true, err
				}
				return true, nil
			}
			t := &types.SimpleTimeout{}
			t.Value = misc.SerializeTime(field.Elem().Int(), s.Options.PreferredTimeSuffix)
			if err := s.set(pName, t); err != nil {
				return true, err
			}
		}
		return true, nil
	}
	return false, nil
}

func (s *SectionObject) checkOptions(fieldName string, field reflect.Value) (bool, error) {
	if pName := "option " + misc.DashCase(fieldName); s.Parser.HasParser(s.Section, pName) {
		if valueIsNil(field) {
			if err := s.set(pName, nil); err != nil {
				return true, err
			}
			return true, nil
		}
		o := &types.SimpleOption{}
		if field.Kind() == reflect.String {
			if field.String() == "disabled" {
				o.NoOption = true
			}
		}
		if err := s.set(pName, o); err != nil {
			return true, err
		}
		return true, nil
	}
	return false, nil
}

func (s *SectionObject) checkSingleLine(fieldName string, field reflect.Value) (bool, error) {
	if pName := misc.DashCase(fieldName); s.Parser.HasParser(s.Section, pName) {
		if valueIsNil(field) {
			if err := s.set(pName, nil); err != nil {
				return true, err
			}
			return true, nil
		}
		d := translateToParserData(field)
		if d == nil {
			return true, fmt.Errorf("cannot parse type for %s %s: %s", s.Section, s.Name, fieldName)
		}
		if err := s.set(pName, d); err != nil {
			return true, err
		}
		return true, nil
	}
	return false, nil
}

func (s *SectionObject) set(attribute string, data interface{}) error {
	return s.Parser.Set(s.Section, s.Name, attribute, data)
}

func (s *SectionObject) from(field reflect.Value) error {
	if s.Section == parser.Frontends || s.Section == parser.Backends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			return s.Parser.SectionsDefaultsFromSet(s.Section, s.Name, "")
		}
		return s.Parser.SectionsDefaultsFromSet(s.Section, s.Name, field.String())
	}
	return nil
}

func (s *SectionObject) httplog(field reflect.Value) error {
	if s.Section == parser.Frontends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			// check if clflog is active, if yes, do nothing
			d, err := s.Parser.Get(s.Section, s.Name, "option httplog", false)
			if err != nil {
				if !errors.Is(err, parser_errors.ErrFetch) {
					return err
				}
				return nil
			}
			o := d.(*types.OptionHTTPLog)
			if !o.Clf {
				if err := s.set("option httplog", nil); err != nil {
					return err
				}
			}
			return nil
		}
		o := &types.OptionHTTPLog{}
		if err := s.set("option httplog", o); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) clflog(field reflect.Value) error {
	if s.Section == parser.Frontends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			// check if httplog exists, if not do nothing
			d, err := s.Parser.Get(s.Section, s.Name, "option httplog", false)
			if err != nil {
				if !errors.Is(err, parser_errors.ErrFetch) {
					return err
				}
				return nil
			}
			o := d.(*types.OptionHTTPLog)
			if o.Clf {
				o.Clf = false
				if err := s.set("option httplog", o); err != nil {
					return err
				}
			}
			return nil
		}
		o := &types.OptionHTTPLog{Clf: true}
		if err := s.set("option httplog", o); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) uniqueIDHeader(field reflect.Value) error {
	if s.Section != parser.Defaults && s.Section != parser.Frontends {
		return nil
	}
	if valueIsNil(field) {
		return s.set("unique-id-header", nil)
	}
	d := types.UniqueIDHeader{
		Name: field.String(),
	}
	return s.set("unique-id-header", &d)
}

func (s *SectionObject) uniqueIDFormat(field reflect.Value) error {
	if s.Section != parser.Defaults && s.Section != parser.Frontends {
		return nil
	}
	if valueIsNil(field) {
		return s.set("unique-id-format", nil)
	}
	d := types.UniqueIDFormat{
		LogFormat: field.String(),
	}
	return s.set("unique-id-format", &d)
}

func (s *SectionObject) httpReuse(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			return s.set("http-reuse", nil)
		}

		b := field.String()
		d := types.HTTPReuse{
			ShareType: b,
		}

		if err := s.set("http-reuse", &d); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) httpConnectionMode(field reflect.Value) error {
	for _, opt := range []string{"httpclose", "http-server-close", "http-keep-alive"} {
		attribute := "option " + opt

		if err := s.set(attribute, nil); err != nil {
			return err
		}
	}
	// Deprecated, delete if exists
	_ = s.set("option forceclose", nil)

	if !valueIsNil(field) {
		pName := fmt.Sprintf("option %v", field.String())
		d := &types.SimpleOption{
			NoOption: false,
		}
		if err := s.set(pName, d); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) defaultBackend(field reflect.Value) error {
	if s.Section == parser.Frontends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			return s.set("default_backend", nil)
		}
		bck := field.String()
		d := &types.StringC{
			Value: bck,
		}
		if err := s.set("default_backend", d); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) externalCheckCommand(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		pExtCmd := &types.ExternalCheckCommand{}
		if valueIsNil(field) {
			pExtCmd = nil
		} else {
			str, ok := field.Interface().(string)
			if !ok {
				return misc.CreateTypeAssertError("external-check command")
			}
			pExtCmd.Command = str
		}
		if err := s.set("external-check command", pExtCmd); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) externalCheckPath(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		pExtPath := &types.ExternalCheckPath{}
		if valueIsNil(field) {
			pExtPath = nil
		} else {
			str, ok := field.Interface().(string)
			if !ok {
				return misc.CreateTypeAssertError("external-check path")
			}
			pExtPath.Path = str
		}
		if err := s.set("external-check path", pExtPath); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) externalCheck(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		pExternalCheck := &types.SimpleOption{}
		if valueIsNil(field) {
			pExternalCheck = nil
		} else if field.String() == "disabled" {
			pExternalCheck.NoOption = true
		}
		if err := s.set("option external-check", pExternalCheck); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) allbackups(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		allbackups := &types.SimpleOption{}
		if valueIsNil(field) {
			allbackups = nil
		} else if field.String() == "disabled" {
			allbackups.NoOption = true
		}
		if err := s.set("option allbackups", allbackups); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) logasap(field reflect.Value) error {
	if s.Section == parser.Frontends || s.Section == parser.Defaults {
		logasap := &types.SimpleOption{}
		if valueIsNil(field) {
			logasap = nil
		} else if field.String() == "disabled" {
			logasap.NoOption = true
		}
		if err := s.set("option logasap", logasap); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) advCheck(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		if err := s.resetCheckOptions(); err != nil {
			return err
		}

		if !valueIsNil(field) {
			pName := fmt.Sprintf("option %v", field.String())
			d, err := s.getCheckData(pName)
			if err != nil {
				return err
			}
			if err := s.set(pName, d); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *SectionObject) resetCheckOptions() error {
	if err := s.set("option ssl-hello-chk", nil); err != nil {
		return err
	}
	if err := s.set("option smtpchk", nil); err != nil {
		return err
	}
	if err := s.set("option ldap-check", nil); err != nil {
		return err
	}
	if err := s.set("option mysql-check", nil); err != nil {
		return err
	}
	if err := s.set("option pgsql-check", nil); err != nil {
		return err
	}
	if err := s.set("option tcp-check", nil); err != nil {
		return err
	}
	if err := s.set("option redis-check", nil); err != nil {
		return err
	}
	return s.set("option httpchk", nil)
}

func (s *SectionObject) getCheckData(pName string) (common.ParserData, error) {
	switch pName {
	case "option smtpchk":
		return s.getSmtpchkData()
	case "option mysql-check":
		return s.getMysqlCheckData()
	case "option pgsql-check":
		return s.getPgsqlCheckData()
	case "option httpchk":
		return s.getHTTPChkData()
	default:
		return &types.SimpleOption{
			NoOption: false,
		}, nil
	}
}

func (s *SectionObject) getSmtpchkData() (common.ParserData, error) {
	data := s.getFieldByName("SmtpchkParams")
	if data == nil {
		return &types.OptionSmtpchk{
			NoOption: false,
		}, nil
	}
	params, ok := data.(models.SmtpchkParams)
	if !ok {
		return nil, misc.CreateTypeAssertError("SmtpchkParams")
	}
	return &types.OptionSmtpchk{
		NoOption: false,
		Hello:    params.Hello,
		Domain:   params.Domain,
	}, nil
}

func (s *SectionObject) getMysqlCheckData() (common.ParserData, error) {
	data := s.getFieldByName("MysqlCheckParams")
	if data == nil {
		return &types.OptionMysqlCheck{
			NoOption: false,
		}, nil
	}
	params := data.(models.MysqlCheckParams)
	return &types.OptionMysqlCheck{
		NoOption:      false,
		ClientVersion: params.ClientVersion,
		User:          params.Username,
	}, nil
}

func (s *SectionObject) getPgsqlCheckData() (common.ParserData, error) {
	data := s.getFieldByName("PgsqlCheckParams")
	if data == nil {
		return errors.New("adv_check value pgsql-check requires pgsql_check_params"), nil
	}
	params, ok := data.(models.PgsqlCheckParams)
	if !ok {
		return nil, misc.CreateTypeAssertError("adv_check value pgsql-check requires pgsql_check_params")
	}
	if params.Username == "" {
		return errors.New("adv_check value pgsql-check requires username in pgsql_check_params"), nil
	}
	return &types.OptionPgsqlCheck{
		NoOption: false,
		User:     params.Username,
	}, nil
}

func (s *SectionObject) getHTTPChkData() (common.ParserData, error) {
	data := s.getFieldByName("HttpchkParams")
	if data == nil {
		return &types.OptionHttpchk{
			NoOption: false,
		}, nil
	}
	params, ok := data.(models.HttpchkParams)
	if !ok {
		return nil, misc.CreateTypeAssertError("HttpchkParams")
	}
	return &types.OptionHttpchk{
		NoOption: false,
		Method:   params.Method,
		Version:  params.Version,
		URI:      params.URI,
		Host:     params.Host,
	}, nil
}

func (s *SectionObject) getFieldByName(fieldName string) interface{} {
	objValue := reflect.ValueOf(s.Object).Elem()
	elem := objValue.FieldByName(fieldName)
	if elem.IsNil() {
		return nil
	}
	return elem.Elem().Interface()
}

func (s *SectionObject) stickTable(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Frontends || s.Section == parser.Peers {
		if valueIsNil(field) {
			return s.set("stick-table", nil)
		}
		st, ok := field.Elem().Interface().(models.ConfigStickTable)
		if !ok {
			return misc.CreateTypeAssertError("stick-table")
		}
		d := types.StickTable{
			Type:    st.Type,
			Store:   st.Store,
			Peers:   st.Peers,
			NoPurge: st.Nopurge,
		}

		if st.Keylen != nil {
			d.Length = strconv.FormatInt(*st.Keylen, 10)
		}
		if st.Expire != nil {
			d.Expire = misc.SerializeTime(*st.Expire, s.Options.PreferredTimeSuffix)
		}
		if st.Size != nil {
			d.Size = misc.SerializeSize(*st.Size)
		}
		if st.Srvkey != nil {
			d.SrvKey = *st.Srvkey
		}
		if st.WriteTo != nil {
			d.WriteTo = *st.WriteTo
		}
		if err := s.set("stick-table", d); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) defaultServer(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults || s.Section == parser.Peers {
		if valueIsNil(field) {
			return s.set("default-server", nil)
		}
		ds, ok := field.Elem().Interface().(models.DefaultServer)
		if !ok {
			return misc.CreateTypeAssertError("default-server")
		}
		dServers := []types.DefaultServer{{}}
		dServers[0].Params = SerializeServerParams(ds.ServerParams, s.Options)
		if err := s.set("default-server", dServers); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) loadServerStateFromFile(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			return s.set("load-server-state-from-file", nil)
		}

		b := field.String()
		d := types.LoadServerStateFromFile{
			Argument: b,
		}

		if err := s.set("load-server-state-from-file", &d); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) errorFiles(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("errorfile", nil)
	}
	efs, ok := field.Interface().([]*models.Errorfile)
	if !ok {
		return nil
	}
	errorFiles := []types.ErrorFile{}
	for _, ef := range efs {
		errorFiles = append(errorFiles, types.ErrorFile{Code: strconv.FormatInt(ef.Code, 10), File: ef.File})
	}
	return s.set("errorfile", errorFiles)
}

func (s *SectionObject) errorFilesFromHTTPErrors(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("errorfiles", nil)
	}
	efs, ok := field.Interface().([]*models.Errorfiles)
	if !ok {
		return nil
	}
	errorFiles := make([]types.ErrorFiles, len(efs))
	for i, ef := range efs {
		errorFiles[i] = types.ErrorFiles{Codes: ef.Codes, Name: ef.Name}
	}
	if len(efs) == 0 {
		errorFiles = nil
	}
	return s.set("errorfiles", errorFiles)
}

func (s *SectionObject) hashType(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			return s.set("hash-type", nil)
		}
		b, ok := field.Elem().Interface().(models.HashType)
		if !ok {
			return misc.CreateTypeAssertError("hash-type")
		}
		d := types.HashType{
			Method:   b.Method,
			Function: b.Function,
			Modifier: b.Modifier,
		}
		if err := s.set("hash-type", &d); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) cookie(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			return s.set("cookie", nil)
		}
		d, ok := field.Elem().Interface().(models.Cookie)
		if !ok {
			return misc.CreateTypeAssertError("cookie")
		}
		domains := make([]string, len(d.Domains))
		for i, domain := range d.Domains {
			domains[i] = domain.Value
		}
		if len(d.Domains) == 0 {
			domains = nil
		}
		attrs := make([]string, len(d.Attrs))
		for i, attr := range d.Attrs {
			attrs[i] = attr.Value
		}
		if len(d.Attrs) == 0 {
			attrs = nil
		}
		data := types.Cookie{
			Attr:     attrs,
			Domain:   domains,
			Dynamic:  d.Dynamic,
			Httponly: d.Httponly,
			Indirect: d.Indirect,
			Maxidle:  d.Maxidle,
			Maxlife:  d.Maxlife,
			Name:     *d.Name,
			Nocache:  d.Nocache,
			Postonly: d.Postonly,
			Preserve: d.Preserve,
			Type:     d.Type,
			Secure:   d.Secure,
		}
		if err := s.set("cookie", &data); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) persistRule(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			return s.set("persist", nil)
		}
		b, ok := field.Elem().Interface().(models.PersistRule)
		if !ok {
			return misc.CreateTypeAssertError("persist")
		}
		d := types.Persist{
			Type: *b.Type,
		}
		switch *b.Type {
		case "rdp-cookie":
			d.Params = &params.PersistRdpCookie{
				Name: b.RdpCookieName,
			}
		default:
		}
		if err := s.set("persist", &d); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) balance(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			return s.set("balance", nil)
		}
		b, ok := field.Elem().Interface().(models.Balance)
		if !ok {
			return misc.CreateTypeAssertError("balance")
		}
		d := types.Balance{
			Algorithm: *b.Algorithm,
		}

		switch *b.Algorithm {
		case "uri":
			d.Params = &params.BalanceURI{
				Depth:    b.URIDepth,
				Len:      b.URILen,
				Whole:    b.URIWhole,
				PathOnly: b.URIPathOnly,
			}
		case "url_param":
			d.Params = &params.BalanceURLParam{
				Param:     b.URLParam,
				CheckPost: b.URLParamCheckPost,
				MaxWait:   b.URLParamMaxWait,
			}
		case "hdr":
			d.Params = &params.BalanceHdr{
				Name:          b.HdrName,
				UseDomainOnly: b.HdrUseDomainOnly,
			}
		case "random":
			d.Params = &params.BalanceRandom{
				Draws: b.RandomDraws,
			}
		case "rdp-cookie":
			d.Params = &params.BalanceRdpCookie{
				Name: b.RdpCookieName,
			}
		case "hash":
			d.Params = &params.BalanceHash{
				Expression: b.HashExpression,
			}
		}
		if err := s.set("balance", &d); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) redispatch(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			return s.set("option redispatch", nil)
		}
		br, ok := field.Elem().Interface().(models.Redispatch)
		if !ok {
			return misc.CreateTypeAssertError("option redispatch")
		}
		d := &types.OptionRedispatch{
			Interval: br.Interval,
			NoOption: false,
		}
		if *br.Enabled == "disabled" {
			d.NoOption = true
		}
		if err := s.set("option redispatch", d); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) forwardfor(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("option forwardfor", nil)
	}
	ff, ok := field.Elem().Interface().(models.Forwardfor)
	if !ok {
		return misc.CreateTypeAssertError("option forwardfor")
	}
	d := &types.OptionForwardFor{
		Except: ff.Except,
		Header: ff.Header,
		IfNone: ff.Ifnone,
	}
	return s.set("option forwardfor", d)
}

func (s *SectionObject) monitorURI(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("monitor-uri", nil)
	}
	v := field.String()
	return s.set("monitor-uri", types.MonitorURI{URI: v})
}

func (s *SectionObject) monitorFail(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("monitor fail", nil)
	}
	opt, ok := field.Elem().Interface().(models.MonitorFail)
	if !ok {
		return misc.CreateTypeAssertError("monitor fail")
	}
	return s.set("monitor fail", types.MonitorFail{
		Condition: *opt.Cond,
		ACLList:   strings.Split(*opt.CondTest, " "),
	})
}

func (s *SectionObject) emailAlert(field reflect.Value) error {
	if valueIsNil(field) {
		return nil
	}
	ea, ok := field.Elem().Interface().(models.EmailAlert)
	if !ok {
		return misc.CreateTypeAssertError("email-alert")
	}

	list := []types.EmailAlert{}

	if ea.From != nil {
		e := types.EmailAlert{
			Attribute: "from",
			Value:     *ea.From,
		}
		list = append(list, e)
	}
	if ea.To != nil {
		e := types.EmailAlert{
			Attribute: "to",
			Value:     *ea.To,
		}
		list = append(list, e)
	}
	if ea.Level != "" {
		e := types.EmailAlert{
			Attribute: "level",
			Value:     ea.Level,
		}
		list = append(list, e)
	}
	if ea.Myhostname != "" {
		e := types.EmailAlert{
			Attribute: "myhostname",
			Value:     ea.Myhostname,
		}
		list = append(list, e)
	}
	if ea.Mailers != nil {
		e := types.EmailAlert{
			Attribute: "mailers",
			Value:     *ea.Mailers,
		}
		list = append(list, e)
	}

	return s.set("email-alert", list)
}

func (s *SectionObject) statsOptions(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("stats", nil)
	}
	opt, ok := field.Elem().Interface().(models.StatsOptions)
	if !ok {
		return misc.CreateTypeAssertError("stats options")
	}
	ss := []types.StatsSettings{}

	if opt.StatsEnable {
		s := &stats.OneWord{
			Name: "enable",
		}
		ss = append(ss, s)
	}
	if opt.StatsHideVersion {
		s := &stats.OneWord{
			Name: "hide-version",
		}
		ss = append(ss, s)
	}
	if opt.StatsShowLegends {
		s := &stats.OneWord{
			Name: "show-legends",
		}
		ss = append(ss, s)
	}
	if opt.StatsShowDesc != nil {
		s := &stats.ShowDesc{
			Desc: *opt.StatsShowDesc,
		}
		ss = append(ss, s)
	}
	if opt.StatsRefreshDelay != nil {
		s := &stats.Refresh{
			Delay: strconv.FormatInt(*opt.StatsRefreshDelay, 10),
		}
		ss = append(ss, s)
	}
	if opt.StatsShowNodeName != nil {
		s := &stats.ShowNode{
			Name: *opt.StatsShowNodeName,
		}
		ss = append(ss, s)
	}
	if opt.StatsURIPrefix != "" {
		s := &stats.URI{
			Prefix: opt.StatsURIPrefix,
		}
		ss = append(ss, s)
	}
	if opt.StatsMaxconn > 0 {
		d := &types.Int64C{
			Value: opt.StatsMaxconn,
		}
		s := &stats.MaxConn{}
		s.Maxconn = &parsers.MaxConn{}
		if err := s.Maxconn.Set(d, 0); err != nil {
			return err
		}
		ss = append(ss, s)
	}
	if opt.StatsAdmin {
		s := &stats.Admin{
			Cond:     opt.StatsAdminCond,
			CondTest: opt.StatsAdminCondTest,
		}
		ss = append(ss, s)
	}
	if opt.StatsShowModules {
		s := &stats.OneWord{
			Name: "show-modules",
		}
		ss = append(ss, s)
	}
	if opt.StatsRealm {
		s := &stats.Realm{
			Realm: *opt.StatsRealmRealm,
		}
		ss = append(ss, s)
	}

	for _, auth := range opt.StatsAuths {
		s := &stats.Auth{
			User:     *auth.User,
			Password: *auth.Passwd,
		}
		ss = append(ss, s)
	}
	for _, httpRequest := range opt.StatsHTTPRequests {
		reqType := *httpRequest.Type
		if reqType == "auth" && httpRequest.Realm != "" {
			reqType = "auth realm " + httpRequest.Realm
		}
		s := &stats.HTTPRequest{
			Type:     reqType,
			Cond:     httpRequest.Cond,
			CondTest: httpRequest.CondTest,
		}
		ss = append(ss, s)
	}
	return s.set("stats", ss)
}

func (s *SectionObject) compression(field reflect.Value) error { //nolint:gocognit
	var err error
	if valueIsNil(field) {
		err = s.set("compression algo", nil)
		if err != nil {
			return err
		}
		err = s.set("compression algo-req", nil)
		if err != nil {
			return err
		}
		err = s.set("compression algo-res", nil)
		if err != nil {
			return err
		}
		err = s.set("compression type", nil)
		if err != nil {
			return err
		}
		err = s.set("compression type-req", nil)
		if err != nil {
			return err
		}
		err = s.set("compression type-res", nil)
		if err != nil {
			return err
		}
		err = s.set("compression offload", nil)
		if err != nil {
			return err
		}

		err = s.set("compression direction", nil)
		if err != nil {
			// compression direction does not exist on Frontends
			if errors.Is(err, parser_errors.ErrAttributeNotFound) {
				return nil
			}
			return err
		}
		return nil
	}
	compression, ok := field.Elem().Interface().(models.Compression)
	if !ok {
		return errors.New("error casting compression model")
	}

	if len(compression.Algorithms) > 0 {
		err = s.set("compression algo", &types.StringSliceC{Value: compression.Algorithms})
		if err != nil {
			return err
		}
	}
	if len(compression.AlgoReq) > 0 {
		err = s.set("compression algo-req", &types.StringC{Value: compression.AlgoReq})
		if err != nil {
			return err
		}
	}
	if len(compression.AlgosRes) > 0 {
		err = s.set("compression algo-res", &types.StringSliceC{Value: compression.AlgosRes})
		if err != nil {
			return err
		}
	}
	if len(compression.Types) > 0 {
		err = s.set("compression type", &types.StringSliceC{Value: compression.Types})
		if err != nil {
			return err
		}
	}
	if len(compression.TypesReq) > 0 {
		err = s.set("compression type-req", &types.StringSliceC{Value: compression.TypesReq})
		if err != nil {
			return err
		}
	}
	if len(compression.TypesRes) > 0 {
		err = s.set("compression type-res", &types.StringSliceC{Value: compression.TypesRes})
		if err != nil {
			return err
		}
	}
	if compression.Offload {
		err = s.set("compression offload", &types.Enabled{})
		if err != nil {
			return err
		}
	}
	if len(compression.Direction) > 0 {
		err = s.set("compression direction", &types.StringC{Value: compression.Direction})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) clitcpkaIdle(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("clitcpka-idle", nil)
	}
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}
	v := field.Int()
	str := misc.SerializeTime(v, s.Options.PreferredTimeSuffix)
	return s.set("clitcpka-idle", types.StringC{Value: str})
}

func (s *SectionObject) clitcpkaIntvl(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("clitcpka-intvl", nil)
	}
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}
	v := field.Int()
	str := misc.SerializeTime(v, s.Options.PreferredTimeSuffix)
	return s.set("clitcpka-intvl", types.StringC{Value: str})
}

func (s *SectionObject) srvtcpkaIdle(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("srvtcpka-idle", nil)
	}
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}
	v := field.Int()
	str := misc.SerializeTime(v, s.Options.PreferredTimeSuffix)
	return s.set("srvtcpka-idle", types.StringC{Value: str})
}

func (s *SectionObject) srvtcpkaIntvl(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("srvtcpka-intvl", nil)
	}
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}
	v := field.Int()
	str := misc.SerializeTime(v, s.Options.PreferredTimeSuffix)
	return s.set("srvtcpka-intvl", types.StringC{Value: str})
}

func (s *SectionObject) serverStateFileName(field reflect.Value) error {
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}
	v := field.String()
	if v == "" {
		return nil
	}
	return s.set("server-state-file-name", types.StringC{Value: v})
}

func (s *SectionObject) description(field reflect.Value) error {
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}
	v := field.String()
	if v == "" {
		return nil
	}
	return s.set("description", types.StringC{Value: v})
}

func (s *SectionObject) errorloc302(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("errorloc302", nil)
	}
	errorLoc, ok := field.Elem().Interface().(models.Errorloc)
	if !ok {
		return misc.CreateTypeAssertError("errorloc302")
	}

	e := &types.ErrorLoc302{
		Code: strconv.FormatInt(*errorLoc.Code, 10),
		URL:  *errorLoc.URL,
	}
	return s.set("errorloc302", e)
}

func (s *SectionObject) errorloc303(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("errorloc303", nil)
	}
	errorLoc, ok := field.Elem().Interface().(models.Errorloc)
	if !ok {
		return misc.CreateTypeAssertError("errorloc303")
	}

	e := &types.ErrorLoc303{
		Code: strconv.FormatInt(*errorLoc.Code, 10),
		URL:  *errorLoc.URL,
	}

	return s.set("errorloc303", e)
}

func (s *SectionObject) httpRestrictReqHdrNames(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("option http-restrict-req-hdr-names", nil)
	}
	t := &types.OptionHTTPRestrictReqHdrNames{Policy: field.String()}
	return s.set("option http-restrict-req-hdr-names", t)
}

func (s *SectionObject) defaultBind(field reflect.Value) error {
	if s.Section != parser.Peers {
		return nil
	}
	if valueIsNil(field) {
		return s.set("default-bind", nil)
	}
	db, ok := field.Elem().Interface().(models.DefaultBind)
	if !ok {
		return misc.CreateTypeAssertError("default-bind")
	}
	dBind := &types.DefaultBind{
		Params: serializeBindParams(db.BindParams, ""),
	}

	return s.set("default-bind", dBind)
}

func (s *SectionObject) httpSendNameHeader(field reflect.Value) error {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return s.set("http-send-name-header", nil)
		}
		field = field.Elem()
	}
	v := field.String()
	return s.set("http-send-name-header", types.HTTPSendNameHeader{Name: v})
}

func (s *SectionObject) forcePersistList(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("force-persist", nil)
	}
	data, ok := field.Interface().([]*models.ForcePersist)
	if !ok {
		return misc.CreateTypeAssertError("force-persist")
	}

	items := make([]types.ForcePersist, len(data))
	for i := range data {
		items[i] = types.ForcePersist{
			Cond:     *data[i].Cond,
			CondTest: *data[i].CondTest,
		}
	}
	if len(data) == 0 {
		items = nil
	}
	return s.set("force-persist", items)
}

func (s *SectionObject) ignorePersistList(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("ignore-persist", nil)
	}
	data, ok := field.Interface().([]*models.IgnorePersist)
	if !ok {
		return misc.CreateTypeAssertError("ignore-persist")
	}

	items := make([]types.IgnorePersist, len(data))
	for i := range data {
		items[i] = types.IgnorePersist{
			Cond:     *data[i].Cond,
			CondTest: *data[i].CondTest,
		}
	}
	if len(data) == 0 {
		items = nil
	}
	return s.set("ignore-persist", items)
}

func (s *SectionObject) source(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("source", nil)
	}
	so, ok := field.Elem().Interface().(models.Source)
	if !ok {
		return misc.CreateTypeAssertError("source")
	}
	source := types.Source{
		Address:       *so.Address,
		Port:          so.Port,
		AddressSecond: so.AddressSecond,
		PortSecond:    so.PortSecond,
		Hdr:           so.Hdr,
		Occ:           so.Occ,
		Interface:     so.Interface,
	}
	switch so.Usesrc {
	case models.SourceUsesrcClient:
		source.Client = true
	case models.SourceUsesrcClientip:
		source.ClientIP = true
	case models.SourceUsesrcHdrIP:
		source.HdrIP = true
	}
	// source parser serialization expects UseSrc flag to be set in order to include related arguments.
	if source.Client || source.ClientIP || source.HdrIP || source.AddressSecond != "" {
		source.UseSrc = true
	}
	return s.set("source", source)
}

func (s *SectionObject) shard(field reflect.Value) error {
	if s.Section == parser.Peers {
		if valueIsNil(field) {
			return s.set("shards", nil)
		}
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		v := field.Int()
		return s.set("shards", types.Int64C{Value: v})
	}
	return nil
}

func (s *SectionObject) originalto(field reflect.Value) error {
	if !(s.Section == parser.Defaults || s.Section == parser.Frontends || s.Section == parser.Backends) {
		return nil
	}
	if valueIsNil(field) {
		return s.set("option originalto", nil)
	}
	originalto, ok := field.Elem().Interface().(models.Originalto)
	if !ok {
		return misc.CreateTypeAssertError("option originalto")
	}
	d := &types.OptionOriginalTo{
		Except: originalto.Except,
		Header: originalto.Header,
	}
	return s.set("option originalto", d)
}

func (s *SectionObject) logSteps(field reflect.Value) error {
	if !(s.Section == parser.Defaults || s.Section == parser.Frontends) {
		return nil
	}
	if valueIsNil(field) {
		return s.set("log-steps", nil)
	}
	logSteps, ok := field.Interface().([]string)
	if !ok {
		return misc.CreateTypeAssertError("log-steps")
	}
	d := strings.Join(logSteps, ",")
	if len(d) == 0 {
		return s.set("log-steps", nil)
	}
	return s.set("log-steps", d)
}

func (c *client) deleteSection(section parser.Section, name string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !c.checkSectionExists(section, name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", section, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err := p.SectionsDelete(section, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) editSection(section parser.Section, name string, data interface{}, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !c.checkSectionExists(section, name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", section, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err := CreateEditSection(data, section, name, p, &c.ConfigurationOptions); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) createSection(section parser.Section, name string, data interface{}, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if c.checkSectionExists(section, name, p) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", section, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err := p.SectionsCreate(section, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err := CreateEditSection(data, section, name, p, &c.ConfigurationOptions); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) checkSectionExists(section parser.Section, sectionName string, p parser.Parser) bool {
	sections, err := p.SectionsGet(section)
	if err != nil {
		return false
	}

	if misc.StringInSlice(sectionName, sections) {
		return true
	}
	return false
}

func (c *client) loadDataForChange(transactionID string, version int64) (parser.Parser, string, error) {
	t, err := c.TransactionClient.CheckTransactionOrVersion(transactionID, version)
	if err != nil {
		// if transactionID is implicit, return err and delete transaction
		if transactionID == "" && t != "" {
			return nil, "", c.ErrAndDeleteTransaction(err, t)
		}
		return nil, "", err
	}

	p, err := c.GetParser(t)
	if err != nil {
		if transactionID == "" && t != "" {
			return nil, "", c.ErrAndDeleteTransaction(err, t)
		}
		return nil, "", err
	}
	return p, t, nil
}

func valueIsNil(v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive,nolintlint
	case reflect.Int64:
		return v.Int() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr:
		return !v.Elem().IsValid()
	default:
		return false
	}
}

func translateToParserData(field reflect.Value) common.ParserData {
	switch field.Kind() { //nolint:exhaustive
	case reflect.Int64:
		return types.Int64C{Value: field.Int()}
	case reflect.String:
		return types.StringC{Value: field.String()}
	case reflect.Ptr:
		return types.Int64C{Value: field.Elem().Int()}
	case reflect.Bool:
		return types.Enabled{}
	default:
		return nil
	}
}

func parseOption(d interface{}) interface{} {
	switch v := d.(type) {
	case *types.StringC:
		return v.Value
	case *types.Int64C:
		return &v.Value
	case *types.Enabled:
		return "enabled"
	case *types.SimpleOption:
		if v.NoOption {
			return "disabled"
		}
		return "enabled"
	default:
		return nil
	}
}

func translateTimeout(mName string) string {
	mName = strings.TrimSuffix(mName, "Timeout")
	return "timeout " + misc.DashCase(mName)
}
