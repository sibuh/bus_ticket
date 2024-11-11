-- name: GetTokenData :one
select u.id,
    t.id
from users u
    INNER JOIN tickets t on t.user_id = u.id;