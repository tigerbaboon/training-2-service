SET statement_timeout = 0;

CREATE TABLE zones (
    id VARCHAR DEFAULT gen_random_uuid() PRIMARY KEY,
    zone_name VARCHAR,
    image_id VARCHAR,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_by VARCHAR,
    deleted_at TIMESTAMPTZ
)
