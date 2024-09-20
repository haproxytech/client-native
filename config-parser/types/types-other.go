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

package types

import "github.com/haproxytech/client-native/v6/config-parser/common"

//name:section
//no:sections
//dir:extra
//no:init
type Section struct {
	Name         string
	FromDefaults string
	Comment      string
}

//name:config-version
//no:sections
//dir:extra
//no:init
//no:get
type ConfigVersion struct {
	Value int64
}

//name:config-hash
//no:sections
//dir:extra
//no:init
//no:get
type ConfigHash struct {
	Value string
}

//name:comments
//no:sections
//dir:extra
//is:multiple
//no:init
//no:parse
type Comments struct {
	Value string
}

//name:unprocessed
//no:sections
//dir:extra
//is:multiple
//no:init
//no:parse
//test:skip
type UnProcessed struct {
	Value string
}

//name:simple-option
//no:sections
//struct:name:Option
//dir:simple
//no:init
type SimpleOption struct {
	NoOption bool
	Comment  string
}

//name:simple-timeout
//no:sections
//struct:name:Timeout
//dir:simple
//no:init
type SimpleTimeout struct {
	Value   string
	Comment string
}

//name:simple-word
//no:sections
//struct:name:Word
//dir:simple
//parser:type:StringC
type SimpleWord struct{}

//name:simple-number
//no:sections
//struct:name:Number
//dir:simple
//parser:type:Int64C
type SimpleNumber struct{}

//name:simple-string
//no:sections
//struct:name:String
//dir:simple
//parser:type:StringC
type SimpleString struct{}

//name:simple-on-off
//no:sections
//struct:name:OnOff
//dir:simple
//parser:type:StringC
type SimpleOnOff struct{}

//name:simple-auto-on-off
//no:sections
//struct:name:AutoOnOff
//dir:simple
//parser:type:StringC
type SimpleAutoOnOff struct{}

//name:simple-string-slice
//no:sections
//struct:name:StringSlice
//dir:simple
//parser:type:StringSliceC
type SimpleStringSlice struct{}

//name:simple-string-kv
//no:sections
//struct:name:StringKeyValue
//dir:simple
//parser:type:StringKeyValueC
type SimpleStringKeyValue struct{}

//name:array-string-kv
//is:multiple
//no:sections
//no:parse
//struct:name:ArrayKeyValue
//dir:simple
//parser:type:StringKeyValueC
type ArrayStringKeyValue struct{}

//name:simple-time
//no:sections
//struct:name:Time
//dir:simple
//parser:type:StringC
type SimpleTime struct{}

//name:simple-size
//no:sections
//struct:name:Size
//dir:simple
//parser:type:StringC
type SimpleSize struct{}

//name:simple-enabled
//no:sections
//struct:name:Enabled
//dir:simple
//parser:type:Enabled
type SimpleEnabled struct{}

//name:simple-time-two-words
//no:sections
//struct:name:TimeTwoWords
//dir:simple
//no:init
//parser:type:StringC
//test:skip
type TimeTwoWords struct{}

type Filter interface {
	Parse(parts []string, comment string) error
	Result() common.ReturnResultLine
}

//name:filter
//no:sections
//dir:filters
//is:multiple
//parser:type:Filter
//is:interface
//no:init
//no:parse
//test:ok:filter bwlim-in name default-limit 1024 default-period 10
//test:ok:filter bwlim-in name default-limit 1024 default-period 10 min-size 32
//test:ok:filter bwlim-in name limit 1024 key name(arg1)
//test:ok:filter bwlim-in name limit 1024 key name(arg1) table st_src_global
//test:ok:filter bwlim-in name limit 1024 key name(arg1) table st_src_global min-size 32
//test:ok:filter bwlim-out name default-limit 1024 default-period 10
//test:ok:filter bwlim-out name default-limit 1024 default-period 10 min-size 32
//test:ok:filter bwlim-out name limit 1024 key name(arg1)
//test:ok:filter bwlim-out name limit 1024 key name(arg1) table st_src_global
//test:ok:filter bwlim-out name limit 1024 key name(arg1) table st_src_global min-size 32
//test:ok:filter opentracing id qwerty-1234-uiop-567890 config file
//test:ok:filter opentracing config file
//test:ok:filter fcgi-app my-application
//test:ok:filter compression
//test:ok:filter spoe config file
//test:ok:filter spoe engine name config file
//test:ok:filter trace name name random-parsing random-forwarding hexdump
//test:ok:filter trace random-parsing random-forwarding hexdump
//test:ok:filter trace random-forwarding hexdump
//test:ok:filter trace hexdump
//test:ok:filter trace
//test:fail:filter bwlim-in
//test:fail:filter bwlim-in name
//test:fail:filter bwlim-in name default-limit
//test:fail:filter bwlim-in name default-limit 1024
//test:fail:filter bwlim-in name default-limit 1024 default-period
//test:fail:filter bwlim-in name default-limit 1024 default-period 10 min-size
//test:fail:filter bwlim-in name default-limit 1024 key name(arg1)
//test:fail:filter bwlim-in name limit 1024 default-period 100
//test:fail:filter bwlim-out
//test:fail:filter bwlim-out name
//test:fail:filter bwlim-out name default-limit
//test:fail:filter bwlim-out name default-limit 1024
//test:fail:filter bwlim-out name default-limit 1024 default-period
//test:fail:filter bwlim-out name default-limit 1024 default-period 10 min-size
//test:fail:filter bwlim-out name limit
//test:fail:filter bwlim-out name limit 1024
//test:fail:filter bwlim-out name limit 1024 key
//test:fail:filter bwlim-out name limit 1024 key name(arg1) table
//test:fail:filter bwlim-out name limit 1024 key name(arg1) table st_src_global min-size
//test:fail:filter opentracing
//test:fail:filter opentracing id
//test:fail:filter opentracing id qwerty-1234-uiop-567890 config
//test:fail:filter opentracing id qwerty-1234-uiop-567890 config file extra
//test:fail:filter opentracing id qwerty-1234-uiop-567890 config file extra option
//test:fail:filter fcgi-app
//test:fail:filter fcgi-app first second
//test:fail:filter compression false
//test:fail:filter spoe
//test:fail:filter spoe config
//test:fail:filter spoe engine
//test:fail:filter spoe engine config
//test:fail:filter trace name
//test:fail:filter trace 0 name
//test:fail:filter spoe l : d 8 t 8 t t c t t t 8 t 8 t t t 8 t t t 8 t 8 t t 8 t t t 8 8 t config
//test:fail:filter cache
type Filters struct{}

type ParserType int

const (
	HTTP ParserType = iota
	TCP
)

type Action interface {
	Parse(parts []string, parserType ParserType, comment string) error
	String() string
	GetComment() string
}

