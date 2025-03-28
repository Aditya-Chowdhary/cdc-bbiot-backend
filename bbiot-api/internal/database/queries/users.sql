-- -- name: GetUserByEmail :one
-- SELECT *
-- FROM users
-- WHERE email = $1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: MapUserIDToLocationID :one
UPDATE users
SET location_id = $1
WHERE id = $2
RETURNING *;

-- name: MapUserToLocationName :one
UPDATE users
SET location_id = (SELECT l.id from locations l where l.location = $1)
WHERE users.id = (SELECT u.id from users u where u.username=$2)
RETURNING *;

-- name: NewUser :one
INSERT INTO users(username, password_hash, email) 
VALUES ($1, $2, $3)
RETURNING *;