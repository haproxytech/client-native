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

const configDefaultsNoName1 = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults
  log global

frontend http
  mode http
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend

backend default_backend
  mode http
  http-request deny deny_status 400 # deny
`

const configDefaultsNoName1Result = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults unnamed_defaults_1
  log global

frontend http from unnamed_defaults_1
  mode http
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend

backend default_backend from unnamed_defaults_1
  mode http
  http-request deny deny_status 400 # deny
`

const configDefaultsNoName2 = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults
  log global

defaults
  mode tcp

frontend http
  mode http
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend

backend default_backend
  mode http
  http-request deny deny_status 400 # deny
`

const configDefaultsNoName2Result = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults unnamed_defaults_1
  log global

defaults unnamed_defaults_2
  mode tcp

frontend http from unnamed_defaults_2
  mode http
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend

backend default_backend from unnamed_defaults_2
  mode http
  http-request deny deny_status 400 # deny
`

const configDefaultsNoName3 = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults A
  log global

defaults
  mode tcp

frontend http
  mode http
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend

backend default_backend from A
  mode http
  http-request deny deny_status 400 # deny
`

const configDefaultsNoName3Result = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults A
  log global

defaults unnamed_defaults_1
  mode tcp

frontend http from unnamed_defaults_1
  mode http
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend

backend default_backend from A
  mode http
  http-request deny deny_status 400 # deny
`

const configDefaultsNoName4 = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults
  log global

defaults A
  mode tcp

frontend http
  mode http
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend

backend default_backend from A
  mode http
  http-request deny deny_status 400 # deny
`

const configDefaultsNoName4Result = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults A
  mode tcp

defaults unnamed_defaults_1
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
