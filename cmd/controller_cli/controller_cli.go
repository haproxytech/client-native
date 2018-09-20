package main

import (
	"fmt"

	"github.com/haproxytech/client-native"
	"github.com/haproxytech/client-native/configuration"
	"github.com/haproxytech/client-native/stats"
	_ "github.com/haproxytech/models"

)

func main() {
	confClient := configuration.NewLBCTLClient("/home/mjuraga/projects/lbctl/haproxy.conf", "/home/mjuraga/projects/lbctl/lbctl", "/tmp/lbctl")
	statsClient := stats.DefaultStatsClient()

	client := client_native.New(confClient, statsClient)

	b, err := client.Configuration.GetBackends()

	bJson, err := b.MarshalBinary()

	if err != nil{
		fmt.Println(err)
	}

	fmt.Println("Backends: \n" + string(bJson) + "\n")

}