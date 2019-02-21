package configuration

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/haproxytech/config-parser/common"

	"github.com/haproxytech/client-native/configuration/cache"
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
	//DefaultUseCache sane default using caching in client native
	DefaultUseCache bool = false
	// DefaultTransactionDir sane default for path for transactions
	DefaultTransactionDir string = "/tmp/haproxy"
)

// ClientParams is just a placeholder for all client options
type ClientParams struct {
	ConfigurationFile string
	Haproxy           string
	UseValidation     bool
	UseCache          bool
	TransactionDir    string
}

// Client configuration client
type Client struct {
	ClientParams
	cache.Cache
	ConfigParser parser.Parser
}

// DefaultClient returns Client with sane defaults
func DefaultClient() (*Client, error) {
	p := ClientParams{
		ConfigurationFile: DefaultConfigurationFile,
		Haproxy:           DefaultHaproxy,
		UseValidation:     DefaultUseValidation,
		UseCache:          DefaultUseCache,
		TransactionDir:    DefaultTransactionDir,
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

	c.Cache = cache.Cache{}

	if c.UseCache {
		v, err := c.GetVersion("")
		if err != nil {
			return err
		}
		c.Cache.Init(v)
	}

	return nil
}

// GetVersion returns configuration file version
func (c *Client) GetVersion(transaction string) (int64, error) {
	if c.Cache.Enabled() {
		v := c.Cache.Version.Get(transaction)
		if v != 0 {
			return v, nil
		}
	}
	v, err := c.getVersion(transaction)
	if err == nil && c.Cache.Enabled() {
		c.Cache.Version.Set(v, transaction)
	}
	return v, err
}

func (c *Client) getVersion(transaction string) (int64, error) {
	var configFile string
	if transaction == "" {
		configFile = c.ConfigurationFile
	} else {
		configFile = c.getTransactionFile(transaction)
	}

	err := c.ConfigParser.LoadData(configFile)
	if err != nil {
		return 0, err
	}
	data, _ := c.ConfigParser.Get(parser.Comments, parser.CommentsSectionName, "# _version", true)
	ver, _ := data.(*types.ConfigVersion)
	return ver.Value, nil
}

func (c *Client) incrementVersion() error {
	err := c.ConfigParser.LoadData(c.ConfigurationFile)
	if err != nil {
		return err
	}
	data, _ := c.ConfigParser.Get(parser.Comments, parser.CommentsSectionName, "# _version", true)
	ver, _ := data.(*types.ConfigVersion)
	ver.Value = ver.Value + 1

	if c.Cache.Enabled() {
		c.Version.Set(ver.Value, "")
	}
	return c.ConfigParser.Save(c.ConfigurationFile)
}

func (c *Client) checkTransactionOrVersion(transactionID string, version int64, startTransaction bool) (string, error) {
	// start an implicit transaction for delete site (multiple operations required) if not already given
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
		if startTransaction {
			transaction, err := c.startTransaction(version, false)
			if err != nil {
				return "", err
			}
			t = transaction.ID
		}
	}
	return t, nil
}

func (c *Client) parseSection(object interface{}, p parser.Section, pName string) error {
	objValue := reflect.ValueOf(object).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		typeField := objValue.Type().Field(i)
		field := objValue.FieldByName(typeField.Name)
		if typeField.Name != "Name" && typeField.Name != "ID" {
			val := c.parseField(p, pName, typeField.Name)
			if val != nil {
				if field.Kind() == reflect.Bool {
					if val == "enabled" {
						field.Set(reflect.ValueOf(true))
					} else {
						field.Set(reflect.ValueOf(false))
					}
				} else {
					field.Set(reflect.ValueOf(val))
				}
			}
		}
	}

	return nil
}

