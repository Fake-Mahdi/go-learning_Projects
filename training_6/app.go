package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func enableCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "This is not The wanted Method", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "This is not The wanted Method", http.StatusUnsupportedMediaType)
		return
	}

	data_message := map[string]interface{}{
		"Message": true,
	}
	encode := json.NewEncoder(w)
	encode.Encode(data_message)
}
func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "This is not The wanted Method", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "This is not The wanted Method", http.StatusUnsupportedMediaType)
		return
	}

	var req User
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	data_message := map[string]interface{}{
		"Message": true,
	}
	encode := json.NewEncoder(w)
	encode.Encode(data_message)
}
func main() {
	mux := http.NewServeMux()
	fmt.Println("The server Listen On Port : 8080")
	mux.HandleFunc("/testgetapi", handleGetRequest)
	mux.HandleFunc("/testpostapi", handlePostRequest)
	http.ListenAndServe(":8080", mux)
}
