-- name: AddEvent :one
INSERT INTO events (title,description,user_id,start_date,end_date,price) VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: FetchEvents :many
SELECT * FROM events;

-- name: FetchEvent :one
SELECT * FROM events where id=$1;
