-- name: UpdateTicketStatus :one
UPDATE tickets SET status ='Onhold' WHERE ticket_no=$1 AND bus_no=$2 AND trip_id=$3 RETURNING *;