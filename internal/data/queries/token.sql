-- name: GetTokenData :one
select u.id,
    t.id
from users u
    INNER JOIN tickets t on t.user_id = u.id;

-- name: GetTicketInfo :one
SELECT u.first_name,u.last_name,t.trip_id FROM users u
INNER JOIN  tickets t 
ON u.id=t.user_id
WHERE u.id=$1
AND t.id=$2;