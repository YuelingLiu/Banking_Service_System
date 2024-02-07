-- name: CreateUser :one
INSERT INTO
    users (username, hashed_password, full_name, email)
VALUES
    (
        sqlc.arg('username'),
        sqlc.arg('hashed_password'),
        sqlc.arg('full_name'),
        sqlc.arg('email')
    ) RETURNING *;

-- name: GetUser :one
SELECT
    *
FROM
    users
WHERE
    username = sqlc.arg('username')
LIMIT
    1;