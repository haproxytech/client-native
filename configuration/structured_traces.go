// Copyright 2024 HAProxy Technologies
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
	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/models"
)

type StructuredTraces interface {
	GetStructuredTraces(transactionID string) (int64, *models.Traces, error)
	PushStructuredTraces(data *models.Traces, transactionID string, version int64) error
}

func (c *client) GetStructuredTraces(transactionID string) (int64, *models.Traces, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	traces, err := ParseTraces(p)

	return v, traces, err
}

func (c *client) PushStructuredTraces(data *models.Traces, transactionID string, version int64) error {
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

	// Delete the existing section.
	if c.checkSectionExists(parser.Traces, parser.TracesSectionName, p) {
		if err = p.SectionsDelete(parser.Traces, parser.TracesSectionName); err != nil {
			return c.HandleError(TracesParentName, "", "", t, transactionID == "", err)
		}
	}

	if err = p.SectionsCreate(parser.Traces, parser.TracesSectionName); err != nil {
		return c.HandleError(TracesParentName, "", "", t, transactionID == "", err)
	}

	if err = serializeStructuredTraces(StructuredToParserArgs{
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

func serializeStructuredTraces(a StructuredToParserArgs, traces *models.Traces) error {
	return SerializeTraces(*a.Parser, traces)
}
