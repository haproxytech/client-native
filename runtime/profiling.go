// Copyright 2026 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package runtime

import (
	"fmt"
	"strings"

	native_errors "github.com/haproxytech/client-native/v6/errors"
	"github.com/haproxytech/client-native/v6/models"
)

// profiling modes, as accepted by "set profiling tasks|memory".
const (
	profilingOn   = "on"
	profilingOff  = "off"
	profilingAuto = "auto"
)

// labels of "show profiling status" output lines.
const (
	profilingLabelTasks  = "Per-task CPU profiling"
	profilingLabelMemory = "Memory usage profiling"
)

// ProfilingSort selects the ordering of "show profiling tasks|memory" output.
type ProfilingSort string

const (
	ProfilingSortUsage   ProfilingSort = "" // HAProxy default ordering
	ProfilingSortAddress ProfilingSort = "byaddr"
	ProfilingSortTime    ProfilingSort = "bytime"
	ProfilingSortContext ProfilingSort = "byctx"
)

// ProfilingDumpOptions are the optional arguments of "show profiling tasks|memory".
// The zero value means HAProxy defaults.
type ProfilingDumpOptions struct {
	Sort      ProfilingSort
	Aggregate bool // "aggr": aggregate by callee
	MaxLines  int  // 0 = unlimited
}

// args returns the CLI suffix of "show profiling tasks|memory", e.g. " byaddr aggr 20".
func (o ProfilingDumpOptions) args() (string, error) {
	var args strings.Builder

	switch o.Sort {
	case ProfilingSortUsage:
		// HAProxy default ordering, no argument.
	case ProfilingSortAddress, ProfilingSortTime, ProfilingSortContext:
		args.WriteString(" ")
		args.WriteString(string(o.Sort))
	default:
		return "", fmt.Errorf("invalid profiling sort %q: %w", o.Sort, native_errors.ErrGeneral)
	}

	if o.Aggregate {
		args.WriteString(" aggr")
	}

	if o.MaxLines != 0 {
		if o.MaxLines < 0 {
			return "", fmt.Errorf("invalid profiling max lines %d: %w", o.MaxLines, native_errors.ErrGeneral)
		}
		fmt.Fprintf(&args, " %d", o.MaxLines)
	}

	return args.String(), nil
}

// ShowProfilingStatus returns the profiling status, as reported by "show profiling status".
func (s *SingleRuntime) ShowProfilingStatus() (*models.Profiling, error) {
	resp, err := s.ExecuteWithResponse("show profiling status")
	if err != nil {
		return nil, err
	}
	return parseProfilingStatus(resp)
}

// ShowProfilingTasks returns the verbatim output of "show profiling tasks".
func (s *SingleRuntime) ShowProfilingTasks(opts ProfilingDumpOptions) (string, error) {
	args, err := opts.args()
	if err != nil {
		return "", err
	}
	return s.ExecuteWithResponse("show profiling tasks" + args)
}

// ShowProfilingMemory returns the verbatim output of "show profiling memory".
func (s *SingleRuntime) ShowProfilingMemory(opts ProfilingDumpOptions) (string, error) {
	args, err := opts.args()
	if err != nil {
		return "", err
	}
	return s.ExecuteWithResponse("show profiling memory" + args)
}

// SetProfilingTasks sets the per-task CPU profiling mode, one of "on", "auto" or "off".
func (s *SingleRuntime) SetProfilingTasks(mode string) error {
	switch mode {
	case profilingOn, profilingAuto, profilingOff:
	default:
		return fmt.Errorf("invalid profiling tasks mode %q: %w", mode, native_errors.ErrGeneral)
	}
	return s.Execute("set profiling tasks " + mode)
}

// SetProfilingMemory sets the memory usage profiling mode, one of "on" or "off".
func (s *SingleRuntime) SetProfilingMemory(mode string) error {
	switch mode {
	case profilingOn, profilingOff:
	default:
		return fmt.Errorf("invalid profiling memory mode %q: %w", mode, native_errors.ErrGeneral)
	}
	return s.Execute("set profiling memory " + mode)
}

// SetProfiling applies the profiling settings present in p, tasks first then memory.
// TasksActive is read-only and ignored.
func (s *SingleRuntime) SetProfiling(p *models.Profiling) error {
	if p == nil {
		return nil
	}
	if p.Tasks != "" {
		if err := s.SetProfilingTasks(p.Tasks); err != nil {
			return err
		}
	}
	if p.Memory != "" {
		if err := s.SetProfilingMemory(p.Memory); err != nil {
			return err
		}
	}
	return nil
}

// parseProfilingStatus parses the output of "show profiling status".
func parseProfilingStatus(resp string) (*models.Profiling, error) {
	p := &models.Profiling{}
	seen := false

	firstLine, _, _ := strings.Cut(resp, "\n")

	var parseErr error
	strings.SplitSeq(resp, "\n")(func(line string) bool {
		label, rest, ok := strings.Cut(line, ":")
		if !ok {
			return true
		}
		value, _, _ := strings.Cut(rest, "#")
		label = strings.TrimSpace(label)
		value = strings.TrimSpace(value)

		switch label {
		case profilingLabelTasks:
			seen = true
			switch value {
			case profilingOn:
				p.Tasks = profilingOn
				p.TasksActive = true
			case profilingOff:
				p.Tasks = profilingOff
				p.TasksActive = false
			case "auto-on":
				p.Tasks = profilingAuto
				p.TasksActive = true
			case "auto-off":
				p.Tasks = profilingAuto
				p.TasksActive = false
			default:
				parseErr = fmt.Errorf("invalid per-task CPU profiling value %q: %w", value, native_errors.ErrGeneral)
				return false
			}
		case profilingLabelMemory:
			seen = true
			switch value {
			case profilingOn, profilingOff:
				p.Memory = value
			default:
				parseErr = fmt.Errorf("invalid memory usage profiling value %q: %w", value, native_errors.ErrGeneral)
				return false
			}
		default:
			// unknown labels are ignored, for forward compatibility
		}
		return true
	})

	if parseErr != nil {
		return nil, parseErr
	}
	if !seen {
		return nil, fmt.Errorf("unexpected 'show profiling status' response: %q: %w", firstLine, native_errors.ErrGeneral)
	}
	return p, nil
}
