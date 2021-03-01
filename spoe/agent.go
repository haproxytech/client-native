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

package spoe

import (
	"fmt"
	"strconv"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v3"
	"github.com/haproxytech/config-parser/v3/spoe"
	"github.com/haproxytech/config-parser/v3/types"

	conf "github.com/haproxytech/client-native/v2/configuration"
	"github.com/haproxytech/client-native/v2/misc"
	"github.com/haproxytech/client-native/v2/models"
)

// GetAgents returns configuration version and an array of
// configured agents. Returns error on fail.
func (c *SingleSpoe) GetAgents(scope, transactionID string) (int64, models.SpoeAgents, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	aNames, err := p.SectionsGet(scope, parser.SPOEAgent)
	if err != nil {
		return v, nil, err
	}

	agents := models.SpoeAgents{}
	for _, name := range aNames {
		_, a, err := c.GetAgent(scope, name, transactionID)
		if err == nil {
			agents = append(agents, a)
		}
	}

	return v, agents, nil
}

// GetAgent returns configuration version and a requested agent.
// Returns error on fail or if agent does not exist.
func (c *SingleSpoe) GetAgent(scope, name, transactionID string) (int64, *models.SpoeAgent, error) { //nolint:gocognit,gocyclo
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}
	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(scope, parser.SPOEAgent, name, p) {
		return v, nil, conf.NewConfError(conf.ErrObjectDoesNotExist, fmt.Sprintf("agent %s does not exist", name))
	}

	agent := &models.SpoeAgent{Name: &name}

	// timeouts
	data, err := p.Get(scope, parser.SPOEAgent, name, "timeout hello", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.StringC); ok {
		if d.Value != "" {
			agent.HelloTimeout = *misc.ParseTimeout(d.Value)
		}
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "timeout idle", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.StringC); ok {
		if d.Value != "" {
			agent.IdleTimeout = *misc.ParseTimeout(d.Value)
		}
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "timeout processing", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.StringC); ok {
		if d.Value != "" {
			agent.ProcessingTimeout = *misc.ParseTimeout(d.Value)
		}
	}

	// options
	data, err = p.Get(scope, parser.SPOEAgent, name, "option continue-on-error", true)
	if err != nil {
		return v, nil, err
	}
	if !data.(*types.SimpleOption).NoOption {
		agent.ContinueOnError = "enabled"
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "option force-set-var", true)
	if err != nil {
		return v, nil, err
	}
	if !data.(*types.SimpleOption).NoOption {
		agent.ForceSetVar = "enabled"
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "option set-on-error", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.StringC); ok {
		agent.OptionSetOnError = d.Value
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "option set-process-time", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.StringC); ok {
		agent.OptionSetProcessTime = d.Value
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "option set-total-time", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.StringC); ok {
		agent.OptionSetTotalTime = d.Value
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "option var-prefix", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.StringC); ok {
		agent.OptionVarPrefix = d.Value
	}

	// others
	data, err = p.Get(scope, parser.SPOEAgent, name, "register-var-names", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.StringC); ok {
		agent.RegisterVarNames = d.Value
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "use-backend", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.StringC); ok {
		agent.UseBackend = d.Value
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "maxconnrate", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.Int64C); ok {
		agent.Maxconnrate = d.Value
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "maxerrrate", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.Int64C); ok {
		agent.Maxerrrate = d.Value
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "max-frame-size", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.Int64C); ok {
		agent.MaxFrameSize = d.Value
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "max-waiting-frames", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.Int64C); ok {
		agent.MaxWaitingFrames = d.Value
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "messages", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.StringC); ok {
		agent.Messages = d.Value
	}

	// simple options
	data, err = p.Get(scope, parser.SPOEAgent, name, "option async", true)
	if err != nil {
		return v, nil, err
	}
	if data.(*types.SimpleOption).NoOption {
		agent.Async = "disabled"
	} else {
		agent.Async = "enabled"
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "option dontlog-normal", true)
	if err != nil {
		return v, nil, err
	}
	if data.(*types.SimpleOption).NoOption {
		agent.DontlogNormal = "disabled"
	} else {
		agent.DontlogNormal = "enabled"
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "option pipelining", true)
	if err != nil {
		return v, nil, err
	}
	if data.(*types.SimpleOption).NoOption {
		agent.Pipelining = "disabled"
	} else {
		agent.Pipelining = "enabled"
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "option send-frag-payload", true)
	if err != nil {
		return v, nil, err
	}
	if data.(*types.SimpleOption).NoOption {
		agent.SendFragPayload = "disabled"
	} else {
		agent.SendFragPayload = "enabled"
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "groups", true)
	if err != nil {
		return v, nil, err
	}
	if d, ok := data.(*types.StringC); ok {
		agent.Groups = d.Value
	}

	data, err = p.Get(scope, parser.SPOEAgent, name, "log", true)
	if err != nil {
		return v, nil, err
	}
	if logs, ok := data.([]types.Log); ok {
		for i, l := range logs {
			indx := int64(i)
			d := &models.LogTarget{
				Address:  l.Address,
				Facility: l.Facility,
				Format:   l.Format,
				Global:   l.Global,
				Length:   l.Length,
				Level:    l.Level,
				Minlevel: l.MinLevel,
				Nolog:    l.NoLog,
				Index:    &indx,
			}
			agent.Log = append(agent.Log, d)
		}
	}
	return v, agent, nil
}

