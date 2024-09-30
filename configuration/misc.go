package configuration

import (
	"github.com/haproxytech/client-native/v5/config-parser/parsers/http/actions"

	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

func actionHdr2ModelHdr(hdrs []*actions.Hdr) []*models.ReturnHeader {
	if len(hdrs) == 0 {
		return nil
	}
	headers := []*models.ReturnHeader{}
	for _, h := range hdrs {
		hdr := models.ReturnHeader{
			Fmt:  misc.Ptr(h.Fmt),
			Name: misc.Ptr(h.Name),
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
