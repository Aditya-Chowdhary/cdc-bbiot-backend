-- name: ListLocations :many
SELECT id, location
FROM locations
ORDER BY location;

-- name: AddLocation :one
INSERT INTO locations(location)
VALUEs ($1)
RETURNING *;

-- name: GetLocationByName :one
SELECT id
FROM locations
WHERE location = $1;