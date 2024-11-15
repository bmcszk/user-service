-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
  name, other, created_at
) VALUES (
  $1, $2, now()
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
  set name = $2,
  other = $3,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
