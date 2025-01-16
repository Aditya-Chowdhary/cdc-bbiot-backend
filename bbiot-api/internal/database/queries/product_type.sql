-- name: NewProduct :one
INSERT INTO product_types (name, code)
VALUES ($1, $2)
RETURNING *;