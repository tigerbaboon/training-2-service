SET statement_timeout = 0;

CREATE TABLE houses (
    id VARCHAR DEFAULT gen_random_uuid() PRIMARY KEY,
    house_name VARCHAR NOT NULL,
    house_type VARCHAR NOT NULL,
    zone_id VARCHAR NOT NULL,
    sell_type VARCHAR NOT NULL,
    amenity_id VARCHAR NOT NULL,
    size FLOAT NOT NULL,
    floor FLOAT NOT NULL,
    price FLOAT NOT NULL,
    number_of_rooms INT NOT NULL,
    number_of_bathrooms INT NOT NULL,
    water_rate FLOAT NOT NULL,
    electricity_rate FLOAT NOT NULL,
    description VARCHAR NOT NULL,
    address VARCHAR NOT NULL,
    location_latitute FLOAT NOT NULL,
    location_longitute FLOAT NOT NULL,
    is_recommend BOOLEAN DEFAULT FALSE,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_by VARCHAR,
    deleted_at TIMESTAMPTZ
);
