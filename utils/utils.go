package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	tireappErrors "github.com/nathaniel-alvin/tireappBE/error"
)

type ApiResponse struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var Validate = validator.New()

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if r.Body == nil {
		return v, fmt.Errorf("missing request body")
	}

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func Encode(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	response := &ApiResponse{
		Data: v,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func WriteError(w http.ResponseWriter, err error) error {
	code, message := tireappErrors.ErrorCode(err), tireappErrors.ErrorMessage(err)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(tireappErrors.ErrorStatusCode(code))

	response := &ErrorResponse{
		Error: message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	return nil
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}
