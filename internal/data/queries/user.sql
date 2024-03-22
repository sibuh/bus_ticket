
-- name: RegisterPayedUser :one
INSERT INTO users (first_name,last_name,phone,email,nonce,session_id) VALUES ($1,$2,$3,$4,$5,$6)
RETURNING *;

-- name: UpdatePaymentStatus :one 
UPDATE users SET payment_status=$1 WHERE session_id=$2 RETURNING *; 

-- name: GetUser :one
SELECT * FROM users WHERE nonce=$1;