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
	"strconv"

	"github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/config-parser/v3"
	"github.com/haproxytech/config-parser/v3/common"
	"github.com/haproxytech/config-parser/v3/types"

	"github.com/haproxytech/client-native/v2/misc"
	"github.com/haproxytech/client-native/v2/models"
)

// GetResolvers returns configuration version and an array of
// configured resolvers. Returns error on fail.
func (c *Client) GetResolvers(transactionID string) (int64, models.Resolvers, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	fNames, err := p.SectionsGet(parser.Resolvers)
	if err != nil {
		return v, nil, err
	}

	var resolver *models.Resolver
	resolvers := []*models.Resolver{}
	for _, name := range fNames {
		if v, resolver, err = c.GetResolver(name, transactionID); err == nil {
			resolvers = append(resolvers, resolver)
		}
	}

	return v, resolvers, nil
}

// GetResolver returns configuration version and a requested resolver.
// Returns error on fail or if resolver does not exist.
func (c *Client) GetResolver(name string, transactionID string) (int64, *models.Resolver, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !c.checkSectionExists(parser.Resolvers, name, p) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("Resolver %s does not exist", name))
	}

	resolver := &models.Resolver{Name: name}
	if err = ParseResolverSection(p, resolver); err != nil {
		return 0, nil, err
	}

	return v, resolver, nil
}

// DeleteResolver deletes a resolver in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) DeleteResolver(name string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !c.checkSectionExists(parser.Resolvers, name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.Resolvers, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err := p.SectionsDelete(parser.Resolvers, name); err != nil {
		return c.HandleError(name, "", "", t, transactionID == "", err)
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// EditResolver edits a resolver in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) EditResolver(name string, data *models.Resolver, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !c.checkSectionExists(parser.Resolvers, name, p) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exist", parser.Resolvers, name))
		return c.HandleError(name, "", "", t, transactionID == "", e)
	}

	if err = SerializeResolverSection(p, data); err != nil {
		return err
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

// CreateResolver creates a resolver in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *Client) CreateResolver(data *models.Resolver, transactionID string, version int64) error {
	if c.UseValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if c.checkSectionExists(parser.Resolvers, data.Name, p) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.Resolvers, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(parser.Resolvers, data.Name); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if err = SerializeResolverSection(p, data); err != nil {
		return err
	}

	if err := c.SaveData(p, t, transactionID == ""); err != nil {
		return err
	}

	return nil
}

func ParseResolverSection(p *parser.Parser, resolver *models.Resolver) error { //nolint:gocognit,gocyclo
	var err error
	var data common.ParserData
	name := resolver.Name

	if data, err = p.Get(parser.Resolvers, name, "accepted_payload_size", false); err == nil {
		d, ok := data.(*types.StringC)
		if ok && d != nil {
			if n, errInt := strconv.ParseInt(d.Value, 10, 64); errInt == nil {
				resolver.AcceptedPayloadSize = n
			}
		}
	}
	if data, err = p.Get(parser.Resolvers, name, "hold nx", false); err == nil {
		d, ok := data.(*types.StringC)
		if ok && d != nil {
			resolver.HoldNx = misc.ParseTimeout(d.Value)
		}
	}
	if data, err = p.Get(parser.Resolvers, name, "hold obsolete", false); err == nil {
		d, ok := data.(*types.StringC)
		if ok && d != nil {
			resolver.HoldObsolete = misc.ParseTimeout(d.Value)
		}
	}
	if data, err = p.Get(parser.Resolvers, name, "hold other", false); err == nil {
		d, ok := data.(*types.StringC)
		if ok && d != nil {
			resolver.HoldOther = misc.ParseTimeout(d.Value)
		}
	}
	if data, err = p.Get(parser.Resolvers, name, "hold refused", false); err == nil {
		d, ok := data.(*types.StringC)
		if ok && d != nil {
			resolver.HoldRefused = misc.ParseTimeout(d.Value)
		}
	}
	if data, err = p.Get(parser.Resolvers, name, "hold timeout", false); err == nil {
		d, ok := data.(*types.StringC)
		if ok && d != nil {
			resolver.HoldTimeout = misc.ParseTimeout(d.Value)
		}
	}
	if data, err = p.Get(parser.Resolvers, name, "hold valid", false); err == nil {
		d, ok := data.(*types.StringC)
		if ok && d != nil {
			resolver.HoldValid = misc.ParseTimeout(d.Value)
		}
	}
	if data, err = p.Get(parser.Resolvers, name, "resolve_retries", false); err == nil {
		d, ok := data.(*types.StringC)
		if ok && d != nil {
			if n, errInt := strconv.ParseInt(d.Value, 10, 64); errInt == nil {
				resolver.ResolveRetries = n
			}
		}
	}
	if data, err = p.Get(parser.Resolvers, name, "parse-resolv-conf", false); err == nil {
		d, ok := data.(*types.StringC)
		if ok && d != nil {
			resolver.ParseResolvConf = true
		}
	}
	if data, err = p.Get(parser.Resolvers, name, "timeout resolve", false); err == nil {
		d, ok := data.(*types.SimpleTimeout)
		if ok && d != nil {
			tOut := misc.ParseTimeout(d.Value)
			if tOut != nil {
				resolver.TimeoutResolve = *tOut
			}
		}
	}
	if data, err = p.Get(parser.Resolvers, name, "timeout retry", false); err == nil {
		d, ok := data.(*types.SimpleTimeout)
		if ok && d != nil {
			tOut := misc.ParseTimeout(d.Value)
			if tOut != nil {
				resolver.TimeoutRetry = *tOut
			}
		}
	}

	return err
}

