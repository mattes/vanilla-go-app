package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/templarbit/vanilla-go-app/server"
)

func main() {
	listenFlag := flag.String("listen", ":8080", "")
	httpKeepAliveFlag := flag.Bool("http-keep-alive", true, "")

	flag.Parse()

	srv := &http.Server{
		Addr:    *listenFlag,
		Handler: server.Server(),
	}

	// By default, keep-alives are always enabled
	// https://golang.org/pkg/net/http/#Server.SetKeepAlivesEnabled
	if !*httpKeepAliveFlag {
		srv.SetKeepAlivesEnabled(false)
	}

	log.Printf("Start listening on %v", *listenFlag)
	log.Fatal(srv.ListenAndServe())
}
