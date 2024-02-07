package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret-key-replace-with-your-own"))

type CoursesEnrolled struct {
	CourseID    int    `json:"courseid,omitempty"`
	CourseTitle string `json:"title,omitempty"`
}

var db *sql.DB
var err error

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
