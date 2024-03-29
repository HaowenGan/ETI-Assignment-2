// Ong Jia Yuan / S10227735B
// /user-service/main.go

package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("secret-key-replace-with-your-own"))

func init() {
	// Set SameSite and other cookie attributes here.
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode, // or http.SameSiteLaxMode if you want less strict settings
		Secure:   true,                    // set to true if you're using https
	}
}

// CORS Middleware
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("CORS Middleware hit")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// User structure
type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Usertype  string `json:"usertype"`
}

var db *sql.DB
var err error

// Establishes a connection to the MySQL database
func connectDatabase() {
	db, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/eti_asg2")
	//db, err = sql.Open("mysql", "user:password@tcp(host.docker.internal:3306)/eti_asg2")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

// Create a new user in the database
func registerUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Explicitly set Usertype to "student"
	user.Usertype = "student"

	// Hashing password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users(firstName, lastName, email, username, password, usertype) VALUES(?, ?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.Username, string(hashedPassword), user.Usertype)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// loginUser authenticates a user
func loginUser(w http.ResponseWriter, r *http.Request) {
	var user User
	var hashedPassword string
	var id int

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//err = db.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&hashedPassword)
	err = db.QueryRow("SELECT id, firstName, lastName, email, password, usertype FROM users WHERE username = ?", user.Username).Scan(&id, &user.FirstName, &user.LastName, &user.Email, &hashedPassword, &user.Usertype)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)) == nil {
		// Create a new session and set the session values
		session, _ := store.Get(r, "user-session")
		session.Values["authenticated"] = true
		session.Values["userID"] = id
		session.Values["username"] = user.Username
		session.Values["firstName"] = user.FirstName
		session.Values["lastName"] = user.LastName
		session.Values["email"] = user.Email
		session.Values["usertype"] = user.Usertype
		session.Save(r, w)

		// Save the session before writing to the response
		if err := session.Save(r, w); err != nil {
			log.Printf("Error saving session: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Println("Session saved for user:", user.Username)

		log.Printf("Session saved with userID: %d", id)
		// Send a response to the client that authentication was successful
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User logged in successfully"))
	} else {
		// Handle invalid credentials
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
}

// Endpoint to get the current user's details
func getCurrentUser(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Error retrieving session", http.StatusInternalServerError)
		log.Printf("Error retrieving session: %v", err)
		return
	}

	// Check if the user is authenticated
	auth, ok := session.Values["authenticated"].(bool)
	if !ok || !auth {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Retrieve user details from the session
	firstName, ok1 := session.Values["firstName"].(string)
	lastName, ok2 := session.Values["lastName"].(string)
	email, ok3 := session.Values["email"].(string)
	username, ok4 := session.Values["username"].(string)
	usertype, ok5 := session.Values["usertype"].(string)

	// If any of the type assertions failed, handle the error
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		http.Error(w, "Error retrieving user details from session", http.StatusInternalServerError)
		return
	}

	// Create a map to hold the user details
	userDetails := map[string]string{
		"firstName": firstName,
		"lastName":  lastName,
		"email":     email,
		"username":  username,
		"usertype":  usertype,
	}

	// Set the header and encode the userDetails map to JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userDetails); err != nil {
		http.Error(w, "Error encoding user details to JSON", http.StatusInternalServerError)
		log.Printf("Error encoding user details to JSON: %v", err)
	}
}

// logoutUser logs a user out by destroying the session
func logoutUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")

	// Revoke users authentication
	session.Values["authenticated"] = false
	delete(session.Values, "userID")
	delete(session.Values, "username")
	delete(session.Values, "firstName")
	delete(session.Values, "lastName")
	delete(session.Values, "email")
	delete(session.Values, "usertype")

	session.Options.MaxAge = -1 // Immediately delete the session

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged out successfully"))
}

