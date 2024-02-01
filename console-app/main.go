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

func main() {
	var choice string
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("User Management Console Application")
	fmt.Println("1. Register")
	fmt.Println("2. Login")
	fmt.Print("Enter choice: ")
	choice, _ = reader.ReadString('\n')
	choice = strings.TrimSpace(choice) // Trims the newline character from the input

	switch choice {
	case "1":
		register()
	case "2":
		login()
	default:
		fmt.Println("Invalid choice")
	}
}
