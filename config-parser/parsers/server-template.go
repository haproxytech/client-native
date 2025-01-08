package parsers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/params"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type ServerTemplate struct {
	data        []types.ServerTemplate
	preComments []string // comments that appear before the actual line
}

func (h *ServerTemplate) parse(line string, parts []string, comment string) (*types.ServerTemplate, error) {
	if len(parts) < 4 {
		return nil, &errors.ParseError{Parser: "ServerTemplate", Line: line}
	}

	data := &types.ServerTemplate{}
	data.Prefix = parts[1]
	data.NumOrRange = parts[2]
	data.Comment = comment

	address, p, found := common.CutRight(parts[3], ":")
	if found {
		if port, err := strconv.ParseInt(p, 10, 64); err == nil {
			data.Fqdn = address
			data.Port = port
		}
	} else {
		data.Fqdn = parts[3]
	}
	if len(parts) >= 4 {
		sp, err := params.ParseServerOptions(parts[4:])
		if err != nil {
			return nil, err
		}
		data.Params = sp
	}
	return data, nil
}

func (h *ServerTemplate) Result() ([]common.ReturnResultLine, error) {
	if len(h.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(h.data))
	for index, template := range h.data {
		var sb strings.Builder
		sb.WriteString("server-template ")
		sb.WriteString(template.Prefix)
		sb.WriteString(" ")
		sb.WriteString(template.NumOrRange)
		sb.WriteString(" ")
		sb.WriteString(template.Fqdn)
		if template.Port != 0 {
			sb.WriteString(fmt.Sprintf(":%d", template.Port))
		}
		params := params.ServerOptionsString(template.Params)
		if params != "" {
			sb.WriteString(" ")
			sb.WriteString(params)
		}
		result[index] = common.ReturnResultLine{
			Data:    sb.String(),
			Comment: template.Comment,
		}
	}
	return result, nil
}
