-- name: CreateLine :one
INSERT INTO lines (destination, departure, price, schedule)
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: GetLine :one
SELECT *
from lines
WHERE id = $1;
-- name: UpdateLine :one
UPDATE lines
SET price = $2,
    schedule = $3
WHERE id = $1
RETURNING *;
-- name: CreateTrip :one
INSERT INTO line_trips (line, date)
VALUES ($1, $2)
RETURNING *;