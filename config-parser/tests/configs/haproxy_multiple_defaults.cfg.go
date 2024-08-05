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

const configDefaults = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults A
  log global

frontend http from A
  mode http
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend

backend default_backend from A
  mode http
  http-request deny deny_status 400 # deny
`

const configDefaultsTwo = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults A
  log global

defaults withName
  log global

frontend http from A
  mode http
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend

backend default_backend from withName
  mode http
  http-request deny deny_status 400 # deny
`

const configDefaultsThree = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults A
  log global

defaults withName
  log global

defaults withName2
  log global

frontend http from A
  mode http
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend

backend default_backend from A
  mode http
  http-request deny deny_status 400 # deny
`

const configDefaultsTwoSecond = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults A
  log global

defaults B
  log global

frontend http1 from A
  mode http
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend1

frontend http2 from B
  mode http
  bind 0.0.0.0:81 name bind_1
  bind :::81 v4v6 name bind_2
  default_backend default_backend2

backend default_backend1 from B
  mode http
  http-request deny deny_status 400 # deny

backend default_backend2 from A
  mode http
  http-request deny deny_status 400 # deny
`

const configDefaultsWithFrom = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults A
  log global

defaults B from A
  log global

frontend http1 from A
  mode http
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend1

frontend http2 from B
  mode http
  bind 0.0.0.0:81 name bind_1
  bind :::81 v4v6 name bind_2
  default_backend default_backend2

backend default_backend1 from B
  mode http
  http-request deny deny_status 400 # deny

backend default_backend2 from A
  mode http
  http-request deny deny_status 400 # deny
`

const configDefaultsWithFromResult2 = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults A
  log global

defaults B from A
  log global

frontend http1 from B
  mode http
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend1

frontend http2 from B
  mode http
  bind 0.0.0.0:81 name bind_1
  bind :::81 v4v6 name bind_2
  default_backend default_backend2

backend default_backend1 from B
  mode http
  http-request deny deny_status 400 # deny

backend default_backend2 from B
  mode http
  http-request deny deny_status 400 # deny
`
