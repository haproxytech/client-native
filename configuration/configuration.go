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
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	parser "github.com/haproxytech/config-parser/v4"
	"github.com/haproxytech/config-parser/v4/common"
	parser_errors "github.com/haproxytech/config-parser/v4/errors"
	parser_options "github.com/haproxytech/config-parser/v4/options"
	"github.com/haproxytech/config-parser/v4/params"
	"github.com/haproxytech/config-parser/v4/parsers"
	stats "github.com/haproxytech/config-parser/v4/parsers/stats/settings"
	"github.com/haproxytech/config-parser/v4/types"
	"github.com/pkg/errors"

	"github.com/haproxytech/client-native/v4/misc"
	"github.com/haproxytech/client-native/v4/models"
)

// ClientParams is just a placeholder for all client options
type ClientParams struct {
	ConfigurationFile         string
	Haproxy                   string
	TransactionDir            string
	BackupsNumber             int
	UseValidation             bool
	PersistentTransactions    bool
	ValidateConfigurationFile bool
	MasterWorker              bool
	SkipFailedTransactions    bool
	UseMd5Hash                bool

	// ValidateCmd allows specifying a custom script to validate the transaction file.
	// The injected environment variable DATAPLANEAPI_TRANSACTION_FILE must be used to get the location of the file.
	ValidateCmd               string
	ValidateConfigFilesBefore []string
	ValidateConfigFilesAfter  []string
}

// Client configuration client
// Parser is the config parser instance that loads "master" configuration file on Init
// and when transaction is committed it gets replaced with the parser from parsers map.
// parsers map contains a config parser for each transaction, which loads data from
// transaction files on StartTransaction, and deletes on CommitTransaction. We save
// data to file on every change for persistence.
type client struct {
	Transaction
	parsers  map[string]parser.Parser
	services map[string]*Service
	parser   parser.Parser
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
		return nil, NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("Transaction %s does not exist", transactionID))
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
		return NewConfError(ErrTransactionAlreadyExists, fmt.Sprintf("Transaction %s already exists", transactionID))
	}

	parserOptions := []parser_options.ParserOption{}
	parserOptions = append(parserOptions, parser_options.NoNamedDefaultsFrom)
	if c.ConfigurationOptions.UseMd5Hash {
		parserOptions = append(parserOptions, parser_options.UseMd5Hash)
	}

	tFile := ""
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
		return NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read %s", tFile))
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
		return NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("Transaction %s does not exist", transactionID))
	}
	delete(c.parsers, transactionID)
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
		return NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("Transaction %s does not exist", transactionID))
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
		return 0, NewConfError(ErrCannotReadVersion, fmt.Sprintf("Cannot read version: %s", err.Error()))
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
		return NewConfError(ErrCannotSetVersion, fmt.Sprintf("Cannot set version: %s", err.Error()))
	}
	return nil
}

