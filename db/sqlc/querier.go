// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error)
	CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteProject(ctx context.Context, id int64) error
	DeleteTask(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	GetProject(ctx context.Context, id int64) (Project, error)
	GetProjectTasks(ctx context.Context, projectID int64) ([]Task, error)
	GetTask(ctx context.Context, id int64) (Task, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserTasks(ctx context.Context, assigneeID int64) ([]Task, error)
	ListProjects(ctx context.Context) ([]Project, error)
	ListTasks(ctx context.Context) ([]Task, error)
	ListUsers(ctx context.Context) ([]User, error)
	SearchProjectsByManager(ctx context.Context, managerID int64) ([]Project, error)
	SearchProjectsByTitle(ctx context.Context, dollar_1 sql.NullString) ([]Project, error)
	SearchTasksByAssignee(ctx context.Context, assigneeID int64) ([]Task, error)
	SearchTasksByPriority(ctx context.Context, priority TaskPriority) ([]Task, error)
	SearchTasksByProject(ctx context.Context, projectID int64) ([]Task, error)
	SearchTasksByStatus(ctx context.Context, status TaskStatus) ([]Task, error)
	SearchTasksByTitle(ctx context.Context, dollar_1 sql.NullString) ([]Task, error)
	SearchUsersByEmail(ctx context.Context, dollar_1 sql.NullString) ([]User, error)
	SearchUsersByName(ctx context.Context, dollar_1 sql.NullString) ([]User, error)
	UpdateProject(ctx context.Context, arg UpdateProjectParams) (Project, error)
	UpdateTask(ctx context.Context, arg UpdateTaskParams) (Task, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
