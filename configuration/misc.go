package configuration

import (
	"encoding/json"

	"github.com/haproxytech/client-native/v6/config-parser/parsers/http/actions"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
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

func parseMetadata(comment string) map[string]interface{} {
	if comment == "" {
		return nil
	}
	metadata := make(map[string]interface{})
	err := json.Unmarshal([]byte(comment), &metadata)
	if err != nil {
		metadata["comment"] = comment
		return metadata
	}
	return metadata
}

func serializeMetadata(metadata map[string]interface{}) (string, error) {
	if metadata == nil {
		return "", nil
	}
	b, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