//sections:frontend,backend,defaults
//name:http-request
//struct:name:Requests
//dir:http
//is:multiple
//parser:type:Action
//is:interface
//no:init
//no:parse
//test:fail:http-request
//test:fail:http-request capture req.cook_cnt(FirstVisit),bool strlen 10
//test:frontend-ok:http-request capture req.cook_cnt(FirstVisit),bool len 10
//test:ok:http-request set-map(map.lst) %[src] %[req.hdr(X-Value)] if value
//test:ok:http-request set-map(map.lst) %[src] %[req.hdr(X-Value)]
//test:fail:http-request set-map(map.lst) %[src]
//test:ok:http-request add-acl(map.lst) [src]
//test:fail:http-request add-acl(map.lst)
//test:ok:http-request add-header X-value value
//test:quote_ok:http-request add-header Authorization Basic\ eC1oYXByb3h5LXJlY3J1aXRzOlBlb3BsZSB3aG8gZGVjb2RlIG1lc3NhZ2VzIG9mdGVuIGxvdmUgd29ya2luZyBhdCBIQVByb3h5LiBEbyBub3QgYmUgc2h5LCBjb250YWN0IHVz
//test:quote_ok:http-request add-header Authorisation "Basic eC1oYXByb3h5LXJlY3J1aXRzOlBlb3BsZSB3aG8gZGVjb2RlIG1lc3NhZ2VzIG9mdGVuIGxvdmUgd29ya2luZyBhdCBIQVByb3h5LiBEbyBub3QgYmUgc2h5LCBjb250YWN0IHVz"
//test:fail:http-request add-header X-value
//test:ok:http-request cache-use cache-name
//test:ok:http-request cache-use cache-name if FALSE
//test:fail:http-request cache-use
//test:fail:http-request cache-use if FALSE
//test:ok:http-request del-acl(map.lst) [src]
//test:fail:http-request del-acl(map.lst)
//test:ok:http-request allow
//test:ok:http-request auth
//test:ok:http-request del-header X-value
//test:ok:http-request del-header X-value if TRUE
//test:ok:http-request del-header X-value -m str if TRUE
//test:fail:http-request del-header
//test:fail:http-request del-header X-value -m if TRUE
//test:fail:http-request del-header X-value bla
//test:ok:http-request del-map(map.lst) %[src] if ! value
//test:ok:http-request del-map(map.lst) %[src]
//test:fail:http-request del-map(map.lst)
//test:ok:http-request deny
//test:ok:http-request deny deny_status 400
//test:ok:http-request deny if TRUE
//test:ok:http-request deny deny_status 400 if TRUE
//test:ok:http-request deny deny_status 400 content-type application/json if TRUE
//test:ok:http-request deny deny_status 400 content-type application/json
//test:ok:http-request deny deny_status 400 content-type application/json default-errorfiles
//test:ok:http-request deny deny_status 400 content-type application/json errorfile errors
//test:ok:http-request deny deny_status 400 content-type application/json string error if TRUE
//test:ok:http-request deny deny_status 400 content-type application/json lf-string error hdr host google.com if TRUE
//test:ok:http-request deny deny_status 400 content-type application/json file /var/errors.file
//test:ok:http-request deny deny_status 400 content-type application/json lf-file /var/errors.file
//test:ok:http-request deny deny_status 400 content-type application/json string error hdr host google.com if TRUE
//test:ok:http-request deny deny_status 400 content-type application/json string error hdr host google.com hdr x-value bla if TRUE
//test:ok:http-request deny deny_status 400 content-type application/json string error hdr host google.com hdr x-value bla
//test:fail:http-request deny test test
//test:ok:http-request disable-l7-retry
//test:ok:http-request disable-l7-retry if FALSE
//test:ok:http-request early-hint hint %[src]
//test:ok:http-request early-hint hint %[src] if FALSE
//test:ok:http-request early-hint if FALSE
//test:fail:http-request early-hint hint
//test:fail:http-request early-hint hint if FALSE
//test:ok:http-request lua.foo
//test:ok:http-request lua.foo if FALSE
//test:ok:http-request lua.foo param
//test:ok:http-request lua.foo param param2
//test:fail:http-request lua.
//test:fail:http-request lua. if FALSE
//test:fail:http-request lua. param
//test:ok:http-request normalize-uri fragment-encode
//test:ok:http-request normalize-uri fragment-encode if TRUE
//test:ok:http-request normalize-uri fragment-strip
//test:ok:http-request normalize-uri fragment-strip if TRUE
//test:ok:http-request normalize-uri path-merge-slashes
//test:ok:http-request normalize-uri path-merge-slashes if TRUE
//test:ok:http-request normalize-uri path-strip-dot
//test:ok:http-request normalize-uri path-strip-dot if TRUE
//test:ok:http-request normalize-uri path-strip-dotdot
//test:ok:http-request normalize-uri path-strip-dotdot full
//test:ok:http-request normalize-uri path-strip-dotdot if TRUE
//test:ok:http-request normalize-uri path-strip-dotdot full if TRUE
//test:ok:http-request normalize-uri percent-decode-unreserved
//test:ok:http-request normalize-uri percent-decode-unreserved if TRUE
//test:ok:http-request normalize-uri percent-decode-unreserved strict
//test:ok:http-request normalize-uri percent-decode-unreserved strict if TRUE
//test:ok:http-request normalize-uri percent-to-uppercase
//test:ok:http-request normalize-uri percent-to-uppercase if TRUE
//test:ok:http-request normalize-uri percent-to-uppercase strict
//test:ok:http-request normalize-uri percent-to-uppercase strict if TRUE
//test:ok:http-request normalize-uri query-sort-by-name
//test:ok:http-request normalize-uri query-sort-by-name if TRUE
//test:fail:http-request normalize-uri bla
//test:fail:http-request normalize-uri path-strip-dot strict
//test:fail:http-request normalize-uri path-strip-dot full
//test:fail:http-request normalize-uri if TRUE
//test:fail:http-request normalize-uri
//test:ok:http-request redirect prefix https://mysite.com
//test:fail:http-request redirect prefix
//test:ok:http-request reject
//test:ok:http-request replace-header User-agent curl foo
//test:fail:http-request replace-header User-agent curl
//test:ok:http-request replace-path (.*) /foo
//test:fail:http-request replace-path (.*)
//test:ok:http-request replace-path (.*) /foo if TRUE
//test:fail:http-request replace-path (.*) if TRUE
//test:ok:http-request replace-pathq (.*) /foo
//test:fail:http-request replace-pathq (.*)
//test:ok:http-request replace-pathq (.*) /foo if TRUE
//test:fail:http-request replace-pathq (.*) if TRUE
//test:ok:http-request replace-uri ^http://(.*) https://1
//test:ok:http-request replace-uri ^http://(.*) https://1 if FALSE
//test:fail:http-request replace-uri ^http://(.*)
//test:fail:http-request replace-uri
//test:fail:http-request replace-uri ^http://(.*) if FALSE
//test:ok:http-request replace-value X-Forwarded-For ^192.168.(.*)$ 172.16.1
//test:fail:http-request replace-value X-Forwarded-For ^192.168.(.*)$
//test:ok:http-request sc-add-gpc(1,2) 1
//test:ok:http-request sc-add-gpc(1,2) 1 if is-error
//test:fail:http-request sc-add-gpc
//test:ok:http-request sc-inc-gpc(1,2)
//test:ok:http-request sc-inc-gpc(1,2) if FALSE
//test:fail:http-request sc-inc-gpc
//test:ok:http-request sc-inc-gpc0(1)
//test:ok:http-request sc-inc-gpc0(1) if FALSE
//test:fail:http-request sc-inc-gpc0
//test:ok:http-request sc-inc-gpc1(1)
//test:ok:http-request sc-inc-gpc1(1) if FALSE
//test:fail:http-request sc-inc-gpc1
//test:ok:http-request sc-set-gpt(1,2) hdr(Host),lower if FALSE
//test:ok:http-request sc-set-gpt0(1) hdr(Host),lower
//test:ok:http-request sc-set-gpt0(1) 10
//test:ok:http-request sc-set-gpt0(1) hdr(Host),lower if FALSE
//test:fail:http-request sc-set-gpt0(1)
//test:fail:http-request sc-set-gpt0
//test:fail:http-request sc-set-gpt0(1) if FALSE
//test:ok:http-request send-spoe-group engine group
//test:fail:http-request send-spoe-group engine
//test:ok:http-request set-header X-value value
//test:fail:http-request set-header X-value
//test:ok:http-request set-log-level silent
//test:fail:http-request set-log-level
//test:ok:http-request set-mark 20
//test:ok:http-request set-mark 0x1Ab
//test:fail:http-request set-mark
//test:ok:http-request set-nice 0
//test:ok:http-request set-nice 0 if FALSE
//test:fail:http-request set-nice
//test:ok:http-request set-method POST
//test:ok:http-request set-method POST if FALSE
//test:fail:http-request set-method
//test:ok:http-request set-path /%[hdr(host)]%[path]
//test:fail:http-request set-path
//test:ok:http-request set-pathq /%[hdr(host)]%[path]
//test:fail:http-request set-pathq
//test:ok:http-request set-priority-class req.hdr(priority)
//test:ok:http-request set-priority-class req.hdr(priority) if FALSE
//test:fail:http-request set-priority-class
//test:ok:http-request set-priority-offset req.hdr(offset)
//test:ok:http-request set-priority-offset req.hdr(offset) if FALSE
//test:fail:http-request set-priority-offset
//test:ok:http-request set-query %[query,regsub(%3D,=,g)]
//test:fail:http-request set-query
//test:ok:http-request set-src hdr(src)
//test:ok:http-request set-src hdr(src) if FALSE
//test:fail:http-request set-src
//test:ok:http-request set-src-port hdr(port)
//test:ok:http-request set-src-port hdr(port) if FALSE
//test:fail:http-request set-src-port
//test:ok:http-request set-timeout server 20
//test:ok:http-request set-timeout tunnel 20
//test:ok:http-request set-timeout tunnel 20s if TRUE
//test:ok:http-request set-timeout server 20s if TRUE
//test:ok:http-request set-timeout client 20
//test:ok:http-request set-timeout client 20s if TRUE
//test:fail:http-request set-timeout fake-timeout 20s if TRUE
//test:ok:http-request set-tos 0 if FALSE
//test:ok:http-request set-tos 0
//test:fail:http-request set-tos
//test:ok:http-request set-uri /%[hdr(host)]%[path]
//test:fail:http-request set-uri
//test:ok:http-request set-var(req.my_var) req.fhdr(user-agent),lower
//test:fail:http-request set-var(req.my_var)
//test:ok:http-request set-var-fmt(req.my_var) req.fhdr(user-agent),lower
//test:fail:http-request set-var-fmt(req.my_var)
//test:ok:http-request silent-drop
//test:ok:http-request silent-drop if FALSE
//test:ok:http-request strict-mode on
//test:ok:http-request strict-mode on if FALSE
//test:fail:http-request strict-mode
//test:fail:http-request strict-mode if FALSE
//test:ok:http-request tarpit
//test:ok:http-request tarpit deny_status 400
//test:ok:http-request tarpit if TRUE
//test:ok:http-request tarpit deny_status 400 if TRUE
//test:ok:http-request tarpit deny_status 400 content-type application/json if TRUE
//test:ok:http-request tarpit deny_status 400 content-type application/json
//test:ok:http-request tarpit deny_status 400 content-type application/json default-errorfiles
//test:ok:http-request tarpit deny_status 400 content-type application/json errorfile errors
//test:ok:http-request tarpit deny_status 400 content-type application/json string error if TRUE
//test:ok:http-request tarpit deny_status 400 content-type application/json lf-string error hdr host google.com if TRUE
//test:ok:http-request tarpit deny_status 400 content-type application/json file /var/errors.file
//test:ok:http-request tarpit deny_status 400 content-type application/json lf-file /var/errors.file
//test:ok:http-request tarpit deny_status 400 content-type application/json string error hdr host google.com if TRUE
//test:ok:http-request tarpit deny_status 400 content-type application/json string error hdr host google.com hdr x-value bla if TRUE
//test:ok:http-request tarpit deny_status 400 content-type application/json string error hdr host google.com hdr x-value bla
//test:fail:http-request tarpit test test
//test:ok:http-request track-sc0 src
//test:fail:http-request track-sc0
//test:ok:http-request track-sc1 src
//test:fail:http-request track-sc1
//test:ok:http-request track-sc2 src
//test:fail:http-request track-sc2
//test:ok:http-request track-sc5 src
//test:ok:http-request track-sc5 src table a_table
//test:ok:http-request track-sc5 src table a_table if some_cond
//test:ok:http-request track-sc5 src if some_cond
//test:fail:http-request track-sc
//test:fail:http-request track-sc5
//test:fail:http-request track-sc5 src table
//test:fail:http-request track-sc5 src if
//test:fail:http-request track-sc src if some_cond
//test:fail:http-request track-sc src table a_table if some_cond
//test:ok:http-request unset-var(req.my_var)
//test:ok:http-request unset-var(req.my_var) if FALSE
//test:fail:http-request unset-var(req.)
//test:fail:http-request unset-var(req)
//test:ok:http-request wait-for-body time 20s
//test:ok:http-request wait-for-body time 20s if TRUE
//test:ok:http-request wait-for-body time 20s at-least 100k
//test:ok:http-request wait-for-body time 20s at-least 100k if TRUE
//test:fail:http-request wait-for-body 20s at-least 100k
//test:fail:http-request wait-for-body time 2000 test
//test:ok:http-request wait-for-handshake
//test:ok:http-request wait-for-handshake if FALSE
//test:ok:http-request do-resolve(txn.myip,mydns) hdr(Host),lower
//test:ok:http-request do-resolve(txn.myip,mydns) hdr(Host),lower if { var(txn.myip) -m found }
//test:ok:http-request do-resolve(txn.myip,mydns) hdr(Host),lower unless { var(txn.myip) -m found }
//test:ok:http-request do-resolve(txn.myip,mydns,ipv4) hdr(Host),lower
//test:ok:http-request do-resolve(txn.myip,mydns,ipv6) hdr(Host),lower
//test:fail:http-request do-resolve(txn.myip)
//test:fail:http-request do-resolve(txn.myip,mydns)
//test:fail:http-request do-resolve(txn.myip,mydns,ipv4)
//test:ok:http-request set-dst var(txn.myip)
//test:ok:http-request set-dst var(txn.myip) if { var(txn.myip) -m found }
//test:ok:http-request set-dst var(txn.myip) unless { var(txn.myip) -m found }
//test:fail:http-request set-dst
//test:ok:http-request set-dst-port hdr(x-port)
//test:ok:http-request set-dst-port hdr(x-port) if { var(txn.myip) -m found }
//test:ok:http-request set-dst-port hdr(x-port) unless { var(txn.myip) -m found }
//test:ok:http-request set-dst-port int(4000)
//test:fail:http-request set-dst-port
//test:quote_ok:http-request return status 200 content-type "text/plain" string "My content" if { var(txn.myip) -m found }
//test:quote_ok:http-request return status 200 content-type "text/plain" string "My content" unless { var(txn.myip) -m found }
//test:quote_ok:http-request return content-type "text/plain" string "My content" if { var(txn.myip) -m found }
//test:quote_ok:http-request return content-type 'text/plain' string 'My content' if { var(txn.myip) -m found }
//test:quote_ok:http-request return content-type "text/plain" lf-string "Hello, you are: %[src]" if { var(txn.myip) -m found }
//test:quote_ok:http-request return content-type "text/plain" file /my/fancy/response/file if { var(txn.myip) -m found }
//test:quote_ok:http-request return content-type "text/plain" lf-file /my/fancy/lof/format/response/file if { var(txn.myip) -m found }
//test:quote_ok:http-request return content-type "text/plain" string "My content" hdr X-value value if { var(txn.myip) -m found }
//test:quote_ok:http-request return content-type "text/plain" string "My content" hdr X-value x-value hdr Y-value y-value if { var(txn.myip) -m found }
//test:ok:http-request return status 400 default-errorfiles if { var(txn.myip) -m found }
//test:ok:http-request return status 400 errorfile /my/fancy/errorfile if { var(txn.myip) -m found }
//test:ok:http-request return status 400 errorfiles myerror if { var(txn.myip) -m found }
//test:quote_ok:http-request return content-type "text/plain" lf-string "Hello, you are: %[src]"
//test:fail:http-request return 8 t hdr
//test:fail:http-request return hdr
//test:fail:http-request return hdr one
//test:fail:http-request return errorfile
//test:fail:http-request return 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 file
//test:fail:http-request return 0 hdr
//test:fail:http-request return 0 0 hdr 0
//test:fail:http-request return e r s n s c m	t e r  s c t e s t e r s c v e hdr ï
//test:quote_ok:http-request redirect location /file.html if { var(txn.routecookie) "ROUTEMP" }:1
//test:ok:http-request redirect location /file.html if { var(txn.routecookie) -m found } !{ var(txn.pod),nbsrv -m found }:1]
//test:fail:http-request redirect location if { var(txn.routecookie) -m found } !{ var(txn.pod),nbsrv -m found }:1]
//test:fail:http-request redirect location /file.html code if { var(txn.routecookie) -m found } !{ var(txn.pod),nbsrv -m found }:1]
//test:ok:http-request set-bandwidth-limit my-limit
//test:ok:http-request set-bandwidth-limit my-limit limit 1m period 10s
//test:ok:http-request set-bandwidth-limit my-limit period 10s
//test:ok:http-request set-bandwidth-limit my-limit limit 1m
//test:fail:http-request set-bandwidth-limit my-limit limit
//test:fail:http-request set-bandwidth-limit my-limit period
//test:fail:http-request set-bandwidth-limit my-limit 10s
//test:fail:http-request set-bandwidth-limit my-limit period 10s limit
//test:fail:http-request set-bandwidth-limit my-limit limit period 10s
//test:ok:http-request set-bc-mark 123
//test:ok:http-request set-bc-mark 0xffffffff
//test:ok:http-request set-bc-mark hdr(port) if FALSE
//test:ok:http-request set-bc-tos 10
//test:ok:http-request set-fc-mark 0
//test:ok:http-request set-fc-tos 0xff if TRUE
type HTTPRequests struct{}

