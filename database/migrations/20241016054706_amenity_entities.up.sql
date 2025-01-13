SET statement_timeout = 0;

CREATE TABLE amenity_entities (
  id VARCHAR DEFAULT gen_random_uuid() PRIMARY KEY,
  amenity_name VARCHAR NOT NULL,
  created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
  updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
  updated_by VARCHAR,
  deleted_at TIMESTAMP
);
