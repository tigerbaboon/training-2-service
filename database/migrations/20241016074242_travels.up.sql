SET statement_timeout = 0;

CREATE TABLE travels (
    id VARCHAR DEFAULT gen_random_uuid() PRIMARY KEY,
    travel_title VARCHAR,
    travel_detail VARCHAR,
    status BOOLEAN DEFAULT TRUE,
    address VARCHAR NOT NULL,
    location_latitute FLOAT NOT NULL,
    location_longitute FLOAT NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_by VARCHAR,
    deleted_at TIMESTAMPTZ
)