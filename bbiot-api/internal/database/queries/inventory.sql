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

-- name: ListInventoryByLocation :many
SELECT i.id, i.product_type_id, i.location_id, i.stock, pt.name, pt.code, pt.img
FROM inventory i
inner join product_types pt on i.product_type_id = pt.id
WHERE location_id = $1;