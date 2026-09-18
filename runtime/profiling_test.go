package runtime

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	native_errors "github.com/haproxytech/client-native/v6/errors"
	"github.com/haproxytech/client-native/v6/models"
)

// profilingStatusText renders "show profiling status" output the way HAProxy
// prints it: labels padded to 36 columns, values to 14, then a "#" hint.
func profilingStatusText(tasks, memory string) string {
	return fmt.Sprintf("Per-task CPU profiling              : %-14s# set profiling tasks {on|auto|off}\n"+
		"Memory usage profiling              : %-14s# set profiling memory {on|off}\n", tasks, memory)
}

func TestParseProfilingStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    models.Profiling
		wantErr bool
	}{
		{
			name:  "auto-off / off",
			input: profilingStatusText("auto-off", "off"),
			want:  models.Profiling{Tasks: "auto", TasksActive: false, Memory: "off"},
		},
		{
			name:  "auto-on / on",
			input: profilingStatusText("auto-on", "on"),
			want:  models.Profiling{Tasks: "auto", TasksActive: true, Memory: "on"},
		},
		{
			name:  "on / off",
			input: profilingStatusText("on", "off"),
			want:  models.Profiling{Tasks: "on", TasksActive: true, Memory: "off"},
		},
		{
			name:  "off / on",
			input: profilingStatusText("off", "on"),
			want:  models.Profiling{Tasks: "off", TasksActive: false, Memory: "on"},
		},
		{
			name:    "invalid tasks value",
			input:   profilingStatusText("weird", "off"),
			wantErr: true,
		},
		{
			name:    "empty response",
			input:   "",
			wantErr: true,
		},
		{
			name:    "unknown command",
			input:   "Unknown command: 'show profiling'...",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseProfilingStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseProfilingStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !errors.Is(err, native_errors.ErrGeneral) {
					t.Errorf("parseProfilingStatus() error = %v, want %v", err, native_errors.ErrGeneral)
				}
				return
			}
			if !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("parseProfilingStatus() = %v, want %v", got, &tt.want)
			}
		})
	}
}

