package configuration

import (
	"github.com/haproxytech/config-parser/v4/parsers/http/actions"

	"github.com/haproxytech/client-native/v4/models"
)

func actionHdr2ModelHdr(hdrs []*actions.Hdr) []*models.ReturnHeader {
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

func modelHdr2ActionHdr(hdrs []*models.ReturnHeader) []*actions.Hdr {
	if len(hdrs) == 0 {
		return nil
	}
	headers := []*actions.Hdr{}
	for _, h := range hdrs {
		hdr := actions.Hdr{
			Name: *h.Name,
			Fmt:  *h.Fmt,
		}
		headers = append(headers, &hdr)
	}
	return headers
}
