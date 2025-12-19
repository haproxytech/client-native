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

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/models"
)

type StructuredUserlist interface {
	GetStructuredUserLists(transactionID string) (int64, models.Userlists, error)
	GetStructuredUserList(name string, transactionID string) (int64, *models.Userlist, error)
	CreateStructuredUserList(data *models.Userlist, transactionID string, version int64) error
}

// GetStructuredUserList returns configuration version and a requested userlist with all its child resources.
// Returns error on fail or if userlist does not exist.
func (c *client) GetStructuredUserList(name string, transactionID string) (int64, *models.Userlist, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !p.SectionExists(parser.UserList, name) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Userlist %s does not exist", name))
	}

	f, err := parseUserlistsSection(name, p)

	return v, f, err
}

func (c *client) GetStructuredUserLists(transactionID string) (int64, models.Userlists, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	userlists, err := parseUserlistsSections(p)
	if err != nil {
		return 0, nil, err
	}

	return v, userlists, nil
}

// CreateStructuredUserList creates a userlist and all it's child resources in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateStructuredUserList(data *models.Userlist, transactionID string, version int64) error {
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

	if p.SectionExists(parser.UserList, data.Name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s already exist", parser.UserList, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = serializeUserlistSection(StructuredToParserArgs{
		TID:         transactionID,
		Parser:      &p,
		Options:     &c.ConfigurationOptions,
		HandleError: c.HandleError,
	}, data); err != nil {
		return err
	}
	return c.SaveData(p, t, transactionID == "")
}

func parseUserlistsSections(p parser.Parser) (models.Userlists, error) {
	names, err := p.SectionsGet(parser.UserList)
	if err != nil {
		return nil, err
	}
	userlists := []*models.Userlist{}
	for _, name := range names {
		f, err := parseUserlistsSection(name, p)
		if err != nil {
			return nil, err
		}
		userlists = append(userlists, f)
	}
	return userlists, nil
}

func parseUserlistsSection(name string, p parser.Parser) (*models.Userlist, error) {
	u := &models.Userlist{
		UserlistBase: models.UserlistBase{Name: name},
	}
	if err := ParseUserlistSection(p, u); err != nil {
		return nil, err
	}

	userlist, err := ParseUsers(name, p)
	if err != nil {
		return nil, err
	}
	userlista, errula := namedResourceArrayToMapWithKey(userlist, "Username")
	if errula != nil {
		return nil, errula
	}
	u.Users = userlista

	// groups
	groups, err := ParseGroups(name, p)
	if err != nil {
		return nil, err
	}
	groupsa, errga := namedResourceArrayToMap(groups)
	if errga != nil {
		return nil, errga
	}
	u.Groups = groupsa

	return u, nil
}

func serializeUserlistSection(a StructuredToParserArgs, u *models.Userlist) error {
	p := *a.Parser
	var err error

	err = p.SectionsCreate(parser.UserList, u.Name)
	if err != nil {
		return err
	}
	if err = SerializeUserlistSection(p, u, a.Options); err != nil {
		return err
	}

	for _, user := range u.Users {
		if err = p.Insert(parser.UserList, u.Name, "user", SerializeUser(user), -1); err != nil {
			return a.HandleError(user.Username, "userlist", u.Name, a.TID, a.TID == "", err)
		}
	}
	for _, group := range u.Groups {
		if err = p.Insert(parser.UserList, u.Name, "group", SerializeGroup(group), -1); err != nil {
			return a.HandleError(group.Name, "userlist", u.Name, a.TID, a.TID == "", err)
		}
	}

	return nil
}
