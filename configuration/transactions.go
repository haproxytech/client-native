package configuration

import (
	"strings"
	"github.com/haproxytech/models"
)

func (self *LBCTLConfigurationClient) GetTransactions(status string) (*models.Transactions, error) {
	// ts := &models.Transactions{}

	// response, err := self.executeLBCTL("transaction-list")
	// for _, id := range(strings.Split("\n")) {
	// 	if strings.TrimSpace(id) == "" {
	// 		continue
	// 	}
	// 	t := &models.Transaction{
	// 		ID: strings.TrimSpace(id),
	// 	}
	// 	ts = append(ts, t)
	// } 
	return &models.Transactions{}, nil
}

func (self *LBCTLConfigurationClient) GetTransaction(id string) (*models.Transaction, error) {
	return &models.Transaction{}, nil
}

func (self *LBCTLConfigurationClient) StartTransaction(version int64) (*models.Transaction, error) {
	t := &models.Transaction{}

	response, err := self.executeLBCTL("transaction-begin", "")
	if err != nil {
		return nil, err
	}

	t.ID = strings.TrimSpace(response)
	t.Version = version
	t.Status = "in_progress"

	return t, nil
}

func (self *LBCTLConfigurationClient) CommitTransaction(id string) error {
	_, err := self.executeLBCTL("-T", id, "transaction-commit")
	if err != nil {
		return err
	}
	err = self.incrementVersion()
	if err != nil {
		return err
	}
	return nil
}