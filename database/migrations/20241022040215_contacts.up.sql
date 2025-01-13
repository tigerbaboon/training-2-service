SET statement_timeout = 0;

CREATE TABLE contacts (
    id VARCHAR DEFAULT gen_random_uuid() PRIMARY KEY,
    house_id VARCHAR NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    phone_number VARCHAR,
    line_id VARCHAR,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    deleted_at TIMESTAMPTZ
);
