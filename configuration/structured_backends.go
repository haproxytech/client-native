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
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/models"
)

type StructuredBackend interface {
	GetStructuredBackends(transactionID string) (int64, models.Backends, error)
	GetStructuredBackend(name string, transactionID string) (int64, *models.Backend, error)
	EditStructuredBackend(name string, data *models.Backend, transactionID string, version int64) error
	CreateStructuredBackend(data *models.Backend, transactionID string, version int64) error
}

// GetStructuredBackend returns configuration version and a requested backend with all its child resources.
// Returns error on fail or if backend does not exist.
func (c *client) GetStructuredBackend(name string, transactionID string) (int64, *models.Backend, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !p.SectionExists(parser.Backends, name) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Backend %s does not exist", name))
	}

	f, err := parseBackendsSection(name, p)

	return v, f, err
}

func (c *client) GetStructuredBackends(transactionID string) (int64, models.Backends, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	backends, err := parseBackendsSections(p)
	if err != nil {
		return 0, nil, err
	}

	return v, backends, nil
}

// EditStructuredBackend replaces a backend and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditStructuredBackend(name string, data *models.Backend, transactionID string, version int64) error {
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

	if !p.SectionExists(parser.Backends, name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.Backends, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.Backends, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err = serializeBackendSection(StructuredToParserArgs{
		TID:         transactionID,
		Parser:      &p,
		Options:     &c.ConfigurationOptions,
		HandleError: c.HandleError,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

// CreateStructuredBackend creates a backend and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateStructuredBackend(data *models.Backend, transactionID string, version int64) error {
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

	if p.SectionExists(parser.Backends, data.Name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s already exist", parser.Backends, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = serializeBackendSection(StructuredToParserArgs{
		TID:         transactionID,
		Parser:      &p,
		Options:     &c.ConfigurationOptions,
		HandleError: c.HandleError,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func parseBackendsSections(p parser.Parser) (models.Backends, error) {
	names, err := p.SectionsGet(parser.Backends)
	if err != nil {
		return nil, err
	}
	backends := []*models.Backend{}
	for _, name := range names {
		f, err := parseBackendsSection(name, p)
		if err != nil {
			return nil, err
		}
		backends = append(backends, f)
	}
	return backends, nil
}

func parseBackendsSection(name string, p parser.Parser) (*models.Backend, error) {
	b := &models.Backend{BackendBase: models.BackendBase{Name: name}}
	if err := ParseSection(&b.BackendBase, parser.Backends, name, p); err != nil {
		return nil, err
	}

	// acls
	acls, err := ParseACLs(BackendParentName, name, p)
	if err != nil {
		return nil, err
	}
	b.ACLList = acls

	// filters
	filters, err := ParseFilters(BackendParentName, name, p)
	if err != nil {
		return nil, err
	}
	b.FilterList = filters

	// http after response rules
	arules, err := ParseHTTPAfterRules(BackendParentName, name, p)
	if err != nil {
		return nil, err
	}
	b.HTTPAfterResponseRuleList = arules

	// http checks
	hchecks, err := ParseHTTPChecks(BackendParentName, name, p)
	if err != nil {
		return nil, err
	}
	b.HTTPCheckList = hchecks

	// http error rules
	errorRules, err := ParseHTTPErrorRules(BackendParentName, name, p)
	if err != nil {
		return nil, err
	}
	b.HTTPErrorRuleList = errorRules

	// http request rules
	rules, err := ParseHTTPRequestRules(BackendParentName, name, p)
	if err != nil {
		return nil, err
	}
	b.HTTPRequestRuleList = rules

	// http response rules
	resprules, err := ParseHTTPResponseRules(BackendParentName, name, p)
	if err != nil {
		return nil, err
	}
	b.HTTPResponseRuleList = resprules

	// log targets
	lt, err := ParseLogTargets(BackendParentName, name, p)
	if err != nil {
		return nil, err
	}
	b.LogTargetList = lt

	// server switching rules
	sr, err := ParseServerSwitchingRules(name, p)
	if err != nil {
		return nil, err
	}
	b.ServerSwitchingRuleList = sr

	// server templates
	st, err := ParseServerTemplates(name, p)
	if err != nil {
		return nil, err
	}
	sta, errCst := namedResourceArrayToMapWithKey(st, "Prefix")
	if errCst != nil {
		return nil, errCst
	}
	b.ServerTemplates = sta

	// servers
	servers, err := ParseServers(BackendParentName, name, p)
	if err != nil {
		return nil, err
	}
	serversa, errC := namedResourceArrayToMap(servers)
	if errC != nil {
		return nil, errC
	}
	b.Servers = serversa

	// stick rules
	stickRules, err := ParseStickRules(name, p)
	if err != nil {
		return nil, err
	}
	b.StickRuleList = stickRules

	// tcp check rules
	tchecks, err := ParseTCPChecks(BackendParentName, name, p)
	if err != nil {
		return nil, err
	}
	b.TCPCheckRuleList = tchecks

	// tcp request rules
	tcpRules, err := ParseTCPRequestRules(BackendParentName, name, p)
	if err != nil {
		return nil, err
	}
	b.TCPRequestRuleList = tcpRules

	// tcp response rules
	tcpRespRules, err := ParseTCPResponseRules(BackendParentName, name, p)
	if err != nil {
		return nil, err
	}
	b.TCPResponseRuleList = tcpRespRules

	return b, nil
}

func serializeBackendSection(a StructuredToParserArgs, b *models.Backend) error { //nolint:gocognit
	p := *a.Parser
	var err error

	err = p.SectionsCreate(parser.Backends, b.Name)
	if err != nil {
		return a.HandleError(b.Name, "", "", a.TID, a.TID == "", err)
	}
	if err = CreateEditSection(&b.BackendBase, parser.Backends, b.Name, p, a.Options); err != nil {
		return a.HandleError(b.Name, "", "", a.TID, a.TID == "", err)
	}

	for _, server := range b.Servers {
		if err = p.Insert(parser.Backends, b.Name, "server", SerializeServer(server, a.Options), -1); err != nil {
			return a.HandleError(server.Name, BackendParentName, b.Name, a.TID, a.TID == "", err)
		}
	}
	for i, filter := range b.FilterList {
		if err = p.Insert(parser.Backends, b.Name, "filter", SerializeFilter(*filter, a.Options), i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, b.Name, a.TID, a.TID == "", err)
		}
	}
	for i, acl := range b.ACLList {
		if err = p.Insert(parser.Backends, b.Name, "acl", *acl, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, b.Name, a.TID, a.TID == "", err)
		}
	}
	for i, rule := range b.ServerSwitchingRuleList {
		if err = p.Insert(parser.Backends, b.Name, "use-server", SerializeServerSwitchingRule(*rule), i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, b.Name, a.TID, a.TID == "", err)
		}
	}
	for _, template := range b.ServerTemplates {
		if err = p.Insert(parser.Backends, b.Name, "server-template", SerializeServerTemplate(template, a.Options), -1); err != nil {
			return a.HandleError(template.Prefix, BackendParentName, b.Name, a.TID, a.TID == "", err)
		}
	}
	for i, rule := range b.HTTPAfterResponseRuleList {
		var s types.Action
		s, err = SerializeHTTPAfterRule(*rule)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Backends, b.Name, "http-after-response", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, b.Name, a.TID, a.TID == "", err)
		}
	}
	for i, httpCheck := range b.HTTPCheckList {
		var s types.Action
		s, err = SerializeHTTPCheck(*httpCheck)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Backends, b.Name, "http-check", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, "", a.TID, a.TID == "", err)
		}
	}
	for i, rule := range b.HTTPRequestRuleList {
		var s types.Action
		s, err = SerializeHTTPRequestRule(*rule, a.Options)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Backends, b.Name, "http-request", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, b.Name, a.TID, a.TID == "", err)
		}
	}
	for i, rule := range b.HTTPResponseRuleList {
		var s types.Action
		s, err = SerializeHTTPResponseRule(*rule, a.Options)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Backends, b.Name, "http-response", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, b.Name, a.TID, a.TID == "", err)
		}
	}
	for i, httpErrorRule := range b.HTTPErrorRuleList {
		var s types.Action
		s, err = SerializeHTTPErrorRule(*httpErrorRule)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Backends, b.Name, "http-error", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, b.Name, a.TID, a.TID == "", err)
		}
	}
	for i, log := range b.LogTargetList {
		if err = p.Insert(parser.Backends, b.Name, "log", SerializeLogTarget(*log), i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, b.Name, a.TID, a.TID == "", err)
		}
	}
	for i, tcpCheck := range b.TCPCheckRuleList {
		var s types.Action
		s, err = SerializeTCPCheck(*tcpCheck)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Backends, b.Name, "tcp-check", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, "", a.TID, a.TID == "", err)
		}
	}
	for i, rule := range b.TCPRequestRuleList {
		var s types.TCPType
		s, err = SerializeTCPRequestRule(*rule, a.Options)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Backends, b.Name, "tcp-request", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, b.Name, a.TID, a.TID == "", err)
		}
	}
	for i, rule := range b.TCPResponseRuleList {
		var s types.TCPType
		s, err = SerializeTCPResponseRule(*rule, a.Options)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Backends, b.Name, "tcp-response", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, b.Name, a.TID, a.TID == "", err)
		}
	}
	for i, rule := range b.StickRuleList {
		if err = p.Insert(parser.Backends, b.Name, "stick", SerializeStickRule(*rule), i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), BackendParentName, b.Name, a.TID, a.TID == "", err)
		}
	}

	return nil
}
