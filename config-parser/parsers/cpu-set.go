package parsers

import (
	"fmt"
	"slices"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

// Supported directives for cpu-set keyword
var cpuSetSupportedDirectives = []string{ //nolint:gochecknoglobals
	CPUSetResetDirective,
	"drop-cpu",
	"only-cpu",
	"drop-node",
	"only-node",
	"drop-cluster",
	"only-cluster",
	"drop-core",
	"only-core",
	"drop-thread",
	"only-thread",
}

const CPUSetResetDirective = "reset"

type CPUSet struct {
	data        []types.CPUSet
	preComments []string // comments that appear before the actual line
}

func (c *CPUSet) parse(line string, parts []string, comment string) (*types.CPUSet, error) {
	// Check if we have at least the command (cpu-set) and directive
	if len(parts) < 2 {
		return nil, &errors.ParseError{Parser: "CPUSet", Line: line, Message: "Parse error"}
	}

	// Validate the directive
	directive := parts[1]
	if !slices.Contains(cpuSetSupportedDirectives, directive) {
		return nil, &errors.ParseError{Parser: "CPUSet", Line: line, Message: fmt.Sprintf("Unknown directive %q, supported directives %v", directive, cpuSetSupportedDirectives)}
	}

	// Handle reset directive (no 'set' parameters)
	if directive == CPUSetResetDirective {
		if len(parts) != 2 {
			return nil, &errors.ParseError{Parser: "CPUSet", Line: line, Message: "cpu-set reset directive does not accept 'set' parameter"}
		}
		return &types.CPUSet{Directive: directive, Comment: comment}, nil
	}

	// Handle other directives (require one parameter)
	if len(parts) != 3 {
		return nil, &errors.ParseError{Parser: "CPUSet", Line: line, Message: fmt.Sprintf("cpu-set %s missing 'set' parameter", directive)}
	}

	return &types.CPUSet{Directive: directive, Set: parts[2], Comment: comment}, nil
}

func (c *CPUSet) Result() ([]common.ReturnResultLine, error) {
	if len(c.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(c.data))
	for index, cpuSet := range c.data {
		data := "cpu-set " + cpuSet.Directive
		if cpuSet.Directive != CPUSetResetDirective { // reset is the only directive that has no 'set' parameter
			data += " " + cpuSet.Set
		}
		result[index] = common.ReturnResultLine{
			Data:    data,
			Comment: cpuSet.Comment,
		}
	}
	return result, nil
}
