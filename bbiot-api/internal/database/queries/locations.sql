-- name: ListLocations :many
SELECT id, location, address, site_desc, additional_notes, img_url
FROM locations
ORDER BY location;

-- name: AddLocation :one
INSERT INTO locations(location,address,site_desc,additional_notes, img_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetLocationByName :one
SELECT id
FROM locations
WHERE location = $1;