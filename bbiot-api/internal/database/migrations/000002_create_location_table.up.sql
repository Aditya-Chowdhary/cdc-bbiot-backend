CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS locations (
	id bigserial NOT NULL,
    location text NOT NULL UNIQUE,
    -- created_at timestamp with time zone NOT NULL DEFAULT NOW(),
	-- updated_at timestamp with time zone DEFAULT Now() NOT NULL,
    PRIMARY KEY (id)
)