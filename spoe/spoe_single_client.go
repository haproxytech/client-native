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

package spoe

import (
	"fmt"

	parser "github.com/haproxytech/config-parser/v3"
	"github.com/haproxytech/config-parser/v3/spoe"
	"github.com/haproxytech/config-parser/v3/types"

	conf "github.com/haproxytech/client-native/v2/configuration"
	"github.com/haproxytech/client-native/v2/misc"
	"github.com/haproxytech/client-native/v2/models"
)

// SingleSpoe configuration client
// Parser is the SPOE parser instance that loads SPOE configuration file on Init
// and when transaction is committed it gets replaced with the SPOE parser from parsers map
// parsers map contains a SPOE parser for each transaction, which loads data from
// transaction files on StartTransaction, and deletes on CommitTransaction
// We save data to file on every change for persistence
type SingleSpoe struct {
	parsers     map[string]*spoe.Parser
	Parser      *spoe.Parser
	Transaction *conf.Transaction
}

type Params struct {
	SpoeDir                string
	UseValidation          *bool
	PersistentTransactions *bool
	SkipFailedTransactions *bool
	TransactionDir         string
	BackupsNumber          int
	ConfigurationFile      string
}

// newSingleSpoe returns Spoe with default options
func newSingleSpoe(params Params) (*SingleSpoe, error) {
	if params.ConfigurationFile == "" {
		return nil, fmt.Errorf("configuration file missing")
	}
	ss := &SingleSpoe{}
	ss.Transaction = &conf.Transaction{}
	ss.Transaction.TransactionClient = ss
	useValidation := true
	if params.UseValidation != nil {
		useValidation = *params.UseValidation
	}
	persistentTransactions := true
	if params.PersistentTransactions != nil {
		persistentTransactions = *params.PersistentTransactions
	}
	skipFailedTransactions := true
	if params.SkipFailedTransactions != nil {
		skipFailedTransactions = *params.SkipFailedTransactions
	}
	ss.Transaction.ClientParams = conf.ClientParams{
		ConfigurationFile:      params.ConfigurationFile,
		TransactionDir:         params.TransactionDir,
		BackupsNumber:          params.BackupsNumber,
		UseValidation:          useValidation,
		PersistentTransactions: persistentTransactions,
		SkipFailedTransactions: skipFailedTransactions,
	}

	ss.parsers = make(map[string]*spoe.Parser)
	if err := ss.InitTransactionParsers(); err != nil {
		return nil, err
	}

	ss.Parser = &spoe.Parser{}
	if err := ss.Parser.LoadData(params.ConfigurationFile); err != nil {
		return nil, conf.NewConfError(conf.ErrCannotReadConfFile, fmt.Sprintf("cannot read %s", ss.Transaction.ConfigurationFile))
	}

	return ss, nil
}

func (c *SingleSpoe) CheckTransactionOrVersion(transactionID string, version int64) (string, error) {
	return c.Transaction.CheckTransactionOrVersion(transactionID, version)
}

// HasParser checks whether transaction exists in parser
func (c *SingleSpoe) HasParser(transactionID string) bool {
	_, ok := c.parsers[transactionID]
	return ok
}

// GetParserTransactions returns parser transactions
func (c *SingleSpoe) GetParserTransactions() models.Transactions {
	transactions := models.Transactions{}
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
	return transactions
}

// GetParser returns a parser for given transactionID, if transactionID is "", it returns "master" parser
func (c *SingleSpoe) GetParser(transactionID string) (*spoe.Parser, error) {
	if transactionID == "" {
		return c.Parser, nil
	}
	p, ok := c.parsers[transactionID]
	if !ok {
		return nil, conf.NewConfError(conf.ErrTransactionDoesNotExist, fmt.Sprintf("transaction %s does not exist", transactionID))
	}
	return p, nil
}

// AddParser adds parser to parser map
func (c *SingleSpoe) AddParser(transactionID string) error {
	if transactionID == "" {
		return conf.NewConfError(conf.ErrValidationError, "not a valid transaction")
	}
	_, ok := c.parsers[transactionID]
	if ok {
		return conf.NewConfError(conf.ErrTransactionAlreadyExists, fmt.Sprintf("transaction %s already exists", transactionID))
	}

	p := &spoe.Parser{}
	tFile := ""
	var err error
	if c.Transaction.PersistentTransactions {
		tFile, err = c.Transaction.GetTransactionFile(transactionID)
		if err != nil {
			return err
		}
	} else {
		tFile = c.Transaction.ConfigurationFile
	}
	if err := p.LoadData(tFile); err != nil {
		return conf.NewConfError(conf.ErrCannotReadConfFile, fmt.Sprintf("cannot read %s", tFile))
	}
	c.parsers[transactionID] = p
	return nil
}

// DeleteParser deletes parser from parsers map
func (c *SingleSpoe) DeleteParser(transactionID string) error {
	if transactionID == "" {
		return conf.NewConfError(conf.ErrValidationError, "not a valid transaction")
	}
	_, ok := c.parsers[transactionID]
	if !ok {
		return conf.NewConfError(conf.ErrTransactionDoesNotExist, fmt.Sprintf("transaction %s does not exist", transactionID))
	}
	delete(c.parsers, transactionID)
	return nil
}

