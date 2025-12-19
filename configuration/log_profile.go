// Copyright 2025 HAProxy Technologies
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
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type LogProfile interface {
	GetLogProfiles(transactionID string) (int64, models.LogProfiles, error)
	GetLogProfile(name, transactionID string) (int64, *models.LogProfile, error)
	DeleteLogProfile(name, transactionID string, version int64) error
	CreateLogProfile(data *models.LogProfile, transactionID string, version int64) error
	EditLogProfile(name string, data *models.LogProfile, transactionID string, version int64) error
}

func (c *client) GetLogProfiles(transactionID string) (int64, models.LogProfiles, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	names, err := p.SectionsGet(parser.LogProfile)
	if err != nil {
		return v, nil, err
	}

	lps := make([]*models.LogProfile, len(names))

	for i, name := range names {
		_, lp, err := c.GetLogProfile(name, transactionID)
		if err != nil {
			return v, nil, err
		}
		lps[i] = lp
	}

	return v, lps, nil
}

func (c *client) GetLogProfile(name, transactionID string) (int64, *models.LogProfile, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !p.SectionExists(parser.LogProfile, name) {
		return v, nil, NewConfError(ErrObjectDoesNotExist,
			fmt.Sprintf("%s section '%s' does not exist", LogProfileParentName, name))
	}

	lp, err := ParseLogProfile(p, name)
	if err != nil {
		return 0, nil, err
	}

	return v, lp, nil
}

func (c *client) DeleteLogProfile(name, transactionID string, version int64) error {
	return c.deleteSection(parser.LogProfile, name, transactionID, version)
}

func (c *client) CreateLogProfile(data *models.LogProfile, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if p.SectionExists(parser.LogProfile, data.Name) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.LogProfile, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(parser.LogProfile, data.Name); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if err = SerializeLogProfile(p, data); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) EditLogProfile(name string, data *models.LogProfile, transactionID string, version int64) error { //nolint:revive
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

	if !p.SectionExists(parser.LogProfile, data.Name) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s does not exists", parser.LogProfile, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = SerializeLogProfile(p, data); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseLogProfile(p parser.Parser, name string) (*models.LogProfile, error) {
	lp := &models.LogProfile{Name: name}

	if data, err := p.SectionGet(parser.LogProfile, name); err == nil {
		d, ok := data.(types.Section)
		if ok {
			lp.Metadata = misc.ParseMetadata(d.Comment)
		}
	}
	// get optional log-tag
	logTag, err := p.Get(parser.LogProfile, name, "log-tag", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return lp, nil
		}
		return nil, err
	}
	sc, ok := logTag.(*types.StringC)
	if !ok {
		return nil, misc.CreateTypeAssertError("log-tag")
	}
	lp.LogTag = sc.Value

	// get optional step handlers ("on")
	steps, err := p.Get(parser.LogProfile, name, "on", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return lp, nil
		}
		return nil, err
	}
	tsteps, ok := steps.([]types.OnLogStep)
	if !ok {
		return nil, misc.CreateTypeAssertError("log-profile on")
	}

	lp.Steps = make(models.LogProfileSteps, len(tsteps))

	for i, s := range tsteps {
		step := &models.LogProfileStep{
			Step:     s.Step,
			Drop:     models.LogProfileStepDropDisabled,
			Format:   strings.Trim(s.Format, `"`),
			Sd:       strings.Trim(s.Sd, `"`),
			Metadata: misc.ParseMetadata(s.Comment),
		}
		if s.Drop {
			step.Drop = models.LogProfileStepDropEnabled
			step.Format = ""
			step.Sd = ""
		}
		lp.Steps[i] = step
	}

	return lp, nil
}

func SerializeLogProfile(p parser.Parser, lp *models.LogProfile) error {
	if lp == nil {
		return fmt.Errorf("empty %s section", LogProfileParentName)
	}
	if lp.Metadata != nil {
		comment, err := misc.SerializeMetadata(lp.Metadata)
		if err != nil {
			return err
		}
		if err := p.SectionCommentSet(parser.LogForward, lp.Name, comment); err != nil {
			return err
		}
	}
	logTag := types.StringC{Value: lp.LogTag}
	if err := p.Set(parser.LogProfile, lp.Name, "log-tag", logTag); err != nil {
		return err
	}

	if len(lp.Steps) == 0 {
		_ = p.Delete(parser.LogProfile, lp.Name, "on")
		return nil
	}

	newSteps := make([]types.OnLogStep, len(lp.Steps))
	for i, step := range lp.Steps {
		newSteps[i] = *SerializeLogProfileStep(step)
	}

	return p.Set(parser.LogProfile, lp.Name, "on", newSteps)
}

func SerializeLogProfileStep(step *models.LogProfileStep) *types.OnLogStep {
	// This is way too simplistic.
	q := func(s string) string {
		if strings.ContainsRune(s, ' ') {
			return `"` + s + `"`
		}
		return s
	}

	comment, err := misc.SerializeMetadata(step.Metadata)
	if err != nil {
		comment = ""
	}

	return &types.OnLogStep{
		Step:    step.Step,
		Drop:    step.Drop == models.LogProfileStepDropEnabled,
		Format:  q(step.Format),
		Sd:      q(step.Sd),
		Comment: comment,
	}
}
