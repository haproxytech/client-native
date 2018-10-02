package configuration

import (
	"fmt"
	"strings"

	"github.com/haproxytech/models"
)

// GetTransactions returns an array of transactions
func (c *LBCTLConfigurationClient) GetTransactions(status string) (*models.Transactions, error) {
	ts := &models.Transactions{}

	// response, err := c.executeLBCTL("transaction-list", "")
	// if err != nil {
	// 	return nil, err
	// }
	// for _, id := range strings.Split(response, "\n") {
	// 	if strings.TrimSpace(id) == "" {
	// 		continue
	// 	}
	// 	t := &models.Transaction{
	// 		ID: strings.TrimSpace(id),
	// 	}
	// 	ts = append(ts, t)
	// }
	return ts, nil
}

// GetTransaction returns transaction information by id
func (c *LBCTLConfigurationClient) GetTransaction(id string) (*models.Transaction, error) {
	return &models.Transaction{}, nil
}

// StartTransaction starts a new empty lbctl transaction
func (c *LBCTLConfigurationClient) StartTransaction(version int64) (*models.Transaction, error) {
	t := &models.Transaction{}

	v, err := c.GetVersion()
	if err != nil {
		return nil, err
	}

	if version != v {
		return nil, NewConfError(ErrVersionMismatch, fmt.Sprintf("Version in configuration file is %v, given version is %v", v, version))
	}

	response, err := c.executeLBCTL("transaction-begin", "")
	if err != nil {
		return nil, err
	}

	t.ID = strings.TrimSpace(response)
	t.Version = version
	t.Status = "in_progress"

	return t, nil
}

// CommitTransaction commits a transaction by id.
func (c *LBCTLConfigurationClient) CommitTransaction(id string) error {
	// do a version check before commiting
	_, err := c.executeLBCTL("transaction-commit", id)
	if err != nil {
		return err
	}

	err = c.incrementVersion()
	if err != nil {
		return err
	}
	return nil
}

// DeleteTransaction deletes a transaction by id.
func (c *LBCTLConfigurationClient) DeleteTransaction(id string) error {
	_, err := c.executeLBCTL("transaction-cancel", id)
	if err != nil {
		return err
	}
	return nil
}
