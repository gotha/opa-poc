package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", loggingMiddleware(helloHandler))
	http.HandleFunc("/resource1", loggingMiddleware(resource1Handler))
	http.HandleFunc("/resource2", loggingMiddleware(resource2Handler))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func resource1Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from resource 1")
}

func resource2Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello from resource 2")
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"[REQUEST] %s %s from %s | Query: %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			r.URL.RawQuery,
		)
		next.ServeHTTP(w, r)
	})
}
