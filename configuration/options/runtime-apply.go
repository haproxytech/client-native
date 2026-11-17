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

package options

import "github.com/haproxytech/client-native/v6/runtime"

// CommitResult tells the caller what a committed transaction still requires.
type CommitResult struct {
	// TransactionID is the committed transaction.
	TransactionID string
	// Version is the configuration version the commit was based on.
	Version int64
	// Reload is true when HAProxy must be reloaded to match the configuration file.
	Reload bool
	// Err is why applying through the runtime API was abandoned, if it was.
	// It is diagnostic only: the configuration file was written either way.
	Err error
}

type runtimeApply struct {
	rt       func() runtime.Runtime
	onCommit func(CommitResult)
}

func (r runtimeApply) Set(p *ConfigurationOptions) error {
	p.RuntimeApply = r.rt
	p.OnCommit = r.onCommit
	return nil
}

// RuntimeApply makes every committed transaction try to apply its changes
// through the runtime API instead of requiring a reload. rt returns the runtime
// to use, or nil when the running HAProxy cannot be assumed to match the
// configuration file as it was before the commit; the commit is then reported
// as needing a reload without being compared. onCommit receives the outcome of
// every successful commit.
func RuntimeApply(rt func() runtime.Runtime, onCommit func(CommitResult)) ConfigurationOption {
	return runtimeApply{rt: rt, onCommit: onCommit}
}
