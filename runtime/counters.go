package runtime

// ClearCounters sets all max counters to zero using the runtime API
// To reset all counters see ClearAllCounters()
// REF: https://www.haproxy.com/documentation/haproxy-runtime-api/reference/clear-counters/
func (s *SingleRuntime) ClearCounters() error {
	_, err := s.ExecuteWithResponse("clear counters")
	return err
}

// ClearAllCounters sets all counters to zero using the runtime API
// REF: https://www.haproxy.com/documentation/haproxy-runtime-api/reference/clear-counters-all/
func (s *SingleRuntime) ClearAllCounters() error {
	_, err := s.ExecuteWithResponse("clear counters all")
	return err
}