//name:http-response
//sections:frontend,backend,defaults
//struct:name:Responses
//dir:http
//is:multiple
//parser:type:Action
//is:interface
//no:init
//no:parse
//test:fail:http-response
//test:frontend-ok:http-response capture res.hdr(Server) id 0
//test:ok:http-response set-map(map.lst) %[src] %[res.hdr(X-Value)] if value
//test:ok:http-response set-map(map.lst) %[src] %[res.hdr(X-Value)]
//test:fail:http-response set-map(map.lst) %[src]
//test:ok:http-response add-acl(map.lst) [src]
//test:fail:http-response add-acl(map.lst)
//test:ok:http-response add-header X-value value
//test:fail:http-response add-header X-value
//test:ok:http-response del-acl(map.lst) [src]
//test:fail:http-response del-acl(map.lst)
//test:ok:http-response allow
//test:ok:http-response cache-store cache-name
//test:ok:http-response cache-store cache-name if FALSE
//test:fail:http-response cache-store
//test:fail:http-response cache-store if FALSE
//test:ok:http-response del-header X-value
//test:fail:http-response del-header
//test:ok:http-response del-map(map.lst) %[src] if ! value
//test:ok:http-response del-map(map.lst) %[src]
//test:fail:http-response del-map(map.lst)
//test:ok:http-response deny
//test:ok:http-response deny deny_status 400
//test:ok:http-response deny if TRUE
//test:ok:http-response deny deny_status 400 if TRUE
//test:ok:http-response deny deny_status 400 content-type application/json if TRUE
//test:ok:http-response deny deny_status 400 content-type application/json
//test:ok:http-response deny deny_status 400 content-type application/json default-errorfiles
//test:ok:http-response deny deny_status 400 content-type application/json errorfile errors
//test:ok:http-response deny deny_status 400 content-type application/json string error if TRUE
//test:ok:http-response deny deny_status 400 content-type application/json lf-string error hdr host google.com if TRUE
//test:ok:http-response deny deny_status 400 content-type application/json file /var/errors.file
//test:ok:http-response deny deny_status 400 content-type application/json lf-file /var/errors.file
//test:ok:http-response deny deny_status 400 content-type application/json string error hdr host google.com if TRUE
//test:ok:http-response deny deny_status 400 content-type application/json string error hdr host google.com hdr x-value bla if TRUE
//test:ok:http-response deny deny_status 400 content-type application/json string error hdr host google.com hdr x-value bla
//test:fail:http-response deny test test
//test:ok:http-response lua.foo
//test:ok:http-response lua.foo if FALSE
//test:ok:http-response lua.foo param
//test:ok:http-response lua.foo param param2
//test:fail:http-response lua.
//test:fail:http-response lua. if FALSE
//test:fail:http-response lua. param
//test:ok:http-response redirect prefix https://mysite.com
//test:fail:http-response redirect prefix
//test:ok:http-response replace-header User-agent curl foo
//test:fail:http-response replace-header User-agent curl
//test:ok:http-response replace-value X-Forwarded-For ^192.168.(.*)$ 172.16.1
//test:fail:http-response replace-value X-Forwarded-For ^192.168.(.*)$
//test:quote_ok:http-response return status 200 content-type "text/plain" string "My content" if { var(txn.myip) -m found }
//test:quote_ok:http-response return status 200 content-type "text/plain" string "My content" unless { var(txn.myip) -m found }
//test:quote_ok:http-response return content-type "text/plain" string "My content" if { var(txn.myip) -m found }
//test:quote_ok:http-response return content-type 'text/plain' string 'My content' if { var(txn.myip) -m found }
//test:quote_ok:http-response return content-type "text/plain" lf-string "Hello, you are: %[src]" if { var(txn.myip) -m found }
//test:quote_ok:http-response return content-type "text/plain" file /my/fancy/response/file if { var(txn.myip) -m found }
//test:quote_ok:http-response return content-type "text/plain" lf-file /my/fancy/lof/format/response/file if { var(txn.myip) -m found }
//test:quote_ok:http-response return content-type "text/plain" string "My content" hdr X-value value if { var(txn.myip) -m found }
//test:quote_ok:http-response return content-type "text/plain" string "My content" hdr X-value x-value hdr Y-value y-value if { var(txn.myip) -m found }
//test:ok:http-response return status 400 default-errorfiles if { var(txn.myip) -m found }
//test:ok:http-response return status 400 errorfile /my/fancy/errorfile if { var(txn.myip) -m found }
//test:ok:http-response return status 400 errorfiles myerror if { var(txn.myip) -m found }
//test:quote_ok:http-response return content-type "text/plain" lf-string "Hello, you are: %[src]"
//test:fail:http-response return 8 t hdr
//test:fail:http-response return hdr
//test:fail:http-response return hdr one
//test:fail:http-response return errorfile
//test:fail:http-response return 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 file
//test:fail:http-response return 0 hdr
//test:fail:http-response return 0 0 hdr 0
//test:fail:http-response return e r s n s c m	t e r  s c t e s t e r s c v e hdr ï
//test:ok:http-response sc-add-gpc(1,2) 1
//test:ok:http-response sc-add-gpc(1,2) 1 if is-error
//test:fail:http-response sc-add-gpc
//test:ok:http-response sc-inc-gpc(1,2)
//test:ok:http-response sc-inc-gpc(1,2) if FALSE
//test:fail:http-response sc-inc-gpc
//test:ok:http-response sc-inc-gpc0(1)
//test:ok:http-response sc-inc-gpc0(1) if FALSE
//test:fail:http-response sc-inc-gpc0
//test:ok:http-response sc-inc-gpc1(1)
//test:ok:http-response sc-inc-gpc1(1) if FALSE
//test:fail:http-response sc-inc-gpc1
//test:ok:http-response sc-set-gpt(1,2) hdr(Host),lower if FALSE
//test:ok:http-response sc-set-gpt0(1) hdr(Host),lower
//test:ok:http-response sc-set-gpt0(1) 10
//test:ok:http-response sc-set-gpt0(1) hdr(Host),lower if FALSE
//test:fail:http-response sc-set-gpt0(1)
//test:fail:http-response sc-set-gpt0
//test:fail:http-response sc-set-gpt0(1) if FALSE
//test:ok:http-response send-spoe-group engine group
//test:fail:http-response send-spoe-group engine
//test:ok:http-response set-header X-value value
//test:fail:http-response set-header X-value
//test:ok:http-response set-log-level silent
//test:fail:http-response set-log-level
//test:ok:http-response set-mark 20
//test:ok:http-response set-mark 0x1Ab
//test:fail:http-response set-mark
//test:ok:http-response set-nice 0
//test:ok:http-response set-nice 0 if FALSE
//test:fail:http-response set-nice
//test:ok:http-response set-status 503
//test:fail:http-response set-status
//test:ok:http-response set-timeout server 20
//test:ok:http-response set-timeout tunnel 20
//test:ok:http-response set-timeout tunnel 20s if TRUE
//test:ok:http-response set-timeout server 20s if TRUE
//test:ok:http-response set-timeout client 20
//test:ok:http-response set-timeout client 20s if TRUE
//test:fail:http-response set-timeout fake-timeout 20s if TRUE
//test:ok:http-response set-tos 0 if FALSE
//test:ok:http-response set-tos 0
//test:fail:http-response set-tos
//test:ok:http-response set-var(req.my_var) res.fhdr(user-agent),lower
//test:fail:http-response set-var(req.my_var)
//test:ok:http-response set-var-fmt(req.my_var) res.fhdr(user-agent),lower
//test:fail:http-response set-var-fmt(req.my_var)
//test:ok:http-response silent-drop
//test:ok:http-response silent-drop if FALSE
//test:ok:http-response unset-var(req.my_var)
//test:ok:http-response unset-var(req.my_var) if FALSE
//test:fail:http-response unset-var(req.)
//test:fail:http-response unset-var(req)
//test:ok:http-response track-sc0 src if FALSE
//test:ok:http-response track-sc0 src table tr if FALSE
//test:ok:http-response track-sc0 src
//test:fail:http-response track-sc0
//test:ok:http-response track-sc1 src if FALSE
//test:ok:http-response track-sc1 src table tr if FALSE
//test:ok:http-response track-sc1 src
//test:fail:http-response track-sc1
//test:ok:http-response track-sc2 src if FALSE
//test:ok:http-response track-sc2 src table tr if FALSE
//test:ok:http-response track-sc2 src
//test:fail:http-response track-sc2
//test:ok:http-response track-sc5 src
//test:ok:http-response track-sc5 src table a_table
//test:ok:http-response track-sc5 src table a_table if some_cond
//test:ok:http-response track-sc5 src if some_cond
//test:fail:http-response track-sc
//test:fail:http-response track-sc5
//test:fail:http-response track-sc5 src table
//test:fail:http-response track-sc5 src if
//test:ok:http-response strict-mode on
//test:ok:http-response strict-mode on if FALSE
//test:fail:http-response strict-mode
//test:fail:http-response strict-mode if FALSE
//test:ok:http-response wait-for-body time 20s
//test:ok:http-response wait-for-body time 20s if TRUE
//test:ok:http-response wait-for-body time 20s at-least 100k
//test:ok:http-response wait-for-body time 20s at-least 100k if TRUE
//test:fail:http-response wait-for-body 20s at-least 100k
//test:fail:http-response wait-for-body time 2000 test
//test:ok:http-response set-bandwidth-limit my-limit
//test:ok:http-response set-bandwidth-limit my-limit limit 1m period 10s
//test:ok:http-response set-bandwidth-limit my-limit period 10s
//test:ok:http-response set-bandwidth-limit my-limit limit 1m
//test:fail:http-response set-bandwidth-limit my-limit limit
//test:fail:http-response set-bandwidth-limit my-limit period
//test:fail:http-response set-bandwidth-limit my-limit 10s
//test:fail:http-response set-bandwidth-limit my-limit period 10s limit
//test:fail:http-response set-bandwidth-limit my-limit limit period 10s
//test:ok:http-response set-fc-mark 2000
//test:ok:http-response set-fc-tos 200
type HTTPResponses struct{}

