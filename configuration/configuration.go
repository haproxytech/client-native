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

	"github.com/haproxytech/config-parser/common"
	"github.com/pkg/errors"

	"github.com/haproxytech/client-native/misc"
	parser "github.com/haproxytech/config-parser"
	parser_errors "github.com/haproxytech/config-parser/errors"
	"github.com/haproxytech/config-parser/params"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

const (
	//DefaultConfigurationFile sane default for path to haproxy configuration file
	DefaultConfigurationFile string = "/etc/haproxy/haproxy.cfg"
	//DefaultHaproxy sane default for path to haproxy executable
	DefaultHaproxy string = "/usr/sbin/haproxy"
	//DefaultUseValidation sane default using validation in client native
	DefaultUseValidation bool = true
	//DefaultPersistentTransactions sane default using persistent transactions in client native
	DefaultPersistentTransactions bool = true
	// DefaultTransactionDir sane default for path for transactions
	DefaultTransactionDir string = "/tmp/haproxy"
)

// ClientParams is just a placeholder for all client options
type ClientParams struct {
	ConfigurationFile      string
	Haproxy                string
	UseValidation          bool
	PersistentTransactions bool
	TransactionDir         string
	BackupsNumber          int
	MasterWorker           bool
}

// Client configuration client
// Parser is the config parser instance that loads "master" configuration file on Init
// and when transaction is committed it gets replaced with the parser from parsers map.
// parsers map contains a config parser for each transaction, which loads data from
// transaction files on StartTransaction, and deletes on CommitTransaction. We save
// data to file on every change for persistence.
type Client struct {
	ClientParams
	parsers map[string]*parser.Parser
	Parser  *parser.Parser
	mu      sync.Mutex
}

// DefaultClient returns Client with sane defaults
func DefaultClient() (*Client, error) {
	p := ClientParams{
		ConfigurationFile:      DefaultConfigurationFile,
		Haproxy:                DefaultHaproxy,
		UseValidation:          DefaultUseValidation,
		PersistentTransactions: DefaultPersistentTransactions,
		TransactionDir:         DefaultTransactionDir,
		MasterWorker:           false,
	}
	c := &Client{}
	err := c.Init(p)

	if err != nil {
		return nil, err
	}

	return c, err
}

// Init initializes a Client
func (c *Client) Init(options ClientParams) error {
	if options.TransactionDir == "" {
		options.TransactionDir = DefaultTransactionDir
	}

	if options.ConfigurationFile == "" {
		options.ConfigurationFile = DefaultConfigurationFile
	}

	if options.Haproxy == "" {
		options.Haproxy = DefaultHaproxy
	}

	c.ClientParams = options

	c.parsers = make(map[string]*parser.Parser)
	if err := c.InitTransactionParsers(); err != nil {
		return err
	}

	c.Parser = &parser.Parser{}
	if err := c.Parser.LoadData(options.ConfigurationFile); err != nil {
		return NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read %s", c.ConfigurationFile))
	}

	return nil
}

// GetParser returns a parser for given transaction, if transaction is "", it returns "master" parser
func (c *Client) GetParser(transaction string) (*parser.Parser, error) {
	if transaction == "" {
		return c.Parser, nil
	}
	p, ok := c.parsers[transaction]
	if !ok {
		return nil, NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("Transaction %s does not exist", transaction))
	}
	return p, nil
}

//AddParser adds parser to parser map
func (c *Client) AddParser(transaction string) error {
	if transaction == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Not a valid transaction"))
	}
	_, ok := c.parsers[transaction]
	if ok {
		return NewConfError(ErrTransactionAlreadyExists, fmt.Sprintf("Transaction %s already exists", transaction))
	}

	p := &parser.Parser{}
	tFile := ""
	var err error
	if c.PersistentTransactions {
		tFile, err = c.getTransactionFile(transaction)
		if err != nil {
			return err
		}
	} else {
		tFile = c.ConfigurationFile
	}
	if err := p.LoadData(tFile); err != nil {
		return NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read %s", tFile))
	}
	c.parsers[transaction] = p
	return nil
}

