SET statement_timeout = 0;

CREATE TABLE promotes (
  id VARCHAR DEFAULT gen_random_uuid() PRIMARY KEY,
  promote_name VARCHAR NOT NULL,
  promote_type VARCHAR NOT NULL,
  status BOOLEAN DEFAULT TRUE,
  link VARCHAR,
  created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
  updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
  updated_by VARCHAR,
  deleted_at TIMESTAMP
);
