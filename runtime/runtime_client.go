package runtime

import (
	"fmt"

	"github.com/haproxytech/models"
)

//Client handles multiple HAProxy clients
type Client struct {
	runtimes []SingleRuntime
}

const (
	// DefaultSocketPath sane default for runtime API socket path
	DefaultSocketPath string = "/var/run/haproxy.sock"
	// DefaultSocketAutoRecconect sane default for runtime API autoReconnect
	DefaultSocketAutoRecconect bool = true
)

// DefaultClient return runtime Client with sane defaults
func DefaultClient() (*Client, error) {
	c := &Client{}
	err := c.Init([]string{DefaultSocketPath}, DefaultSocketAutoRecconect)
	if err != nil {
		return nil, err
	}
	return c, nil
}

//Init must be given path to runtime socket
func (c *Client) Init(socketPath []string, autoReconnect bool) error {
	c.runtimes = make([]SingleRuntime, len(socketPath))
	for index, path := range socketPath {
		runtime := SingleRuntime{}
		err := runtime.Init(path, autoReconnect)
		if err != nil {
			return err
		}
		c.runtimes[index] = runtime
	}
	return nil
}

func (c *Client) GetStats() ([]models.NativeStats, error) {
	result := make([]models.NativeStats, len(c.runtimes))
	for index, runtime := range c.runtimes {
		stats, err := runtime.GetStats()
		if err != nil {
			return nil, err
		}
		result[index] = stats
	}
	return result, nil
}

func (c *Client) GetInfo() ([]string, error) {
	result := make([]string, len(c.runtimes))
	for index, runtime := range c.runtimes {
		stats, err := runtime.GetInfo()
		if err != nil {
			return nil, err
		}
		result[index] = stats
	}
	return result, nil
}

//SetServerAddr set ip [port] for server
func (c *Client) SetServerAddr(backend, server string, ip string, port int) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerAddr(backend, server, ip, port)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//SetServerState set state for server
func (c *Client) SetServerState(backend, server string, state string) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerState(backend, server, state)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//SetServerWeight set weight for server
func (c *Client) SetServerWeight(backend, server string, weight string) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerWeight(backend, server, weight)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//SetServerHealth set health for server
func (c *Client) SetServerHealth(backend, server string, health string) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerHealth(backend, server, health)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}
