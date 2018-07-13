package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	connectionDrainingTimeoutFlag := flag.Duration("connection-draining-timeout", 30*time.Second, "")
	shutdownDelayFlag := flag.Duration("shutdown-delay", 15*time.Second, "")

	flag.Parse()

	// Register SIGINT and SIGTERM termination calls
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

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

	go func() {
		log.Println("Start listening", *listenFlag)
		lx := server.NewTcpKeepAliveListener(l.(*net.TCPListener), *tcpKeepAliveFlag, *tcpIdleTimeoutFlag)
		err := httpServer.Serve(lx)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-shutdown
	log.Println("Shutting down ...")

	// 1. Delay shutdown
	// This is useful when run inside Docker container with Kubernetes as
	// scheduler for example: Kubernetes will send SIGTERM to let container
	// know end of life is coming soon, but expect the container to still serve
	// HTTP request until Kubernetes has updated all proxies and the container
	// is finally removed.
	time.Sleep(*shutdownDelayFlag)

	// 2. Close all open listeners
	// 3. Close all idle connections
	// 4. Wait up to connectionDrainingTimeoutFlag for existing connections
	//    to return to idle and then shut down
	// 5. After connectionDrainingTimeoutFlag, forcefully shutdown
	// Docs: https://golang.org/pkg/net/http/#Server.Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), *connectionDrainingTimeoutFlag)
	defer cancel()
	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bye")
}
