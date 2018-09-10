package main

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver"
)

func main() {
	server := golangredisserver.New()
	server.Init()

	err := server.Start()
	if err != nil {
		fmt.Errorf("Fatal error - unable to start server %v", err)
	}
}
