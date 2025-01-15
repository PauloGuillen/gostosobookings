CREATE TABLE users (
    id BIGINT PRIMARY KEY,                          -- Sonyflake ID (generated as BIGINT)
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'customer',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE users
ADD CONSTRAINT valid_roles CHECK (role IN ('customer', 'admin', 'business_admin', 'business_manager'));


CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    expires_at BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
CREATE INDEX idx_user_id ON refresh_tokens (user_id);


-- Table for storing properties
CREATE TABLE properties (
    id BIGINT PRIMARY KEY,                          -- Sonyflake ID (generated as BIGINT)
    name VARCHAR(255) NOT NULL,
    description TEXT,
    address VARCHAR(255),
    contact_email VARCHAR(100),
    contact_phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table for mapping users to properties
CREATE TABLE user_properties (
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    property_id BIGINT REFERENCES properties(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'admin',
    PRIMARY KEY (user_id, property_id)
);