//name:http-after-response
//sections:frontend,backend,defaults
//struct:name:AfterResponses
//dir:http
//is:multiple
//parser:type:Action
//is:interface
//no:parse
//no:init
//test:fail:http-after-response
//test:fail:http-after-response set-header
//test:fail:http-after-response set-header x-foo
//test:ok:http-after-response allow
//test:ok:http-after-response allow if acl
//test:ok:http-after-response set-header Strict-Transport-Security \"max-age=31536000\"
//test:ok:http-after-response add-header X-Header \"foo=bar\"
//test:ok:http-after-response add-header X-Header \"foo=bar\" if acl
//test:ok:http-after-response add-header X-Header \"foo=bar\" unless acl
//test:ok:http-after-response allow unless acl
//test:ok:http-after-response del-header X-Value
//test:ok:http-after-response del-header X-Value -m GET
//test:ok:http-after-response del-header X-Value -m GET if acl
//test:ok:http-after-response del-header X-Value -m GET unless acl
//test:fail:http-after-response del-header
//test:ok:http-after-response replace-header Set-Cookie (C=[^;]*);(.*) \\1;ip=%bi;\\2
//test:ok:http-after-response replace-header Set-Cookie (C=[^;]*);(.*) \\1;ip=%bi;\\2 if acl
//test:fail:http-after-response replace-header Set-Cookie
//test:fail:http-after-response replace-header Set-Cookie (C=[^;]*);(.*)
//test:ok:http-after-response replace-value Cache-control ^public$ private
//test:ok:http-after-response replace-value Cache-control ^public$ private if acl
//test:fail:http-after-response replace-value Cache-control
//test:fail:http-after-response replace-value Cache-control ^public$
//test:ok:http-after-response set-status 431
//test:ok:http-after-response set-status 503 reason \"SlowDown\"
//test:ok:http-after-response set-status 500 reason \"ServiceUnavailable\" if acl
//test:ok:http-after-response set-status 500 reason \"ServiceUnavailable\" unless acl
//test:fail:http-after-response set-status
//test:fail:http-after-response set-status error
//test:ok:http-after-response set-var(sess.last_redir) res.hdr(location)
//test:ok:http-after-response set-var(sess.last_redir) res.hdr(location) if acl
//test:ok:http-after-response set-var(sess.last_redir) res.hdr(location) unless acl
//test:fail:http-after-response set-var(sess.last_redir)
//test:ok:http-after-response strict-mode on
//test:ok:http-after-response strict-mode off
//test:fail:http-after-response strict-mode
//test:fail:http-after-response strict-mode 1
//test:fail:http-after-response strict-mode 0
//test:ok:http-after-response unset-var(sess.last_redir)
//test:ok:http-after-response unset-var(sess.last_redir) if acl
//test:ok:http-after-response unset-var(sess.last_redir) unless acl
//test:fail:http-after-response unset-var()
//test:ok:http-after-response set-map(map.lst) %[src] %[res.hdr(X-Value)] if value
//test:ok:http-after-response set-map(map.lst) %[src] %[res.hdr(X-Value)]
//test:fail:http-after-response set-map(map.lst) %[src]
//test:ok:http-after-response del-acl(map.lst) [src]
//test:fail:http-after-response del-acl(map.lst)
//test:ok:http-after-response del-map(map.lst) %[src] if ! value
//test:ok:http-after-response del-map(map.lst) %[src]
//test:fail:http-after-response del-map(map.lst)
//test:ok:http-after-response sc-add-gpc(1,2) 1
//test:ok:http-after-response sc-add-gpc(1,2) 1 if is-error
//test:fail:http-after-response sc-add-gpc
//test:ok:http-after-response sc-inc-gpc(1,2)
//test:ok:http-after-response sc-inc-gpc(1,2) if is-error
//test:fail:http-after-response sc-inc-gpc
//test:ok:http-after-response sc-inc-gpc0(1)
//test:ok:http-after-response sc-inc-gpc0(1) if FALSE
//test:fail:http-after-response sc-inc-gpc0
//test:ok:http-after-response sc-inc-gpc1(1)
//test:ok:http-after-response sc-inc-gpc1(1) if FALSE
//test:fail:http-after-response sc-inc-gpc1
//test:ok:http-after-response sc-set-gpt(1,2) 10
//test:ok:http-after-response sc-set-gpt0(1) hdr(Host),lower
//test:ok:http-after-response sc-set-gpt0(1) 10
//test:ok:http-after-response sc-set-gpt0(1) hdr(Host),lower if FALSE
//test:fail:http-after-response sc-set-gpt0(1)
//test:fail:http-after-response sc-set-gpt0
//test:fail:http-after-response sc-set-gpt0(1) if FALSE
type HTTPAfterResponse struct{}

