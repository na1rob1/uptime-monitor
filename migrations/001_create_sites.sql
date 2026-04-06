CREATE TABLE Sites (
    id SERIAL PRIMARY KEY,
    url varchar(255),
    name varchar(255),
    status BOOL DEFAULT false,
    uptime REAL,
    checked_at timestamptz,
    created_at timestamptz DEFAULT NOW(),
    total_checks INT DEFAULT 0,
    up_checks INT DEFAULT 0
)