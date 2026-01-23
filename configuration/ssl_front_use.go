// Copyright 2025 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/params"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type SSLFrontUse interface {
	GetSSLFrontUses(parentType string, parentName string, transactionID string) (int64, models.SSLFrontUses, error)
	GetSSLFrontUse(number int64, parentType string, parentName string, transactionID string) (int64, *models.SSLFrontUse, error)
	DeleteSSLFrontUse(number int64, parentType string, parentName string, transactionID string, version int64) error
	CreateSSLFrontUse(parentType string, parentName string, data *models.SSLFrontUse, transactionID string, version int64) error
	EditSSLFrontUse(number int64, parentType string, parentName string, data *models.SSLFrontUse, transactionID string, version int64) error
}

func (c *client) GetSSLFrontUses(parentType string, parentName string, transactionID string) (int64, models.SSLFrontUses, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	sfu, err := ParseSSLFrontUses(parentType, parentName, p)
	if err != nil {
		return v, nil, c.HandleError("", parentType, parentName, "", false, err)
	}

	return v, sfu, nil
}

func (c *client) GetSSLFrontUse(number int64, parentType string, parentName string, transactionID string) (int64, *models.SSLFrontUse, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	sfu := GetSSLFrontUseByNumber(number, parentType, parentName, p)
	if sfu == nil {
		return v, nil, NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("SSLFrontUse %d does not exist in %s %s", number, parentName, parentType))
	}

	return v, sfu, nil
}

