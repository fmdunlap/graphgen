package database

import (
	"context"
	"graphgen/internal/config"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testDbConfig = &config.DatabaseConfig{}
)

func mustStartPostgresContainer() (func(context.Context) error, *config.DatabaseConfig, error) {
	dbConfig := config.DatabaseConfig{
		Database: "database",
		Password: "password",
		Schema:   "public",
		Username: "user",
	}

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbConfig.Database),
		postgres.WithUsername(dbConfig.Username),
		postgres.WithPassword(dbConfig.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, nil, err
	}

	dbConfig.Host, err = dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, nil, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, nil, err
	}
	dbConfig.Port = dbPort.Int()

	return dbContainer.Terminate, &dbConfig, err
}

func TestMain(m *testing.M) {
	teardown, dbConfig, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	testDbConfig = dbConfig

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
}

func TestNew(t *testing.T) {
	srv := New(testDbConfig)
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New(testDbConfig)

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	srv := New(testDbConfig)

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}
