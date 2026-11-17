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
	"errors"
	"fmt"
	"net"
	"slices"
	"strconv"
	"strings"
	"time"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/params"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/runtime"
)

// registry lists everything that can be applied through the runtime API
// instead of a reload. To make another attribute or object runtime-capable,
// add an entry here; the engine in runtime_apply.go needs no change. Any
// difference without an entry, and any applier returning false, makes the
// commit fall back to a reload.
//
// A plain attribute (keyword of a section, such as the frontend maxconn or the
// global process limits) gets an apply function that turns the old and new
// parser data into runtime calls.
//
// An object list (named lines such as servers) gets an object spec: items
// lists the lines of a section with their effective options, fields maps an
// option name to the applier run when that option changes on a kept line,
// add and del handle added and removed lines (nil = reload).
var registry = runtimeRegistry{ //nolint:gochecknoglobals
	parser.Global: {
		// maxconn 0 is automatic in the file but restores the startup value at runtime.
		"maxconn":     {apply: globalNumber(1, func(rt runtime.Runtime, v int64) error { return rt.SetGlobalMaxConn(int(v)) })},
		"maxconnrate": {apply: globalNumber(0, func(rt runtime.Runtime, v int64) error { return rt.SetRateLimitConnectionsGlobal(uint64(v)) })},
		"maxsessrate": {apply: globalNumber(0, func(rt runtime.Runtime, v int64) error { return rt.SetRateLimitSessionsGlobal(uint64(v)) })},
		"maxsslrate":  {apply: globalNumber(0, func(rt runtime.Runtime, v int64) error { return rt.SetRateLimitSSLSessionGlobal(uint64(v)) })},
		"maxcomprate": {apply: globalNumber(0, func(rt runtime.Runtime, v int64) error { return rt.SetRateLimitHTTPCompressionGlobal(uint64(v)) })},
	},
	parser.Frontends: {
		"maxconn": {apply: frontendMaxconn},
	},
	parser.Backends: {
		"server": {object: &objectSpec{
			items: serverItems,
			fields: map[string]itemApplier{
				"address":     serverAddress,
				"weight":      serverWeight,
				"disabled":    serverState,
				"enabled":     serverEnabled,
				"port":        serverCheckPort,
				"agent-check": serverAgentCheck,
				"agent-addr":  serverAgentAddr,
				"agent-send":  serverAgentSend,
				"agent-port":  serverAgentPort,
				"maxconn":     serverMaxconn,
				"addr":        serverCheckAddr,
			},
			add: addServer,
			del: delServer,
		}},
	},
}

// serverRemovableTimeout is how long a removed server may take to drop its
// connections before the commit falls back to a reload.
const serverRemovableTimeout = 500 * time.Millisecond

// frontendMaxconn needs a reload for maxconn 0: it means the global limit in
// the file but stops the frontend from accepting at runtime.
func frontendMaxconn(c change) ([]step, bool) {
	v, ok := c.new.(*types.Int64C)
	if !ok || v == nil || v.Value <= 0 {
		return nil, false
	}
	return []step{func(rt runtime.Runtime) error { return rt.SetFrontendMaxConn(c.section, int(v.Value)) }}, true
}

// globalNumber applies a global numeric keyword with set; a removed keyword
// or a value below lowest means reload.
func globalNumber(lowest int64, set func(rt runtime.Runtime, v int64) error) fieldApplier {
	return func(c change) ([]step, bool) {
		v, ok := c.new.(*types.Int64C)
		if !ok || v == nil || v.Value < lowest {
			return nil, false
		}
		n := v.Value
		return []step{func(rt runtime.Runtime) error { return set(rt, n) }}, true
	}
}

// serverItems returns the server lines of a section with their effective
// options: the default-server lines of the defaults chain (root first) and of
// the section, then the address and the options of the server line, each
// option keeping only its last occurrence.
func serverItems(p parser.Parser, t parser.Section, section string) ([]item, error) {
	data, err := p.Get(t, section, "server")
	if errors.Is(err, parser_errors.ErrFetch) {
		return []item{}, nil
	}
	if err != nil {
		return nil, err
	}
	servers, ok := data.([]types.Server)
	if !ok {
		return nil, fmt.Errorf("unexpected server data %T", data)
	}
	var inherited []option
	for _, d := range defaultsChain(p, t, section) {
		if inherited, err = addDefaultServers(p, parser.Defaults, d, inherited); err != nil {
			return nil, err
		}
	}
	if inherited, err = addDefaultServers(p, t, section, inherited); err != nil {
		return nil, err
	}
	items := make([]item, 0, len(servers))
	for _, srv := range servers {
		opts := addOption(slices.Clone(inherited), option{"address", srv.Address})
		opts = addParams(opts, srv.Params)
		items = append(items, item{name: srv.Name, options: opts})
	}
	return items, nil
}

