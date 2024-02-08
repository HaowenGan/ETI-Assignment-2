package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/charge"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/sub"
)

type Course struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Price    int64  `json:"price"`
	Currency string `json:"currency"`
}

type Enrollment struct {
	CourseID   string `json:"course_id"`
	CustomerID string `json:"customer_id"`
}

var courses []Course
var enrollments []Enrollment

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/charge", handleCharge).Methods("POST")
	r.HandleFunc("/payment/{id}", getPaymentDetails).Methods("GET")
	r.HandleFunc("/webhook", handleWebhook).Methods("POST")
	r.HandleFunc("/customer", createCustomer).Methods("POST")
	r.HandleFunc("/subscription", createSubscription).Methods("POST")
	r.HandleFunc("/courses", getCourses).Methods("GET")
	r.HandleFunc("/courses", createCourse).Methods("POST")
	r.HandleFunc("/enroll", enrollInCourse).Methods("POST")
	r.HandleFunc("/purchases", getUserPurchases).Methods("GET")

	authMiddleware := authMiddleware()
	r.Use(authMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func handleCharge(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Parse the request body
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the token from the request
	token := r.Form.Get("stripeToken")

	// Create a charge using the token
	params := &stripe.ChargeParams{
		Amount:   stripe.Int64(1000), // Amount in cents
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Source:   &stripe.SourceParams{Token: stripe.String(token)},
	}

	charge, err := charge.New(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Associate the purchase with the course
	courseID := r.Form.Get("course_id")
	customerID := r.Form.Get("customer_id")
	enrollments = append(enrollments, Enrollment{CourseID: courseID, CustomerID: customerID})

	// Return success response with charge ID
	response := map[string]string{"status": "success", "charge_id": charge.ID}
	json.NewEncoder(w).Encode(response)
}

func getPaymentDetails(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Get the payment ID from the request path parameters
	vars := mux.Vars(r)
	paymentID := vars["id"]

	// Retrieve the charge details from Stripe
	charge, err := charge.Get(paymentID, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return payment details
	json.NewEncoder(w).Encode(charge)
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming webhook payload
	payload, err := stripe.ConstructEvent(r.Body, r.Header.Get("Stripe-Signature"), os.Getenv("STRIPE_WEBHOOK_SECRET"))
	if err != nil {
		log.Printf("Error parsing webhook payload: %s\n", err)
		http.Error(w, "Webhook Error", http.StatusBadRequest)
		return
	}

	// Handle the event based on its type
	switch payload.Type {
	case "charge.failed":
		// Handle failed charge event
		log.Println("Charge failed:", payload.Data.Object)
	case "charge.succeeded":
		// Handle successful charge event
		log.Println("Charge succeeded:", payload.Data.Object)
	default:
		log.Printf("Unhandled event type: %s\n", payload.Type)
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Parse the request body
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a customer
	customerParams := &stripe.CustomerParams{
		Name:  stripe.String(r.Form.Get("name")),
		Email: stripe.String(r.Form.Get("email")),
	}

	cust, err := customer.New(customerParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response with customer ID
	response := map[string]string{"status": "success", "customer_id": cust.ID}
	json.NewEncoder(w).Encode(response)
}

func createSubscription(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("sk_test_51KUQmKFmWCKB31rkHzO9NyfXjCvC2piQASkRfaoTHZwc3jSBOo504yqcQ7EFfjuJNO36Hjdd1cRow4ELEr90tx7D00ztlS6wm4")

	// Parse the request body
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a subscription
	subParams := &stripe.SubscriptionParams{
		Customer: stripe.String(r.Form.Get("customer_id")),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(r.Form.Get("price_id")),
			},
		},
	}

	_, err = sub.New(subParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]string{"status": "success"}
	json.NewEncoder(w).Encode(response)
}

func getCourses(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(courses)
}

func createCourse(w http.ResponseWriter, r *http.Request) {
	var newCourse Course
	err := json.NewDecoder(r.Body).Decode(&newCourse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	courses = append(courses, newCourse)

	json.NewEncoder(w).Encode(newCourse)
}

func enrollInCourse(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var enrollment Enrollment
	err := json.NewDecoder(r.Body).Decode(&enrollment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the course exists
	courseID := enrollment.CourseID
	var foundCourse bool
	for _, course := range courses {
		if course.ID == courseID {
			foundCourse = true
			break
		}
	}
	if !foundCourse {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	// Associate the enrollment with the course
	enrollments = append(enrollments, enrollment)

	// Return success response
	response := map[string]string{"status": "success"}
	json.NewEncoder(w).Encode(response)
}

func getUserPurchases(w http.ResponseWriter, r *http.Request) {
	// Get user ID from authentication
	userID := r.Context().Value("user_id").(string)

	// Find purchases associated with the user
	var userPurchases []Course
	for _, enrollment := range enrollments {
		if enrollment.CustomerID == userID {
			for _, course := range courses {
				if course.ID == enrollment.CourseID {
					userPurchases = append(userPurchases, course)
				}
			}
		}
	}

	// Return user purchases
	json.NewEncoder(w).Encode(userPurchases)
}

func authMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Mock authentication logic (validate token)
			userID, err := authenticateToken(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Add user ID to request context
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func authenticateToken(token string) (string, error) {
	// Mock authentication logic
	if token == "valid_token" {
		return "user_id", nil
	}
	return "", errors.New("invalid token")
}
