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
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/haproxytech/client-native/v3/configuration/options"
)

const testConf = `
# _version=1
global
	daemon
	nbproc 4
	maxconn 2000
	external-check
  ca-base /etc/ssl/certs
  crt-base /etc/ssl/private
	stats socket /var/run/haproxy.sock level admin mode 0660
  lua-prepend-path /usr/share/haproxy-lua/?/init.lua
  lua-prepend-path /usr/share/haproxy-lua/?.lua cpath
	lua-load /etc/foo.lua
	lua-load /etc/bar.lua
  h1-case-adjust-file /etc/headers.adjust
  h1-case-adjust host Host
  h1-case-adjust content-type Content-Type
  uid 1
  gid 1
  profiling.memory on
  ssl-mode-async
  tune.buffers.limit 11
  tune.buffers.reserve 12
  tune.bufsize 13
  tune.comp.maxlevel 14
  tune.fail-alloc
  tune.fd.edge-triggered on
  tune.h2.header-table-size 15
  tune.h2.initial-window-size 16
  tune.h2.max-concurrent-streams 17
  tune.h2.max-frame-size 18
  tune.http.cookielen 19
  tune.http.logurilen 20
  tune.http.maxhdr 21
  tune.idle-pool.shared on
  tune.idletimer 22
  tune.listener.multi-queue on
  tune.lua.forced-yield 23
  tune.lua.maxmem
  tune.lua.session-timeout 25
  tune.lua.task-timeout 26
  tune.lua.service-timeout 27
  tune.maxaccept 28
  tune.maxpollevents 29
  tune.maxrewrite 30
  tune.pattern.cache-size 31
  tune.pipesize 32
  tune.pool-high-fd-ratio 33
  tune.pool-low-fd-ratio 34
  tune.rcvbuf.client 35
  tune.rcvbuf.server 36
  tune.recv_enough 37
  tune.runqueue-depth 38
  tune.sched.low-latency on
  tune.sndbuf.client 39
  tune.sndbuf.server 40
  tune.ssl.cachesize 41
  tune.ssl.force-private-cache
  tune.ssl.keylog on
  tune.ssl.lifetime 43
  tune.ssl.maxrecord 44
  tune.ssl.default-dh-param 45
  tune.ssl.ssl-ctx-cache-size 46
  tune.ssl.capture-buffer-size 47
  tune.vars.global-max-size 49
  tune.vars.proc-max-size 50
  tune.vars.reqres-max-size 51
  tune.vars.sess-max-size 52
  tune.vars.txn-max-size 53
  tune.zlib.memlevel 54
  tune.zlib.windowsize 55
  busy-polling
  max-spread-checks 1
  maxconnrate 2
  maxcomprate 3
  maxcompcpuusage 4
  maxpipes 5
  maxsessrate 6
  maxsslconn 7
  maxsslrate 8
  maxzlibmem 9
  noepoll
  nokqueue
  noevports
  nopoll
  nosplice
  nogetaddrinfo
  noreuseport
  profiling.tasks on
  spread-checks 10
  wurfl-data-file path
  wurfl-information-list wurfl_id,wurfl_root_id,wurfl_isdevroot,wurfl_useragent,wurfl_api_version,wurfl_info,wurfl_last_load_time,wurfl_normalized_useragent
  wurfl-information-list-separator ,
  wurfl-patch-file path1,path2
  wurfl-cache-size 64
  ssl-default-bind-curves X25519:P-256
  ssl-skip-self-issued-ca
  node node
  description description
  expose-experimental-directives
  insecure-fork-wanted
  insecure-setuid-wanted
  issuers-chain-path issuers-chain-path
  h2-workaround-bogus-websocket-clients
  lua-load-per-thread file.ext
  mworker-max-reloads 5
  numa-cpu-mapping
  pp2-never-send-local
  ulimit-n 10
  set-dumpable
  strict-limits
  grace 10s
  chroot /var/www
  ssl-default-server-ciphers ECDH+AESGCM:ECDH+CHACHA20:ECDH+AES256:ECDH+AES128:!aNULL:!SHA1:!AESCCM
  ssl-default-server-ciphersuites TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256
  hard-stop-after 2s
  localpeer test
  user thomas
  group anderson
  nbthread 128
  pidfile pidfile.text
  ssl-default-bind-ciphers ECDH+AESGCM:ECDH+CHACHA20
  ssl-default-bind-ciphersuites TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384
  ssl-default-server-options ssl-min-ver TLSv1.1 no-tls-tickets
  thread-groups 1
  thread-group first 1-16
  stats maxconn 20
  ssl-load-extra-files bundle
  deviceatlas-json-file atlas.json
  deviceatlas-log-level 1
  deviceatlas-separator -
  deviceatlas-properties-cookie chocolate
  51degrees-data-file 51.file
  51degrees-property-name-list first second third fourth fifth
  51degrees-property-separator /
  51degrees-cache-size 51
  quiet
  zero-warning
  ssl-engine first
  ssl-engine second RSA,DSA,DH,EC,RAND
  ssl-engine third CIPHERS,DIGESTS,PKEY,PKEY_CRYPTO,PKEY_ASN1
  ssl-dh-param-file ssl-dh-param-file.txt
  ssl-server-verify required
  set-var proc.current_state str(primary)
  set-var proc.prio int(100)
  set-var proc.threshold int(200),sub(proc.prio)
  set-var-fmt proc.bootid "%pid|%t"
  set-var-fmt proc.current_state "primary"

defaults
  maxconn 2000
  backlog 1024
  mode http
  bind-process 1-4
  balance roundrobin
  option tcpka
  option srvtcpka
  option clitcpka
  option dontlognull
  option forwardfor header X-Forwarded-For
  option http-use-htx
  option httpclose
  option httplog
  option disable-h2-upgrade
  option log-health-checks
  no option checkcache
  no option http-ignore-probes
  no option http-use-proxy-header
  no option httpslog
  no option independent-streams
  no option nolinger
  no option originalto
  option persist
  option prefer-last-server
  option socket-stats
  option tcp-smart-accept
  option tcp-smart-connect
  option transparent
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
  http-check send-state
  http-check disable-on-404
  option accept-invalid-http-request
  option accept-invalid-http-response
  option h1-case-adjust-bogus-client
  option h1-case-adjust-bogus-server
  compression offload
  clitcpka-cnt 10
  clitcpka-idle 10s
  clitcpka-intvl 10

frontend test
  mode http
  backlog 2048
  bind 192.168.1.1:80 name webserv
  bind 192.168.1.1:8080 name webserv2
  bind-process odd
  option httplog
  option dontlognull
  option contstats
  option log-separate-errors
  option http-ignore-probes
  option http-use-proxy-header
  option httpslog
  option independent-streams
  option nolinger
  option originalto
  option socket-stats
  option tcp-smart-accept
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
  http-request redirect scheme https if !{ ssl_fc }
  http-request redirect location https://%[hdr(host),field(1,:)]:443%[capture.req.uri] code 302
  http-request deny unless src 192.168.0.0/16
  http-request deny deny_status 400 content-type application/json lf-file /var/errors.file
  http-request wait-for-body time 20s at-least 100k
  http-request set-timeout server 20
  http-request set-timeout tunnel 20
  http-response allow if src 192.168.0.0/16
  http-response set-header X-SSL %[ssl_fc]
  http-response set-var(req.my_var) req.fhdr(user-agent),lower
  http-response set-map(map.lst) %[src] %[res.hdr(X-Value)]
  http-response del-map(map.lst) %[src] if FALSE
  http-response cache-store cache-name if FALSE
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
  http-response deny deny_status 400 content-type application/json lf-file /var/errors.file
  http-response wait-for-body time 20s at-least 100k
  http-after-response set-header Strict-Transport-Security "max-age=31536000"
  http-after-response replace-header Set-Cookie (C=[^;]*);(.*) \1;ip=%bi;\2
  http-after-response replace-value Cache-control ^public$ private
  http-after-response set-status 503 reason "SlowDown"
  http-after-response set-var(sess.last_redir) res.hdr(location)
  http-after-response unset-var(sess.last_redir)
  tcp-request connection accept if TRUE
  tcp-request connection reject if FALSE
  tcp-request content accept if TRUE
  tcp-request content reject if FALSE
  tcp-request connection silent-drop
  tcp-request connection silent-drop if TRUE
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
  option tcpka
  option clitcpka
  unique-id-format %{+X}o%ci:%cp_%fi:%fp_%Ts_%rt:%pid
  unique-id-header X-Unique-ID
  no option accept-invalid-http-request
  no option h1-case-adjust-bogus-client
  compression algo identity gzip
  compression type text/plain
  compression offload
  clitcpka-cnt 10
  clitcpka-idle 10s
  clitcpka-intvl 10

frontend test_2
  mode http
  bind-process even
  option httplog
  option dontlognull
  option contstats
  option log-separate-errors
  no option http-ignore-probes
  no option http-use-proxy-header
  no option httpslog
  no option independent-streams
  no option nolinger
  no option originalto
  no option socket-stats
  no option tcp-smart-accept
  log-tag bla
  option httpclose
  timeout http-request 2s
  timeout http-keep-alive 3s
  maxconn 2000
  backlog 2048
  default_backend test_2
  timeout client 4s
  option tcpka
  option clitcpka
  http-request capture req.cook_cnt(FirstVisit),bool len 10
  http-request capture req.cook_cnt(FirstVisit),bool id 0
  http-response capture res.header id 0
  unique-id-format %{+X}o%ci:%cp_%fi:%fp_%Ts_%rt
  unique-id-header X-Unique-ID-test-2
  clitcpka-cnt 10
  clitcpka-idle 10s
  clitcpka-intvl 10

backend test
  mode http
  balance roundrobin
  bind-process all
  hash-type consistent sdbm avalanche
  log-tag bla
  option http-keep-alive
  option forwardfor header X-Forwarded-For
  option httpchk HEAD /
  option tcpka
  option srvtcpka
  option checkcache
  option independent-streams
  option nolinger
  option originalto
  option persist
  option prefer-last-server
  option spop-check
  option tcp-smart-connect
  option transparent
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
  http-check connect
  http-check send meth GET uri / ver HTTP/1.1 hdr host haproxy.1wt.eu
  http-check expect status 200-399
  http-check connect port 443 ssl sni haproxy.1wt.eu
  http-check expect status 200,201,300-310
  http-check expect header name "set-cookie" value -m beg "sessid="
  http-check expect ! string SQL\ Error
  http-check expect ! rstatus ^5
  http-check expect rstring <!--tag:[0-9a-f]*--></html>
  http-check unset-var(check.port)
  http-check set-var(check.port) int(1234)
  http-check set-var-fmt(check.port) int(1234)
  http-check send-state
  http-check disable-on-404
  server-template srv 1-3 google.com:80 check
  server-template site 1-10 google.com:8080 check backup
  server-template website 10-100 google.com:443 check no-backup
  server-template test 5 test.com check backup
  no option accept-invalid-http-response
  no option h1-case-adjust-bogus-server
  compression type application/json text/plain

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

cache test
  total-max-size 1024
  max-object-size 8
  max-age 60
  process-vary on
  max-secondary-entries 10

ring myring
  description "My local buffer"
  format rfc3164
  maxlen 1200
  size 32764
  timeout connect 5s
  timeout server 10s
  server mysyslogsrv 127.0.0.1:6514 log-proto octet-count

backend test_2
  mode http
  balance roundrobin
  bind-process all
  hash-type consistent sdbm avalanche
  log-tag bla
  option http-keep-alive
  option forwardfor header X-Forwarded-For
  option httpchk HEAD /
  option tcpka
  option srvtcpka
  no option checkcache
  no option independent-streams
  no option nolinger
  no option originalto
  no option persist
  no option prefer-last-server
  no option spop-check
  no option tcp-smart-connect
  no option transparent
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
	clientTest Configuration
	version    int64 = 1
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		var err error

		if err = prepareTestFile(testConf, testPath); err != nil {
			fmt.Println("Could not prepare tests")
			return 1
		}

		defer func() { _ = deleteTestFile(testPath) }()
		clientTest, err = prepareClient(testPath)
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

	file, err = os.OpenFile(path, os.O_RDWR, 0o644)
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

func prepareClient(path string) (c Configuration, err error) {
	c, err = New(context.Background(),
		options.ConfigurationFile(path),
		options.HAProxyBin("echo"),
		options.UseModelsValidation,
		options.UsePersistentTransactions,
		options.TransactionsDir("/tmp/haproxy-test"),
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}
