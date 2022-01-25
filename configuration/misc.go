package configuration

import (
	"github.com/haproxytech/client-native/v2/models"
	http_actions "github.com/haproxytech/config-parser/v4/parsers/http/actions"
)

func actionHdr2ModelHdr(hdrs []*http_actions.Hdr) []*models.ReturnHeader {
	if len(hdrs) == 0 {
		return nil
	}
	headers := []*models.ReturnHeader{}
	for _, h := range hdrs {
		hdr := models.ReturnHeader{
			Fmt:  &h.Fmt,
			Name: &h.Name,
		}
		headers = append(headers, &hdr)
	}
	return headers
}

func modelHdr2ActionHdr(hdrs []*models.ReturnHeader) []*http_actions.Hdr {
	if len(hdrs) == 0 {
		return nil
	}
	headers := []*http_actions.Hdr{}
	for _, h := range hdrs {
		hdr := http_actions.Hdr{
			Name: *h.Name,
			Fmt:  *h.Fmt,
		}
		headers = append(headers, &hdr)
	}
	return headers
}
