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

type StructuredHealthcheck interface {
	GetStructuredHealthchecks(transactionID string) (int64, models.Healthchecks, error)
	GetStructuredHealthcheck(name string, transactionID string) (int64, *models.HealthCheck, error)
	CreateStructuredHealthcheck(data *models.HealthCheck, transactionID string, version int64) error
	EditStructuredHealthcheck(name string, data *models.HealthCheck, transactionID string, version int64) error
}

// GetStructuredHealthcheck returns configuration version and a requested healthcheck with all its child resources.
// Returns error on fail or if healthcheck does not exist.
func (c *client) GetStructuredHealthcheck(name string, transactionID string) (int64, *models.HealthCheck, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !p.SectionExists(parser.HealthChecks, name) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("HealthCheck %s does not exist", name))
	}

	healthCheck := &models.HealthCheck{HealthCheckBase: models.HealthCheckBase{Name: name}}
	err = ParseHealthcheckSection(p, healthCheck)
	if err != nil {
		return v, nil, err
	}

	// parse the HTTP check rules for this health check
	hchecks, err := ParseHTTPChecks(HealthcheckParentName, name, p)
	if err != nil {
		return v, nil, err
	}
	healthCheck.HTTPCheckRuleList = hchecks

	// parse the TCP check rules for this health check
	tchecks, err := ParseTCPChecks(HealthcheckParentName, name, p)
	if err != nil {
		return v, nil, err
	}
	healthCheck.TCPCheckRuleList = tchecks

	return v, healthCheck, nil
}

func (c *client) GetStructuredHealthchecks(transactionID string) (int64, models.Healthchecks, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	hNames, err := p.SectionsGet(parser.HealthChecks)
	if err != nil {
		return v, nil, err
	}

	healthChecks := []*models.HealthCheck{}
	for _, name := range hNames {
		_, healthCheck, err := c.GetHealthcheck(name, transactionID)
		if err == nil {
			// parse the HTTP check rules for this health check
			hchecks, err := ParseHTTPChecks(HealthcheckParentName, name, p)
			if err != nil {
				return v, nil, err
			}
			healthCheck.HTTPCheckRuleList = hchecks

			// parse the TCP check rules for this health check
			tchecks, err := ParseTCPChecks(HealthcheckParentName, name, p)
			if err != nil {
				return v, nil, err
			}
			healthCheck.TCPCheckRuleList = tchecks

			healthChecks = append(healthChecks, healthCheck)
		}
	}

	return v, healthChecks, nil
}

// EditStructuredHealthcheck replaces a healthcheck and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditStructuredHealthcheck(name string, data *models.HealthCheck, transactionID string, version int64) error {
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

	if !p.SectionExists(parser.HealthChecks, name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.HealthChecks, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsDelete(parser.HealthChecks, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err = serializeHealthcheckSection(StructuredToParserArgs{
		TID:         transactionID,
		Parser:      &p,
		Options:     &c.ConfigurationOptions,
		HandleError: c.HandleError,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

// CreateStructuredHealthcheck creates a healthcheck and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateStructuredHealthcheck(data *models.HealthCheck, transactionID string, version int64) error {
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

	if p.SectionExists(parser.HealthChecks, data.Name) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exist", parser.HealthChecks, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = serializeHealthcheckSection(StructuredToParserArgs{
		TID:         transactionID,
		Parser:      &p,
		Options:     &c.ConfigurationOptions,
		HandleError: c.HandleError,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func serializeHealthcheckSection(a StructuredToParserArgs, h *models.HealthCheck) error {
	p := *a.Parser
	var err error
	err = p.SectionsCreate(parser.HealthChecks, h.Name)
	if err != nil {
		return err
	}
	if err = SerializeHealthCheckSection(p, h); err != nil {
		return err
	}

	for i, httpCheck := range h.HTTPCheckRuleList {
		var s types.Action
		s, err = SerializeHTTPCheck(*httpCheck)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.HealthChecks, h.Name, "http-check", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), HealthcheckParentName, h.Name, a.TID, a.TID == "", err)
		}
	}

	for i, tcpCheck := range h.TCPCheckRuleList {
		var s types.Action
		s, err = SerializeTCPCheck(*tcpCheck)
		if err != nil {
			return err
		}
		if err = p.Insert(parser.HealthChecks, h.Name, "tcp-check", s, i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), HealthcheckParentName, h.Name, a.TID, a.TID == "", err)
		}
	}

	return nil
}
