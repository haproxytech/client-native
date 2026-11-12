// Copyright 2026 HAProxy Technologies
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
	"strconv"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/parsers"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type FilterSequence interface {
	GetFilterSequences(parentType, parentName string, transactionID string) (int64, models.FilterSequences, error)
	GetFilterSequence(id int64, parentType, parentName string, transactionID string) (int64, *models.FilterSequence, error)
	DeleteFilterSequence(id int64, parentType string, parentName string, transactionID string, version int64) error
	CreateFilterSequence(id int64, parentType string, parentName string, data *models.FilterSequence, transactionID string, version int64) error
	EditFilterSequence(id int64, parentType string, parentName string, data *models.FilterSequence, transactionID string, version int64) error
	ReplaceFilterSequences(parentType string, parentName string, data models.FilterSequences, transactionID string, version int64) error
}

func (c *client) GetFilterSequences(parentType, parentName string, transactionID string) (int64, models.FilterSequences, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	sequences, err := ParseFilterSequences(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, sequences, nil
}

func (c *client) GetFilterSequence(id int64, parentType, parentName string, transactionID string) (int64, *models.FilterSequence, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	section, parentName, err := getParserFromParent(p, "filter-sequence", parentType, parentName)
	if err != nil {
		return v, nil, err
	}

	data, err := p.GetOne(section, parentName, "filter-sequence", int(id))
	if err != nil {
		return v, nil, c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, "", false, err)
	}

	sequence, ok := data.(models.FilterSequence)
	if !ok {
		return v, nil, misc.CreateTypeAssertError("filter-sequence")
	}

	return v, &sequence, nil
}

func (c *client) DeleteFilterSequence(id int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	section, parentName, err := getParserFromParent(p, "filter-sequence", parentType, parentName)
	if err != nil {
		return err
	}

	if err := p.Delete(section, parentName, "filter-sequence", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) CreateFilterSequence(id int64, parentType string, parentName string, data *models.FilterSequence, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent(p, "filter-sequence", parentType, parentName)
	if err != nil {
		return err
	}

	if err := p.Insert(section, parentName, "filter-sequence", *data, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) EditFilterSequence(id int64, parentType string, parentName string, data *models.FilterSequence, transactionID string, version int64) error {
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

	section, parentName, err := getParserFromParent(p, "filter-sequence", parentType, parentName)
	if err != nil {
		return err
	}

	if _, err := p.GetOne(section, parentName, "filter-sequence", int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	if err := p.Set(section, parentName, "filter-sequence", *data, int(id)); err != nil {
		return c.HandleError(strconv.FormatInt(id, 10), parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) ReplaceFilterSequences(parentType string, parentName string, data models.FilterSequences, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	if err := ValidateFilterSequenceList(data); err != nil {
		return NewConfError(ErrValidationError, err.Error())
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	section, parentName, err := getParserFromParent(p, "filter-sequence", parentType, parentName)
	if err != nil {
		return err
	}

	sequences, err := ParseFilterSequences(parentType, parentName, p)
	if err != nil {
		return c.HandleError("", parentType, parentName, "", false, err)
	}

	for i := range sequences {
		if err := p.Delete(section, parentName, "filter-sequence", 0); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	for i, sequence := range data {
		if err := p.Insert(section, parentName, "filter-sequence", *sequence, i); err != nil {
			return c.HandleError(strconv.FormatInt(int64(i), 10), parentType, parentName, t, transactionID == "", err)
		}
	}

	return c.SaveData(p, t, transactionID == "")
}

// ValidateFilterSequenceList validates every sequence line independently and
// delegates to the parser's single cycle implementation.
func ValidateFilterSequenceList(data models.FilterSequences) error {
	for _, sequence := range data {
		if sequence == nil {
			continue
		}
		if err := parsers.ValidateFilterSequence(*sequence); err != nil {
			return err
		}
	}
	return nil
}

func ParseFilterSequences(parentType, parentName string, p parser.Parser) (models.FilterSequences, error) {
	section, parentName, err := getParserFromParent(p, "filter-sequence", parentType, parentName)
	if err != nil {
		return nil, err
	}

	data, err := p.Get(section, parentName, "filter-sequence", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return nil, nil
		}
		return nil, err
	}

	sequences, ok := data.([]models.FilterSequence)
	if !ok {
		return nil, misc.CreateTypeAssertError("filter-sequence")
	}
	result := make(models.FilterSequences, 0, len(sequences))
	for i := range sequences {
		sequence := sequences[i]
		result = append(result, &sequence)
	}
	return result, nil
}
