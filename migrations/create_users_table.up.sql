CREATE TYPE roles AS ENUM('admin', 'user');

CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(255) PRIMARY KEY,
    user_firstname VARCHAR(50) NOT NULL,
    user_lastname VARCHAR(50) NULL,
    user_email VARCHAR(100) NOT NULL,
    user_address TEXT NULL,
    user_password VARCHAR(255) NOT NULL,
    user_role roles DEFAULT 'user',
    created_at DATETIME NULL,
    updated_at DATETIME NULL,
    deleted_at DATETIME NULL
);