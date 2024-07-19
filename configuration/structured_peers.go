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
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v6/models"
	parser "github.com/haproxytech/config-parser/v5"
)

type StructuredPeerSection interface {
	GetStructuredPeerSections(transactionID string) (int64, models.PeerSections, error)
	GetStructuredPeerSection(name string, transactionID string) (int64, *models.PeerSection, error)
	CreateStructuredPeerSection(data *models.PeerSection, transactionID string, version int64) error
	EditStructuredPeerSection(data *models.PeerSection, transactionID string, version int64) error
}

// GetStructuredPeerSection returns configuration version and a requested peer section with all its child resources.
// Returns error on fail or if peer section does not exist.
func (c *client) GetStructuredPeerSection(name string, transactionID string) (int64, *models.PeerSection, error) {
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

	f, err := parsePeerSection(name, p)

	return v, f, err
}

func (c *client) GetStructuredPeerSections(transactionID string) (int64, models.PeerSections, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	peerSections, err := parsePeerSections(p)
	if err != nil {
		return 0, nil, err
	}

	return v, peerSections, nil
}

// EditStructuredPeerSection replaces a peer section and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditStructuredPeerSection(data *models.PeerSection, transactionID string, version int64) error {
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

	if !c.checkSectionExists(parser.Peers, data.Name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.Peers, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.Peers, data.Name); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if err = serializePeerSection(StructuredToParserArgs{
		TID:                transactionID,
		Parser:             &p,
		Options:            &c.ConfigurationOptions,
		HandleError:        c.HandleError,
		CheckSectionExists: c.checkSectionExists,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

// CreateStructuredPeerSection creates a peer section and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateStructuredPeerSection(data *models.PeerSection, transactionID string, version int64) error {
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

	if c.checkSectionExists(parser.Peers, data.Name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s already exist", parser.Peers, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = serializePeerSection(StructuredToParserArgs{
		TID:                transactionID,
		Parser:             &p,
		Options:            &c.ConfigurationOptions,
		HandleError:        c.HandleError,
		CheckSectionExists: c.checkSectionExists,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func parsePeerSections(p parser.Parser) (models.PeerSections, error) {
	names, err := p.SectionsGet(parser.Peers)
	if err != nil {
		return nil, err
	}
	peerSections := []*models.PeerSection{}
	for _, name := range names {
		f, err := parsePeerSection(name, p)
		if err != nil {
			return nil, err
		}
		peerSections = append(peerSections, f)
	}
	return peerSections, nil
}

func parsePeerSection(name string, p parser.Parser) (*models.PeerSection, error) {
	ps := &models.PeerSection{PeerSectionBase: models.PeerSectionBase{Name: name}}
	if err := ParseSection(&ps.PeerSectionBase, parser.Peers, name, p); err != nil {
		return nil, err
	}

	// bind
	b, err := ParseBinds(PeersParentName, name, p)
	if err != nil {
		return nil, err
	}
	ba, errba := namedResourceArrayToMap(b)
	if errba != nil {
		return nil, errba
	}
	ps.Binds = ba

	// log targets
	logTargets, err := ParseLogTargets(PeersParentName, name, p)
	if err != nil {
		return nil, err
	}
	ps.LogTargetList = logTargets
	// peer entries
	entries, err := ParsePeerEntries(name, p)
	if err != nil {
		return nil, err
	}
	entriesa, errea := namedResourceArrayToMap(entries)
	if errea != nil {
		return nil, errea
	}
	ps.PeerEntries = entriesa

	// servers
	servers, err := ParseServers(PeersParentName, name, p)
	if err != nil {
		return nil, err
	}
	serversa, errsa := namedResourceArrayToMap(servers)
	if errsa != nil {
		return nil, errsa
	}
	ps.Servers = serversa
	return ps, nil
}

func serializePeerSection(a StructuredToParserArgs, ps *models.PeerSection) error {
	p := *a.Parser
	var err error
	err = p.SectionsCreate(parser.Peers, ps.Name)
	if err != nil {
		return err
	}
	if err = CreateEditSection(&ps.PeerSectionBase, parser.Peers, ps.Name, p, a.Options); err != nil {
		return a.HandleError(ps.Name, "", "", a.TID, a.TID == "", err)
	}
	for _, entry := range ps.PeerEntries {
		if err = p.Insert(parser.Peers, ps.Name, "peer", SerializePeerEntry(entry), -1); err != nil {
			return a.HandleError(entry.Name, PeersParentName, ps.Name, a.TID, a.TID == "", err)
		}
	}
	for _, bind := range ps.Binds {
		if err = p.Insert(parser.Peers, ps.Name, "bind", SerializeBind(bind), -1); err != nil {
			return a.HandleError(bind.Name, PeersParentName, ps.Name, a.TID, a.TID == "", err)
		}
	}
	for _, server := range ps.Servers {
		if err = p.Insert(parser.Peers, ps.Name, "server", SerializeServer(server, a.Options), -1); err != nil {
			return a.HandleError(server.Name, PeersParentName, ps.Name, a.TID, a.TID == "", err)
		}
	}
	for i, log := range ps.LogTargetList {
		if err = p.Insert(parser.Peers, ps.Name, "log", SerializeLogTarget(*log), i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), PeersParentName, ps.Name, a.TID, a.TID == "", err)
		}
	}
	return nil
}
