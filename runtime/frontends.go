package runtime

import (
	"fmt"
	"strconv"
)

//SetFrontendMaxConn set maxconn for frontend
func (s *SingleRuntime) SetFrontendMaxConn(frontend string, maxconn int) error {
	cmd := fmt.Sprintf("set maxconn frontend %s %s", frontend, strconv.FormatInt(int64(maxconn), 10))
	return s.Execute(cmd)
}