// CommitParser commits transaction parser, deletes it from parsers map, and replaces master Parser
func (c *SingleSpoe) CommitParser(transactionID string) error {
	if transactionID == "" {
		return conf.NewConfError(conf.ErrValidationError, "not a valid transaction")
	}
	p, ok := c.parsers[transactionID]
	if !ok {
		return conf.NewConfError(conf.ErrTransactionDoesNotExist, fmt.Sprintf("transaction %s does not exist", transactionID))
	}
	c.Parser = p
	delete(c.parsers, transactionID)
	return nil
}

// InitTransactionParsers checks transactions and initializes parsers map with transactions in_progress
func (c *SingleSpoe) InitTransactionParsers() error {
	transactions, err := c.Transaction.GetTransactions("in_progress")
	if err != nil {
		return err
	}

	for _, t := range *transactions {
		if err := c.AddParser(t.ID); err != nil {
			continue
		}
		p, err := c.GetParser(t.ID)
		if err != nil {
			continue
		}
		tFile, err := c.Transaction.GetTransactionFile(t.ID)
		if err != nil {
			return err
		}
		if err := p.LoadData(tFile); err != nil {
			return conf.NewConfError(conf.ErrCannotReadConfFile, fmt.Sprintf("cannot read %s", tFile))
		}
	}
	return nil
}

func (c *SingleSpoe) IncrementVersion() error {
	data, _ := c.Parser.Get("", parser.Comments, parser.CommentsSectionName, "# _version", true)
	ver, _ := data.(*types.ConfigVersion)
	ver.Value++

	if err := c.Parser.Save(c.Transaction.ConfigurationFile); err != nil {
		return conf.NewConfError(conf.ErrCannotSetVersion, fmt.Sprintf("cannot set version: %s", err.Error()))
	}
	return nil
}

func (c *SingleSpoe) IncrementTransactionVersion(transactionID string) error {
	if transactionID == "" {
		return c.incrementTransactionVersion(c.Parser)
	}
	p, err := c.GetParser(transactionID)
	if err != nil {
		return err
	}
	return c.incrementTransactionVersion(p)
}

func (c *SingleSpoe) incrementTransactionVersion(p *spoe.Parser) error {
	data, err := p.Get("", parser.Comments, parser.CommentsSectionName, "# _version", true)
	if err != nil {
		return err
	}
	ver, _ := data.(*types.ConfigVersion)
	ver.Value++
	return nil
}

func (c *SingleSpoe) LoadData(filename string) error {
	err := c.Parser.LoadData(filename)
	if err != nil {
		return conf.NewConfError(conf.ErrCannotReadConfFile, fmt.Sprintf("cannot read %s", filename))
	}
	return nil
}

func (c *SingleSpoe) Save(transactionFile, transactionID string) error {
	if transactionID == "" {
		return c.Parser.Save(transactionFile)
	}
	p, err := c.GetParser(transactionID)
	if err != nil {
		return err
	}
	return p.Save(transactionFile)
}

func (c *SingleSpoe) GetFailedParserTransactionVersion(transactionID string) (int64, error) {
	p := &spoe.Parser{}
	if err := p.LoadData(transactionID); err != nil {
		return 0, conf.NewConfError(conf.ErrCannotReadConfFile, fmt.Sprintf("cannot read %s", transactionID))
	}

	data, _ := p.Get("", parser.Comments, parser.CommentsSectionName, "# _version", false)

	ver, ok := data.(*types.ConfigVersion)
	if !ok {
		return 0, conf.NewConfError(conf.ErrCannotReadVersion, "cannot read version")
	}
	return ver.Value, nil
}

// GetVersion returns configuration file version
func (c *SingleSpoe) GetVersion(transactionID string) (int64, error) {
	return c.getVersion(transactionID)
}

func (c *SingleSpoe) getVersion(transactionID string) (int64, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, conf.NewConfError(conf.ErrCannotReadVersion, fmt.Sprintf("cannot read version: %s", err.Error()))
	}
	data, err := p.Get("", parser.Comments, parser.CommentsSectionName, "# _version", true)
	if err != nil {
		return 0, conf.NewConfError(conf.ErrCannotReadVersion, fmt.Sprintf("cannot read version: %s", err.Error()))
	}
	ver, ok := data.(*types.ConfigVersion)
	if !ok {
		return 0, conf.NewConfError(conf.ErrCannotReadVersion, fmt.Sprintf("cannot read version: %s", err.Error()))
	}
	return ver.Value, nil
}

func (c *SingleSpoe) deleteSection(scope string, section parser.Section, name string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !c.checkSectionExists(scope, section, name, p) {
		e := conf.NewConfError(conf.ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", section, name))
		return c.Transaction.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err := p.SectionsDelete(scope, section, name); err != nil {
		return c.Transaction.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err := c.Transaction.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func (c *SingleSpoe) checkSectionExists(scope string, section parser.Section, sectionName string, p *spoe.Parser) bool {
	sections, err := p.SectionsGet(scope, section)
	if err != nil {
		return false
	}

	if misc.StringInSlice(sectionName, sections) {
		return true
	}
	return false
}

func (c *SingleSpoe) loadDataForChange(transactionID string, version int64) (*spoe.Parser, string, error) {
	t, err := c.CheckTransactionOrVersion(transactionID, version)
	if err != nil {
		// if transaction is implicit, return err and delete transaction
		if transactionID == "" && t != "" {
			return nil, "", c.Transaction.ErrAndDeleteTransaction(err, t)
		}
		return nil, "", err
	}

	p, err := c.GetParser(t)
	if err != nil {
		if transactionID == "" && t != "" {
			return nil, "", c.Transaction.ErrAndDeleteTransaction(err, t)
		}
		return nil, "", err
	}
	return p, t, nil
}
