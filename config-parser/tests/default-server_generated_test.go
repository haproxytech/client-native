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

package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/haproxytech/client-native/v6/config-parser/parsers"
)

func TestDefaultServer(t *testing.T) {
	tests := map[string]bool{
		"default-server addr 127.0.0.1":                      true,
		"default-server addr ::1":                            true,
		"default-server agent-check":                         true,
		"default-server agent-send name":                     true,
		"default-server agent-inter 1000ms":                  true,
		"default-server agent-addr 127.0.0.1":                true,
		"default-server agent-addr site.com":                 true,
		"default-server agent-port 1":                        true,
		"default-server agent-port 65535":                    true,
		"default-server allow-0rtt":                          true,
		"default-server alpn h2":                             true,
		"default-server alpn http/1.1":                       true,
		"default-server alpn h2,http/1.1":                    true,
		"default-server backup":                              true,
		"default-server ca-file cert.crt":                    true,
		"default-server check":                               true,
		"default-server check-send-proxy":                    true,
		"default-server check-alpn http/1.0":                 true,
		"default-server check-alpn http/1.1,http/1.0":        true,
		"default-server check-proto h2":                      true,
		"default-server check-ssl":                           true,
		"default-server check-via-socks4":                    true,
		"default-server ciphers ECDHE-RSA-AES128-GCM-SHA256": true,
		"default-server ciphers ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA384:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES256-SHA384:ECDHE-ECDSA-AES256-SHA:ECDHE-RSA-AES256-SHA:DHE-RSA-AES128-SHA256:DHE-RSA-AES128-SHA:DHE-RSA-AES256-SHA256:DHE-RSA-AES256-SHA:ECDHE-ECDSA-DES-CBC3-SHA:ECDHE-RSA-DES-CBC3-SHA:EDH-RSA-DES-CBC3-SHA:AES128-GCM-SHA256:AES256-GCM-SHA384:AES128-SHA256:AES256-SHA256:AES128-SHA:AES256-SHA:DES-CBC3-SHA:!DSS":      true,
		"default-server ciphersuites ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA384:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES256-SHA384:ECDHE-ECDSA-AES256-SHA:ECDHE-RSA-AES256-SHA:DHE-RSA-AES128-SHA256:DHE-RSA-AES128-SHA:DHE-RSA-AES256-SHA256:DHE-RSA-AES256-SHA:ECDHE-ECDSA-DES-CBC3-SHA:ECDHE-RSA-DES-CBC3-SHA:EDH-RSA-DES-CBC3-SHA:AES128-GCM-SHA256:AES256-GCM-SHA384:AES128-SHA256:AES256-SHA256:AES128-SHA:AES256-SHA:DES-CBC3-SHA:!DSS": true,
		"default-server cookie value":                              true,
		"default-server crl-file file.pem":                         true,
		"default-server crt cert.pem":                              true,
		"default-server disabled":                                  true,
		"default-server enabled":                                   true,
		"default-server error-limit 50":                            true,
		"default-server fall 30":                                   true,
		"default-server fall 1 rise 2 inter 3s port 4444":          true,
		"default-server force-sslv3":                               true,
		"default-server force-tlsv10":                              true,
		"default-server force-tlsv11":                              true,
		"default-server force-tlsv12":                              true,
		"default-server force-tlsv13":                              true,
		"default-server init-addr last,libc,none":                  true,
		"default-server init-addr last,libc,none,127.0.0.1":        true,
		"default-server inter 1500ms":                              true,
		"default-server inter 1000 weight 13":                      true,
		"default-server fastinter 2500ms":                          true,
		"default-server fastinter unknown":                         true,
		"default-server downinter 3500ms":                          true,
		"default-server log-proto legacy":                          true,
		"default-server log-proto octet-count":                     true,
		"default-server maxconn 1":                                 true,
		"default-server maxconn 50":                                true,
		"default-server maxqueue 0":                                true,
		"default-server maxqueue 1000":                             true,
		"default-server max-reuse -1":                              true,
		"default-server max-reuse 0":                               true,
		"default-server max-reuse 1":                               true,
		"default-server minconn 1":                                 true,
		"default-server minconn 50":                                true,
		"default-server namespace test":                            true,
		"default-server no-agent-check":                            true,
		"default-server no-backup":                                 true,
		"default-server no-check":                                  true,
		"default-server no-check-ssl":                              true,
		"default-server no-send-proxy-v2":                          true,
		"default-server no-send-proxy-v2-ssl":                      true,
		"default-server no-send-proxy-v2-ssl-cn":                   true,
		"default-server no-ssl":                                    true,
		"default-server no-ssl-reuse":                              true,
		"default-server no-sslv3":                                  true,
		"default-server no-tls-tickets":                            true,
		"default-server no-tlsv10":                                 true,
		"default-server no-tlsv11":                                 true,
		"default-server no-tlsv12":                                 true,
		"default-server no-tlsv13":                                 true,
		"default-server no-verifyhost":                             true,
		"default-server no-tfo":                                    true,
		"default-server non-stick":                                 true,
		"default-server npn http/1.1,http/1.0":                     true,
		"default-server observe layer4":                            true,
		"default-server observe layer7":                            true,
		"default-server on-error fastinter":                        true,
		"default-server on-error fail-check":                       true,
		"default-server on-error sudden-death":                     true,
		"default-server on-error mark-down":                        true,
		"default-server on-marked-down shutdown-sessions":          true,
		"default-server on-marked-up shutdown-backup-session":      true,
		"default-server pool-max-conn -1":                          true,
		"default-server pool-max-conn 0":                           true,
		"default-server pool-max-conn 100":                         true,
		"default-server pool-purge-delay 0":                        true,
		"default-server pool-purge-delay 5":                        true,
		"default-server pool-purge-delay 500":                      true,
		"default-server port 27015":                                true,
		"default-server port 27016":                                true,
		"default-server proto h2":                                  true,
		"default-server redir http://image1.mydomain.com":          true,
		"default-server redir https://image1.mydomain.com":         true,
		"default-server rise 2":                                    true,
		"default-server rise 200":                                  true,
		"default-server resolve-opts allow-dup-ip":                 true,
		"default-server resolve-opts ignore-weight":                true,
		"default-server resolve-opts allow-dup-ip,ignore-weight":   true,
		"default-server resolve-opts prevent-dup-ip,ignore-weight": true,
		"default-server resolve-prefer ipv4":                       true,
		"default-server resolve-prefer ipv6":                       true,
		"default-server resolve-net 10.0.0.0/8":                    true,
		"default-server resolve-net 10.0.0.0/8,10.0.0.0/16":        true,
		"default-server resolvers mydns":                           true,
		"default-server send-proxy":                                true,
		"default-server send-proxy-v2":                             true,
		"default-server proxy-v2-options ssl":                      true,
		"default-server proxy-v2-options ssl,cert-cn":              true,
		"default-server proxy-v2-options ssl,cert-cn,ssl-cipher,cert-sig,cert-key,authority,crc32c,unique-id": true,
		"default-server send-proxy-v2-ssl":    true,
		"default-server send-proxy-v2-ssl-cn": true,
		"default-server slowstart 2000ms":     true,
		"default-server sni TODO":             true,
		"default-server source TODO":          true,
		"default-server ssl":                  true,
		"default-server ssl-max-ver SSLv3":    true,
		"default-server ssl-max-ver TLSv1.0":  true,
		"default-server ssl-max-ver TLSv1.1":  true,
		"default-server ssl-max-ver TLSv1.2":  true,
		"default-server ssl-max-ver TLSv1.3":  true,
		"default-server ssl-min-ver SSLv3":    true,
		"default-server ssl-min-ver TLSv1.0":  true,
		"default-server ssl-min-ver TLSv1.1":  true,
		"default-server ssl-min-ver TLSv1.2":  true,
		"default-server ssl-min-ver TLSv1.3":  true,
		"default-server ssl-reuse":            true,
		"default-server stick":                true,
		"default-server socks4 127.0.0.1:81":  true,
		"default-server tcp-ut 20ms":          true,
		"default-server tfo":                  true,
		"default-server track TODO":           true,
		"default-server tls-tickets":          true,
		"default-server verify none":          true,
		"default-server verify required":      true,
		"default-server verifyhost site.com":  true,
		"default-server weight 1":             true,
		"default-server weight 128":           true,
		"default-server weight 256":           true,
		"default-server pool-low-conn 384":    true,
		"default-server ws h1":                true,
		"default-server ws h2":                true,
		"default-server ws auto":              true,
		"default-server log-bufsize 10":       true,
		"default-server":                      false,
		"---":                                 false,
		"--- ---":                             false,
	}
	parser := &parsers.DefaultServer{}
	for command, shouldPass := range tests {
		t.Run(command, func(t *testing.T) {
			line := strings.TrimSpace(command)
			lines := strings.SplitN(line, "\n", -1)
			var err error
			parser.Init()
			if len(lines) > 1 {
				for _, line = range lines {
					line = strings.TrimSpace(line)
					if err = ProcessLine(line, parser); err != nil {
						break
					}
				}
			} else {
				err = ProcessLine(line, parser)
			}
			if shouldPass {
				if err != nil {
					t.Error(err)
					return
				}
				result, err := parser.Result()
				if err != nil {
					t.Error(err)
					return
				}
				var returnLine string
				if result[0].Comment == "" {
					returnLine = result[0].Data
				} else {
					returnLine = fmt.Sprintf("%s # %s", result[0].Data, result[0].Comment)
				}
				if command != returnLine {
					t.Errorf("error: has [%s] expects [%s]", returnLine, command)
				}
			} else {
				if err == nil {
					t.Errorf("error: did not throw error for line [%s]", line)
				}
				_, parseErr := parser.Result()
				if parseErr == nil {
					t.Errorf("error: did not throw error on result for line [%s]", line)
				}
			}
		})
	}
}
