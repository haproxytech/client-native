// Code generated by go generate; DO NOT EDIT.
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

package integration_test

import (
	"bytes"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/options"
)

func TestWholeConfigsSectionsGlobal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		Name, Config string
	}{
		{"global_cpumap1403", global_cpumap1403},
		{"global_cpumap1all03", global_cpumap1all03},
		{"global_cpumapauto140123", global_cpumapauto140123},
		{"global_cpumapauto1403", global_cpumapauto1403},
		{"global_defaultpathconfig", global_defaultpathconfig},
		{"global_defaultpathcurrent", global_defaultpathcurrent},
		{"global_defaultpathcurrentcomment", global_defaultpathcurrentcomment},
		{"global_defaultpathoriginsomepath", global_defaultpathoriginsomepath},
		{"global_defaultpathoriginsomepathcomment", global_defaultpathoriginsomepathcomment},
		{"global_defaultpathparent", global_defaultpathparent},
		{"global_h1caseadjustcontenttypeContentTy", global_h1caseadjustcontenttypeContentTy},
		{"global_httpclientresolverspreferipv4", global_httpclientresolverspreferipv4},
		{"global_httpclientresolverspreferipv6", global_httpclientresolverspreferipv6},
		{"global_httpclientsslverify", global_httpclientsslverify},
		{"global_httpclientsslverifynone", global_httpclientsslverifynone},
		{"global_httpclientsslverifyrequired", global_httpclientsslverifyrequired},
		{"global_httperrcodes400402444446480490", global_httperrcodes400402444446480490},
		{"global_httperrcodes400408comment", global_httperrcodes400408comment},
		{"global_httperrcodes400499450500", global_httperrcodes400499450500},
		{"global_httpfailcodes400402444446480490", global_httpfailcodes400402444446480490},
		{"global_httpfailcodes400408comment", global_httpfailcodes400408comment},
		{"global_httpfailcodes400499450500", global_httpfailcodes400499450500},
		{"global_lualoadetchaproxyluafoolua", global_lualoadetchaproxyluafoolua},
		{"global_luaprependpathusrsharehaproxylua", global_luaprependpathusrsharehaproxylua},
		{"global_luaprependpathusrsharehaproxylua_", global_luaprependpathusrsharehaproxylua_},
		{"global_nonumacpumapping", global_nonumacpumapping},
		{"global_numacpumapping", global_numacpumapping},
		{"global_setvarfmtprocbootidpidt", global_setvarfmtprocbootidpidt},
		{"global_setvarfmtproccurrentstateprimary", global_setvarfmtproccurrentstateprimary},
		{"global_setvarproccurrentstatestrprimary", global_setvarproccurrentstatestrprimary},
		{"global_setvarprocprioint100", global_setvarprocprioint100},
		{"global_setvarprocthresholdint200subproc", global_setvarprocthresholdint200subproc},
		{"global_sslenginerdrand", global_sslenginerdrand},
		{"global_sslenginerdrandALL", global_sslenginerdrandALL},
		{"global_sslenginerdrandRSADSA", global_sslenginerdrandRSADSA},
		{"global_sslmodeasync", global_sslmodeasync},
		{"global_statssocket1270018080", global_statssocket1270018080},
		{"global_statssocket1270018080modeadmin", global_statssocket1270018080modeadmin},
		{"global_statssocketsomepathtosocket", global_statssocketsomepathtosocket},
		{"global_statssocketsomepathtosocketmodea", global_statssocketsomepathtosocketmodea},
		{"global_threadgroupname10", global_threadgroupname10},
		{"global_threadgroupname110", global_threadgroupname110},
		{"global_tunequicsocketownerconnection", global_tunequicsocketownerconnection},
		{"global_tunequicsocketownerlistener", global_tunequicsocketownerlistener},
		{"global_unixbindprefixpre", global_unixbindprefixpre},
		{"global_unixbindprefixpremodetest", global_unixbindprefixpremodetest},
		{"global_unixbindprefixpremodetestusergga", global_unixbindprefixpremodetestusergga},
		{"global_unixbindprefixpremodetestusergga_", global_unixbindprefixpremodetestusergga_},
		{"global_unixbindprefixpremodetestusergga__", global_unixbindprefixpremodetestusergga__},
		{"global_unixbindprefixpremodetestusergga___", global_unixbindprefixpremodetestusergga___},
	}
	for _, config := range tests {
		t.Run(config.Name, func(t *testing.T) {
			t.Parallel()
			var buffer bytes.Buffer
			buffer.WriteString(config.Config)
			p, err := parser.New(options.Reader(&buffer))
			if err != nil {
				t.Fatal(err.Error())
			}
			result := p.String()
			if result != config.Config {
				compare(t, config.Config, result)
				t.Error("======== ORIGINAL =========")
				t.Error(config.Config)
				t.Error("======== RESULT ===========")
				t.Error(result)
				t.Error("===========================")
				t.Fatal("configurations does not match")
			}
		})
	}
}
