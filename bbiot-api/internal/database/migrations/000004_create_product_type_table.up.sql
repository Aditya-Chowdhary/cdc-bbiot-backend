CREATE TABLE IF NOT EXISTS product_types (
	id bigserial NOT NULL,
    name text NOT NULL,
    code text NOT NULL UNIQUE,
    -- created_at timestamp with time zone NOT NULL DEFAULT NOW(),
	-- updated_at timestamp with time zone DEFAULT Now() NOT NULL,
    PRIMARY KEY (id)
)