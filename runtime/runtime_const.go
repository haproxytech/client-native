package runtime

import (
	"strconv"
	"strings"
	"sync"
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
