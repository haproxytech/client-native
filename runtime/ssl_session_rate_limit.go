package runtime

import (
	"fmt"
	"strings"

	native_errors "github.com/haproxytech/client-native/v6/errors"
)

// SetRateLimitSSLSessionGlobal sets the SSL session global rate limit
func (s *SingleRuntime) SetRateLimitSSLSessionGlobal(value uint64) error {
	cmd := fmt.Sprintf("set rate-limit ssl-sessions global %d\n", value)
	response, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return fmt.Errorf("%s %w", err.Error(), native_errors.ErrNotFound)
	}
	if strings.Contains(response, "Value out of range") {
		return fmt.Errorf("%s %w", strings.TrimSpace(response), native_errors.ErrGeneral)
	}
	return nil
}
