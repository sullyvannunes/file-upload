-- name: CreateFile :one
insert into files (file, name, mimetype) values ($1, $2, $3) RETURNING *;

-- name: GetFiles :many
select id, name from files;

-- name: GetFile :one
select * from files where id=$1;
