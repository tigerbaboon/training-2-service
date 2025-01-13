SET statement_timeout = 0;

CREATE TABLE images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    s_id VARCHAR NOT NULL,
    image_url VARCHAR NOT NULL,
    type VARCHAR NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_by VARCHAR,
    deleted_at TIMESTAMPTZ
);