func (c *client) LoadData(filename string) error {
	err := c.parser.LoadData(filename)
	if err != nil {
		return NewConfError(ErrCannotReadConfFile, fmt.Sprintf("cannot read %s", filename))
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

// SectionParser is used set fields of a section based on the provided parser
type SectionParser struct {
	Object  interface{}
	Section parser.Section
	Name    string
	Parser  parser.Parser
}

// Parse parses the sections fields and sets their values with the data from the parser
func (s *SectionParser) Parse() error {
	objValue := reflect.ValueOf(s.Object).Elem()
	for i := 0; i < objValue.NumField(); i++ {
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

	if match, data := s.checkSingleLine(fieldName); match {
		return data
	}

	if match, data := s.checkOptions(fieldName); match {
		return data
	}

	return nil
}

func (s *SectionParser) checkSpecialFields(fieldName string) (match bool, data interface{}) {
	switch fieldName {
	case "MonitorFail":
		return true, s.monitorFail()
	case "MonitorURI":
		return true, s.monitorURI()
	case "StatsOptions":
		return true, s.statsOptions()
	case "HTTPCheck":
		return true, s.httpCheck()
	case "Forwardfor":
		return true, s.forwardfor()
	case "Redispatch":
		return true, s.redispatch()
	case "Balance":
		return true, s.balance()
	case "BindProcess":
		return true, s.bindProcess()
	case "Cookie":
		return true, s.cookie()
	case "HashType":
		return true, s.hashType()
	case "ErrorFiles":
		return true, s.errorFiles()
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
	case "Originalto":
		return true, s.originalto()
	default:
		return false, nil
	}
}

func (s *SectionParser) checkTimeouts(fieldName string) (match bool, data interface{}) {
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

func (s *SectionParser) checkSingleLine(fieldName string) (match bool, data interface{}) {
	if pName := misc.DashCase(fieldName); s.Parser.HasParser(s.Section, pName) {
		data, err := s.get(pName, false)
		if err != nil {
			return true, nil
		}
		return true, parseOption(data)
	}
	return false, nil
}

func (s *SectionParser) checkOptions(fieldName string) (match bool, data interface{}) {
	if pName := fmt.Sprintf("option %s", misc.DashCase(fieldName)); s.Parser.HasParser(s.Section, pName) {
		data, err := s.get(pName, false)
		if err != nil {
			return true, nil
		}
		return true, parseOption(data)
	}
	return false, nil
}

func (s *SectionParser) get(attribute string, createIfNotExists ...bool) (data common.ParserData, err error) {
	return s.Parser.Get(s.Section, s.Name, attribute, createIfNotExists...)
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

func (s *SectionParser) getSslChkData() (found bool, data interface{}) {
	data, err := s.get("option ssl-hello-chk", false)
	if err == nil {
		d := data.(*types.SimpleOption)
		if !d.NoOption {
			return true, "ssl-hello-chk"
		}
	}
	return false, nil
}

func (s *SectionParser) getSMTPChkData() (found bool, data interface{}) {
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

func (s *SectionParser) getLdapCheckData() (found bool, data interface{}) {
	data, err := s.get("option ldap-check", false)
	if err == nil {
		d := data.(*types.SimpleOption)
		if !d.NoOption {
			return true, "ldap-check"
		}
	}
	return false, nil
}

func (s *SectionParser) getMysqlCheckData() (found bool, data interface{}) {
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

func (s *SectionParser) getPgsqlCheckData() (found bool, data interface{}) {
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

func (s *SectionParser) getTCPCheckData() (found bool, data interface{}) {
	data, err := s.get("option tcp-check", false)
	if err == nil {
		d := data.(*types.SimpleOption)
		if !d.NoOption {
			return true, "tcp-check"
		}
	}
	return false, nil
}

func (s *SectionParser) getRedisCheckData() (found bool, data interface{}) {
	data, err := s.get("option redis-check", false)
	if err == nil {
		d := data.(*types.SimpleOption)
		if !d.NoOption {
			return true, "redis-check"
		}
	}
	return false, nil
}

func (s *SectionParser) getHttpchkData() (found bool, data interface{}) {
	data, err := s.get("option httpchk", false)
	if err == nil {
		d := data.(*types.OptionHttpchk)
		if !d.NoOption {
			s.setField("HttpchkParams", &models.HttpchkParams{
				Method:  d.Method,
				URI:     d.URI,
				Version: d.Version,
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
	return bst
}

func (s *SectionParser) defaultServer() interface{} { //nolint:gocognit,gocyclo,cyclop,maintidx
	data, err := s.get("default-server", false)
	if err != nil {
		return nil
	}
	d := data.([]types.DefaultServer)
	dServer := &models.DefaultServer{}
	for _, ds := range d {
		dsParams := ds.Params
		for _, p := range dsParams {
			switch v := p.(type) {
			case *params.ServerOptionWord:
				switch v.Name {
				case "agent-check":
					dServer.AgentCheck = "enabled"
				case "no-agent-check":
					dServer.AgentCheck = "disabled"
				case "allow-0rtt":
					dServer.Allow0rtt = true
				case "backup":
					dServer.Backup = "enabled"
				case "no-backup":
					dServer.Backup = "disabled"
				case "check":
					dServer.Check = "enabled"
				case "no-check":
					dServer.Check = "disabled"
				case "check-send-proxy":
					dServer.CheckSendProxy = "enabled"
				case "check-ssl":
					dServer.CheckSsl = "enabled"
				case "no-check-ssl":
					dServer.CheckSsl = "disabled"
				case "check-via-socks4":
					dServer.CheckViaSocks4 = "enabled"
				case "no-check-via-socks4":
					dServer.CheckViaSocks4 = "disabled"
				case "disabled":
					dServer.Disabled = "enabled"
				case "enabled":
					dServer.Enabled = "enabled"
				case "force-sslv3":
					dServer.ForceSslv3 = "enabled"
				case "no-force-sslv3":
					dServer.ForceSslv3 = "disabled"
				case "force-tlsv10":
					dServer.ForceTlsv10 = "enabled"
				case "no-force-tlsv10":
					dServer.ForceTlsv10 = "disabled"
				case "force-tlsv11":
					dServer.ForceTlsv11 = "enabled"
				case "no-force-tlsv11":
					dServer.ForceTlsv11 = "disabled"
				case "force-tlsv12":
					dServer.ForceTlsv12 = "enabled"
				case "no-force-tlsv12":
					dServer.ForceTlsv12 = "disabled"
				case "force-tlsv13":
					dServer.ForceTlsv13 = "enabled"
				case "no-force-tlsv13":
					dServer.ForceTlsv13 = "disabled"
				case "send-proxy":
					dServer.SendProxy = "enabled"
				case "no-send-proxy":
					dServer.SendProxy = "disabled"
				case "send-proxy-v2":
					dServer.SendProxyV2 = "enabled"
				case "no-send-proxy-v2":
					dServer.SendProxyV2 = "disabled"
				case "send-proxy-v2-ssl":
					dServer.SendProxyV2Ssl = "enabled"
				case "no-send-proxy-v2-ssl":
					dServer.SendProxyV2 = "disabled"
				case "send-proxy-v2-ssl-cn":
					dServer.SendProxyV2SslCn = "enabled"
				case "no-send-proxy-v2-ssl-cn":
					dServer.SendProxyV2SslCn = "disabled"
				case "ssl":
					dServer.Ssl = "enabled"
				case "no-ssl":
					dServer.Ssl = "disabled"
				case "ssl-reuse":
					dServer.SslReuse = "enabled"
				case "no-ssl-reuse":
					dServer.SslReuse = "disabled"
				case "no-sslv3":
					dServer.NoSslv3 = "disabled"
				case "tls-tickets":
					dServer.TLSTickets = "enabled"
				case "no-tls-tickets":
					dServer.TLSTickets = "disabled"
				case "no-tlsv10":
					dServer.NoTlsv10 = "enabled"
				case "no-tlsv11":
					dServer.NoTlsv11 = "enabled"
				case "no-tlsv12":
					dServer.NoTlsv12 = "enabled"
				case "no-tlsv13":
					dServer.NoTlsv13 = "enabled"
				case "no-verifyhost":
					dServer.NoVerifyhost = "enabled"
				case "no-tfo":
					dServer.Tfo = "disabled"
				case "non-stick":
					dServer.Stick = "disabled"
				case "stick":
					dServer.Stick = "enabled"
				case "tfo":
					dServer.Tfo = "enabled"
				}
			case *params.ServerOptionValue:
				switch v.Name {
				case "agent-send":
					dServer.AgentSend = v.Value
				case "agent-inter":
					dServer.AgentInter = misc.ParseTimeout(v.Value)
				case "agent-addr":
					dServer.AgentAddr = v.Value
				case "agent-port":
					port, err := strconv.ParseInt(v.Value, 10, 64)
					if err == nil && port != 0 {
						dServer.AgentPort = &port
					}
				case "alpn":
					dServer.Alpn = v.Value
				case "ca-file":
					dServer.CaFile = v.Value
				case "check-alpn":
					dServer.CheckAlpn = v.Value
				case "check-proto":
					dServer.CheckProto = v.Value
				case "check-sni":
					dServer.CheckSni = v.Value
				case "ciphers":
					dServer.Ciphers = v.Value
				case "ciphersuites":
					dServer.Ciphersuites = v.Value
				case "cookie":
					dServer.Cookie = v.Value
				case "crl-file":
					dServer.CrlFile = v.Value
				case "crt":
					dServer.SslCertificate = v.Value
				case "error-limit":
					count, err := strconv.ParseInt(v.Value, 10, 64)
					if err == nil {
						dServer.ErrorLimit = count
					}
				case "fall":
					dServer.Fall = misc.ParseTimeout(v.Value)
				case "init-addr":
					dServer.InitAddr = &v.Value
				case "inter":
					dServer.Inter = misc.ParseTimeout(v.Value)
				case "fastinter":
					dServer.Fastinter = misc.ParseTimeout(v.Value)
				case "downinter":
					dServer.Downinter = misc.ParseTimeout(v.Value)
				case "log-proto":
					dServer.LogProto = v.Value
				case "maxconn":
					maxconn, err := strconv.ParseInt(v.Value, 10, 64)
					if err == nil {
						dServer.Maxconn = &maxconn
					}
				case "maxqueue":
					maxqueue, err := strconv.ParseInt(v.Value, 10, 64)
					if err == nil {
						dServer.Maxqueue = &maxqueue
					}
				case "max-reuse":
					count, err := strconv.ParseInt(v.Value, 10, 64)
					if err == nil {
						dServer.MaxReuse = &count
					}
				case "minconn":
					minconn, err := strconv.ParseInt(v.Value, 10, 64)
					if err == nil {
						dServer.Minconn = &minconn
					}
				case "namespace":
					dServer.Namespace = v.Value
				case "npn":
					dServer.Npn = v.Value
				case "observe":
					dServer.Observe = v.Value
				case "on-error":
					dServer.OnError = v.Value
				case "on-marked-down":
					dServer.OnMarkedDown = v.Value
				case "on-marked-up":
					dServer.OnMarkedUp = v.Value
				case "pool-low-conn":
					max, err := strconv.ParseInt(v.Value, 10, 64)
					if err == nil {
						dServer.PoolLowConn = &max
					}
				case "pool-max-conn":
					max, err := strconv.ParseInt(v.Value, 10, 64)
					if err == nil {
						dServer.PoolMaxConn = &max
					}
				case "pool-purge-delay":
					delay, err := strconv.ParseInt(v.Value, 10, 64)
					if err == nil {
						dServer.PoolPurgeDelay = &delay
					}
				case "port":
					port, err := strconv.ParseInt(v.Value, 10, 64)
					if err == nil {
						dServer.Port = &port
					}
				case "proto":
					dServer.Proto = v.Value
				case "redir":
					dServer.Redir = v.Value
				case "rise":
					dServer.Rise = misc.ParseTimeout(v.Value)
				case "resolve-opts":
					dServer.ResolveOpts = v.Value
				case "resolve-prefer":
					dServer.ResolvePrefer = v.Value
				case "resolve-net":
					dServer.ResolveNet = v.Value
				case "resolvers":
					dServer.Resolvers = v.Value
				case "proxy-v2-options":
					dServer.ProxyV2Options = strings.Split(v.Value, ",")
				case "slowstart":
					dServer.Slowstart = misc.ParseTimeout(v.Value)
				case "sni":
					dServer.Sni = v.Value
				case "source":
					dServer.Source = v.Value
				case "ssl-max-ver":
					dServer.SslMaxVer = v.Value
				case "ssl-min-ver":
					dServer.SslMinVer = v.Value
				case "socks4":
					dServer.Socks4 = v.Value
				case "tcp-ut":
					delay, err := strconv.ParseInt(v.Value, 10, 64)
					if err == nil {
						dServer.TCPUt = delay
					}
				case "track":
					dServer.Track = v.Value
				case "verify":
					dServer.Verify = v.Value
				case "verifyhost":
					dServer.Verifyhost = v.Value
				case "weight":
					weight, err := strconv.ParseInt(v.Value, 10, 64)
					if err == nil {
						dServer.Weight = &weight
					}

				}
			}
		}
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
	if s.Section == parser.Defaults {
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
		return dEFiles
	}
	return nil
}

func (s *SectionParser) hashType() interface{} {
	data, err := s.get("hash-type", false)
	if err != nil {
		return nil
	}
	d := data.(*types.HashType)
	return &models.BackendHashType{
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
	return &models.Cookie{
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

func (s *SectionParser) bindProcess() interface{} {
	data, err := s.get("bind-process", false)
	if err != nil {
		return nil
	}
	d := data.(*types.BindProcess)
	return d.Process
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
	case *params.BalanceURLParam:
		b.URLParam = prm.Param
		b.URLParamCheckPost = prm.CheckPost
		b.URLParamMaxWait = prm.MaxWait
	}
	return b
}

func (s *SectionParser) redispatch() interface{} {
	data, err := s.get("option redispatch", false)
	if err != nil {
		return nil
	}
	d := data.(*types.OptionRedispatch)
	br := &models.Redispatch{}
	if d.Interval != nil {
		br.Interval = *d.Interval
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

func (s *SectionParser) httpCheck() interface{} {
	data, err := s.get("http-check", false)
	if err != nil {
		return nil
	}
	d := data.([]types.Action)
	if s.Section == parser.Defaults || s.Section == parser.Backends {
		for _, h := range d {
			httpCheck, err := ParseHTTPCheck(h)
			if err != nil {
				continue
			}
			httpCheck.Index = misc.Int64P(0)
			return httpCheck
		}
	}
	return nil
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
			if v.Name != "" {
				opt.StatsShowNodeName = misc.StringP(v.Name)
			}
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

func (s *SectionParser) compression() interface{} {
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

	data, err = s.get("compression type", false)
	if err == nil {
		d, ok := data.(*types.StringSliceC)
		if ok && d != nil && len(d.Value) > 0 {
			compressionFound = true
			compression.Types = d.Value
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

// SectionObject represents a configuration section
type SectionObject struct {
	Object  interface{}
	Section parser.Section
	Name    string
	Parser  parser.Parser
}

// CreateEditSection creates or updates a section in the parser based on the provided object
func CreateEditSection(object interface{}, section parser.Section, pName string, p parser.Parser) error {
	so := SectionObject{
		Object:  object,
		Section: section,
		Name:    pName,
		Parser:  p,
	}
	return so.CreateEditSection()
}

// CreateEditSection creates or updates a section in the parser based on the provided object
func (s *SectionObject) CreateEditSection() error {
	objValue := reflect.ValueOf(s.Object)
	if objValue.Kind() == reflect.Ptr {
		objValue = reflect.ValueOf(s.Object).Elem()
	}
	for i := 0; i < objValue.NumField(); i++ {
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

	return errors.Errorf("Cannot parse option for %s %s: %s", s.Section, s.Name, fieldName)
}

func (s *SectionObject) checkParams(fieldName string) (match bool) {
	return strings.HasSuffix(fieldName, "Params")
}

func (s *SectionObject) checkSpecialFields(fieldName string, field reflect.Value) (match bool, err error) {
	switch fieldName {
	case "MonitorURI":
		return true, s.monitorURI(field)
	case "MonitorFail":
		return true, s.monitorFail(field)
	case "StatsOptions":
		return true, s.statsOptions(field)
	case "HTTPCheck":
		return true, s.httpCheck(field)
	case "Forwardfor":
		return true, s.forwardfor(field)
	case "Redispatch":
		return true, s.redispatch(field)
	case "Balance":
		return true, s.balance(field)
	case "BindProcess":
		return true, s.bindProcess(field)
	case "Cookie":
		return true, s.cookie(field)
	case "HashType":
		return true, s.hashType(field)
	case "ErrorFiles":
		return true, s.errorFiles(field)
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
	case "Originalto":
		return true, s.originalto(field)
	default:
		return false, nil
	}
}

func (s *SectionObject) checkTimeouts(fieldName string, field reflect.Value) (match bool, err error) {
	if strings.HasSuffix(fieldName, "Timeout") {
		if pName := translateTimeout(fieldName); s.Parser.HasParser(s.Section, pName) {
			if valueIsNil(field) {
				if err := s.set(pName, nil); err != nil {
					return true, err
				}
				return true, nil
			}
			t := &types.SimpleTimeout{}
			t.Value = strconv.FormatInt(field.Elem().Int(), 10)
			if err := s.set(pName, t); err != nil {
				return true, err
			}
		}
		return true, nil
	}
	return false, nil
}

func (s *SectionObject) checkOptions(fieldName string, field reflect.Value) (match bool, err error) {
	if pName := fmt.Sprintf("option %s", misc.DashCase(fieldName)); s.Parser.HasParser(s.Section, pName) {
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

func (s *SectionObject) checkSingleLine(fieldName string, field reflect.Value) (match bool, err error) {
	if pName := misc.DashCase(fieldName); s.Parser.HasParser(s.Section, pName) {
		if valueIsNil(field) {
			if err := s.set(pName, nil); err != nil {
				return true, err
			}
			return true, nil
		}
		d := translateToParserData(field)
		if d == nil {
			return true, errors.Errorf("Cannot parse type for %s %s: %s", s.Section, s.Name, fieldName)
		}
		if err := s.set(pName, d); err != nil {
			return true, err
		}
		return true, nil
	}
	return false, nil
}

func (s *SectionObject) isNoOption(attribute string) bool {
	data, err := s.Parser.Get(s.Section, s.Name, attribute)
	if err != nil {
		return false
	}

	simpleOpt, ok := data.(*types.SimpleOption)
	if !ok {
		return false
	}

	return simpleOpt.NoOption
}

func (s *SectionObject) set(attribute string, data interface{}) error {
	return s.Parser.Set(s.Section, s.Name, attribute, data)
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
		if err := s.set("unique-id-header", nil); err != nil {
			return err
		}
		return nil
	}
	d := types.UniqueIDHeader{
		Name: field.String(),
	}
	if err := s.set("unique-id-header", &d); err != nil {
		return err
	}
	return nil
}

func (s *SectionObject) uniqueIDFormat(field reflect.Value) error {
	if s.Section != parser.Defaults && s.Section != parser.Frontends {
		return nil
	}
	if valueIsNil(field) {
		if err := s.set("unique-id-format", nil); err != nil {
			return err
		}
		return nil
	}
	d := types.UniqueIDFormat{
		LogFormat: field.String(),
	}
	if err := s.set("unique-id-format", &d); err != nil {
		return err
	}
	return nil
}

func (s *SectionObject) httpReuse(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			if err := s.set("http-reuse", nil); err != nil {
				return err
			}
			return nil
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
		attribute := fmt.Sprintf("option %s", opt)

		if s.isNoOption(attribute) {
			continue
		}

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
			if err := s.set("default_backend", nil); err != nil {
				return err
			}
			return nil
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
	if err := s.set("option httpchk", nil); err != nil {
		return err
	}
	return nil
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
	if s.Section == parser.Backends || s.Section == parser.Frontends {
		if valueIsNil(field) {
			if err := s.set("stick-table", nil); err != nil {
				return err
			}
			return nil
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
			d.Expire = strconv.FormatInt(*st.Expire, 10)
		}
		if st.Size != nil {
			d.Size = strconv.FormatInt(*st.Size, 10)
		}
		if err := s.set("stick-table", d); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) defaultServer(field reflect.Value) error { //nolint:gocognit,gocyclo,cyclop,maintidx
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			if err := s.set("default-server", nil); err != nil {
				return err
			}
			return nil
		}
		ds, ok := field.Elem().Interface().(models.DefaultServer)
		if !ok {
			return misc.CreateTypeAssertError("default-server")
		}
		dServers := []types.DefaultServer{{}}

		ps := make([]params.ServerOption, 0, 4)

		// ServerOptionWord
		if ds.AgentCheck == "enabled" {
			param := &params.ServerOptionWord{
				Name: "agent-check",
			}
			ps = append(ps, param)
		}
		if ds.Allow0rtt {
			param := &params.ServerOptionWord{
				Name: "allow-0rtt",
			}
			ps = append(ps, param)
		}
		if ds.Backup == "enabled" {
			param := &params.ServerOptionWord{
				Name: "backup",
			}
			ps = append(ps, param)
		}
		if ds.Check == "enabled" {
			param := &params.ServerOptionWord{
				Name: "check",
			}
			ps = append(ps, param)
		}
		if ds.CheckSendProxy == "enabled" {
			param := &params.ServerOptionWord{
				Name: "check-send-proxy",
			}
			ps = append(ps, param)
		}
		if ds.CheckSsl == "enabled" {
			param := &params.ServerOptionWord{
				Name: "check-ssl",
			}
			ps = append(ps, param)
		}
		if ds.CheckViaSocks4 == "enabled" {
			param := &params.ServerOptionWord{
				Name: "check-via-socks4",
			}
			ps = append(ps, param)
		}
		if ds.ForceSslv3 == "enabled" {
			param := &params.ServerOptionWord{
				Name: "force-sslv3",
			}
			ps = append(ps, param)
		}
		if ds.ForceTlsv10 == "enabled" {
			param := &params.ServerOptionWord{
				Name: "force-tlsv10",
			}
			ps = append(ps, param)
		}
		if ds.ForceTlsv11 == "enabled" {
			param := &params.ServerOptionWord{
				Name: "force-tlsv11",
			}
			ps = append(ps, param)
		}
		if ds.ForceTlsv12 == "enabled" {
			param := &params.ServerOptionWord{
				Name: "force-tlsv12",
			}
			ps = append(ps, param)
		}
		if ds.ForceTlsv13 == "enabled" {
			param := &params.ServerOptionWord{
				Name: "force-tlsv13",
			}
			ps = append(ps, param)
		}
		if ds.SendProxy == "enabled" {
			param := &params.ServerOptionWord{
				Name: "send-proxy",
			}
			ps = append(ps, param)
		}
		if ds.SendProxyV2 == "enabled" {
			param := &params.ServerOptionWord{
				Name: "send-proxy-v2",
			}
			ps = append(ps, param)
		}
		if ds.SendProxyV2Ssl == "enabled" {
			param := &params.ServerOptionWord{
				Name: "send-proxy-v2-ssl",
			}
			ps = append(ps, param)
		}
		if ds.SendProxyV2SslCn == "enabled" {
			param := &params.ServerOptionWord{
				Name: "send-proxy-v2-ssl-cn",
			}
			ps = append(ps, param)
		}
		if ds.Ssl == "enabled" {
			param := &params.ServerOptionWord{
				Name: "ssl",
			}
			ps = append(ps, param)
		}
		if ds.SslReuse == "enabled" {
			param := &params.ServerOptionWord{
				Name: "ssl-reuse",
			}
			ps = append(ps, param)
		}
		if ds.TLSTickets == "enabled" {
			param := &params.ServerOptionWord{
				Name: "tls-tickets",
			}
			ps = append(ps, param)
		}
		if ds.Stick == "enabled" {
			param := &params.ServerOptionWord{
				Name: "stick",
			}
			ps = append(ps, param)
		}
		if ds.Tfo == "enabled" {
			param := &params.ServerOptionWord{
				Name: "tfo",
			}
			ps = append(ps, param)
		}
		// ServerOptionValue
		if ds.AgentSend != "" {
			param := &params.ServerOptionValue{
				Name:  "agent-send",
				Value: ds.AgentSend,
			}
			ps = append(ps, param)
		}
		if ds.AgentInter != nil {
			param := &params.ServerOptionValue{
				Name:  "agent-inter",
				Value: strconv.FormatInt(*ds.AgentInter, 10),
			}
			ps = append(ps, param)
		}
		if ds.AgentAddr != "" {
			param := &params.ServerOptionValue{
				Name:  "agent-addr",
				Value: ds.AgentAddr,
			}
			ps = append(ps, param)
		}
		if ds.AgentPort != nil {
			param := &params.ServerOptionValue{
				Name:  "agent-port",
				Value: strconv.FormatInt(*ds.AgentPort, 10),
			}
			ps = append(ps, param)
		}
		if ds.Alpn != "" {
			param := &params.ServerOptionValue{
				Name:  "alpn",
				Value: ds.Alpn,
			}
			ps = append(ps, param)
		}
		if ds.CheckAlpn != "" {
			param := &params.ServerOptionValue{
				Name:  "check-alpn",
				Value: ds.CheckAlpn,
			}
			ps = append(ps, param)
		}
		if ds.CheckProto != "" {
			param := &params.ServerOptionValue{
				Name:  "check-proto",
				Value: ds.CheckProto,
			}
			ps = append(ps, param)
		}
		if ds.CheckSni != "" {
			param := &params.ServerOptionValue{
				Name:  "check-sni",
				Value: ds.CheckSni,
			}
			ps = append(ps, param)
		}
		if ds.Ciphers != "" {
			param := &params.ServerOptionValue{
				Name:  "ciphers",
				Value: ds.Ciphers,
			}
			ps = append(ps, param)
		}
		if ds.Ciphersuites != "" {
			param := &params.ServerOptionValue{
				Name:  "ciphersuites",
				Value: ds.Ciphersuites,
			}
			ps = append(ps, param)
		}
		if ds.Cookie != "" {
			param := &params.ServerOptionValue{
				Name:  "cookie",
				Value: ds.Cookie,
			}
			ps = append(ps, param)
		}
		if ds.CrlFile != "" {
			param := &params.ServerOptionValue{
				Name:  "crl-file",
				Value: ds.CrlFile,
			}
			ps = append(ps, param)
		}
		if ds.SslCertificate != "" { // crt command
			param := &params.ServerOptionValue{
				Name:  "crt",
				Value: ds.SslCertificate,
			}
			ps = append(ps, param)
		}
		if ds.ErrorLimit != 0 {
			param := &params.ServerOptionValue{
				Name:  "error-limit",
				Value: strconv.FormatInt(ds.ErrorLimit, 10),
			}
			ps = append(ps, param)
		}
		if ds.Fall != nil {
			param := &params.ServerOptionValue{
				Name:  "fall",
				Value: strconv.FormatInt(*ds.Fall, 10),
			}
			ps = append(ps, param)
		}
		if ds.InitAddr != nil {
			param := &params.ServerOptionValue{
				Name:  "init-addr",
				Value: *ds.InitAddr,
			}
			ps = append(ps, param)
		}
		if ds.Inter != nil {
			param := &params.ServerOptionValue{
				Name:  "inter",
				Value: strconv.FormatInt(*ds.Inter, 10),
			}
			ps = append(ps, param)
		}
		if ds.Fastinter != nil {
			param := &params.ServerOptionValue{
				Name:  "fastinter",
				Value: strconv.FormatInt(*ds.Fastinter, 10),
			}
			ps = append(ps, param)
		}
		if ds.Downinter != nil {
			param := &params.ServerOptionValue{
				Name:  "downinter",
				Value: strconv.FormatInt(*ds.Downinter, 10),
			}
			ps = append(ps, param)
		}
		if ds.LogProto != "" {
			param := &params.ServerOptionValue{
				Name:  "log-proto",
				Value: ds.LogProto,
			}
			ps = append(ps, param)
		}
		if ds.Maxconn != nil {
			param := &params.ServerOptionValue{
				Name:  "maxconn",
				Value: strconv.FormatInt(*ds.Maxconn, 10),
			}
			ps = append(ps, param)
		}
		if ds.Maxqueue != nil {
			param := &params.ServerOptionValue{
				Name:  "maxqueue",
				Value: strconv.FormatInt(*ds.Maxqueue, 10),
			}
			ps = append(ps, param)
		}
		if ds.MaxReuse != nil {
			param := &params.ServerOptionValue{
				Name:  "max-reuse",
				Value: strconv.FormatInt(*ds.MaxReuse, 10),
			}
			ps = append(ps, param)
		}
		if ds.Minconn != nil {
			param := &params.ServerOptionValue{
				Name:  "minconn",
				Value: strconv.FormatInt(*ds.Minconn, 10),
			}
			ps = append(ps, param)
		}
		if ds.Namespace != "" {
			param := &params.ServerOptionValue{
				Name:  "namespace",
				Value: ds.Namespace,
			}
			ps = append(ps, param)
		}
		if ds.Npn != "" {
			param := &params.ServerOptionValue{
				Name:  "npn",
				Value: ds.Npn,
			}
			ps = append(ps, param)
		}
		if ds.Observe != "" {
			param := &params.ServerOptionValue{
				Name:  "observe",
				Value: ds.Observe,
			}
			ps = append(ps, param)
		}
		if ds.OnError != "" {
			param := &params.ServerOptionValue{
				Name:  "on-error",
				Value: ds.OnError,
			}
			ps = append(ps, param)
		}
		if ds.OnMarkedDown != "" {
			param := &params.ServerOptionValue{
				Name:  "on-marked-down",
				Value: ds.OnMarkedDown,
			}
			ps = append(ps, param)
		}
		if ds.OnMarkedUp != "" {
			param := &params.ServerOptionValue{
				Name:  "on-marked-up",
				Value: ds.OnMarkedUp,
			}
			ps = append(ps, param)
		}
		if ds.PoolLowConn != nil {
			param := &params.ServerOptionValue{
				Name:  "pool-low-conn",
				Value: strconv.FormatInt(*ds.PoolLowConn, 10),
			}
			ps = append(ps, param)
		}
		if ds.PoolMaxConn != nil {
			param := &params.ServerOptionValue{
				Name:  "pool-max-conn",
				Value: strconv.FormatInt(*ds.PoolMaxConn, 10),
			}
			ps = append(ps, param)
		}
		if ds.PoolPurgeDelay != nil {
			param := &params.ServerOptionValue{
				Name:  "pool-purge-delay",
				Value: strconv.FormatInt(*ds.PoolPurgeDelay, 10),
			}
			ps = append(ps, param)
		}
		if ds.Port != nil {
			param := &params.ServerOptionValue{
				Name:  "port",
				Value: strconv.FormatInt(*ds.Port, 10),
			}
			ps = append(ps, param)
		}
		if ds.Proto != "" {
			param := &params.ServerOptionValue{
				Name:  "proto",
				Value: ds.Proto,
			}
			ps = append(ps, param)
		}
		if ds.Redir != "" {
			param := &params.ServerOptionValue{
				Name:  "redir",
				Value: ds.Redir,
			}
			ps = append(ps, param)
		}
		if ds.Rise != nil {
			param := &params.ServerOptionValue{
				Name:  "rise",
				Value: strconv.FormatInt(*ds.Rise, 10),
			}
			ps = append(ps, param)
		}
		if ds.ResolveOpts != "" {
			param := &params.ServerOptionValue{
				Name:  "resolve-opts",
				Value: ds.ResolveOpts,
			}
			ps = append(ps, param)
		}
		if ds.ResolvePrefer != "" {
			param := &params.ServerOptionValue{
				Name:  "resolve-prefer",
				Value: ds.ResolvePrefer,
			}
			ps = append(ps, param)
		}
		if ds.ResolveNet != "" {
			param := &params.ServerOptionValue{
				Name:  "resolve-net",
				Value: ds.ResolveNet,
			}
			ps = append(ps, param)
		}
		if ds.Resolvers != "" {
			param := &params.ServerOptionValue{
				Name:  "resolvers",
				Value: ds.Resolvers,
			}
			ps = append(ps, param)
		}
		if len(ds.ProxyV2Options) > 0 {
			param := &params.ServerOptionValue{
				Name:  "proxy-v2-options",
				Value: strings.Join(ds.ProxyV2Options, ","),
			}
			ps = append(ps, param)
		}
		if ds.Slowstart != nil {
			param := &params.ServerOptionValue{
				Name:  "slowstart",
				Value: strconv.FormatInt(*ds.Slowstart, 10),
			}
			ps = append(ps, param)
		}
		if ds.Sni != "" {
			param := &params.ServerOptionValue{
				Name:  "sni",
				Value: ds.Sni,
			}
			ps = append(ps, param)
		}
		if ds.Source != "" {
			param := &params.ServerOptionValue{
				Name:  "source",
				Value: ds.Source,
			}
			ps = append(ps, param)
		}
		if ds.SslMaxVer != "" {
			param := &params.ServerOptionValue{
				Name:  "ssl-max-ver",
				Value: ds.SslMaxVer,
			}
			ps = append(ps, param)
		}
		if ds.SslMinVer != "" {
			param := &params.ServerOptionValue{
				Name:  "ssl-min-ver",
				Value: ds.SslMinVer,
			}
			ps = append(ps, param)
		}
		if ds.Socks4 != "" {
			param := &params.ServerOptionValue{
				Name:  "socks4",
				Value: ds.Socks4,
			}
			ps = append(ps, param)
		}
		if ds.TCPUt != 0 {
			param := &params.ServerOptionValue{
				Name:  "tcp-ut",
				Value: strconv.FormatInt(ds.TCPUt, 10),
			}
			ps = append(ps, param)
		}
		if ds.Track != "" {
			param := &params.ServerOptionValue{
				Name:  "track",
				Value: ds.Track,
			}
			ps = append(ps, param)
		}
		if ds.Verify != "" {
			param := &params.ServerOptionValue{
				Name:  "verify",
				Value: ds.Verify,
			}
			ps = append(ps, param)
		}
		if ds.Verifyhost != "" {
			param := &params.ServerOptionValue{
				Name:  "verifyhost",
				Value: ds.Verifyhost,
			}
			ps = append(ps, param)
		}
		if ds.Weight != nil {
			param := &params.ServerOptionValue{
				Name:  "weight",
				Value: strconv.FormatInt(*ds.Weight, 10),
			}
			ps = append(ps, param)
		}

		dServers[0].Params = ps
		if err := s.set("default-server", dServers); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) loadServerStateFromFile(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			if err := s.set("load-server-state-from-file", nil); err != nil {
				return err
			}
			return nil
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
	if s.Section == parser.Defaults {
		if valueIsNil(field) {
			if err := s.set("errorfile", nil); err != nil {
				return err
			}
			return nil
		}
		efs, ok := field.Interface().([]*models.Errorfile)
		if !ok {
			return nil
		}
		errorFiles := []types.ErrorFile{}
		for _, ef := range efs {
			errorFiles = append(errorFiles, types.ErrorFile{Code: strconv.FormatInt(ef.Code, 10), File: ef.File})
		}
		if err := s.set("errorfile", errorFiles); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) hashType(field reflect.Value) error {
	if s.Section == parser.Backends {
		if valueIsNil(field) {
			if err := s.set("hash-type", nil); err != nil {
				return err
			}
			return nil
		}
		b, ok := field.Elem().Interface().(models.BackendHashType)
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
			if err := s.set("cookie", nil); err != nil {
				return err
			}
			return nil
		}
		d, ok := field.Elem().Interface().(models.Cookie)
		if !ok {
			return misc.CreateTypeAssertError("cookie")
		}
		domains := make([]string, len(d.Domains))
		for i, domain := range d.Domains {
			domains[i] = domain.Value
		}
		data := types.Cookie{
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

func (s *SectionObject) bindProcess(field reflect.Value) error {
	if s.Section == parser.Defaults || s.Section == parser.Frontends || s.Section == parser.Backends {
		if valueIsNil(field) {
			if err := s.set("bind-process", nil); err != nil {
				return err
			}
			return nil
		}
		b := field.String()
		d := &types.BindProcess{
			Process: b,
		}
		if err := s.set("bind-process", d); err != nil {
			return err
		}
	}
	return nil
}

func (s *SectionObject) balance(field reflect.Value) error {
	if s.Section == parser.Backends || s.Section == parser.Defaults {
		if valueIsNil(field) {
			if err := s.set("balance", nil); err != nil {
				return err
			}
			return nil
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
			if err := s.set("option redispatch", nil); err != nil {
				return err
			}
			return nil
		}
		br, ok := field.Elem().Interface().(models.Redispatch)
		if !ok {
			return misc.CreateTypeAssertError("option redispatch")
		}
		d := &types.OptionRedispatch{
			Interval: &br.Interval,
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
		if err := s.set("option forwardfor", nil); err != nil {
			return err
		}
		return nil
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
	if err := s.set("option forwardfor", d); err != nil {
		return err
	}
	return nil
}

func (s *SectionObject) httpCheck(field reflect.Value) error {
	if s.Section == parser.Defaults || s.Section == parser.Backends {
		hc, ok := field.Interface().(*models.HTTPCheck)
		if !ok {
			return misc.CreateTypeAssertError("http-check")
		}

		if hc == nil {
			return nil
		}
		httpChecks, err := ParseHTTPChecks(string(s.Section), s.Name, s.Parser)
		if err != nil {
			return err
		}
		if hc != nil {
			for _, httpCheck := range httpChecks {
				if reflect.DeepEqual(httpCheck, hc) {
					return nil
				}
			}
			check, err := SerializeHTTPCheck(*hc)
			if err != nil {
				return err
			}
			if err = s.Parser.Insert(s.Section, s.Name, "http-check", check, 0); err != nil {
				return err
			}
		}
	}
	return nil
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

func (s *SectionObject) statsOptions(field reflect.Value) error {
	if valueIsNil(field) {
		if err := s.set("stats", nil); err != nil {
			return err
		}
		return nil
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
			reqType = fmt.Sprintf("auth realm %s", httpRequest.Realm)
		}
		s := &stats.HTTPRequest{
			Type:     reqType,
			Cond:     httpRequest.Cond,
			CondTest: httpRequest.CondTest,
		}
		ss = append(ss, s)
	}
	if err := s.set("stats", ss); err != nil {
		return err
	}
	return nil
}

func (s *SectionObject) compression(field reflect.Value) error {
	var err error
	if valueIsNil(field) {
		err = s.set("compression algo", nil)
		if err != nil {
			return err
		}
		err = s.set("compression type", nil)
		if err != nil {
			return err
		}
		err = s.set("compression offload", nil)
		if err != nil {
			return err
		}
		return nil
	}
	compression, ok := field.Elem().Interface().(models.Compression)
	if !ok {
		return fmt.Errorf("error casting compression model")
	}

	if len(compression.Algorithms) > 0 {
		err = s.set("compression algo", &types.StringSliceC{Value: compression.Algorithms})
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
	if compression.Offload {
		err = s.set("compression offload", &types.Enabled{})
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
	return s.set("clitcpka-idle", types.StringC{Value: fmt.Sprintf("%dms", v)})
}

func (s *SectionObject) clitcpkaIntvl(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("clitcpka-intvl", nil)
	}
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}
	v := field.Int()
	return s.set("clitcpka-intvl", types.StringC{Value: fmt.Sprintf("%dms", v)})
}

func (s *SectionObject) srvtcpkaIdle(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("srvtcpka-idle", nil)
	}
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}
	v := field.Int()
	return s.set("srvtcpka-idle", types.StringC{Value: fmt.Sprintf("%dms", v)})
}

func (s *SectionObject) srvtcpkaIntvl(field reflect.Value) error {
	if valueIsNil(field) {
		return s.set("srvtcpka-intvl", nil)
	}
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}
	v := field.Int()
	return s.set("srvtcpka-intvl", types.StringC{Value: fmt.Sprintf("%dms", v)})
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

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
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

	if err := CreateEditSection(data, section, name, p); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
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

	if err := CreateEditSection(data, section, name, p); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
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
	switch v.Kind() { //nolint:exhaustive
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
	return fmt.Sprintf("timeout %s", misc.DashCase(mName))
}
