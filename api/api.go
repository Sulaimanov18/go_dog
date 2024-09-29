package main

import (
	"encoding/json"
	"net/http"

	"github.com/cucumber/godog"
)

// Handler for /version endpoint
func getVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fail(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := struct {
		Version string `json:"version"`
	}{Version: godog.Version}

	ok(w, data)
}

// Helper to handle errors
func fail(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	data := struct {
		Error string `json:"error"`
	}{Error: msg}
	resp, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// Helper to send OK response
func ok(w http.ResponseWriter, data interface{}) {
	resp, err := json.Marshal(data)
	if err != nil {
		fail(w, "Oops something went wrong", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// Start the server
func main() {
	http.HandleFunc("/version", getVersion)
	http.ListenAndServe(":8080", nil)
}
