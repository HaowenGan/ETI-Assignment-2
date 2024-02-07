// main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

// Course structure
type Course struct {
	ID       int      `json:"id,omitempty"`
	Title    string   `json:"title,omitempty"`
	Content  string   `json:"content,omitempty"`
	Price    float64  `json:"price,omitempty"`
	Sections []Section `json:"sections,omitempty"`
}

// Section structure within a course
type Section struct {
	ID      int    `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

var db *sql.DB

func main() {
	// Initialize MySQL connection
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/eti_asg2")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize router
	router := mux.NewRouter()

	// Define API endpoints
	router.HandleFunc("/courses", GetCourses).Methods("GET")
	router.HandleFunc("/courses/{id}", GetCourse).Methods("GET")
	router.HandleFunc("/courses", CreateCourse).Methods("POST")
	router.HandleFunc("/courses/{id}", UpdateCourse).Methods("PUT")
	router.HandleFunc("/courses/{id}", DeleteCourse).Methods("DELETE")

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8080", router))
}

// GetCourses retrieves all courses
func GetCourses(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, content FROM courses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		err := rows.Scan(&course.ID, &course.Title, &course.Content)
		if err != nil {
			log.Fatal(err)
		}

		// Fetch sections for each course
		course.Sections, err = getSections(course.ID)
		if err != nil {
			log.Fatal(err)
		}

		courses = append(courses, course)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// GetCourse retrieves a specific course by ID
func GetCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var course Course
	err := db.QueryRow("SELECT id, title, content FROM courses WHERE id=?", params["id"]).Scan(&course.ID, &course.Title, &course.Content)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Fetch sections for the course
	course.Sections, err = getSections(course.ID)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

// CreateCourse creates a new course
func CreateCourse(w http.ResponseWriter, r *http.Request) {
    var course Course
    if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
        log.Printf("Error decoding JSON request: %v", err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // Insert course into the database
    result, err := db.Exec("INSERT INTO courses(title, content, price) VALUES(?, ?, ?)", course.Title, course.Content, course.Price)
    if err != nil {
        log.Printf("Error inserting course into the database: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // Get the last inserted ID
    courseID, err := result.LastInsertId()
    if err != nil {
        log.Printf("Error retrieving last inserted ID: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // Insert sections for the course if available
    for _, section := range course.Sections {
        _, err := db.Exec("INSERT INTO sections(course_id, title, content) VALUES(?, ?, ?)", courseID, section.Title, section.Content)
        if err != nil {
            log.Printf("Error inserting section into the database: %v", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(courseID)
}

// UpdateCourse updates an existing course by ID
func UpdateCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedCourse Course
	json.NewDecoder(r.Body).Decode(&updatedCourse)

	// Update course details
	_, err := db.Exec("UPDATE courses SET title=?, content=?, price=? WHERE id=?", updatedCourse.Title, updatedCourse.Content, updatedCourse.Price, params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Delete existing sections for the course
	_, err = db.Exec("DELETE FROM sections WHERE course_id=?", params["id"])
	if err != nil {
		log.Fatal(err)
	}

	// Insert updated sections for the course if available
	for _, section := range updatedCourse.Sections {
		_, err := db.Exec("INSERT INTO sections(course_id, title, content) VALUES(?, ?, ?)", params["id"], section.Title, section.Content)
		if err != nil {
			log.Fatal(err)
		}
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteCourse deletes a course by ID
func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Delete course and associated sections
	_, err := db.Exec("DELETE FROM courses WHERE id=?", params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = db.Exec("DELETE FROM sections WHERE course_id=?", params["id"])
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
}

// getSections retrieves sections for a specific course ID
func getSections(courseID int) ([]Section, error) {
	rows, err := db.Query("SELECT id, title, content FROM sections WHERE course_id=?", courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []Section
	for rows.Next() {
		var section Section
		err := rows.Scan(&section.ID, &section.Title, &section.Content)
		if err != nil {
			return nil, err
		}
		sections = append(sections, section)
	}

	return sections, nil
}
