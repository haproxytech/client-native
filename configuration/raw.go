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
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Raw interface {
	GetRawConfigurationWithClusterData(transactionID string, version int64) (int64, int64, string, string, error)
	GetRawConfiguration(transactionID string, version int64) (int64, string, error)
	PostRawConfiguration(config *string, version int64, skipVersionCheck bool, onlyValidate ...bool) error
}

func (c *client) GetRawConfiguration(transactionID string, version int64) (int64, string, error) {
	v, _, _, data, err := c.getRawConfiguration(transactionID, version)
	if err != nil {
		return 0, "", err
	}
	return v, data, nil
}

func (c *client) GetRawConfigurationWithClusterData(transactionID string, version int64) (int64, int64, string, string, error) {
	return c.getRawConfiguration(transactionID, version)
}

// GetRawConfiguration returns configuration version and a
// string containing raw config file
func (c *client) getRawConfiguration(transactionID string, version int64) (int64, int64, string, string, error) { //nolint: gocognit
	config := c.ConfigurationFile
	var err error
	if transactionID != "" && version != 0 {
		return 0, 0, "", "", NewConfError(ErrBothVersionTransaction, "Both version and transactionID specified, specify only one")
	}
	if transactionID != "" {
		config, err = c.GetTransactionFile(transactionID)
		if err != nil {
			return 0, 0, "", "", err
		}
	} else if version != 0 {
		config, err = c.getBackupFile(version)
		if err != nil {
			return 0, 0, "", "", err
		}
	}
	file, err := os.Open(config)
	if err != nil {
		return 0, 0, "", "", NewConfError(ErrCannotReadConfFile, err.Error())
	}
	defer file.Close()

	dataStr := ""
	ondiskV := int64(0)
	ondiskClusterV := int64(0)
	ondiskMD5Hash := ""
	scanner := bufio.NewScanner(file)
	// parse out version
	for scanner.Scan() {
		switch line := scanner.Text(); {
		case strings.HasPrefix(line, "# _version="):
			w := strings.Split(line, "=")
			if len(w) == 2 {
				ondiskV, err = strconv.ParseInt(w[1], 10, 64)
				if err != nil {
					ondiskV = int64(0)
				}
			}
			dataStr += line + "\n"
		case strings.HasPrefix(line, "# _md5hash="):
			w := strings.Split(line, "=")
			if len(w) == 2 {
				ondiskMD5Hash = strings.TrimSpace(w[1])
			}
			dataStr += line + "\n"
		case strings.HasPrefix(line, "# _cluster_version="):
			w := strings.Split(line, "=")
			if len(w) == 2 {
				ondiskClusterV, err = strconv.ParseInt(w[1], 10, 64)
				if err != nil {
					ondiskClusterV = int64(0)
				}
			}
			dataStr += line + "\n"
		default:
			dataStr += line + "\n"
		}
	}
	if err = scanner.Err(); err != nil {
		return ondiskV, 0, "", "", NewConfError(ErrCannotReadConfFile, err.Error())
	}

	return ondiskV, ondiskClusterV, ondiskMD5Hash, dataStr, nil
}

// PostRawConfiguration pushes given string to the config file if the version
// matches
func (c *client) PostRawConfiguration(config *string, version int64, skipVersionCheck bool, onlyValidate ...bool) error {
	if len(onlyValidate) > 0 && onlyValidate[0] {
		f, err := ioutil.TempFile("/tmp", "onlyvalidate")
		if err != nil {
			return NewConfError(ErrGeneralError, err.Error())
		}
		defer os.Remove(f.Name())
		_, err = f.WriteString(*config)
		if err != nil {
			return NewConfError(ErrGeneralError, err.Error())
		}
		err = f.Sync()
		if err != nil {
			return NewConfError(ErrGeneralError, err.Error())
		}
		err = c.validateConfigFile(f.Name())
		if err != nil {
			return err
		}
		return nil
	}
	t := ""
	if skipVersionCheck {
		// Create impicit transaction
		transaction, err := c.startTransaction(version, skipVersionCheck)
		if err != nil {
			return err
		}
		t = transaction.ID
	} else {
		// Create implicit transaction and check version
		var err error
		t, err = c.CheckTransactionOrVersion("", version)
		if err != nil {
			// if transaction is implicit, return err and delete transaction
			if t != "" {
				return c.ErrAndDeleteTransaction(err, t)
			}
			return err
		}
	}

	tFile, err := c.GetTransactionFile(t)
	if err != nil {
		return err
	}
	// Write the transaction file directly
	tmp, err := os.OpenFile(tFile, os.O_RDWR|os.O_TRUNC, 0o777)
	defer func() { _ = tmp.Close() }()

	if err != nil {
		return NewConfError(ErrCannotReadConfFile, err.Error())
	}

	w := bufio.NewWriter(tmp)
	if !skipVersionCheck {
		_, _ = w.WriteString(fmt.Sprintf("# _version=%v\n%v", version, *config))
	} else {
		_, _ = w.WriteString(*config)
	}
	_ = w.Flush()

	// Load the data into the transaction parser
	p, err := c.GetParser(t)
	if err != nil {
		return err
	}

	if err := p.LoadData(tFile); err != nil {
		return NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read %s", tFile))
	}

	// Do a regular commit of the transaction
	if _, err := c.commitTransaction(t, skipVersionCheck); err != nil {
		return err
	}

	return nil
}

func (c *client) validateConfigFile(confFile string) error {
	// #nosec G204
	cmd := exec.Command(c.Haproxy)
	cmd.Args = append(cmd.Args, "-c")

	if confFile != "" {
		cmd.Args = append(cmd.Args, "-f")
		cmd.Args = append(cmd.Args, confFile)
	} else {
		cmd.Args = append(cmd.Args, "-f")
		cmd.Args = append(cmd.Args, c.ConfigurationFile)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return NewConfError(ErrValidationError, err.Error())
	}
	return nil
}