func TestProfilingDumpOptions_args(t *testing.T) {
	tests := []struct {
		name      string
		opts      ProfilingDumpOptions
		want      string
		wantErr   bool
		wantErrIs error // non-nil: the error must wrap it
	}{
		{
			name: "zero value means HAProxy defaults",
			opts: ProfilingDumpOptions{},
			want: "",
		},
		{
			name: "sort by address with aggregation and max lines",
			opts: ProfilingDumpOptions{Sort: ProfilingSortAddress, Aggregate: true, MaxLines: 20},
			want: " byaddr aggr 20",
		},
		{
			name:      "unknown sort",
			opts:      ProfilingDumpOptions{Sort: "bogus"},
			wantErr:   true,
			wantErrIs: native_errors.ErrGeneral,
		},
		{
			name:      "negative max lines",
			opts:      ProfilingDumpOptions{MaxLines: -1},
			wantErr:   true,
			wantErrIs: native_errors.ErrGeneral,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.opts.args()
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfilingDumpOptions.args() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Errorf("ProfilingDumpOptions.args() error = %v, wantErrIs %v", err, tt.wantErrIs)
			}
			if got != tt.want {
				t.Errorf("ProfilingDumpOptions.args() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSingleRuntime_ShowProfilingStatus(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	tests := []struct {
		name           string
		fields         fields
		want           *models.Profiling
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "profiling status",
			fields: fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			want: &models.Profiling{
				Tasks:       "auto",
				TasksActive: false,
				Memory:      "off",
			},
			socketResponse: map[string]string{
				"show profiling status\n": profilingStatusText("auto-off", "off"),
			},
		},
		{
			name:    "unknown command",
			fields:  fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			wantErr: true,
			socketResponse: map[string]string{
				"show profiling status\n": "[3]: Unknown command: 'show profiling'...",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.masterWorkerMode)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			got, err := s.ShowProfilingStatus()
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowProfilingStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SingleRuntime.ShowProfilingStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleRuntime_ShowProfilingTasks(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	tests := []struct {
		name           string
		fields         fields
		opts           ProfilingDumpOptions
		want           string
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:   "dump bytime 10",
			fields: fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			opts:   ProfilingDumpOptions{Sort: ProfilingSortTime, MaxLines: 10},
			want: `# calls  tot_us  avg_us  function
2        1234    617     p_ha_cli_conn
7        891     127     p_ha_srv_conn`,
			socketResponse: map[string]string{
				"show profiling tasks bytime 10\n": `# calls  tot_us  avg_us  function
2        1234    617     p_ha_cli_conn
7        891     127     p_ha_srv_conn
`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.masterWorkerMode)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			got, err := s.ShowProfilingTasks(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowProfilingTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != strings.TrimSpace(tt.want) {
				t.Errorf("SingleRuntime.ShowProfilingTasks() = %q, want %q", got, strings.TrimSpace(tt.want))
			}
		})
	}
}

func TestSingleRuntime_ShowProfilingMemory(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	tests := []struct {
		name           string
		fields         fields
		opts           ProfilingDumpOptions
		want           string
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			// a build without USE_MEMORY_PROFILING returns an empty string
			name:           "memory profiling not compiled in",
			fields:         fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			opts:           ProfilingDumpOptions{},
			want:           "",
			socketResponse: map[string]string{"show profiling memory\n": ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.masterWorkerMode)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			got, err := s.ShowProfilingMemory(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.ShowProfilingMemory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SingleRuntime.ShowProfilingMemory() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSingleRuntime_SetProfilingTasks(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	tests := []struct {
		name           string
		fields         fields
		mode           string
		wantErr        bool
		wantErrIs      error // non-nil: the error must wrap it
		socketResponse map[string]string
	}{
		{
			name:           "on",
			fields:         fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			mode:           "on",
			socketResponse: map[string]string{"set profiling tasks on\n": ""},
		},
		{
			name:           "auto",
			fields:         fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			mode:           "auto",
			socketResponse: map[string]string{"set profiling tasks auto\n": ""},
		},
		{
			name:           "off",
			fields:         fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			mode:           "off",
			socketResponse: map[string]string{"set profiling tasks off\n": ""},
		},
		{
			name:      "invalid mode",
			fields:    fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			mode:      "enabled",
			wantErr:   true,
			wantErrIs: native_errors.ErrGeneral,
			// no response registered: proves the mode is rejected before any socket call
			socketResponse: map[string]string{},
		},
		{
			name:    "haproxy rejects the mode",
			fields:  fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			mode:    "on",
			wantErr: true,
			socketResponse: map[string]string{
				"set profiling tasks on\n": "[3]: Expects either 'on', 'auto', or 'off'.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.masterWorkerMode)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			err = s.SetProfilingTasks(tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.SetProfilingTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Errorf("SingleRuntime.SetProfilingTasks() error = %v, wantErrIs %v", err, tt.wantErrIs)
			}
		})
	}
}

func TestSingleRuntime_SetProfilingMemory(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	tests := []struct {
		name           string
		fields         fields
		mode           string
		wantErr        bool
		wantErrIs      error // non-nil: the error must wrap it
		socketResponse map[string]string
	}{
		{
			name:           "on",
			fields:         fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			mode:           "on",
			socketResponse: map[string]string{"set profiling memory on\n": ""},
		},
		{
			name:           "off",
			fields:         fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			mode:           "off",
			socketResponse: map[string]string{"set profiling memory off\n": ""},
		},
		{
			name:           "invalid mode",
			fields:         fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			mode:           "auto",
			wantErr:        true,
			wantErrIs:      native_errors.ErrGeneral,
			socketResponse: map[string]string{},
		},
		{
			name:    "memory profiling not compiled in",
			fields:  fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			mode:    "on",
			wantErr: true,
			socketResponse: map[string]string{
				"set profiling memory on\n": "[3]: Memory profiling not compiled in.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.masterWorkerMode)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			err = s.SetProfilingMemory(tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.SetProfilingMemory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Errorf("SingleRuntime.SetProfilingMemory() error = %v, wantErrIs %v", err, tt.wantErrIs)
			}
		})
	}
}

func TestSingleRuntime_SetProfiling(t *testing.T) {
	haProxy := NewHAProxyMock(t)
	haProxy.Start()
	defer haProxy.Stop()

	type fields struct {
		socketPath       string
		masterWorkerMode bool
	}
	tests := []struct {
		name           string
		fields         fields
		profiling      *models.Profiling
		wantErr        bool
		socketResponse map[string]string
	}{
		{
			name:           "tasks only",
			fields:         fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			profiling:      &models.Profiling{Tasks: "auto"},
			socketResponse: map[string]string{"set profiling tasks auto\n": ""},
		},
		{
			name:      "memory profiling not compiled in",
			fields:    fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			profiling: &models.Profiling{Tasks: "on", Memory: "on"},
			wantErr:   true,
			socketResponse: map[string]string{
				"set profiling tasks on\n":  "",
				"set profiling memory on\n": "[3]: Memory profiling not compiled in.",
			},
		},
		{
			name:           "nil applies nothing",
			fields:         fields{socketPath: haProxy.Addr().String(), masterWorkerMode: false},
			profiling:      nil,
			socketResponse: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haProxy.SetResponses(&tt.socketResponse)
			s := &SingleRuntime{}
			err := s.Init(tt.fields.socketPath, tt.fields.masterWorkerMode)
			if err != nil {
				t.Errorf("SingleRuntime.Init() error = %v", err)
				return
			}
			err = s.SetProfiling(tt.profiling)
			if (err != nil) != tt.wantErr {
				t.Errorf("SingleRuntime.SetProfiling() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
