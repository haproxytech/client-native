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

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v3"
	parser_errors "github.com/haproxytech/config-parser/v3/errors"
	"github.com/haproxytech/config-parser/v3/types"

	"github.com/haproxytech/client-native/v2/models"
)

// GetPeerEntries returns configuration version and an array of
// configured binds in the specified peers section. Returns error on fail.
func (c *Client) GetPeerEntries(peerSection string, transactionID string) (int64, models.PeerEntries, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	peerEntries, err := ParsePeerEntries(peerSection, p)
	if err != nil {
		return v, nil, c.HandleError("", "peers", peerSection, "", false, err)
	}

	return v, peerEntries, nil
}

// GetPeerEntry returns configuration version and a requested peer entry
// in the specified peer section. Returns error on fail or if bind does not exist.
func (c *Client) GetPeerEntry(name string, peerSection string, transactionID string) (int64, *models.PeerEntry, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	peerEntry, _ := GetPeerEntryByName(name, peerSection, p)
	if peerEntry == nil {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("PeerEntry %s does not exist in peer section %s", name, peerSection))
	}

	return v, peerEntry, nil
}

// DeletePeerEntry deletes an peer entry in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeletePeerEntry(name string, peerSection string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	peerEntry, i := GetPeerEntryByName(name, peerSection, p)
	if peerEntry == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("PeerEntry %s does not exist in peer section %s", name, peerSection))
		return c.HandleError(name, "peers", peerSection, t, transactionID == "", e)
	}

	if err := p.Delete(parser.Peers, peerSection, "peer", i); err != nil {
		return c.HandleError(name, "peers", peerSection, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// CreatePeerEntry creates a peer entry in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreatePeerEntry(peerSection string, data *models.PeerEntry, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	peerEntry, _ := GetPeerEntryByName(data.Name, peerSection, p)
	if peerEntry != nil {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("PeerEntry %s already exists in peer section %s", data.Name, peerSection))
		return c.HandleError(data.Name, "peers", peerSection, t, transactionID == "", e)
	}

	if err := p.Insert(parser.Peers, peerSection, "peer", SerializePeerEntry(*data), -1); err != nil {
		return c.HandleError(data.Name, "peers", peerSection, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// EditPeerEntry edits a peer entry in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditPeerEntry(name string, peerSection string, data *models.PeerEntry, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	peerEntry, i := GetPeerEntryByName(name, peerSection, p)
	if peerEntry == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("PeerEntry %v does not exist in peer section %s", name, peerSection))
		return c.HandleError(data.Name, "peers", peerSection, t, transactionID == "", e)
	}

	if err := p.Set(parser.Peers, peerSection, "peer", SerializePeerEntry(*data), i); err != nil {
		return c.HandleError(data.Name, "peers", peerSection, t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func ParsePeerEntries(peerSection string, p *parser.Parser) (models.PeerEntries, error) {
	peerEntry := models.PeerEntries{}

	data, err := p.Get(parser.Peers, peerSection, "peer", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return peerEntry, nil
		}
		return nil, err
	}

	peerEntries := data.([]types.Peer)
	for _, e := range peerEntries {
		pe := ParsePeerEntry(e)
		if pe != nil {
			peerEntry = append(peerEntry, pe)
		}
	}
	return peerEntry, nil
}

func ParsePeerEntry(p types.Peer) *models.PeerEntry {
	return &models.PeerEntry{
		Address: &p.IP,
		Port:    &p.Port,
		Name:    p.Name,
	}
}

func SerializePeerEntry(pe models.PeerEntry) types.Peer {
	return types.Peer{
		Name: pe.Name,
		IP:   *pe.Address,
		Port: *pe.Port,
	}
}

func GetPeerEntryByName(name string, peerSection string, p *parser.Parser) (*models.PeerEntry, int) {
	peerEntries, err := ParsePeerEntries(peerSection, p)
	if err != nil {
		return nil, 0
	}

	for i, b := range peerEntries {
		if b.Name == name {
			return b, i
		}
	}
	return nil, 0
}
