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
	"github.com/haproxytech/client-native/v6/models"
)

type StructuredHTTPHardening interface {
	GetStructuredHTTPHardenings(transactionID string) (int64, models.HTTPHardenings, error)
	GetStructuredHTTPHardening(name, transactionID string) (int64, *models.HTTPHardening, error)
	CreateStructuredHTTPHardening(data *models.HTTPHardening, transactionID string, version int64) error
	EditStructuredHTTPHardening(name string, data *models.HTTPHardening, transactionID string, version int64) error
}

func (c *client) GetStructuredHTTPHardenings(transactionID string) (int64, models.HTTPHardenings, error) {
	return c.GetHTTPHardenings(transactionID)
}

func (c *client) GetStructuredHTTPHardening(name, transactionID string) (int64, *models.HTTPHardening, error) {
	return c.GetHTTPHardening(name, transactionID)
}

func (c *client) CreateStructuredHTTPHardening(data *models.HTTPHardening, transactionID string, version int64) error {
	return c.CreateHTTPHardening(data, transactionID, version)
}

func (c *client) EditStructuredHTTPHardening(name string, data *models.HTTPHardening, transactionID string, version int64) error {
	return c.EditHTTPHardening(name, data, transactionID, version)
}
