package main

import (
	"log"
	"net/http"
	"github.com/random-number-api/pkg"
)



func main() {
	srv := pkg.StartServer()
	
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	
	pkg.GracefulShutdown(srv)
}