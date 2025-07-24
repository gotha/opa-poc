package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/resource1", validApiKeyMiddleware(resource1Handler))
	http.HandleFunc("/resource2", validApiKeyMiddleware(resource2Handler))

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

func validApiKeyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		allowed, err := checkOPAPolicy(authHeader, r.URL.Path)
		if err != nil {
			http.Error(w, "Policy check failed", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

type OPARequest struct {
	Input struct {
		Request struct {
			Headers map[string]string `json:"headers"`
			Path    string            `json:"path"`
		} `json:"request"`
	} `json:"input"`
}

type OPAResponse struct {
	Result bool `json:"result"`
}

func checkOPAPolicy(authHeader, path string) (bool, error) {
	opaReq := OPARequest{}
	opaReq.Input.Request.Headers = map[string]string{
		"authorization": authHeader,
	}
	opaReq.Input.Request.Path = path

	jsonData, _ := json.Marshal(opaReq)

	resp, err := http.Post("http://localhost:8181/v1/data/myapp/allow",
		"application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var opaResp OPAResponse
	json.NewDecoder(resp.Body).Decode(&opaResp)
	return opaResp.Result, nil
}
