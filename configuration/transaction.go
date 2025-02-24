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
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	spoe "github.com/haproxytech/client-native/v6/config-parser/spoe"

	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
)

type TransactionClient interface {
	GetVersion(transactionID string) (int64, error)
	AddParser(transactionID string) error
	CommitParser(transactionID string) error
	DeleteParser(transactionID string) error
	IncrementTransactionVersion(transactionID string) error
	LoadData(filename string) error
	Save(transactionFile, transactionID string) error
	HasParser(transactionID string) bool
	GetParserTransactions() models.Transactions
	GetFailedParserTransactionVersion(transactionID string) (int64, error)
	CheckTransactionOrVersion(transactionID string, version int64) (string, error)
	SetValidateConfigFiles(before, after []string)
}

// transactionCleanerHandler is just a type dealing with a transaction file:
// actually implemented moving to the `failed` or `outdated` folder.
type transactionCleanerHandler func(transactionId, configurationFile string)

type Transaction struct {
	TransactionClient TransactionClient
	options.ConfigurationOptions
	mu                  sync.Mutex
	noNamedDefaultsFrom bool
}

type Transactions interface {
	GetTransactions(status string) (*models.Transactions, error)
	GetTransaction(transactionID string) (*models.Transaction, error)
	StartTransaction(version int64) (*models.Transaction, error)
	DeleteTransaction(transactionID string) error
	CommitTransaction(transactionID string) (*models.Transaction, error)
	MarkTransactionOutdated(transactionID string) (err error)
	SetValidateConfigFiles(before, after []string)
}

// GetTransactions returns an array of transactions
func (t *Transaction) GetTransactions(status string) (*models.Transactions, error) {
	return t.parseTransactions(status)
}

// GetTransaction returns transaction information by id
func (t *Transaction) GetTransaction(transactionID string) (*models.Transaction, error) {
	// check if parser exists, if not, look for files
	ok := t.TransactionClient.HasParser(transactionID)
	if !ok {
		tFile, err := t.GetTransactionFile(transactionID)
		if err != nil {
			return nil, NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("transaction %v does not exist", transactionID))
		}
		return t.parseTransactionFile(tFile), nil
	}
	v, _ := t.TransactionClient.GetVersion(transactionID)

	return &models.Transaction{ID: transactionID, Status: models.TransactionStatusInProgress, Version: v}, nil
}

// StartTransaction starts a new empty lbctl transaction
func (t *Transaction) StartTransaction(version int64) (*models.Transaction, error) {
	return t.startTransaction(version, false)
}

func (t *Transaction) startTransaction(version int64, skipVersion bool) (*models.Transaction, error) {
	m := &models.Transaction{}

	if !skipVersion {
		v, err := t.TransactionClient.GetVersion("")
		if err != nil {
			return nil, err
		}

		if version != v {
			return nil, NewConfError(ErrVersionMismatch, fmt.Sprintf("version in configuration file is %v, given version is %v", v, version))
		}
	}

	m.ID = uuid.New().String()

	if t.PersistentTransactions {
		err := t.createTransactionFiles(m.ID)
		if err != nil {
			return nil, err
		}
	}

	m.Version = version
	m.Status = models.TransactionStatusInProgress

	if err := t.TransactionClient.AddParser(m.ID); err != nil {
		if t.PersistentTransactions {
			_ = t.deleteTransactionFiles(m.ID)
		}
		return nil, err
	}
	return m, nil
}

// CommitTransaction commits a transaction by id.
func (t *Transaction) CommitTransaction(transactionID string) (*models.Transaction, error) {
	return t.commitTransaction(transactionID, false)
}

