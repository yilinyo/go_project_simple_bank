-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email
) VALUES(
            $1,$2,$3,$4
        )RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET hashed_password = COALESCE(SQLC.narg(hashed_password), hashed_password),
    full_name = COALESCE(SQLC.narg(full_name), full_name),
    email = COALESCE(SQLC.narg(email), email)
WHERE username = SQLC.arg(username)
RETURNING *;