func (c *Client) parseField(p parser.Section, sectionName string, fieldName string) interface{} {
	//Handle special cases
	if fieldName == "Httpchk" {
		data, err := c.ConfigParser.Get(p, sectionName, "option httpchk", false)
		if err != nil {
			return nil
		}
		d := data.(*types.OptionHttpchk)
		if p == parser.Backends {
			return &models.BackendHttpchk{
				Method:  d.Method,
				URI:     d.Uri,
				Version: d.Version,
			}
		}
	}
	if fieldName == "Balance" {
		data, err := c.ConfigParser.Get(p, sectionName, "balance", false)
		if err != nil {
			return nil
		}
		d := data.(*types.Balance)
		return &models.BackendBalance{
			Algorithm: d.Algorithm,
			Arguments: d.Arguments,
		}
	}
	if fieldName == "DefaultServer" {
		data, err := c.ConfigParser.Get(p, sectionName, "default-server", false)
		if err != nil {
			return nil
		}
		d := data.([]types.DefaultServer)
		if p == parser.Backends {
			dServer := &models.BackendDefaultServer{}
			for _, ds := range d {
				dsParams := ds.Params
				for _, p := range dsParams {
					v, ok := p.(*params.ServerOptionValue)
					if ok {
						switch v.Name {
						case "fall":
							dServer.Fall = parseTimeout(v.Value)
						case "inter":
							dServer.Inter = parseTimeout(v.Value)
						case "rise":
							dServer.Rise = parseTimeout(v.Value)
						case "port":
							p, err := strconv.ParseInt(v.Value, 10, 64)
							if err == nil {
								dServer.Port = &p
							}
						}
					}
				}
			}
			return dServer
		}
		return nil
	}
	if fieldName == "StickTable" {
		data, err := c.ConfigParser.Get(p, sectionName, "stick-table", false)
		if err != nil {
			return nil
		}
		d := data.(*types.StickTable)
		if p == parser.Backends {
			st := &models.BackendStickTable{
				Type:   d.Type,
				Size:   parseSize(d.Size),
				Store:  d.Store,
				Expire: parseTimeout(d.Expire),
				Peers:  d.Peers,
			}
			k, err := strconv.ParseInt(d.Length, 10, 64)
			if err == nil {
				st.Keylen = &k
			}
			if d.NoPurge {
				st.Nopurge = true
			}
		}
		return nil
	}
	if fieldName == "AdvCheck" {
		data, err := c.ConfigParser.Get(p, sectionName, "option ssl-hello-check", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "ssl-hello-check"
			}
		}

		data, err = c.ConfigParser.Get(p, sectionName, "option smtpchk", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "smtpchk"
			}
		}

		data, err = c.ConfigParser.Get(p, sectionName, "option ldap-check", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "ldap-check"
			}
		}

		data, err = c.ConfigParser.Get(p, sectionName, "option mysql-check", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "mysql-check"
			}
		}

		data, err = c.ConfigParser.Get(p, sectionName, "option pgsql-check", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "pgsql-check"
			}
		}

		data, err = c.ConfigParser.Get(p, sectionName, "option tcp-check", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "tcp-check"
			}
		}

		data, err = c.ConfigParser.Get(p, sectionName, "option redis-check", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "redis-check"
			}
		}
	}
	if fieldName == "DefaultBackend" {
		data, err := c.ConfigParser.Get(p, sectionName, "default_backend", false)
		if err != nil {
			return nil
		}
		bck := data.(*types.StringC)
		return bck.Value
	}
	if fieldName == "Clflog" {
		data, err := c.ConfigParser.Get(p, sectionName, "option httplog", false)
		if err == nil {
			d := data.(*types.OptionHTTPLog)
			if !d.NoOption {
				return d.Clf
			}
		}
		return nil
	}
	if fieldName == "Httplog" {
		data, err := c.ConfigParser.Get(p, sectionName, "option httplog", false)
		if err == nil {
			d := data.(*types.OptionHTTPLog)
			if !d.NoOption {
				return !d.Clf
			}
		}
		return nil
	}
	if fieldName == "HTTPConnectionMode" {
		data, err := c.ConfigParser.Get(p, sectionName, "option http-tunnel", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "http-tunnel"
			}
		}

		data, err = c.ConfigParser.Get(p, sectionName, "option httpclose", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "httpclose"
			}
		}

		data, err = c.ConfigParser.Get(p, sectionName, "option forceclose", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "force-close"
			}
		}

		data, err = c.ConfigParser.Get(p, sectionName, "option http-server-close", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "http-server-close"
			}
		}

		data, err = c.ConfigParser.Get(p, sectionName, "option http-keep-alive", false)
		if err == nil {
			d := data.(*types.SimpleOption)
			if !d.NoOption {
				return "http-keep-alive"
			}
		}
	}
	//Check Timeouts
	if strings.HasSuffix(fieldName, "Timeout") {
		if pName := translateTimeout(fieldName); c.ConfigParser.HasParser(p, pName) {
			data, err := c.ConfigParser.Get(p, sectionName, pName, false)
			if err != nil {
				return nil
			}
			timeout := data.(*types.SimpleTimeout)
			return parseTimeout(timeout.Value)
		}
	}
	//Check single line
	if pName := misc.DashCase(fieldName); c.ConfigParser.HasParser(p, pName) {
		data, err := c.ConfigParser.Get(p, sectionName, pName, false)
		if err != nil {
			return nil
		}
		return parseOption(data)
	}
	//Check options
	if pName := fmt.Sprintf("option %s", misc.DashCase(fieldName)); c.ConfigParser.HasParser(p, pName) {
		data, err := c.ConfigParser.Get(p, sectionName, pName, false)
		if err != nil {
			return nil
		}
		return parseOption(data)
	}
	return nil
}

