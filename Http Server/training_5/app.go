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

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "This is not the needed Request Method", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "This is not the needed Header", http.StatusUnsupportedMediaType)
		return
	}

	var request User
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&request); err != nil {
		http.Error(w, "BAD REQUEST", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	sendData := map[string]interface{}{
		"Message": "HAHAHAAHHA just for fun",
	}
	json.NewEncoder(w).Encode(sendData)

	fmt.Printf("%s and %s", request.Name, request.Password)
}
func handleMainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world From Mahdi")
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "This is not the method we search for", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Fauls Form", http.StatusUnsupportedMediaType)
		return
	}

	send_data := map[string]interface{}{
		"name":     "Mahdi",
		"lastname": "Boukhouima",
	}
	encode := json.NewEncoder(w)
	encode.Encode(send_data)
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleMainPage)
	mux.HandleFunc("/user", handlePostRequest)
	mux.HandleFunc("/GetData", handleGetRequest)
	fmt.Println("The Server Is Listening on Port : 8080")
	http.ListenAndServe(":8080", mux)

}
