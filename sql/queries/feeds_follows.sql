-- name: CreateFeedFollow :one
INSERT INTO feeds_follows (id, created_at, updated_at, feed_id,user_id)
VALUES ($1, $2, $3, $4,$5)
RETURNING *;

-- name: DeleteFeedFollow :exec
DELETE FROM feeds_follows
WHERE id = $1;

-- name: GetAllUserFeedFollows :many
SELECT * FROM feeds_follows
WHERE user_id = $1
ORDER BY updated_at;

