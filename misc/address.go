// Copyright 2019 HAProxy Technologies
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

package misc

import (
	"fmt"
	"net"
	"strings"
)

func IsPrefixed(address string) bool {
	switch {
	case strings.HasPrefix(address, "ipv4@"):
		return true
	case strings.HasPrefix(address, "ipv6@"):
		return true
	case strings.HasPrefix(address, "udp@"):
		return true
	case strings.HasPrefix(address, "udp6"):
		return true
	case strings.HasPrefix(address, "unix@"):
		return true
	case strings.HasPrefix(address, "abns@"):
		return true
	case strings.HasPrefix(address, "fd@<n>"):
		return true
	case strings.HasPrefix(address, "sockpair@<n>"):
		return true
	case strings.HasPrefix(address, "quicv4@"):
		return true
	case strings.HasPrefix(address, "quicv6@"):
		return true
	}
	return false
}

func IsSocket(address string) bool {
	return strings.HasPrefix(address, "/")
}

func HasNoPath(address string) bool {
	return address == "*"
}

func ParseAddress(address string, portRequired bool) (string, error) {
	if IsIPv4(address) {

		colons := strings.Count(address, ":")
		hasPort := colons >= 1

		if portRequired && !hasPort {
			return "", fmt.Errorf("given address does not contain a Port number")
		}

		if !hasPort {
			ip := net.ParseIP(address)
			if ip == nil {
				return "", fmt.Errorf("invalid IP address given")
			}
			return ip.String(), nil
		}

		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return "", err
		}

		if hasPort {
			return fmt.Sprintf("%s:%s", host, port), nil
		}
		return host, nil
	}

	if strings.HasPrefix(address, "[") {
		if strings.HasSuffix(address, "]") {
			return address, nil
		}
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[%s]:%s", host, port), nil
	}

	ip := net.ParseIP(address)
	if ip == nil {
		return "", fmt.Errorf("invalid IP address given")
	}
	return address, nil
}

func IsIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

func IsIPv6(address string) bool {
	return !IsIPv4(address)
}
