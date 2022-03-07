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
	parser "github.com/haproxytech/config-parser/v4"
	parser_errors "github.com/haproxytech/config-parser/v4/errors"
	"github.com/haproxytech/config-parser/v4/types"

	"github.com/haproxytech/client-native/v3/models"
)

type User interface {
	GetUsers(userlist string, transactionID string) (int64, models.Users, error)
	GetUser(username string, userlist string, transactionID string) (int64, *models.User, error)
	DeleteUser(username string, userlist string, transactionID string, version int64) error
	CreateUser(userlist string, data *models.User, transactionID string, version int64) error
	EditUser(username string, userlist string, data *models.User, transactionID string, version int64) error
}

// GetUsers returns configuration version and an array of configured User lines in the specified userlist.
// Returns error on fail.
func (c *client) GetUsers(userlist string, transactionID string) (int64, models.Users, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}
	users, err := ParseUsers(userlist, p)
	if err != nil {
		return v, nil, c.HandleError("", "userlist", userlist, "", false, err)
	}
	return v, users, nil
}

// GetUser returns configuration version and a requested User line in the specified userlist.
// Returns error on fail or if User does not exist
func (c *client) GetUser(username string, userlist string, transactionID string) (int64, *models.User, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}
	user, _ := GetUserByUsername(username, userlist, p)
	if err != nil {
		return v, nil, c.HandleError(username, "userlist", userlist, "", false, err)
	}
	return v, user, nil
}

// DeleteUser deletes a User line in the configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success
func (c *client) DeleteUser(username string, userlist string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}
	user, i := GetUserByUsername(username, userlist, p)
	if user == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("User %s does not exist in userlist %s", username, userlist))
		return c.HandleError(username, "userlist", userlist, "", false, e)
	}
	if err := p.Delete("userlist", userlist, "user", i); err != nil {
		return c.HandleError(username, "userlist", userlist, t, transactionID == "", err)
	}
	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// CreateUser creates a User line in the configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success
func (c *client) CreateUser(userlist string, data *models.User, transactionID string, version int64) error {
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
	user, _ := GetUserByUsername(data.Username, userlist, p)
	if user == nil {
		return c.HandleError(data.Username, "userlist", userlist, "", false, err)
	}
	if err := p.Insert("userlist", userlist, "user", serializeUser(*data), -1); err != nil {
		return c.HandleError(data.Username, "userlist", userlist, t, transactionID == "", err)
	}
	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

// EditUser edits a User line in the configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) EditUser(username string, userlist string, data *models.User, transactionID string, version int64) error {
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
	user, i := GetUserByUsername(data.Username, userlist, p)
	if user == nil {
		return c.HandleError(data.Username, "userlist", userlist, "", false, err)
	}
	if err := p.Set(parser.UserList, userlist, "user", serializeUser(*data), i); err != nil {
		return c.HandleError(data.Username, "userlist", userlist, t, transactionID == "", err)
	}
	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}
	return nil
}

func ParseUsers(userlist string, p parser.Parser) (models.Users, error) {
	users := models.Users{}
	data, err := p.Get("userlist", userlist, "user", false)
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return users, nil
		}
		return nil, err
	}
	items := data.([]types.User)
	for _, item := range items {
		user := ParseUser(item)
		users = append(users, user)
	}
	return users, nil
}

func ParseUser(u types.User) *models.User {
	securePassword := !u.IsInsecure
	return &models.User{
		Username:       u.Name,
		Password:       u.Password,
		SecurePassword: &securePassword,
		Groups:         strings.Join(u.Groups, ","),
	}
}

func serializeUser(u models.User) types.User {
	return types.User{
		Name:       u.Username,
		Password:   u.Password,
		IsInsecure: *u.SecurePassword,
		Groups:     strings.Split(u.Groups, ","),
	}
}

func GetUserByUsername(username string, userlist string, p parser.Parser) (*models.User, int) {
	users, err := ParseUsers(userlist, p)
	if err != nil {
		return nil, 0
	}
	for i, user := range users {
		if user.Username == username {
			return user, i
		}
	}
	return nil, 0
}
