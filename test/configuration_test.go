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

package test

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"testing"

	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/configuration/options"
)

const testConf = `
# _version=1
global
	daemon
	default-path origin /some/path
	maxconn 2000
	external-check
  ca-base /etc/ssl/certs
  crt-base /etc/ssl/private
  cluster-secret my_secret
	stats socket /var/run/haproxy.sock level admin mode 0660
  lua-prepend-path /usr/share/haproxy-lua/?/init.lua
  lua-prepend-path /usr/share/haproxy-lua/?.lua cpath
	lua-load /etc/foo.lua
	lua-load /etc/bar.lua
  h1-case-adjust-file /etc/headers.adjust
  h1-case-adjust host Host
  h1-case-adjust content-type Content-Type
  limited-quic
  uid 1
  gid 1
  profiling.memory on
  ssl-mode-async
  tune.buffers.limit 11
  tune.buffers.reserve 12
  tune.bufsize 13
  tune.comp.maxlevel 14
  tune.disable-zero-copy-forwarding
  tune.events.max-events-at-once 10
  tune.fail-alloc
  tune.fd.edge-triggered on
  tune.h1.zero-copy-fwd-recv on
  tune.h1.zero-copy-fwd-send on
  tune.h2.header-table-size 15
  tune.h2.initial-window-size 16
  tune.h2.max-concurrent-streams 17
  tune.h2.max-frame-size 18
  tune.h2.zero-copy-fwd-send on
  tune.http.cookielen 19
  tune.http.logurilen 20
  tune.http.maxhdr 21
  tune.idle-pool.shared on
  tune.idletimer 22
  tune.listener.default-shards by-process
  tune.listener.multi-queue on
  tune.lua.forced-yield 23
  tune.lua.log.loggers on
  tune.lua.log.stderr auto
  tune.lua.maxmem
  tune.lua.session-timeout 25
  tune.lua.task-timeout 26
  tune.lua.service-timeout 27
  tune.max-checks-per-thread 0
  tune.maxaccept 28
  tune.maxpollevents 29
  tune.maxrewrite 30
  tune.pattern.cache-size 31
  tune.pipesize 32
  tune.pool-high-fd-ratio 33
  tune.pool-low-fd-ratio 34
  tune.pt.zero-copy-forwarding on
  tune.rcvbuf.backend 1024
  tune.rcvbuf.client 35
  tune.rcvbuf.frontend 2048
  tune.rcvbuf.server 36
  tune.recv_enough 37
  tune.runqueue-depth 38
  tune.sched.low-latency on
  tune.sndbuf.backend 1024
  tune.sndbuf.client 39
  tune.sndbuf.frontend 2048
  tune.sndbuf.server 40
  tune.ssl.cachesize 41
  tune.ssl.force-private-cache
  tune.ssl.keylog on
  tune.ssl.lifetime 43
  tune.ssl.maxrecord 44
  tune.ssl.default-dh-param 45
  tune.ssl.ssl-ctx-cache-size 46
  tune.ssl.capture-buffer-size 47
  tune.ssl.ocsp-update.maxdelay 48
  tune.ssl.ocsp-update.mindelay 49
  tune.stick-counters 50
  tune.vars.global-max-size 51
  tune.vars.proc-max-size 52
  tune.vars.reqres-max-size 53
  tune.vars.sess-max-size 54
  tune.vars.txn-max-size 55
  tune.quic.frontend.conn-tx-buffers.limit 10
  tune.quic.frontend.max-idle-timeout 10000
  tune.quic.frontend.max-streams-bidi 100
  tune.quic.max-frame-loss 5
  tune.quic.retry-threshold 5
  tune.quic.socket-owner connection
  tune.zlib.memlevel 54
  tune.zlib.windowsize 55
  tune.memory.hot-size 56
  busy-polling
  max-spread-checks 1ms
  close-spread-time 1s
  maxconnrate 2
  maxcomprate 3
  maxcompcpuusage 4
  maxpipes 5
  maxsessrate 6
  maxsslconn 7
  maxsslrate 8
  maxzlibmem 9
  no-quic
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
  ssl-default-server-curves brainpoolP384r1,brainpoolP512r1
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
  httpclient.resolvers.disabled on
  httpclient.resolvers.prefer ipv4
  httpclient.resolvers.id resolver_1
  httpclient.retries 3
  httpclient.ssl.ca-file my_test_file.ca
  httpclient.ssl.verify none
  httpclient.timeout.connect 2s
  prealloc-fd
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
  server-state-base /path
  server-state-file serverstatefile
  presetenv first order
  setenv third sister
  resetenv first second
  unsetenv third fourth
  anonkey 25
  tune.peers.max-updates-at-once 200
  tune.h2.be.initial-window-size 201
  tune.h2.be.max-concurrent-streams 202
  tune.h2.fe.initial-window-size 203
  tune.h2.fe.max-concurrent-streams 204
  tune.lua.burst-timeout 205
  ssl-default-bind-sigalgs RSA+SHA256
  ssl-default-bind-client-sigalgs ECDSA+SHA256:RSA+SHA256
  ssl-default-server-sigalgs RSA+SHA256
  ssl-default-server-client-sigalgs ECDSA+SHA256:RSA+SHA256
  ssl-propquery provider
  ssl-provider default
  ssl-provider-path test
  setcap cap_net_raw,cap_net_bind_service
  module-path /tmp/modules/path
  module-load modsecurity.so
  module-load test.so
  waf-load /tmp/rules.conf
  waf-body-limit 15000
  waf-json-levels 8
  modsecurity-deny-blocking-io
  maxmind-update url CITY http://192.168.122.1/GeoIP2-City.mmdb url ISP http://192.168.122.1/GeoIP2-ISP.mmdb delay 24h checksum hash log
  maxmind-load mlock_max 512000000 CITY /etc/hapee-2.5/GeoIP2-City.mmdb ISP /etc/hapee-2.5/GeoIP2-ISP.mmdb
  maxmind-cache-size 200000
  maxmind-debug
  fingerprint_ssl_bufsize 56
  log 127.0.0.1:10001 sample 1:4 local0
  log 127.0.0.1:10002 sample 2:4 local0

defaults test_defaults
  maxconn 2000
  backlog 1024
  mode http
  balance roundrobin
  hash-balance-factor 150

defaults test_defaults_2 from test_defaults
  option srvtcpka
  option clitcpka

defaults
  maxconn 2000
  backlog 1024
  mode http
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
  option originalto
  option persist
  option prefer-last-server
  option socket-stats
  option tcp-smart-accept
  option tcp-smart-connect
  option transparent
  option dontlog-normal
  option http-no-delay
  option splice-auto
  option splice-request
  option splice-response
  option idle-close-on-response
  option http-restrict-req-hdr-names reject
  timeout queue 900
  timeout server 2s
  timeout check 2s
  timeout client 4s
  timeout connect 5s
  timeout http-request 2s
  timeout http-keep-alive 3s
  timeout server-fin 1000
  timeout client-fin 1000
  timeout tarpit 2000
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
  http-error status 503 content-type "application/json" file /test/503
  http-error status 429 content-type application/json file /test/429
  option accept-invalid-http-request
  option accept-invalid-http-response
  option h1-case-adjust-bogus-client
  option h1-case-adjust-bogus-server
  compression offload
  clitcpka-cnt 10
  clitcpka-idle 10s
  clitcpka-intvl 10
  srvtcpka-cnt 10
  srvtcpka-idle 10s
  srvtcpka-intvl 10
  stats auth admin:AdMiN123
  stats auth admin2:AdMiN1234
  stats show-modules
  stats realm HAProxy\\ Statistics
  email-alert from srv01@example.com
  email-alert to support@example.com
  email-alert level err
  email-alert myhostname srv01
  email-alert mailers localmailer1
  load-server-state-from-file global
  fullconn 10
  max-keep-alive-queue 100
  retry-on 503 504
  http-send-name-header
  persist rdp-cookie
  source 192.168.1.200:80 usesrc 192.168.1.201:443
	tcp-check send GET\ /\ HTTP/2.0\r\n

frontend test
  mode http
  backlog 2048
  bind 192.168.1.1:80 name webserv thread all sigalgs RSA+SHA256 client-sigalgs ECDSA+SHA256:RSA+SHA256 ca-verify-file ca.pem nice 789
  bind 192.168.1.1:8080 name webserv2 thread 1/all
  bind 192.168.1.2:8080 name webserv3 thread 1/1
  bind [2a01:c9c0:a3:8::3]:80 name ipv6 thread 1/1-1
  bind 192.168.1.1:80 name test-quic quic-socket connection thread 1/1
  bind 192.168.1.1:80 name testnbcon thread 1/all nbconn 6
  option httplog
  option dontlognull
  option contstats
  option log-separate-errors
  option http-ignore-probes
  option http-use-proxy-header
  option httpslog
  option independent-streams
  option nolinger
  option originalto except 127.0.0.1
  option socket-stats
  option tcp-smart-accept
  option dontlog-normal
  option http-no-delay
  option splice-auto
  option splice-request
  option splice-response
  option idle-close-on-response
  option http-restrict-req-hdr-names delete
  acl invalid_src  src          0.0.0.0/7 224.0.0.0/3
  acl invalid_src  src_port     0:1023
  acl local_dst    hdr(host) -i localhost
  acl waf_wafTest_drop var(txn.wafTest.drop),bool
  monitor-uri /healthz
  monitor fail if site_dead
  filter trace name BEFORE-HTTP-COMP random-parsing hexdump
  filter compression
  filter trace name AFTER-HTTP-COMP random-forwarding
  filter fcgi-app my-app
  filter bwlim-in in default-limit 1k default-period 1s min-size 1m
  filter bwlim-out out limit 1024 key name(arg1) table st_src_global min-size 32
  http-request allow if src 192.168.0.0/16
  http-request set-header X-SSL %[ssl_fc]
  http-request set-var(req.my_var) req.fhdr(user-agent),lower
  http-request set-map(map.lst) %[src] %[req.hdr(X-Value)]
  http-request del-map(map.lst) %[src] if FALSE
  http-request del-acl(map.lst) %[src] if FALSE
  http-request cache-use cache-name if FALSE
  http-request disable-l7-retry if FALSE
  http-request early-hint hint-name %[src] if FALSE
  http-request replace-uri ^http://(.*) https://1 if FALSE
  http-request sc-add-gpc(0,1) 1 if FALSE
  http-request sc-inc-gpc(0,1) if FALSE
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
  http-request set-timeout client 20
  http-request set-bandwidth-limit my-limit limit 1m period 10s
  http-request set-bandwidth-limit my-limit-reverse period 20s limit 2m
  http-request set-bandwidth-limit my-limit-cond limit 3m if FALSE
  http-request track-sc0 src table tr0 if TRUE
  http-request track-sc1 src table tr1 if TRUE
  http-request track-sc2 src table tr2 if TRUE
  http-request track-sc5 src table test if TRUE
  http-request set-bc-mark 123 if TRUE
  http-request set-bc-tos 0x22
  http-request set-fc-mark hdr(port)
  http-request set-fc-tos 255 if FALSE
  http-request sc-set-gpt(1,2) hdr(Host),lower if FALSE
  http-response allow if src 192.168.0.0/16
  http-response set-header X-SSL %[ssl_fc]
  http-response set-var(req.my_var) req.fhdr(user-agent),lower
  http-response set-map(map.lst) %[src] %[res.hdr(X-Value)]
  http-response del-map(map.lst) %[src] if FALSE
  http-response del-acl(map.lst) %[src] if FALSE
  http-response cache-store cache-name if FALSE
  http-response sc-add-gpc(0,1) 1 if FALSE
  http-response sc-inc-gpc(0,1) if FALSE
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
  http-response set-bandwidth-limit my-limit limit 1m period 10s
  http-response set-bandwidth-limit my-limit-reverse period 20s limit 2m
  http-response set-bandwidth-limit my-limit-cond limit 3m if FALSE
  http-response track-sc0 src table tr0 if TRUE
  http-response track-sc1 src table tr1 if TRUE
  http-response track-sc2 src table tr2 if TRUE
  http-response track-sc5 src table test if TRUE
  http-response set-timeout server 20
  http-response set-timeout tunnel 20
  http-response set-timeout client 20
  http-response sc-set-gpt(1,2) 1234 if FALSE
  http-after-response set-map(map.lst) %[src] %[res.hdr(X-Value)]
  http-after-response del-map(map.lst) %[src] if FALSE
  http-after-response del-acl(map.lst) %[src] if FALSE
  http-after-response sc-add-gpc(0,1) 1 if FALSE
  http-after-response sc-inc-gpc(0,1) if FALSE
  http-after-response sc-inc-gpc0(0) if FALSE
  http-after-response sc-inc-gpc1(0) if FALSE
  http-after-response sc-set-gpt0(1) hdr(Host),lower if FALSE
  http-after-response sc-set-gpt0(1) 20 if FALSE
  http-after-response set-header Strict-Transport-Security "max-age=31536000"
  http-after-response set-log-level silent if FALSE
  http-after-response replace-header Set-Cookie (C=[^;]*);(.*) \1;ip=%bi;\2
  http-after-response replace-value Cache-control ^public$ private
  http-after-response set-status 503 reason "SlowDown"
  http-after-response set-var(sess.last_redir) res.hdr(location)
  http-after-response unset-var(sess.last_redir)
  http-after-response sc-set-gpt(1,2) hdr(port) if FALSE
  http-error status 400 content-type application/json lf-file /var/errors.file
  tcp-request connection accept if TRUE
  tcp-request connection reject if FALSE
  tcp-request content accept if TRUE
  tcp-request content reject if FALSE
  tcp-request connection silent-drop
  tcp-request connection silent-drop if TRUE
  tcp-request connection lua.foo param1 param2 if FALSE
  tcp-request connection sc-add-gpc(0,1) 1 if FALSE
  tcp-request content lua.foo param1 param2 if FALSE
  tcp-request content sc-add-gpc(0,1) 1 if FALSE
  tcp-request content set-bandwidth-limit my-limit limit 1m period 10s
  tcp-request content set-bandwidth-limit my-limit-reverse period 20s limit 2m
  tcp-request content set-bandwidth-limit my-limit-cond limit 3m if FALSE
  tcp-request connection set-mark 0x1Ab if FALSE
  tcp-request connection set-src-port hdr(port) if FALSE
  tcp-request connection set-tos 1 if FALSE
  tcp-request content set-log-level silent if FALSE
  tcp-request content set-mark 0x1Ac if FALSE
  tcp-request content set-nice 2 if FALSE
  tcp-request content set-src-port hdr(port) if FALSE
  tcp-request content set-tos 3 if FALSE
  tcp-request content set-var-fmt(req.tn) ssl_c_s_tn if FALSE
  tcp-request content switch-mode http proto my-proto if FALSE
  tcp-request session sc-add-gpc(0,1) 1 if FALSE
  tcp-request content track-sc0 src table tr0 if TRUE
  tcp-request connection track-sc0 src table tr0 if TRUE
  tcp-request session track-sc0 src table tr0 if TRUE
  tcp-request content track-sc5 src table test if TRUE
  tcp-request connection track-sc5 src table test if TRUE
  tcp-request session track-sc5 src table test if TRUE
  tcp-request session attach-srv srv1
  tcp-request session attach-srv srv2 name example.com
  tcp-request session attach-srv srv3 if is_cached
  tcp-request connection set-fc-mark 0xffffffff
  tcp-request connection set-fc-tos 0
  tcp-request session set-fc-mark 0
  tcp-request session set-fc-tos 0xff
  tcp-request content set-bc-mark 899 if TRUE
  tcp-request content set-bc-tos 2 if FALSE
  tcp-request content set-fc-mark hdr(port) if TRUE
  tcp-request content set-fc-tos req.hdr_cnt("X-Secret")
  tcp-request connection set-var-fmt(txn.ip_port) %%[dst]:%%[dst_port]
  tcp-request connection sc-set-gpt(1,2) 1234 if FALSE
  tcp-request content sc-set-gpt(1,2) hdr(port) if FALSE
  tcp-request session sc-set-gpt(1,2) 1234
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
  timeout client-fin 1000
  timeout tarpit 2000
  option tcpka
  option clitcpka
  unique-id-format %{+X}o%ci:%cp_%fi:%fp_%Ts_%rt:%pid
  unique-id-header X-Unique-ID
  no option accept-invalid-http-request
  no option h1-case-adjust-bogus-client
  compression algo identity gzip
  compression type text/plain
  compression offload
  compression algo-req raw-deflate
  compression algo-res raw-deflate identity
  compression type-req text/plain application/json
  compression type-res text/plain
  clitcpka-cnt 10
  clitcpka-idle 10s
  clitcpka-intvl 10
  stats auth admin:AdMiN123
  stats auth admin2:AdMiN1234
  stats show-modules
  stats realm HAProxy\\ Statistics
  email-alert from srv01@example.com
  email-alert to problems@example.com
  email-alert level warning
  email-alert myhostname srv01
  email-alert mailers localmailer1
  description this is a frontend description
  disabled
  id 123
  errorfile 403 /test/403.html
  errorfile 500 /test/500.html
  errorfile 429 /test/429.html
  errorloc302 404 http://www.myawesomesite.com/not_found
  errorloc303 404 http://www.myawesomesite.com/not_found
  error-log-format %T\ %t\ Some\ Text

frontend test_2 from test_defaults
  mode http
  option httplog
  option dontlognull
  option contstats
  option log-separate-errors
  no option http-ignore-probes
  no option http-use-proxy-header
  no option httpslog
  no option independent-streams
  no option nolinger
  option originalto except 127.0.0.1 header X-Client-Dst
  no option socket-stats
  no option tcp-smart-accept
  no option dontlog-normal
  no option http-no-delay
  no option splice-auto
  no option splice-request
  no option splice-response
  no option idle-close-on-response
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
  stats auth admin:AdMiN123
  stats auth admin2:AdMiN1234
  stats show-modules
  stats realm HAProxy\\ Statistics

backend test
  mode http
  balance roundrobin
  hash-type consistent sdbm avalanche
  hash-balance-factor 150
  log-tag bla
  option http-keep-alive
  option forwardfor header X-Forwarded-For
  option httpchk HEAD /
  option tcpka
  option srvtcpka
  option checkcache
  option independent-streams
  option nolinger
  option originalto header X-Client-Dst
  option persist
  option prefer-last-server
  option spop-check
  option tcp-smart-connect
  option transparent
  option splice-auto
  option splice-request
  option splice-response
  option http-restrict-req-hdr-names preserve
  default-server fall 2s rise 4s inter 5s port 8888 ws auto pool-low-conn 128 log-bufsize 6
  stick store-request src table test
  stick match src table test
  stick on src table test
  stick store-response src
  stick store-response src_port table test_port
  stick store-response src table test if TRUE
  tcp-response content accept if TRUE
  tcp-response content reject if FALSE
  tcp-response content lua.foo param1 param2 if FALSE
  tcp-response content set-bandwidth-limit my-limit limit 1m period 10s
  tcp-response content set-bandwidth-limit my-limit-reverse period 20s limit 2m
  tcp-response content set-bandwidth-limit my-limit-cond limit 3m if FALSE
  tcp-response content close if FALSE
  tcp-response content sc-add-gpc(0,1) 1 if FALSE
  tcp-response content sc-inc-gpc0(1) if FALSE
  tcp-response content sc-inc-gpc1(2) if FALSE
  tcp-response content sc-set-gpt0(3) hdr(Host),lower if FALSE
  tcp-response content send-spoe-group engine group if FALSE
  tcp-response content set-log-level silent if FALSE
  tcp-response content set-mark 0x1Ab if FALSE
  tcp-response content set-nice 1 if FALSE
  tcp-response content set-tos 2 if FALSE
  tcp-response content silent-drop if FALSE
  tcp-response content unset-var(req.my_var) if FALSE
  tcp-response content set-fc-mark 7676 if TRUE
  tcp-response content set-fc-tos 0xab if FALSE
  tcp-response content sc-set-gpt(1,2) 1234
  option contstats
  timeout check 2s
  timeout tunnel 5s
  timeout server 3s
  timeout server-fin 1000
  timeout tarpit 2000
  cookie BLA rewrite httponly nocache
  option external-check
  external-check command /bin/false
  use-server webserv if TRUE
  use-server webserv2 unless TRUE
  server webserv 192.168.1.1:9200 maxconn 1000 ssl weight 10 inter 2s cookie BLAH slowstart 6000 proxy-v2-options authority,crc32c ws h1 pool-low-conn 128 id 1234 pool-purge-delay 10s tcp-ut 2s curves secp384r1 client-sigalgs ECDSA+SHA256:RSA+SHA256 sigalgs ECDSA+SHA256 log-bufsize 10 set-proxy-v2-tlv-fmt(0x20) %[fc_pp_tlv(0x20)]
  server webserv2 192.168.1.1:9300 maxconn 1000 ssl weight 10 inter 2s cookie BLAH slowstart 6000 proxy-v2-options authority,crc32c ws h1 pool-low-conn 128
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
  compression direction both
  compression algo-req raw-deflate
  compression algo-res raw-deflate identity
  compression type-req text/plain application/json
  compression type-res text/plain
  srvtcpka-cnt 10
  srvtcpka-idle 10s
  srvtcpka-intvl 10
  stats http-request auth realm HAProxy\\ Statistics
  stats http-request allow if something
  stats auth admin:AdMiN123
  stats auth admin2:AdMiN1234
  stats show-modules
  stats realm HAProxy\\ Statistics
  email-alert from prod01@example.com
  email-alert to sre@example.com
  email-alert level warning
  email-alert myhostname prod01
  email-alert mailers localmailer1
  load-server-state-from-file local
  server-state-file-name use-backend-name
  description this is a backend description
  use-fcgi-app app-name
  enabled
  id 456
  errorfile 403 /test/403.html
  errorfile 500 /test/500.html
  errorfile 429 /test/429.html
  errorfiles my_errors 404 401 500
  errorfiles other_errors
  errorfiles another_errors 501
  errorloc302 404 http://www.myawesomesite.com/not_found
  errorloc303 404 http://www.myawesomesite.com/not_found
  error-log-format %T\ %t\ Some\ Text
  fullconn 11
  max-keep-alive-queue 101
  ignore-persist if acl-name
  ignore-persist unless local_dst
  force-persist unless acl-name-2
  force-persist if acl-name-3
  retry-on 504 505
  http-send-name-header X-My-Awesome-Header
  persist rdp-cookie(name)
  source 192.168.1.222 usesrc hdr_ip(hdr,occ)
  http-response set-fc-mark 123
  http-response set-fc-tos 1 if TRUE

peers mycluster
  enabled
  default-server fall 2s rise 4s inter 5s port 8888 slowstart 6000
  default-bind v4v6 ssl crt /etc/haproxy/site.pem alpn h2,http/1.1
  peer hapee 192.168.1.1:1023 shard 1
  peer aggregator HARDCODEDCLUSTERIP:10023
  shards 3
  table t1 type string len 1000 size 200k expire 5m nopurge store gpc0,conn_rate(30s)
  table t2 type string len 1000 size 200k expire 5m nopurge store gpc0 store gpc1,conn_rate(30s)
  table t9 type string len 1000 size 200k expire 5m write-to t2 nopurge store gpc0,conn_rate(30s)

program test
  command echo "Hello, World!"
  user hapee-lb
  group hapee
  option start-on-reload

program test_2
  command echo "Hello, World!"
  no option start-on-reload

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

fcgi-app test
  acl invalid_src src 0.0.0.0/7 224.0.0.0/3
  acl invalid_src src_port 0:1023
  acl local_dst hdr(host) -i localhost
  docroot /path/to/chroot
  index index.php
  path-info ^(/.+\.php)(/.*)?$
  option get-values
  no option keep-conn
  no option mpxs-conns
  option max-reqs 1024
  set-param name fmt if acl
  set-param name fmt unless acl
  set-param name fmt
  pass-header x-header unless acl
  pass-header x-header if acl
  pass-header x-header
  log-stderr 127.0.0.1:1515 len 8192 format rfc5424 sample 1,2-5:6 local2 info debug
  log-stderr 127.0.0.1:1515 len 8192 format rfc5424 sample 1,2-5:6 local2 info
  log-stderr 127.0.0.1:1515 local2
  log-stderr global


fcgi-app test_2
  acl invalid_src src 0.0.0.0/7 224.0.0.0/3
  acl invalid_src src_port 0:1023
  acl local_dst hdr(host) -i localhost
  docroot /path/to/chroot
  index index.php
  path-info ^(/.+\.php)(/.*)?$
  option get-values
  no option keep-conn
  no option mpxs-conns
  set-param name fmt if acl
  set-param name fmt unless acl
  set-param name fmt
  option max-reqs 1024
  pass-header x-header unless acl
  pass-header x-header if acl
  pass-header x-header
  log-stderr 127.0.0.1:1515 len 8192 format rfc5424 sample 1,2-5:6 local2 info debug
  log-stderr 127.0.0.1:1515 len 8192 format rfc5424 sample 1,2-5:6 local2 info
  log-stderr 127.0.0.1:1515 local2
  log-stderr global

ring myring
  description "My local buffer"
  format rfc3164
  maxlen 1200
  size 32764
  timeout connect 5s
  timeout server 10s
  server mysyslogsrv 127.0.0.1:6514 log-proto octet-count
  server s1 192.168.1.1:80 check resolve-opts allow-dup-ip,ignore-weight resolve-net 10.0.0.0/8,10.200.200.0/12

log-forward sylog-loadb
  dgram-bind 127.0.0.1:1514 transparent name webserv
  bind 127.0.0.1:1514
  backlog 10
  maxconn 1000
  timeout client 10s
  log global
  log ring@myring local0
  log 127.0.0.1:10001 sample 1:4 local0
  log 127.0.0.1:10002 sample 2:4 local0
  log 127.0.0.1:10003 sample 3:4 local0
  log 127.0.0.1:10004 sample 4:4 local0

mailers localmailer1
  mailer smtp1 10.0.10.1:514
  mailer smtp2 10.0.10.2:514
  timeout mail 15s

http-errors website-1
  errorfile 400 /etc/haproxy/errorfiles/site1/400.http
  errorfile 404 /etc/haproxy/errorfiles/site1/404.http
  errorfile 408 /dev/null  # work around Chrome pre-connect bug

http-errors website-2
  errorfile 400 /etc/haproxy/errorfiles/site2/400.http
  errorfile 404 /etc/haproxy/errorfiles/site2/404.http
  errorfile 501 /etc/haproxy/errorfiles/site2/501.http

backend test_2 from test_defaults_2
  mode http
  balance roundrobin
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
  option originalto header X-Client-Dst except 127.0.0.1
  no option persist
  no option prefer-last-server
  no option spop-check
  no option tcp-smart-connect
  no option transparent
  no option splice-auto
  no option splice-request
  no option splice-response
  default-server fall 2s rise 4s inter 5s port 8888 slowstart 6000
  option contstats
  timeout check 2s
  timeout tunnel 5s
  timeout server 3s
  cookie BLA rewrite httponly nocache
  stick-table type ip size 100k expire 1h peers mycluster write-to t99 store http_req_rate(10s)
  http-check expect rstatus some-pattern
  http-error status 200 content-type "text/plain" string "My content" hdr Some-Header value
  http-error status 503 content-type application/json string "My content" hdr Additional-Header value1 hdr Some-Header value
  srvtcpka-cnt 10
  srvtcpka-idle 10s
  srvtcpka-intvl 10
  stats http-request auth realm HAProxy\\ Statistics
  stats http-request allow if something
  stats auth admin:AdMiN123
  stats auth admin2:AdMiN1234
  stats show-modules
  stats realm HAProxy\\ Statistics
  email-alert from prod01@example.com
  email-alert to sre@example.com
  email-alert level warning
  email-alert myhostname prod01
  email-alert mailers localmailer1

userlist first
	group G1 users tiger,scott
	group G2 users scott
	user tiger password $6$k6y3o.eP$JlKBx9za9667qe4xHSwRv6J.C0/D7cV91
	user scott insecure-password elgato

userlist second
	group one
	group two
  group three
	user neo password $6$k6y3o.eP$JlKBxxHSwRv6J.C0/D7cV91 groups one
	user thomas insecure-password white-rabbit groups one,two
	user anderson insecure-password hello groups two
`
const testPath = "/tmp/haproxy-test.cfg"

//nolint:gochecknoglobals
var (
	clientTest configuration.Configuration
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

func prepareClient(path string) (c configuration.Configuration, err error) {
	c, err = configuration.New(context.Background(),
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
