package server

import (
	"backend/internal/repository/postgres"
	"backend/internal/repository/redis"
	"backend/internal/server/handlers"
	"backend/internal/server/router"
	"fmt"
	"log"
	"net/http"
)

func Run() error {
	// initializing storage to mem (later to SQLite)
	repo, err := postgres.NewStorage("postgres://postgres:password@127.0.0.1:5444/metrics_db?sslmode=disable")
	// repo, err := repository.NewSqlStorage("metrics.db")
	if err != nil {
		return fmt.Errorf("failed to create reposql: %w", err)
	}
	cache, err := redis.NewStorage("localhost:6379", "", 0)
	if err != nil {
		return fmt.Errorf("failed to create reporedis: %w", err)
	}

	// init handlers and give them repo
	h := &handlers.Handler{
		Repo:  repo,
		Cache: cache,
	}

	// creating router
	r := router.NewRouter(h)

	log.Println("server is up on :8080")
	// up port for listening
	if err := http.ListenAndServe("127.0.0.1:8080", r); err != nil {
		log.Fatalf("error starting server: %s", err)
	}

	return nil
}
