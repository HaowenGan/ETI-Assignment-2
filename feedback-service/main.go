package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// User represents a user structure
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// Review represents a user review structure
type Review struct {
	ID       int    `json:"id"`
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
	// Decode the JSON request body into a new review
	var newReview Review
	err := json.NewDecoder(r.Body).Decode(&newReview)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the user exists and get the user ID
	var user User
	err = db.QueryRow("SELECT id FROM users WHERE id = ?", newReview.UserID).Scan(&user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Insert the new review into the database
	result, err := db.Exec("INSERT INTO reviews (user_id, course_id, rating, comment) VALUES (?, ?, ?, ?)",
		newReview.UserID, newReview.CourseID, newReview.Rating, newReview.Comment)
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

	// Define routes
	router.HandleFunc("/submit-review", SubmitReviewHandler).Methods("POST")

	// Start the server
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", router)
}
