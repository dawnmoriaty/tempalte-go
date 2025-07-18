-- name: CreateUser :one
INSERT INTO users (
  email,
  name,
  password,
  role_id
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: CheckUserExists :one
SELECT COUNT(*) > 0 FROM users
WHERE email = $1;

-- name: GetUserWithRoleByEmail :one
SELECT
  u.id,
  u.email,
  u.name,
  u.password,
  r.name as role_name
FROM
  users u
JOIN
  roles r ON u.role_id = r.id
WHERE
  u.email = $1
LIMIT 1;