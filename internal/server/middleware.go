package server

import (
	"fmt"
	"log"
	"net/http"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%v] %v", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if len(authorizationHeader) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			ErrorResp(w, "no authorization token provided")
			return
		}

		tokenString := authorizationHeader[len("Bearer "):]
		_, claims, err := s.auth.VerifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			ErrorResp(w, fmt.Sprintf("error verifying token: %v", err))
			return
		}

		r = SetAuthTokenContext(r, claims)

		next.ServeHTTP(w, r)
	})
}
