package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/templarbit/vanilla-go-app/server"
)

func main() {
	listenFlag := flag.String("listen", ":8080", "")
	httpKeepAliveFlag := flag.Bool("http-keep-alive", true, "")
	httpReadTimeoutFlag := flag.Duration("http-read-timeout", 10*time.Second, "")
	httpWriteTimeoutFlag := flag.Duration("http-write-timeout", 10*time.Second, "")
	httpIdleTimeoutFlag := flag.Duration("http-idle-timeout", 30*time.Second, "")
	tcpKeepAliveFlag := flag.Bool("tcp-keep-alive", true, "")
	tcpIdleTimeoutFlag := flag.Duration("tcp-idle-timeout", 60*time.Second, "")

	flag.Parse()

	httpServer := &http.Server{
		Handler:      server.Server(),
		ReadTimeout:  *httpReadTimeoutFlag,
		WriteTimeout: *httpWriteTimeoutFlag,
		IdleTimeout:  *httpIdleTimeoutFlag,
	}

	// By default, http keep-alives are enabled
	// https://golang.org/pkg/net/http/#Server.SetKeepAlivesEnabled
	if !*httpKeepAliveFlag {
		httpServer.SetKeepAlivesEnabled(false)
	}

	l, err := net.Listen("tcp", *listenFlag)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start listening", *listenFlag)
	lt := TcpKeepAliveListener{l.(*net.TCPListener), *tcpKeepAliveFlag, *tcpIdleTimeoutFlag}
	log.Fatal(httpServer.Serve(lt))
}

// TcpKeepAliveListener is more or less copied from:
// https://github.com/golang/go/blob/release-branch.go1.10/src/net/http/server.go#L3211
type TcpKeepAliveListener struct {
	*net.TCPListener
	KeepAlive bool
	Timeout   time.Duration
}

func (ln TcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(ln.KeepAlive)
	tc.SetKeepAlivePeriod(ln.Timeout)
	return tc, nil
}
