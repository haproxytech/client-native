// Copyright 2025 HAProxy Technologies
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
	"sort"
	"strings"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

type Acme interface {
	GetAcmeProviders(transactionID string) (int64, models.AcmeProviders, error)
	GetAcmeProvider(name, transactionID string) (int64, *models.AcmeProvider, error)
	CreateAcmeProvider(data *models.AcmeProvider, transactionID string, version int64) error
	EditAcmeProvider(name string, data *models.AcmeProvider, transactionID string, version int64) error
	DeleteAcmeProvider(name, transactionID string, version int64) error
}

func (c *client) GetAcmeProviders(transactionID string) (int64, models.AcmeProviders, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	names, err := p.SectionsGet(parser.Acme)
	if err != nil {
		return v, nil, err
	}

	acmes := make(models.AcmeProviders, 0, len(names))

	for _, name := range names {
		a, err := ParseAcmeProvider(p, name)
		if err == nil {
			acmes = append(acmes, a)
		}
	}

	return v, acmes, nil
}

func (c *client) GetAcmeProvider(name, transactionID string) (int64, *models.AcmeProvider, error) {
	p, err := c.GetParser(transactionID)
	if err != nil {
		return 0, nil, err
	}

	v, err := c.GetVersion(transactionID)
	if err != nil {
		return 0, nil, err
	}

	if !p.SectionExists(parser.Acme, name) {
		return v, nil, NewConfError(ErrObjectDoesNotExist,
			fmt.Sprintf("%s section '%s' does not exist", AcmeParentName, name))
	}

	acme, err := ParseAcmeProvider(p, name)
	if err != nil {
		return 0, nil, err
	}

	return v, acme, nil
}

func (c *client) DeleteAcmeProvider(name, transactionID string, version int64) error {
	return c.deleteSection(parser.Acme, name, transactionID, version)
}

func (c *client) CreateAcmeProvider(data *models.AcmeProvider, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if p.SectionExists(parser.Acme, data.Name) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s already exists", parser.Acme, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = p.SectionsCreate(parser.Acme, data.Name); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	if err = SerializeAcmeProvider(p, data); err != nil {
		return c.HandleError(data.Name, "", "", t, transactionID == "", err)
	}

	return c.SaveData(p, t, transactionID == "")
}

func (c *client) EditAcmeProvider(name string, data *models.AcmeProvider, transactionID string, version int64) error {
	if c.UseModelsValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return NewConfError(ErrValidationError, validationErr.Error())
		}
	}

	if data.Name == "" {
		data.Name = name
	}

	p, t, err := c.loadDataForChange(transactionID, version)
	if err != nil {
		return err
	}

	if !p.SectionExists(parser.Acme, data.Name) {
		e := NewConfError(ErrObjectAlreadyExists, fmt.Sprintf("%s %s does not exists", parser.Acme, data.Name))
		return c.HandleError(data.Name, "", "", t, transactionID == "", e)
	}

	if err = SerializeAcmeProvider(p, data); err != nil {
		return err
	}

	return c.SaveData(p, t, transactionID == "")
}

func ParseAcmeProvider(p parser.Parser, name string) (*models.AcmeProvider, error) {
	acme := &models.AcmeProvider{Name: name}

	if data, err := p.SectionGet(parser.Acme, name); err == nil {
		d, ok := data.(types.Section)
		if ok {
			acme.Metadata = parseMetadata(d.Comment)
		}
	}

	var varsStr string

	stringAttr := map[string]*string{
		"account-key":   &acme.AccountKey,
		"acme-provider": &acme.AcmeProvider,
		"acme-vars":     &varsStr,
		"challenge":     &acme.Challenge,
		"contact":       &acme.Contact,
		"curves":        &acme.Curves,
		"directory":     &acme.Directory,
		"keytype":       &acme.Keytype,
		"map":           &acme.Map,
		"reuse-key":     &acme.ReuseKey,
	}

	for kw, dest := range stringAttr {
		val, err := p.Get(parser.Acme, name, kw)
		if err != nil {
			if errors.Is(err, parser_errors.ErrFetch) {
				continue
			}
			return nil, err
		}
		str, ok := val.(*types.StringC)
		if !ok {
			return nil, misc.CreateTypeAssertError(kw)
		}
		*dest = str.Value
	}

	acme.ReuseKey = onOff(acme.ReuseKey)

	// bits
	val, err := p.Get(parser.Acme, name, "bits")
	if err != nil {
		if !errors.Is(err, parser_errors.ErrFetch) {
			return nil, err
		}
	} else {
		ic, ok := val.(*types.Int64C)
		if !ok {
			return nil, misc.CreateTypeAssertError("bits")
		}
		acme.Bits = misc.Ptr(ic.Value)
	}

	// acme-vars
	acme.AcmeVars = ParseAcmeVars(varsStr)

	return acme, nil
}

