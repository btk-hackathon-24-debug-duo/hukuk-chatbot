CREATE USER btk WITH PASSWORD 'your-beautiful-password';

CREATE DATABASE btk;

GRANT ALL PRIVILEGES ON DATABASE btk TO btk;

/* You need to reconnect db before next line with user btk and db btk. */

CREATE TABLE IF NOT EXISTS users ( 
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);