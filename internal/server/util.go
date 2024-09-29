package server

import (
	"context"
	"encoding/json"
	"fmt"
	"graphgen/internal/auth"
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

func SetAuthTokenContext(r *http.Request, claims *auth.JWTClaims) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "claims", claims)
	r = r.WithContext(ctx)
	return r
}

func GetAuthenticatedUser(r *http.Request) (string, error) {
	claims, ok := r.Context().Value("claims").(*auth.JWTClaims)
	if claims == nil || !ok {
		return "", NoAuthTokenError
	}
	subject, err := claims.GetSubject()
	if err != nil {
		return "", NoAuthTokenError
	}
	return subject, nil
}
