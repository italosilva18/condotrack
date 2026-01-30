-- Migration: 002_additional_tables.sql
-- Description: Additional tables for users, courses, agenda, tasks, suppliers, team, inspections, and routine plans
-- Created: 2026-01-28
-- Author: Backend Architect

-- ============================================================================
-- TABLE: users
-- Description: Authentication and user management table
-- ============================================================================
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role ENUM('admin', 'gestor', 'supervisor', 'zelador', 'manutencao', 'asg', 'student', 'instructor') NOT NULL,
    phone VARCHAR(20) NULL,
    cpf VARCHAR(14) NULL,
    avatar_url VARCHAR(500) NULL,
    is_active BOOLEAN DEFAULT TRUE,
    last_login TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT uq_users_email UNIQUE (email),
    CONSTRAINT uq_users_cpf UNIQUE (cpf)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Indexes for users table
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_role_active ON users(role, is_active);

-- ============================================================================
-- TABLE: courses
-- Description: Course catalog for training platform
-- ============================================================================
CREATE TABLE IF NOT EXISTS courses (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    instructor_id VARCHAR(36) NULL,
    duration_hours INT DEFAULT 40,
    price DECIMAL(10,2) NOT NULL,
    discount_price DECIMAL(10,2) NULL,
    thumbnail_url VARCHAR(500) NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_courses_instructor
        FOREIGN KEY (instructor_id) REFERENCES users(id)
        ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Indexes for courses table
CREATE INDEX idx_courses_instructor ON courses(instructor_id);
CREATE INDEX idx_courses_is_active ON courses(is_active);
CREATE INDEX idx_courses_price ON courses(price);

-- ============================================================================
-- TABLE: agenda
-- Description: Calendar events and scheduling
-- ============================================================================
CREATE TABLE IF NOT EXISTS agenda (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,
    event_type ENUM('audit', 'inspection', 'meeting', 'task', 'other') NOT NULL DEFAULT 'other',
    start_datetime TIMESTAMP NOT NULL,
    end_datetime TIMESTAMP NULL,
    all_day BOOLEAN DEFAULT FALSE,
    location VARCHAR(255) NULL,
    contract_id VARCHAR(36) NULL,
    user_id VARCHAR(36) NULL,
    recurrence_rule VARCHAR(255) NULL,
    color VARCHAR(7) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_agenda_contract
        FOREIGN KEY (contract_id) REFERENCES contratos(id)
        ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_agenda_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Indexes for agenda table
CREATE INDEX idx_agenda_contract ON agenda(contract_id);
CREATE INDEX idx_agenda_user ON agenda(user_id);
CREATE INDEX idx_agenda_event_type ON agenda(event_type);
CREATE INDEX idx_agenda_start_datetime ON agenda(start_datetime);
CREATE INDEX idx_agenda_date_range ON agenda(start_datetime, end_datetime);

-- ============================================================================
-- TABLE: tasks
-- Description: Task management and tracking
-- ============================================================================
CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,
    status ENUM('pending', 'in_progress', 'completed', 'cancelled') DEFAULT 'pending',
    priority ENUM('low', 'medium', 'high', 'urgent') DEFAULT 'medium',
    due_date TIMESTAMP NULL,
    contract_id VARCHAR(36) NULL,
    assigned_to VARCHAR(36) NULL,
    created_by VARCHAR(36) NULL,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_tasks_contract
        FOREIGN KEY (contract_id) REFERENCES contratos(id)
        ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_tasks_assigned_to
        FOREIGN KEY (assigned_to) REFERENCES users(id)
        ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_tasks_created_by
        FOREIGN KEY (created_by) REFERENCES users(id)
        ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Indexes for tasks table
CREATE INDEX idx_tasks_contract ON tasks(contract_id);
CREATE INDEX idx_tasks_assigned_to ON tasks(assigned_to);
CREATE INDEX idx_tasks_created_by ON tasks(created_by);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
CREATE INDEX idx_tasks_status_priority ON tasks(status, priority);

-- ============================================================================
-- TABLE: suppliers
-- Description: Supplier and vendor management
-- ============================================================================
CREATE TABLE IF NOT EXISTS suppliers (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cnpj VARCHAR(18) NULL,
    email VARCHAR(255) NULL,
    phone VARCHAR(20) NULL,
    address TEXT NULL,
    category VARCHAR(100) NULL,
    is_active BOOLEAN DEFAULT TRUE,
    notes TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT uq_suppliers_cnpj UNIQUE (cnpj)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Indexes for suppliers table
CREATE INDEX idx_suppliers_name ON suppliers(name);
CREATE INDEX idx_suppliers_cnpj ON suppliers(cnpj);
CREATE INDEX idx_suppliers_category ON suppliers(category);
CREATE INDEX idx_suppliers_is_active ON suppliers(is_active);

-- ============================================================================
-- TABLE: team_members
-- Description: Team assignment to contracts
-- ============================================================================
CREATE TABLE IF NOT EXISTS team_members (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    contract_id VARCHAR(36) NOT NULL,
    role VARCHAR(50) NULL,
    start_date DATE NULL,
    end_date DATE NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_team_members_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_team_members_contract
        FOREIGN KEY (contract_id) REFERENCES contratos(id)
        ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Indexes for team_members table
CREATE INDEX idx_team_members_user ON team_members(user_id);
CREATE INDEX idx_team_members_contract ON team_members(contract_id);
CREATE INDEX idx_team_members_role ON team_members(role);
CREATE INDEX idx_team_members_is_active ON team_members(is_active);
CREATE INDEX idx_team_members_user_contract ON team_members(user_id, contract_id);
CREATE INDEX idx_team_members_dates ON team_members(start_date, end_date);

-- ============================================================================
-- TABLE: inspections
-- Description: Property inspections and audits
-- ============================================================================
CREATE TABLE IF NOT EXISTS inspections (
    id VARCHAR(36) PRIMARY KEY,
    contract_id VARCHAR(36) NOT NULL,
    inspector_id VARCHAR(36) NULL,
    inspection_date DATE NOT NULL,
    inspection_type ENUM('routine', 'preventive', 'corrective', 'emergency') NOT NULL DEFAULT 'routine',
    status ENUM('scheduled', 'in_progress', 'completed', 'cancelled') DEFAULT 'scheduled',
    findings TEXT NULL,
    recommendations TEXT NULL,
    photos JSON NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_inspections_contract
        FOREIGN KEY (contract_id) REFERENCES contratos(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_inspections_inspector
        FOREIGN KEY (inspector_id) REFERENCES users(id)
        ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Indexes for inspections table
CREATE INDEX idx_inspections_contract ON inspections(contract_id);
CREATE INDEX idx_inspections_inspector ON inspections(inspector_id);
CREATE INDEX idx_inspections_date ON inspections(inspection_date);
CREATE INDEX idx_inspections_type ON inspections(inspection_type);
CREATE INDEX idx_inspections_status ON inspections(status);
CREATE INDEX idx_inspections_contract_date ON inspections(contract_id, inspection_date);

-- ============================================================================
-- TABLE: routine_plans
-- Description: Generated cleaning/maintenance routine plans
-- ============================================================================
CREATE TABLE IF NOT EXISTS routine_plans (
    id VARCHAR(36) PRIMARY KEY,
    contract_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    plan_type ENUM('cleaning', 'maintenance', 'security', 'custom') NOT NULL DEFAULT 'cleaning',
    schedule JSON NULL COMMENT 'JSON structure containing days, times, and areas for the routine',
    is_active BOOLEAN DEFAULT TRUE,
    created_by VARCHAR(36) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_routine_plans_contract
        FOREIGN KEY (contract_id) REFERENCES contratos(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_routine_plans_created_by
        FOREIGN KEY (created_by) REFERENCES users(id)
        ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Indexes for routine_plans table
CREATE INDEX idx_routine_plans_contract ON routine_plans(contract_id);
CREATE INDEX idx_routine_plans_type ON routine_plans(plan_type);
CREATE INDEX idx_routine_plans_is_active ON routine_plans(is_active);
CREATE INDEX idx_routine_plans_created_by ON routine_plans(created_by);

-- ============================================================================
-- End of Migration 002
-- ============================================================================
