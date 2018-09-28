package configuration

import (
	"fmt"
	"os"
	"testing"
)

const testConf = `
# _version=1
frontend test
  mode http                                                  #alctl: protocol analyser
  bind 192.168.1.1:80 name webserv
  bind 192.168.1.1:8080 name webserv2
  log global                                                 #alctl: log activation
  option httplog                                             #alctl: log format
  option dontlognull
  option contstats
  option log-separate-errors
  log-tag bla
  option httpclose
  timeout http-request 2s
  timeout http-keep-alive 3s
  maxconn 2000
  default_backend test
  use-backend test_2 if TRUE
  timeout client 4s
  option clitcpka

frontend test_2
  mode http                                                  #alctl: protocol analyser
  log global                                                 #alctl: log activation
  option httplog                                             #alctl: log format
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
  mode http                                                  #alctl: protocol analyser
  balance roundrobin
  log global                                                 #alctl: log activation
  log-tag bla
  option httplog                                             #alctl: log format
  option http-keep-alive                                     #alctl: http connection mode
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
  server webserv 192.168.1.1:9200 maxconn 1000 ssl weight 10 cookie BLAH
  server webserv2 192.168.1.1:9300 maxconn 1000 ssl weight 10 cookie BLAH

backend test_2
  mode http                                                  #alctl: protocol analyser
  balance roundrobin
  log global  
  log-tag bla                                                #alctl: log activation
  option httplog                                             #alctl: log format
  option http-keep-alive                                     #alctl: http connection mode
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

var client Client
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

func prepareClient(path string) Client {
	return NewLBCTLClient(path, "/usr/sbin/lbctl", "/tmp/lbctl")
}
