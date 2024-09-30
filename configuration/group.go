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
	"errors"
	"fmt"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v5/config-parser"
	parser_errors "github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"

	"github.com/haproxytech/client-native/v5/models"
)

type Group interface {
	GetGroups(userlist string, transactionID string) (int64, models.Groups, error)
	GetGroup(name string, userlist string, transactionID string) (int64, *models.Group, error)
	DeleteGroup(name string, userlist string, transactionID string, version int64) error
	CreateGroup(userlist string, data *models.Group, transactionID string, version int64) error
	EditGroup(name string, userlist string, data *models.Group, transactionID string, version int64) error
}

// GetGroups returns configuration version and an array of configured Group lines in the specified userlist.
// Returns error on fail.
func (c *client) GetGroups(userlist string, transactionID string) (int64, models.Groups, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}
	groups, err := ParseGroups(userlist, p)
	if err != nil {
		return v, nil, c.HandleError("", "userlist", userlist, "", false, err)
	}
	return v, groups, nil
}

// GetGroup returns configuration version and a requested Group line in the specified userlist.
// Returns error on fail or if Group does not exist
func (c *client) GetGroup(name string, userlist string, transactionID string) (int64, *models.Group, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}
	group, _ := GetGroupByName(name, userlist, p)
	if err != nil {
		return v, nil, c.HandleError(name, "userlist", userlist, "", false, err)
	}
	return v, group, nil
}

// DeleteGroup deletes a Group line in the configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success
func (c *client) DeleteGroup(name string, userlist string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}
	if _, _, err := c.GetUserList(userlist, transactionID); err != nil {
		return err
	}
	group, i := GetGroupByName(name, userlist, p)
	if group == nil {
		return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("group %s does not exist", name))
	}
	if err := p.Delete("userlist", userlist, "group", i); err != nil {
		return c.HandleError(name, "userlist", userlist, t, transactionID == "", err)
	}
	return c.SaveData(p, t, transactionID == "")
}

// CreateGroup creates a Group line in the configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success
func (c *client) CreateGroup(userlist string, data *models.Group, transactionID string, version int64) error {
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
	if _, _, err := c.GetUserList(userlist, transactionID); err != nil {
		return err
	}
	group, _ := GetGroupByName(data.Name, userlist, p)
	if group != nil {
		return NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("group %s already exists", data.Name))
	}
	if err := p.Insert("userlist", userlist, "group", SerializeGroup(*data), -1); err != nil {
		return c.HandleError(data.Name, "userlist", userlist, t, transactionID == "", err)
	}
	return c.SaveData(p, t, transactionID == "")
}

// EditGroup edits a Group line in the configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) EditGroup(name string, userlist string, data *models.Group, transactionID string, version int64) error { //nolint:revive
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
	if _, _, err := c.GetUserList(userlist, transactionID); err != nil {
		return err
	}
	group, i := GetGroupByName(data.Name, userlist, p)
	if group == nil {
		return NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("group %s does not exist", data.Name))
	}
	if _, err := p.GetOne("userlist", userlist, "group", i); err != nil {
		return c.HandleError(data.Name, "userlist", userlist, t, transactionID == "", err)
	}
	return c.SaveData(p, t, transactionID == "")
}

func ParseGroups(userlist string, p parser.Parser) (models.Groups, error) {
	groups := models.Groups{}
	data, err := p.Get("userlist", userlist, "group", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return groups, nil
		}
		return nil, err
	}
	items := data.([]types.Group)
	for _, item := range items {
		group := ParseGroup(item)
		groups = append(groups, group)
	}
	return groups, nil
}

func ParseGroup(u types.Group) *models.Group {
	return &models.Group{
		Name:  u.Name,
		Users: strings.Join(u.Users, ","),
	}
}

func SerializeGroup(u models.Group) types.Group {
	var users []string
	if u.Users != "" {
		users = strings.Split(u.Users, ",")
	}

	return types.Group{
		Name:  u.Name,
		Users: users,
	}
}

func GetGroupByName(name string, userlist string, p parser.Parser) (*models.Group, int) {
	groups, err := ParseGroups(userlist, p)
	if err != nil {
		return nil, 0
	}
	for i, group := range groups {
		if group.Name == name {
			return group, i
		}
	}
	return nil, 0
}
