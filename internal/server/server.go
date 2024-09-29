package server

import (
	"context"
	"fmt"
	"graphgen/internal/auth"
	"graphgen/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"graphgen/internal/database"
)

type Server struct {
	config *config.Config

	httpServer *http.Server

	auth auth.Service
	db   database.Service
}

func NewServer(config *config.Config) *Server {
	newServer := &Server{
		config: config,

		auth: auth.New(&config.Auth),
		db:   database.New(&config.Database),
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%v:%d", config.Server.Address, config.Server.Port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	newServer.httpServer = httpServer

	return newServer
}

func (s *Server) StartAndBlock() {
	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), s.config.Server.GetShutdownTimeoutDuration())
	defer cancel()
	_ = s.httpServer.Shutdown(ctx)
}
