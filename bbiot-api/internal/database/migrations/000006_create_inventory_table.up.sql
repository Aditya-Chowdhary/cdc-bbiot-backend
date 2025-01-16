CREATE TABLE IF NOT EXISTS inventory (
    id bigserial PRIMARY KEY,
    product_type_id bigint NOT NULL REFERENCES product_types(id) ON DELETE CASCADE ON UPDATE CASCADE,
    location_id bigint NOT NULL REFERENCES locations(id) ON DELETE CASCADE ON UPDATE CASCADE,
    stock integer NOT NULL DEFAULT(0)
    -- created_at timestamp with time zone NOT NULL DEFAULT NOW(),
	-- updated_at timestamp with time zone DEFAULT Now() NOT NULL,
)