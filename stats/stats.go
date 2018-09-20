package stats

import (
	"github.com/haproxytech/models"
)

const (
	DefaultSocketFile string = "/var/run/haproxy"
)

func NewStatsClient(socketFile string) *StatsClient {
	if socketFile == "" {
		socketFile = DefaultSocketFile
	}

	client := new(StatsClient)
	client.socketFile = socketFile

	return client
}

func DefaultStatsClient() *StatsClient {
	return NewStatsClient("")
}

func (self *StatsClient) SocketFile() string {
	return self.socketFile
}

func (self *StatsClient) GetStats() (*models.Stats, error) {
	return nil, nil
}

type StatsClient struct {
	socketFile string
}

type Client interface {
	GetStats() (*models.Stats, error)
}