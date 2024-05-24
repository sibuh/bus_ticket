
-- name: CreateUser :one
INSERT INTO users (first_name,last_name,phone,email,username,password) VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE username=$1;