//DeleteParser deletes parser from parsers map
func (c *Client) DeleteParser(transaction string) error {
	if transaction == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Not a valid transaction"))
	}
	_, ok := c.parsers[transaction]
	if !ok {
		return NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("Transaction %s does not exist", transaction))
	}
	delete(c.parsers, transaction)
	return nil
}

//CommitParser commits transaction parser, deletes it from parsers map, and replaces master Parser
func (c *Client) CommitParser(transaction string) error {
	if transaction == "" {
		return NewConfError(ErrValidationError, fmt.Sprintf("Not a valid transaction"))
	}
	p, ok := c.parsers[transaction]
	if !ok {
		return NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("Transaction %s does not exist", transaction))
	}
	c.Parser = p
	delete(c.parsers, transaction)
	return nil
}

//InitTransactionParsers checks transactions and initializes parsers map with transactions in_progress
func (c *Client) InitTransactionParsers() error {
	transactions, err := c.GetTransactions("in_progress")
	if err != nil {
		return err
	}

	for _, t := range *transactions {
		if err := c.AddParser(t.ID); err != nil {
			continue
		}
		p, err := c.GetParser(t.ID)
		if err != nil {
			continue
		}
		tFile, err := c.getTransactionFile(t.ID)
		if err != nil {
			return err
		}
		if err := p.LoadData(tFile); err != nil {
			return NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read %s", tFile))
		}
	}
	return nil
}

// GetVersion returns configuration file version
func (c *Client) GetVersion(transaction string) (int64, error) {
	return c.getVersion(transaction)
}

func (c *Client) getVersion(transaction string) (int64, error) {
	p, err := c.GetParser(transaction)
	if err != nil {
		return 0, NewConfError(ErrCannotReadVersion, fmt.Sprintf("Cannot read version: %s", err.Error()))
	}

	data, _ := p.Get(parser.Comments, parser.CommentsSectionName, "# _version", true)
	ver, _ := data.(*types.ConfigVersion)
	return ver.Value, nil
}

func (c *Client) incrementVersion() error {
	data, _ := c.Parser.Get(parser.Comments, parser.CommentsSectionName, "# _version", true)
	ver, _ := data.(*types.ConfigVersion)
	ver.Value = ver.Value + 1

	if err := c.Parser.Save(c.ConfigurationFile); err != nil {
		return NewConfError(ErrCannotSetVersion, fmt.Sprintf("Cannot set version: %s", err.Error()))
	}
	return nil
}

func (c *Client) checkTransactionOrVersion(transactionID string, version int64) (string, error) {
	// start an implicit transaction if transaction is not already given
	t := ""
	if transactionID != "" && version != 0 {
		return "", NewConfError(ErrBothVersionTransaction, "Both version and transaction specified, specify only one")
	} else if transactionID == "" && version == 0 {
		return "", NewConfError(ErrNoVersionTransaction, "Version or transaction not specified, specify only one")
	} else if transactionID != "" {
		t = transactionID
	} else {
		v, err := c.GetVersion("")
		if err != nil {
			return "", err
		}
		if version != v {
			return "", NewConfError(ErrVersionMismatch, fmt.Sprintf("Version in configuration file is %v, given version is %v", v, version))
		}

		transaction, err := c.startTransaction(version)
		if err != nil {
			return "", err
		}
		t = transaction.ID

	}
	return t, nil
}

func (c *Client) parseSection(object interface{}, section parser.Section, pName string, p *parser.Parser) error {
	objValue := reflect.ValueOf(object).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		typeField := objValue.Type().Field(i)
		field := objValue.FieldByName(typeField.Name)
		if typeField.Name != "Name" && typeField.Name != "ID" {
			val := c.parseField(section, pName, typeField.Name, p)
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
	}

	return nil
}

