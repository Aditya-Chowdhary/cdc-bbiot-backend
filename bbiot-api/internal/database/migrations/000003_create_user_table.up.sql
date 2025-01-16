CREATE TABLE IF NOT EXISTS users(
	id bigserial NOT NULL,
    username citext UNIQUE NOT NULL,
	email citext UNIQUE,
    password_hash bytea NOT NULL,
    role text NOT NULL DEFAULT 'user',
    location_id bigint REFERENCES locations(id),
    -- created_at timestamp with time zone NOT NULL DEFAULT NOW(),
	-- updated_at timestamp with time zone DEFAULT Now() NOT NULL,
    PRIMARY KEY (id)
);