func SerializeAcmeProvider(p parser.Parser, acme *models.AcmeProvider) error {
	if acme == nil {
		return fmt.Errorf("empty %s section", AcmeParentName)
	}

	if acme.Metadata != nil {
		comment, err := serializeMetadata(acme.Metadata)
		if err != nil {
			return err
		}
		if err := p.SectionCommentSet(parser.Acme, acme.Name, comment); err != nil {
			return err
		}
	}

	acmeVars, err := serializeAcmeVars(acme.AcmeVars)
	if err != nil {
		return fmt.Errorf("acme %s: %w", acme.Name, err)
	}

	stringAttr := map[string]string{
		"account-key":   acme.AccountKey,
		"acme-provider": acme.AcmeProvider,
		"acme-vars":     acmeVars,
		"challenge":     acme.Challenge,
		"contact":       acme.Contact,
		"curves":        acme.Curves,
		"directory":     acme.Directory,
		"keytype":       acme.Keytype,
		"map":           acme.Map,
		"reuse-key":     onOff(acme.ReuseKey),
	}

	for kw, val := range stringAttr {
		if val != "" {
			if err := p.Set(parser.Acme, acme.Name, kw, types.StringC{Value: val}); err != nil {
				return err
			}
		} else {
			_ = p.Delete(parser.Acme, acme.Name, kw)
		}
	}

	if acme.Bits != nil && *acme.Bits != 0 {
		if err := p.Set(parser.Acme, acme.Name, "bits", types.Int64C{Value: *acme.Bits}); err != nil {
			return err
		}
	} else {
		_ = p.Delete(parser.Acme, acme.Name, "bits")
	}

	return nil
}

// acme-vars "key=value,foo=\"bar baz\""
func serializeAcmeVars(vars map[string]string) (string, error) {
	if len(vars) == 0 {
		return "", nil
	}
	var sb strings.Builder
	first := true

	// Extract and sort the keys
	keys := make([]string, 0, len(vars))
	for name := range vars {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	sb.WriteByte('"')
	for _, k := range keys {
		v := vars[k]
		if len(k) == 0 {
			continue
		}
		if !acmeValidKey(k) {
			return "", fmt.Errorf("acme-vars: invalid character found in key '%s'", k)
		}
		if first {
			first = false
		} else {
			sb.WriteByte(',')
		}
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(acmeVarEscape(v))
	}
	sb.WriteByte('"')

	return sb.String(), nil
}

// Exported because used in dataplaneapi.
func ParseAcmeVars(vars string) map[string]string {
	n := len(vars)
	if n == 0 {
		return nil
	}

	if vars[0] == '"' && vars[n-1] == '"' {
		vars = vars[1 : n-1]
	}

	vars = strings.TrimSpace(vars)
	if len(vars) == 0 {
		return nil
	}

	vlist := acmeVarSplit(vars)
	vmap := make(map[string]string, len(vlist))
	for _, keyval := range vlist {
		if k, v, found := strings.Cut(strings.TrimSpace(keyval), "="); found {
			if len(k) > 0 {
				vmap[k] = acmeVarUnescape(v)
			}
		}
	}

	if len(vmap) == 0 {
		return nil
	}
	return vmap
}

// Split string by ',' but not escaped commas "\,".
func acmeVarSplit(s string) []string {
	s = strings.ReplaceAll(s, `\,`, "\x00")
	tokens := strings.Split(s, ",")
	for i, token := range tokens {
		tokens[i] = strings.ReplaceAll(token, "\x00", `\,`)
	}
	return tokens
}

func acmeVarEscape(s string) string {
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, `,`, `\,`)
	return s
}

func acmeVarUnescape(s string) string {
	s = strings.ReplaceAll(s, `\"`, `"`)
	s = strings.ReplaceAll(s, `\,`, `,`)
	return s
}

// Variable keys must also be valid Go variable names.
func acmeValidKey(key string) bool {
	for _, c := range key {
		match := ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z') ||
			('0' <= c && c <= '9') || c == '_'
		if !match {
			return false
		}
	}
	return true
}

func onOff(s string) string {
	switch len(s) {
	case 2: // on
		return "enabled"
	case 3: // off
		return "disabled"
	case 7: // enabled
		return "on"
	case 8: // disabled
		return "off"
	default:
		return s
	}
}
