// Ong Jia Yuan / S10227735B
// /console-app/main.go

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// User struct to match the one used in the microservice
type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type Review struct {
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	CourseID int    `json:"courseId"`
	Rating   int    `json:"rating"`
	Comment  string `json:"comment"`
}

// helper function to make an HTTP POST request
func makePostRequest(url string, data []byte) error {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
	return nil
}

// register function to simulate user registration
func register() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter first name: ")
	firstName, _ := reader.ReadString('\n')

	fmt.Print("Enter last name: ")
	lastName, _ := reader.ReadString('\n')

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')

	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')

	user := User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Username:  username,
		Password:  password,
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	err = makePostRequest("http://localhost:5000/register", userJson)
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
}

// login function to simulate user login
func login() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')

	user := User{
		Username: username,
		Password: password,
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	err = makePostRequest("http://localhost:5000/login", userJson)
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
}

// addReview function to simulate adding a review
func addReview() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	var username string
	fmt.Scan(&username)
	fmt.Scanln() // Read the newline character left in the buffer

	fmt.Print("Enter course ID: ")
	var courseID int
	fmt.Scan(&courseID)
	fmt.Scanln() // Read the newline character left in the buffer

	fmt.Print("Enter rating (1-5): ")
	var rating int
	fmt.Scan(&rating)
	fmt.Scanln() // Read the newline character left in the buffer

	fmt.Print("Enter comment: ")
	comment, _ := reader.ReadString('\n')

	review := Review{
		Username: username,
		CourseID: courseID,
		Rating:   rating,
		Comment:  strings.TrimSpace(comment),
	}

	reviewJSON, err := json.Marshal(review)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	err = makePostRequest("http://localhost:8080/api/submit-review", reviewJSON)
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
}

func main() {
	var choice string
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("User Management Console Application")
	fmt.Println("1. Register")
	fmt.Println("2. Login")
	fmt.Println("3. Add Review (Not working as of yet)")
	fmt.Print("Enter choice: ")
	choice, _ = reader.ReadString('\n')
	choice = strings.TrimSpace(choice) // Trims the newline character from the input

	switch choice {
	case "1":
		register()
	case "2":
		login()
	case "3":
		addReview()
	default:
		fmt.Println("Invalid choice")
	}
}
