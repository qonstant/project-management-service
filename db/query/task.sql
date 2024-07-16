-- name: GetTask :one
SELECT * FROM tasks
WHERE id = $1 LIMIT 1;

-- name: ListTasks :many
SELECT * FROM tasks
ORDER BY creation_date ASC;

-- name: CreateTask :one
INSERT INTO tasks (
    title, description, priority, status, assignee_id, project_id, creation_date, completion_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET 
    title = $2,
    description = $3,
    priority = $4,
    status = $5,
    assignee_id = $6,
    project_id = $7,
    creation_date = $8,
    completion_date = $9
WHERE id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;

-- name: SearchTasksByTitle :many
SELECT * FROM tasks
WHERE title ILIKE '%' || $1 || '%'
ORDER BY creation_date ASC;

-- name: SearchTasksByStatus :many
SELECT * FROM tasks
WHERE status = $1
ORDER BY creation_date ASC;

-- name: SearchTasksByPriority :many
SELECT * FROM tasks
WHERE priority = $1
ORDER BY creation_date ASC;

-- name: SearchTasksByAssignee :many
SELECT * FROM tasks
WHERE assignee_id = $1
ORDER BY creation_date ASC;

-- name: SearchTasksByProject :many
SELECT * FROM tasks
WHERE project_id = $1
ORDER BY creation_date ASC;