-- name: AddProductTagList :one
INSERT INTO product_tag_list (product_type, epcHex)
VALUES ($1, $2)
RETURNING *;

-- name: ListProductTagList :many
SELECT ptl.product_type, ptl.epcHex
FROM product_tag_list ptl
ORDER BY product_type;