package handler

// import (
// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-chi/chi/v5/middleware"
// 	httpSwagger "github.com/swaggo/http-swagger"
// 	"project-management/internal/config"
// 	"project-management/internal/handler/http"
// 	"project-management/internal/service/project"
// 	"project-management/internal/service/task"
// 	"project-management/internal/service/user"
// )

// type Dependencies struct {
// 	Configs        config.Configs
// 	UserService    *user.Service
// 	TaskService    *task.Service
// 	ProjectService *project.Service
// }

// // Configuration is an alias for a function that will take in a pointer to a Handler and modify it
// type Configuration func(h *Handler) error

// // Handler is an implementation of the Handler
// type Handler struct {
// 	dependencies Dependencies
// 	HTTP         *chi.Mux
// }

// // New takes a variable amount of Configuration functions and returns a new Handler
// func New(d Dependencies, configs ...Configuration) (h *Handler, err error) {
// 	h = &Handler{dependencies: d}
// 	for _, cfg := range configs {
// 		if err = cfg(h); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return h, nil
// }

// // WithHTTPHandler applies an HTTP handler to the Handler
// func WithHTTPHandler() Configuration {
// 	return func(h *Handler) error {
// 		// Create the HTTP router
// 		h.HTTP = chi.NewRouter()

// 		// Apply middleware
// 		h.HTTP.Use(middleware.Logger)
// 		h.HTTP.Use(middleware.Recoverer)
// 		h.HTTP.Use(middleware.Timeout(h.dependencies.Configs.APP.Timeout))

// 		// Swagger documentation endpoint
// 		h.HTTP.Get("/swagger/*", httpSwagger.WrapHandler)

// 		// Initialize service handlers
// 		userHandler := http.NewUserHandler(h.dependencies.UserService)
// 		taskHandler := http.NewTaskHandler(h.dependencies.TaskService)
// 		projectHandler := http.NewProjectHandler(h.dependencies.ProjectService)

// 		// Register routes
// 		h.HTTP.Route("/", func(r chi.Router) {
// 			r.Mount("/users", userHandler.Routes())
// 			r.Mount("/tasks", taskHandler.Routes())
// 			r.Mount("/projects", projectHandler.Routes())
// 		})

// 		return nil
// 	}
// }
