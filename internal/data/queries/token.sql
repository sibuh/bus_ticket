-- name: GetTokenData :one
select id,
    t.id
from users u
    INNER JOIN tickets t on t.user_id = id;