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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	parser "github.com/haproxytech/config-parser"
	"github.com/haproxytech/config-parser/types"
	"github.com/haproxytech/models"
)

// GetTransactions returns an array of transactions
func (c *Client) GetTransactions(status string) (*models.Transactions, error) {
	return c.parseTransactions(status)
}

// GetTransaction returns transaction information by id
func (c *Client) GetTransaction(id string) (*models.Transaction, error) {
	// check if parser exists, if not, look for files
	_, ok := c.parsers[id]
	if !ok {
		tFile, err := c.getTransactionFile(id)
		if err != nil {
			return nil, NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("Transaction %v does not exist", id))
		}
		return c.parseTransactionFile(tFile), nil
	}
	v, _ := c.GetVersion(id)

	return &models.Transaction{ID: id, Status: "in_progress", Version: v}, nil
}

// StartTransaction starts a new empty lbctl transaction
func (c *Client) StartTransaction(version int64) (*models.Transaction, error) {
	return c.startTransaction(version)
}

func (c *Client) startTransaction(version int64) (*models.Transaction, error) {
	t := &models.Transaction{}

	v, err := c.GetVersion("")
	if err != nil {
		return nil, err
	}

	if version != v {
		return nil, NewConfError(ErrVersionMismatch, fmt.Sprintf("Version in configuration file is %v, given version is %v", v, version))
	}

	t.ID = uuid.New().String()

	if c.PersistentTransactions {
		err = c.createTransactionFiles(t.ID)
		if err != nil {
			return nil, err
		}
	}

	t.Version = version
	t.Status = "in_progress"

	if err := c.AddParser(t.ID); err != nil {
		if c.PersistentTransactions {
			c.deleteTransactionFiles(t.ID)
		}
		return nil, err
	}
	return t, nil
}

// CommitTransaction commits a transaction by id.
func (c *Client) CommitTransaction(id string) (*models.Transaction, error) {
	// check if parser exists and if transaction exists
	c.mu.Lock()
	defer c.mu.Unlock()

	p, err := c.GetParser(id)
	if err != nil {
		return nil, err
	}

	// do a version check before commiting
	version, err := c.GetVersion("")
	if err != nil {
		return nil, err
	}

	tVersion, err := c.GetVersion(id)
	if err != nil {
		return nil, err
	}

	if tVersion != version {
		c.failTransaction(id)
		return nil, NewConfError(ErrVersionMismatch, fmt.Sprintf("Version mismatch, transaction version: %v, configured version: %v", tVersion, version))
	}

	// create transaction file now if transactions are not persistent
	if !c.PersistentTransactions {
		err = c.createTransactionFiles(id)
		if err != nil {
			return nil, err
		}
	}

	transactionFile, err := c.getTransactionFile(id)
	if err != nil {
		return nil, err
	}

	// save to transaction file if transactions are not persistent
	if !c.PersistentTransactions {
		if err := p.Save(transactionFile); err != nil {
			c.failTransaction(id)
			return nil, NewConfError(ErrErrorChangingConfig, err.Error())
		}
	}

	if err := c.checkTransactionFile(id); err != nil {
		c.failTransaction(id)
		return nil, err
	}

	// Fail backing up and cleaning backups silently
	if c.BackupsNumber > 0 {
		copyFile(c.ConfigurationFile, fmt.Sprintf("%v.%v", c.ConfigurationFile, version))
		backupToDel := fmt.Sprintf("%v.%v", c.ConfigurationFile, strconv.Itoa(int(version)-c.BackupsNumber))
		os.Remove(backupToDel)
	}

	if err := copyFile(transactionFile, c.ConfigurationFile); err != nil {
		c.failTransaction(id)
		return nil, err
	}

	c.deleteTransactionFiles(id)

	if err := c.CommitParser(id); err != nil {
		c.Parser.LoadData(c.ConfigurationFile)
		return nil, err
	}

	if err := c.incrementVersion(); err != nil {
		return nil, err
	}

	return &models.Transaction{ID: id, Version: tVersion, Status: "success"}, nil
}

