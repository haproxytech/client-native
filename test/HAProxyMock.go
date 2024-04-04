package test

import (
	"bufio"
	"io"
	"io/ioutil"
	"net"
	"os"
	"testing"
)

// HAProxyMock - Mock HAProxy Server for testing the socket communication
type HAProxyMock struct {
	net.Listener
	responses map[string]string
	t         *testing.T
	running   bool
}

// NewHAProxyMock - create new haproxy mock
func NewHAProxyMock(t *testing.T) *HAProxyMock {
	haProxyMock := &HAProxyMock{}
	haProxyMock.t = t
	l, err := net.Listen("unix", socket())
	if err != nil {
		t.Fatal(err)
	}
	haProxyMock.Listener = l
	return haProxyMock
}

// Stop - stop mock
func (haproxy *HAProxyMock) Stop() {
	haproxy.running = false
}

// Start - start mock
func (haproxy *HAProxyMock) Start() {
	haproxy.running = true
	go func() {
		for {
			if haproxy.running == false {
				return
			}
			conn, err := haproxy.Accept()
			if err != nil {
				haproxy.t.Error(err)
				return
			}
			haproxy.handleConnection(conn)
		}
	}()
}

// SetResponses - setting the expected responses
func (haproxy *HAProxyMock) SetResponses(responses *map[string]string) {
	haproxy.responses = *responses
}

func (haproxy *HAProxyMock) handleConnection(conn net.Conn) {
	go func() {
		defer func() {
			_ = conn.Close()
		}()

		r := bufio.NewReader(conn)
		w := bufio.NewWriter(conn)
		s, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			haproxy.t.Log(err)
			return
		}

		response := haproxy.responses[s]
		_, err = w.WriteString(response)
		if err != nil {
			haproxy.t.Log(err)
			return
		}
		err = w.Flush()
		if err != nil {
			haproxy.t.Log(err)
			return
		}
	}()
}

func socket() string {
	f, err := ioutil.TempFile("", "haproxy-sock")
	if err != nil {
		panic(err)
	}
	addr := f.Name()
	_ = f.Close()
	_ = os.Remove(addr)
	return addr
}
