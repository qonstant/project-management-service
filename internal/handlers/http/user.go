package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/render"
	"project-management-service/db/sqlc"
	"project-management-service/pkg/server/response"
)

type UserHandler struct {
	db *db.Queries
}

func NewUserHandler(conn *sql.DB) *UserHandler {
	return &UserHandler{
		db: db.New(conn), // Use the correct constructor from your db package
	}
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.add)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
	})

	return r
}

// list, add, get, update, and delete functions stay the same...

func (h *UserHandler) list(w http.ResponseWriter, r *http.Request) {
	users, err := h.db.ListUsers(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
	response.OK(w, r, users)
}

func (h *UserHandler) add(w http.ResponseWriter, r *http.Request) {
	var req db.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	user, err := h.db.CreateUser(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, user)
}

func (h *UserHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	user, err := h.db.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, user)
}

func (h *UserHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	var req db.UpdateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	req.ID = id

	user, err := h.db.UpdateUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, user)
}

func (h *UserHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	if err := h.db.DeleteUser(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.NoContent(w, r)
}
