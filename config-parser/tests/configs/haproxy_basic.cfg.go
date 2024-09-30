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

const configBasic1 = `# _version=1
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

const configBasicMissingNewLineOnEnd = `# _version=1
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
  http-request deny deny_status 400 # deny`

const configBasicWithComments = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults A
  # line comment 1
  log global

frontend http from A
  mode http # inline comment #1
  # line comment 2
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2 # inline comment #2
  default_backend default_backend

backend default_backend from A
  # line comment 3
  mode http
  http-request deny deny_status 400 # deny
`

const configBasicWithLineComments = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults A
  # line comment 1
  log global

frontend http from A
  mode http
  # line comment 2
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend

backend default_backend from A
  # line comment 3
  mode http
  http-request deny deny_status 400 # deny
`

const configSnippet = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  master-worker

defaults A
  # line comment 1
  log global

frontend http from A
  mode http
  ###_config-snippet_### BEGIN
  ### origin:ingress:haproxy-controller/openspeedtest ###
  ### a comment
  tune.ssl.default-dh-param 2048 # inline a long comment
  tune.bufsize 32768
  ###_config-snippet_### END
  default_backend default_backend

backend default_backend from A
  # line comment 3
  mode http
  http-request deny deny_status 400 # deny
`
