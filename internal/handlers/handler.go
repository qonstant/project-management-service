package handler

import (
	"database/sql"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hellofresh/health-go/v5"
	healthPg "github.com/hellofresh/health-go/v5/checks/postgres"
	httpSwagger "github.com/swaggo/http-swagger"

	"project-management-service/docs"
	"project-management-service/internal/config"
	"project-management-service/internal/handlers/http"
)

type Dependencies struct {
	DB      *sql.DB
	Configs config.Config
}

// Configuration is an alias for a function that modifies the Handler
type Configuration func(h *Handler) error

// Handler is an implementation of the Handler
type Handler struct {
	dependencies Dependencies
	HTTP         *chi.Mux
}

// New creates a new Handler
func New(d Dependencies, configs ...Configuration) (h *Handler, err error) {
	// Create the handler
	h = &Handler{
		dependencies: d,
	}

	// Apply all Configurations passed in
	for _, cfg := range configs {
		if err = cfg(h); err != nil {
			return
		}
	}

	return
}

// WithHTTPHandler applies an HTTP handler to the Handler
func WithHTTPHandler() Configuration {
	return func(h *Handler) error {
		// Create the HTTP handler
		h.HTTP = chi.NewRouter()

		// Init swagger handler
		docs.SwaggerInfo.BasePath = h.dependencies.Configs.BaseURL
		h.HTTP.Get("/swagger/*", httpSwagger.WrapHandler)

		// Init service handlers
		userHandler := http.NewUserHandler(h.dependencies.DB)
		projectHandler := http.NewProjectHandler(h.dependencies.DB)
		taskHandler := http.NewTaskHandler(h.dependencies.DB)

		h.HTTP.Route("/", func(r chi.Router) {
			r.Mount("/users", userHandler.Routes())
			r.Mount("/projects", projectHandler.Routes())
			r.Mount("/tasks", taskHandler.Routes())
		})

		// Set up health checks
		healthHandler, _ := health.New(health.WithComponent(health.Component{
			Name:    "project-management-service",
			Version: "v1.0",
		}), health.WithChecks(
			health.Config{
				Name:      "postgres",
				Timeout:   time.Second * 10,
				SkipOnErr: false,
				Check: healthPg.New(healthPg.Config{
					DSN: os.Getenv("DB_SOURCE"),
				}),
			},
		))

		// Register health check endpoint
		h.HTTP.Get("/status", healthHandler.HandlerFunc)

		return nil
	}
}
