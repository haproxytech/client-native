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

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type HTTPHardening interface {
	GetHTTPHardenings(transactionID string) (int64, models.HTTPHardenings, error)
	GetHTTPHardening(name, transactionID string) (int64, *models.HTTPHardening, error)
	CreateHTTPHardening(data *models.HTTPHardening, transactionID string, version int64) error
	EditHTTPHardening(name string, data *models.HTTPHardening, transactionID string, version int64) error
	DeleteHTTPHardening(name, transactionID string, version int64) error
}

func (c *client) GetHTTPHardenings(transactionID string) (int64, models.HTTPHardenings, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	names, err := p.SectionsGet(parser.HTTPHardening)
	if err != nil {
		return v, nil, err
	}

	profiles := make(models.HTTPHardenings, len(names))
	for i, name := range names {
		_, profile, err := c.GetHTTPHardening(name, transactionID)
		if err != nil {
			return v, nil, err
		}
		profiles[i] = profile
	}
	return v, profiles, nil
}

func (c *client) GetHTTPHardening(name, transactionID string) (int64, *models.HTTPHardening, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !p.SectionExists(parser.HTTPHardening, name) {
		return v, nil, NewConfError(ErrObjectDoesNotExist,
			fmt.Sprintf("%s section '%s' does not exist", HTTPHardeningParentName, name))
	}

	profile, err := ParseHTTPHardening(p, name)
	if err != nil {
		return 0, nil, err
	}
	return v, profile, nil
}

func (c *client) CreateHTTPHardening(data *models.HTTPHardening, transactionID string, version int64) error {
	if data == nil {
		return NewConfError(ErrValidationError, "empty http-hardening section")
	}
	if data.Name == "auto" {
		return NewConfError(ErrValidationError, `"auto" is reserved and cannot be declared`)
	}
	if c.UseModelsValidation {
		if err := data.Validate(strfmt.Default); err != nil {
			return NewConfError(ErrValidationError, err.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if p.SectionExists(parser.HTTPHardening, data.Name) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.HTTPHardening, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}
	if err = validateHTTPHardeningParent(p, data.Name, data.From); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if err = p.SectionsCreate(parser.HTTPHardening, data.Name); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}
	if err = SerializeHTTPHardening(p, data); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}
	return c.SaveData(p, t, transactionID == "")
}

func (c *client) EditHTTPHardening(name string, data *models.HTTPHardening, transactionID string, version int64) error {
	if data == nil {
		return NewConfError(ErrValidationError, "empty http-hardening section")
	}
	if data.Name == "auto" {
		return NewConfError(ErrValidationError, `"auto" is reserved and cannot be declared`)
	}
	if name != data.Name {
		return NewConfError(ErrValidationError,
			fmt.Sprintf("path name %q does not match body name %q", name, data.Name))
	}
	if c.UseModelsValidation {
		if err := data.Validate(strfmt.Default); err != nil {
			return NewConfError(ErrValidationError, err.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !p.SectionExists(parser.HTTPHardening, name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.HTTPHardening, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}
	if err = validateHTTPHardeningParent(p, name, data.From); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}
	if err = SerializeHTTPHardening(p, data); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}
	return c.SaveData(p, t, transactionID == "")
}

func (c *client) DeleteHTTPHardening(name, transactionID string, version int64) error {
	return c.deleteSection(parser.HTTPHardening, name, transactionID, version)
}

func ParseHTTPHardening(p parser.Parser, name string) (*models.HTTPHardening, error) {
	profile := &models.HTTPHardening{Name: name}

	if data, err := p.SectionGet(parser.HTTPHardening, name); err == nil {
		d, ok := data.(types.Section)
		if !ok {
			return nil, misc.CreateTypeAssertError(HTTPHardeningParentName)
		}
		profile.Metadata = misc.ParseMetadata(d.Comment)
		profile.From = d.From
	} else if !errors.Is(err, parser_errors.ErrFetch) {
		return nil, err
	}

	for keyword, field := range httpHardeningOptionFields(profile) {
		data, err := p.Get(parser.HTTPHardening, name, keyword, false)
		if err != nil {
			if errors.Is(err, parser_errors.ErrFetch) {
				continue
			}
			return nil, err
		}
		option, ok := data.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError(keyword)
		}
		*field = option.Value
	}
	return profile, nil
}

func SerializeHTTPHardening(p parser.Parser, profile *models.HTTPHardening) error {
	if profile == nil {
		return fmt.Errorf("empty %s section", HTTPHardeningParentName)
	}
	if err := validateHTTPHardeningParent(p, profile.Name, profile.From); err != nil {
		return err
	}

	comment, err := misc.SerializeMetadata(profile.Metadata)
	if err != nil {
		return err
	}

	if err := p.SectionsDefaultsFromSet(parser.HTTPHardening, profile.Name, profile.From); err != nil {
		return err
	}
	if err := p.SectionCommentSet(parser.HTTPHardening, profile.Name, comment); err != nil {
		return err
	}

	for keyword, field := range httpHardeningOptionFields(profile) {
		if *field == "" {
			if err := p.Set(parser.HTTPHardening, profile.Name, keyword, nil); err != nil {
				return err
			}
			continue
		}
		value := &types.StringC{Value: *field}
		if err := p.Set(parser.HTTPHardening, profile.Name, keyword, value); err != nil {
			return err
		}
	}
	return nil
}

func httpHardeningOptionFields(profile *models.HTTPHardening) map[string]*string {
	return map[string]*string{
		"all":                                &profile.All,
		"all-reject":                         &profile.AllReject,
		"all-delete":                         &profile.AllDelete,
		"delete-all-trailer-fields":          &profile.DeleteAllTrailerFields,
		"delete-underscore-in-header-name":   &profile.DeleteUnderscoreInHeaderName,
		"delete-underscore-in-trailer-name":  &profile.DeleteUnderscoreInTrailerName,
		"delete-unregistered-header-fields":  &profile.DeleteUnregisteredHeaderFields,
		"delete-unregistered-trailer-fields": &profile.DeleteUnregisteredTrailerFields,
		"reject-content-length-leading-zero": &profile.RejectContentLengthLeadingZero,
		"reject-content-length-with-te":      &profile.RejectContentLengthWithTe,
		"reject-ctl-in-header-trailer":       &profile.RejectCtlInHeaderTrailer,
		"reject-duplicate-host-header":       &profile.RejectDuplicateHostHeader,
		"reject-header-trailer-folding":      &profile.RejectHeaderTrailerFolding,
		"reject-lowercase-method":            &profile.RejectLowercaseMethod,
		"reject-underscore-in-header-name":   &profile.RejectUnderscoreInHeaderName,
		"reject-underscore-in-trailer-name":  &profile.RejectUnderscoreInTrailerName,
	}
}

func validateHTTPHardeningParent(p parser.Parser, name, parent string) error {
	if parent == "" {
		return nil
	}
	if parent == "auto" {
		return NewConfError(ErrValidationError, "auto cannot be used as a parent profile")
	}
	if parent == name {
		return NewConfError(ErrValidationError, fmt.Sprintf("profile %q cannot inherit from itself", name))
	}
	if parent != "trusted" && parent != "untrusted" && !p.SectionExists(parser.HTTPHardening, parent) {
		return NewConfError(ErrValidationError, fmt.Sprintf("parent profile %q does not exist", parent))
	}
	return nil
}
