package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
)

// FormParsingFailedResponse - json response notifying user that form parsing failed.
func FormParsingFailedResponse(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	ErrorResp(w, fmt.Sprintf("failed to parse form: %v", err))
}

// BadFormValuesResponse - sends a json response with an error and a map of field:problem.
func BadFormValuesResponse(w http.ResponseWriter, fieldProblems map[string]string) {
	w.WriteHeader(http.StatusBadRequest)

	resp := make(map[string]any)

	fields := make([]string, 0)
	for fieldName := range fieldProblems {
		fields = append(fields, fieldName)
	}

	resp["error"] = fmt.Sprintf("error parsing fields [%v]", strings.Join(fields, ","))

	resp["fields"] = fieldProblems

	JsonResp(w, resp)
}

func ErrorResp(w http.ResponseWriter, errorMessage string) {
	resp := make(map[string]string)
	resp["error"] = errorMessage
	JsonResp(w, resp)
}

// JsonResp - marshals the provided data sends via the writer
func JsonResp(w http.ResponseWriter, data any) {
	jsonResp, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func SetAuthTokenContext(r *http.Request, token *jwt.Token) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "token", token)
	r = r.WithContext(ctx)
	return r
}

func GetAuthToken(r *http.Request) (*jwt.Token, error) {
	token, ok := r.Context().Value("token").(*jwt.Token)
	if token == nil || !ok {
		return nil, NoAuthTokenError
	}
	return token, nil
}
