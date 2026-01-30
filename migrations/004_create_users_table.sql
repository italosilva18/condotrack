-- Migration: Create users table for authentication
-- Run this migration to add user authentication support

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nome VARCHAR(255) NOT NULL,
    role ENUM('admin', 'manager', 'instructor', 'student', 'user') NOT NULL DEFAULT 'user',
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    telefone VARCHAR(20) NULL,
    avatar_url VARCHAR(500) NULL,
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_email (email),
    INDEX idx_users_role (role),
    INDEX idx_users_ativo (ativo)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert default admin user (password: admin123)
-- The password hash is bcrypt with cost 10 for 'admin123'
INSERT INTO users (id, email, password_hash, nome, role, ativo) VALUES
(UUID(), 'admin@condotrack.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Administrador', 'admin', TRUE)
ON DUPLICATE KEY UPDATE id = id;