func (c *Client) parseField(section parser.Section, sectionName string, fieldName string, p *parser.Parser) interface{} {
	//Handle special cases
	if fieldName == "Httpchk" {
		data, err := p.Get(section, sectionName, "option httpchk", false)
		if err != nil {
			return nil
		}
		d := data.(*types.OptionHttpchk)
		return &models.Httpchk{
			Method:  d.Method,
			URI:     d.URI,
			Version: d.Version,
		}
	}
	if fieldName == "Forwardfor" {
		data, err := p.Get(section, sectionName, "option forwardfor", false)
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
	if fieldName == "Redispatch" {
		data, err := p.Get(section, sectionName, "option redispatch", false)
		if err != nil {
			return nil
		}
		d := data.(*types.OptionRedispatch)
		br := &models.Redispatch{}
		if d.Interval != nil {
			br.Interval = *d.Interval
		}
		if d.NoOption == true {
			d := "disabled"
			br.Enabled = &d
		} else {
			e := "enabled"
			br.Enabled = &e
		}
		return br
	}
	if fieldName == "Balance" {
		data, err := p.Get(section, sectionName, "balance", false)
		if err != nil {
			return nil
		}
		d := data.(*types.Balance)
		return &models.Balance{
			Algorithm: d.Algorithm,
			Arguments: d.Arguments,
		}
	}
	if fieldName == "HashType" {
		data, err := p.Get(section, sectionName, "hash-type", false)
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
	if fieldName == "ErrorFiles" {
		data, err := p.Get(section, sectionName, "errorfile", false)
		if err != nil {
			return nil
		}
		d := data.([]types.ErrorFile)
		if section == parser.Defaults {
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
	if fieldName == "DefaultServer" {
		data, err := p.Get(section, sectionName, "default-server", false)
		if err != nil {
			return nil
		}
		d := data.([]types.DefaultServer)
		dServer := &models.DefaultServer{}
		for _, ds := range d {
			dsParams := ds.Params
			for _, p := range dsParams {
				v, ok := p.(*params.ServerOptionValue)
				if ok {
					switch v.Name {
					case "fall":
						dServer.Fall = misc.ParseTimeout(v.Value)
					case "inter":
						dServer.Inter = misc.ParseTimeout(v.Value)
					case "rise":
						dServer.Rise = misc.ParseTimeout(v.Value)
					case "port":
						port, err := strconv.ParseInt(v.Value, 10, 64)
						if err == nil {
							dServer.Port = &port
						}
					}
				}
			}
			return dServer
		}
		return nil
	}
	if fieldName == "StickTable" {
		data, err := p.Get(section, sectionName, "stick-table", false)
		if err != nil {
			return nil
		}
		d := data.(*types.StickTable)
		bst := &models.BackendStickTable{}

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
	if fieldName == "AdvCheck" {
		data, err := p.Get(section, sectionName, "option ssl-hello-chk", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "ssl-hello-chk"
			}
		}

		data, err = p.Get(section, sectionName, "option smtpchk", false)
		if err == nil {
			d := data.(*types.OptionSmtpchk)
			if !d.NoOption {
				return "smtpchk"
			}
		}

		data, err = p.Get(section, sectionName, "option ldap-check", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "ldap-check"
			}
		}

		data, err = p.Get(section, sectionName, "option mysql-check", false)
		if err == nil {
			d := data.(*types.OptionMysqlCheck)
			if !d.NoOption {
				return "mysql-check"
			}
		}

		data, err = p.Get(section, sectionName, "option pgsql-check", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "pgsql-check"
			}
		}

		data, err = p.Get(section, sectionName, "option tcp-check", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "tcp-check"
			}
		}

		data, err = p.Get(section, sectionName, "option redis-check", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "redis-check"
			}
		}
	}

	if fieldName == "ExternalCheck" {
		data, err := p.Get(section, sectionName, "option external-check", false)
		if err != nil {
			return nil
		}
		if data.(*types.SimpleOption).NoOption {
			return "disabled"
		}
		return "enabled"
	}
	if fieldName == "ExternalCheckPath" {
		data, err := p.Get(section, sectionName, "external-check path", false)
		if err != nil {
			return nil
		}
		d := data.(*types.ExternalCheckPath)
		return d.Path
	}
	if fieldName == "ExternalCheckCommand" {
		data, err := p.Get(section, sectionName, "external-check command", false)
		if err != nil {
			return nil
		}
		d := data.(*types.ExternalCheckCommand)
		return d.Command
	}
	if fieldName == "DefaultBackend" {
		data, err := p.Get(section, sectionName, "default_backend", false)
		if err != nil {
			return nil
		}
		bck := data.(*types.StringC)
		return bck.Value
	}
	if fieldName == "Clflog" {
		data, err := p.Get(section, sectionName, "option httplog", false)
		if err == nil {
			d := data.(*types.OptionHTTPLog)
			if !d.NoOption {
				return d.Clf
			}
		}
		return nil
	}
	if fieldName == "Httplog" {
		data, err := p.Get(section, sectionName, "option httplog", false)
		if err == nil {
			d := data.(*types.OptionHTTPLog)
			if !d.NoOption {
				return !d.Clf
			}
		}
		return nil
	}
	if fieldName == "HTTPConnectionMode" {
		data, err := p.Get(section, sectionName, "option http-tunnel", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "http-tunnel"
			}
		}

		data, err = p.Get(section, sectionName, "option httpclose", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "httpclose"
			}
		}
		// deprecated option, alias for httpclose
		data, err = p.Get(section, sectionName, "option forceclose", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "httpclose"
			}
		}

		data, err = p.Get(section, sectionName, "option http-server-close", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "http-server-close"
			}
		}

		data, err = p.Get(section, sectionName, "option http-keep-alive", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "http-keep-alive"
			}
		}
	}
	//Check Timeouts
	if strings.HasSuffix(fieldName, "Timeout") {
		if pName := translateTimeout(fieldName); p.HasParser(section, pName) {
			data, err := p.Get(section, sectionName, pName, false)
			if err != nil {
				return nil
			}
			timeout := data.(*types.SimpleTimeout)
			return misc.ParseTimeout(timeout.Value)
		}
	}
	//Check single line
	if pName := misc.DashCase(fieldName); p.HasParser(section, pName) {
		data, err := p.Get(section, sectionName, pName, false)
		if err != nil {
			return nil
		}
		return parseOption(data)
	}
	//Check options
	if pName := fmt.Sprintf("option %s", misc.DashCase(fieldName)); p.HasParser(section, pName) {
		data, err := p.Get(section, sectionName, pName, false)
		if err != nil {
			return nil
		}
		return parseOption(data)
	}
	return nil
}

