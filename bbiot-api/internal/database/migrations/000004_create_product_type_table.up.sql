CREATE TABLE IF NOT EXISTS product_types (
	id bigserial NOT NULL,
    name text NOT NULL,
    code citext NOT NULL UNIQUE,
    img_url text NOT NULL,
    -- created_at timestamp with time zone NOT NULL DEFAULT NOW(),
	-- updated_at timestamp with time zone DEFAULT Now() NOT NULL,
    PRIMARY KEY (id)
)