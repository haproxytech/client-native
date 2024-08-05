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
package configs

const configBasicUseV2HTTPCheck = `
frontend http
  default_backend default_backend

backend api.domain.com
  balance leastconn
  option httpchk
  http-check send meth HEAD uri / ver HTTP/1.1 hdr Host api.domain.com hdr User-Agent healthcheck/haproxy
  http-check expect rstatus ((2|3)[0-9][0-9]|404)
  http-check connect default
  default-server check inter 5s fall 3 rise 2 maxconn 256 ssl verify none
  http-request add-header X-Forwarded-For %[src]
  cookie SERVERID indirect nocache insert
  server api01 127.0.0.1:443 cookie 01
`
