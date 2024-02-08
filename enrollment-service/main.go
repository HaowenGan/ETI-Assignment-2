package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Course represents a course object
type Course struct {
	ID      int     `json:"id,omitempty"`
	Title   string  `json:"title,omitempty"`
	Content string  `json:"content,omitempty"`
	Price   float64 `json:"price,omitempty"`
}

// UserData represents user enrollment data
type UserData struct {
	UserID    int   `json:"user_id"`
	CourseIDs []int `json:"course_ids"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:Joeykkh101204@tcp(127.0.0.1:3306)/dbname")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/courses", getCourses)
	http.HandleFunc("/enroll", enrollCourse)
	http.HandleFunc("/user/", getUserCourses)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getCourses(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, content, name FROM courses")
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

func enrollCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data UserData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID := data.UserID
	courseIDs := data.CourseIDs

	for _, courseID := range courseIDs {
		_, err := db.Exec("INSERT INTO user_courses (user_id, course_id) VALUES (?, ?)", userID, courseID)
		if err != nil {
			http.Error(w, "Failed to enroll user in courses", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Enrolled successfully")
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