//name:http-error
//sections:defaults,frontend,backend
//struct:name:HTTPErrors
//dir:http
//is:multiple
//parser:type:Action
//is:interface
//no:init
//no:parse
//test:fail:http-error
//test:ok:http-error status 400
//test:fail:http-error status 402
//test:fail:http-error status
//test:quote_ok:http-error status 200 content-type "text/plain" string "My content"
//test:fail:http-error status 200 content-type 'text/plain' string 'My content' if { var(txn.myip) -m found }
//test:quote_ok:http-error status 400 content-type "text/plain" lf-string "Hello, you are: %[src]"
//test:quote_ok:http-error status 400 content-type "text/plain" file /my/fancy/response/file
//test:quote_ok:http-error status 400 content-type "text/plain" lf-file /my/fancy/lof/format/response/file
//test:quote_ok:http-error status 400 content-type "text/plain" string "My content" hdr X-value value
//test:quote_ok:http-error status 400 content-type "text/plain" string "My content" hdr X-value x-value hdr Y-value y-value
//test:ok:http-error status 400 default-errorfiles
//test:ok:http-error status 400 errorfile /my/fancy/errorfile
//test:ok:http-error status 400 errorfiles myerror
type HTTPErrors struct{}

//name:http-check
//sections:defaults,backend
//struct:name:Checks
//dir:http
//is:multiple
//parser:type:Action
//is:interface
//no:init
//no:parse
//test:ok:http-check comment testcomment
//test:ok:http-check connect
//test:ok:http-check connect default
//test:ok:http-check connect port 8080
//test:ok:http-check connect addr 8.8.8.8
//test:ok:http-check connect send-proxy
//test:ok:http-check connect via-socks4
//test:ok:http-check connect ssl
//test:ok:http-check connect sni haproxy.1wt.eu
//test:ok:http-check connect alpn h2,http/1.1
//test:ok:http-check connect proto h2
//test:ok:http-check connect linger
//test:ok:http-check connect comment testcomment
//test:ok:http-check connect port 443 addr 8.8.8.8 send-proxy via-socks4 ssl sni haproxy.1wt.eu alpn h2,http/1.1 linger proto h2 comment testcomment
//test:ok:http-check disable-on-404
//test:ok:http-check expect status 200
//test:ok:http-check expect min-recv 50 status 200
//test:ok:http-check expect comment testcomment status 200
//test:ok:http-check expect ok-status L7OK status 200
//test:ok:http-check expect error-status L7RSP status 200
//test:ok:http-check expect tout-status L7TOUT status 200
//test:ok:http-check expect on-success \"my-log-format\" status 200
//test:ok:http-check expect on-error \"my-log-format\" status 200
//test:ok:http-check expect status-code \"500\" status 200
//test:ok:http-check expect ! string SQL\\ Error
//test:ok:http-check expect ! rstatus ^5
//test:ok:http-check expect rstring <!--tag:[0-9a-f]*--></html>
//test:ok:http-check send meth GET
//test:ok:http-check send uri /health
//test:ok:http-check send ver \"HTTP/1.1\"
//test:ok:http-check send comment testcomment
//test:ok:http-check send meth GET uri /health ver \"HTTP/1.1\" hdr Host example.com hdr Accept-Encoding gzip body '{\"key\":\"value\"}'
//test:ok:http-check send uri-lf my-log-format body-lf 'my-log-format'
//test:ok:http-check send-state
//test:fail:http-check
//test:fail:http-check comment
//test:fail:http-check expect
//test:fail:http-check expect status
//test:fail:http-check expect comment testcomment
//test:fail:http-check set-var(check.port)
//test:quote_ok:http-check set-var(check.port) int(1234)
//test:fail:http-check set-var(check.port) int(1234) if x
//test:fail:http-check set-var-fmt(check.port)
//test:quote_ok:http-check set-var-fmt(check.port) int(1234)
//test:fail:http-check set-var-fmt(check.port) int(1234) if x
//test:quote_ok:http-check unset-var(txn.from)
//test:fail:http-check unset-var(txn.from) if x
type HTTPCheck struct{}

