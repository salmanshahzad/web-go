-- name: CountUsersWithUsername :one
SELECT COUNT(*) FROM "user" WHERE username = $1;

-- name: CreateUser :one
INSERT INTO "user" (username, password) VALUES ($1, $2) RETURNING id;

-- name: DeleteUser :exec
DELETE FROM "user" WHERE id = $1;

-- name: GetUser :one
SELECT * FROM "user" WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM "user" WHERE username = $1;

-- name: UpdatePassword :exec
UPDATE "user" SET password = $2 WHERE id = $1;

-- name: UpdateUsername :exec
UPDATE "user" SET username = $2 WHERE id = $1;
