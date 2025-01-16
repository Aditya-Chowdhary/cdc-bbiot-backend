-- name: NewInventoryProduct :one
INSERT INTO inventory (product_type_id, location_id, stock)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateInventoryStock :one
UPDATE inventory
SET stock = stock + $3
WHERE product_type_id = $1
AND location_id = $2
RETURNING *;