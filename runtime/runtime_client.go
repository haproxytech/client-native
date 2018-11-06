package runtime

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/haproxytech/models"
)

//Task has command to execute on runtime api, and response channel for result
type Task struct {
	command  string
	response chan string
}

//Client handles multiple HAProxy clients
type Client struct {
	SocketsPath []string
}

//SingleRuntime handles one runtime API
type SingleRuntime struct {
	socketOpen bool
	jobs       chan Task
	socketPath string
}

//Init must be given path to runtime socket
func (s *SingleRuntime) Init(socketPath string) error {
	s.socketPath = socketPath
	c, err := net.Dial("unix", socketPath)
	if err != nil {
		return err
	}
	s.jobs = make(chan Task)
	go s.handleIncommingJobs(c)
	return nil
}

func (s *SingleRuntime) handleIncommingJobs(c net.Conn) {
	_, err := c.Write([]byte(fmt.Sprintf("prompt\n")))
	if err != nil {
		return
	}
	log.Println("start")
	for {
		select {
		case job := <-s.jobs:
			result, err := s.readFromSocket(c, job.command)
			if err != nil {
				job.response <- ""
			} else {
				job.response <- result
			}
		case <-time.After(time.Duration(60) * time.Second):
			log.Println(s.readFromSocket(c, "show env"))
		}
	}
	defer c.Close()
}

func (s *SingleRuntime) readFromSocket(c net.Conn, command string) (string, error) {
	_, err := c.Write([]byte(fmt.Sprintf("%s\n", command)))
	if err != nil {
		return "", nil
	}
	time.Sleep(1e9)
	bufferSize := 1024
	buf := make([]byte, bufferSize)
	var data strings.Builder
	for {
		n, err := c.Read(buf[:])
		if err != nil {
			break
		}
		data.Write(buf[0:n])
		if n < bufferSize {
			break
		}
	}
	return data.String(), nil
}

func (s *SingleRuntime) readFromSocketClean(command string) (string, error) {
	c, err := net.Dial("unix", s.socketPath)
	if err != nil {
		return "", err
	}
	defer c.Close()

	_, err = c.Write([]byte(fmt.Sprintf("%s\n", command)))
	if err != nil {
		return "", nil
	}
	time.Sleep(1e9)
	buf := make([]byte, 1024)
	var data strings.Builder
	for {
		n, err := c.Read(buf[:])
		if err != nil {
			break
		}
		data.Write(buf[0:n])
	}
	return data.String(), nil
}

//ExecuteRaw executes command on runtime API and returns raw result
func (s *SingleRuntime) ExecuteRaw(command string) (string, error) {
	response := make(chan string)
	task := Task{
		command:  command,
		response: response,
	}
	s.jobs <- task
	return <-response, nil
}

//GetStats fetches HAProxy stats from runtime API
func (s *SingleRuntime) GetStats() (models.NativeStats, error) {
	rawdata, err := s.ExecuteRaw("show stat")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(rawdata[2:], "\n")
	result := models.NativeStats{}
	keys := strings.Split(lines[0], ",")
	//data := []map[string]string{}
	for i := 1; i < len(lines); i++ {
		data := map[string]string{}
		line := strings.Split(lines[i], ",")
		if len(line) < len(keys) {
			continue
		}
		for index, key := range keys {
			if len(line[index]) > 0 {
				data[key] = line[index]
			}
		}
		oneLineData := &models.NativeStatsItems{
			Name: line[0],
			Type: strings.ToLower(line[1]),
		}
		var st models.NativeStatsItemsStats
		err := mapstructure.WeakDecode(data, &st)
		if err != nil {
			continue
		}
		oneLineData.Stats = &st
		result = append(result, oneLineData)
	}
	return result, nil
}

//GetInfo fetches HAProxy info from runtime API
func (s *SingleRuntime) GetInfo() (string, error) {
	data, err := s.ExecuteRaw("show stat")
	if err != nil {
		return "", err
	}
	return data, nil
}
