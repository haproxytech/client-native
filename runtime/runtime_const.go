package runtime

import (
	"strconv"
	"strings"
	"sync"
)

//HTTPLikeStatus status simulates http status responses
type HTTPLikeStatus int

//HTTPLikeStatus possible response
const (
	HTTPLikeStatusOK                  HTTPLikeStatus = 200
	HTTPLikeStatusCreated             HTTPLikeStatus = 201
	HTTPLikeStatusAccepted            HTTPLikeStatus = 202
	HTTPLikeStatusNoContent           HTTPLikeStatus = 204
	HTTPLikeStatusNotModified         HTTPLikeStatus = 304
	HTTPLikeStatusBadRequest          HTTPLikeStatus = 400
	HTTPLikeStatusUnauthorized        HTTPLikeStatus = 401
	HTTPLikeStatusForbidden           HTTPLikeStatus = 403
	HTTPLikeStatusNotFound            HTTPLikeStatus = 404
	HTTPLikeStatusTimeout             HTTPLikeStatus = 408
	HTTPLikeStatusConflict            HTTPLikeStatus = 409
	HTTPLikeStatusInternalServerError HTTPLikeStatus = 500
	HTTPLikeStatusNotImplemented      HTTPLikeStatus = 501
)

var possibleStates map[string]struct{}
var oncePossibleStates sync.Once

//ServerStateValid checks if server state is valid
func ServerStateValid(state string) bool {
	oncePossibleStates.Do(func() {
		possibleStates = map[string]struct{}{
			"ready": {},
			"drain": {},
			"maint": {},
		}
	})
	_, ok := possibleStates[state]
	return ok
}

//ServerHealthValid checks if server state is valid
func ServerHealthValid(state string) bool {
	oncePossibleStates.Do(func() {
		possibleStates = map[string]struct{}{
			"on":       {},
			"stopping": {},
			"down":     {},
		}
	})
	_, ok := possibleStates[state]
	return ok
}

//ServerWeightValid checks if server state is valid
func ServerWeightValid(weight string) bool {
	var n int64
	var err error
	if strings.HasSuffix(weight, "%") {
		percent := strings.TrimSuffix(weight, "%")
		if n, err = strconv.ParseInt(percent, 10, 64); err != nil {
			return false
		}
		if n > -1 && n < 101 {
			return true
		}
	}
	if n, err = strconv.ParseInt(weight, 10, 64); err != nil {
		return false
	}
	if n > -1 && n < 257 {
		return true
	}
	return false
}
