package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Age      int    `json:"age"`
}

var sliceOfUser []User
var db *sql.DB

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
func initDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/train"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Connection is Established")
}

func handleRootPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This Is The Root Page")
}

func handleSavedData(name string, lastname string, age int) {
	query := "insert into user(name , last , age) values(? , ? , ?)"
	result, err := db.Exec(query, name, lastname, age)
	if err != nil {
		panic(err)
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Println("Rows inserted:", rowsAffected)
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Bad Method Request", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Bad Method Request", http.StatusUnsupportedMediaType)
		return
	}

	var request User
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&request); err != nil {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	handleSavedData(request.Name, request.Lastname, request.Age)
	request_answer := map[string]string{
		"message": "ok",
	}
	encode := json.NewEncoder(w)
	encode.Encode(request_answer)
}
func handleBringData() {
	result, err := db.Query("select name ,last ,age from user")
	if err != nil {
		panic(err)
	}
	for result.Next() {
		var name string
		var lastname string
		var age int
		result.Scan(&name, &lastname, &age)
		mapOfUser := User{Name: name,
			Lastname: lastname,
			Age:      age}

		sliceOfUser = append(sliceOfUser, mapOfUser)
	}

	for _, user := range sliceOfUser {
		fmt.Printf("name is : %s lastname is : %s and age is : %d \n", user.Name, user.Lastname, user.Age)
	}
}
func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Bad-Method", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Bad Content-type", http.StatusUnsupportedMediaType)
		return
	}
	handleBringData()

	encode := json.NewEncoder(w)
	encode.Encode(sliceOfUser)
}

func main() {
	initDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRootPage)
	mux.HandleFunc("/AddUser", handlePostRequest)
	mux.HandleFunc("/SendUser", handleGetRequest)
	fmt.Println("The Server is Listening...")
	http.ListenAndServe(":8080", mux)

}
