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

func SanitizeIPv6Address(addr string) string {
	switch {
	case strings.HasPrefix(addr, "ipv6@"),
		strings.HasPrefix(addr, "udp6@"),
		strings.HasPrefix(addr, "quic6@"):
		addrSlice := strings.SplitN(addr, "@", 2)
		if strings.HasPrefix(addrSlice[1], "[") {
			return addr
		}
		return fmt.Sprintf("%s@[%s]", addrSlice[0], addrSlice[1])
	case strings.Count(addr, ":") > 1:
		if strings.HasPrefix(addr, "[") {
			return addr
		}
		return fmt.Sprintf("[%s]", addr)
	default:
		return addr
	}
}

func ParseBindAddress(path string) (string, string, error) {
	switch {
	// environment variables, port can be part of it or not
	case strings.HasPrefix(path, "\"$"), strings.HasPrefix(path, "$"):
		if last := strings.LastIndex(path, ":"); last > 0 {
			return path[:last], path[last+1:], nil
		}
		return path, "", nil
	// unix socket, abstract namespace or file descriptor, no port available
	case strings.HasPrefix(path, "/"),
		strings.HasPrefix(path, "unix@"),
		strings.HasPrefix(path, "absn@"),
		strings.HasPrefix(path, "fd@"),
		strings.HasPrefix(path, "sockpair@"):

		return path, "", nil
	// ipv6 address and port is mandatory
	case strings.HasPrefix(path, "ipv6@"),
		strings.HasPrefix(path, "udp6@"),
		strings.HasPrefix(path, "quic6@"),
		strings.HasPrefix(path, "["),
		strings.Count(path, ":") > 1:

		pathSlice := strings.SplitN(path, "@", 2)
		prefix := ""
		address := ""
		if len(pathSlice) > 1 {
			prefix = fmt.Sprintf("%s@", pathSlice[0])
			address = pathSlice[1]
		} else {
			address = pathSlice[0]
		}
		if strings.HasPrefix(address, "[") {
			host, port, err := net.SplitHostPort(address)
			if err != nil {
				return "", "", err
			}
			return fmt.Sprintf("%s%s", prefix, host), port, nil
		}
		index := strings.LastIndex(address, ":")
		if index == -1 {
			return "", "", &net.AddrError{Err: "missing port in address", Addr: address}
		}
		port := address[strings.LastIndex(address, ":")+1:]
		host := fmt.Sprintf("[%s]", address[:strings.LastIndex(address, ":")])
		host, port, err := net.SplitHostPort(fmt.Sprintf("%s:%s", host, port))
		if err != nil {
			return "", "", err
		}
		return fmt.Sprintf("%s%s", prefix, host), port, nil

	// ipv4 address and port is mandatory
	case strings.HasPrefix(path, "ipv4@"),
		strings.HasPrefix(path, "udp4@"),
		strings.HasPrefix(path, "quicv4@"):

		pathSlice := strings.SplitN(path, "@", 2)
		prefix := ""
		address := ""
		if len(pathSlice) > 1 {
			prefix = fmt.Sprintf("%s@", pathSlice[0])
			address = pathSlice[1]
		} else {
			address = pathSlice[0]
		}
		// split host/port, validate ip address and return it
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return "", "", err
		}
		return fmt.Sprintf("%s%s", prefix, host), port, nil
	// hostname and port is mandatory
	default:
		// split host/port, validate ip address and return it
		host, port, err := net.SplitHostPort(path)
		if err != nil {
			return "", "", err
		}
		return host, port, nil
	}
}
