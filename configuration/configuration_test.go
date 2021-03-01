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

package configuration

import (
	"fmt"
	"os"
	"testing"
)

const testConf = `
# _version=1
global
	daemon
	nbproc 4
	maxconn 2000
	external-check
	stats socket /var/run/haproxy.sock level admin mode 0660
	lua-load /etc/foo.lua
	lua-load /etc/bar.lua

defaults
  maxconn 2000
  mode http
  bind-process 1-4
  balance roundrobin
  option clitcpka
  option dontlognull
  option forwardfor header X-Forwarded-For
  option http-use-htx
  option httpclose
  option httplog
  timeout queue 900
  timeout server 2s
  timeout check 2s
  timeout client 4s
  timeout connect 5s
  timeout http-request 2s
  timeout http-keep-alive 3s
  default-server fall 2s rise 4s inter 5s port 8888
  default_backend test
  option external-check
  external-check path /bin
  external-check command /bin/true
  errorfile 403 /test/403.html
  errorfile 500 /test/500.html
  errorfile 429 /test/429.html
  monitor-uri /monitor

frontend test
  mode http
  bind 192.168.1.1:80 name webserv
  bind 192.168.1.1:8080 name webserv2
  bind-process odd
  option httplog
  option dontlognull
  option contstats
  option log-separate-errors
  acl invalid_src  src          0.0.0.0/7 224.0.0.0/3
  acl invalid_src  src_port     0:1023
  acl local_dst    hdr(host) -i localhost
  monitor-uri /healthz
  monitor fail if site_dead
  filter trace name BEFORE-HTTP-COMP random-parsing hexdump
  filter compression
  filter trace name AFTER-HTTP-COMP random-forwarding
  http-request allow if src 192.168.0.0/16
  http-request set-header X-SSL %[ssl_fc]
  http-request set-var(req.my_var) req.fhdr(user-agent),lower
  http-request set-map(map.lst) %[src] %[req.hdr(X-Value)]
  http-request del-map(map.lst) %[src] if FALSE
  http-request cache-use cache-name if FALSE
  http-request disable-l7-retry if FALSE
  http-request early-hint hint-name %[src] if FALSE
  http-request replace-uri ^http://(.*) https://1 if FALSE
  http-request sc-inc-gpc0(0) if FALSE
  http-request sc-inc-gpc1(0) if FALSE
  http-request do-resolve(txn.myip,mydns,ipv4) hdr(Host),lower
  http-request sc-set-gpt0(1) hdr(Host),lower if FALSE
  http-request sc-set-gpt0(1) 20 if FALSE
  http-request set-mark 20 if FALSE
  http-request set-nice 20 if FALSE
  http-request set-method POST if FALSE
  http-request set-priority-class req.hdr(class) if FALSE
  http-request set-priority-offset req.hdr(offset) if FALSE
  http-request set-src req.hdr(src) if FALSE
  http-request set-src-port req.hdr(port) if FALSE
  http-request wait-for-handshake if FALSE
  http-request set-tos 0 if FALSE
  http-request silent-drop if FALSE
  http-request unset-var(req.my_var) if FALSE
  http-request strict-mode on if FALSE
  http-request lua.foo param1 param2 if FALSE
  http-request use-service svrs if FALSE
  http-request return status 200 content-type "text/plain" string "My content" hdr Some-Header value if FALSE 
  http-response allow if src 192.168.0.0/16
  http-response set-header X-SSL %[ssl_fc]
  http-response set-var(req.my_var) req.fhdr(user-agent),lower
  http-response set-map(map.lst) %[src] %[res.hdr(X-Value)]
  http-response del-map(map.lst) %[src] if FALSE
  http-response sc-inc-gpc0(0) if FALSE
  http-response sc-inc-gpc1(0) if FALSE
  http-response sc-set-gpt0(1) hdr(Host),lower if FALSE
  http-response sc-set-gpt0(1) 20 if FALSE
  http-response set-mark 20 if FALSE
  http-response set-nice 20 if FALSE
  http-response set-tos 0 if FALSE
  http-response silent-drop if FALSE
  http-response unset-var(req.my_var) if FALSE
  http-response track-sc0 src table tr0 if FALSE
  http-response track-sc1 src table tr1 if FALSE
  http-response track-sc2 src table tr2 if FALSE
  http-response strict-mode on if FALSE
  http-response lua.foo param1 param2 if FALSE
  tcp-request connection accept if TRUE
  tcp-request connection reject if FALSE
  tcp-request content accept if TRUE
  tcp-request content reject if FALSE
  tcp-request connection lua.foo param1 param2 if FALSE
  tcp-request content lua.foo param1 param2 if FALSE
  log global
  no log
  log 127.0.0.1:514 local0 notice notice
  log-tag bla
  option httpclose
  timeout http-request 2s
  timeout http-keep-alive 3s
  maxconn 2000
  default_backend test
  use_backend test_2 if TRUE
  use_backend %[req.cookie(foo)]
  timeout client 4s
  option clitcpka
  unique-id-format %{+X}o%ci:%cp_%fi:%fp_%Ts_%rt:%pid
  unique-id-header X-Unique-ID

frontend test_2
  mode http
  bind-process even
  option httplog
  option dontlognull
  option contstats
  option log-separate-errors
  log-tag bla
  option httpclose
  timeout http-request 2s
  timeout http-keep-alive 3s
  maxconn 2000
  default_backend test_2
  timeout client 4s
  option clitcpka
  http-request capture req.cook_cnt(FirstVisit),bool len 10
  http-request capture req.cook_cnt(FirstVisit),bool id 0
  http-response capture res.header id 0
  unique-id-format %{+X}o%ci:%cp_%fi:%fp_%Ts_%rt
  unique-id-header X-Unique-ID-test-2

backend test
  mode http
  balance roundrobin
  bind-process all
  hash-type consistent sdbm avalanche
  log-tag bla
  option http-keep-alive
  option forwardfor header X-Forwarded-For
  option httpchk HEAD /
  default-server fall 2s rise 4s inter 5s port 8888
  stick store-request src table test
  stick match src table test
  stick on src table test
  stick store-response src
  stick store-response src_port table test_port
  stick store-response src table test if TRUE
  tcp-response content accept if TRUE
  tcp-response content reject if FALSE
  tcp-response content lua.foo param1 param2 if FALSE
  option contstats
  timeout check 2s
  timeout tunnel 5s
  timeout server 3s
  cookie BLA rewrite httponly nocache
  option external-check
  external-check command /bin/false
  use-server webserv if TRUE
  use-server webserv2 unless TRUE
  server webserv 192.168.1.1:9200 maxconn 1000 ssl weight 10 inter 2s cookie BLAH slowstart 6000 proxy-v2-options authority,crc32c
  server webserv2 192.168.1.1:9300 maxconn 1000 ssl weight 10 inter 2s cookie BLAH slowstart 6000 proxy-v2-options authority,crc32c
  http-request set-dst hdr(x-dst)
  http-request set-dst-port int(4000)

peers mycluster
  peer hapee 192.168.1.1:1023
  peer aggregator HARDCODEDCLUSTERIP:10023

resolvers test
  nameserver dns1 10.0.0.1:53
  accepted_payload_size 8192
  resolve_retries       3
  timeout resolve       1s
  timeout retry         1s
  hold other           30s
  hold refused         30s
  hold nx              30s
  hold timeout         30s
  hold valid 5s

backend test_2
  mode http
  balance roundrobin
  bind-process all
  hash-type consistent sdbm avalanche
  log-tag bla
  option http-keep-alive
  option forwardfor header X-Forwarded-For
  option httpchk HEAD /
  default-server fall 2s rise 4s inter 5s port 8888 slowstart 6000
  option contstats
  timeout check 2s
  timeout tunnel 5s
  timeout server 3s
  cookie BLA rewrite httponly nocache
  stick-table type ip size 100k expire 1h peers mycluster store http_req_rate(10s)
  http-check expect rstatus some-pattern
`
const testPath = "/tmp/haproxy-test.cfg"

//nolint:gochecknoglobals
var (
	client  *Client
	version int64 = 1
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var err error

		if err = prepareTestFile(testConf, testPath); err != nil {
			fmt.Println("Could not prepare tests")
			return 1
		}

		defer func() { _ = deleteTestFile(testPath) }()
		client, err = prepareClient(testPath)
		if err != nil {
			fmt.Println("Could not prepare client:", err.Error())
			return 1
		}

		return m.Run()
	}())
}

func prepareTestFile(conf string, path string) error {
	// detect if file exists
	_, err := os.Stat(path)
	var file *os.File
	// create file if not exists
	if os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			return err
		}
	} else {
		// if exists delete it and create again
		err = deleteTestFile(path)
		if err != nil {
			return err
		}
		file, err = os.Create(path)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(conf)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}

func deleteTestFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func prepareClient(path string) (c *Client, err error) {
	c = &Client{}
	p := ClientParams{
		ConfigurationFile:      path,
		Haproxy:                "echo",
		UseValidation:          true,
		PersistentTransactions: true,
		TransactionDir:         "/tmp/haproxy-test",
	}
	err = c.Init(p)
	return
}