func SerializeResolverSection(p *parser.Parser, data *models.Resolver) error { //nolint:gocognit,gocyclo
	var err error

	if data.AcceptedPayloadSize == 0 {
		if err = p.Set(parser.Resolvers, data.Name, "accepted_payload_size", nil); err != nil {
			return err
		}
	} else {
		n := types.StringC{Value: strconv.FormatInt(data.AcceptedPayloadSize, 10)}
		if err = p.Set(parser.Resolvers, data.Name, "accepted_payload_size", n); err != nil {
			return err
		}
	}
	if data.HoldNx == nil {
		if err = p.Set(parser.Resolvers, data.Name, "hold nx", nil); err != nil {
			return err
		}
	} else {
		n := types.StringC{Value: strconv.FormatInt(*data.HoldNx, 10)}
		if err = p.Set(parser.Resolvers, data.Name, "hold nx", n); err != nil {
			return err
		}
	}
	if data.HoldObsolete == nil {
		if err = p.Set(parser.Resolvers, data.Name, "hold obsolete", nil); err != nil {
			return err
		}
	} else {
		n := types.StringC{Value: strconv.FormatInt(*data.HoldObsolete, 10)}
		if err = p.Set(parser.Resolvers, data.Name, "hold obsolete", n); err != nil {
			return err
		}
	}
	if data.HoldOther == nil {
		if err = p.Set(parser.Resolvers, data.Name, "hold other", nil); err != nil {
			return err
		}
	} else {
		n := types.StringC{Value: strconv.FormatInt(*data.HoldOther, 10)}
		if err = p.Set(parser.Resolvers, data.Name, "hold other", n); err != nil {
			return err
		}
	}
	if data.HoldRefused == nil {
		if err = p.Set(parser.Resolvers, data.Name, "hold refused", nil); err != nil {
			return err
		}
	} else {
		n := types.StringC{Value: strconv.FormatInt(*data.HoldRefused, 10)}
		if err = p.Set(parser.Resolvers, data.Name, "hold refused", n); err != nil {
			return err
		}
	}
	if data.HoldTimeout == nil {
		if err = p.Set(parser.Resolvers, data.Name, "hold timeout", nil); err != nil {
			return err
		}
	} else {
		n := types.StringC{Value: strconv.FormatInt(*data.HoldTimeout, 10)}
		if err = p.Set(parser.Resolvers, data.Name, "hold timeout", n); err != nil {
			return err
		}
	}
	if data.HoldValid == nil {
		if err = p.Set(parser.Resolvers, data.Name, "hold valid", nil); err != nil {
			return err
		}
	} else {
		n := types.StringC{Value: strconv.FormatInt(*data.HoldValid, 10)}
		if err = p.Set(parser.Resolvers, data.Name, "hold valid", n); err != nil {
			return err
		}
	}
	if data.ParseResolvConf {
		b := types.StringC{Value: strconv.FormatBool(data.ParseResolvConf)}
		if err = p.Set(parser.Resolvers, data.Name, "parse-resolv-conf", b); err != nil {
			return err
		}
	} else if err = p.Set(parser.Resolvers, data.Name, "parse-resolv-conf", nil); err != nil {
		return err
	}
	if data.ResolveRetries == 0 {
		if err = p.Set(parser.Resolvers, data.Name, "resolve_retries", nil); err != nil {
			return err
		}
	} else {
		n := types.StringC{Value: strconv.FormatInt(data.ResolveRetries, 10)}
		if err = p.Set(parser.Resolvers, data.Name, "resolve_retries", n); err != nil {
			return err
		}
	}
	if data.TimeoutResolve == 0 {
		if err = p.Set(parser.Resolvers, data.Name, "timeout resolve", nil); err != nil {
			return err
		}
	} else {
		timeout := types.SimpleTimeout{Value: strconv.FormatInt(data.TimeoutResolve, 10)}
		if err = p.Set(parser.Resolvers, data.Name, "timeout resolve", timeout); err != nil {
			return err
		}
	}
	if data.TimeoutRetry == 0 {
		if err = p.Set(parser.Resolvers, data.Name, "timeout retry", nil); err != nil {
			return err
		}
	} else {
		timeout := types.SimpleTimeout{Value: strconv.FormatInt(data.TimeoutRetry, 10)}
		if err = p.Set(parser.Resolvers, data.Name, "timeout retry", timeout); err != nil {
			return err
		}
	}

	return err
}
