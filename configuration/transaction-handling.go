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
	"fmt"

	parser "github.com/haproxytech/config-parser/v5"
	parser_options "github.com/haproxytech/config-parser/v5/options"
	"github.com/haproxytech/config-parser/v5/types"

	"github.com/haproxytech/client-native/v6/models"
)

type TransactionHandling interface {
	InitTransactionParsers() error
	GetParserTransactions() models.Transactions
	GetFailedParserTransactionVersion(transactionID string) (int64, error)
	IncrementTransactionVersion(transactionID string) error
}

// InitTransactionParsers checks transactions and initializes parsers map with transactions in_progress
func (c *client) InitTransactionParsers() error {
	transactions, err := c.GetTransactions(models.TransactionStatusInProgress)
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
		tFile, err := c.GetTransactionFile(t.ID)
		if err != nil {
			return err
		}
		if err := p.LoadData(tFile); err != nil {
			return NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read %s", tFile))
		}
	}
	return nil
}

// GetParserTransactionIDs returns parser transactionIDs
func (c *client) GetParserTransactionIDs() []string {
	transactionIDs := []string{}
	c.clientMu.Lock()
	defer c.clientMu.Unlock()
	for tID := range c.parsers {
		transactionIDs = append(transactionIDs, tID)
	}
	return transactionIDs
}

// GetParserTransactions returns parser transactions
func (c *client) GetParserTransactions() models.Transactions {
	transactions := models.Transactions{}
	transactionIDs := c.GetParserTransactionIDs()
	for _, tID := range transactionIDs {
		v, err := c.GetVersion(tID)
		if err == nil {
			t := &models.Transaction{
				ID:      tID,
				Status:  models.TransactionStatusInProgress,
				Version: v,
			}
			transactions = append(transactions, t)
		}
	}
	return transactions
}

func (c *client) GetFailedParserTransactionVersion(transactionID string) (int64, error) {
	p, err := parser.New(parser_options.Path(transactionID))
	if err != nil {
		return 0, NewConfError(ErrCannotReadConfFile, fmt.Sprintf("cannot read %s", transactionID))
	}

	data, _ := p.Get(parser.Comments, parser.CommentsSectionName, "# _version", false)

	ver, ok := data.(*types.ConfigVersion)
	if !ok {
		return 0, NewConfError(ErrCannotReadVersion, "")
	}
	return ver.Value, nil
}

func (c *client) IncrementTransactionVersion(transactionID string) error {
	if transactionID == "" {
		return c.incrementTransactionVersion(c.parser)
	}
	p, err := c.GetParser(transactionID)
	if err != nil {
		return err
	}
	return c.incrementTransactionVersion(p)
}

func (c *client) incrementTransactionVersion(p parser.Parser) error {
	data, err := p.Get(parser.Comments, parser.CommentsSectionName, "# _version", true)
	if err != nil {
		return err
	}
	ver, _ := data.(*types.ConfigVersion)
	ver.Value++
	return nil
}
