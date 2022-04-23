package main

import (
	"GreensOne/gsweb"
	"GreensOne/test"
	"net/http"
)

func main() {
	gs := gsweb.NewGreensCore()
	test.Register(gs)
	server := &http.Server{
		Handler: gs,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
