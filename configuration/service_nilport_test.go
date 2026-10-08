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
	"testing"

	"github.com/stretchr/testify/require"
)

// A discovered server without port (consul port 0) sharing the address of an
// active node with a port must not crash the reconciler.
func TestServiceNodesMatch_ServerWithoutPort(t *testing.T) {
	port := int64(80)
	s := &Service{
		name: "svc",
		nodes: []*ServiceNode{
			{address: "10.0.0.1", name: "SRV_1", port: &port},
			{address: "127.0.0.1", name: "SRV_2", port: &port, disabled: true},
			{address: "127.0.0.1", name: "SRV_3", port: &port, disabled: true},
		},
	}
	servers := []ServiceServer{{Address: "10.0.0.1", Port: &port}, {Address: "10.0.0.1"}}
	// the reconcile steps of Update before the config is written
	require.NotPanics(t, func() {
		s.markRemovedNodes(servers)
		require.True(t, s.serverExists(servers[0]))
		require.False(t, s.serverExists(servers[1]))
	})
}
