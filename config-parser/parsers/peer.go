/*
Copyright 2019 HAProxy Technologies

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package parsers

import (
	"fmt"
	"strconv"

	"github.com/haproxytech/client-native/v6/config-parser/common"
	"github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

type Peer struct {
	data        []types.Peer
	preComments []string // comments that appear before the actual line
}

func (l *Peer) parse(line string, parts []string, comment string) (*types.Peer, error) {
	if len(parts) > 2 {
		adr, p, found := common.CutRight(parts[2], ":")
		if found && len(adr) > 0 {
			if port, err := strconv.ParseInt(p, 10, 64); err == nil {
				peer := &types.Peer{
					Name:    parts[1],
					IP:      adr,
					Port:    port,
					Comment: comment,
				}
				if len(parts) > 4 && parts[3] == "shard" {
					peer.Shard = parts[4]
				}
				return peer, nil
			}
		}
	}
	return nil, &errors.ParseError{Parser: "PeerLines", Line: line}
}

func (l *Peer) Result() ([]common.ReturnResultLine, error) {
	if len(l.data) == 0 {
		return nil, errors.ErrFetch
	}
	result := make([]common.ReturnResultLine, len(l.data))
	for index, peer := range l.data {
		str := fmt.Sprintf("peer %s %s:%d", peer.Name, peer.IP, peer.Port)
		if peer.Shard != "" {
			str = fmt.Sprintf("%s shard %s", str, peer.Shard)
		}
		result[index] = common.ReturnResultLine{
			Data:    str,
			Comment: peer.Comment,
		}
	}
	return result, nil
}
