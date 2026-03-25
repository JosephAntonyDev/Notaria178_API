-- ==========================================
-- 0. DATABASE CREATION
-- ==========================================
--CREATE DATABASE notaria178_db;

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
    category VARCHAR(150) NOT NULL DEFAULT 'Sin categoría',
    status user_status DEFAULT 'ACTIVE'
);

CREATE TABLE act_requirements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    act_id UUID NOT NULL REFERENCES act_catalogs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    status user_status DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
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

-- Tabla para guardar los tokens FCM de los dispositivos de cada usuario
CREATE TABLE user_device_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    fcm_token VARCHAR(255) UNIQUE NOT NULL,
    device_type VARCHAR(50) DEFAULT 'web', -- 'web', 'android', 'ios'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para mejorar el rendimiento
CREATE INDEX IF NOT EXISTS idx_user_device_tokens_user_id ON user_device_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_user_device_tokens_fcm_token ON user_device_tokens(fcm_token);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE, -- Owner of the notification
    work_id UUID REFERENCES works(id) ON DELETE CASCADE, -- Direct link to the dossier
    type notification_type NOT NULL,
    title VARCHAR(100), -- Título para la notificación push
    body TEXT, -- Cuerpo detallado para la notificación push
    message VARCHAR(255) NOT NULL, -- Mensaje corto para compatibilidad con código existente
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para la tabla notifications
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_work_id ON notifications(work_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Who made the change
    action VARCHAR(50) NOT NULL, -- e.g., 'CREATE', 'UPDATE', 'APPROVE'
    entity VARCHAR(50) NOT NULL, -- e.g., 'WORK', 'USER', 'DOCUMENT'
    entity_id UUID NOT NULL, -- ID of the modified record
    json_details JSONB, -- Stores "Before" and "After" for deep auditing
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_action ON audit_logs(entity, action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_id ON audit_logs(entity_id);

-- ==========================================
-- 5. INITIAL SEED DATA (Semillas)
-- ==========================================

-- 1. Crear las Sucursales
INSERT INTO branches (id, name, address) 
VALUES ('11111111-1111-1111-1111-111111111111', 'Tuxtla Gutiérrez', 'Centro, Tuxtla Gutiérrez')
ON CONFLICT (id) DO NOTHING;

INSERT INTO branches (id, name, address) 
VALUES ('bbbbbbbb-1111-1111-1111-111111111111', 'San Fernando', 'Centro, San Fernando')
ON CONFLICT (id) DO NOTHING;

INSERT INTO branches (id, name, address) 
VALUES ('cccccccc-1111-1111-1111-111111111111', 'CDMX', 'Polanco, CDMX')
ON CONFLICT (id) DO NOTHING;

-- 2. Crear el Usuario SUPER_ADMIN (Notario Titular)
-- La contraseña será "admin123" (encriptada nativamente con pgcrypto)
INSERT INTO users (id, branch_id, full_name, email, password_hash, role, status)
VALUES (
    '22222222-2222-2222-2222-222222222222',
    NULL, -- El Notario Titular tiene visión global, sin sucursal fija por defecto
    'Notario Titular',
    'admin@notaria178.com',
    crypt('admin123', gen_salt('bf')),
    'SUPER_ADMIN',
    'ACTIVE'
)
ON CONFLICT (email) DO NOTHING;

-- 3. Crear 3 Proyectistas (DRAFTER)
INSERT INTO users (id, branch_id, full_name, email, password_hash, role, status)
VALUES (
    'd1111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    'Proyectista Tuxtla',
    'tuxtla@notaria178.com',
    crypt('password123', gen_salt('bf')),
    'DRAFTER',
    'ACTIVE'
)
ON CONFLICT (email) DO NOTHING;

INSERT INTO users (id, branch_id, full_name, email, password_hash, role, status)
VALUES (
    'd2222222-2222-2222-2222-222222222222',
    'bbbbbbbb-1111-1111-1111-111111111111',
    'Proyectista San Fernando',
    'sanfernando@notaria178.com',
    crypt('password123', gen_salt('bf')),
    'DRAFTER',
    'ACTIVE'
)
ON CONFLICT (email) DO NOTHING;

INSERT INTO users (id, branch_id, full_name, email, password_hash, role, status)
VALUES (
    'd3333333-3333-3333-3333-333333333333',
    'cccccccc-1111-1111-1111-111111111111',
    'Proyectista CDMX',
    'cdmx@notaria178.com',
    crypt('password123', gen_salt('bf')),
    'DRAFTER',
    'ACTIVE'
)
ON CONFLICT (email) DO NOTHING;

-- ==========================================
-- 4. Crear Actos Reales (Catálogo de Actos)
-- ==========================================
INSERT INTO act_catalogs (id, name, description, category, status)
VALUES 
    -- INMOBILIARIOS Y CIVILES
    ('a0000000-0000-0000-0000-000000000001', 'Compraventa de inmueble', 'Transmisión de propiedad de un bien inmueble (casa, terreno, local comercial).', 'Inmobiliarios y Civiles', 'ACTIVE'),
    ('a0000000-0000-0000-0000-000000000002', 'Testamento público abierto', 'Disposición de bienes, derechos y obligaciones para después de la muerte ante notario.', 'Inmobiliarios y Civiles', 'ACTIVE'),
    ('a0000000-0000-0000-0000-000000000003', 'Donación con reserva de usufructo', 'Transmisión gratuita de un bien, reservándose el donante el derecho de uso vitalicio.', 'Inmobiliarios y Civiles', 'ACTIVE'),
    ('a0000000-0000-0000-0000-000000000004', 'Cancelación de hipoteca', 'Liberación jurídica de un gravamen una vez liquidado el crédito (Infonavit, Bancario, etc.).', 'Inmobiliarios y Civiles', 'ACTIVE'),
    ('a0000000-0000-0000-0000-000000000005', 'Adjudicación por herencia', 'Formalización de la transmisión de bienes a los herederos legales o testamentarios.', 'Inmobiliarios y Civiles', 'ACTIVE'),
    ('a0000000-0000-0000-0000-000000000006', 'Régimen de propiedad en condominio', 'Acta constitutiva legal para la creación de edificios de departamentos o fraccionamientos.', 'Inmobiliarios y Civiles', 'ACTIVE'),

    -- PODERES Y MANDATOS
    ('a0000000-0000-0000-0000-000000000007', 'Poder general para actos de dominio', 'Faculta al apoderado para actuar como dueño absoluto (vender, hipotecar, donar).', 'Poderes y Mandatos', 'ACTIVE'),
    ('a0000000-0000-0000-0000-000000000008', 'Poder general para pleitos y cobranzas', 'Otorga facultades únicamente para representar al mandante en juicios y cobrar deudas.', 'Poderes y Mandatos', 'ACTIVE'),
    ('a0000000-0000-0000-0000-000000000009', 'Revocación de poderes', 'Cancelación legal y definitiva de facultades otorgadas previamente a un representante.', 'Poderes y Mandatos', 'ACTIVE'),

    -- CORPORATIVOS Y MERCANTILES
    ('a0000000-0000-0000-0000-000000000010', 'Constitución de sociedades', 'Creación y registro de nuevas empresas (S.A., S. de R.L., S.A.P.I., S.C., etc.).', 'Corporativos y Mercantiles', 'ACTIVE'),
    ('a0000000-0000-0000-0000-000000000011', 'Protocolización de actas de asamblea', 'Formalización notarial de acuerdos tomados por los socios o accionistas de una empresa.', 'Corporativos y Mercantiles', 'ACTIVE'),
    ('a0000000-0000-0000-0000-000000000012', 'Fusión o escisión de sociedades', 'Unión de dos o más empresas, o la división de una sociedad mercantil en varias.', 'Corporativos y Mercantiles', 'ACTIVE'),

    -- CERTIFICACIONES Y FE PÚBLICA
    ('a0000000-0000-0000-0000-000000000013', 'Cotejo y certificación de documentos', 'Certificación notarial de que una copia es fiel y exacta de su documento original.', 'Certificaciones y Fe Pública', 'ACTIVE'),
    ('a0000000-0000-0000-0000-000000000014', 'Fe de hechos', 'Constancia notarial objetiva sobre hechos materiales, lugares o situaciones específicas.', 'Certificaciones y Fe Pública', 'ACTIVE')
ON CONFLICT (name) DO NOTHING;


-- ==========================================
-- 5. Crear Requisitos Base (Checklists)
-- ==========================================
INSERT INTO act_requirements (act_id, name)
VALUES 
    -- 1. Requisitos Compraventa (ID: a0000000-0000-0000-0000-000000000001)
    ('a0000000-0000-0000-0000-000000000001', 'Identificación oficial (INE/Pasaporte) vigente de vendedor y comprador'),
    ('a0000000-0000-0000-0000-000000000001', 'Actas de nacimiento y matrimonio (si aplica) de ambas partes'),
    ('a0000000-0000-0000-0000-000000000001', 'Constancia de Situación Fiscal (RFC) de ambas partes'),
    ('a0000000-0000-0000-0000-000000000001', 'Escritura original del inmueble con sello del Registro Público'),
    ('a0000000-0000-0000-0000-000000000001', 'Boleta Predial pagada del año en curso'),
    ('a0000000-0000-0000-0000-000000000001', 'Recibo de agua sin adeudos (Últimos 3 meses)'),
    
    -- 2. Requisitos Testamento (ID: a0000000-0000-0000-0000-000000000002)
    ('a0000000-0000-0000-0000-000000000002', 'Identificación oficial (INE/Pasaporte) del testador'),
    ('a0000000-0000-0000-0000-000000000002', 'Acta de nacimiento original'),
    ('a0000000-0000-0000-0000-000000000002', 'CURP actualizada'),
    ('a0000000-0000-0000-0000-000000000002', 'Nombres completos de herederos, legatarios y albacea'),

    -- 3. Requisitos Constitución de sociedades (ID: a0000000-0000-0000-0000-000000000010)
    ('a0000000-0000-0000-0000-000000000010', 'Permiso de la Secretaría de Economía (Denominación social)'),
    ('a0000000-0000-0000-0000-000000000010', 'Identificaciones oficiales de todos los socios'),
    ('a0000000-0000-0000-0000-000000000010', 'Constancias de Situación Fiscal (RFC) de los socios'),
    ('a0000000-0000-0000-0000-000000000010', 'Comprobante de domicilio fiscal de la nueva sociedad'),

    -- 4. Requisitos Poder general (ID: a0000000-0000-0000-0000-000000000008)
    ('a0000000-0000-0000-0000-000000000008', 'Identificación oficial del poderdante (quien otorga el poder)'),
    ('a0000000-0000-0000-0000-000000000008', 'Constancia de Situación Fiscal del poderdante'),
    ('a0000000-0000-0000-0000-000000000008', 'Nombre completo y generales del apoderado');