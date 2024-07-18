package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

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

type createProjectRequest struct {
	Description string `json:"description"`
	EndDate     string `json:"end_date"`
	ManagerID   int64  `json:"manager_id"`
	Name        string `json:"name"`
	StartDate   string `json:"start_date"`
}

// @Summary	List of projects from the repository
// @Tags		projects
// @Accept		json
// @Produce	json
// @Success	200	{array}		db.Project
// @Failure	500	{object}	response.Object
// @Router		/projects [get]
func (h *ProjectHandler) list(w http.ResponseWriter, r *http.Request) {
	projects, err := h.db.ListProjects(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
	response.OK(w, r, projects)
}

// @Summary	Add a new project to the repository
// @Tags		projects
// @Accept		json
// @Produce	json
// @Param		request	body	createProjectRequest	true	"Project details"
// @Success	200		{object}	db.Project
// @Failure	400		{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router		/projects [post]
func (h *ProjectHandler) add(w http.ResponseWriter, r *http.Request) {
	var req createProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	params := db.CreateProjectParams{
		Name:        req.Name,
		Description: req.Description,
		StartDate:   startDate,
		EndDate:     endDate,
		ManagerID:   req.ManagerID,
	}

	project, err := h.db.CreateProject(r.Context(), params)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, project)
}

// @Summary	Get a project from the repository
// @Tags		projects
// @Accept		json
// @Produce	json
// @Param		id	path		int	true	"Project ID"
// @Success	200	{object}	db.Project
// @Failure	404	{object}	response.Object
// @Failure	500	{object}	response.Object
// @Router		/projects/{id} [get]
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

// @Summary	Update a project in the repository
// @Tags		projects
// @Accept		json
// @Produce	json
// @Param		id		path		int						true	"Project ID"
// @Param		request	body		db.UpdateProjectParams	true	"Project details"
// @Success	200		{object}	db.Project
// @Failure	400		{object}	response.Object
// @Failure	404		{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router		/projects/{id} [put]
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

// @Summary	Delete a project from the repository
// @Tags		projects
// @Accept		json
// @Produce	json
// @Param		id	path		int	true	"Project ID"
// @Success	204	{object}	response.Object
// @Failure	404	{object}	response.Object
// @Failure	500	{object}	response.Object
// @Router		/projects/{id} [delete]
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
