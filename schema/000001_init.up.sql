CREATE TABLE users (
    guid UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    email VARCHAR(256) NOT NULL,
    refresh_hash VARCHAR(256) NOT NULL
);