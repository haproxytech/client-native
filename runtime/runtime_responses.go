package runtime

import (
	"fmt"
)

type RuntimeResponse struct {
	HTTPLikeStatus HTTPLikeStatus
	Message        string
}

func (r RuntimeResponse) String() {
	fmt.Printf("[%d] %s", r.HTTPLikeStatus, r.Message)
}
