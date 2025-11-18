package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
)

var db *sql.DB
var secretKey = []byte("Mahdicode1234")

type User struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Age      int    `json:"age"`
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

func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func VerifyToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Not Authorized to Have access ", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	tokenString := parts[1]
	claims, err := VerifyToken(tokenString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
		return
	}
	fmt.Println(claims)
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
	fmt.Println("connection was been established")
}

func handleRootPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This Is The Root Page")
}

func SavedIntoDatabase(name string, lastname string, age int) {
	query := "insert into user values (? , ?, ?)"
	rows, err := db.Exec(query, name, lastname, age)
	if err != nil {
		panic(err)
	}
	affectedRows, _ := rows.RowsAffected()
	fmt.Println(affectedRows)
}
func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if enableCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Bad Method Request", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-type") != "application/json" {
		http.Error(w, "Bad Method Request", http.StatusUnsupportedMediaType)
		return
	}
	var req User
	decode := json.NewDecoder(r.Body)
	if err := decode.Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	SavedIntoDatabase(req.Name, req.Lastname, req.Age)
	token, err := GenerateToken(req.Name)
	if err != nil {
		panic(err)
	}
	request_answer := map[string]string{
		"message": "The data was Received",
		"token":   token,
	}
	encode := json.NewEncoder(w)
	encode.Encode(&request_answer)

}
func main() {
	initDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRootPage)
	mux.HandleFunc("/insertData", handlePostRequest)
	fmt.Println("Server is Listining in port : 8080")
	http.ListenAndServe(":8080", mux)

}
