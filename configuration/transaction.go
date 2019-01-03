package configuration

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/haproxytech/models"
)

// GetTransactions returns an array of transactions
func (c *LBCTLClient) GetTransactions(status string) (*models.Transactions, error) {
	ts := models.Transactions{}

	response, err := c.executeLBCTL("transaction-list", "")
	if err != nil {
		return nil, err
	}
	for _, id := range strings.Split(response, "\n") {
		if strings.TrimSpace(id) == "" {
			continue
		}
		t, err := c.parseTransaction(id)
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	return &ts, nil
}

// GetTransaction returns transaction information by id
func (c *LBCTLClient) GetTransaction(id string) (*models.Transaction, error) {
	return c.parseTransaction(id)
}

// StartTransaction starts a new empty lbctl transaction
func (c *LBCTLClient) StartTransaction(version int64) (*models.Transaction, error) {
	t := &models.Transaction{}

	err := c.GlobalParser.LoadData(c.ClientParams.GlobalConfigurationFile())
	if err != nil {
		return nil, err
	}

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

	if c.Cache.Enabled() {
		c.Cache.InitTransactionCache(t.ID)
	}

	return t, nil
}

// CommitTransaction commits a transaction by id.
func (c *LBCTLClient) CommitTransaction(id string) error {
	// do a version check before commiting
	version, err := c.GetVersion()
	if err != nil {
		return err
	}
	tVersion := c.getTransactionVersion(id)
	if tVersion != version {
		return NewConfError(ErrVersionMismatch, fmt.Sprintf("Version mismatch, transaction version: %v, configured version: %v", tVersion, version))
	}

	_, err = c.executeLBCTL("transaction-commit", id)
	if err != nil {
		return err
	}

	err = c.incrementVersion()
	if err != nil {
		return err
	}

	if c.Cache.Enabled() {
		c.Cache.DeleteTransactionCache(id)
		c.Cache.InvalidateCache()
	}

	/*err = c.GlobalParser.Save(c.ClientParams.GlobalConfigurationFile())
	if err != nil {
		return err
	}*/

	return nil
}

// DeleteTransaction deletes a transaction by id.
func (c *LBCTLClient) DeleteTransaction(id string) error {
	_, err := c.executeLBCTL("transaction-cancel", id)
	if err != nil {
		return err
	}
	if c.Cache.Enabled() {
		c.Cache.DeleteTransactionCache(id)
	}
	return nil
}

func (c *LBCTLClient) parseTransaction(id string) (*models.Transaction, error) {
	tPath := c.LBCTLTmpPath + "/" + "tmp." + id
	info, err := os.Stat(tPath)

	// check if transaction exists and is a dir
	if os.IsNotExist(err) {
		return nil, NewConfError(ErrTransactionDoesNotExist, "Transaction with id "+id+" does not exist")
	}
	if !info.IsDir() {
		return nil, NewConfError(ErrTransactionDoesNotExist, "Transaction with id "+id+" is not a dir")
	}

	t := &models.Transaction{ID: id, Status: "in_progress"}

	// check status, if it has merged dir it is a failed transaction
	info, err = os.Stat(tPath + "/l7/merged")
	if err == nil && info.IsDir() {
		t.Status = "failed"
	}
	t.Version = c.getTransactionVersion(id)
	// populate operations
	file, err := os.Open(tPath + "/history")
	if err != nil {
		return nil, NewConfError(ErrCannotReadConfFile, fmt.Sprintf("Cannot read operations file for transaction %v: %v", id, err.Error()))
	}
	defer file.Close()

	ops := []*models.TransactionOperationsItems{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		op, err := c.parseOperation(line)
		if err != nil {
			continue
		}
		ops = append(ops, op)
	}
	t.Operations = ops
	return t, nil
}

func (c *LBCTLClient) parseOperation(operationStr string) (*models.TransactionOperationsItems, error) {
	op := &models.TransactionOperationsItems{}
	// Operation and options are split by ^K - rune with ascii code 11
	words := strings.Split(operationStr, string(11))
	if len(words) < 2 {
		return nil, NewConfError(ErrCannotParseTransaction, "Operation "+operationStr+" Cannot be parsed")
	}

	// parse operation
	w := strings.Split(words[0], "_")
	if len(w) < 3 {
		return nil, NewConfError(ErrCannotParseTransaction, "Operation "+operationStr+" Cannot be parsed")
	}
	if w[0] != "l7" {
		return nil, NewConfError(ErrCannotParseTransaction, "Operation "+operationStr+" Cannot be parsed")
	}

	objType := ""
	opType := ""
	parentType := ""
	if len(w) == 3 {
		objType = lbctlTypeToType(w[1])
		opType = lbctlTypeToType(w[2])
	} else {
		parentType = lbctlTypeToType(w[1])
		objType = lbctlTypeToType(w[2])
		opType = lbctlOpToOp(w[3])
	}

	// Get operation data
	op.Operation = opType + objType
	data := make(map[string]string)

	if parentType != "" {
		data["parentType"] = parentType
	}

	var options []string
	if strings.HasPrefix(words[1], "--") && strings.HasPrefix(words[2], "--") {
		data["parent"] = words[1]
		data["name"] = words[2]
		options = words[2:]
	} else if !strings.HasPrefix(words[1], "--") {
		data["name"] = words[1]
		options = words[1:]
	}

	for i := 0; i < len(options); i++ {
		if strings.HasPrefix(options[i], "--reset") {
			continue
		} else if strings.HasPrefix(options[i], "--") {
			fieldName := options[i][2:]
			i++
			fieldValue := options[i]
			data[fieldName] = fieldValue
		}
	}
	op.Data = data

	return op, nil
}

func lbctlOpToOp(lOp string) string {
	if lOp == "update" {
		return "replace"
	}
	return lOp
}

func (c *LBCTLClient) getTransactionVersion(id string) int64 {
	// get original version from the ex file
	tPath := c.LBCTLTmpPath + "/" + "tmp." + id
	file, err := os.Open(tPath + "/l7/ctx/haproxy.cfg.ex")
	if err != nil {
		return 1
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read only first line, version MUST BE on the first line, if not set or malformatted, use 1
	lineNo := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			if lineNo == 0 {
				if strings.HasPrefix(line, "# _version=") {
					w := strings.Split(line, "=")
					if len(w) != 2 {
						return 1
					}
					version, err := strconv.ParseInt(w[1], 10, 64)
					if err != nil {
						return 1
					}
					return version
				}
			} else {
				break
			}
			lineNo++
		}
	}

	return 1
}
