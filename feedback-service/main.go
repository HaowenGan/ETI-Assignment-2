// /feedback-service/main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret-key-replace-with-your-own"))

// A simple CORS middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5000") // allow requests from your front-end
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// If it's a preflight OPTIONS request, send a 200 response
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Serve the next handler
		next.ServeHTTP(w, r)
	})
}

// User represents a user structure
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// Review represents a user review structure
type Review struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	UserID   int    `json:"userId"`
	CourseID int    `json:"courseId"`
	Rating   int    `json:"rating"`
	Comment  string `json:"comment"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/eti_asg2")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

// SubmitReviewHandler handles the submission of new reviews
func SubmitReviewHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil || session.IsNew {
		http.Error(w, "Unauthorized - Session not found", http.StatusUnauthorized)
		return
	}
	// Decode the JSON request body into a new review
	var newReview Review
	err = json.NewDecoder(r.Body).Decode(&newReview)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the user ID from the session
	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Error(w, "Session does not contain user ID", http.StatusInternalServerError)
		return
	}

	// Insert the new review into the database
	userName, ok := session.Values["username"].(string)
	if !ok {
		http.Error(w, "Session does not contain username", http.StatusInternalServerError)
		return
	}

	result, err := db.Exec("INSERT INTO reviews (user_id, username, course_id, rating, comment) VALUES (?, ?, ?, ?, ?)",
		userID, userName, newReview.CourseID, newReview.Rating, newReview.Comment)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to submit review", http.StatusInternalServerError)
		return
	}

	// Get the ID of the newly inserted review
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get review ID", http.StatusInternalServerError)
		return
	}
	newReview.ID = int(lastInsertID)
	newReview.UserID = userID // Ensure the review struct reflects the correct userID
	newReview.Username = userName

	// Respond with the created review
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newReview)
}

func main() {
	// Initialize the database connection
	initDB()
	defer db.Close()

	// Initialize the Gorilla Mux router
	router := mux.NewRouter()
	// Wrap the router with the CORS middleware
	handlerWithCORS := enableCORS(router)

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/submit-review", SubmitReviewHandler).Methods("POST")

	// Start the server
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", handlerWithCORS)
}
