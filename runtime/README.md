# HAProxy runtime client

## usage

```go
package main

import (
	"log"

	"github.com/haproxytech/client-native/v6/runtime"
	runtime_options "github.com/haproxytech/client-native/v6/runtime/options"
)

func main() {
	nbproc := 8
	ms := runtime_options.MasterSocket("/var/run/haproxy-mw.sock", nbproc)
	client, err = runtime_api.New(ctx, ms)
	if err != nil {
		return nil, fmt.Errorf("error setting up runtime client: %s", err.Error())
	}
	// or if not using master-worker
	socketList := map[int]string{
		1: "/var/run/haproxy-runtime-api.sock"
	}
	sockets := runtime_options.Sockets(socketList)
	client, err = runtime_api.New(ctx, mapsDir, sockets)
	if err != nil {
		return nil, fmt.Errorf("error setting up runtime client: %s", err.Error())
	}

	statsCollection := client.GetStats()
	if statsCollection.Error != "" {
		log.Println(err)
	}
	log.Println(statsCollection.Stats)

	processInfo := client.GetInfo()
	if processInfo.Error != "" {
		log.Println(err)
	}
	log.Println(processInfo.Info)

	env, err := client.ExecuteRaw("show env")
	if err != nil {
		log.Println(err)
	}
	log.Println(env)
}

```
