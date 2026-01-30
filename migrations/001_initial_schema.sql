-- CondoTrack API - Initial Schema
-- Run this migration to create the database schema

-- ===========================================
-- GESTORES (Managers)
-- ===========================================
CREATE TABLE IF NOT EXISTS gestores (
    id VARCHAR(36) PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    telefone VARCHAR(20),
    cpf VARCHAR(14),
    ativo BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_gestores_email (email),
    INDEX idx_gestores_ativo (ativo)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- CONTRATOS (Contracts)
-- ===========================================
CREATE TABLE IF NOT EXISTS contratos (
    id VARCHAR(36) PRIMARY KEY,
    gestor_id VARCHAR(36) NOT NULL,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    endereco VARCHAR(255),
    cidade VARCHAR(100),
    estado VARCHAR(2),
    cep VARCHAR(10),
    total_unidades INT DEFAULT 0,
    meta_score DECIMAL(5,2) DEFAULT 80.00,
    data_inicio DATE,
    data_fim DATE,
    ativo BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (gestor_id) REFERENCES gestores(id) ON DELETE RESTRICT,
    INDEX idx_contratos_gestor (gestor_id),
    INDEX idx_contratos_ativo (ativo)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- AUDIT CATEGORIES
-- ===========================================
CREATE TABLE IF NOT EXISTS audit_categories (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    weight DECIMAL(5,2) DEFAULT 1.00,
    order_num INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- AUDITS
-- ===========================================
CREATE TABLE IF NOT EXISTS audits (
    id VARCHAR(36) PRIMARY KEY,
    contract_id VARCHAR(36) NOT NULL,
    auditor_name VARCHAR(255) NOT NULL,
    audit_date DATE NOT NULL,
    score DECIMAL(5,2) NOT NULL,
    target_score DECIMAL(5,2) DEFAULT 80.00,
    previous_score DECIMAL(5,2),
    status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
    observations TEXT,
    data_json JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (contract_id) REFERENCES contratos(id) ON DELETE CASCADE,
    INDEX idx_audits_contract (contract_id),
    INDEX idx_audits_date (audit_date),
    INDEX idx_audits_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- AUDIT ITEMS
-- ===========================================
CREATE TABLE IF NOT EXISTS audit_items (
    id VARCHAR(36) PRIMARY KEY,
    audit_id VARCHAR(36) NOT NULL,
    category_id VARCHAR(36),
    item_name VARCHAR(255) NOT NULL,
    score DECIMAL(5,2) NOT NULL,
    max_score DECIMAL(5,2) NOT NULL,
    observation TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (audit_id) REFERENCES audits(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES audit_categories(id) ON DELETE SET NULL,
    INDEX idx_audit_items_audit (audit_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- ENROLLMENTS (Matriculas)
-- ===========================================
CREATE TABLE IF NOT EXISTS enrollments (
    id VARCHAR(36) PRIMARY KEY,
    student_id VARCHAR(36) NOT NULL,
    student_name VARCHAR(255) NOT NULL,
    student_email VARCHAR(255) NOT NULL,
    student_cpf VARCHAR(14),
    student_phone VARCHAR(20),
    course_id VARCHAR(36) NOT NULL,
    course_name VARCHAR(255) NOT NULL,
    instructor_id VARCHAR(36),
    instructor_name VARCHAR(255),
    payment_id VARCHAR(36),
    payment_status ENUM('pending', 'confirmed', 'failed', 'refunded') DEFAULT 'pending',
    amount DECIMAL(10,2) NOT NULL,
    discount_amount DECIMAL(10,2) DEFAULT 0.00,
    final_amount DECIMAL(10,2) NOT NULL,
    payment_method VARCHAR(50),
    enrollment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completion_date TIMESTAMP NULL,
    expiration_date TIMESTAMP NULL,
    status ENUM('pending', 'active', 'completed', 'cancelled', 'expired') DEFAULT 'pending',
    progress DECIMAL(5,2) DEFAULT 0.00,
    certificate_id VARCHAR(36),
    asaas_customer_id VARCHAR(50),
    asaas_payment_id VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_enrollments_student (student_id),
    INDEX idx_enrollments_course (course_id),
    INDEX idx_enrollments_status (status),
    INDEX idx_enrollments_payment_status (payment_status),
    INDEX idx_enrollments_asaas_payment (asaas_payment_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- CERTIFICATES
-- ===========================================
CREATE TABLE IF NOT EXISTS certificates (
    id VARCHAR(36) PRIMARY KEY,
    enrollment_id VARCHAR(36) NOT NULL UNIQUE,
    student_id VARCHAR(36) NOT NULL,
    student_name VARCHAR(255) NOT NULL,
    student_cpf VARCHAR(14),
    course_id VARCHAR(36) NOT NULL,
    course_name VARCHAR(255) NOT NULL,
    course_hours INT DEFAULT 40,
    instructor_name VARCHAR(255),
    completion_date DATE NOT NULL,
    issue_date DATE NOT NULL,
    validation_code VARCHAR(20) NOT NULL UNIQUE,
    status ENUM('active', 'revoked', 'expired') DEFAULT 'active',
    download_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (enrollment_id) REFERENCES enrollments(id) ON DELETE CASCADE,
    INDEX idx_certificates_student (student_id),
    INDEX idx_certificates_validation (validation_code),
    INDEX idx_certificates_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- NOTIFICATIONS
-- ===========================================
CREATE TABLE IF NOT EXISTS notifications (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    type ENUM('payment', 'enrollment', 'audit', 'certificate', 'system') NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSON,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_notifications_user (user_id),
    INDEX idx_notifications_read (is_read),
    INDEX idx_notifications_type (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- REVENUE SPLITS
-- ===========================================
CREATE TABLE IF NOT EXISTS revenue_splits (
    id VARCHAR(36) PRIMARY KEY,
    enrollment_id VARCHAR(36) NOT NULL,
    payment_id VARCHAR(50) NOT NULL,
    gross_amount DECIMAL(10,2) NOT NULL,
    net_amount DECIMAL(10,2) NOT NULL,
    platform_fee DECIMAL(10,2) NOT NULL,
    payment_fee DECIMAL(10,2) NOT NULL,
    instructor_amount DECIMAL(10,2) NOT NULL,
    platform_amount DECIMAL(10,2) NOT NULL,
    instructor_id VARCHAR(36),
    payment_method VARCHAR(50) NOT NULL,
    status ENUM('pending', 'processed', 'failed') DEFAULT 'pending',
    processed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (enrollment_id) REFERENCES enrollments(id) ON DELETE CASCADE,
    INDEX idx_revenue_splits_enrollment (enrollment_id),
    INDEX idx_revenue_splits_instructor (instructor_id),
    INDEX idx_revenue_splits_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- MEDIA FILES
-- ===========================================
CREATE TABLE IF NOT EXISTS media_files (
    id VARCHAR(36) PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    bucket VARCHAR(50) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    url VARCHAR(500) NOT NULL,
    entity_type VARCHAR(50),
    entity_id VARCHAR(36),
    uploaded_by VARCHAR(36),
    is_public BOOLEAN DEFAULT TRUE,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_media_files_bucket (bucket),
    INDEX idx_media_files_entity (entity_type, entity_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- AUDIT EVIDENCE
-- ===========================================
CREATE TABLE IF NOT EXISTS audit_evidence (
    id VARCHAR(36) PRIMARY KEY,
    audit_id VARCHAR(36) NOT NULL,
    audit_item_id VARCHAR(36),
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size BIGINT NOT NULL,
    description TEXT,
    uploaded_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (audit_id) REFERENCES audits(id) ON DELETE CASCADE,
    FOREIGN KEY (audit_item_id) REFERENCES audit_items(id) ON DELETE SET NULL,
    INDEX idx_audit_evidence_audit (audit_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ===========================================
-- DEFAULT DATA
-- ===========================================

-- Insert default audit categories
INSERT INTO audit_categories (id, name, description, weight, order_num) VALUES
('cat-001', 'Limpeza e Conservação', 'Avaliação de limpeza das áreas comuns', 1.0, 1),
('cat-002', 'Segurança', 'Avaliação dos itens de segurança', 1.5, 2),
('cat-003', 'Manutenção Predial', 'Estado de conservação do edifício', 1.2, 3),
('cat-004', 'Jardinagem', 'Condições das áreas verdes', 0.8, 4),
('cat-005', 'Documentação', 'Conformidade documental', 1.0, 5);

-- Insert sample gestor
INSERT INTO gestores (id, nome, email, telefone, ativo) VALUES
('gest-001', 'Administrador Padrão', 'admin@condotrack.com', '11999999999', TRUE);

-- Insert sample contrato
INSERT INTO contratos (id, gestor_id, nome, descricao, cidade, estado, total_unidades, meta_score, ativo) VALUES
('cont-001', 'gest-001', 'Condomínio Exemplo', 'Condomínio de demonstração', 'São Paulo', 'SP', 50, 85.00, TRUE);
