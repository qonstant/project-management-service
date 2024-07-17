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

type TaskHandler struct {
	db *db.Queries
}

func NewTaskHandler(conn *sql.DB) *TaskHandler {
	return &TaskHandler{
		db: db.New(conn),
	}
}

func (h *TaskHandler) Routes() chi.Router {
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

// @Summary	List of tasks from the repository
// @Tags		tasks
// @Accept		json
// @Produce	json
// @Success	200	{array}		db.Task
// @Failure	500	{object}	response.Object
// @Router		/tasks [get]
func (h *TaskHandler) list(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.db.ListTasks(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
	response.OK(w, r, tasks)
}

// @Summary	Add a new task to the repository
// @Tags		tasks
// @Accept		json
// @Produce	json
// @Param		request	body		db.CreateTaskParams	true	"Task details"
// @Success	200		{object}	db.Task
// @Failure	400		{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router		/tasks [post]
func (h *TaskHandler) add(w http.ResponseWriter, r *http.Request) {
	var req db.CreateTaskParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	task, err := h.db.CreateTask(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, task)
}

// @Summary	Get a task from the repository
// @Tags		tasks
// @Accept		json
// @Produce	json
// @Param		id	path		int	true	"Task ID"
// @Success	200	{object}	db.Task
// @Failure	404	{object}	response.Object
// @Failure	500	{object}	response.Object
// @Router		/tasks/{id} [get]
func (h *TaskHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	task, err := h.db.GetTask(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, task)
}

// @Summary	Update a task in the repository
// @Tags		tasks
// @Accept		json
// @Produce	json
// @Param		id		path		int						true	"Task ID"
// @Param		request	body		db.UpdateTaskParams	true	"Task details"
// @Success	200		{object}	db.Task
// @Failure	400		{object}	response.Object
// @Failure	404		{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router		/tasks/{id} [put]
func (h *TaskHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	var req db.UpdateTaskParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	req.ID = id

	task, err := h.db.UpdateTask(r.Context(), req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.OK(w, r, task)
}

// @Summary	Delete a task from the repository
// @Tags		tasks
// @Accept		json
// @Produce	json
// @Param		id	path		int	true	"Task ID"
// @Success	204	{object}	response.Object
// @Failure	404	{object}	response.Object
// @Failure	500	{object}	response.Object
// @Router		/tasks/{id} [delete]
func (h *TaskHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	if err := h.db.DeleteTask(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(w, r, err)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	response.NoContent(w, r)
}
