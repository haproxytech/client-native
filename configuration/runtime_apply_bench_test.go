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

package configuration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/runtime"
)

func genConfig(backends, servers int) string {
	var sb strings.Builder
	sb.WriteString("global\n  maxconn 1000\n\n")
	sb.WriteString("defaults d\n  mode http\n  timeout client 30s\n  timeout server 30s\n\n")
	sb.WriteString("frontend fe from d\n  bind :80\n  default_backend be0\n\n")
	for i := range backends {
		fmt.Fprintf(&sb, "backend be%d from d\n  balance roundrobin\n  option httpchk\n", i)
		for j := range servers {
			fmt.Fprintf(&sb, "  server s%d 10.%d.%d.%d:8080 check weight 10\n", j, i/256%256, i%256, j%256)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

type benchSize struct {
	name              string
	backends, servers int
}

var benchSizes = []benchSize{{"100x100", 100, 100}, {"10000x10", 10000, 10}}

func BenchmarkParserString(b *testing.B) {
	for _, size := range benchSizes {
		p := newRuntimeApplyParser(b, genConfig(size.backends, size.servers))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				_ = p.String()
			}
		})
	}
}

// benchCases returns the configuration pairs the runtime steps benchmarks run.
func benchCases(b *testing.B, size benchSize) (string, []struct{ name, cfg string }) {
	b.Helper()
	cfg := genConfig(size.backends, size.servers)
	return cfg, []struct{ name, cfg string }{
		{"nochange", cfg},
		{"weight-last", changeLastBackendWeight(b, cfg, size.backends)},
		{"global", strings.Replace(cfg, "maxconn 1000", "maxconn 2000", 1)},
		{"frontend", strings.Replace(cfg, "bind :80", "bind :81", 1)},
	}
}

func BenchmarkRuntimeSteps(b *testing.B) {
	for _, size := range benchSizes {
		cfg, cases := benchCases(b, size)
		oldP := newRuntimeApplyParser(b, cfg)
		for _, c := range cases {
			newP := newRuntimeApplyParser(b, c.cfg)
			b.Run(size.name+"/"+c.name, func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for range b.N {
					runtimeSteps(registry, oldP, newP, nil)
				}
			})
		}
	}
}

func BenchmarkRuntimeStepsPrefiltered(b *testing.B) {
	for _, size := range benchSizes {
		cfg, cases := benchCases(b, size)
		oldP := newRuntimeApplyParser(b, cfg)
		oldText := oldP.String()
		for _, c := range cases {
			newP := newRuntimeApplyParser(b, c.cfg)
			newText := newP.String()
			b.Run(size.name+"/"+c.name, func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for range b.N {
					runtimeSteps(registry, oldP, newP, func() map[string]struct{} { return sameSections(oldText, newText) })
				}
			})
		}
	}
}

// noopRuntime accepts the weight change every benchmark commit makes.
type noopRuntime struct{ runtime.Runtime }

func (noopRuntime) SetServerWeight(string, string, string) error { return nil }

// BenchmarkRuntimeCommitRaw measures a whole PostRawConfiguration commit with
// and without runtime apply; every commit changes one server weight.
func BenchmarkRuntimeCommitRaw(b *testing.B) {
	for _, size := range benchSizes {
		cfg := genConfig(size.backends, size.servers)
		cfgs := [2]string{changeLastBackendWeight(b, cfg, size.backends), cfg}
		for _, mode := range []string{"disabled", "enabled"} {
			b.Run(size.name+"/"+mode, func(b *testing.B) {
				var opts []options.ConfigurationOption
				if mode == "enabled" {
					opts = append(opts, options.RuntimeApply(func() runtime.Runtime { return &noopRuntime{} }, func(res options.CommitResult) {
						if res.Reload {
							b.Fatalf("runtime apply fell back to a reload: %v", res.Err)
						}
					}))
				}
				c := newCommitClient(b, cfg, opts...)
				b.ReportAllocs()
				b.ResetTimer()
				for i := range b.N {
					version, err := c.GetVersion("")
					if err != nil {
						b.Fatal(err)
					}
					if err := c.PostRawConfiguration(&cfgs[i%2], version, false); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

// changeLastBackendWeight changes the weight of the first server of the
// backend whose name sorts last.
func changeLastBackendWeight(b *testing.B, cfg string, backends int) string {
	b.Helper()
	last := ""
	for i := range backends {
		last = max(last, fmt.Sprintf("be%d", i))
	}
	at := strings.Index(cfg, "backend "+last+" from d\n")
	if at < 0 {
		b.Fatalf("backend %s not found", last)
	}
	return cfg[:at] + strings.Replace(cfg[at:], "weight 10", "weight 20", 1)
}
