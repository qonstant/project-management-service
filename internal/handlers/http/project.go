package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"project-management-service/db/sqlc"
	"project-management-service/pkg/server/response"
)

type ProjectHandler struct {
	db *db.Queries
}

func NewProjectHandler(conn *sql.DB) *ProjectHandler {
	return &ProjectHandler{
		db: db.New(conn),
	}
}

func (h *ProjectHandler) Routes() chi.Router {
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

func (h *ProjectHandler) list(w http.ResponseWriter, r *http.Request) {
	projects, err := h.db.ListProjects(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
	response.OK(w, r, projects)
}

func (h *ProjectHandler) add(w http.ResponseWriter, r *http.Request) {
	var req db.CreateProjectParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	project, err := h.db.CreateProject(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, project)
}

func (h *ProjectHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	project, err := h.db.GetProject(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, project)
}

func (h *ProjectHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	var req db.UpdateProjectParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	req.ID = id

	project, err := h.db.UpdateProject(r.Context(), req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, project)
}

func (h *ProjectHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	if err := h.db.DeleteProject(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.NoContent(w, r)
}