// CommitTransaction commits a transaction by id.
func (t *Transaction) commitTransaction(transactionID string, skipVersion bool) (*models.Transaction, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	// do a version check before committing
	version, err := t.TransactionClient.GetVersion("")
	if err != nil {
		return nil, err
	}

	tVersion, err := t.TransactionClient.GetVersion(transactionID)
	if err != nil {
		return nil, err
	}

	if !skipVersion {
		if tVersion != version {
			t.failTransaction(transactionID, t.writeOutdatedTransaction)
			return nil, NewConfError(ErrVersionMismatch, fmt.Sprintf("version mismatch, transaction version: %v, configured version: %v", tVersion, version))
		}
	}

	// create transaction file now if transactions are not persistent
	if !t.PersistentTransactions {
		err = t.createTransactionFiles(transactionID)
		if err != nil {
			return nil, err
		}
	}

	transactionFile, err := t.GetTransactionFile(transactionID)
	if err != nil {
		return nil, err
	}

	// Always save parsed transaction file in order to validate the exact same
	// configuration that will be deployed
	if err := t.TransactionClient.Save(transactionFile, transactionID); err != nil {
		t.failTransaction(transactionID, t.writeFailedTransaction)
		return nil, NewConfError(ErrErrorChangingConfig, err.Error())
	}

	if !skipVersion {
		if err := t.TransactionClient.IncrementTransactionVersion(transactionID); err != nil {
			return nil, err
		}
	}

	if err := t.checkTransactionFile(transactionID); err != nil {
		t.failTransaction(transactionID, t.writeFailedTransaction)
		return nil, err
	}

	// Fail backing up and cleaning backups silently
	if t.BackupsNumber > 0 {
		t.backupCfgAndCleanup(version)
	}

	if err := t.TransactionClient.Save(t.ConfigurationFile, transactionID); err != nil {
		t.failTransaction(transactionID, t.writeFailedTransaction)
		return nil, err
	}

	_ = t.deleteTransactionFiles(transactionID)

	if err := t.TransactionClient.CommitParser(transactionID); err != nil {
		_ = t.TransactionClient.LoadData(t.ConfigurationFile)
		return nil, err
	}

	return &models.Transaction{ID: transactionID, Version: tVersion, Status: "success"}, nil
}

func (t *Transaction) backupCfgAndCleanup(version int64) {
	backupFilePrefix := filepath.Join(t.BackupsDir, filepath.Base(t.ConfigurationFile))

	backupConfFile := fmt.Sprintf("%v.%v", backupFilePrefix, strconv.Itoa(int(version)))
	_ = t.TransactionClient.Save(backupConfFile, "")
	backupToDel := fmt.Sprintf("%v.%v", backupFilePrefix, strconv.Itoa(int(version)-t.BackupsNumber))
	os.Remove(backupToDel)
}

func (t *Transaction) checkTransactionFile(transactionID string) error {
	// check only against HAProxy file
	_, ok := t.TransactionClient.(*client)
	if !ok {
		return nil
	}
	// there are some cases when we don't want to validate a config file,
	// such as if want to use different HAProxy (community, enterprise, aloha)
	// where different options are supported.
	// By disabling validation we can still use DPAPI
	if t.ConfigurationOptions.SkipConfigurationFileValidation {
		return nil
	}

	transactionFile, err := t.GetTransactionFile(transactionID)
	if err != nil {
		return err
	}

	return checkHaproxyConfiguration(t.ConfigurationOptions, transactionFile, transactionID)
}

func (t *Transaction) CheckTransactionOrVersion(transactionID string, version int64) (string, error) {
	// start an implicit transaction if transaction is not already given
	var tID string
	if transactionID != "" && version != 0 {
		return "", NewConfError(ErrBothVersionTransaction, "")
	}
	if transactionID == "" && version == 0 {
		return "", NewConfError(ErrNoVersionTransaction, "version or transaction not specified, specify only one")
	}
	if transactionID != "" {
		tID = transactionID
	} else {
		v, err := t.TransactionClient.GetVersion("")
		if err != nil {
			return "", err
		}
		if version != v {
			return "", NewConfError(ErrVersionMismatch, fmt.Sprintf("version in configuration file is %v, given version is %v", v, version))
		}

		transaction, err := t.StartTransaction(version)
		if err != nil {
			return "", err
		}
		tID = transaction.ID
	}
	return tID, nil
}

