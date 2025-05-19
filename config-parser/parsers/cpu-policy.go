package parsers

import (
	"fmt"
	"slices"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type CPUPolicy struct {
	data        *types.StringC
	preComments []string // comments that appear before the actual line
}

func (p *CPUPolicy) Parse(line string, parts []string, comment string) (string, error) {
	if parts[0] == "cpu-policy" {
		if len(parts) < 2 {
			return "", &errors.ParseError{Parser: "CPUPolicy", Line: line, Message: "Parse error"}
		}

		policyList := []string{
			"none",
			"efficiency",
			"first-usable-node",
			"group-by-2-ccx",
			"group-by-2-clusters",
			"group-by-3-ccx",
			"group-by-3-clusters",
			"group-by-4-ccx",
			"group-by-4-cluster",
			"group-by-ccx",
			"group-by-cluster",
			"performance",
			"resource",
		}
		if slices.Contains(policyList, parts[1]) {
			p.data = &types.StringC{
				Value:   parts[1],
				Comment: comment,
			}
			return "", nil
		}

		return "", &errors.ParseError{Parser: "CPUPolicy", Line: line, Message: fmt.Sprintf("Unknown policy %q. Must be one of %v", parts[1], policyList)}
	}
	return "", &errors.ParseError{Parser: "CPUPolicy", Line: line}
}

func (p *CPUPolicy) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}
	return []common.ReturnResultLine{
		{
			Data:    "cpu-policy " + p.data.Value,
			Comment: p.data.Comment,
		},
	}, nil
}