// defaultsChain returns the defaults sections a section inherits from, root
// first. It stops at a cycle or a missing section.
func defaultsChain(p parser.Parser, t parser.Section, section string) []string {
	var chain []string
	name, err := p.SectionsDefaultsFromGet(t, section)
	for err == nil && name != "" && !slices.Contains(chain, name) {
		var next string
		if next, err = p.SectionsDefaultsFromGet(parser.Defaults, name); err == nil {
			chain = append(chain, name)
			name = next
		}
	}
	slices.Reverse(chain)
	return chain
}

// addDefaultServers adds the options of every default-server line of a section.
func addDefaultServers(p parser.Parser, t parser.Section, section string, opts []option) ([]option, error) {
	data, err := p.Get(t, section, "default-server")
	if errors.Is(err, parser_errors.ErrFetch) {
		return opts, nil
	}
	if err != nil {
		return nil, err
	}
	lines, ok := data.([]types.DefaultServer)
	if !ok {
		return nil, fmt.Errorf("unexpected default-server data %T", data)
	}
	for _, line := range lines {
		opts = addParams(opts, line.Params)
	}
	return opts, nil
}

func addParams(opts []option, ps []params.ServerOption) []option {
	for _, param := range ps {
		if s := param.String(); s != "" {
			name, args, _ := strings.Cut(s, " ")
			opts = addOption(opts, option{name, args})
		}
	}
	return opts
}

// addOption appends o after removing an earlier option of the same name, so
// the last writer is last and opposite keywords keep their relative order.
func addOption(opts []option, o option) []option {
	return append(slices.DeleteFunc(opts, func(x option) bool { return x.name == o.name }), o)
}

// plain reports whether s is a single token that the runtime socket reads
// like the configuration file does; quoted or escaped values need a reload.
func plain(s string) bool {
	return s != "" && !strings.ContainsAny(s, " \t\"'\\")
}

// lastOf returns whichever of the keywords a and b comes last in it, "" if none.
func lastOf(it item, a, b string) string {
	for _, o := range slices.Backward(it.options) {
		if o.name == a || o.name == b {
			return o.name
		}
	}
	return ""
}

// hostPort parses an IP literal with an optional port: 10.0.0.1, 10.0.0.1:80,
// [::1]:80. Port 0 means no port. Host names, sockets, prefixes, port offsets
// and ranges, and IPv6 without brackets and port are rejected.
func hostPort(addr string) (string, int, bool) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		ip := net.ParseIP(addr)
		if ip == nil || ip.To4() == nil || strings.Contains(addr, ":") {
			return "", 0, false
		}
		return addr, 0, true
	}
	if net.ParseIP(host) == nil || portStr == "" || strings.Trim(portStr, "0123456789") != "" {
		return "", 0, false
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port > 65535 {
		return "", 0, false
	}
	return host, port, true
}

func serverWeight(section string, _, cur item) ([]step, bool) {
	weight, ok := cur.get("weight")
	if !ok || !plain(weight) {
		return nil, false
	}
	return []step{func(rt runtime.Runtime) error { return rt.SetServerWeight(section, cur.name, weight) }}, true
}

func serverAddress(section string, old, cur item) ([]step, bool) {
	curAddr, _ := cur.get("address")
	ip, port, ok := hostPort(curAddr)
	if !ok {
		return nil, false
	}
	// An old host name may still be resolved and overwrite the new address.
	oldAddr, _ := old.get("address")
	if _, oldPort, oldOK := hostPort(oldAddr); !oldOK || (port == 0 && oldPort != 0) {
		return nil, false
	}
	return []step{func(rt runtime.Runtime) error { return rt.SetServerAddr(section, cur.name, ip, port) }}, true
}

func serverState(section string, _, cur item) ([]step, bool) {
	state := "ready"
	if lastOf(cur, "enabled", "disabled") == "disabled" {
		state = "maint"
	}
	return []step{func(rt runtime.Runtime) error { return rt.SetServerState(section, cur.name, state) }}, true
}

// serverEnabled leaves the state to serverState when disabled changed too.
func serverEnabled(section string, old, cur item) ([]step, bool) {
	_, oldDisabled := old.get("disabled")
	_, curDisabled := cur.get("disabled")
	if oldDisabled != curDisabled {
		return nil, true
	}
	return serverState(section, old, cur)
}

func serverCheckPort(section string, _, cur item) ([]step, bool) {
	v, ok := cur.get("port")
	if !ok || !plain(v) {
		return nil, false
	}
	port, err := strconv.Atoi(v)
	if err != nil {
		return nil, false
	}
	return []step{func(rt runtime.Runtime) error { return rt.SetServerCheckPort(section, cur.name, port) }}, true
}

func serverAgentCheck(section string, _, cur item) ([]step, bool) {
	if lastOf(cur, "agent-check", "no-agent-check") == "agent-check" {
		return []step{func(rt runtime.Runtime) error { return rt.EnableAgentCheck(section, cur.name) }}, true
	}
	return []step{func(rt runtime.Runtime) error { return rt.DisableAgentCheck(section, cur.name) }}, true
}

