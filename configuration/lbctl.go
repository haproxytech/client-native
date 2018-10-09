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

	"github.com/haproxytech/client-native/misc"
)

// LBCTLConfigurationClient configuration.Client implementation using lbctl
type LBCTLConfigurationClient struct {
	*ClientParams
	LBCTLPath    string
	LBCTLTmpPath string
}

const (
	// DefaultLBCTLPath sane default for path to lbctl
	DefaultLBCTLPath string = "/usr/sbin/lbctl"
	// DefaultLBCTLTmpPath sane default for path for lbctl transactions
	DefaultLBCTLTmpPath string = "/tmp/lbctl"
)

// NewLBCTLClient constructor
func NewLBCTLClient(configurationFile string, LBCTLPath string, LBCTLTmpPath string) *LBCTLConfigurationClient {
	if LBCTLPath == "" {
		LBCTLPath = DefaultLBCTLPath
	}

	if LBCTLTmpPath == "" {
		LBCTLTmpPath = DefaultLBCTLTmpPath
	}

	return &LBCTLConfigurationClient{NewConfigurationClientParams(configurationFile), LBCTLPath, LBCTLTmpPath}
}

// DefaultLBCTLClient returns LBCTLConfigurationClient with sane defaults
func DefaultLBCTLClient() *LBCTLConfigurationClient {
	return NewLBCTLClient("", "", "")
}

func (c *LBCTLConfigurationClient) executeLBCTL(command string, transaction string, args ...string) (string, error) {
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
	cmd.Env = append(cmd.Env, "LBCTL_L7_SVC_CHECK_CMD=true")
	cmd.Env = append(cmd.Env, "LBCTL_L7_SVC_APPLY_CMD=true")

	if c.ConfigurationFile() != "" {
		cmd.Env = append(cmd.Env, "LBCTL_L7_HAPROXY_CONFIG="+c.ConfigurationFile())
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

func (c *LBCTLConfigurationClient) createObject(name string, objType string, parent string, parentType string, data interface{}, skipFields []string, transactionID string, version int64) error {
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

func (c *LBCTLConfigurationClient) editObject(name string, objType string, parent string, parentType string, data interface{}, ondisk interface{}, skipFields []string, transactionID string, version int64) error {
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

func (c *LBCTLConfigurationClient) deleteObject(name string, objType string, parent string, parentType string, transactionID string, version int64) error {
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

// GetVersion returns configuration file version
func (c *LBCTLConfigurationClient) GetVersion() (int64, error) {
	file, err := os.Open(c.ConfigurationFile())
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
						return c.setInitialVersion(true)
					}
					version, err := strconv.ParseInt(w[1], 10, 64)
					if err != nil {
						return c.setInitialVersion(true)
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
		return c.setInitialVersion(false)
	}
	return c.setInitialVersion(true)
}

func (c *LBCTLConfigurationClient) setInitialVersion(hasVersion bool) (int64, error) {
	input, err := ioutil.ReadFile(c.ConfigurationFile())

	if err != nil {
		return 0, NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read configuration file %v: %v", c.ConfigurationFile(), err.Error()))
	}

	inputStr := string(input)

	if hasVersion {
		inputStrArr := strings.SplitAfterN(inputStr, "\n", 2)
		if len(inputStrArr) == 2 {
			inputStr = inputStrArr[1]
		}
	}

	output := fmt.Sprintf("# _version=1\n%s", inputStr)

	if err = ioutil.WriteFile(c.ConfigurationFile(), []byte(output), 0666); err != nil {
		return 0, NewConfError(ErrCannotSetVersion, fmt.Sprintf("Cannot set initial version in file %v: %v", c.ConfigurationFile(), err.Error()))
	}
	return 0, nil
}

func (c *LBCTLConfigurationClient) parseObject(str string, obj interface{}) {
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

func (c *LBCTLConfigurationClient) serializeObject(obj interface{}, ondisk interface{}, skipFields []string) []string {
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

func (c *LBCTLConfigurationClient) checkTransactionOrVersion(transactionID string, version int64, startTransaction bool) (string, error) {
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

func (c *LBCTLConfigurationClient) checkVersion(version int64) (bool, error) {
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

func (c *LBCTLConfigurationClient) incrementVersion() error {
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
