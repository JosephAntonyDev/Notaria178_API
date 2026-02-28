-- ==========================================
-- 0. DATABASE CREATION
-- ==========================================
CREATE DATABASE notaria178_db;

-- (Execute the following command only if you are in psql to change database)
-- \c notaria178_db;

-- ==========================================
-- 1. EXTENSIONS & ENUMS
-- ==========================================
-- Enable native UUIDs
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Enums for standardization
CREATE TYPE user_role AS ENUM ('SUPER_ADMIN', 'LOCAL_ADMIN', 'DRAFTER', 'DATA_ENTRY');
CREATE TYPE user_status AS ENUM ('ACTIVE', 'INACTIVE');
CREATE TYPE work_status AS ENUM ('PENDING', 'IN_PROGRESS', 'READY_FOR_REVIEW', 'REJECTED', 'APPROVED');
CREATE TYPE document_category AS ENUM ('DRAFT_DEED', 'FINAL_DEED', 'CLIENT_REQUIREMENT', 'OTHER');
CREATE TYPE notification_type AS ENUM ('NEW_COMMENT', 'ASSIGNMENT', 'STATUS_CHANGE', 'SYSTEM');

-- ==========================================
-- 2. MAIN TABLES (Branches, Users, Clients)
-- ==========================================

CREATE TABLE branches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    branch_id UUID REFERENCES branches(id) ON DELETE RESTRICT,
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    role user_role NOT NULL DEFAULT 'DRAFTER',
    status user_status NOT NULL DEFAULT 'ACTIVE',
    hire_date DATE DEFAULT CURRENT_DATE,
    start_time TIME,
    end_time TIME,
    photo_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE attendances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    check_in_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    check_out_time TIMESTAMP
    -- UNIQUE constraint removed to allow multiple shifts.
    -- Time logic will be handled by the Go API.
);

CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(200) NOT NULL,
    rfc VARCHAR(13),
    phone VARCHAR(20),
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE act_catalogs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) UNIQUE NOT NULL,
    description TEXT,
    status user_status DEFAULT 'ACTIVE'
);

-- ==========================================
-- 3. CORE: WORKS & ACTS (Expedientes)
-- ==========================================

CREATE TABLE works (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    branch_id UUID REFERENCES branches(id) ON DELETE RESTRICT,
    client_id UUID REFERENCES clients(id) ON DELETE RESTRICT,
    main_drafter_id UUID REFERENCES users(id) ON DELETE SET NULL,
    folio VARCHAR(50) UNIQUE, -- Can be null initially and assigned later
    status work_status NOT NULL DEFAULT 'PENDING',
    deadline DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Pivot table: A work can have multiple acts (e.g., Sale + Usufruct)
CREATE TABLE work_acts (
    work_id UUID REFERENCES works(id) ON DELETE CASCADE,
    act_id UUID REFERENCES act_catalogs(id) ON DELETE RESTRICT,
    PRIMARY KEY (work_id, act_id)
);

-- Pivot table: Multiple drafters collaborating on a single work
CREATE TABLE work_collaborators (
    work_id UUID REFERENCES works(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (work_id, user_id)
);

-- ==========================================
-- 4. DOCUMENTS, HISTORY & COMMUNICATION
-- ==========================================

CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID REFERENCES clients(id) ON DELETE CASCADE, -- Link to client's "Vault"
    work_id UUID REFERENCES works(id) ON DELETE CASCADE, -- Work it currently belongs to
    user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Uploaded by
    document_name VARCHAR(200) NOT NULL,
    category document_category NOT NULL,
    version INT DEFAULT 1,
    file_path TEXT NOT NULL, -- Physical path in Ubuntu Server (/var/uploads/...)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE work_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_id UUID REFERENCES works(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE, -- Owner of the notification
    work_id UUID REFERENCES works(id) ON DELETE CASCADE, -- Direct link to the dossier
    type notification_type NOT NULL,
    message VARCHAR(255) NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Who made the change
    action VARCHAR(50) NOT NULL, -- e.g., 'CREATE', 'UPDATE', 'APPROVE'
    entity VARCHAR(50) NOT NULL, -- e.g., 'WORK', 'USER', 'DOCUMENT'
    entity_id UUID NOT NULL, -- ID of the modified record
    json_details JSONB, -- Stores "Before" and "After" for deep auditing
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);