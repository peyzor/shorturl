-- name: CreateURL :one
insert into urls (url, short, created_at)
values ($1, $2, now())
returning *;