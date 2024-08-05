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
	"strconv"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/models"
)

type StructuredGlobal interface {
	GetStructuredGlobalConfiguration(transactionID string) (int64, *models.Global, error)
	PushStructuredGlobalConfiguration(data *models.Global, transactionID string, version int64) error
}

// GetStructuredGlobalConfiguration returns configuration version and global configuration with all its child resources.
// Returns error on fail or if frontend does not exist.
func (c *client) GetStructuredGlobalConfiguration(transactionID string) (int64, *models.Global, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	f, err := parseGlobalSection(p)

	return v, f, err
}

// PushStructuredGlobalConfiguration replaces global configurationand all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) PushStructuredGlobalConfiguration(data *models.Global, transactionID string, version int64) error {
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

	if err = serializeGlobalSection(StructuredToParserArgs{
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

func parseGlobalSection(p parser.Parser) (*models.Global, error) {
	g, err := ParseGlobalSection(p)
	if err != nil {
		return nil, err
	}

	lt, err := ParseLogTargets(GlobalParentName, "", p)
	if err != nil {
		return nil, err
	}
	g.LogTargetList = lt

	return g, nil
}

func serializeGlobalSection(a StructuredToParserArgs, g *models.Global) error {
	p := *a.Parser
	if err := SerializeGlobalSection(p, g, a.Options); err != nil {
		return err
	}

	if err := p.Set(parser.Global, parser.GlobalSectionName, "log", []types.Log{}); err != nil {
		return a.HandleError(GlobalParentName, GlobalParentName, "", a.TID, a.TID == "", err)
	}

	for i, log := range g.LogTargetList {
		if err := p.Insert(parser.Global, parser.GlobalSectionName, "log", SerializeLogTarget(*log), i); err != nil {
			return a.HandleError(strconv.FormatInt(int64(i), 10), GlobalParentName, "", a.TID, a.TID == "", err)
		}
	}

	return nil
}
