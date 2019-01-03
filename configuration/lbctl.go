package configuration

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/haproxytech/client-native/configuration/cache"

	"github.com/haproxytech/client-native/misc"
	parser "github.com/haproxytech/config-parser"
)

// LBCTLClient configuration.Client implementation using lbctl
type LBCTLClient struct {
	ClientParams
	LBCTLPath    string
	LBCTLTmpPath string
	GlobalParser parser.Parser
}

const (
	// DefaultLBCTLPath sane default for path to lbctl
	DefaultLBCTLPath string = "/usr/sbin/lbctl"
	// DefaultLBCTLTmpPath sane default for path for lbctl transactions
	DefaultLBCTLTmpPath string = "/tmp/lbctl"
)

// DefaultLBCTLClient returns LBCTLClient with sane defaults
func DefaultLBCTLClient() (*LBCTLClient, error) {
	c := &LBCTLClient{}
	err := c.Init("", "", "", true, false, "", "")

	if err != nil {
		return nil, err
	}

	return c, err
}

// Init initializes a LBCTLClient
func (c *LBCTLClient) Init(configurationFile string, globalConfigurationFile string, haproxy string, useValidation bool, useCache bool, LBCTLPath string, LBCTLTmpPath string) error {
	if LBCTLPath == "" {
		LBCTLPath = DefaultLBCTLPath
	}

	if LBCTLTmpPath == "" {
		LBCTLTmpPath = DefaultLBCTLTmpPath
	}

	c.configurationFile = configurationFile
	c.globalConfigurationFile = globalConfigurationFile
	c.haproxy = haproxy
	c.useValidation = useValidation
	c.LBCTLPath = LBCTLPath
	c.LBCTLTmpPath = LBCTLTmpPath

	c.Cache = cache.Cache{}
	v, err := c.GetVersion()
	if err == nil {
		c.Cache.Init(v)
	}

	return nil
}

func (c *LBCTLClient) executeLBCTL(command string, transaction string, args ...string) (string, error) {
	// fmt.Println("executeLBCTL: command:" + command)
	// fmt.Println("executeLBCTL: transaction:" + transaction)
	// fmt.Printf("executeLBCTL: args: %v \n", args)

	var lbctlArgs []string
	if transaction == "" {
		lbctlArgs = []string{"-S", "root", command}
	} else {
		lbctlArgs = []string{"-T", transaction, command}
	}
	lbctlArgs = append(lbctlArgs, args...)

	cmd := exec.Command(c.LBCTLPath, lbctlArgs...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "LBCTL_MODULES=l7")
	cmd.Env = append(cmd.Env, "LBCTL_L7_SVC_APPLY_CMD=true")

	if c.ConfigurationFile() != "" {
		cmd.Env = append(cmd.Env, "LBCTL_L7_HAPROXY_CONFIG="+c.ConfigurationFile())
		if c.Haproxy() != "" {
			cmd.Env = append(cmd.Env, "LBCTL_L7_SVC_CHECK_CMD="+c.Haproxy()+" -f "+c.GlobalConfigurationFile()+" -f "+c.ConfigurationFile()+" -c")
		} else {
			cmd.Env = append(cmd.Env, "LBCTL_L7_SVC_CHECK_CMD=true")
		}
	} else {
		cmd.Env = append(cmd.Env, "LBCTL_L7_SVC_CHECK_CMD=true")
	}

	if c.LBCTLTmpPath != "" {
		cmd.Env = append(cmd.Env, "LBCTL_TRANS_DIR="+c.LBCTLTmpPath)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			errMsg := err.Error()
			if strings.HasPrefix(errMsg, "exit status") {
				words := strings.Split(errMsg, " ")
				sCode := words[len(words)-1]
				code, err := strconv.ParseInt(sCode, 10, 64)
				if err != nil {
					code = ErrGeneralError
				}
				return "", NewLBCTLError(int(code), cmd.Path, command, string(stderr.Bytes()))
			}
		default:
			return "", NewLBCTLError(ErrGeneralError, cmd.Path, command, string(stderr.Bytes()))
		}
	}

	return string(stdout.Bytes()), nil
}

