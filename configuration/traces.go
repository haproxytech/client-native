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
	"errors"
	"fmt"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/common"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type Traces interface {
	GetTraces(transactionID string) (int64, *models.Traces, error)
	CreateTraces(data *models.Traces, transactionID string, version int64) error
	EditTraces(data *models.Traces, transactionID string, version int64) error
	DeleteTraces(transactionID string, version int64) error
	CreateTraceEntry(data *models.TraceEntry, transactionID string, version int64) error
	DeleteTraceEntry(data *models.TraceEntry, transactionID string, version int64) error
}

func (c *client) GetTraces(transactionID string) (int64, *models.Traces, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	traces, err := ParseTraces(p)
	if err != nil {
		return v, nil, c.HandleError(TracesParentName, "", "", transactionID, transactionID == "", err)
	}

	return v, traces, err
}

func (c *client) CreateTraces(data *models.Traces, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return c.HandleError(TracesParentName, "", "", t, transactionID == "", err)
	}

	if c.checkSectionExists(parser.Traces, parser.TracesSectionName, p) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s section already exists", parser.Traces))
		return c.HandleError(TracesParentName, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(parser.Traces, parser.TracesSectionName); err != nil {
		return c.HandleError(TracesParentName, "", "", t, transactionID == "", err)
	}

	if err = SerializeTraces(p, data); err != nil {
		return c.HandleError(TracesParentName, "", "", t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) EditTraces(data *models.Traces, transactionID string, version int64) error {
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

	if err = SerializeTraces(p, data); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) DeleteTraces(transactionID string, version int64) error {
	return c.deleteSection(parser.Traces, parser.TracesSectionName, transactionID, version)
}

// Add an entry to the traces section.
func (c *client) CreateTraceEntry(data *models.TraceEntry, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return c.HandleError("trace", TracesParentName, "", t, transactionID == "", err)
	}

	// Count existing entries. This Get() will create the section if needed.
	traceEntries, err := p.Get(parser.Traces, parser.TracesSectionName, "trace", true)
	if err != nil {
		if !errors.Is(err, parser_errors.ErrFetch) {
			return err
		}
	}
	entries, ok := traceEntries.([]types.Trace)
	if !ok {
		return misc.CreateTypeAssertError("trace entries")
	}
	i := len(entries)

	err = p.Insert(parser.Traces, parser.TracesSectionName, "trace", convertTraceEntry(data), i)
	if err != nil {
		return c.HandleError("trace", TracesParentName, "", t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

// Delete a trace entry from the traces section.
func (c *client) DeleteTraceEntry(data *models.TraceEntry, transactionID string, version int64) error {
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

	if !c.checkSectionExists(parser.Traces, parser.TracesSectionName, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s section does not exists", parser.Traces))
		return c.HandleError(TracesParentName, "", "", t, transactionID == "", e)
	}

	// Look for the trace entry to remove
	traces, err := ParseTraces(p)
	if err != nil {
		return c.HandleError(TracesParentName, "", "", t, transactionID == "", err)
	}
	for i, ent := range traces.Entries {
		if ent.Trace == data.Trace {
			err = p.Delete(parser.Traces, parser.TracesSectionName, "trace", i)
			if err != nil {
				return c.HandleError(TracesParentName, "", "", t, transactionID == "", err)
			}
			return c.SaveData(p, t, transactionID == "")
		}
	}

	// Entry not found.
	e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("trace entry does not exists: '%s'", data.Trace))
	return c.HandleError(TracesParentName, "", "", t, transactionID == "", e)
}

func ParseTraces(p parser.Parser) (*models.Traces, error) {
	traces := new(models.Traces)

	traceEntries, err := p.Get(parser.Traces, parser.TracesSectionName, "trace", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return traces, nil
		}
		return nil, err
	}

	entries, ok := traceEntries.([]types.Trace)
	if !ok {
		return nil, misc.CreateTypeAssertError("trace entries")
	}
	if len(entries) == 0 {
		return traces, nil
	}

	traces.Entries = make(models.TraceEntries, len(entries))
	for i, t := range entries {
		traces.Entries[i] = &models.TraceEntry{Trace: strings.Join(t.Params, " ")}
	}

	return traces, nil
}

func SerializeTraces(p parser.Parser, traces *models.Traces) error {
	if traces == nil {
		return fmt.Errorf("empty %s section", TracesParentName)
	}

	for i, t := range traces.Entries {
		err := p.Insert(parser.Traces, parser.TracesSectionName, "trace", convertTraceEntry(t), i)
		if err != nil {
			return err
		}
	}

	return nil
}

func convertTraceEntry(entry *models.TraceEntry) types.Trace {
	result := types.Trace{}
	if entry != nil {
		result.Params = common.StringSplitIgnoreEmpty(entry.Trace, ' ', '	')
	}
	return result
}
