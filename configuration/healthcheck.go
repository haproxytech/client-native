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

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parsererrors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/misc"

	"github.com/haproxytech/client-native/v6/models"
)

type Healthcheck interface {
	GetHealthchecks(transactionID string) (int64, models.Healthchecks, error)
	GetHealthcheck(name string, transactionID string) (int64, *models.HealthCheck, error)
	DeleteHealthcheck(name string, transactionID string, version int64) error
	CreateHealthcheck(data *models.HealthCheck, transactionID string, version int64) error
	EditHealthcheck(name string, data *models.HealthCheck, transactionID string, version int64) error
}

// GetHealthchecks returns configuration version and an array of
// configured health checks. Returns error on fail.
func (c *client) GetHealthchecks(transactionID string) (int64, models.Healthchecks, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	hNames, err := p.SectionsGet(parser.HealthChecks)
	if err != nil {
		return v, nil, err
	}

	healthChecks := []*models.HealthCheck{}
	for _, name := range hNames {
		_, a, err := c.GetHealthcheck(name, transactionID)
		if err == nil {
			healthChecks = append(healthChecks, a)
		}
	}

	return v, healthChecks, nil
}

// GetHealthcheck returns configuration version and a requested health check.
// Returns error on fail or if health check does not exist.
func (c *client) GetHealthcheck(name string, transactionID string) (int64, *models.HealthCheck, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !p.SectionExists(parser.HealthChecks, name) {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("health check %s does not exist", name))
	}

	healthCheck := &models.HealthCheck{HealthCheckBase: models.HealthCheckBase{Name: name}}
	err = ParseHealthcheckSection(p, healthCheck)
	if err != nil {
		return v, nil, err
	}

	return v, healthCheck, nil
}

func ParseHealthcheckSection(p parser.Parser, hc *models.HealthCheck) error { //nolint:gocognit
	// simple types
	sslHelloChk, err := p.Get(parser.HealthChecks, hc.Name, "type ssl-hello-chk", false)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		_, ok := sslHelloChk.(*types.SimpleType)
		if !ok {
			return misc.CreateTypeAssertError("type ssl-hello-chk")
		}
		hc.Type = "ssl-hello-chk"
	}

	redisCheck, err := p.Get(parser.HealthChecks, hc.Name, "type redis-check", false)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		_, ok := redisCheck.(*types.SimpleType)
		if !ok {
			return misc.CreateTypeAssertError("type redis-check")
		}
		hc.Type = "redis-check"
	}

	ldapCheck, err := p.Get(parser.HealthChecks, hc.Name, "type ldap-check", false)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		_, ok := ldapCheck.(*types.SimpleType)
		if !ok {
			return misc.CreateTypeAssertError("type ldap-check")
		}
		hc.Type = "ldap-check"
	}

	spopCheck, err := p.Get(parser.HealthChecks, hc.Name, "type spop-check", false)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		_, ok := spopCheck.(*types.SimpleType)
		if !ok {
			return misc.CreateTypeAssertError("type spop-check")
		}
		hc.Type = "spop-check"
	}

	tcpCheck, err := p.Get(parser.HealthChecks, hc.Name, "type tcp-check", false)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		_, ok := tcpCheck.(*types.SimpleType)
		if !ok {
			return misc.CreateTypeAssertError("type tcp-check")
		}
		hc.Type = "tcp-check"
	}

	// complex types
	pgsqlchk, err := p.Get(parser.HealthChecks, hc.Name, "type pgsql-check", false)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		pgsqlchkData, ok := pgsqlchk.(*types.TypePgsqlCheck)
		if !ok {
			return misc.CreateTypeAssertError("type pgsql-check")
		}

		hc.PgsqlCheckParams = &models.PgsqlCheckParams{
			Username: pgsqlchkData.User,
		}
		hc.Type = "pgsql-check"
	}

	smtpchk, err := p.Get(parser.HealthChecks, hc.Name, "type smtpchk", false)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		smtpchkData, ok := smtpchk.(*types.TypeSmtpchk)
		if !ok {
			return misc.CreateTypeAssertError("type smtpchk")
		}

		hc.SmtpchkParams = &models.SmtpchkParams{
			Domain: smtpchkData.Domain,
			Hello:  smtpchkData.Hello,
		}
		hc.Type = "smtpchk"
	}

	mysqlchk, err := p.Get(parser.HealthChecks, hc.Name, "type mysql-check", false)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		mysqlchkData, ok := mysqlchk.(*types.TypeMysqlCheck)
		if !ok {
			return misc.CreateTypeAssertError("type mysql-check")
		}

		hc.MysqlCheckParams = &models.MysqlCheckParams{
			ClientVersion: mysqlchkData.ClientVersion,
			Username:      mysqlchkData.User,
		}
		hc.Type = "mysql-check"
	}

	httpchk, err := p.Get(parser.HealthChecks, hc.Name, "type httpchk", false)
	if err != nil && !errors.Is(err, parsererrors.ErrFetch) {
		return err
	}
	if err == nil {
		httpchkData, ok := httpchk.(*types.TypeHttpchk)
		if !ok {
			return misc.CreateTypeAssertError("type httpchk")
		}

		hc.HttpchkParams = &models.HttpchkParams{
			Host:    httpchkData.Host,
			Method:  httpchkData.Method,
			URI:     httpchkData.URI,
			Version: httpchkData.Version,
		}
		hc.Type = "httpchk"
	}

	return nil
}

