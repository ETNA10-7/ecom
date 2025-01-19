package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

// Function to CHeck whether the body is empty
func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("MIssing Request Body")
	}
	// json
	return json.NewDecoder(r.Body).Decode(v)
	//This only fills User struct and returns nil
	//To the err object that called this func
}

// Function to write JSON Payload to user
// type any is interface
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

// func GetTokenFromRequest(r *http.Request) string {
// 	tokenAuth := r.Header.Get("Authorization")
// 	//tokenQuery := r.URL.Query().Get("token")

// 	if tokenAuth != "" {
// 		return tokenAuth
// 	}
// 	// if tokenQuery != "" {
// 	// 	return tokenQuery
// 	// }

// 	return ""
// }

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth != "" && strings.HasPrefix(tokenAuth, "Bearer ") {
		return strings.TrimPrefix(tokenAuth, "Bearer ")
	}
	return ""
}
