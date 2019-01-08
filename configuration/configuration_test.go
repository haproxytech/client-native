package configuration

import (
	"fmt"
	"os"
	"testing"
)

const testConf = `
# _version=1
frontend test
  mode http
  bind 192.168.1.1:80 name webserv
  bind 192.168.1.1:8080 name webserv2
  log global
  option httplog
  option dontlognull
  option contstats
  option log-separate-errors
  filter trace name BEFORE-HTTP-COMP random-parsing hexdump
  filter compression
  filter trace name AFTER-HTTP-COMP random-forwarding
  http-request allow if src 192.168.0.0/16
  http-request set-header X-SSL %[ssl_fc]
  http-request set-var(req.my_var) req.fhdr(user-agent),lower
  http-response allow if src 192.168.0.0/16
  http-response set-header X-SSL %[ssl_fc]
  http-response set-var(req.my_var) req.fhdr(user-agent),lower
  tcp-request connection accept if TRUE
  tcp-request connection reject if FALSE
  tcp-request content accept if TRUE
  tcp-request content reject if FALSE
  log-tag bla
  option httpclose
  timeout http-request 2s
  timeout http-keep-alive 3s
  maxconn 2000
  default_backend test
  use_backend test_2 if TRUE
  timeout client 4s
  option clitcpka

frontend test_2
  mode http
  log global
  option httplog
  option dontlognull
  option contstats
  option log-separate-errors
  log-tag bla
  option httpclose
  timeout http-request 2s
  timeout http-keep-alive 3s
  maxconn 2000
  default_backend test_2
  timeout client 4s
  option clitcpka

backend test
  mode http
  balance roundrobin
  log global
  log-tag bla
  option httplog
  option http-keep-alive
  option forwardfor header X-Forwarded-For
  option httpchk HEAD /
  default-server fall 2
  default-server rise 4
  default-server inter 5s
  default-server port 8888
  stick store-request src table test
  stick match src table test
  stick on src table test
  stick store-response src
  stick store-response src_port table test_port
  stick store-response src table test if TRUE
  tcp-response content accept if TRUE
  tcp-response content reject if FALSE
  option contstats
  timeout check 2s
  timeout tunnel 5s
  timeout server 3s
  cookie BLA
  use-server webserv if TRUE
  use-server webserv2 unless TRUE
  server webserv 192.168.1.1:9200 maxconn 1000 ssl weight 10 cookie BLAH
  server webserv2 192.168.1.1:9300 maxconn 1000 ssl weight 10 cookie BLAH

backend test_2
  mode http
  balance roundrobin
  log global  
  log-tag bla
  option httplog
  option http-keep-alive
  option forwardfor header X-Forwarded-For
  option httpchk HEAD /
  default-server fall 2
  default-server rise 4
  default-server inter 5s
  default-server port 8888
  option contstats
  timeout check 2s
  timeout tunnel 5s
  timeout server 3s
  cookie BLA
`
const testPath = "/tmp/haproxy-test.cfg"
const haproxyExec = "/usr/sbin/haproxy"

var client *Client
var version int64 = 1

func TestMain(m *testing.M) {
	err := prepareTestFile(testConf, testPath)
	if err != nil {
		fmt.Println("Could not prepare tests")
		os.Exit(1)
	}
	defer deleteTestFile(testPath)
	client = prepareClient(testPath)

	os.Exit(m.Run())
}

func prepareTestFile(conf string, path string) error {
	// detect if file exists
	var _, err = os.Stat(path)
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

	file, err = os.OpenFile(path, os.O_RDWR, 0644)
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

func prepareClient(path string) *Client {
	c := Client{}
	p := ClientParams{
		ConfigurationFile:       path,
		GlobalConfigurationFile: "",
		Haproxy:                 "echo",
		UseValidation:           true,
		UseCache:                true,
		LBCTLPath:               "/usr/sbin/lbctl",
		LBCTLTmpPath:            "/tmp/lbctl",
	}
	c.Init(p)
	return &c
}
