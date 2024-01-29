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
	parser "github.com/haproxytech/config-parser/v5"

	"github.com/haproxytech/client-native/v6/models"
)

type PeerSection interface {
	GetPeerSections(transactionID string) (int64, models.PeerSections, error)
	GetPeerSection(name string, transactionID string) (int64, *models.PeerSection, error)
	DeletePeerSection(name string, transactionID string, version int64) error
	CreatePeerSection(data *models.PeerSection, transactionID string, version int64) error
	EditPeerSection(data *models.PeerSection, transactionID string, version int64) error
}

// GetPeerSections returns configuration version and an array of
// configured peer sections. Returns error on fail.
func (c *client) GetPeerSections(transactionID string) (int64, models.PeerSections, error) {
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
		peerSection := &models.PeerSection{Name: name}
		if err := ParseSection(peerSection, parser.Peers, name, p); err != nil {
			continue
		}
		peerSections = append(peerSections, peerSection)
	}

	return v, peerSections, nil
}

// GetPeerSection returns configuration version and a requested peer section.
// Returns error on fail or if peer section does not exist.
func (c *client) GetPeerSection(name string, transactionID string) (int64, *models.PeerSection, error) {
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
	if err := ParseSection(peerSection, parser.Peers, name, p); err != nil {
		return v, nil, err
	}

	return v, peerSection, nil
}

// DeletePeerSection deletes a peerSection in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeletePeerSection(name string, transactionID string, version int64) error {
	return c.deleteSection(parser.Peers, name, transactionID, version)
}

// CreatePeerSection creates a peerSection in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreatePeerSection(data *models.PeerSection, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	return c.createSection(parser.Peers, data.Name, data, transactionID, version)
}

// EditPeerSection edits a peer section in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditPeerSection(data *models.PeerSection, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	return c.editSection(parser.Peers, data.Name, data, transactionID, version)
}
