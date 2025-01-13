SET statement_timeout = 0;

CREATE TABLE logs (
    id VARCHAR DEFAULT gen_random_uuid() PRIMARY KEY,        
    menager_id VARCHAR NOT NULL,               
    actiontype VARCHAR NOT NULL,               
    description TEXT,                      
    record_id VARCHAR,                
    table_name VARCHAR,               
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),                     
    updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_by VARCHAR,                    
    deleted_at TIMESTAMP              
);

