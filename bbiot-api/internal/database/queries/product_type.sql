-- name: NewProduct :one
INSERT INTO product_types (name, code, img_url)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListProducts :many
SELECT *
FROM product_types
ORDER BY id;

-- name: GetProductTypeByName :one
SELECT *
FROM product_types
WHERE name=$1;