// MarkTransactionOutdated is marking the transaction by ID as outdated due to a newer commit,
// moving it to the `outdated` folder, as well cleaning from the current parsers.
func (t *Transaction) MarkTransactionOutdated(transactionID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// retrieving current version
	var version int64
	version, err := t.TransactionClient.GetVersion("")
	if err != nil {
		return err
	}
	// retrieving transaction version: needed for comparison
	var tVersion int64
	tVersion, err = t.TransactionClient.GetVersion(transactionID)
	if err != nil {
		return err
	}

	switch {
	case tVersion > version:
		return fmt.Errorf("transaction %s version (%d) is greater than the current (%d) (are you back from the future?)", transactionID, tVersion, version)
	case tVersion == version:
		return fmt.Errorf("transaction %s version is in even with the current one (%d), rather perform deletion", transactionID, version)
	}

	t.failTransaction(transactionID, t.writeOutdatedTransaction)

	return nil
}

// DeleteTransaction deletes a transaction by id.
func (t *Transaction) DeleteTransaction(transactionID string) error {
	if transactionID == "" {
		return nil
	}

	if t.PersistentTransactions {
		if err := t.deleteTransactionFiles(transactionID); err != nil {
			return err
		}
	}

	// Parsers from transactions with `failed` status are deleted when CommitParser implementation is invoked.
	// Because of that, we should not try to delete already deleted parser.
	// Parser with `in_progress` status transaction should be deleted.
	ts := t.TransactionClient.GetParserTransactions()
	for _, tr := range ts {
		if tr.ID == transactionID && tr.Status == models.TransactionStatusInProgress {
			err := t.TransactionClient.DeleteParser(transactionID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Transaction) parseTransactions(status string) (*models.Transactions, error) { //nolint:gocognit
	confFileName := filepath.Base(t.ConfigurationFile)

	_, err := os.Stat(t.TransactionDir)
	if err != nil && os.IsNotExist(err) {
		errMkdir := os.MkdirAll(t.TransactionDir, 0o755)
		if errMkdir != nil {
			return nil, errMkdir
		}
		return &models.Transactions{}, nil
	}

	transactions := models.Transactions{}
	files, err := os.ReadDir(t.TransactionDir)
	if err != nil {
		return nil, err
	}

	readDirAndAppend := func(f fs.DirEntry) error {
		var ffiles []fs.DirEntry
		ffiles, err = os.ReadDir(filepath.Join(t.TransactionDir, f.Name()))
		if err != nil {
			return err
		}
		for _, ff := range ffiles {
			if !ff.IsDir() {
				if strings.HasPrefix(ff.Name(), confFileName) {
					transactions = append(transactions, t.parseTransactionFile(filepath.Join(t.TransactionDir, f.Name(), ff.Name())))
				}
			}
		}
		return nil
	}

	for _, f := range files {
		switch {
		// regular file
		case !f.IsDir() && t.PersistentTransactions && (status == "" || status == "in_progress"):
			if strings.HasPrefix(f.Name(), confFileName) {
				transactions = append(transactions, t.parseTransactionFile(filepath.Join(t.TransactionDir, f.Name())))
			}
		case status == models.TransactionStatusFailed:
			if f.Name() == models.TransactionStatusFailed {
				if err = readDirAndAppend(f); err != nil {
					return nil, err
				}
			}
		case status == models.TransactionStatusOutdated:
			if f.Name() == models.TransactionStatusOutdated {
				if err = readDirAndAppend(f); err != nil {
					return nil, err
				}
			}
		case f.IsDir() && status == "":
			if err = readDirAndAppend(f); err != nil {
				return nil, err
			}
		}
	}

	if !t.PersistentTransactions && status != models.TransactionStatusFailed {
		pt := t.TransactionClient.GetParserTransactions()
		if len(pt) > 0 {
			transactions = append(transactions, pt...)
		}
	}
	return &transactions, nil
}

func (t *Transaction) parseTransactionFile(filePath string) *models.Transaction {
	parts := strings.Split(filePath, string(filepath.Separator))
	f := parts[len(parts)-1]
	status := models.TransactionStatusInProgress

	if len(parts) > 1 {
		switch parts[len(parts)-2] {
		case models.TransactionStatusFailed:
			status = models.TransactionStatusFailed
		case models.TransactionStatusOutdated:
			status = models.TransactionStatusOutdated
		}
	}

	s := strings.Split(f, ".")
	tID := s[len(s)-1]

	v, err := t.TransactionClient.GetVersion(tID)
	if err != nil {
		v, _ = t.getFailedTransactionVersion(tID)
	}

	m := &models.Transaction{
		ID:      tID,
		Status:  status,
		Version: v,
	}
	return m
}

func (t *Transaction) createTransactionFiles(transactionID string) error {
	transDir, err := os.Stat(t.TransactionDir)

	if err != nil && os.IsNotExist(err) {
		errMkdir := os.MkdirAll(t.TransactionDir, 0o755)
		if errMkdir != nil {
			return errMkdir
		}
	} else if !transDir.Mode().IsDir() {
		return fmt.Errorf("transaction dir %s is a file", t.TransactionDir)
	}

	confFilePath := filepath.Join(t.TransactionDir, t.getTransactionFileName(transactionID))

	err = t.TransactionClient.Save(confFilePath, "")
	if err != nil {
		return err
	}

	return nil
}

func (t *Transaction) deleteTransactionFiles(transactionID string) error {
	confFilePath, err := t.GetTransactionFile(transactionID)
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

func (t *Transaction) getTransactionFileName(transactionID string) string {
	baseFileName := filepath.Base(filepath.Clean(t.ConfigurationFile))
	return baseFileName + "." + transactionID
}

func (t *Transaction) GetTransactionFile(transactionID string) (string, error) {
	if transactionID == "" {
		return t.ConfigurationFile, nil
	}
	// First find failed transaction file
	transactionFileName := t.getTransactionFileName(transactionID)

	var fPath string
	fPath = filepath.Join(t.TransactionDir, models.TransactionStatusOutdated, transactionFileName)
	if _, err := os.Stat(fPath); err == nil {
		return fPath, nil
	}
	fPath = filepath.Join(t.TransactionDir, models.TransactionStatusFailed, transactionFileName)
	if _, err := os.Stat(fPath); err == nil {
		return fPath, nil
	}
	// Return in progress transaction file if exists, else empty string
	fPath = filepath.Join(t.TransactionDir, transactionFileName)
	if _, err := os.Stat(fPath); err == nil {
		return fPath, nil
	}
	return "", NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("transaction file %v does not exist", transactionID))
}

func (t *Transaction) getTransactionFile(transactionID, status string) string {
	baseFileName := filepath.Base(filepath.Clean(t.ConfigurationFile))
	transactionFileName := baseFileName + "." + transactionID

	return filepath.Join(t.TransactionDir, status, transactionFileName)
}

func (t *Transaction) getBackupFile(version int64) (string, error) {
	if version == 0 {
		return t.ConfigurationFile, nil
	}
	fileName := fmt.Sprintf("%v.%v", t.ConfigurationFile, version)
	backupFileName := filepath.Join(t.BackupsDir, filepath.Base(fileName))

	if _, err := os.Stat(backupFileName); err == nil {
		return backupFileName, nil
	}
	return "", NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("backup file for version %v does not exist", version))
}

func (t *Transaction) failTransaction(transactionID string, txHandler transactionCleanerHandler) {
	configFile, err := t.GetTransactionFile(transactionID)
	if err != nil {
		return
	}

	if t.SkipFailedTransactions {
		os.Remove(configFile)
	} else {
		txHandler(transactionID, configFile)
	}
	_ = t.TransactionClient.DeleteParser(transactionID)
}

func (t *Transaction) writeOutdatedTransaction(transactionID, configFile string) {
	outdatedDir := filepath.Join(t.TransactionDir, models.TransactionStatusOutdated)
	if _, err := os.Stat(outdatedDir); os.IsNotExist(err) {
		_ = os.Mkdir(outdatedDir, 0o755)
	}
	outdatedConfigFile := t.getTransactionFile(transactionID, models.TransactionStatusOutdated)
	if err := moveFile(configFile, outdatedConfigFile); err != nil {
		_ = os.Remove(configFile)
	}
}

func (t *Transaction) writeFailedTransaction(transactionID, configFile string) {
	failedDir := filepath.Join(t.TransactionDir, models.TransactionStatusFailed)
	if _, err := os.Stat(failedDir); os.IsNotExist(err) {
		_ = os.Mkdir(failedDir, 0o755)
	}
	failedConfigFile := t.getTransactionFile(transactionID, models.TransactionStatusFailed)
	if err := moveFile(configFile, failedConfigFile); err != nil {
		_ = os.Remove(configFile)
	}
}

func (t *Transaction) getFailedTransactionVersion(transactionID string) (int64, error) {
	fName := t.getTransactionFileName(transactionID)
	failedDir := filepath.Join(t.TransactionDir, models.TransactionStatusFailed)
	if _, err := os.Stat(failedDir); os.IsNotExist(err) {
		return 0, NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("transaction %v not failed", transactionID))
	}
	fPath := filepath.Join(failedDir, fName)
	if _, err := os.Stat(fPath); os.IsNotExist(err) {
		return 0, NewConfError(ErrTransactionDoesNotExist, fmt.Sprintf("transaction %v not failed", transactionID))
	}

	_, err := parser.New(
		parser_options.Path(fPath),
	)
	if err != nil {
		return 0, NewConfError(ErrCannotReadConfFile, "cannot read "+fPath)
	}

	ver, err := t.TransactionClient.GetFailedParserTransactionVersion(transactionID)
	if err != nil {
		return 0, NewConfError(ErrCannotReadVersion, "")
	}
	return ver, nil
}

func moveFile(src, dest string) error {
	return os.Rename(src, dest)
}

func (t *Transaction) SaveData(prsr interface{}, tID string, commitImplicit bool) error {
	if t.PersistentTransactions {
		tFile, err := t.GetTransactionFile(tID)
		if err != nil {
			return err
		}
		switch p := prsr.(type) {
		case *spoe.Parser:
			err = p.Save(tFile)
		case parser.Parser:
			err = p.Save(tFile)
		default:
			return fmt.Errorf("provided parser %s not supported", p)
		}
		if err != nil {
			e := NewConfError(ErrErrorChangingConfig, err.Error())
			if commitImplicit {
				return t.ErrAndDeleteTransaction(e, tID)
			}
			return err
		}
	}
	if commitImplicit {
		if _, err := t.CommitTransaction(tID); err != nil {
			return err
		}
	}
	return nil
}

func (t *Transaction) ErrAndDeleteTransaction(err error, tID string) error {
	// Just a safety to not delete the master files by mistake
	if tID != "" {
		_ = t.DeleteTransaction(tID)
		return err
	}
	return err
}

func (t *Transaction) HandleError(id, parentType, parentName, transactionID string, implicit bool, err error) error {
	var e error
	switch {
	case errors.Is(err, parser_errors.ErrSectionMissing):
		if parentType != "" {
			e = NewConfError(ErrParentDoesNotExist, fmt.Sprintf("%s %s does not exist", parentType, parentName))
		} else {
			e = NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Object %s does not exist", id))
		}
	case errors.Is(err, parser_errors.ErrFromDefaultsSectionMissing):
		e = NewConfError(ErrValidationError, fmt.Sprintf("Object %s references missing defaults section in from", id))
	case errors.Is(err, parser_errors.ErrSectionAlreadyExists):
		e = NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("Object %s already exists", id))
	case errors.Is(err, parser_errors.ErrFetch):
		e = NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Object %v does not exist in %s %s", id, parentType, parentName))
	case errors.Is(err, parser_errors.ErrIndexOutOfRange):
		e = NewConfError(ErrObjectIndexOutOfRange, fmt.Sprintf("Object with id %v in %s %s out of range", id, parentType, parentName))
	default:
		e = err
	}

	if implicit {
		return t.ErrAndDeleteTransaction(e, transactionID)
	}
	return e
}
