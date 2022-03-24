package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lauslim12/basic"
)

// Response to send back JSON response to the client.
type Response struct {
	Status     string        `json:"status"`
	StatusCode int           `json:"statusCode"`
	StatusText string        `json:"statusText"`
	Message    string        `json:"message"`
	Data       []interface{} `json:"data,omitempty"`
	Error      string        `json:"error,omitempty"`
	ErrorCode  string        `json:"errorCode,omitempty"`
}

// NewSuccessJSON creates a new, successful JSON response.
func NewSuccessJSON(statusCode int, message string, data []interface{}) Response {
	return Response{
		Status:     "success",
		StatusCode: statusCode,
		StatusText: http.StatusText(statusCode),
		Message:    message,
		Data:       data,
	}
}

// NewFailureJSON creates a new, failed JSON response.
func NewFailureJSON(statusCode int, message string, err string, errorCode string) Response {
	return Response{
		Status:     "fail",
		StatusCode: statusCode,
		StatusText: http.StatusText(statusCode),
		Message:    message,
		Error:      err,
		ErrorCode:  errorCode,
	}
}

// SendSuccess sends back a successful response to the client.
func SendSuccess(w http.ResponseWriter, code int, message string, data []interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(NewSuccessJSON(code, message, data))
}

// SendFailure sends back a failure response to the client.
func SendFailure(w http.ResponseWriter, code int, err string, errorCode string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(NewFailureJSON(code, "Error!", err, errorCode))
}

// Special simple auth route for authenticated users.
func simpleAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the private endpoint!"))
}

// Special complex auth route for authenticated users.
func complexAuth(w http.ResponseWriter, r *http.Request) {
	SendSuccess(w, http.StatusOK, "Successfully authenticated!", make([]interface{}, 0))
}

// An example middleware.
func middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware says hello!")
		next.ServeHTTP(w, r)
	}
}

// Driver code.
func main() {
	// Sample users.
	users := map[string]string{"gerysantoso": "gerysantoso"}

	// Create simple configuration for Basic Authentication.
	basicAuthSimple := basic.NewCustomBasicAuth(nil, "UTF-8", nil, nil, "Private", users)

	// Create complex configuration for Basic Authentication.
	basicAuthComplex := basic.NewCustomBasicAuth(nil, "UTF-8", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SendFailure(w, http.StatusUnauthorized, "Invalid authentication scheme!", "BAUTH: E0001")
	}), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SendFailure(w, http.StatusUnauthorized, "Invalid username and/or password!", "BAUTH: E0002")
	}), "Secret", users)

	// Hello endpoint!
	http.HandleFunc("/simple", basicAuthSimple.Authenticate(simpleAuth))
	http.HandleFunc("/complex", basicAuthComplex.Authenticate(complexAuth))
	http.HandleFunc("/middleware", middleware(basicAuthSimple.Authenticate(simpleAuth)))

	// Listen and serve.
	log.Println("Golang server powered by 'net/http' is listening at port 5000.")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