// DeleteAgent deletes an agent in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *SingleSpoe) DeleteAgent(scope, name, transactionID string, version int64) error {
	if err := c.deleteSection(scope, parser.SPOEAgent, name, transactionID, version); err != nil {
		return err
	}
	return nil
}

// CreateAgent creates a agent in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *SingleSpoe) CreateAgent(scope string, data *models.SpoeAgent, transactionID string, version int64) error {
	if c.Transaction.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return conf.NewConfError(conf.ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if c.checkSectionExists(scope, parser.SPOEAgent, *data.Name, p) {
		e := conf.NewConfError(conf.ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.SPOEAgent, *data.Name))
		return c.Transaction.HandleError(*data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(scope, parser.SPOEAgent, *data.Name); err != nil {
		return err
	}

	err = c.createEditAgent(scope, data, t, transactionID, p)
	if err != nil {
		return err
	}

	if err := c.Transaction.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// EditAgent edits a agent in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *SingleSpoe) EditAgent(scope string, data *models.SpoeAgent, transactionID string, version int64) error {
	if c.Transaction.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return conf.NewConfError(conf.ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !c.checkSectionExists(scope, parser.SPOEAgent, *data.Name, p) {
		e := conf.NewConfError(conf.ErrObjectAlreadyExists, fmt.Sprintf("%s %s does not exists", parser.SPOEAgent, *data.Name))
		return c.Transaction.HandleError(*data.Name, "", "", t, transactionID == "", e)
	}

	err = c.createEditAgent(scope, data, t, transactionID, p)
	if err != nil {
		return err
	}

	if err := c.Transaction.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func (c *SingleSpoe) createEditAgent(scope string, data *models.SpoeAgent, t string, transactionID string, p *spoe.Parser) error { //nolint:gocognit,gocyclo
	if data == nil {
		return fmt.Errorf("spoe agent not initialized")
	}
	name := *data.Name

	if data.DontlogNormal == "enabled" {
		d := &types.SimpleOption{}
		if err := p.Set(scope, parser.SPOEAgent, name, "option dontlog-normal", d); err != nil {
			return c.Transaction.HandleError("option dontlog-normal", "", "", t, transactionID == "", err)
		}
	} else if data.DontlogNormal == "disabled" {
		d := &types.SimpleOption{NoOption: true}
		if err := p.Set(scope, parser.SPOEAgent, name, "option dontlog-normal", d); err != nil {
			return c.Transaction.HandleError("option dontlog-normal", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "option dontlog-normal", nil); err != nil {
		return err
	}

	if data.Async == "enabled" {
		d := &types.SimpleOption{}
		if err := p.Set(scope, parser.SPOEAgent, name, "option async", d); err != nil {
			return c.Transaction.HandleError("option async", "", "", t, transactionID == "", err)
		}
	} else if data.Async == "disabled" {
		d := &types.SimpleOption{NoOption: true}
		if err := p.Set(scope, parser.SPOEAgent, name, "option async", d); err != nil {
			return c.Transaction.HandleError("option async", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "option async", nil); err != nil {
		return err
	}

	if data.ContinueOnError == "enabled" {
		d := &types.SimpleOption{}
		if err := p.Set(scope, parser.SPOEAgent, name, "option continue-on-error", d); err != nil {
			return c.Transaction.HandleError("option continue-on-error", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "option continue-on-error", nil); err != nil {
		return err
	}

	if data.Pipelining == "enabled" {
		d := &types.SimpleOption{}
		if err := p.Set(scope, parser.SPOEAgent, name, "option pipelining", d); err != nil {
			return c.Transaction.HandleError("option pipelining", "", "", t, transactionID == "", err)
		}
	} else if data.Pipelining == "disabled" {
		d := &types.SimpleOption{NoOption: true}
		if err := p.Set(scope, parser.SPOEAgent, name, "option pipelining", d); err != nil {
			return c.Transaction.HandleError("option pipelining", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "option pipelining", nil); err != nil {
		return err
	}

	if data.SendFragPayload == "enabled" {
		d := &types.SimpleOption{}
		if err := p.Set(scope, parser.SPOEAgent, name, "option send-frag-payload", d); err != nil {
			return c.Transaction.HandleError("option send-frag-payload", "", "", t, transactionID == "", err)
		}
	} else if data.SendFragPayload == "disabled" {
		d := &types.SimpleOption{NoOption: true}
		if err := p.Set(scope, parser.SPOEAgent, name, "option send-frag-payload", d); err != nil {
			return c.Transaction.HandleError("option send-frag-payload", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "option send-frag-payload", nil); err != nil {
		return err
	}

	if data.Maxconnrate > 0 {
		d := &types.Int64C{Value: data.Maxconnrate}
		if err := p.Set(scope, parser.SPOEAgent, name, "maxconnrate", d); err != nil {
			return c.Transaction.HandleError("maxconnrate", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "maxconnrate", nil); err != nil {
		return err
	}

	if data.Maxerrrate > 0 {
		d := &types.Int64C{Value: data.Maxerrrate}
		if err := p.Set(scope, parser.SPOEAgent, name, "maxerrrate", d); err != nil {
			return c.Transaction.HandleError("maxerrrate", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "maxerrrate", nil); err != nil {
		return err
	}

	if data.MaxFrameSize > 0 {
		d := &types.Int64C{Value: data.MaxFrameSize}
		if err := p.Set(scope, parser.SPOEAgent, name, "max-frame-size", d); err != nil {
			return c.Transaction.HandleError("max-frame-size", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "max-frame-size", nil); err != nil {
		return err
	}

	if data.MaxWaitingFrames > 0 {
		d := &types.Int64C{Value: data.MaxWaitingFrames}
		if err := p.Set(scope, parser.SPOEAgent, name, "max-waiting-frames", d); err != nil {
			return c.Transaction.HandleError("max-waiting-frames", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "max-waiting-frames", nil); err != nil {
		return err
	}

	if data.Messages != "" {
		d := &types.StringC{Value: data.Messages}
		if err := p.Set(scope, parser.SPOEAgent, name, "messages", d); err != nil {
			return c.Transaction.HandleError("messages", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "messages", nil); err != nil {
		return err
	}

	if data.ForceSetVar == "enabled" {
		d := &types.SimpleOption{}
		if err := p.Set(scope, parser.SPOEAgent, name, "option force-set-var", d); err != nil {
			return c.Transaction.HandleError("option force-set-var", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "option force-set-var", nil); err != nil {
		return err
	}

	if data.OptionSetOnError != "" {
		d := &types.StringC{Value: data.OptionSetOnError}
		if err := p.Set(scope, parser.SPOEAgent, name, "option set-on-error", d); err != nil {
			return c.Transaction.HandleError("option set-on-error", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "option set-on-error", nil); err != nil {
		return err
	}

	if data.OptionSetProcessTime != "" {
		d := &types.StringC{Value: data.OptionSetProcessTime}
		if err := p.Set(scope, parser.SPOEAgent, name, "option set-process-time", d); err != nil {
			return c.Transaction.HandleError("option set-process-time", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "option set-process-time", nil); err != nil {
		return err
	}

	if data.OptionSetTotalTime != "" {
		d := &types.StringC{Value: data.OptionSetTotalTime}
		if err := p.Set(scope, parser.SPOEAgent, name, "option set-total-time", d); err != nil {
			return c.Transaction.HandleError("option set-total-time", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "option set-total-time", nil); err != nil {
		return err
	}

	if data.OptionVarPrefix != "" {
		d := &types.StringC{Value: data.OptionVarPrefix}
		if err := p.Set(scope, parser.SPOEAgent, name, "option var-prefix", d); err != nil {
			return c.Transaction.HandleError("option var-prefix", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "option var-prefix", nil); err != nil {
		return err
	}

	if data.RegisterVarNames != "" {
		d := &types.StringC{Value: data.RegisterVarNames}
		if err := p.Set(scope, parser.SPOEAgent, name, "register-var-names", d); err != nil {
			return c.Transaction.HandleError("register-var-names", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "register-var-names", nil); err != nil {
		return err
	}

	if data.HelloTimeout > 0 {
		d := &types.StringC{Value: strconv.FormatInt(data.HelloTimeout, 10)}
		if err := p.Set(scope, parser.SPOEAgent, name, "timeout hello", d); err != nil {
			return c.Transaction.HandleError(d.Value, "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "timeout hello", nil); err != nil {
		return err
	}

	if data.IdleTimeout > 0 {
		d := &types.StringC{Value: strconv.FormatInt(data.IdleTimeout, 10)}
		if err := p.Set(scope, parser.SPOEAgent, name, "timeout idle", d); err != nil {
			return c.Transaction.HandleError(d.Value, "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "timeout idle", nil); err != nil {
		return err
	}

	if data.ProcessingTimeout > 0 {
		d := &types.StringC{Value: strconv.FormatInt(data.ProcessingTimeout, 10)}
		if err := p.Set(scope, parser.SPOEAgent, name, "timeout processing", d); err != nil {
			return c.Transaction.HandleError(d.Value, "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "timeout processing", nil); err != nil {
		return err
	}

	if data.UseBackend != "" {
		d := &types.StringC{Value: data.UseBackend}
		if err := p.Set(scope, parser.SPOEAgent, name, "use-backend", d); err != nil {
			return c.Transaction.HandleError("use-backend", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "use-backend", nil); err != nil {
		return err
	}

	if data.Groups != "" {
		d := &types.StringC{Value: data.Groups}
		if err := p.Set(scope, parser.SPOEAgent, name, "groups", d); err != nil {
			return c.Transaction.HandleError("groups", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "groups", nil); err != nil {
		return err
	}

	if len(data.Log) > 0 {
		logs := []types.Log{}
		for _, d := range data.Log {
			log := types.Log{
				Global:   d.Global,
				Address:  d.Address,
				Facility: d.Facility,
				Format:   d.Format,
				Length:   d.Length,
				Level:    d.Level,
				MinLevel: d.Minlevel,
				NoLog:    d.Nolog,
			}
			logs = append(logs, log)
		}
		if err := p.Set(scope, parser.SPOEAgent, name, "log", logs); err != nil {
			return c.Transaction.HandleError("log", "", "", t, transactionID == "", err)
		}
	} else if err := p.Set(scope, parser.SPOEAgent, name, "log", nil); err != nil {
		return err
	}
	return nil
}
