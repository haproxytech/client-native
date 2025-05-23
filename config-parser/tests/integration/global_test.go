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
		{"global_cpusetdropcluster05", global_cpusetdropcluster05},
		{"global_cpusetdropcluster13", global_cpusetdropcluster13},
		{"global_cpusetdropcluster1somecomment", global_cpusetdropcluster1somecomment},
		{"global_cpusetdropcore05", global_cpusetdropcore05},
		{"global_cpusetdropcore13", global_cpusetdropcore13},
		{"global_cpusetdropcore1somecomment", global_cpusetdropcore1somecomment},
		{"global_cpusetdropcpu05", global_cpusetdropcpu05},
		{"global_cpusetdropcpu13", global_cpusetdropcpu13},
		{"global_cpusetdropcpu1somecomment", global_cpusetdropcpu1somecomment},
		{"global_cpusetdropnode05", global_cpusetdropnode05},
		{"global_cpusetdropnode13", global_cpusetdropnode13},
		{"global_cpusetdropnode1somecomment", global_cpusetdropnode1somecomment},
		{"global_cpusetdropthread05", global_cpusetdropthread05},
		{"global_cpusetdropthread13", global_cpusetdropthread13},
		{"global_cpusetdropthread1somecomment", global_cpusetdropthread1somecomment},
		{"global_cpusetonlycluster05", global_cpusetonlycluster05},
		{"global_cpusetonlycluster13", global_cpusetonlycluster13},
		{"global_cpusetonlycluster1somecomment", global_cpusetonlycluster1somecomment},
		{"global_cpusetonlycore05", global_cpusetonlycore05},
		{"global_cpusetonlycore13", global_cpusetonlycore13},
		{"global_cpusetonlycore1somecomment", global_cpusetonlycore1somecomment},
		{"global_cpusetonlycpu05", global_cpusetonlycpu05},
		{"global_cpusetonlycpu13", global_cpusetonlycpu13},
		{"global_cpusetonlycpu1somecomment", global_cpusetonlycpu1somecomment},
		{"global_cpusetonlynode05", global_cpusetonlynode05},
		{"global_cpusetonlynode13", global_cpusetonlynode13},
		{"global_cpusetonlynode1somecomment", global_cpusetonlynode1somecomment},
		{"global_cpusetonlythread05", global_cpusetonlythread05},
		{"global_cpusetonlythread13", global_cpusetonlythread13},
		{"global_cpusetonlythread1somecomment", global_cpusetonlythread1somecomment},
		{"global_cpusetreset", global_cpusetreset},
		{"global_cpusetresetsomecomment", global_cpusetresetsomecomment},
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
