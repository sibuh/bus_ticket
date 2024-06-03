-- name: RecordPayment :one
INSERT INTO payments (user_id,event_id,payment_status,intent_id,check_in_status) VALUES($1,$2,$3,$4,$5) RETURNING *;