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

type NullableTime struct {
	Time  time.Time
	Valid bool
}

func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	var aux interface{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch v := aux.(type) {
	case string:
		if v == "" {
			nt.Valid = false
			return nil
		}
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		nt.Time = t
		nt.Valid = true
	case map[string]interface{}:
		if v["time"] == nil {
			nt.Valid = false
			return nil
		}
		timeStr, ok := v["time"].(string)
		if !ok {
			return errors.New("invalid format for NullableTime")
		}
		t, err := time.Parse("2006-01-02", timeStr)
		if err != nil {
			return err
		}
		nt.Time = t
		nt.Valid = v["valid"].(bool)
	case nil:
		nt.Valid = false
	default:
		return errors.New("invalid format for NullableTime")
	}

	return nil
}

type createTaskRequest struct {
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	Priority       string       `json:"priority"`
	Status         string       `json:"status"`
	AssigneeID     int64        `json:"assignee_id"`
	ProjectID      int64        `json:"project_id"`
	CompletionDate NullableTime `json:"completion_date"`
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
// @Param		request	body	createTaskRequest	true	"Task details"
// @Success	200		{object}	db.Task
// @Failure	400		{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router		/tasks [post]
func (h *TaskHandler) add(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	params := db.CreateTaskParams{
		Title:         req.Title,
		Description:   req.Description,
		Priority:      db.TaskPriority(req.Priority),
		Status:        db.TaskStatus(req.Status),
		AssigneeID:    req.AssigneeID,
		ProjectID:     req.ProjectID,
		CompletionDate: sql.NullTime{Time: req.CompletionDate.Time, Valid: req.CompletionDate.Valid},
	}

	task, err := h.db.CreateTask(r.Context(), params)
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
