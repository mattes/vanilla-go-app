package server

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Server() *http.ServeMux {
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
		if r.Header.Get("Content-Type") != "" {
			w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}
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
