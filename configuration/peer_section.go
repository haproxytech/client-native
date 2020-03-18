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

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v2"
	"github.com/haproxytech/models"
)

// GetPeerSections returns configuration version and an array of
// configured peer sections. Returns error on fail.
func (c *Client) GetPeerSections(transactionID string) (int64, models.PeerSections, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	names, err := p.SectionsGet(parser.Peers)
	if err != nil {
		return v, nil, err
	}

	peerSections := []*models.PeerSection{}
	for _, name := range names {
		f := &models.PeerSection{Name: name}
		peerSections = append(peerSections, f)
	}

	return v, peerSections, nil
}

// GetPeerSection returns configuration version and a requested peer section.
// Returns error on fail or if peer section does not exist.
func (c *Client) GetPeerSection(name string, transactionID string) (int64, *models.PeerSection, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.Peers, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("PeerSection %s does not exist", name))
	}

	peerSection := &models.PeerSection{Name: name}

	return v, peerSection, nil
}

// DeletePeerSection deletes a peerSection in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeletePeerSection(name string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !c.checkSectionExists(parser.Peers, name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.Peers, name))
		return c.handleError(name, "", "", t, transactionID == "", e)
	}

	if err := p.SectionsDelete(parser.Peers, name); err != nil {
		return c.handleError(name, "", "", t, transactionID == "", err)
	}

	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// EditPeerSection edits a peerSection in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditPeerSection(name string, data *models.PeerSection, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	if err := c.DeletePeerSection(name, transactionID, version); err != nil {
		return err
	}
	if err := c.CreatePeerSection(data, transactionID, version); err != nil {
		return err
	}

	return nil
}

// CreatePeerSection creates a peerSection in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreatePeerSection(data *models.PeerSection, transactionID string, version int64) error {
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
	if err := p.SectionsCreate(parser.Peers, data.Name); err != nil {
		return c.handleError(data.Name, "", "", t, transactionID == "", err)
	}
	if err := c.saveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}
