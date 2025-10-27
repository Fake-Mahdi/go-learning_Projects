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
func handleRootPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This Is The Root Page")
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
	fmt.Println("Connection into database Has Been Established")
}

func savedDataToSql(name string, lastname string, age int) {
	query := " insert into user values(?, ? , ?)"
	row, err := db.Exec(query, name, lastname, age)
	if err != nil {
		panic(err)
	}

	rowsAffected, _ := row.RowsAffected()
	fmt.Println("rows Inserted : ", rowsAffected)
}
func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "http this is not the Method", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsuported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	var req User
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	savedDataToSql(req.Name, req.Lastname, req.Age)
	response := map[string]string{
		"Message": "The Data was Received",
	}

	encode := json.NewEncoder(w)
	encode.Encode(&response)
}
func getDataFromDatabase() {
	sliceOfUser = nil
	query := "select * from user"
	result, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer result.Close()
	var user User
	for result.Next() {
		if err := result.Scan(&user.Name, &user.Lastname, &user.Age); err != nil {
			panic(err)
		}
		sliceOfUser = append(sliceOfUser, user)
	}
	if err := result.Err(); err != nil {
		panic(err)
	}
}
func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "http this is not the Method", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsuported Media Type", http.StatusUnsupportedMediaType)
		return
	}
	getDataFromDatabase()
	encode := json.NewEncoder(w)
	encode.Encode(sliceOfUser)
}

func main() {
	initDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRootPage)
	mux.HandleFunc("/insertData", handlePostRequest)
	mux.HandleFunc("/selectData", handleGetRequest)
	fmt.Println("The Server is Listening...")
	http.ListenAndServe(":8080", mux)

}
