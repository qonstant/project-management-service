-- name: GetProject :one
SELECT * FROM projects
WHERE id = $1 LIMIT 1;

-- name: ListProjects :many
SELECT * FROM projects
ORDER BY start_date ASC;

-- name: CreateProject :one
INSERT INTO projects (
    name, description, start_date, end_date, manager_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateProject :one
UPDATE projects
SET 
    name = $2,
    description = $3,
    start_date = $4,
    end_date = $5,
    manager_id = $6
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;

-- name: GetProjectTasks :many
SELECT * FROM tasks
WHERE project_id = $1
ORDER BY creation_date ASC;

-- name: SearchProjectsByTitle :many
SELECT * FROM projects
WHERE name ILIKE '%' || $1 || '%'
ORDER BY start_date ASC;

-- name: SearchProjectsByManager :many
SELECT * FROM projects
WHERE manager_id = $1
ORDER BY start_date ASC;
