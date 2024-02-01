package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Course structure
type Course struct {
	ID       string    `json:"id,omitempty" bson:"_id,omitempty"`
	Title    string    `json:"title,omitempty" bson:"title,omitempty"`
	Content  string    `json:"content,omitempty" bson:"content,omitempty"`
	Sections []Section `json:"sections,omitempty" bson:"sections,omitempty"`
}

// Section structure within a course
type Section struct {
	Title   string `json:"title,omitempty" bson:"title,omitempty"`
	Content string `json:"content,omitempty" bson:"content,omitempty"`
}

var client *mongo.Client

func main() {
	// Initialize MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Initialize router
	router := mux.NewRouter()

	// Define API endpoints
	router.HandleFunc("/courses", GetCourses).Methods("GET")
	router.HandleFunc("/courses/{id}", GetCourse).Methods("GET")
	router.HandleFunc("/courses", CreateCourse).Methods("POST")
	router.HandleFunc("/courses/{id}", UpdateCourse).Methods("PUT")
	router.HandleFunc("/courses/{id}", DeleteCourse).Methods("DELETE")

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8081", router))
}

// GetCourses retrieves all courses
func GetCourses(w http.ResponseWriter, r *http.Request) {
	var courses []Course
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := client.Database("coursedb").Collection("courses").Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var course Course
		cursor.Decode(&course)
		courses = append(courses, course)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// GetCourse retrieves a specific course by ID
func GetCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var course Course
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := client.Database("coursedb").Collection("courses").FindOne(ctx, bson.M{"_id": params["id"]}).Decode(&course)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

// CreateCourse creates a new course
func CreateCourse(w http.ResponseWriter, r *http.Request) {
	var course Course
	json.NewDecoder(r.Body).Decode(&course)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := client.Database("coursedb").Collection("courses").InsertOne(ctx, course)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.InsertedID)
}

// UpdateCourse updates an existing course by ID
func UpdateCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedCourse Course
	json.NewDecoder(r.Body).Decode(&updatedCourse)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": params["id"]}
	update := bson.M{"$set": updatedCourse}
	_, err := client.Database("coursedb").Collection("courses").UpdateOne(ctx, filter, update)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteCourse deletes a course by ID
func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Database("coursedb").Collection("courses").DeleteOne(ctx, bson.M{"_id": params["id"]})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
