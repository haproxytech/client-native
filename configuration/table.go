// Copyright 2023 HAProxy Technologies
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

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type Table interface {
	GetTables(peerSection string, transactionID string) (int64, models.Tables, error)
	GetTable(name string, peerSection string, transactionID string) (int64, *models.Table, error)
	DeleteTable(name string, peerSection string, transactionID string, version int64) error
	CreateTable(peerSection string, data *models.Table, transactionID string, version int64) error
	EditTable(name string, peerSection string, data *models.Table, transactionID string, version int64) error
}

// GetTables returns configuration version and an array of
// configured tables in the specified peers section. Returns error on fail.
func (c *client) GetTables(peerSection string, transactionID string) (int64, models.Tables, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	Tables, err := ParseTables(peerSection, p)
	if err != nil {
		return v, nil, c.HandleError("", PeersParentName, peerSection, "", false, err)
	}

	return v, Tables, nil
}

// GetTable returns configuration version and a requested table
// in the specified peer section. Returns error on fail or if table does not exist.
func (c *client) GetTable(name string, peerSection string, transactionID string) (int64, *models.Table, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	Table, _ := GetTableByName(name, peerSection, p)
	if Table == nil {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Table %s does not exist in peer section %s", name, peerSection))
	}

	return v, Table, nil
}

// DeleteTable deletes a table in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteTable(name string, peerSection string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	Table, i := GetTableByName(name, peerSection, p)
	if Table == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Table %s does not exist in peer section %s", name, peerSection))
		return c.HandleError(name, PeersParentName, peerSection, t, transactionID == "", e)
	}

	if err := p.Delete(parser.Peers, peerSection, "table", i); err != nil {
		return c.HandleError(name, PeersParentName, peerSection, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// CreateTable creates a table in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateTable(peerSection string, data *models.Table, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	Table, _ := GetTableByName(data.Name, peerSection, p)
	if Table != nil {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("Table %s already exists in peer section %s", data.Name, peerSection))
		return c.HandleError(data.Name, "tables", peerSection, t, transactionID == "", e)
	}

	if err := p.Insert(parser.Peers, peerSection, "table", SerializeTable(*data), -1); err != nil {
		return c.HandleError(data.Name, "tables", peerSection, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditTable edits a table in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditTable(name string, peerSection string, data *models.Table, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	Table, i := GetTableByName(name, peerSection, p)
	if Table == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Table %v does not exist in peer section %s", name, peerSection))
		return c.HandleError(data.Name, PeersParentName, peerSection, t, transactionID == "", e)
	}

	if err := p.Set(parser.Peers, peerSection, "table", SerializeTable(*data), i); err != nil {
		return c.HandleError(data.Name, PeersParentName, peerSection, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseTables(peerSection string, p parser.Parser) (models.Tables, error) {
	var tables models.Tables

	data, err := p.Get(parser.Peers, peerSection, "table", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return tables, nil
		}
		return nil, err
	}

	Tables, ok := data.([]types.Table)
	if !ok {
		return nil, misc.CreateTypeAssertError("table")
	}
	for _, t := range Tables {
		table := ParseTable(t)
		if table != nil {
			tables = append(tables, table)
		}
	}
	return tables, nil
}

func ParseTable(t types.Table) *models.Table {
	table := &models.Table{
		Name:     t.Name,
		Type:     t.Type,
		Size:     t.Size,
		Metadata: parseMetadata(t.Comment),
	}
	if t.Expire != "" {
		table.Expire = &t.Expire
	}
	if t.NoPurge {
		table.NoPurge = t.NoPurge
	}
	if t.RecvOnly {
		table.RecvOnly = true
	}
	if t.Store != "" {
		table.Store = t.Store
	}
	if t.Type != "" {
		table.Type = t.Type
	}
	if t.TypeLen != 0 {
		table.TypeLen = &t.TypeLen
	}
	if t.WriteTo != "" {
		table.WriteTo = &t.WriteTo
	}

	return table
}

func SerializeTable(t models.Table) types.Table {
	comment, err := serializeMetadata(t.Metadata)
	if err != nil {
		comment = ""
	}
	table := types.Table{
		Name:     t.Name,
		Type:     t.Type,
		Size:     t.Size,
		Comment:  comment,
		RecvOnly: t.RecvOnly,
	}
	if t.Expire != nil {
		table.Expire = *t.Expire
	}
	if t.NoPurge {
		table.NoPurge = t.NoPurge
	}
	if t.Store != "" {
		table.Store = t.Store
	}
	if t.Type != "" {
		table.Type = t.Type
	}
	if t.TypeLen != nil {
		table.TypeLen = *t.TypeLen
	}
	if t.WriteTo != nil {
		table.WriteTo = *t.WriteTo
	}
	return table
}

func GetTableByName(name string, peerSection string, p parser.Parser) (*models.Table, int) {
	Tables, err := ParseTables(peerSection, p)
	if err != nil {
		return nil, 0
	}

	for i, b := range Tables {
		if b.Name == name {
			return b, i
		}
	}
	return nil, 0
}
