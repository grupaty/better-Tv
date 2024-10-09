package main

import (
	"github.com/random-number-api/pkg"
)



func main() {
	resultStore := pkg.NewResultStore()
	
	server := pkg.NewServer(resultStore)
	server.Start()
	
	server.GracefulShutdown()
}