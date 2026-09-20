-- name: CreateURL :one
INSERT INTO urls (
    short_code,
    original_url
) VALUES (
    $1, $2
) RETURNING *;


-- name: GetURLByShortCode :one
SELECT *
FROM urls
WHERE short_code = $1;


-- name: DeleteURL :exec
DELETE FROM urls
WHERE short_code = $1;


