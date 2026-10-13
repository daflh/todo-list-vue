package todo

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:         300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/todo", func(r chi.Router) {
		// GET /todo/getAll - List all todos
		r.Get("/getAll", h.list)

		// POST /todo/create - Create a new todo
		r.Post("/create", h.create)

		// PUT /todo/update/:id - Update a todo by ID
		r.Put("/update/{id}", h.update)

		// PUT /todo/update/:id/status - Update a todo's status by ID
		r.Put("/update/{id}/status", h.updateStatus)

		// DELETE /todo/delete - Delete todos by IDs
		r.Delete("/delete", h.delete)
	})

	return r
}
