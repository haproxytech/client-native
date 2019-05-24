// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this files except in compliance with the License.
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
	"os/exec"
	"strconv"
	"strings"
)

// GetRawConfiguration returns configuration version and a
// string containing raw config file
func (c *Client) GetRawConfiguration(transactionID string) (int64, string, error) {
	config := c.ConfigurationFile
	if transactionID != "" {
		config = c.getTransactionFile(transactionID)
		if config == "" {
			return 0, "", NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("Transaction %s does not exist", transactionID))
		}
	}
	file, err := os.Open(config)
	if err != nil {
		return 0, "", NewConfError(ErrCannotReadConfFile, err.Error())
	}
	defer file.Close()

	dataStr := ""
	ondiskV := int64(0)
	scanner := bufio.NewScanner(file)
	// parse out version
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# _version=") {
			w := strings.Split(line, "=")
			if len(w) == 2 {
				ondiskV, err = strconv.ParseInt(w[1], 10, 64)
				if err != nil {
					ondiskV = int64(0)
				}
			}
		} else {
			dataStr += line + "\n"
		}
	}
	if err = scanner.Err(); err != nil {
		return 0, "", NewConfError(ErrCannotReadConfFile, err.Error())
	}

	return ondiskV, dataStr, nil
}

// PostRawConfiguration pushes given string to the config file if the version
// matches
func (c *Client) PostRawConfiguration(config *string, version int64) error {
	// Create implicit transaction and check version
	t, err := c.checkTransactionOrVersion("", version)
	if err != nil {
		// if transaction is implicit, return err and delete transaction
		if t != "" {
			return c.errAndDeleteTransaction(err, t)
		}
		return err
	}

	// Write the transaction file directly
	tmp, err := os.OpenFile(c.getTransactionFile(t), os.O_RDWR|os.O_TRUNC, 0777)
	defer tmp.Close()
	if err != nil {
		return NewConfError(ErrCannotReadConfFile, err.Error())
	}

	w := bufio.NewWriter(tmp)
	w.WriteString(fmt.Sprintf("# _version=%v\n%v", version, *config))
	w.Flush()

	// Load the data into the transaction parser
	p, err := c.GetParser(t)
	if err != nil {
		return err
	}

	if err := p.LoadData(c.getTransactionFile(t)); err != nil {
		return NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read %s", c.getTransactionFile(t)))
	}

	// Do a regular commit of the transaction
	if _, err := c.CommitTransaction(t); err != nil {
		return err
	}

	return nil
}

func (c *Client) validateConfigFile(confFile string) error {
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
