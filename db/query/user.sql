-- name: CreateUser :one
INSERT INTO users (username, email, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: CheckUserNameExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: AddRoleToUser :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2);

-- name: GetRolesForUser :many
SELECT r.name
FROM roles r
JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = $1;

