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
	"os"
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
func (c *client) getRawConfiguration(transactionID string, version int64) (int64, int64, string, string, error) {
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
	ondiskV, ondiskClusterV, ondiskMD5Hash, metaErr := c.getConfigurationMetaData(config)
	if metaErr != nil {
		return 0, 0, "", "", metaErr
	}

	data, err := os.ReadFile(config)
	if err != nil {
		return 0, 0, "", "", NewConfError(ErrCannotReadConfFile, err.Error())
	}

	return ondiskV, ondiskClusterV, ondiskMD5Hash, string(data), nil
}

func (c *client) getConfigurationMetaData(config string) (int64, int64, string, error) {
	file, err := os.Open(config)
	if err != nil {
		return 0, 0, "", NewConfError(ErrCannotReadConfFile, err.Error())
	}
	defer file.Close()

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
		case strings.HasPrefix(line, "# _md5hash="):
			w := strings.Split(line, "=")
			if len(w) == 2 {
				ondiskMD5Hash = strings.TrimSpace(w[1])
			}
		case strings.HasPrefix(line, "# _cluster_version="):
			w := strings.Split(line, "=")
			if len(w) == 2 {
				ondiskClusterV, err = strconv.ParseInt(w[1], 10, 64)
				if err != nil {
					ondiskClusterV = int64(0)
				}
			}
		}
		if ondiskV != 0 && ondiskMD5Hash != "" && ondiskClusterV != 0 {
			break
		}
	}
	if err = scanner.Err(); err != nil {
		return ondiskV, 0, "", NewConfError(ErrCannotReadConfFile, err.Error())
	}

	return ondiskV, ondiskClusterV, ondiskMD5Hash, nil
}

// PostRawConfiguration pushes given string to the config file if the version
// matches
func (c *client) PostRawConfiguration(config *string, version int64, skipVersionCheck bool, onlyValidate ...bool) error {
	if len(onlyValidate) > 0 && onlyValidate[0] {
		f, err := os.CreateTemp("/tmp", "onlyvalidate")
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
		return checkHaproxyConfiguration(c.ConfigurationOptions, f.Name())
	}
	var t string
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
		_, _ = w.WriteString(fmt.Sprintf("# _version=%d\n%s", version, c.dropVersionFromRaw(*config)))
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
		return NewConfError(ErrCannotReadConfFile, "Cannot read "+tFile)
	}

	// Do a regular commit of the transaction
	if _, err := c.commitTransaction(t, skipVersionCheck); err != nil {
		return err
	}

	return nil
}

// dropVersionFromRaw is used when force pushing a raw configuration with version check:
// if the provided user input has already a version metadata it must be withdrawn.
func (c *client) dropVersionFromRaw(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))

	var sanitized strings.Builder

	for scanner.Scan() {
		t := scanner.Bytes()

		if bytes.HasPrefix(t, []byte("# _version=")) {
			continue
		}

		sanitized.Write(t)
		sanitized.WriteByte('\n')
	}

	return sanitized.String()
}
