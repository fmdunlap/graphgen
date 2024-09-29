package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) registerAuthRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/auth").Subrouter()
	subrouter.HandleFunc("/new-token", s.issueTokenHandler).Methods(http.MethodPost)

	verifyRouter := subrouter.PathPrefix("/verify").Subrouter()
	verifyRouter.Use(s.authMiddleware)
	verifyRouter.HandleFunc("", s.verifyTokenHandler)
}

func (s *Server) issueTokenHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		FormParsingFailedResponse(w, err)
		return
	}
	username := r.Form.Get("username")
	if username == "" {
		BadFormValuesResponse(w, map[string]string{
			"username": "required",
		})
		return
	}

	token, err := s.auth.CreateToken(username)
	if err != nil {
		ErrorResp(w, fmt.Sprintf("failed to create token: %v", err))
	}

	w.WriteHeader(http.StatusOK)
	JsonResp(w, map[string]string{
		"token": token,
	})
}

func (s *Server) verifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		FormParsingFailedResponse(w, err)
		return
	}

	user, err := GetAuthenticatedUser(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorResp(w, "no auth token found for request")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hello %v!", user)))
}
