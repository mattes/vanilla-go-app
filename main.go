package main

import (
	"crypto/rand"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	listenFlag := flag.String("listen", ":8080", "")
	httpKeepAliveFlag := flag.Bool("http-keep-alive", true, "")

	flag.Parse()

	srv := &http.Server{
		Addr:    *listenFlag,
		Handler: server(),
	}

	// By default, keep-alives are always enabled
	// https://golang.org/pkg/net/http/#Server.SetKeepAlivesEnabled
	if !*httpKeepAliveFlag {
		srv.SetKeepAlivesEnabled(false)
	}

	log.Printf("Start listening on %v", *listenFlag)
	log.Fatal(srv.ListenAndServe())
}

func server() *http.ServeMux {
	mux := http.NewServeMux()

	// Default route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200) // OK
	})

	// Routes that return fixed KB of binary data
	kb := []int{0, 1, 10, 100, 1000}
	for _, v := range kb {
		b := make([]byte, v*1000)
		_, err := rand.Read(b)
		if err != nil {
			log.Fatal(err)
		}

		mux.HandleFunc("/bin/"+strconv.Itoa(v)+"KB",
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/octet-stream")
				w.WriteHeader(200) // OK
				w.Write(b)
			})
	}

	// Route that reads full body but discards input
	mux.HandleFunc("/readall", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(204) // No Content
	})

	// Route that reads full body and echos it back to client
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(201) // Created
		io.Copy(w, r.Body)
		r.Body.Close()
	})

	// Route that will sleep some time
	mux.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		msParam := r.URL.Query().Get("ms")
		ms, err := strconv.Atoi(msParam)
		if err != nil {
			w.WriteHeader(400) // Bad Request
			return
		}
		time.Sleep(time.Duration(ms) * time.Millisecond)
		w.WriteHeader(202) // Accepted
	})

	return mux
}
