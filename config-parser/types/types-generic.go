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

//nolint:godot
package types

// Enabled is used by parsers Daemon, MasterWorker, ExternalCheck, NoSplice, CompressionOffload
//
//generate:type:Daemon
//name:daemon
//create:type:bool
//test:ok:daemon
//test:ok:daemon # comment
//generate:type:MasterWorker
//name:master-worker
//create:type:bool
//test:ok:master-worker
//test:ok:master-worker # comment
//generate:type:ExternalCheck
//name:external-check
//create:type:bool
//test:ok:external-check
//test:ok:external-check # comment
//generate:type:NoSplice
//name:nosplice
//create:type:bool
//test:ok:nosplice
//test:ok:nosplice # comment
//generate:type:CompressionOffload
//name:compression offload
//create:type:bool
//test:ok:compression offload
//test:ok:compression offload # comment
type Enabled struct {
	Comment string
}

// Int64 is used by parsers MaxConn, NbProc, NbThread
//
//generate:type:MaxConn
//name:maxconn
//test:ok:maxconn 10000
//test:ok:maxconn 10000 # comment
//test:fail:maxconn
//generate:type:NbProc
//name:nbproc
//test:ok:nbproc 4
//test:ok:nbproc 4 # comment
//test:fail:nbproc
//generate:type:NbThread
//name:nbthread
//test:ok:nbthread 4
//test:ok:nbthread 4 # comment
//test:fail:nbthread
//generate:type:StatsMaxconn
//name:stats maxconn
//test:ok:stats maxconn 10
//test:fail:stats
//test:fail:maxconn
//test:fail:stats maxconn
//test:fail:stats maxconn string
type Int64C struct {
	Value   int64
	Comment string
}

// String is used by parsers Mode, DefaultBackend, SimpleTimeTwoWords, StatsTimeout, CompressionDirection, CompressionAlgoReq
//
//generate:type:Mode
//name:mode
//test:ok:mode tcp
//test:ok:mode http
//test:ok:mode tcp # comment
//test:fail:mode
//generate:type:DefaultBackend
//name:default_backend
//test:ok:default_backend http
//test:fail:default_backend
//generate:type:StatsTimeout
//name:stats timeout
//test:ok:stats timeout 4
//test:ok:stats timeout 4 # comment
//test:fail:stats timeout
//test:fail:stats
//test:fail:timeout
//generate:type:LogSendHostName
//name:log-send-hostname
//test:ok:log-send-hostname
//test:ok:log-send-hostname something
//generate:type:CompressionDirection
//name:compression direction
//test:ok:compression direction both
//test:fail:compression direction
//generate:type:CompressionAlgoReq
//name:compression algo-req
//test:ok:compression algo-req gzip
//test:fail:compression algo-req
type StringC struct {
	Value   string
	Comment string
}

// StringSliceC is used by ConfigSnippet, CompressionAlgo, CompressionType, CompressionTypeReq, CompressionTypeRes, CompressionAlgoRes
//
//generate:type:ConfigSnippet
//name:config-snippet
//test:ok:###_config-snippet_### BEGIN\n  tune.ssl.default-dh-param 2048\n  tune.bufsize 32768\n  ###_config-snippet_### END
//test:fail:tune.ssl.default-dh-param 2048\ntune.bufsize 32768
//generate:type:CompressionAlgo
//name:compression algo
//test:ok:compression algo identity
//test:ok:compression algo identity raw-deflate
//test:fail:compression algo
//generate:type:CompressionType
//name:compression type
//test:ok:compression type text/plain
//test:ok:compression type text/plain application/json
//test:fail:compression type
//generate:type:CompressionTypeReq
//name:compression type-req
//test:ok:compression type-req text/plain
//test:ok:compression type-req text/plain application/json
//test:fail:compression type-req
//generate:type:CompressionTypeRes
//name:compression type-res
//test:ok:compression type-res text/plain
//test:ok:compression type-res text/plain application/json
//test:fail:compression type-res
//generate:type:CompressionAlgoRes
//name:compression algo-res
//test:ok:compression algo-res gzip raw-deflate
//test:fail:compression algo-res
type StringSliceC struct {
	Value   []string
	Comment string
}

// StringKeyValueC is a simple key value, for example environment variables.
type StringKeyValueC struct {
	Key     string
	Value   string
	Comment string
}

// Filters are not here, see parsers/filters
// ==============================================================================
