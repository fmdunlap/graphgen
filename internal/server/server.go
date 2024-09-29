package server

import (
	"fmt"
	"graphgen/internal/config"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"graphgen/internal/database"
)

type Server struct {
	address string
	port    int
	db      database.Service
}

func NewServer(config *config.Config) *http.Server {
	NewServer := &Server{
		address: config.Server.Address,
		port:    config.Server.Port,

		db: database.New(&config.Database),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf("%v:%d", NewServer.address, NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