func (c *LBCTLClient) createObject(name string, objType string, parent string, parentType string, data interface{}, skipFields []string, transactionID string, version int64) error {
	var args []string

	t, err := c.checkTransactionOrVersion(transactionID, version, false)
	if err != nil {
		return err
	}

	cmd := "l7-"
	if parentType != "" {
		cmd = cmd + parentType + "-"
		if parent != "" {
			args = append(args, parent)
		} else {
			return NewConfError(ErrNoParentSpecified, fmt.Sprintf("No parent specified when parent type is %v", parentType))
		}
	} else if parent != "" {
		args = append(args, parent)
	}

	cmd = cmd + objType + "-create"
	args = append(args, name)
	args = append(args, c.serializeObject(data, nil, skipFields)...)

	_, err = c.executeLBCTL(cmd, t, args...)
	if err != nil {
		return err
	}

	if t == "" {
		err = c.incrementVersion()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *LBCTLClient) editObject(name string, objType string, parent string, parentType string, data interface{}, ondisk interface{}, skipFields []string, transactionID string, version int64) error {
	var args []string

	t, err := c.checkTransactionOrVersion(transactionID, version, false)
	if err != nil {
		return err
	}

	cmd := "l7-"
	if parentType != "" {
		cmd = cmd + parentType + "-"
		if parent != "" {
			args = append(args, parent)
		} else {
			return NewConfError(ErrNoParentSpecified, fmt.Sprintf("No parent specified when parent type is %v", parentType))
		}
	} else if parent != "" {
		args = append(args, parent)
	}

	cmd = cmd + objType + "-update"
	args = append(args, name)
	args = append(args, c.serializeObject(data, ondisk, skipFields)...)

	_, err = c.executeLBCTL(cmd, t, args...)
	if err != nil {
		return err
	}
	if t == "" {
		err = c.incrementVersion()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *LBCTLClient) deleteObject(name string, objType string, parent string, parentType string, transactionID string, version int64) error {
	args := make([]string, 0, 1)

	t, err := c.checkTransactionOrVersion(transactionID, version, false)
	if err != nil {
		return err
	}

	cmd := "l7-"
	if parentType != "" {
		cmd = cmd + parentType + "-"
		if parent != "" {
			args = append(args, parent)
		} else {
			return NewConfError(ErrNoParentSpecified, fmt.Sprintf("No parent specified when parent type is %v", parentType))
		}
	} else if parent != "" {
		args = append(args, parent)
	}

	cmd = cmd + objType + "-delete"
	args = append(args, name)

	_, err = c.executeLBCTL(cmd, t, args...)
	if err != nil {
		return err
	}
	if t == "" {
		err = c.incrementVersion()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetGlobalVersion returns global configuration file version
func (c *LBCTLClient) GetGlobalVersion() (int64, error) {
	return c.getVersion("global")
}

// GetVersion returns configuration file version
func (c *LBCTLClient) GetVersion() (int64, error) {
	return c.getVersion("config")
}

func (c *LBCTLClient) getVersion(t string) (int64, error) {
	var file *os.File
	var err error
	if t == "global" {
		file, err = os.Open(c.GlobalConfigurationFile())
	} else {
		file, err = os.Open(c.ConfigurationFile())
	}

	if err != nil {
		return 0, NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read configuration file %v: %v", c.ConfigurationFile(), err.Error()))
	}
	defer file.Close()

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
						return c.setInitialVersion(true, t)
					}
					version, err := strconv.ParseInt(w[1], 10, 64)
					if err != nil {
						return c.setInitialVersion(true, t)
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
		return c.setInitialVersion(false, t)
	}
	return c.setInitialVersion(false, t)
}

func (c *LBCTLClient) setInitialVersion(hasVersion bool, t string) (int64, error) {
	var cFile string
	var err error
	if t == "global" {
		cFile = c.GlobalConfigurationFile()
	} else {
		cFile = c.ConfigurationFile()
	}

	input, err := ioutil.ReadFile(cFile)

	if err != nil {
		return 0, NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read configuration file %v: %v", cFile, err.Error()))
	}

	inputStr := string(input)

	if hasVersion {
		inputStrArr := strings.SplitAfterN(inputStr, "\n", 2)
		if len(inputStrArr) == 2 {
			inputStr = inputStrArr[1]
		}
	}

	output := fmt.Sprintf("# _version=1\n%s", inputStr)
	if err = ioutil.WriteFile(cFile, []byte(output), 0666); err != nil {
		return 0, NewConfError(ErrCannotSetVersion, fmt.Sprintf("Cannot set initial version in file %v: %v", cFile, err.Error()))
	}
	return 0, nil
}

func (c *LBCTLClient) parseObject(str string, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem()
	for _, line := range strings.Split(str, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "+") {
			lineArr := strings.SplitN(line, " ", 2)
			fieldName := lineArr[0]
			fieldValue := lineArr[1]
			fieldName = misc.CamelCase(fieldName[1:], true)
			field := objValue.FieldByName(fieldName)
			if field.IsValid() {
				if field.CanSet() {
					if field.Kind() == reflect.Int64 {
						fieldValueInt, err := strconv.ParseInt(fieldValue, 10, 64)
						if err != nil {
							continue
						}
						field.SetInt(fieldValueInt)
					} else if field.Kind() == reflect.String {
						field.SetString(fieldValue)
					} else if field.Kind() == reflect.Float64 {
						fieldValueFl, err := strconv.ParseFloat(fieldValue, 64)
						if err != nil {
							continue
						}
						field.SetFloat(fieldValueFl)
					} else if field.Kind() == reflect.Ptr {
						fieldValueInt, err := strconv.ParseInt(fieldValue, 10, 64)
						if err == nil {
							field.Set(reflect.ValueOf(&fieldValueInt))
						} else {
							fieldValueFl, err := strconv.ParseFloat(fieldValue, 64)
							if err == nil {
								field.Set(reflect.ValueOf(&fieldValueFl))
							} else {
								field.Set(reflect.ValueOf(&fieldValue))
							}
						}
					}
				}
			}
		}
	}
}

func (c *LBCTLClient) serializeObject(obj interface{}, ondisk interface{}, skipFields []string) []string {
	var argsArr []string
	objValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		typeField := objValue.Type().Field(i)
		field := objValue.FieldByName(typeField.Name)
		if field.IsValid() {
			if skipFields != nil {
				if misc.StringInSlice(typeField.Name, skipFields) {
					continue
				}
			}
			if typeField.Name != "Name" && typeField.Name != "ID" {
				// fmt.Printf("serializeObject: typeField.Name: %v\n", typeField.Name)
				if field.Kind() == reflect.Int64 {
					// fmt.Printf("serializeObject: field.Kind(): Int64\n")
					if field.Int() != 0 {
						argsArr = append(argsArr, "--"+misc.SnakeCase(typeField.Name))
						argsArr = append(argsArr, strconv.FormatInt(field.Int(), 10))
					}
				} else if field.Kind() == reflect.String {
					// fmt.Printf("serializeObject: field.Kind(): String\n")
					if field.String() != "" {
						argsArr = append(argsArr, "--"+misc.SnakeCase(typeField.Name))
						argsArr = append(argsArr, field.String())
					}
				} else if field.Kind() == reflect.Float64 {
					// fmt.Printf("serializeObject: field.Kind(): Float64\n")
					if field.Float() != 0.0 {
						argsArr = append(argsArr, "--"+misc.SnakeCase(typeField.Name))
						argsArr = append(argsArr, strconv.FormatFloat(field.Float(), 'f', -1, 64))
					}
				} else if field.Kind() == reflect.Ptr {
					// fmt.Printf("serializeObject: field.Kind(): Ptr\n")
					if field.Pointer() != 0 {
						argsArr = append(argsArr, "--"+misc.SnakeCase(typeField.Name))
						p := (*int64)(unsafe.Pointer(field.Pointer()))

						argsArr = append(argsArr, fmt.Sprintf("%v", *p))
					}
				}
				// fmt.Printf("serializeObject: argsArr: %v\n", argsArr)
			}
		}
	}
	//delete options
	if ondisk != nil {
		ondiskObjValue := reflect.ValueOf(ondisk).Elem()
		// fmt.Printf("serializeObject: ondisk is not nil, check for deletion\n")
		for i := 0; i < ondiskObjValue.NumField(); i++ {
			fName := objValue.Type().Field(i).Name
			if fName != "Name" && fName != "LogFormatCustom" {
				// fmt.Printf("serializeObject: fName: %v\n", fName)
				if misc.StringInSlice(fName, skipFields) {
					// fmt.Printf("serializeObject: %v in ignored fields\n", fName)
					continue
				}
				field := objValue.FieldByName(fName)
				// fmt.Printf("serializeObject: field: %v\n", field)
				if misc.IsZeroValue(field) {
					// fmt.Printf("serializeObject: %v is not valid add reset arg\n", fName)
					argsArr = append(argsArr, "--reset-"+misc.SnakeCase(fName))
				}
			}
		}
	}

	// fmt.Printf("serializeObject: argsArr: %v\n", argsArr)

	return argsArr
}

func (c *LBCTLClient) checkTransactionOrVersion(transactionID string, version int64, startTransaction bool) (string, error) {
	// start an implicit transaction for delete site (multiple operations required) if not already given
	t := ""
	if transactionID != "" && version != 0 {
		return "", NewConfError(ErrBothVersionTransaction, "Both version and transaction specified, specify only one")
	} else if transactionID == "" && version == 0 {
		return "", NewConfError(ErrNoVersionTransaction, "Version or transaction not specified, specify only one")
	} else if transactionID != "" {
		t = transactionID
	} else {
		v, err := c.GetVersion()
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

func (c *LBCTLClient) checkVersion(version int64) (bool, error) {
	v, err := c.GetVersion()
	return v == version, err
}

func splitHeaderLine(obj string) (name string, parent string) {
	w := strings.SplitN(obj, "\n", 2)
	if len(w) < 2 {
		return "", ""
	}
	header := w[0]

	headerSpl := strings.SplitN(header, " ", 4)
	if len(headerSpl) < 4 {
		return "", ""
	}
	return headerSpl[1], headerSpl[3]
}

func (c *LBCTLClient) incrementVersion() error {
	input, err := ioutil.ReadFile(c.ConfigurationFile())

	if err != nil {
		return NewConfError(ErrCannotReadVersion, fmt.Sprintf("Cannot read version from file %v: %v", c.ConfigurationFile(), err.Error()))
	}

	v, err := c.GetVersion()
	if err != nil {
		return err
	}

	toReplace := fmt.Sprintf("# _version=%v", v)
	replace := fmt.Sprintf("# _version=%v", v+1)

	output := bytes.Replace(input, []byte(toReplace), []byte(replace), -1)

	if err = ioutil.WriteFile(c.ConfigurationFile(), output, 0666); err != nil {
		return NewConfError(ErrCannotSetVersion, fmt.Sprintf("Cannot increment version in file %v: %v", c.ConfigurationFile(), err.Error()))
	}
	if c.Cache.Enabled() {
		c.Cache.Version.Set(v + 1)
	}
	return nil
}

func (c *LBCTLClient) incrementGlobalVersion() error {
	input, err := ioutil.ReadFile(c.GlobalConfigurationFile())

	if err != nil {
		return NewConfError(ErrCannotReadVersion, fmt.Sprintf("Cannot read version from file %v: %v", c.GlobalConfigurationFile(), err.Error()))
	}

	v, err := c.GetGlobalVersion()
	if err != nil {
		return err
	}

	toReplace := fmt.Sprintf("# _version=%v", v)
	replace := fmt.Sprintf("# _version=%v", v+1)

	output := bytes.Replace(input, []byte(toReplace), []byte(replace), -1)

	if err = ioutil.WriteFile(c.GlobalConfigurationFile(), output, 0666); err != nil {
		return NewConfError(ErrCannotSetVersion, fmt.Sprintf("Cannot increment version in file %v: %v", c.GlobalConfigurationFile(), err.Error()))
	}
	return nil
}

func lbctlTypeToType(lType string) string {
	switch lType {
	case "farm":
		return "Backend"
	case "service":
		return "Frontend"
	case "usefarm":
		return "BackendSwitchingRule"
	case "useserver":
		return "ServerSwitchingRule"
	case "stickreq":
		return "StickRequestRule"
	case "stickrsp":
		return "StickResponseRule"
	case "httpreq":
		return "HTTPRequestRule"
	case "httprsp":
		return "HTTPResponseRule"
	case "tcpreqconn":
		return "TCPConnectionRule"
	case "tcpreqcont":
	case "tcprspcont":
		return "TCPConnectionRule"
	default:
		return misc.CamelCase(lType, true)
	}
	return misc.CamelCase(lType, true)
}

func typeToLbctlType(oType string) string {
	switch oType {
	case "frontend":
		return "service"
	case "backend":
		return "farm"
	default:
		return ""
	}
}
