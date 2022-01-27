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

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v4"

	"github.com/haproxytech/client-native/v3/models"
)

type Userlist interface {
	GetUserLists(transactionID string) (int64, models.Userlists, error)
	GetUserList(name string, transactionID string) (int64, *models.Userlist, error)
	DeleteUserList(name string, transactionID string, version int64) error
	CreateUserList(data *models.Userlist, transactionID string, version int64) error
}

// GetUserlists returns configuration version and an array of configured userlists.
// Returns error on fail.
func (c *client) GetUserLists(transactionID string) (int64, models.Userlists, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}
	names, err := p.SectionsGet(parser.UserList)
	if err != nil {
		return v, nil, err
	}
	userlists := []*models.Userlist{}
	for _, name := range names {
		userlists = append(userlists, &models.Userlist{Name: name})
	}
	return v, userlists, nil
}

// GetUserList returns configuration version and a requested userlist.
// Returns error on fail or if userlist does not exist.
func (c *client) GetUserList(name string, transactionID string) (int64, *models.Userlist, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}
	if !c.checkSectionExists(parser.UserList, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Userlist %s does not exist", name))
	}
	return v, &models.Userlist{Name: name}, nil
}

// DeleteUserList deletes a userlist in configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) DeleteUserList(name string, transactionID string, version int64) error {
	return c.deleteSection(parser.UserList, name, transactionID, version)
}

// CreateUserList creates a userlist in configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) CreateUserList(data *models.Userlist, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}
	if err := c.createSection(parser.UserList, data.Name, data, transactionID, version); err != nil {
		return err
	}
	return nil
}
