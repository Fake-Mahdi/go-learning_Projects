package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Age      int    `json:"age"`
}

func enableCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}
func connectIntoDatabase() {
	dsn := "root:@tcp(127.0.0.1:3306)/train"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("There was an error connecting into the database", err)
		return
	}
	if err := db.Ping(); err != nil {
		fmt.Println("There Was no Answer from the database server")
		return
	}
	fmt.Println("The connection was Been Established with The database")
}
func GetFromDataBase() (chan []byte, chan error) {
	jsonByte := make(chan []byte)
	errChan := make(chan error, 1)

	go func() {
		query := "select * from user"
		result, err := db.Query(query)
		if err != nil {
			fmt.Println("There was en Error in Query", err)
		}
		defer close(jsonByte)
		defer close(errChan)
		for result.Next() {
			var user User
			if err := result.Scan(&user.Name, &user.Lastname, &user.Age); err != nil {
				fmt.Println("error in scaning a certain row", err)
				continue
			}
			intoByte, err := json.Marshal(user)
			if err != nil {
				fmt.Println("Error in parsing into byte", err)
				errChan <- err
				continue
			}
			jsonByte <- intoByte

		}
		if err := result.Err(); err != nil {
			fmt.Println("the errors in in each row are", err)
		}
	}()
	return jsonByte, errChan
}

func HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong Request Type", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-type") != "application/json" {
		http.Error(w, "Wrong Request Type", http.StatusUnsupportedMediaType)
		return
	}
	jsonByteData, err := GetFromDataBase()
	if err != nil {
		fmt.Println("There was an error during Bringing Data from database", err)
	}

	for data := range jsonByteData {
		w.Write(data)
		w.Write([]byte("\n"))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

func main() {
	connectIntoDatabase()
	mux := http.NewServeMux()
	mux.HandleFunc("/GetData", HandleGetRequest)
	fmt.Println("Server Is Listening")
	http.ListenAndServe(":8080", mux)
}
