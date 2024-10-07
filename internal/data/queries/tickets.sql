-- name: UpdateTicketStatus :one
UPDATE tickets SET status =$1 WHERE id=$2 RETURNING *;
-- name: GetTicket :one
SELECT * from tickets where id=$1;