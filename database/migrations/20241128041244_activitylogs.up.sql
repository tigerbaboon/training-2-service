SET statement_timeout = 0;

CREATE TABLE ActivityLogs (
    ID VARCHAR DEFAULT gen_random_uuid() PRIMARY KEY,
    Section TEXT NOT NULL,
    EventType TEXT NOT NULL,
    StatusCode TEXT NOT NULL,
    Detail TEXT NOT NULL,
    Request JSONB,
    Responses JSONB,
    IpAddress TEXT NOT NULL,
    UserAgent TEXT NOT NULL,
    CreatedBy TEXT NOT NULL,
    CreatedAt BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()) NOT NULL
);
