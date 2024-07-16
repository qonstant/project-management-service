-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY full_name ASC;

-- name: CreateUser :one
INSERT INTO users (
    full_name, email, registration_date, role
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET 
    full_name = $2,
    email = $3,
    role = $4,
    registration_date = $5
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUserTasks :many
SELECT * FROM tasks
WHERE assignee_id = $1
ORDER BY creation_date ASC;

-- name: SearchUsersByName :many
SELECT * FROM users
WHERE full_name ILIKE '%' || $1 || '%'
ORDER BY full_name ASC;

-- name: SearchUsersByEmail :many
SELECT * FROM users
WHERE email ILIKE '%' || $1 || '%'
ORDER BY email ASC;
