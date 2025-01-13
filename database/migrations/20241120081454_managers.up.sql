SET statement_timeout = 0;

CREATE TABLE managers(
    id VARCHAR DEFAULT gen_random_uuid() PRIMARY KEY,
    username VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    manager_name VARCHAR NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_by VARCHAR,
    deleted_at TIMESTAMPTZ
);
