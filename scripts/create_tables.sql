CREATE TABLE users (
    id BIGINT PRIMARY KEY,                          -- Sonyflake ID (generated as BIGINT)
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'customer',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    expires_at BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
CREATE INDEX idx_user_id ON refresh_tokens (user_id);

CREATE TABLE user_businesses (
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,  -- Referência ao usuário
    business_id BIGINT REFERENCES businesses(id) ON DELETE CASCADE,  -- Referência ao empreendimento
    role VARCHAR(50) DEFAULT 'admin',  -- Papel do usuário: 'admin' (administrador) ou outro (por exemplo, 'manager')
    PRIMARY KEY (user_id, business_id)  -- Chave composta, garantindo que o usuário só possa ser associado a um único empreendimento por vez
);

