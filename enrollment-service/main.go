package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret-key-replace-with-your-own"))

// Course represents a course object
type Course struct {
	ID      int     `json:"id,omitempty"`
	Title   string  `json:"title,omitempty"`
	Content string  `json:"content,omitempty"`
	Price   float64 `json:"price,omitempty"`
}

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Usertype  string `json:"usertype"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbname")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/courses", getCourses)
	http.HandleFunc("/enroll", enrollCurrentUser) // Change handler to enrollCurrentUser
	http.HandleFunc("/user/", getUserCourses)

	fmt.Println("Server is running on port 6000...")
	log.Fatal(http.ListenAndServe(":6000", nil))
}

// enrollCurrentUser enrolls the current logged-in user into a course
func enrollCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the current user's ID from the session
	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Error retrieving session", http.StatusInternalServerError)
		log.Printf("Error retrieving session: %v", err)
		return
	}
	userID, ok := session.Values["userID"].(int)
	if !ok {
		http.Error(w, "User ID not found in session", http.StatusInternalServerError)
		return
	}

	// Decode the request body to get the course ID(s) to enroll the user into
	var courseIDs []int
	if err := json.NewDecoder(r.Body).Decode(&courseIDs); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Iterate over each course ID and check if it exists in the courses table
	for _, courseID := range courseIDs {
		// Check if the course exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM courses WHERE id = ?)", courseID).Scan(&exists)
		if err != nil {
			http.Error(w, "Error checking course existence", http.StatusInternalServerError)
			log.Printf("Error checking course existence: %v", err)
			return
		}
		if !exists {
			http.Error(w, fmt.Sprintf("Course with ID %d does not exist", courseID), http.StatusBadRequest)
			return
		}

		// Enroll the user into the course
		_, err = db.Exec("INSERT INTO user_courses (user_id, course_id) VALUES (?, ?)", userID, courseID)
		if err != nil {
			http.Error(w, "Failed to enroll user in courses", http.StatusInternalServerError)
			return
		}
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Enrolled successfully")
}

func getCourses(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, content, price FROM courses")
	if err != nil {
		http.Error(w, "Failed to retrieve courses", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.ID, &course.Title, &course.Content, &course.Price); err != nil {
			http.Error(w, "Failed to retrieve courses", http.StatusInternalServerError)
			return
		}
		courses = append(courses, course)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getUserCourses(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Path[len("/user/"):]

	rows, err := db.Query("SELECT c.id, c.title, c.content, c.price FROM courses c JOIN user_courses uc ON c.id = uc.course_id WHERE uc.user_id = ?", userID)
	if err != nil {
		http.Error(w, "Failed to retrieve user courses", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.ID, &course.Title, &course.Content, &course.Price); err != nil {
			http.Error(w, "Failed to retrieve user courses", http.StatusInternalServerError)
			return
		}
		courses = append(courses, course)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}
