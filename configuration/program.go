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

package configuration

import (
	"errors"
	"fmt"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v4"
	parsererrors "github.com/haproxytech/config-parser/v4/errors"
	"github.com/haproxytech/config-parser/v4/types"

	"github.com/haproxytech/client-native/v4/misc"
	"github.com/haproxytech/client-native/v4/models"
)

type Program interface {
	GetPrograms(transactionID string) (int64, models.Programs, error)
	GetProgram(name string, transactionID string) (int64, *models.Program, error)
	DeleteProgram(name string, transactionID string, version int64) error
	CreateProgram(data *models.Program, transactionID string, version int64) error
	EditProgram(name string, data *models.Program, transactionID string, version int64) error
}

func (c *client) GetPrograms(transactionID string) (int64, models.Programs, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	programNames, err := p.SectionsGet(parser.Program)
	if err != nil {
		return v, nil, err
	}

	programs := make(models.Programs, 0, len(programNames))

	for _, name := range programNames {
		program, parseErr := ParseProgram(p, name)
		if parseErr != nil {
			continue
		}

		programs = append(programs, program)
	}

	return v, programs, nil
}

func (c *client) GetProgram(name string, transactionID string) (int64, *models.Program, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.Program, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Program %s does not exist", name))
	}

	program, parseErr := ParseProgram(p, name)
	if parseErr != nil {
		return 0, nil, parseErr
	}

	return v, program, nil
}

func (c *client) DeleteProgram(name string, transactionID string, version int64) error {
	if err := c.deleteSection(parser.Program, name, transactionID, version); err != nil {
		return err
	}

	return nil
}

func (c *client) CreateProgram(data *models.Program, transactionID string, version int64) error {
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

	if c.checkSectionExists(parser.Program, data.Name, p) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.Program, data.Name))

		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(parser.Program, data.Name); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if err = SerializeProgramSection(p, data); err != nil {
		return err
	}

	if err = c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func (c *client) EditProgram(name string, data *models.Program, transactionID string, version int64) error {
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

	if !c.checkSectionExists(parser.Program, name, p) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s does not exists", parser.Program, data.Name))

		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = SerializeProgramSection(p, data); err != nil {
		return err
	}

	if err = c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func SerializeProgramSection(p parser.Parser, data *models.Program) error {
	if data == nil {
		return fmt.Errorf("empty program")
	}

	if data.Command == nil {
		return fmt.Errorf("command must be set")
	}
	if err := p.Set(parser.Program, data.Name, "command", types.Command{Args: *data.Command}); err != nil {
		return err
	}

	user := &types.StringC{Value: data.User}
	if data.User == "" {
		user = nil
	}
	if err := p.Set(parser.Program, data.Name, "user", user); err != nil {
		return err
	}

	group := &types.StringC{Value: data.Group}
	if data.Group == "" {
		group = nil
	}
	if err := p.Set(parser.Program, data.Name, "group", group); err != nil {
		return err
	}

	if err := p.Set(parser.Program, data.Name, "option start-on-reload", serializeSimpleOption(data.StartOnReload)); err != nil {
		return err
	}

	return nil
}

func ParseProgram(p parser.Parser, name string) (*models.Program, error) {
	program := models.Program{
		Name: name,
	}

	data, err := p.Get(parser.Program, name, "command")
	if err != nil {
		if errors.Is(err, parsererrors.ErrFetch) {
			data = types.Command{}
		} else {
			return nil, err
		}
	}

	program.Command = misc.StringP((data.(*types.Command)).Args)

	data, err = p.Get(parser.Program, name, "user")
	if err == nil {
		program.User = data.(*types.StringC).Value
	}

	data, err = p.Get(parser.Program, name, "group")
	if err == nil {
		program.Group = data.(*types.StringC).Value
	}

	data, err = p.Get(parser.Program, name, "option start-on-reload")
	if err == nil {
		opt := data.(*types.SimpleOption)
		if opt.NoOption {
			program.StartOnReload = models.ProgramStartOnReloadDisabled
		} else {
			program.StartOnReload = models.ProgramStartOnReloadEnabled
		}
	}

	return &program, nil
}
