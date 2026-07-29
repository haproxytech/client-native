package configuration

import (
	"strings"
	"testing"

	parser "github.com/haproxytech/client-native/v5/config-parser"
	"github.com/haproxytech/client-native/v5/config-parser/options"
	"github.com/haproxytech/client-native/v5/models"

	"github.com/stretchr/testify/require"
)

// statsRefreshRoundTrip parses a frontend holding a single 'stats refresh'
// keyword and serializes it back, returning the regenerated keyword and the
// milliseconds the model held in between.
func statsRefreshRoundTrip(t *testing.T, delay string) (string, *int64) {
	t.Helper()

	pIn, err := parser.New(options.String("frontend stats\n  stats refresh " + delay + "\n"))
	require.NoError(t, err)

	fe := &models.Frontend{Name: "stats"}
	require.NoError(t, ParseSection(fe, parser.Frontends, "stats", pIn))
	require.NotNil(t, fe.StatsOptions)

	pOut, err := parser.New(options.String("frontend stats\n"))
	require.NoError(t, err)
	require.NoError(t, CreateEditSection(fe, parser.Frontends, "stats", pOut))

	var got string
	for _, line := range strings.Split(pOut.String(), "\n") {
		if line = strings.TrimSpace(line); strings.HasPrefix(line, "stats refresh ") {
			got = strings.TrimPrefix(line, "stats refresh ")
		}
	}
	return got, fe.StatsOptions.StatsRefreshDelay
}

// TestStatsRefreshDelayDefaultsToSeconds checks that a 'stats refresh' value
// written without a suffix is read as seconds, as the keyword documents, and
// that the delay is always serialized with an explicit unit. The model stores
// milliseconds, so writing the bare number would scale the value by 1000 on
// reload and make the round trip non-idempotent.
func TestStatsRefreshDelayDefaultsToSeconds(t *testing.T) {
	tests := []struct {
		name     string
		delay    string
		expected string
		ms       int64
	}{
		{name: "no suffix means seconds", delay: "30", expected: "30000ms", ms: 30000},
		{name: "large value with no suffix", delay: "30000", expected: "30000000ms", ms: 30000000},
		{name: "seconds", delay: "30s", expected: "30000ms", ms: 30000},
		{name: "milliseconds", delay: "5000ms", expected: "5000ms", ms: 5000},
		{name: "minutes", delay: "2m", expected: "120000ms", ms: 120000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ms := statsRefreshRoundTrip(t, tt.delay)
			require.Equal(t, tt.expected, got)
			require.NotNil(t, ms)
			require.Equal(t, tt.ms, *ms)

			// serializing the regenerated value again must be a no-op,
			// otherwise every reload scales the delay by 1000
			stable, msAgain := statsRefreshRoundTrip(t, got)
			require.Equal(t, tt.expected, stable)
			require.NotNil(t, msAgain)
			require.Equal(t, tt.ms, *msAgain)
		})
	}
}
