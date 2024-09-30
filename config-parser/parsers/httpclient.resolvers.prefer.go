package parsers

import (
	"fmt"

	"github.com/haproxytech/client-native/v5/config-parser/common"
	"github.com/haproxytech/client-native/v5/config-parser/errors"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

type HTTPClientResolversPrefer struct {
	data        *types.HTTPClientResolversPrefer
	preComments []string
}

func (p *HTTPClientResolversPrefer) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) == 2 && parts[0] == "httpclient.resolvers.prefer" {
		p.data = &types.HTTPClientResolversPrefer{}
		switch parts[1] {
		case "ipv4", "ipv6":
			p.data.Type = parts[1]
		default:
			return "", &errors.ParseError{Parser: "httpclient.resolvers.prefer", Line: line}
		}
		return "", nil
	}
	return "", &errors.ParseError{Parser: "httpclient.resolvers.prefer", Line: line}
}

func (p *HTTPClientResolversPrefer) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil || len(p.data.Type) == 0 {
		return nil, errors.ErrFetch
	}
	data := fmt.Sprintf("httpclient.resolvers.prefer %s", p.data.Type)
	return []common.ReturnResultLine{
		{
			Data: data,
		},
	}, nil
}
