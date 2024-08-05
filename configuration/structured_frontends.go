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

type StructuredFrontend interface {
	GetStructuredFrontends(transactionID string) (int64, models.Frontends, error)
	GetStructuredFrontend(name string, transactionID string) (int64, *models.Frontend, error)
	EditStructuredFrontend(name string, data *models.Frontend, transactionID string, version int64) error
	CreateStructuredFrontend(data *models.Frontend, transactionID string, version int64) error
}

// GetStructuredFrontend returns configuration version and a requested frontend with all its child resources.
// Returns error on fail or if frontend does not exist.
func (c *client) GetStructuredFrontend(name string, transactionID string) (int64, *models.Frontend, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.Frontends, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Frontend %s does not exist", name))
	}

	f, err := parseFrontendsSection(name, p)

	return v, f, err
}

func (c *client) GetStructuredFrontends(transactionID string) (int64, models.Frontends, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	frontends, err := parseFrontendsSections(p)
	if err != nil {
		return 0, nil, err
	}

	return v, frontends, nil
}

// EditStructuredFrontend replaces a frontend and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditStructuredFrontend(name string, data *models.Frontend, transactionID string, version int64) error {
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

	if !c.checkSectionExists(parser.Frontends, name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.Frontends, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.Frontends, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err = serializeFrontendSection(StructuredToParserArgs{
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

// CreateStructuredFrontend creates a frontend and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateStructuredFrontend(data *models.Frontend, transactionID string, version int64) error {
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

	if c.checkSectionExists(parser.Frontends, data.Name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s already exist", parser.Frontends, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = serializeFrontendSection(StructuredToParserArgs{
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

func parseFrontendsSections(p parser.Parser) (models.Frontends, error) {
	names, err := p.SectionsGet(parser.Frontends)
	if err != nil {
		return nil, err
	}
	frontends := []*models.Frontend{}
	for _, name := range names {
		f, err := parseFrontendsSection(name, p)
		if err != nil {
			return nil, err
		}
		frontends = append(frontends, f)
	}
	return frontends, nil
}

func parseFrontendsSection(name string, p parser.Parser) (*models.Frontend, error) {
	f := &models.Frontend{FrontendBase: models.FrontendBase{Name: name}}
	if err := ParseSection(&f.FrontendBase, parser.Frontends, name, p); err != nil {
		return nil, err
	}

	// acl
	acls, err := ParseACLs(FrontendParentName, name, p)
	if err != nil {
		return nil, err
	}
	f.ACLList = acls

	// backend switching rule
	bsr, err := ParseBackendSwitchingRules(name, p)
	if err != nil {
		return nil, err
	}
	f.BackendSwitchingRuleList = bsr

	// bind
	binds, err := ParseBinds(FrontendParentName, name, p)
	if err != nil {
		return nil, err
	}
	b, err := namedResourceArrayToMap(binds)
	if err != nil {
		return nil, err
	}
	f.Binds = b

	// captures
	captures, err := ParseDeclareCaptures(name, p)
	if err != nil {
		return nil, err
	}
	f.CaptureList = captures

	// filters
	filters, err := ParseFilters(FrontendParentName, name, p)
	if err != nil {
		return nil, err
	}
	f.FilterList = filters

	// http after response rules
	httpAfterResponseRules, err := ParseHTTPAfterRules(FrontendParentName, name, p)
	if err != nil {
		return nil, err
	}
	f.HTTPAfterResponseRuleList = httpAfterResponseRules

	// http error rules
	httpErrorRules, err := ParseHTTPErrorRules(FrontendParentName, name, p)
	if err != nil {
		return nil, err
	}
	f.HTTPErrorRuleList = httpErrorRules

	// http request rules
	httpRequestRules, err := ParseHTTPRequestRules(FrontendParentName, name, p)
	if err != nil {
		return nil, err
	}
	f.HTTPRequestRuleList = httpRequestRules

	// http response rules
	httpResponseRules, err := ParseHTTPResponseRules(FrontendParentName, name, p)
	if err != nil {
		return nil, err
	}
	f.HTTPResponseRuleList = httpResponseRules

	// log targets
	logTargets, err := ParseLogTargets(FrontendParentName, name, p)
	if err != nil {
		return nil, err
	}
	f.LogTargetList = logTargets

	// tcp request rules
	tcpRequestRules, err := ParseTCPRequestRules(FrontendParentName, name, p)
	if err != nil {
		return nil, err
	}
	f.TCPRequestRuleList = tcpRequestRules

	return f, nil
}

func serializeFrontendSection(a StructuredToParserArgs, f *models.Frontend) error { //nolint:gocognit
	p := *a.Parser
	var err error

	err = p.SectionsCreate(parser.Frontends, f.Name)
	if err != nil {
		return a.HandleError(f.Name, "", "", a.TID, a.TID == "", err)
	}
	if err = CreateEditSection(&f.FrontendBase, parser.Frontends, f.Name, p, a.Options); err != nil {
		return a.HandleError(f.Name, "", "", a.TID, a.TID == "", err)
	}
	for _, bind := range f.Binds {
		if err = p.Insert(parser.Frontends, f.Name, "bind", SerializeBind(bind), -1); err != nil {
			return a.HandleError(bind.Name, FrontendParentName, f.Name, a.TID, a.TID == "", err)
		}
	}
	for i, filter := range f.FilterList {
		if err = p.Insert(parser.Frontends, f.Name, "filter", SerializeFilter(*filter, a.Options), i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), FrontendParentName, f.Name, a.TID, a.TID == "", err)
		}
	}
	for i, acl := range f.ACLList {
		if err = p.Insert(parser.Frontends, f.Name, "acl", SerializeACL(*acl), i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), FrontendParentName, f.Name, a.TID, a.TID == "", err)
		}
	}
	for i, capture := range f.CaptureList {
		if err = p.Insert(parser.Frontends, f.Name, "declare capture", SerializeDeclareCapture(*capture), i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), DefaultsParentName, "", a.TID, a.TID == "", err)
		}
	}
	for i, bsr := range f.BackendSwitchingRuleList {
		if err = p.Insert(parser.Frontends, f.Name, "use_backend", SerializeBackendSwitchingRule(*bsr), i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), FrontendParentName, f.Name, a.TID, a.TID == "", err)
		}
	}
	for i, rule := range f.HTTPRequestRuleList {
		var s types.Action
		s, err = SerializeHTTPRequestRule(*rule, a.Options)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Frontends, f.Name, "http-request", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), FrontendParentName, f.Name, a.TID, a.TID == "", err)
		}
	}
	for i, rule := range f.HTTPAfterResponseRuleList {
		var s types.Action
		s, err = SerializeHTTPAfterRule(*rule)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Frontends, f.Name, "http-after-response", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), FrontendParentName, f.Name, a.TID, a.TID == "", err)
		}
	}
	for i, rule := range f.HTTPResponseRuleList {
		var s types.Action
		s, err = SerializeHTTPResponseRule(*rule, a.Options)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Frontends, f.Name, "http-response", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), FrontendParentName, f.Name, a.TID, a.TID == "", err)
		}
	}
	for i, httpErrorRule := range f.HTTPErrorRuleList {
		var s types.Action
		s, err = SerializeHTTPErrorRule(*httpErrorRule)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Frontends, f.Name, "http-error", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), FrontendParentName, f.Name, a.TID, a.TID == "", err)
		}
	}
	for i, log := range f.LogTargetList {
		if err = p.Insert(parser.Frontends, f.Name, "log", SerializeLogTarget(*log), i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), FrontendParentName, f.Name, a.TID, a.TID == "", err)
		}
	}
	for i, rule := range f.TCPRequestRuleList {
		var s types.TCPType
		s, err = SerializeTCPRequestRule(*rule, a.Options)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.Frontends, f.Name, "tcp-request", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), FrontendParentName, f.Name, a.TID, a.TID == "", err)
		}
	}

	return nil
}
