-- name: StoreCheckoutSession :one
INSERT INTO sessions (id,ticket_id,payment_status,payment_url,cancel_url,amount,created_at)VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING *;