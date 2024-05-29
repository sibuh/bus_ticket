-- name: AddEvent :one
INSERT INTO events (title,description,user_id,start_date,end_date) VALUES ($1,$2,$3,$4,$5)
RETURNING *;

-- name: GetEvent :one
SELECT * FROM events;