SET statement_timeout = 0;

CREATE TABLE zone_entities (
    id VARCHAR DEFAULT gen_random_uuid(),
    zone_name VARCHAR Unique,
    lat FLOAT,
    long FLOAT,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_by VARCHAR,
    deleted_at TIMESTAMPTZ
);