func (c *Client) checkTransactionFile(id string) error {
	transactionFile, err := c.getTransactionFile(id)
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	if c.MasterWorker {
		cmd = exec.Command(c.Haproxy, "-W", "-f", transactionFile, "-c")
	} else {
		cmd = exec.Command(c.Haproxy, "-f", transactionFile, "-c")
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return NewConfError(ErrValidationError, c.parseHAProxyCheckError(stderr.Bytes(), id))
	}
	return nil
}

func (c *Client) parseHAProxyCheckError(output []byte, id string) string {
	oStr := string(output)
	var b strings.Builder
	b.WriteString(fmt.Sprintf("ERR transactionId=%s \n", id))

	for _, line := range strings.Split(oStr, "\n") {
		line := strings.TrimSpace(line)
		if strings.HasPrefix(line, "[ALERT]") {
			if strings.HasSuffix(line, "Fatal errors found in configuration.") {
				continue
			}
			if strings.Contains(line, "Error(s) found in configuration file : ") {
				continue
			}

			parts := strings.Split(line, " : ")
			if len(parts) > 2 && strings.HasPrefix(strings.TrimSpace(parts[1]), "parsing [") {
				fParts := strings.Split(strings.TrimSpace(parts[1]), ":")
				var msgB strings.Builder
				for i := 2; i < len(parts); i++ {
					msgB.WriteString(parts[i])
					msgB.WriteString(" ")
				}
				if len(fParts) > 1 {
					lNo, err := strconv.ParseInt(strings.TrimSuffix(fParts[1], "]"), 10, 64)
					if err == nil {
						b.WriteString(fmt.Sprintf("line=%d msg=\"%s\"\n", lNo, strings.TrimSpace(msgB.String())))
					} else {
						b.WriteString(fmt.Sprintf("msg=\"%s\"\n", strings.TrimSpace(msgB.String())))
					}
				}
			} else if len(parts) > 1 {
				var msgB strings.Builder
				for i := 1; i < len(parts); i++ {
					msgB.WriteString(parts[i])
					msgB.WriteString(" ")
				}
				b.WriteString(fmt.Sprintf("msg=\"%s\"\n", strings.TrimSpace(msgB.String())))
			}
		}
	}
	return strings.TrimSuffix(b.String(), "\n")
}

// DeleteTransaction deletes a transaction by id.
func (c *Client) DeleteTransaction(id string) error {
	if id != "" {
		if c.PersistentTransactions {
			if err := c.deleteTransactionFiles(id); err != nil {
				return err
			}
		}
		c.DeleteParser(id)
	}
	return nil
}

func (c *Client) parseTransactions(status string) (*models.Transactions, error) {
	confFileName := filepath.Base(c.ConfigurationFile)

	_, err := os.Stat(c.TransactionDir)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(c.TransactionDir, 0755)
		if err != nil {
			return nil, err
		}
		return &models.Transactions{}, nil
	}

	transactions := models.Transactions{}
	files, err := ioutil.ReadDir(c.TransactionDir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if !f.IsDir() && status != "failed" && c.PersistentTransactions {
			if strings.HasPrefix(f.Name(), confFileName) {
				transactions = append(transactions, c.parseTransactionFile(filepath.Join(c.TransactionDir, f.Name())))
			}
		} else {
			if f.Name() == "failed" && status != "in_progress" {
				ffiles, err := ioutil.ReadDir(filepath.Join(c.TransactionDir, "failed"))
				if err != nil {
					return nil, err
				}
				for _, ff := range ffiles {
					if !ff.IsDir() {
						if strings.HasPrefix(ff.Name(), confFileName) {
							transactions = append(transactions, c.parseTransactionFile(filepath.Join(c.TransactionDir, "failed", ff.Name())))
						}
					}
				}
			}
		}
	}

	if !c.PersistentTransactions && status != "failed" {
		for tID := range c.parsers {
			v, err := c.GetVersion(tID)
			if err == nil {
				t := &models.Transaction{
					ID:      tID,
					Status:  "in_progress",
					Version: v,
				}
				transactions = append(transactions, t)
			}
		}
	}
	return &transactions, nil
}

