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

type LBCTLError struct {
	msg    string
	stderr string
	cmd    string
}

func (self *LBCTLError) Error() string {
	return fmt.Sprintf("Error executing: %s, %s. Output: %s", self.cmd, self.msg, self.stderr)
}

type LBCTLConfigurationClient struct {
	*ConfigurationClientParams
	LBCTLPath    string
	LBCTLTmpPath string
}

const (
	DefaultLBCTLPath    string = "/usr/sbin/lbctl"
	DefaultLBCTLTmpPath string = "/tmp/lbctl"
)

func NewLBCTLClient(configurationFile string, LBCTLPath string, LBCTLTmpPath string) *LBCTLConfigurationClient {
	if LBCTLPath == "" {
		LBCTLPath = DefaultLBCTLPath
	}

	if LBCTLTmpPath == "" {
		LBCTLTmpPath = DefaultLBCTLTmpPath
	}

	return &LBCTLConfigurationClient{NewConfigurationClientParams(configurationFile), LBCTLPath, LBCTLTmpPath}
}

func DefaultLBCTLClient() *LBCTLConfigurationClient {
	return NewLBCTLClient("", "", "")
}

func (self *LBCTLConfigurationClient) executeLBCTL(command string, transaction string, args ...string) (string, error) {
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

	cmd := exec.Command(self.LBCTLPath, lbctlArgs...)
	if self.ConfigurationFile() != "" {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "LBCTL_L7_HAPROXY_CONFIG="+self.ConfigurationFile())
	}

	if self.LBCTLTmpPath == "" {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "LBCTL_TRANS_DIR="+self.LBCTLTmpPath)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		output := string(stderr.Bytes())
		if strings.HasPrefix(output, "LBCTL Fatal:") {
			if strings.Contains(output, "Transaction "+transaction+", not found") {
				return "", ErrTransactionNotFound
			} else if strings.HasSuffix(output, "does not exist") {
				return "", ErrNotExists
			}
		}
		return "", &LBCTLError{err.Error(), output, cmd.Path}
	}

	return string(stdout.Bytes()), nil
}

func (self *LBCTLConfigurationClient) createObject(name string, objType string, parent string, parentType string, data interface{}, skipFields []string, transactionID string, version int64) error {
	var args []string

	t, err := self.checkTransactionOrVersion(transactionID, version, false)
	if err != nil {
		return err
	}

	cmd := "l7-"
	if parentType != "" {
		cmd = cmd + parentType + "-"
		if parent != "" {
			args = append(args, parent)
		} else {
			return ErrNoParentSpecified
		}
	}

	cmd = cmd + objType + "-create"
	args = append(args, name)
	args = append(args, self.serializeObject(data, nil, skipFields)...)

	_, err = self.executeLBCTL(cmd, t, args...)
	if err != nil {
		return err
	}

	if t == "" {
		err = self.incrementVersion()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *LBCTLConfigurationClient) editObject(name string, objType string, parent string, parentType string, data interface{}, ondisk interface{}, skipFields []string, transactionID string, version int64) error {
	var args []string

	t, err := self.checkTransactionOrVersion(transactionID, version, false)
	if err != nil {
		return err
	}

	cmd := "l7-"
	if parentType != "" {
		cmd = cmd + parentType + "-"
		if parent != "" {
			args = append(args, parent)
		} else {
			return ErrNoParentSpecified
		}
	}

	cmd = cmd + objType + "-update"
	args = append(args, name)
	args = append(args, self.serializeObject(data, ondisk, skipFields)...)

	_, err = self.executeLBCTL(cmd, t, args...)
	if err != nil {
		return err
	}
	if t == "" {
		err = self.incrementVersion()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *LBCTLConfigurationClient) deleteObject(name string, objType string, parent string, parentType string, transactionID string, version int64) error {
	args := make([]string, 0, 1)

	t, err := self.checkTransactionOrVersion(transactionID, version, false)
	if err != nil {
		return err
	}

	cmd := "l7-"
	if parentType != "" {
		cmd = cmd + parentType + "-"
		if parent != "" {
			args = append(args, parent)
		} else {
			return ErrNoParentSpecified
		}
	}
	cmd = cmd + objType + "-delete"
	args = append(args, name)

	_, err = self.executeLBCTL(cmd, t, args...)
	if err != nil {
		return err
	}
	if t == "" {
		err = self.incrementVersion()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *LBCTLConfigurationClient) GetVersion() (int64, error) {
	file, err := os.Open(self.ConfigurationFile())
	if err != nil {
		return 0, ErrCannotReadConfFile
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# _version=") {
			w := strings.Split(line, "=")
			if len(w) != 2 {
				return 0, ErrCannotReadVersion
			}
			version, err := strconv.ParseInt(w[1], 10, 64)
			if err != nil {
				return 0, ErrCannotReadVersion
			}
			return version, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, ErrCannotReadConfFile
	}
	return 0, ErrCannotReadVersion
}

func (self *LBCTLConfigurationClient) parseObject(str string, obj interface{}) {
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

func (self *LBCTLConfigurationClient) serializeObject(obj interface{}, ondisk interface{}, skipFields []string) []string {
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
			if typeField.Name != "Name" {
				// fmt.Printf("serializeObject: typeField.Name: %v\n", typeField.Name)
				if field.Kind() == reflect.Int64 {
					// fmt.Printf("serializeObject: field.Kind(): Int64\n")
					if field.Int() != 0 {
						argsArr = append(argsArr, "--"+misc.SnakeCase(typeField.Name))
						argsArr = append(argsArr, string(field.Int()))
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

func (self *LBCTLConfigurationClient) checkTransactionOrVersion(transactionID string, version int64, startTransaction bool) (string, error) {
	// start an implicit transaction for delete site (multiple operations required) if not already given
	t := ""
	if transactionID != "" && version != 0 {
		return "", ErrBothVersionTransaction
	} else if transactionID == "" && version == 0 {
		return "", ErrNoVersionTransaction
	} else if transactionID != "" {
		t = transactionID
	} else {
		v, _ := self.GetVersion()
		if version != v {
			return "", ErrVersionMismatch
		}
		if startTransaction {
			transaction, err := self.StartTransaction(version)
			if err != nil {
				return "", err
			}
			t = transaction.ID
		}
	}
	return t, nil
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

func (self *LBCTLConfigurationClient) incrementVersion() error {
	input, err := ioutil.ReadFile(self.ConfigurationFile())

	if err != nil {
		return ErrCannotReadConfFile
	}

	v, err := self.GetVersion()
	if err != nil {
		return err
	}

	toReplace := fmt.Sprintf("# _version=%v", v)
	replace := fmt.Sprintf("# _version=%v", v+1)

	output := bytes.Replace(input, []byte(toReplace), []byte(replace), -1)

	if err = ioutil.WriteFile(self.ConfigurationFile(), output, 0666); err != nil {
		return ErrCannotIncrementVersion
	}
	return nil
}
