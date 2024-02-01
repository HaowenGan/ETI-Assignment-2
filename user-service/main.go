// Ong Jia Yuan / S10227735B
// /user-service/main.go

package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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

// Establishes a connection to the MySQL database
func connectDatabase() {
	db, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/eti_asg2")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

// Create a new user in the database
func registerUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hashing password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users(firstName, lastName, email, username, password) VALUES(?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// loginUser authenticates a user
func loginUser(w http.ResponseWriter, r *http.Request) {
	var user User
	var hashedPassword string

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Sets login.html as the default page
func serveLogin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./front-end/login.html")
}

func main() {
	connectDatabase()
	defer db.Close()
	workDir, _ := os.Getwd()
	frontEndDir := filepath.Join(workDir, "../front-end")

	router := mux.NewRouter()

	router.PathPrefix("/front-end/").Handler(http.StripPrefix("/front-end/", http.FileServer(http.Dir(frontEndDir))))
	// Sets login.html as default page for the root URL
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(frontEndDir, "login.html"))
	})

	// API routes
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/register", registerUser).Methods("POST")
	apiRouter.HandleFunc("/login", loginUser).Methods("POST")

	log.Fatal(http.ListenAndServe(":5000", router))
}
