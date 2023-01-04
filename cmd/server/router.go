package server

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
	"github.com/solutionstack/jobsity-demo/handlers/auth"
	"net/http"
)

func NewRouter(handlers *auth.AuthHandler) *chi.Mux {
	reqLogger := httplog.NewLogger("req-logger", httplog.Options{
		JSON: true,
	})

	// Service
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(reqLogger))
	r.Use(middleware.Heartbeat("/ping"))

	r.Use(cors.Handler(cors.Options{}))

	r.Post("/auth/signup", handlers.Register)
	r.Post("/auth/login", handlers.Login)

	//html dir
	static := http.FileServer(http.Dir("./static"))
	r.Mount("/", static)

	return r
}
