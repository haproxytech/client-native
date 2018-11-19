package main

import (
	"fmt"
	"os"

	"github.com/haproxytech/client-native"
)

func main() {
	c, err := client_native.DefaultClient()
	if err != nil {
		fmt.Println("Panic")
		os.Exit(0)
	}

	b, err := c.Configuration.GetBackends("")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	bJSON, err := b.MarshalBinary()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	fmt.Println(string(bJSON))

	s, err := c.Runtime.GetStats()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	for _, ns := range s {
		for _, st := range ns {
			sJSON, err := st.MarshalBinary()

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}
			fmt.Println(string(sJSON))
		}

	}

}
