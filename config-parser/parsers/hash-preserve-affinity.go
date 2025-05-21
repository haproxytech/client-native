package parsers

import (
	"fmt"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type HashPreserveAffinity struct {
	data        *types.StringC
	preComments []string // comments that appear before the actual line
}

func (p *HashPreserveAffinity) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] == "hash-preserve-affinity" {
		if len(parts) < 2 {
			return "", &errors.ParseError{Parser: "HashPreserveAffinity", Line: line, Message: "Parse error"}
		}

		if parts[1] == "always" || parts[1] == "maxconn" || parts[1] == "maxqueue" {
			p.data = &types.StringC{
				Value:   parts[1],
				Comment: comment,
			}
			return "", nil
		}

		return "", &errors.ParseError{Parser: "HashPreserveAffinity", Line: line, Message: fmt.Sprintf("Unknown data %q", parts[1])}
	}
	return "", &errors.ParseError{Parser: "HashPreserveAffinity", Line: line}
}

func (p *HashPreserveAffinity) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    "hash-preserve-affinity " + p.data.Value,
			Comment: p.data.Comment,
		},
	}, nil
}
