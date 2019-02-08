package configuration

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/configuration/cache"
	parser "github.com/haproxytech/config-parser"
)

const (
	//DefaultConfigurationFile sane default for path to haproxy configuration file
	DefaultConfigurationFile string = "/etc/haproxy/haproxy.cfg"
	//DefaultGlobalConfigurationFile sane default for path to haproxy global configuration file
	DefaultGlobalConfigurationFile string = "/etc/haproxy/global.cfg"
	//DefaultHaproxy sane default for path to haproxy executable
	DefaultHaproxy string = "/usr/sbin/haproxy"
	//DefaultUseValidation sane default using validation in client native
	DefaultUseValidation bool = true
	//DefaultUseCache sane default using caching in client native
	DefaultUseCache bool = false
	// DefaultLBCTLPath sane default for path to lbctl
	DefaultLBCTLPath string = "/usr/sbin/lbctl"
	// DefaultTransactionDir sane default for path for transactions
	DefaultTransactionDir string = "/tmp/haproxy"
)

// ClientParams is just a placeholder for all client options
type ClientParams struct {
	ConfigurationFile       string
	GlobalConfigurationFile string
	Haproxy                 string
	UseValidation           bool
	UseCache                bool
	LBCTLPath               string
	TransactionDir          string
}

// Client configuration client
type Client struct {
	ClientParams
	cache.Cache
	GlobalParser parser.Parser
}

// DefaultClient returns Client with sane defaults
func DefaultClient() (*Client, error) {
	p := ClientParams{
		ConfigurationFile:       DefaultConfigurationFile,
		GlobalConfigurationFile: DefaultGlobalConfigurationFile,
		Haproxy:                 DefaultHaproxy,
		UseValidation:           DefaultUseValidation,
		UseCache:                DefaultUseCache,
		LBCTLPath:               DefaultLBCTLPath,
		TransactionDir:          DefaultTransactionDir,
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
	if options.LBCTLPath == "" {
		options.LBCTLPath = DefaultLBCTLPath
	}

	if options.TransactionDir == "" {
		options.TransactionDir = DefaultTransactionDir
	}

	if options.ConfigurationFile == "" {
		options.ConfigurationFile = DefaultConfigurationFile
	}

	if options.GlobalConfigurationFile == "" {
		options.GlobalConfigurationFile = DefaultGlobalConfigurationFile
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

	err := c.GlobalParser.LoadData(c.GlobalConfigurationFile)
	if err != nil {
		return err
	}

	return nil
}

// GetGlobalVersion returns global configuration file version
func (c *Client) GetGlobalVersion() (int64, error) {
	return c.getVersion("global", "")
}

// GetVersion returns configuration file version
func (c *Client) GetVersion(transaction string) (int64, error) {
	v, err := c.getVersion("config", transaction)
	if err == nil && c.Cache.Enabled() {
		c.Cache.Version.Set(v, transaction)
	}
	return v, err
}

func (c *Client) getVersion(t string, transaction string) (int64, error) {
	var file *os.File
	var err error
	if t == "global" {
		file, err = os.Open(c.GlobalConfigurationFile)
	} else {
		if transaction == "" {
			file, err = os.Open(c.ConfigurationFile)
		} else {
			file, err = os.Open(c.getTransactionFile(c.ConfigurationFile, transaction))
			if err != nil && os.IsNotExist(err) {
				file, err = os.Open(c.getFailedTransactionFile(c.ConfigurationFile, transaction))
			}
		}
	}

	if err != nil {
		return 0, NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read configuration file %v: %v", c.ConfigurationFile, err.Error()))
	}
	defer file.Close()

	return c.readVersionFromFile(file)
}

func (c *Client) readVersionFromFile(file *os.File) (int64, error) {
	scanner := bufio.NewScanner(file)
	// Read only first line, version MUST BE on the first line
	lineNo := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			if lineNo == 0 {
				if strings.HasPrefix(line, "# _version=") {
					w := strings.Split(line, "=")
					if len(w) != 2 {
						return c.setInitialVersion(true, file)
					}
					version, err := strconv.ParseInt(w[1], 10, 64)
					if err != nil {
						return c.setInitialVersion(true, file)
					}
					return version, nil
				}
			} else {
				break
			}
			lineNo++
		}
	}
	if err := scanner.Err(); err != nil {
		return c.setInitialVersion(false, file)
	}
	return c.setInitialVersion(false, file)
}

func (c *Client) setInitialVersion(hasVersion bool, file *os.File) (int64, error) {
	var input []byte
	_, err := file.Read(input)
	if err != nil {
		return 0, NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read configuration file %v: %v", file.Name(), err.Error()))
	}

	inputStr := string(input)

	if hasVersion {
		inputStrArr := strings.SplitAfterN(inputStr, "\n", 2)
		if len(inputStrArr) == 2 {
			inputStr = inputStrArr[1]
		}
	}

	output := fmt.Sprintf("# _version=1\n%s", inputStr)
	if _, err = file.Write([]byte(output)); err != nil {
		return 0, NewConfError(ErrCannotSetVersion, fmt.Sprintf("Cannot set initial version in file %v: %v", file.Name(), err.Error()))
	}
	return 0, nil
}

func (c *Client) incrementVersion() error {
	input, err := ioutil.ReadFile(c.ConfigurationFile)

	if err != nil {
		return NewConfError(ErrCannotReadVersion, fmt.Sprintf("Cannot read version from file %v: %v", c.ConfigurationFile, err.Error()))
	}

	v, err := c.GetVersion("")
	if err != nil {
		return err
	}

	toReplace := fmt.Sprintf("# _version=%v", v)
	replace := fmt.Sprintf("# _version=%v", v+1)

	output := bytes.Replace(input, []byte(toReplace), []byte(replace), -1)

	if err = ioutil.WriteFile(c.ConfigurationFile, output, 0666); err != nil {
		return NewConfError(ErrCannotSetVersion, fmt.Sprintf("Cannot increment version in file %v: %v", c.ConfigurationFile, err.Error()))
	}
	if c.Cache.Enabled() {
		c.Cache.Version.Set(v+1, "")
	}
	return nil
}

func (c *Client) incrementGlobalVersion() error {
	input, err := ioutil.ReadFile(c.GlobalConfigurationFile)

	if err != nil {
		return NewConfError(ErrCannotReadVersion, fmt.Sprintf("Cannot read version from file %v: %v", c.GlobalConfigurationFile, err.Error()))
	}

	v, err := c.GetGlobalVersion()
	if err != nil {
		return err
	}

	toReplace := fmt.Sprintf("# _version=%v", v)
	replace := fmt.Sprintf("# _version=%v", v+1)

	output := bytes.Replace(input, []byte(toReplace), []byte(replace), -1)

	if err = ioutil.WriteFile(c.GlobalConfigurationFile, output, 0666); err != nil {
		return NewConfError(ErrCannotSetVersion, fmt.Sprintf("Cannot increment version in file %v: %v", c.GlobalConfigurationFile, err.Error()))
	}
	return nil
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
			transaction, err := c.StartTransaction(version)
			if err != nil {
				return "", err
			}
			t = transaction.ID
		}
	}
	return t, nil
}