//name:tcp-check
//sections:defaults,backend
//struct:name:Checks
//dir:tcp
//is:multiple
//parser:type:Action
//is:interface
//no:init
//no:parse
//test:ok:tcp-check comment testcomment
//test:ok:tcp-check connect
//test:ok:tcp-check connect port 443 ssl
//test:ok:tcp-check connect port 110 linger
//test:ok:tcp-check connect port 143
//test:quote_ok:tcp-check expect string +OK\ POP3\ ready
//test:quote_ok:tcp-check expect string *\ OK\ IMAP4\ ready
//test:ok:tcp-check expect string +PONG
//test:ok:tcp-check expect string role:master
//test:ok:tcp-check expect string +OK
//test:quote_ok:tcp-check send PING\r\n
//test:quote_ok:tcp-check send PING\r\n comment testcomment
//test:quote_ok:tcp-check send QUIT\r\n
//test:quote_ok:tcp-check send QUIT\r\n comment testcomment
//test:quote_ok:tcp-check send info\ replication\r\n
//test:ok:tcp-check send-lf testfmt
//test:ok:tcp-check send-lf testfmt comment testcomment
//test:ok:tcp-check send-binary testhexstring
//test:ok:tcp-check send-binary testhexstring comment testcomment
//test:ok:tcp-check send-binary-lf testhexfmt
//test:ok:tcp-check send-binary-lf testhexfmt comment testcomment
//test:fail:tcp-check set-var(check.port)
//test:ok:tcp-check set-var(check.port) int(1234)
//test:fail:tcp-check set-var(check.port) int(1234) if x
//test:quote_ok:tcp-check set-var-fmt(check.name) "%H"
//test:quote_ok:tcp-check set-var-fmt(txn.from) "addr=%[src]:%[src_port]"
//test:quote_fail:tcp-check set-var-fmt(txn.from) "addr=%[src]:%[src_port]" if TRUE
//test:quote_ok:tcp-check unset-var(txn.from)
//test:fail:tcp-check unset-var(txn.from) if x
type TCPCheck struct{}

type TCPType interface {
	Parse(parts []string, comment string) error
	String() string
	GetComment() string
}