// DeleteHealthcheck deletes a health check in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) DeleteHealthcheck(name string, transactionID string, version int64) error {
	return c.deleteSection(parser.HealthChecks, name, transactionID, version)
}

// CreateHealthcheck creates a health check in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) CreateHealthcheck(data *models.HealthCheck, transactionID string, version int64) error {
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

	if p.SectionExists(parser.HealthChecks, data.Name) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.HealthChecks, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(parser.HealthChecks, data.Name); err != nil {
		return err
	}

	if err = SerializeHealthCheckSection(p, data); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

// EditHealthcheck edits a health check in configuration. One of version or transactionID is
// mandatory. Returns error on fail, nil on success.
func (c *client) EditHealthcheck(name string, data *models.HealthCheck, transactionID string, version int64) error { //nolint:revive
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

	if !p.SectionExists(parser.HealthChecks, data.Name) {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("%s %s does not exists", parser.HealthChecks, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = SerializeHealthCheckSection(p, data); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

func clearAllTypes(p parser.Parser, data *models.HealthCheck) error {
	var err error
	if err = p.Set(parser.HealthChecks, data.Name, "type smtpchk", nil); err != nil {
		return err
	}

	if err = p.Set(parser.HealthChecks, data.Name, "type mysql-check", nil); err != nil {
		return err
	}

	if err = p.Set(parser.HealthChecks, data.Name, "type pgsql-check", nil); err != nil {
		return err
	}

	if err = p.Set(parser.HealthChecks, data.Name, "type httpchk", nil); err != nil {
		return err
	}

	if err = p.Set(parser.HealthChecks, data.Name, "type ssl-hello-chk", nil); err != nil {
		return err
	}

	if err = p.Set(parser.HealthChecks, data.Name, "type redis-check", nil); err != nil {
		return err
	}

	if err = p.Set(parser.HealthChecks, data.Name, "type ldap-check", nil); err != nil {
		return err
	}

	if err = p.Set(parser.HealthChecks, data.Name, "type spop-check", nil); err != nil {
		return err
	}

	if err = p.Set(parser.HealthChecks, data.Name, "type tcp-check", nil); err != nil {
		return err
	}

	return nil
}

func SerializeHealthCheckSection(p parser.Parser, data *models.HealthCheck) error { //nolint:gocognit
	if data == nil {
		return errors.New("empty health check")
	}

	var err error
	err = clearAllTypes(p, data)
	if err != nil {
		return err
	}

	switch data.Type {
	case "smtpchk":
		if data.SmtpchkParams != nil {
			smtpchck := types.TypeSmtpchk{
				Domain: data.SmtpchkParams.Domain,
				Hello:  data.SmtpchkParams.Hello,
				NoType: false,
			}
			if err := p.Set(parser.HealthChecks, data.Name, "type smtpchk", smtpchck); err != nil {
				return err
			}
		} else {
			if err = p.Set(parser.HealthChecks, data.Name, "type smtpchk", nil); err != nil {
				return err
			}
		}
	case "mysql-check":
		if data.MysqlCheckParams != nil {
			mysqlchk := types.TypeMysqlCheck{
				ClientVersion: data.MysqlCheckParams.ClientVersion,
				User:          data.MysqlCheckParams.Username,
				NoType:        false,
			}
			if err := p.Set(parser.HealthChecks, data.Name, "type mysql-check", mysqlchk); err != nil {
				return err
			}
		} else {
			if err = p.Set(parser.HealthChecks, data.Name, "type mysql-check", nil); err != nil {
				return err
			}
		}
	case "pgsql-check":
		if data.PgsqlCheckParams != nil {
			pgsqlchk := types.TypePgsqlCheck{
				User:   data.PgsqlCheckParams.Username,
				NoType: false,
			}
			if err := p.Set(parser.HealthChecks, data.Name, "type pgsql-check", pgsqlchk); err != nil {
				return err
			}
		} else {
			if err = p.Set(parser.HealthChecks, data.Name, "type pgsql-check", nil); err != nil {
				return err
			}
		}
	case "httpchk":
		if data.HttpchkParams != nil {
			httpchk := types.TypeHttpchk{
				Method:  data.HttpchkParams.Method,
				URI:     data.HttpchkParams.URI,
				Version: data.HttpchkParams.Version,
				NoType:  false,
			}
			if err := p.Set(parser.HealthChecks, data.Name, "type httpchk", httpchk); err != nil {
				return err
			}
		} else {
			if err = p.Set(parser.HealthChecks, data.Name, "type httpchk", nil); err != nil {
				return err
			}
		}
	case "ssl-hello-chk":
		if err = p.Set(parser.HealthChecks, data.Name, "type ssl-hello-chk", &types.SimpleType{}); err != nil {
			return err
		}
	case "redis-check":
		if err = p.Set(parser.HealthChecks, data.Name, "type redis-check", &types.SimpleType{}); err != nil {
			return err
		}
	case "ldap-check":
		if err = p.Set(parser.HealthChecks, data.Name, "type ldap-check", &types.SimpleType{}); err != nil {
			return err
		}
	case "spop-check":
		if err = p.Set(parser.HealthChecks, data.Name, "type spop-check", &types.SimpleType{}); err != nil {
			return err
		}
	case "tcp-check":
		if err = p.Set(parser.HealthChecks, data.Name, "type tcp-check", &types.SimpleType{}); err != nil {
			return err
		}
	default:
		if err = p.Set(parser.HealthChecks, data.Name, "type", nil); err != nil {
			return err
		}
	}

	return nil
}
