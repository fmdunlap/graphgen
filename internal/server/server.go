package server

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"graphgen/internal/database"
)

type Server struct {
	address string
	port    int

	db database.Service
}

func NewServer() *http.Server {
	address := viper.GetString("address")
	port := viper.GetInt("port")
	NewServer := &Server{
		address: address,
		port:    port,

		db: database.New(),
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
