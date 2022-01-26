# HAProxy runtime client

## usage

```go
package main

import (
	"log"

	"github.com/haproxytech/client-native/v3/runtime"
)

func main() {
	client := runtime.SingleRuntime{}
	err := client.Init("/var/run/haproxy-runtime-api.sock", 1, 0)
	if err != nil {
		log.Println(err)
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
