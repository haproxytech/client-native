package parsers

import (
	"fmt"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type HTTPClientSSLVerify struct {
	data        *types.HTTPClientSSLVerify
	preComments []string
}

func (p *HTTPClientSSLVerify) Parse(line string, parts []string, comment string) (string, error) {
	if len(parts) > 0 && parts[0] == "httpclient.ssl.verify" {
		switch len(parts) {
		case 2:
			switch parts[1] {
			case "none", "required":
				p.data = &types.HTTPClientSSLVerify{}
				p.data.Type = parts[1]
			default:
				return "", &errors.ParseError{Parser: "httpclient.ssl.verify", Line: line}
			}
			return "", nil
		case 1:
			p.data = &types.HTTPClientSSLVerify{}
			return "", nil
		default:
			return "", &errors.ParseError{Parser: "httpclient.ssl.verify", Line: line}
		}
	}
	return "", &errors.ParseError{Parser: "httpclient.ssl.verify", Line: line}
}

func (p *HTTPClientSSLVerify) Result() ([]common.ReturnResultLine, error) {
	if p.data == nil {
		return nil, errors.ErrFetch
	}
	data := "httpclient.ssl.verify"
	if len(p.data.Type) > 0 {
		data = fmt.Sprintf("httpclient.ssl.verify %s", p.data.Type)
	}
	return []common.ReturnResultLine{
		{
			Data: data,
		},
	}, nil
}
