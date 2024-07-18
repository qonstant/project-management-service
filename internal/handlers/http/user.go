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

type UserHandler struct {
	db *db.Queries
}

func NewUserHandler(conn *sql.DB) *UserHandler {
	return &UserHandler{
		db: db.New(conn),
	}
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.add)
	r.Get("/search", h.search)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
		r.Get("/tasks", h.getTasks)
	})

	return r
}

// @Summary	List of users from the repository
// @Tags		users
// @Accept		json
// @Produce	json
// @Success	200	{array}		db.User
// @Failure	500	{object}	response.Object
// @Router		/users [get]
func (h *UserHandler) list(w http.ResponseWriter, r *http.Request) {
	users, err := h.db.ListUsers(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}
	response.OK(w, r, users)
}

// @Summary	Add a new user to the repository
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		request	body		db.CreateUserParams	true	"User details"
// @Success	200		{object}	db.User
// @Failure	400		{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router		/users [post]
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

// @Summary	Get a user from the repository
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id	path		int	true	"User ID"
// @Success	200	{object}	db.User
// @Failure	404	{object}	response.Object
// @Failure	500	{object}	response.Object
// @Router		/users/{id} [get]
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

// @Summary	Update a user in the repository
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id		path		int					true	"User ID"
// @Param		request	body		db.UpdateUserParams	true	"User details"
// @Success	200		{object}	db.User
// @Failure	400		{object}	response.Object
// @Failure	404		{object}	response.Object
// @Failure	500	{object}	response.Object
// @Router		/users/{id} [put]
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

// @Summary	Delete a user from the repository
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id	path		int	true	"User ID"
// @Success	204	{object}	response.Object
// @Failure	404	{object}	response.Object
// @Failure	500	{object}	response.Object
// @Router		/users/{id} [delete]
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

// @Summary	Get tasks for a specific user
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		id	path		int	true	"User ID"
// @Success	200	{array}		db.Task
// @Failure	404	{object}	response.Object
// @Failure	500	{object}	response.Object
// @Router		/users/{id}/tasks [get]
func (h *UserHandler) getTasks(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	tasks, err := h.db.GetUserTasks(r.Context(), id)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, tasks)
}

// @Summary	Search users by name or email
// @Tags		users
// @Accept		json
// @Produce	json
// @Param		name	query		string	false	"User name"
// @Param		email	query		string	false	"User email"
// @Success	200	{array}		db.User
// @Failure	500	{object}	response.Object
// @Router		/users/search [get]
func (h *UserHandler) search(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")

	var users []db.User
	var err error

	switch {
	case name != "":
		users, err = h.db.SearchUsersByName(r.Context(), sql.NullString{String: name, Valid: true})
	case email != "":
		users, err = h.db.SearchUsersByEmail(r.Context(), sql.NullString{String: email, Valid: true})
	default:
		response.BadRequest(w, r, errors.New("missing query parameter"), nil)
		return
	}

	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, users)
}
