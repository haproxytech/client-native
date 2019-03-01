package configuration

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/haproxytech/models"
)

// GetTransactions returns an array of transactions
func (c *Client) GetTransactions(status string) (*models.Transactions, error) {
	return c.parseTransactions(status)
}

// GetTransaction returns transaction information by id
func (c *Client) GetTransaction(id string) (*models.Transaction, error) {
	return c.parseTransaction(id)
}

// StartTransaction starts a new empty lbctl transaction
func (c *Client) StartTransaction(version int64) (*models.Transaction, error) {
	return c.startTransaction(version, true)
}

func (c *Client) startTransaction(version int64, initCache bool) (*models.Transaction, error) {
	t := &models.Transaction{}

	v, err := c.GetVersion("")
	if err != nil {
		return nil, err
	}

	if version != v {
		return nil, NewConfError(ErrVersionMismatch, fmt.Sprintf("Version in configuration file is %v, given version is %v", v, version))
	}

	t.ID = uuid.New().String()
	err = c.createTransactionFiles(t.ID)
	if err != nil {
		return nil, err
	}

	t.Version = version
	t.Status = "in_progress"

	if err := c.AddParser(t.ID); err != nil {
		c.deleteTransactionFiles(t.ID)
		return nil, err
	}
	return t, nil
}

// CommitTransaction commits a transaction by id.
func (c *Client) CommitTransaction(id string) error {
	// do a version check before commiting
	version, err := c.GetVersion("")
	if err != nil {
		return err
	}

	t, err := c.parseTransaction(id)
	if err != nil {
		return err
	}

	if t.Version != version {
		c.failTransaction(id)
		return NewConfError(ErrVersionMismatch, fmt.Sprintf("Version mismatch, transaction version: %v, configured version: %v", t.Version, version))
	}

	if err := c.checkTransactionFile(id); err != nil {
		c.failTransaction(id)
		return err
	}

	if err := copyFile(c.getTransactionFile(id), c.ConfigurationFile); err != nil {
		c.failTransaction(id)
		return err
	}

	c.deleteTransactionFiles(id)

	if err := c.CommitParser(id); err != nil {
		c.Parser.LoadData(c.ConfigurationFile)
		return nil
	}

	if err := c.incrementVersion(); err != nil {
		return err
	}

	return nil
}

func (c *Client) checkTransactionFile(id string) error {
	cmd := exec.Command(c.Haproxy, "-f", c.getTransactionFile(id), "-c")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return NewConfError(ErrValidationError, string(stderr.Bytes()))
	}
	return nil
}

// DeleteTransaction deletes a transaction by id.
func (c *Client) DeleteTransaction(id string) error {
	if id != "" {
		if err := c.deleteTransactionFiles(id); err != nil {
			return err
		}
		if err := c.DeleteParser(id); err != nil {
			return err
		}
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
		if !f.IsDir() && status != "failed" {
			if strings.HasPrefix(f.Name(), confFileName) {
				transactions = append(transactions, c.parseTransactionFile(f, "in_progress"))
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
							transactions = append(transactions, c.parseTransactionFile(ff, "failed"))
						}
					}
				}
			}
		}
	}

	return &transactions, nil
}

func (c *Client) parseTransaction(id string) (*models.Transaction, error) {
	_, err := os.Stat(c.TransactionDir)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(c.TransactionDir, 0755)
		if err != nil {
			return nil, err
		}
		return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Transaction %v does not exist", id))
	}

	fName := filepath.Base(c.ConfigurationFile) + "." + id

	//Check if there is a file in transaction directory
	inProgressFile := filepath.Join(c.TransactionDir, fName)
	f, err := os.Stat(inProgressFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if f != nil {
		return c.parseTransactionFile(f, "in_progress"), nil
	}

	//Check if there is a file in failed directory
	failedFile := filepath.Join(c.TransactionDir, "failed", fName)
	f, err = os.Stat(failedFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Transaction %v does not exist", id))
		}
		return nil, err
	}
	if f != nil {
		return c.parseTransactionFile(f, "failed"), nil
	}
	return nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Transaction %v does not exist", id))
}

func (c *Client) parseTransactionFile(f os.FileInfo, status string) *models.Transaction {
	s := strings.Split(f.Name(), ".")
	tID := s[len(s)-1]
	v := int64(0)
	if status == "in_progress" {
		v, _ = c.GetVersion(tID)
	} else {
		v, _ = c.GetVersion(tID)
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

	confFilePath := c.getTransactionFile(transactionID)

	err = copyFile(c.ConfigurationFile, confFilePath)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) deleteTransactionFiles(transactionID string) error {
	confFilePath := c.getTransactionFile(transactionID)

	err := os.Remove(confFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func (c *Client) getTransactionFile(transactionID string) string {
	if transactionID == "" {
		return c.ConfigurationFile
	}
	// First find failed transaction file
	baseFileName := filepath.Base(filepath.Clean(c.ConfigurationFile))
	transactionFileName := baseFileName + "." + transactionID

	fPath := filepath.Join(c.TransactionDir, "failed", transactionFileName)
	if _, err := os.Stat(fPath); err == nil {
		return fPath
	}
	// Return in progress transaction file
	return filepath.Join(c.TransactionDir, transactionFileName)
}

func (c *Client) getTransactionFileFailed(transactionID string) string {
	baseFileName := filepath.Base(filepath.Clean(c.ConfigurationFile))
	transactionFileName := baseFileName + "." + transactionID

	return filepath.Join(c.TransactionDir, "failed", transactionFileName)
}

func (c *Client) failTransaction(id string) {
	failedDir := filepath.Join(c.TransactionDir, "failed")
	if _, err := os.Stat(failedDir); os.IsNotExist(err) {
		os.Mkdir(failedDir, 0755)
	}

	configFile := c.getTransactionFile(id)
	failedConfigFile := c.getTransactionFileFailed(id)
	copyFile(configFile, failedConfigFile)
	os.Remove(configFile)
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
