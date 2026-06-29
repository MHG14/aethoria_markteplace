-- name: GetItem :one
SELECT * FROM items WHERE id = $1;

-- name: GetItemForUpdate :one
SELECT * FROM items WHERE id = $1 FOR UPDATE;

-- name: CreateItem :one
INSERT INTO items (name, type, status, owner_id)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: UpdateItemStatus :one
UPDATE items SET status = $2 WHERE id = $1 RETURNING *;

-- name: UpdateItemOwner :one
UPDATE items SET owner_id = $2, status = 'available' WHERE id = $1 RETURNING *;

-- name: ListItems :many
SELECT * FROM items;

-- name: ListItemsByOwner :many
SELECT * FROM items WHERE owner_id = $1;