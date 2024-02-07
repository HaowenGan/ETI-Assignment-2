// /feedback-service/main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PATCH")

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

// EditReviewHandler handles the editing of existing reviews
func EditReviewHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil || session.IsNew {
		http.Error(w, "Unauthorized - Session not found", http.StatusUnauthorized)
		return
	}

	// Decode the JSON request body into an updated review
	var updatedReview Review
	err = json.NewDecoder(r.Body).Decode(&updatedReview)
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

	// Get the username from the session
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Error(w, "Session does not contain username", http.StatusInternalServerError)
		return
	}

	// Check if the user is authorized to edit the review
	if userID != updatedReview.UserID || username != updatedReview.Username {
		http.Error(w, "Unauthorized - User does not have permission to edit this review", http.StatusUnauthorized)
		return
	}

	// Check if the user is authorized to edit the review (you may want to add additional checks)
	// For example, ensuring that the user is the author of the review.

	// Update the review in the database
	_, err = db.Exec("UPDATE reviews SET rating=?, comment=? WHERE id=? AND user_id=?",
		updatedReview.Rating, updatedReview.Comment, updatedReview.ID, userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update review", http.StatusInternalServerError)
		return
	}

	// Respond with the updated review
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedReview)
}

// DeleteReviewHandler handles the deletion of existing reviews
func DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil || session.IsNew {
		http.Error(w, "Unauthorized - Session not found", http.StatusUnauthorized)
		return
	}

	// Get the user ID from the session
	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Error(w, "Session does not contain user ID", http.StatusInternalServerError)
		return
	}

	// Get the username from the session
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Error(w, "Session does not contain username", http.StatusInternalServerError)
		return
	}

	// Get the review ID from the request parameters
	params := mux.Vars(r)
	reviewID, ok := params["id"]
	if !ok {
		http.Error(w, "Review ID not provided in the request", http.StatusBadRequest)
		return
	}

	// Convert the review ID to an integer
	reviewIDInt, err := strconv.Atoi(reviewID)
	if err != nil {
		http.Error(w, "Invalid Review ID", http.StatusBadRequest)
		return
	}

	// Check if the user is authorized to delete the review
	var reviewAuthorID int
	err = db.QueryRow("SELECT user_id FROM reviews WHERE id=?", reviewIDInt).Scan(&reviewAuthorID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to fetch review information", http.StatusInternalServerError)
		return
	}

	if userID != reviewAuthorID || username != username {
		http.Error(w, "Unauthorized - User does not have permission to delete this review", http.StatusUnauthorized)
		return
	}

	// Delete the review from the database
	_, err = db.Exec("DELETE FROM reviews WHERE id=? AND user_id=?", reviewIDInt, userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to delete review", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Review with ID %d deleted successfully", reviewIDInt)
}

func GetReviewsHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil || session.IsNew {
		log.Println("Session error or new session:", err)
		http.Error(w, "Unauthorized - Session not found", http.StatusUnauthorized)
		return
	}

	// Assuming userID is stored as an integer
	userID, ok := session.Values["userID"].(int)
	if !ok || userID == 0 {
		log.Printf("Session user ID not found or is zero, found: %v", session.Values["userID"])
		http.Error(w, "Session does not contain user ID", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, course_id, rating, comment FROM reviews WHERE user_id=?", userID)
	if err != nil {
		log.Printf("Database query error: %v", err)
		http.Error(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		if err := rows.Scan(&review.ID, &review.CourseID, &review.Rating, &review.Comment); err != nil {
			log.Printf("Error scanning review: %v", err)
			continue // Optionally, handle the error as you prefer
		}
		reviews = append(reviews, review)
	}

	log.Printf("Retrieved userID from session: %d", userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		log.Printf("Error encoding reviews: %v", err)
		http.Error(w, "Failed to encode reviews", http.StatusInternalServerError)
	}
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
	apiRouter.HandleFunc("/edit-review", EditReviewHandler).Methods("PATCH")
	apiRouter.HandleFunc("/delete-review/{id:[0-9]+}", DeleteReviewHandler).Methods("DELETE")
	apiRouter.HandleFunc("/get-reviews", GetReviewsHandler).Methods("GET")

	// Start the server
	fmt.Println("Server is running on :5001")
	http.ListenAndServe(":5001", handlerWithCORS)
}
