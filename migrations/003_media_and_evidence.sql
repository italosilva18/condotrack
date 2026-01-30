-- Migration: Create media_files and audit_evidence tables
-- These tables support file uploads and audit evidence tracking

-- Media files table for general file storage
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

-- Audit evidence table for audit photos/documents
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