// Update the current user's profile
func updateUserProfile(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil || !session.Values["authenticated"].(bool) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := session.Values["userID"].(int)

	// Build SQL query dynamically based on provided fields
	query := "UPDATE users SET "
	queryParams := []interface{}{}
	if user.FirstName != "" {
		query += "firstName = ?, "
		queryParams = append(queryParams, user.FirstName)
	}
	if user.LastName != "" {
		query += "lastName = ?, "
		queryParams = append(queryParams, user.LastName)
	}
	if user.Email != "" {
		query += "email = ?, "
		queryParams = append(queryParams, user.Email)
	}

	// Remove trailing comma and space
	query = query[:len(query)-2]
	query += " WHERE id = ?"
	queryParams = append(queryParams, userID)

	_, err = db.Exec(query, queryParams...)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Profile updated successfully"))
}

func changeUserPassword(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil || !session.Values["authenticated"].(bool) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := session.Values["userID"].(int)
	var currentHashedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&currentHashedPassword)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify current password
	if bcrypt.CompareHashAndPassword([]byte(currentHashedPassword), []byte(payload.CurrentPassword)) != nil {
		http.Error(w, "Invalid current password", http.StatusUnauthorized)
		return
	}

	// Hash new password
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Update password in database
	_, err = db.Exec("UPDATE users SET password = ? WHERE id = ?", newHashedPassword, userID)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password updated successfully"))
}

func main() {
	connectDatabase()
	defer db.Close()
	workDir, _ := os.Getwd()
	staticDir := filepath.Join(workDir, "../") // The parent directory now

	// In Docker, set the static directory directly
	//staticDir := "/root/parent/front-end"

	// Routing
	router := mux.NewRouter()
	// Apply the CORS middleware to the router
	router.Use(CORSMiddleware)

	/*
		// FOR DOCKER
		jsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/javascript")
			http.ServeFile(w, r, filepath.Join(staticDir, "js", r.URL.Path))
		})
		router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", jsHandler))
		cssHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/css")
			http.ServeFile(w, r, filepath.Join(staticDir, "css", r.URL.Path))
		})
		router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", cssHandler))
	*/

	// Custom handler for JavaScript files to ensure the correct Content-Type
	jsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, filepath.Join(staticDir, "front-end/js/", r.URL.Path))
	})
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", jsHandler))

	// Custom handler for CSS files to ensure the correct Content-Type
	cssHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, filepath.Join(staticDir, "front-end/css/", r.URL.Path))
	})
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", cssHandler))

	// HTML files directly from the root of the server
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})

	router.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})
	router.HandleFunc("/dashboard.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "dashboard.html"))
	})

	// User related
	router.HandleFunc("/profile.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "profile.html"))
	})
	router.HandleFunc("/settings.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "settings.html"))
	})
	router.HandleFunc("/register.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "register.html"))
	})
	router.HandleFunc("/login.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "login.html"))
	})

	// Review related
	router.HandleFunc("/MakeReview.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "MakeReview.html"))
	})
	router.HandleFunc("/ViewReviews.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "ViewReviews.html"))
	})
	router.HandleFunc("/EditReview.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "EditReview.html"))
	})

	// Courses related
	router.HandleFunc("/courses.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "courses.html"))
	})
	router.HandleFunc("/Enrollment.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "Enrollment.html"))
	})

	// Payment related
	router.HandleFunc("/payment.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "payment.html"))
	})

	// API routes
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/register", registerUser).Methods("POST")
	apiRouter.HandleFunc("/login", loginUser).Methods("POST")
	apiRouter.HandleFunc("/current-user", getCurrentUser).Methods("GET")
	apiRouter.HandleFunc("/logout", logoutUser).Methods("POST")
	apiRouter.HandleFunc("/update-profile", updateUserProfile).Methods("POST")
	apiRouter.HandleFunc("/change-password", changeUserPassword).Methods("POST")

	log.Fatal(http.ListenAndServe(":5000", router))
}