func (c *Client) parseTransactionFile(filePath string) *models.Transaction {
	parts := strings.Split(filePath, string(filepath.Separator))
	f := parts[len(parts)-1]
	status := "in_progress"

	if len(parts) > 1 {
		if parts[len(parts)-2] == "failed" {
			status = "failed"
		}
	}

	s := strings.Split(f, ".")
	tID := s[len(s)-1]

	v, err := c.GetVersion(tID)
	if err != nil {
		v, _ = c.getFailedTransactionVersion(tID)
	}

	t := &models.Transaction{
		ID:      tID,
		Status:  status,
		Version: v,
	}
	return t
}

func (c *Client) createTransactionFiles(transactionID string) error {
	transDir, err := os.Stat(c.TransactionDir)

	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(c.TransactionDir, 0755)
		if err != nil {
			return err
		}
	} else {
		if !transDir.Mode().IsDir() {
			return fmt.Errorf("Transaction dir %s is a file", c.TransactionDir)
		}
	}

	confFilePath := filepath.Join(c.TransactionDir, c.getTransactionFileName(transactionID))

	err = copyFile(c.ConfigurationFile, confFilePath)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) deleteTransactionFiles(transactionID string) error {
	confFilePath, err := c.getTransactionFile(transactionID)
	if err != nil {
		return err
	}

	err = os.Remove(confFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func (c *Client) getTransactionFileName(transactionID string) string {
	baseFileName := filepath.Base(filepath.Clean(c.ConfigurationFile))
	return baseFileName + "." + transactionID
}

func (c *Client) getTransactionFile(transactionID string) (string, error) {
	if transactionID == "" {
		return c.ConfigurationFile, nil
	}
	// First find failed transaction file
	transactionFileName := c.getTransactionFileName(transactionID)

	fPath := filepath.Join(c.TransactionDir, "failed", transactionFileName)
	if _, err := os.Stat(fPath); err == nil {
		return fPath, nil
	}
	// Return in progress transaction file if exists, else empty string
	fPath = filepath.Join(c.TransactionDir, transactionFileName)
	if _, err := os.Stat(fPath); err == nil {
		return fPath, nil
	}
	return "", NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("Transaction file %v does not exist", transactionID))
}

func (c *Client) getTransactionFileFailed(transactionID string) string {
	baseFileName := filepath.Base(filepath.Clean(c.ConfigurationFile))
	transactionFileName := baseFileName + "." + transactionID

	return filepath.Join(c.TransactionDir, "failed", transactionFileName)
}

func (c *Client) getBackupFile(version int64) (string, error) {
	if version == 0 {
		return c.ConfigurationFile, nil
	}
	backupFileName := fmt.Sprintf("%v.%v", c.ConfigurationFile, version)

	if _, err := os.Stat(backupFileName); err == nil {
		return backupFileName, nil
	}
	return "", NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backup file for version %v does not exist", version))
}

func (c *Client) failTransaction(id string) {
	failedDir := filepath.Join(c.TransactionDir, "failed")
	if _, err := os.Stat(failedDir); os.IsNotExist(err) {
		os.Mkdir(failedDir, 0755)
	}

	configFile, err := c.getTransactionFile(id)
	if err != nil {
		return
	}

	failedConfigFile := c.getTransactionFileFailed(id)
	copyFile(configFile, failedConfigFile)
	os.Remove(configFile)
	c.DeleteParser(id)
}

func (c *Client) getFailedTransactionVersion(id string) (int64, error) {
	fName := c.getTransactionFileName(id)
	failedDir := filepath.Join(c.TransactionDir, "failed")
	if _, err := os.Stat(failedDir); os.IsNotExist(err) {
		return 0, NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("Transaction %v not failed", id))
	}
	fPath := filepath.Join(failedDir, fName)
	if _, err := os.Stat(fPath); os.IsNotExist(err) {
		return 0, NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("Transaction %v not failed", id))
	}

	p := &parser.Parser{}
	if err := p.LoadData(fPath); err != nil {
		return 0, NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read %s", fPath))
	}

	data, _ := p.Get(parser.Comments, parser.CommentsSectionName, "# _version", false)
	ver, ok := data.(*types.ConfigVersion)
	if !ok {
		return 0, NewConfError(ErrCannotReadVersion, "Cannot read version")
	}
	return ver.Value, nil
}

func copyFile(src, dest string) error {
	srcContent, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcContent.Close()

	destContent, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destContent.Close()

	_, err = io.Copy(destContent, srcContent)
	if err != nil {
		return err
	}
	return destContent.Sync()
}