func (c *Client) createEditSection(object interface{}, p parser.Section, pName string) error {
	objValue := reflect.ValueOf(object).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		typeField := objValue.Type().Field(i)
		field := objValue.FieldByName(typeField.Name)
		if typeField.Name != "Name" && typeField.Name != "ID" {
			if err := c.setFieldValue(p, pName, typeField.Name, field); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Client) setFieldValue(p parser.Section, sectionName string, fieldName string, field reflect.Value) error {
	//Handle special cases
	if fieldName == "Httpchk" {
		if valueIsNil(field) {
			if err := c.ConfigParser.Set(p, sectionName, "option httpchk", nil); err != nil {
				return err
			}
			return nil
		}
		hc := field.Elem().Interface().(models.BackendHttpchk)
		d := &types.OptionHttpchk{
			Method:  hc.Method,
			Version: hc.Version,
			Uri:     hc.URI,
		}
		if err := c.ConfigParser.Set(p, sectionName, "option httpchk", d); err != nil {
			return err
		}
		return nil
	}
	if fieldName == "Balance" {
		if valueIsNil(field) {
			if err := c.ConfigParser.Set(p, sectionName, "balance", nil); err != nil {
				return err
			}
			return nil
		}
		b := field.Elem().Interface().(models.BackendBalance)
		d := types.Balance{
			Algorithm: b.Algorithm,
			Arguments: b.Arguments,
		}
		if err := c.ConfigParser.Set(p, sectionName, "balance", &d); err != nil {
			return err
		}
		return nil
	}
	if fieldName == "DefaultServer" {
		return nil
	}
	if fieldName == "StickTable" {
		if valueIsNil(field) {
			if err := c.ConfigParser.Set(p, sectionName, "stick-table", nil); err != nil {
				return err
			}
			return nil
		}
		st := field.Elem().Interface().(models.BackendStickTable)
		d := &types.StickTable{
			Type:   st.Type,
			Size:   strconv.FormatInt(*st.Size, 10),
			Store:  st.Store,
			Expire: strconv.FormatInt(*st.Expire, 10),
			Peers:  st.Peers,
		}
		if err := c.ConfigParser.Set(p, sectionName, "stick-table", d); err != nil {
			return err
		}
		return nil
	}
	if fieldName == "AdvCheck" {
		if err := c.ConfigParser.Set(p, sectionName, "option ssl-hello-check", nil); err != nil {
			return err
		}
		if err := c.ConfigParser.Set(p, sectionName, "option smtpchk", nil); err != nil {
			return err
		}
		if err := c.ConfigParser.Set(p, sectionName, "option ldap-check", nil); err != nil {
			return err
		}
		if err := c.ConfigParser.Set(p, sectionName, "option mysql-check", nil); err != nil {
			return err
		}
		if err := c.ConfigParser.Set(p, sectionName, "option pgsql-check", nil); err != nil {
			return err
		}
		if err := c.ConfigParser.Set(p, sectionName, "option tcp-check", nil); err != nil {
			return err
		}
		if err := c.ConfigParser.Set(p, sectionName, "option redis-check", nil); err != nil {
			return err
		}

		if !valueIsNil(field) {
			pName := fmt.Sprintf("option %v", field.String())
			d := &types.SimpleOption{
				NoOption: false,
			}
			if err := c.ConfigParser.Set(p, sectionName, pName, d); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "DefaultBackend" {
		if valueIsNil(field) {
			if err := c.ConfigParser.Set(p, sectionName, "default_backend", nil); err != nil {
				return err
			}
			return nil
		}
		bck := field.String()
		d := &types.StringC{
			Value: bck,
		}
		if err := c.ConfigParser.Set(p, sectionName, "default_backend", d); err != nil {
			return err
		}
		return nil
	}
	if fieldName == "HTTPConnectionMode" {
		if err := c.ConfigParser.Set(p, sectionName, "option http-tunnel", nil); err != nil {
			return err
		}
		if err := c.ConfigParser.Set(p, sectionName, "option httpclose", nil); err != nil {
			return err
		}
		if err := c.ConfigParser.Set(p, sectionName, "option forceclose", nil); err != nil {
			return err
		}
		if err := c.ConfigParser.Set(p, sectionName, "option http-server-close", nil); err != nil {
			return err
		}
		if err := c.ConfigParser.Set(p, sectionName, "option http-keep-alive", nil); err != nil {
			return err
		}

		if !valueIsNil(field) {
			pName := fmt.Sprintf("option %v", field.String())
			d := &types.SimpleOption{
				NoOption: false,
			}
			if err := c.ConfigParser.Set(p, sectionName, pName, d); err != nil {
				return err
			}
		}
		return nil
	}
	if fieldName == "Clflog" {
		if valueIsNil(field) {
			// check if httplog exists, if not do nothing
			d, err := c.ConfigParser.Get(p, sectionName, "option httplog", false)
			if err != nil {
				if err != parser_errors.FetchError {
					return err
				}
				return nil
			}
			o := d.(*types.OptionHTTPLog)
			if o.Clf {
				o.Clf = false
				if err := c.ConfigParser.Set(p, sectionName, "option httplog", o); err != nil {
					return err
				}
			}
			return nil
		}
		o := &types.OptionHTTPLog{Clf: true}
		if err := c.ConfigParser.Set(p, sectionName, "option httplog", o); err != nil {
			return err
		}
		return nil
	}
	if fieldName == "Httplog" {
		if valueIsNil(field) {
			// check if clflog is active, if yes, do nothing
			d, err := c.ConfigParser.Get(p, sectionName, "option httplog", false)
			if err != nil {
				if err != parser_errors.FetchError {
					return err
				}
				return nil
			}
			o := d.(*types.OptionHTTPLog)
			if !o.Clf {
				if err := c.ConfigParser.Set(p, sectionName, "option httplog", nil); err != nil {
					return err
				}
			}
			return nil
		}
		o := &types.OptionHTTPLog{}
		if err := c.ConfigParser.Set(p, sectionName, "option httplog", o); err != nil {
			return err
		}
		return nil
	}
	// check timeouts
	if strings.HasSuffix(fieldName, "Timeout") {
		if pName := translateTimeout(fieldName); c.ConfigParser.HasParser(p, pName) {
			if valueIsNil(field) {
				if err := c.ConfigParser.Set(p, sectionName, pName, nil); err != nil {
					return err
				}
				return nil
			}
			t := &types.SimpleTimeout{}
			t.Value = strconv.FormatInt(field.Elem().Int(), 10)
			if err := c.ConfigParser.Set(p, sectionName, pName, t); err != nil {
				return err
			}
		}
		return nil
	}
	//Check options
	if pName := fmt.Sprintf("option %s", misc.DashCase(fieldName)); c.ConfigParser.HasParser(p, pName) {
		if valueIsNil(field) {
			if err := c.ConfigParser.Set(p, sectionName, pName, nil); err != nil {
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
		if err := c.ConfigParser.Set(p, sectionName, pName, o); err != nil {
			return err
		}
		return nil
	}
	//Check single line
	if pName := misc.DashCase(fieldName); c.ConfigParser.HasParser(p, pName) {
		if valueIsNil(field) {
			if err := c.ConfigParser.Set(p, sectionName, pName, nil); err != nil {
				return err
			}
			return nil
		}
		d := translateToParserData(field)
		if d == nil {
			return errors.Errorf("Cannot parse type for %s %s: %s", p, sectionName, fieldName)
		}
		if err := c.ConfigParser.Set(p, sectionName, pName, d); err != nil {
			return err
		}
		return nil
	}

	return errors.Errorf("Cannot parse option for %s %s: %s", p, sectionName, fieldName)
}

func (c *Client) errAndDeleteTransaction(err error, tID string, delete bool) error {
	if delete {
		c.DeleteTransaction(tID)
	}
	return err
}

func (c *Client) deleteSection(p parser.Section, name string, transactionID string, version int64) error {
	t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !c.checkSectionExists(p, name) {
		return c.errAndDeleteTransaction(NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", p, name)),
			t, transactionID == "")
	}

	if c.ConfigParser.SectionsDelete(p, name); err != nil {
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}

	return nil
}

func (c *Client) editSection(p parser.Section, name string, data interface{}, transactionID string, version int64) error {
	t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !c.checkSectionExists(p, name) {
		return c.errAndDeleteTransaction(NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", p, name)),
			t, transactionID == "")
	}

	if err := c.createEditSection(data, p, name); err != nil {
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}

	return nil
}

func (c *Client) createSection(p parser.Section, name string, data interface{}, transactionID string, version int64) error {
	t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if c.checkSectionExists(p, name) {
		return c.errAndDeleteTransaction(NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", p, name)),
			t, transactionID == "")
	}

	if err := c.ConfigParser.SectionsCreate(p, name); err != nil {
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.createEditSection(data, p, name); err != nil {
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if err := c.saveData(t, transactionID); err != nil {
		return err
	}

	return nil
}

func (c *Client) checkSectionExists(p parser.Section, sectionName string) bool {
	sections, err := c.ConfigParser.SectionsGet(p)
	if err != nil {
		return false
	}

	if misc.StringInSlice(sectionName, sections) {
		return true
	}
	return false
}

func (c *Client) loadDataForChange(transactionID string, version int64) (string, error) {
	t, err := c.checkTransactionOrVersion(transactionID, version, true)
	if err != nil {
		err, ok := err.(*ConfError)
		if !ok {
			return "", c.errAndDeleteTransaction(err, t, transactionID == "")
		}
		return "", err
	}

	if err := c.ConfigParser.LoadData(c.getTransactionFile(t)); err != nil {
		return "", c.errAndDeleteTransaction(err, t, transactionID == "")
	}
	return t, nil
}

func (c *Client) saveData(t, transactionID string) error {
	if err := c.ConfigParser.Save(c.getTransactionFile(t)); err != nil {
		return c.errAndDeleteTransaction(err, t, transactionID == "")
	}

	if transactionID == "" {
		if err := c.commitTransaction(t, false); err != nil {
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

func parseSize(size string) *int64 {
	var v int64
	if strings.HasSuffix(size, "k") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(size, "k"), 10, 64)
		v = v * 1024
	} else if strings.HasSuffix(size, "m") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(size, "m"), 10, 64)
		v = v * 1024 * 1024
	} else if strings.HasSuffix(size, "g") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(size, "g"), 10, 64)
		v = v * 1024 * 1024 * 1024
	} else {
		v, _ = strconv.ParseInt(size, 10, 64)
	}
	if v != 0 {
		return &v
	}
	return nil
}

func parseTimeout(tOut string) *int64 {
	var v int64
	if strings.HasSuffix(tOut, "ms") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "ms"), 10, 64)
	} else if strings.HasSuffix(tOut, "s") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "s"), 10, 64)
		v = v * 1000
	} else if strings.HasSuffix(tOut, "m") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "m"), 10, 64)
		v = v * 1000 * 60
	} else if strings.HasSuffix(tOut, "h") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "h"), 10, 64)
		v = v * 1000 * 60 * 60
	} else if strings.HasSuffix(tOut, "d") {
		v, _ = strconv.ParseInt(strings.TrimSuffix(tOut, "d"), 10, 64)
		v = v * 1000 * 60 * 60 * 24
	} else {
		v, _ = strconv.ParseInt(tOut, 10, 64)
	}
	if v != 0 {
		return &v
	}
	return nil
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
