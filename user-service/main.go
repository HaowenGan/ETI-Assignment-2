package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// User structure
type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

var db *sql.DB
var err error

// connectDatabase establishes a connection to the MySQL database
func connectDatabase() {
	db, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/elearning")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	connectDatabase()
	defer db.Close()

	router := mux.NewRouter()
	log.Fatal(http.ListenAndServe(":5000", router))
}
