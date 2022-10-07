CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(255) PRIMARY KEY,
    user_firstname VARCHAR(50) NOT NULL,
    user_lastname VARCHAR(50) NULL,
    user_email VARCHAR(100) NOT NULL,
    user_address TEXT NULL,
    user_password VARCHAR(255) NOT NULL,
    user_role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);
