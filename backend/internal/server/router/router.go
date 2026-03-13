package router

import (
	_ "backend/docs"
	"backend/internal/server/handlers"
	"net/http"

	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
	// creating router
	r := chi.NewRouter()

	// CORS for React to connect
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://localhost:3000",
			"http://localhost:8080",
		}, // порты фронтенда
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
	}))

	// adding middlewares
	r.Use(middleware.Logger)    // writes every request to console
	r.Use(middleware.Recoverer) // keeps server up

	r.Get("/", h.Main)
	// for docs
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	// when agent posts, dumpStat is called
	r.Get("/stat", h.GetStat)
	r.Post("/update", h.DumpStat)
	// test router to check in browser
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return r
}
