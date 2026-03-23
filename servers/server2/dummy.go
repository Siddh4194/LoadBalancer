package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	// 1. Register a handler function for the default route ("/")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World from server 2!")
	})


	// 2. Start the HTTP server
	port := ":3001"
	log.Printf("Starting server on port %s", port)
	// http.ListenAndServe listens on the TCP network address port and then calls Serve
	// to handle incoming requests.
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}