func (c *Client) createEditSection(object interface{}, section parser.Section, pName string, p *parser.Parser) error {
	objValue := reflect.ValueOf(object).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		typeField := objValue.Type().Field(i)
		field := objValue.FieldByName(typeField.Name)
		if typeField.Name != "Name" && typeField.Name != "ID" {
			if err := c.setFieldValue(section, pName, typeField.Name, field, p); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Client) setFieldValue(section parser.Section, sectionName string, fieldName string, field reflect.Value, p *parser.Parser) error {
	//Handle special cases
	if fieldName == "Httpchk" {
		if section == parser.Backends || section == parser.Defaults {
			if valueIsNil(field) {
				if err := p.Set(section, sectionName, "option httpchk", nil); err != nil {
					return err
				}
				return nil
			}
			hc := field.Elem().Interface().(models.Httpchk)
			d := &types.OptionHttpchk{
				Method:  hc.Method,
				Version: hc.Version,
				URI:     hc.URI,
			}
			if err := p.Set(section, sectionName, "option httpchk", d); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "Forwardfor" {
		if valueIsNil(field) {
			if err := p.Set(section, sectionName, "option forwardfor", nil); err != nil {
				return err
			}
			return nil
		}
		ff := field.Elem().Interface().(models.Forwardfor)
		d := &types.OptionForwardFor{
			Except: ff.Except,
			Header: ff.Header,
			IfNone: ff.Ifnone,
		}
		if err := p.Set(section, sectionName, "option forwardfor", d); err != nil {
			return err
		}
		return nil
	}
	if fieldName == "Redispatch" {
		if section == parser.Backends || section == parser.Defaults {
			if valueIsNil(field) {
				if err := p.Set(section, sectionName, "option redispatch", nil); err != nil {
					return err
				}
				return nil
			}
			br := field.Elem().Interface().(models.Redispatch)
			d := &types.OptionRedispatch{
				Interval: &br.Interval,
				NoOption: false,
			}
			if *br.Enabled == "disabled" {
				d.NoOption = true
			}
			if err := p.Set(section, sectionName, "option redispatch", d); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "Balance" {
		if section == parser.Backends || section == parser.Defaults {
			if valueIsNil(field) {
				if err := p.Set(section, sectionName, "balance", nil); err != nil {
					return err
				}
				return nil
			}
			b := field.Elem().Interface().(models.Balance)
			d := types.Balance{
				Algorithm: b.Algorithm,
				Arguments: b.Arguments,
			}
			if err := p.Set(section, sectionName, "balance", &d); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "HashType" {
		if section == parser.Backends {
			if valueIsNil(field) {
				if err := p.Set(section, sectionName, "hash-type", nil); err != nil {
					return err
				}
				return nil
			}
			b := field.Elem().Interface().(models.BackendHashType)
			d := types.HashType{
				Method:   b.Method,
				Function: b.Function,
				Modifier: b.Modifier,
			}
			if err := p.Set(section, sectionName, "hash-type", &d); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "ErrorFiles" {
		if section == parser.Defaults {
			if valueIsNil(field) {
				if err := p.Set(section, sectionName, "errorfile", nil); err != nil {
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
			if err := p.Set(section, sectionName, "errorfile", errorFiles); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "DefaultServer" {
		if section == parser.Backends || section == parser.Defaults {
			if valueIsNil(field) {
				if err := p.Set(section, sectionName, "default-server", nil); err != nil {
					return err
				}
				return nil
			}
			ds := field.Elem().Interface().(models.DefaultServer)
			dServers := []types.DefaultServer{types.DefaultServer{}}

			ps := make([]params.ServerOption, 0, 4)
			if ds.Fall != nil {
				param := &params.ServerOptionValue{
					Name:  "fall",
					Value: strconv.FormatInt(*ds.Fall, 10),
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

			if ds.Port != nil {
				param := &params.ServerOptionValue{
					Name:  "port",
					Value: strconv.FormatInt(*ds.Port, 10),
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

			dServers[0].Params = ps
			if err := p.Set(section, sectionName, "default-server", dServers); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "StickTable" {
		if section == parser.Backends {
			if valueIsNil(field) {
				if err := p.Set(section, sectionName, "stick-table", nil); err != nil {
					return err
				}
				return nil
			}
			st := field.Elem().Interface().(models.BackendStickTable)
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
			if err := p.Set(section, sectionName, "stick-table", d); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "AdvCheck" {
		if section == parser.Backends || section == parser.Defaults {
			if err := p.Set(section, sectionName, "option ssl-hello-chk", nil); err != nil {
				return err
			}
			if err := p.Set(section, sectionName, "option smtpchk", nil); err != nil {
				return err
			}
			if err := p.Set(section, sectionName, "option ldap-check", nil); err != nil {
				return err
			}
			if err := p.Set(section, sectionName, "option mysql-check", nil); err != nil {
				return err
			}
			if err := p.Set(section, sectionName, "option pgsql-check", nil); err != nil {
				return err
			}
			if err := p.Set(section, sectionName, "option tcp-check", nil); err != nil {
				return err
			}
			if err := p.Set(section, sectionName, "option redis-check", nil); err != nil {
				return err
			}

			if !valueIsNil(field) {
				var d common.ParserData
				pName := fmt.Sprintf("option %v", field.String())
				if pName == "option smtpchk" {
					d = &types.OptionSmtpchk{
						NoOption: false,
					}
				} else if pName == "option mysql-check" {
					d = &types.OptionMysqlCheck{
						NoOption: false,
					}
				} else {
					d = &types.SimpleOption{
						NoOption: false,
					}
				}
				if err := p.Set(section, sectionName, pName, d); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if fieldName == "ExternalCheck" {
		if section == parser.Backends || section == parser.Defaults {
			pExternalCheck := &types.SimpleOption{}
			if valueIsNil(field) {
				pExternalCheck = nil
			} else if field.String() == "disabled" {
				pExternalCheck.NoOption = true
			}
			if err := p.Set(section, sectionName, "option external-check", pExternalCheck); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "ExternalCheckPath" {
		if section == parser.Backends || section == parser.Defaults {
			pExtPath := &types.ExternalCheckPath{}
			if valueIsNil(field) {
				pExtPath = nil
			} else {
				pExtPath.Path = field.Interface().(string)
			}
			if err := p.Set(section, sectionName, "external-check path", pExtPath); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "ExternalCheckCommand" {
		if section == parser.Backends || section == parser.Defaults {
			pExtCmd := &types.ExternalCheckCommand{}
			if valueIsNil(field) {
				pExtCmd = nil
			} else {
				pExtCmd.Command = field.Interface().(string)
			}
			if err := p.Set(section, sectionName, "external-check command", pExtCmd); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "DefaultBackend" {
		if section == parser.Frontends || section == parser.Defaults {
			if valueIsNil(field) {
				if err := p.Set(section, sectionName, "default_backend", nil); err != nil {
					return err
				}
				return nil
			}
			bck := field.String()
			d := &types.StringC{
				Value: bck,
			}
			if err := p.Set(section, sectionName, "default_backend", d); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "HTTPConnectionMode" {
		if err := p.Set(section, sectionName, "option httpclose", nil); err != nil {
			return err
		}
		if err := p.Set(section, sectionName, "option http-server-close", nil); err != nil {
			return err
		}
		if err := p.Set(section, sectionName, "option http-keep-alive", nil); err != nil {
			return err
		}
		//Deprecated, delete if exists
		p.Set(section, sectionName, "option forceclose", nil)

		if !valueIsNil(field) {
			pName := fmt.Sprintf("option %v", field.String())
			d := &types.SimpleOption{
				NoOption: false,
			}
			if err := p.Set(section, sectionName, pName, d); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "Clflog" {
		if section == parser.Frontends || section == parser.Defaults {
			if valueIsNil(field) {
				// check if httplog exists, if not do nothing
				d, err := p.Get(section, sectionName, "option httplog", false)
				if err != nil {
					if err != parser_errors.ErrFetch {
						return err
					}
					return nil
				}
				o := d.(*types.OptionHTTPLog)
				if o.Clf {
					o.Clf = false
					if err := p.Set(section, sectionName, "option httplog", o); err != nil {
						return err
					}
				}
				return nil
			}
			o := &types.OptionHTTPLog{Clf: true}
			if err := p.Set(section, sectionName, "option httplog", o); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "Httplog" {
		if section == parser.Frontends || section == parser.Defaults {
			if valueIsNil(field) {
				// check if clflog is active, if yes, do nothing
				d, err := p.Get(section, sectionName, "option httplog", false)
				if err != nil {
					if err != parser_errors.ErrFetch {
						return err
					}
					return nil
				}
				o := d.(*types.OptionHTTPLog)
				if !o.Clf {
					if err := p.Set(section, sectionName, "option httplog", nil); err != nil {
						return err
					}
				}
				return nil
			}
			o := &types.OptionHTTPLog{}
			if err := p.Set(section, sectionName, "option httplog", o); err != nil {
				return err
			}
		}
		return nil
	}
	// check timeouts
	if strings.HasSuffix(fieldName, "Timeout") {
		if pName := translateTimeout(fieldName); p.HasParser(section, pName) {
			if valueIsNil(field) {
				if err := p.Set(section, sectionName, pName, nil); err != nil {
					return err
				}
				return nil
			}
			t := &types.SimpleTimeout{}
			t.Value = strconv.FormatInt(field.Elem().Int(), 10)
			if err := p.Set(section, sectionName, pName, t); err != nil {
				return err
			}
		}
		return nil
	}
	//Check options
	if pName := fmt.Sprintf("option %s", misc.DashCase(fieldName)); p.HasParser(section, pName) {
		if valueIsNil(field) {
			if err := p.Set(section, sectionName, pName, nil); err != nil {
				return err
			}
			return nil
		}
		o := &types.SimpleOption{}
		if field.Kind() == reflect.String {
			if field.String() == "disabled" {
				o.NoOption = true
			}
		}
		if err := p.Set(section, sectionName, pName, o); err != nil {
			return err
		}
		return nil
	}
	//Check single line
	if pName := misc.DashCase(fieldName); p.HasParser(section, pName) {
		if valueIsNil(field) {
			if err := p.Set(section, sectionName, pName, nil); err != nil {
				return err
			}
			return nil
		}
		d := translateToParserData(field)
		if d == nil {
			return errors.Errorf("Cannot parse type for %s %s: %s", section, sectionName, fieldName)
		}
		if err := p.Set(section, sectionName, pName, d); err != nil {
			return err
		}
		return nil
	}

	return errors.Errorf("Cannot parse option for %s %s: %s", section, sectionName, fieldName)
}

func (c *Client) handleError(id, parentType, parentName, transaction string, implicit bool, err error) error {
	var e error
	if err == parser_errors.ErrSectionMissing {
		if parentName != "" {
			e = NewConfError(ErrParentDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		} else {
			e = NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Object %s does not exist", id))
		}
	} else if err == parser_errors.ErrSectionAlreadyExists {
		e = NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("Object %s already exists", id))
	} else if err == parser_errors.ErrFetch {
		e = NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Object %v does not exist in %s %s", id, parentType, parentName))
	} else if err == parser_errors.ErrIndexOutOfRange {
		e = NewConfError(ErrObjectIndexOutOfRange, fmt.Sprintf("Object with id %v in %s %s out of range", id, parentType, parentName))
	} else {
		e = err
	}

	if implicit {
		return c.errAndDeleteTransaction(e, transaction)
	}
	return e
}

func (c *Client) errAndDeleteTransaction(err error, tID string) error {
	// Just a safety to not delete the master files by mistake
	if tID != "" {
		c.DeleteTransaction(tID)
		return err
	}
	return err
}

func (c *Client) deleteSection(section parser.Section, name string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !c.checkSectionExists(section, name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", section, name))
		return c.handleError(name, "", "", t, transactionID == "", e)
	}

	if err := p.SectionsDelete(section, name); err != nil {
		return c.handleError(name, "", "", t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func (c *Client) editSection(section parser.Section, name string, data interface{}, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !c.checkSectionExists(section, name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", section, name))
		return c.handleError(name, "", "", t, transactionID == "", e)
	}

	if err := c.createEditSection(data, section, name, p); err != nil {
		return c.handleError(name, "", "", t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func (c *Client) createSection(section parser.Section, name string, data interface{}, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if c.checkSectionExists(section, name, p) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", section, name))
		return c.handleError(name, "", "", t, transactionID == "", e)
	}

	if err := p.SectionsCreate(section, name); err != nil {
		return c.handleError(name, "", "", t, transactionID == "", err)
	}

	if err := c.createEditSection(data, section, name, p); err != nil {
		return c.handleError(name, "", "", t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func (c *Client) checkSectionExists(section parser.Section, sectionName string, p *parser.Parser) bool {
	sections, err := p.SectionsGet(section)
	if err != nil {
		return false
	}

	if misc.StringInSlice(sectionName, sections) {
		return true
	}
	return false
}

func (c *Client) loadDataForChange(transactionID string, version int64) (*parser.Parser, string, error) {
	t, err := c.checkTransactionOrVersion(transactionID, version)
	if err != nil {
		// if transaction is implicit, return err and delete transaction
		if transactionID == "" && t != "" {
			return nil, "", c.errAndDeleteTransaction(err, t)
		}
		return nil, "", err
	}

	p, err := c.GetParser(t)
	if err != nil {
		if transactionID == "" && t != "" {
			return nil, "", c.errAndDeleteTransaction(err, t)
		}
		return nil, "", err
	}
	return p, t, nil
}

func (c *Client) saveData(p *parser.Parser, t string, commitImplicit bool) error {
	if c.PersistentTransactions {
		tFile, err := c.getTransactionFile(t)
		if err != nil {
			return err
		}

		if err := p.Save(tFile); err != nil {
			e := NewConfError(ErrErrorChangingConfig, err.Error())
			if commitImplicit {
				return c.errAndDeleteTransaction(e, t)
			}
			return err
		}
	}

	if commitImplicit {
		if _, err := c.CommitTransaction(t); err != nil {
			return err
		}
	}
	return nil
}

func valueIsNil(v reflect.Value) bool {
	switch v.Kind() {
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
	switch field.Kind() {
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
