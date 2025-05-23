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

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
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
	user, _, err := GetUserByUsername(username, userlist, p)
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
	if _, _, err = c.GetUserList(userlist, transactionID); err != nil {
		return err
	}
	_, i, err := GetUserByUsername(username, userlist, p)
	if err != nil {
		return err
	}
	if err := p.Delete("userlist", userlist, "user", i); err != nil {
		return c.HandleError(username, "userlist", userlist, t, transactionID == "", err)
	}
	return c.SaveData(p, t, transactionID == "")
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

	if _, _, err = c.GetUserList(userlist, transactionID); err != nil {
		return err
	}
	_, _, err = GetUserByUsername(data.Username, userlist, p)
	if err == nil {
		return NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("user %s already exists", data.Username))
	}
	if err := p.Insert("userlist", userlist, "user", SerializeUser(*data), -1); err != nil {
		return c.HandleError(data.Username, "userlist", userlist, t, transactionID == "", err)
	}
	return c.SaveData(p, t, transactionID == "")
}

// EditUser edits a User line in the configuration. One of version or transactionID is mandatory.
// Returns error on fail, nil on success.
func (c *client) EditUser(username string, userlist string, data *models.User, transactionID string, version int64) error { //nolint:revive
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
	if _, _, err = c.GetUserList(userlist, transactionID); err != nil {
		return err
	}
	_, i, err := GetUserByUsername(data.Username, userlist, p)
	if err != nil {
		return err
	}
	if err := p.Set(parser.UserList, userlist, "user", SerializeUser(*data), i); err != nil {
		return c.HandleError(data.Username, "userlist", userlist, t, transactionID == "", err)
	}
	return c.SaveData(p, t, transactionID == "")
}

func ParseUsers(userlist string, p parser.Parser) (models.Users, error) {
	var users models.Users
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
		Metadata:       parseMetadata(u.Comment),
	}
}

func SerializeUser(u models.User) types.User {
	if u.SecurePassword == nil {
		u.SecurePassword = misc.BoolP(false)
	}

	var groups []string
	if u.Groups != "" {
		groups = strings.Split(u.Groups, ",")
	}
	comment, _ := serializeMetadata(u.Metadata)
	return types.User{
		Name:       u.Username,
		Password:   u.Password,
		IsInsecure: !*u.SecurePassword,
		Groups:     groups,
		Comment:    comment,
	}
}

func GetUserByUsername(username string, userlist string, p parser.Parser) (*models.User, int, error) {
	users, err := ParseUsers(userlist, p)
	if err != nil {
		return nil, 0, err
	}
	for i, user := range users {
		if user.Username == username {
			return user, i, nil
		}
	}
	return nil, 0, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("user %s does not exist", username))
}