//name:tcp-request
//sections:frontend,backend,defaults
//struct:name:Requests
//dir:tcp
//is:multiple
//parser:type:TCPType
//is:interface
//no:init
//no:parse
//test:ok:tcp-request content accept
//test:ok:tcp-request content accept if !HTTP
//test:ok:tcp-request content reject
//test:ok:tcp-request content reject if !HTTP
//test:ok:tcp-request content capture req.payload(0,6) len 6
//test:ok:tcp-request content capture req.payload(0,6) len 6 if !HTTP
//test:ok:tcp-request content do-resolve(txn.myip,mydns,ipv6) capture.req.hdr(0),lower
//test:ok:tcp-request content do-resolve(txn.myip,mydns) capture.req.hdr(0),lower
//test:ok:tcp-request content set-priority-class int(1)
//test:ok:tcp-request content set-priority-class int(1) if some_check
//test:ok:tcp-request content set-priority-offset int(10)
//test:ok:tcp-request content set-priority-offset int(10) if some_check
//test:ok:tcp-request content track-sc0 src
//test:ok:tcp-request content track-sc0 src if some_check
//test:ok:tcp-request content track-sc1 src
//test:ok:tcp-request content track-sc1 src if some_check
//test:ok:tcp-request content track-sc2 src
//test:ok:tcp-request content track-sc2 src if some_check
//test:ok:tcp-request content track-sc0 src table foo
//test:ok:tcp-request content track-sc0 src table foo if some_check
//test:ok:tcp-request content track-sc1 src table foo
//test:ok:tcp-request content track-sc1 src table foo if some_check
//test:ok:tcp-request content track-sc2 src table foo
//test:ok:tcp-request content track-sc2 src table foo if some_check
//test:ok:tcp-request content track-sc5 src
//test:ok:tcp-request content track-sc5 src if some_check
//test:ok:tcp-request content track-sc5 src table foo
//test:ok:tcp-request content track-sc5 src table foo if some_check
//test:ok:tcp-request content sc-inc-gpc(1,2)
//test:ok:tcp-request content sc-inc-gpc(1,2) if is-error
//test:fail:tcp-request content sc-inc-gpc
//test:ok:tcp-request content sc-inc-gpc0(2)
//test:ok:tcp-request content sc-inc-gpc0(2) if is-error
//test:fail:tcp-request content sc-inc-gpc0
//test:ok:tcp-request content sc-inc-gpc1(2)
//test:ok:tcp-request content sc-inc-gpc1(2) if is-error
//test:fail:tcp-request content sc-inc-gpc1
//test:ok:tcp-request content sc-set-gpt(x,9) 1337 if exceeds_limit
//test:ok:tcp-request content sc-set-gpt0(0) 1337
//test:ok:tcp-request content sc-set-gpt0(0) 1337 if exceeds_limit
//test:ok:tcp-request content sc-add-gpc(1,2) 1
//test:ok:tcp-request content sc-add-gpc(1,2) 1 if is-error
//test:fail:tcp-request content sc-add-gpc
//test:ok:tcp-request content set-dst ipv4(10.0.0.1)
//test:ok:tcp-request content set-var(sess.src) src
//test:ok:tcp-request content set-var(sess.dn) ssl_c_s_dn
//test:ok:tcp-request content set-var-fmt(sess.src) src
//test:ok:tcp-request content set-var-fmt(sess.dn) ssl_c_s_dn
//test:ok:tcp-request content unset-var(sess.src)
//test:ok:tcp-request content unset-var(sess.dn)
//test:ok:tcp-request content silent-drop
//test:ok:tcp-request content silent-drop if !HTTP
//test:ok:tcp-request content send-spoe-group engine group
//test:ok:tcp-request content use-service lua.deny
//test:ok:tcp-request content use-service lua.deny if !HTTP
//test:ok:tcp-request content lua.foo
//test:ok:tcp-request content lua.foo param if !HTTP
//test:ok:tcp-request content lua.foo param param1
//test:ok:tcp-request connection accept
//test:ok:tcp-request connection accept if !HTTP
//test:ok:tcp-request connection reject
//test:ok:tcp-request connection reject if !HTTP
//test:ok:tcp-request connection expect-proxy layer4 if { src -f proxies.lst }
//test:ok:tcp-request connection expect-netscaler-cip layer4
//test:ok:tcp-request connection expect-netscaler-cip layer4 if TRUE
//test:ok:tcp-request connection capture req.payload(0,6) len 6
//test:ok:tcp-request connection track-sc0 src
//test:ok:tcp-request connection track-sc0 src if some_check
//test:ok:tcp-request connection track-sc1 src
//test:ok:tcp-request connection track-sc1 src if some_check
//test:ok:tcp-request connection track-sc2 src
//test:ok:tcp-request connection track-sc2 src if some_check
//test:ok:tcp-request connection track-sc0 src table foo
//test:ok:tcp-request connection track-sc0 src table foo if some_check
//test:ok:tcp-request connection track-sc1 src table foo
//test:ok:tcp-request connection track-sc1 src table foo if some_check
//test:ok:tcp-request connection track-sc2 src table foo
//test:ok:tcp-request connection track-sc2 src table foo if some_check
//test:ok:tcp-request connection track-sc5 src
//test:ok:tcp-request connection track-sc5 src if some_check
//test:ok:tcp-request connection track-sc5 src table foo
//test:ok:tcp-request connection track-sc5 src table foo if some_check
//test:ok:tcp-request connection sc-add-gpc(1,2) 1
//test:ok:tcp-request connection sc-add-gpc(1,2) 1 if is-error
//test:fail:tcp-request connection sc-add-gpc
//test:ok:tcp-request connection sc-inc-gpc(1,2)
//test:ok:tcp-request connection sc-inc-gpc(1,2) if is-error
//test:fail:tcp-request connection sc-inc-gpc
//test:ok:tcp-request connection sc-inc-gpc0(2)
//test:ok:tcp-request connection sc-inc-gpc0(2) if is-error
//test:fail:tcp-request connection sc-inc-gpc0
//test:ok:tcp-request connection sc-inc-gpc1(2)
//test:ok:tcp-request connection sc-inc-gpc1(2) if is-error
//test:fail:tcp-request connection sc-inc-gpc1
//test:ok:tcp-request connection sc-set-gpt(scx,44) 1337 if exceeds_limit
//test:ok:tcp-request connection sc-set-gpt0(0) 1337
//test:ok:tcp-request connection sc-set-gpt0(0) 1337 if exceeds_limit
//test:ok:tcp-request connection set-src src,ipmask(24)
//test:ok:tcp-request connection set-src src,ipmask(24) if some_check
//test:ok:tcp-request connection set-src hdr(x-forwarded-for)
//test:ok:tcp-request connection set-src hdr(x-forwarded-for) if some_check
//test:fail:tcp-request connection set-src
//test:ok:tcp-request connection silent-drop
//test:ok:tcp-request connection silent-drop if !HTTP
//test:ok:tcp-request connection lua.foo
//test:ok:tcp-request connection lua.foo param if !HTTP
//test:ok:tcp-request connection lua.foo param param1
//test:ok:tcp-request session accept
//test:ok:tcp-request session accept if !HTTP
//test:ok:tcp-request session reject
//test:ok:tcp-request session reject if !HTTP
//test:ok:tcp-request session track-sc0 src
//test:ok:tcp-request session track-sc0 src if some_check
//test:ok:tcp-request session track-sc1 src
//test:ok:tcp-request session track-sc1 src if some_check
//test:ok:tcp-request session track-sc2 src
//test:ok:tcp-request session track-sc2 src if some_check
//test:ok:tcp-request session track-sc0 src table foo
//test:ok:tcp-request session track-sc0 src table foo if some_check
//test:ok:tcp-request session track-sc1 src table foo
//test:ok:tcp-request session track-sc1 src table foo if some_check
//test:ok:tcp-request session track-sc2 src table foo
//test:ok:tcp-request session track-sc2 src table foo if some_check
//test:ok:tcp-request session track-sc5 src
//test:ok:tcp-request session track-sc5 src if some_check
//test:ok:tcp-request session track-sc5 src table foo
//test:ok:tcp-request session track-sc5 src table foo if some_check
//test:ok:tcp-request session sc-add-gpc(1,2) 1
//test:ok:tcp-request session sc-add-gpc(1,2) 1 if is-error
//test:fail:tcp-request session sc-add-gpc
//test:ok:tcp-request session sc-inc-gpc(1,2)
//test:ok:tcp-request session sc-inc-gpc(1,2) if is-error
//test:fail:tcp-request session sc-inc-gpc
//test:ok:tcp-request session sc-inc-gpc0(2)
//test:ok:tcp-request session sc-inc-gpc0(2) if is-error
//test:fail:tcp-request session sc-inc-gpc0
//test:ok:tcp-request session sc-inc-gpc1(2)
//test:fail:tcp-request session sc-inc-gpc1
//test:ok:tcp-request session sc-inc-gpc1(2) if is-error
//test:ok:tcp-request session sc-set-gpt(sc5,1) 1337 if exceeds_limit
//test:ok:tcp-request session sc-set-gpt0(0) 1337
//test:ok:tcp-request session sc-set-gpt0(0) 1337 if exceeds_limit
//test:ok:tcp-request session set-var(sess.src) src
//test:ok:tcp-request session set-var(sess.dn) ssl_c_s_dn
//test:ok:tcp-request session set-var-fmt(sess.src) src
//test:ok:tcp-request session set-var-fmt(sess.dn) ssl_c_s_dn
//test:ok:tcp-request session unset-var(sess.src)
//test:ok:tcp-request session unset-var(sess.dn)
//test:ok:tcp-request session silent-drop
//test:ok:tcp-request session silent-drop if !HTTP
//test:ok:tcp-request session attach-srv srv1
//test:ok:tcp-request session attach-srv srv1 name example.com
//test:ok:tcp-request session attach-srv srv1 name example.com if exceeds_limit
//test:fail:tcp-request session attach-srv
//test:fail:tcp-request session attach-srv srv1 name
//test:fail:tcp-request session attach-srv srv1 if
//test:fail:tcp-request session attach-srv srv1 name example.com unless
//test:ok:tcp-request content set-bandwidth-limit my-limit
//test:ok:tcp-request content set-bandwidth-limit my-limit limit 1m period 10s
//test:ok:tcp-request content set-bandwidth-limit my-limit period 10s
//test:ok:tcp-request content set-bandwidth-limit my-limit limit 1m
//test:ok:tcp-request content set-log-level silent
//test:ok:tcp-request content set-log-level silent if FALSE
//test:ok:tcp-request content set-mark 20
//test:ok:tcp-request content set-mark 0x1Ab if FALSE
//test:ok:tcp-request connection set-mark 20
//test:ok:tcp-request connection set-mark 0x1Ab if FALSE
//test:ok:tcp-request connection set-src-port hdr(port)
//test:ok:tcp-request connection set-src-port hdr(port) if FALSE
//test:ok:tcp-request content set-src-port hdr(port)
//test:ok:tcp-request content set-src-port hdr(port) if FALSE
//test:ok:tcp-request content set-tos 0 if FALSE
//test:ok:tcp-request content set-tos 0
//test:ok:tcp-request connection set-tos 0 if FALSE
//test:ok:tcp-request connection set-tos 0
//test:ok:tcp-request connection set-var-fmt(txn.ip_port) %%[dst]:%%[dst_port]
//test:ok:tcp-request content set-nice 0 if FALSE
//test:ok:tcp-request content set-nice 0
//test:ok:tcp-request content switch-mode http
//test:ok:tcp-request content switch-mode http if FALSE
//test:ok:tcp-request content switch-mode http proto my-proto
//test:fail:tcp-request
//test:fail:tcp-request content
//test:fail:tcp-request connection
//test:fail:tcp-request content do-resolve(txn.myip) capture.req.hdr(0),lower
//test:fail:tcp-request session
//test:fail:tcp-request content lua.
//test:fail:tcp-request content lua. param
//test:fail:tcp-request connection lua.
//test:fail:tcp-request connection lua. param
//test:fail:tcp-request content track-sc0 src table
//test:fail:tcp-request content track-sc0 src table if some_check
//test:fail:tcp-request content track-sc1 src table
//test:fail:tcp-request content track-sc1 src table if some_check
//test:fail:tcp-request content track-sc2 src table
//test:fail:tcp-request content track-sc2 src table if some_check
//test:fail:tcp-request connection track-sc0 src table
//test:fail:tcp-request connection track-sc0 src table if some_check
//test:fail:tcp-request connection track-sc1 src table
//test:fail:tcp-request connection track-sc1 src table if some_check
//test:fail:tcp-request connection track-sc2 src table
//test:fail:tcp-request connection track-sc2 src table if some_check
//test:fail:tcp-request session track-sc0 src table
//test:fail:tcp-request session track-sc0 src table if some_check
//test:fail:tcp-request session track-sc1 src table
//test:fail:tcp-request session track-sc1 src table if some_check
//test:fail:tcp-request session track-sc2 src table
//test:fail:tcp-request session track-sc2 src table if some_check
//test:fail:tcp-request content track-sc5 src table
//test:fail:tcp-request content track-sc5 src table if some_check
//test:fail:tcp-request connection track-sc5 src table
//test:fail:tcp-request connection track-sc5 src table if some_check
//test:fail:tcp-request session track-sc5 src table
//test:fail:tcp-request session track-sc5 src table if some_check
//test:fail:tcp-request content set-bandwidth-limit my-limit limit
//test:fail:tcp-request content set-bandwidth-limit my-limit period
//test:fail:tcp-request content set-bandwidth-limit my-limit 10s
//test:fail:tcp-request content set-bandwidth-limit my-limit period 10s limit
//test:fail:tcp-request content set-bandwidth-limit my-limit limit period 10s
//test:fail:tcp-request content set-log-level
//test:fail:tcp-request connection set-mark
//test:fail:tcp-request content set-mark
//test:fail:tcp-request connection set-src-port
//test:fail:tcp-request content set-src-port
//test:fail:tcp-request connection set-tos
//test:fail:tcp-request content set-tos
//test:fail:tcp-request content set-nice
//test:fail:tcp-request content switch-mode
//test:fail:tcp-request content switch-mode tcp
//test:fail:tcp-request content switch-mode http proto
//test:ok:tcp-request connection set-fc-mark 1
//test:ok:tcp-request connection set-fc-tos 1
//test:ok:tcp-request session set-fc-mark 9999 if some_check
//test:ok:tcp-request session set-fc-tos 255
//test:ok:tcp-request content set-bc-mark hdr(port)
//test:ok:tcp-request content set-bc-tos 0xff if some_check
//test:ok:tcp-request content set-fc-mark 0xffffffff
//test:ok:tcp-request content set-fc-tos 100
type TCPRequests struct{}

