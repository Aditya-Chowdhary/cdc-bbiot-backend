CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS locations (
	id bigserial NOT NULL,
    location citext NOT NULL UNIQUE,
    address text NOT NULL,
    site_desc text NOT NULL,
    additional_notes text,
    img_url text NOT NULL,
    -- created_at timestamp with time zone NOT NULL DEFAULT NOW(),
	-- updated_at timestamp with time zone DEFAULT Now() NOT NULL,
    PRIMARY KEY (id)
)