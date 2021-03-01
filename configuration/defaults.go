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
	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v3"

	"github.com/haproxytech/client-native/v2/models"
)

// GetDefaultsConfiguration returns configuration version and a
// struct representing Defaults configuration
func (c *Client) GetDefaultsConfiguration(transactionID string) (int64, *models.Defaults, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	d := &models.Defaults{}
	_ = ParseSection(d, parser.Defaults, parser.DefaultSectionName, p)

	return v, d, nil
}

// PushDefaultsConfiguration pushes a Defaults config struct to global
// config file
func (c *Client) PushDefaultsConfiguration(data *models.Defaults, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	if err := c.editSection(parser.Defaults, parser.DefaultSectionName, data, transactionID, version); err != nil {
		return err
	}

	return nil
}