func serverAgentAddr(section string, _, cur item) ([]step, bool) {
	addr, ok := cur.get("agent-addr")
	if !ok || !plain(addr) {
		return nil, false
	}
	return []step{func(rt runtime.Runtime) error { return rt.SetServerAgentAddr(section, cur.name, addr) }}, true
}

func serverAgentSend(section string, _, cur item) ([]step, bool) {
	send, ok := cur.get("agent-send")
	if !ok || !plain(send) {
		return nil, false
	}
	return []step{func(rt runtime.Runtime) error { return rt.SetServerAgentSend(section, cur.name, send) }}, true
}

func serverAgentPort(section string, _, cur item) ([]step, bool) {
	v, ok := cur.get("agent-port")
	if !ok || !plain(v) {
		return nil, false
	}
	port, err := strconv.Atoi(v)
	if err != nil {
		return nil, false
	}
	return []step{func(rt runtime.Runtime) error { return rt.SetServerAgentPort(section, cur.name, port) }}, true
}

func serverMaxconn(section string, _, cur item) ([]step, bool) {
	v, ok := cur.get("maxconn")
	if !ok || !plain(v) {
		return nil, false
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return nil, false
	}
	return []step{func(rt runtime.Runtime) error { return rt.SetServerMaxconn(section, cur.name, n) }}, true
}

// serverCheckAddr sets the health check address (addr keyword); the check
// port has its own entry.
func serverCheckAddr(section string, _, cur item) ([]step, bool) {
	addr, ok := cur.get("addr")
	if !ok || !plain(addr) || net.ParseIP(addr) == nil {
		return nil, false
	}
	return []step{func(rt runtime.Runtime) error { return rt.SetServerCheckAddr(section, cur.name, addr, 0) }}, true
}

// dynamicServerDefaults maps keywords that only restate what a runtime-added
// server starts with, and that add server refuses, to the opposite keyword
// each one cancels when it comes later. Any other no-... keyword is sent as
// written; HAProxy rejects it and the commit falls back to a reload.
var dynamicServerDefaults = map[string]string{ //nolint:gochecknoglobals
	"enabled":        "disabled",
	"no-check":       "check",
	"no-agent-check": "agent-check",
}

// addServer adds a server with all its effective options, since a server
// added at runtime does not inherit default-server settings. Keywords of
// dynamicServerDefaults are not sent and drop their earlier opposite. The
// server starts in maintenance with checks off, so those are turned on
// afterwards.
func addServer(section string, it item) ([]step, bool) {
	addr, _ := it.get("address")
	if _, _, ok := hostPort(addr); !ok {
		return nil, false
	}
	var sent []option
	for _, o := range it.options {
		if o.name == "address" {
			continue
		}
		if strings.ContainsAny(o.name+o.args, "\"'\\") {
			return nil, false
		}
		if opposite, found := dynamicServerDefaults[o.name]; found {
			sent = slices.DeleteFunc(sent, func(x option) bool { return x.name == opposite })
			continue
		}
		sent = append(sent, o)
	}
	attrs := []string{addr}
	for _, o := range sent {
		if o.args == "" {
			attrs = append(attrs, o.name)
		} else {
			attrs = append(attrs, o.name+" "+o.args)
		}
	}
	attributes := strings.Join(attrs, " ")
	steps := []step{func(rt runtime.Runtime) error { return rt.AddServer(section, it.name, attributes) }}
	if lastOf(it, "check", "no-check") == "check" {
		steps = append(steps, func(rt runtime.Runtime) error { return rt.EnableServerHealth(section, it.name) })
	}
	if lastOf(it, "agent-check", "no-agent-check") == "agent-check" {
		steps = append(steps, func(rt runtime.Runtime) error { return rt.EnableAgentCheck(section, it.name) })
	}
	if lastOf(it, "enabled", "disabled") != "disabled" {
		steps = append(steps, func(rt runtime.Runtime) error { return rt.SetServerState(section, it.name, "ready") })
	}
	return steps, true
}

// delServer drains and deletes a server; del server needs it in maintenance
// with no connection left.
func delServer(section string, it item) ([]step, bool) {
	var steps []step
	if lastOf(it, "agent-check", "no-agent-check") == "agent-check" {
		steps = append(steps, func(rt runtime.Runtime) error { return rt.DisableAgentCheck(section, it.name) })
	}
	return append(steps,
		func(rt runtime.Runtime) error { return rt.SetServerState(section, it.name, "maint") },
		func(rt runtime.Runtime) error { return rt.ShutdownSessionsServer(section, it.name) },
		func(rt runtime.Runtime) error {
			return rt.WaitServerRemovable(section, it.name, serverRemovableTimeout)
		},
		func(rt runtime.Runtime) error { return rt.DeleteServer(section, it.name) },
	), true
}
