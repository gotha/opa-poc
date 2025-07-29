package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var opaURL = "http://opa.default.svc.cluster.local:8181/v1/data/myapp/allow"

type input struct {
	Request request `json:"request"`
}

type request struct {
	Headers map[string]string `json:"headers"`
	Path    string            `json:"path"`
}

type OPAInput struct {
	Input input `json:"input"`
}

type OPAResponse struct {
	Result bool `json:"result"`
}

func opaHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}

	queryPath := r.URL.Query().Get("path")

	// Construct the input payload
	input := OPAInput{
		Input: input{
			Request: request{
				Headers: map[string]string{
					"authorization": r.Header.Get("Authorization"),
				},
				Path: queryPath,
			},
		},
	}

	fmt.Println(input)

	payloadBytes, err := json.Marshal(input)
	if err != nil {
		http.Error(w, "Failed to encode input", http.StatusInternalServerError)
		return
	}

	// Send to OPA
	resp, err := http.Post(opaURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "OPA request failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read OPA response", http.StatusInternalServerError)
		return
	}

	var opaResp OPAResponse
	if err := json.Unmarshal(body, &opaResp); err != nil {
		http.Error(w, "Invalid OPA response", http.StatusInternalServerError)
		return
	}

	fmt.Println(opaResp)

	// Decision logic
	if opaResp.Result {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func main() {
	if v := os.Getenv("OPA_URL"); v != "" {
		opaURL = v
	}
	http.HandleFunc("/", loggingMiddleware(opaHandler))
	log.Println("Starting OPA proxy on :7080")
	if err := http.ListenAndServe(":7080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		queryPath := r.URL.Query().Get("path")

		log.Printf(
			"[REQUEST] %s %s | Query: %s | Body: %s | QueryPath: %s",
			r.Method,
			r.URL.Path,
			r.URL.RawQuery,
			string(body),
			queryPath,
		)
		next.ServeHTTP(w, r)
	})
}
