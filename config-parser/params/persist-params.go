package params

import (
	"strings"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
)

type PersistParams interface {
	String() string
	Parse(parts []string) (PersistParams, error)
}

type PersistRdpCookie struct {
	Name string
}

func (r *PersistRdpCookie) String() string {
	var result strings.Builder
	if r.Name != "" {
		result.WriteString("(")
		result.WriteString(r.Name)
		result.WriteString(")")
	}
	return result.String()
}

func (r *PersistRdpCookie) Parse(parts []string) (PersistParams, error) {
	if len(parts) > 0 {
		split := common.StringSplitIgnoreEmpty(parts[0], '(', ')')
		if len(split) < 2 {
			return nil, errors.ErrInvalidData
		}
		r.Name = split[1]
	}
	return r, nil
}