func (c *client) DeleteSSLFrontUse(number int64, parentType string, parentName string, transactionID string, version int64) error {
	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	sfu := GetSSLFrontUseByNumber(number, parentType, parentName, p)
	if sfu == nil {
		e := NewConfError(ErrObjectDoesNotExist, fmt.Sprintf("SSLFrontUse %d does not exist in %s %s", number, parentName, parentType))
		return c.HandleError("ssl-f-use", parentType, parentName, t, transactionID == "", e)
	}

	if err := p.Delete(parser.Section(parentType), parentName, "ssl-f-use", int(number)); err != nil {
		return c.HandleError("ssl-f-use", parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) CreateSSLFrontUse(parentType string, parentName string, data *models.SSLFrontUse, transactionID string, version int64) error {
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

	if err := p.Insert(parser.Section(parentType), parentName, "ssl-f-use", SerializeSSLFrontUse(*data)); err != nil {
		return c.HandleError("ssl-f-use", parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) EditSSLFrontUse(number int64, parentType string, parentName string, data *models.SSLFrontUse, transactionID string, version int64) error {
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

	if err := p.Set(parser.Section(parentType), parentName, "ssl-f-use", SerializeSSLFrontUse(*data), int(number)); err != nil {
		return c.HandleError("ssl-f-use", parentType, parentName, t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseSSLFrontUses(parentType string, parentName string, p parser.Parser) (models.SSLFrontUses, error) {
	var uses models.SSLFrontUses

	data, err := p.Get(parser.Section(parentType), parentName, "ssl-f-use")
	if err != nil {
		if errors.Is(err, parser_errors.ErrFetch) {
			return uses, nil
		}
		return nil, err
	}

	ondiskUses := data.([]types.SSLFrontUse) //nolint:forcetypeassert
	uses = make(models.SSLFrontUses, 0, len(ondiskUses))
	for _, ondiskUse := range ondiskUses {
		b := ParseSSLFrontUse(ondiskUse)
		if b != nil {
			uses = append(uses, b)
		}
	}
	return uses, nil
}

func ParseSSLFrontUse(ondisk types.SSLFrontUse) *models.SSLFrontUse {
	u := models.SSLFrontUse{}

	for _, p := range ondisk.Params {
		switch v := p.(type) {
		case *params.BindOptionWord:
			switch v.Name {
			case "allow-0rtt":
				u.Allow0rtt = true
			case "no-alpn":
				u.NoAlpn = true
			case "no-ca-names":
				u.NoCaNames = true
			}
		case *params.BindOptionValue:
			switch v.Name {
			case "alpn":
				u.Alpn = v.Value
			case "ca-file":
				u.CaFile = v.Value
			case "ciphers":
				u.Ciphers = v.Value
			case "ciphersuites":
				u.Ciphersuites = v.Value
			case "client-sigalgs":
				u.ClientSigalgs = v.Value
			case "crl-file":
				u.CrlFile = v.Value
			case "crt":
				u.Certificate = v.Value
			case "curves":
				u.Curves = v.Value
			case "ecdhe":
				u.Ecdhe = v.Value
			case "issuer":
				u.Issuer = v.Value
			case "key":
				u.Key = v.Value
			case "npn":
				u.Npn = v.Value
			case "ocsp":
				u.Ocsp = v.Value
			case "sctl":
				u.Sctl = v.Value
			case "sigalgs":
				u.Sigalgs = v.Value
			case "ssl-max-ver":
				u.SslMaxVer = v.Value
			case "ssl-min-ver":
				u.SslMinVer = v.Value
			case "verify":
				u.Verify = v.Value
			}
		case *params.BindOptionOnOff:
			if v.Name == "ocsp-update" {
				switch v.Value {
				case "on":
					u.OcspUpdate = "enabled"
				case "off":
					u.OcspUpdate = "disabled"
				}
			}
		}
	}

	u.Metadata = misc.ParseMetadata(ondisk.Comment)
	return &u
}

func SerializeSSLFrontUse(m models.SSLFrontUse) types.SSLFrontUse {
	u := types.SSLFrontUse{}
	comment, err := misc.SerializeMetadata(m.Metadata)
	if err != nil {
		comment = ""
	}
	u.Comment = comment

	options := make([]params.SSLBindOption, 0, 8)

	if m.Certificate != "" {
		options = append(options, &params.BindOptionValue{Name: "crt", Value: m.Certificate})
	}
	if m.CaFile != "" {
		options = append(options, &params.BindOptionValue{Name: "ca-file", Value: m.CaFile})
	}
	if m.Verify != "" {
		options = append(options, &params.BindOptionValue{Name: "verify", Value: m.Verify})
	}
	if m.Alpn != "" {
		options = append(options, &params.BindOptionValue{Name: "alpn", Value: m.Alpn})
	}
	if m.Allow0rtt {
		options = append(options, &params.BindOptionWord{Name: "allow-0rtt"})
	}
	if m.Curves != "" {
		options = append(options, &params.BindOptionValue{Name: "curves", Value: m.Curves})
	}
	if m.Ecdhe != "" {
		options = append(options, &params.BindOptionValue{Name: "ecdhe", Value: m.Ecdhe})
	}
	if m.Ciphers != "" {
		options = append(options, &params.BindOptionValue{Name: "ciphers", Value: m.Ciphers})
	}
	if m.Ciphersuites != "" {
		options = append(options, &params.BindOptionValue{Name: "ciphersuites", Value: m.Ciphersuites})
	}
	if m.ClientSigalgs != "" {
		options = append(options, &params.BindOptionValue{Name: "client-sigalgs", Value: m.ClientSigalgs})
	}
	if m.CrlFile != "" {
		options = append(options, &params.BindOptionValue{Name: "crl-file", Value: m.CrlFile})
	}
	if m.NoCaNames {
		options = append(options, &params.BindOptionWord{Name: "no-ca-names"})
	}
	if m.Npn != "" {
		options = append(options, &params.BindOptionValue{Name: "npn", Value: m.Npn})
	}
	if m.Sigalgs != "" {
		options = append(options, &params.BindOptionValue{Name: "sigalgs", Value: m.Sigalgs})
	}
	if m.SslMaxVer != "" {
		options = append(options, &params.BindOptionValue{Name: "ssl-max-ver", Value: m.SslMaxVer})
	}
	if m.SslMinVer != "" {
		options = append(options, &params.BindOptionValue{Name: "ssl-min-ver", Value: m.SslMinVer})
	}
	if m.NoAlpn {
		options = append(options, &params.BindOptionWord{Name: "no-alpn"})
	}

	if m.Key != "" {
		options = append(options, &params.BindOptionValue{Name: "key", Value: m.Key})
	}
	if m.Issuer != "" {
		options = append(options, &params.BindOptionValue{Name: "issuer", Value: m.Issuer})
	}
	if m.Ocsp != "" {
		options = append(options, &params.BindOptionValue{Name: "ocsp", Value: m.Ocsp})
	}
	if m.Sctl != "" {
		options = append(options, &params.BindOptionValue{Name: "sctl", Value: m.Sctl})
	}
	if m.OcspUpdate != "" {
		v := "off"
		if m.OcspUpdate == "enabled" {
			v = "on"
		}
		options = append(options, &params.BindOptionOnOff{Name: "ocsp-update", Value: v})
	}

	u.Params = options
	return u
}

func GetSSLFrontUseByNumber(number int64, parentType string, parentName string, p parser.Parser) *models.SSLFrontUse {
	uses, err := ParseSSLFrontUses(parentType, parentName, p)
	if err != nil || int64(len(uses)) < number+1 {
		return nil
	}
	return uses[number]
}