//name:tcp-response
//sections:frontend,backend
//struct:name:Responses
//dir:tcp
//is:multiple
//parser:type:TCPType
//is:interface
//no:init
//no:parse
//test:ok:tcp-response content lua.foo
//test:ok:tcp-response content lua.foo param if !HTTP
//test:ok:tcp-response content lua.foo param param1
//test:fail:tcp-response
//test:fail:tcp-response content lua.
//test:fail:tcp-response content lua. param
//test:fail:tcp-response content set-priority-class
//test:fail:tcp-response content do-resolve
//test:fail:tcp-response content set-priority-offset
//test:ok:tcp-response content set-dst dest
//test:fail:tcp-response content set-dst
//test:fail:tcp-response content capture
//test:ok:tcp-response content unset-var(sess.my_var)
//test:ok:tcp-response content set-bandwidth-limit my-limit
//test:ok:tcp-response content set-bandwidth-limit my-limit limit 1m period 10s
//test:ok:tcp-response content set-bandwidth-limit my-limit period 10s
//test:ok:tcp-response content set-bandwidth-limit my-limit limit 1m
//test:fail:tcp-response content set-bandwidth-limit my-limit limit
//test:fail:tcp-response content set-bandwidth-limit my-limit period
//test:fail:tcp-response content set-bandwidth-limit my-limit 10s
//test:fail:tcp-response content set-bandwidth-limit my-limit period 10s limit
//test:fail:tcp-response content set-bandwidth-limit my-limit limit period 10s
//test:ok:tcp-response content set-log-level silent
//test:ok:tcp-response content set-log-level silent if FALSE
//test:fail:tcp-response content set-log-level
//test:ok:tcp-response content set-mark 20
//test:ok:tcp-response content set-mark 0x1Ab if FALSE
//test:fail:tcp-response content set-mark
//test:ok:tcp-response content set-tos 0 if FALSE
//test:ok:tcp-response content set-tos 0
//test:fail:tcp-response content set-tos
//test:ok:tcp-response content set-nice 0 if FALSE
//test:ok:tcp-response content set-nice 0
//test:fail:tcp-response content set-nice
//test:ok:tcp-response content close
//test:ok:tcp-response content close if !HTTP
//test:ok:tcp-response content sc-inc-gpc(1,2)
//test:ok:tcp-response content sc-inc-gpc(1,2) if is-error
//test:fail:tcp-response content sc-inc-gpc
//test:ok:tcp-response content sc-inc-gpc0(2)
//test:ok:tcp-response content sc-inc-gpc0(2) if is-error
//test:ok:tcp-response content sc-inc-gpc1(2)
//test:ok:tcp-response content sc-inc-gpc1(2) if is-error
//test:ok:tcp-response content set-fc-mark 123456
//test:ok:tcp-response content set-fc-tos 0x02
type TCPResponses struct{}

//name:redirect
//sections:frontend,backend
//dir:http
//is:multiple
//parser:type:Action
//is:interface
//no:init
//no:parse
//test:fail:redirect
//test:ok:redirect prefix http://www.bar.com code 301 if { hdr(host) -i foo.com }
type Redirect struct{}

type StatsSettings interface {
	Parse(parts []string, comment string) error
	String() string
	GetComment() string
}

//name:stats
//sections:defaults,frontend,backend
//struct:name:Stats
//dir:stats
//is:multiple
//parser:type:StatsSettings
//is:interface
//no:init
//no:parse
//test:fail:stats
//test:frontend-ok:stats admin if LOCALHOST
//test:ok:stats auth admin1:AdMiN123
//test:fail:stats auth admin1:
//test:fail:stats auth
//test:ok:stats enable
//test:ok:stats hide-version
//test:ok:stats show-legends
//test:ok:stats show-modules
//test:fail:stats NON-EXISTS
//test:ok:stats maxconn 10
//test:fail:stats maxconn WORD
//test:ok:stats realm HAProxy\\ Statistics
//test:ok:stats refresh 10s
//test:fail:stats refresh
//test:ok:stats scope .
//test:fail:stats scope
//test:ok:stats show-desc Master node for Europe, Asia, Africa
//test:ok:stats show-node
//test:ok:stats show-node Europe-1
//test:ok:stats uri /admin?stats
//test:fail:stats uri
//test:ok:stats bind-process all
//test:ok:stats bind-process odd
//test:ok:stats bind-process even
//test:ok:stats bind-process 1 2 3 4
//test:ok:stats bind-process 1-4
//test:fail:stats bind-process none
//test:fail:stats bind-process 1+4
//test:fail:stats bind-process none-none
//test:fail:stats bind-process 1-4 1-3
//test:backend-ok:stats http-request auth realm HAProxy\\ Statistics
//test:backend-ok:stats http-request auth realm HAProxy\\ Statistics if something
//test:backend-ok:stats http-request auth if something
//test:backend-ok:stats http-request deny unless something
//test:backend-ok:stats http-request allow
//test:fail:stats http-request
//test:fail:stats http-request none
//test:fail:stats http-request realm HAProxy\\ Statistics
type Stats struct{}
