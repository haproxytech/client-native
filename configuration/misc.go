package configuration

import (
	"fmt"

	"github.com/haproxytech/client-native/v5/config-parser/common"
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

// validateReturnContent ensures a return-content value can be faithfully
// persisted. The value is emitted verbatim after the content-format keyword
// (string/lf-string/file/lf-file/errorfile), and HAProxy treats it as a single
// token: multi-word values must be quoted by the caller. When they are not, the
// serialized line (e.g. `... string Missing Auth`) fails to parse and is
// silently dropped on the next config reload, even though the transaction
// reports success. Reject such values here instead so the caller gets an error.
//
// The check mirrors the serializer's emit condition: content is only written
// when both the content and the format are non-empty and the format is not
// default-errorfiles, so anything that is not actually emitted cannot break the
// round-trip and is left alone.
func validateReturnContent(format, content string) error {
	if content == "" || format == "" || format == "default-errorfiles" {
		return nil
	}
	tokens, comment := common.StringSplitWithCommentIgnoreEmpty(content)
	if comment != "" || len(tokens) != 1 || tokens[0] != content {
		return NewConfError(ErrValidationError,
			fmt.Sprintf("invalid return content %q: multi-word values must be quoted", content))
	}
	return nil
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
