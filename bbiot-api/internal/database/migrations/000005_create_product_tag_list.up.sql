CREATE TABLE IF NOT EXISTS product_tag_list (
	id bigserial NOT NULL PRIMARY KEY,
    -- product_type_id bigint REFERENCES product_types(id),
    product_type text NOT NULL,
    epcHex text NOT NULL UNIQUE
    -- created_at timestamp with time zone NOT NULL DEFAULT NOW(),
	-- updated_at timestamp with time zone DEFAULT Now() NOT NULL
)