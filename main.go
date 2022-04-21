package GreensOne

import (
	"GreensOne/gsweb"
	"net/http"
)

func main() {
	server := &http.Server{
		Handler: gsweb.NewGreensCore(),
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
