-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING id, created_at, updated_at, name;

-- name: GetUser :one
SELECT id, created_at, updated_at, name FROM users WHERE name = $1;

-- name: Reset :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT id, created_at, updated_at, name FROM users;

-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING id, created_at, updated_at, name, url, user_id;

-- name: Feeds :many
SELECT feeds.name AS feed_name, feeds.url, users.name AS user_name
FROM feeds
JOIN users ON feeds.user_id = users.id;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (  
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id
INNER JOIN users ON users.id = inserted_feed_follow.user_id;

-- name: GetFeedByUrl :one
SELECT id, created_at, updated_at, name, url, user_id
FROM feeds
WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT 
    feed_follows.*, 
    feeds.name as feed_name,
    users.name as user_name
FROM feed_follows
JOIN feeds ON feeds.id = feed_follows.feed_id
join users ON users.id = feed_follows.user_id
WHERE feed_follows.user_id = $1;

-- name: DeleteUserFeedCombo :exec
DELETE FROM feed_follows WHERE feed_id IN (SELECT id FROM feeds WHERE feeds.url = $1)
    AND user_id = $2;

