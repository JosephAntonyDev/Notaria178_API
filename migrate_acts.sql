-- 1. Add Category to Act Catalogs
ALTER TABLE act_catalogs ADD COLUMN IF NOT EXISTS category VARCHAR(150) DEFAULT 'General';

-- 2. Create Act Requirements Table
CREATE TABLE IF NOT EXISTS act_requirements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    act_id UUID REFERENCES act_catalogs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Update Existing Seeds (if any, although initially there are no acts seeded in schema.sql but just in case)
-- 
-- 4. Insert Seed Data for Acts and Requirements matching the image
INSERT INTO act_catalogs (id, name, description, category, status)
VALUES 
    ('a1111111-1111-1111-1111-111111111111', 'Constitucion de sociedades', 'Acta constitutiva de sociedad', 'Corporativos y Mercantiles', 'ACTIVE'),
    ('a2222222-2222-2222-2222-222222222222', 'Actas de asamblea', 'Protocolización de actas de asamblea', 'Corporativos y Mercantiles', 'ACTIVE'),
    ('a3333333-3333-3333-3333-333333333333', 'Otorgamiento de poderes', 'Otorgamiento de poder notarial', 'Corporativos y Mercantiles', 'ACTIVE'),
    ('a4444444-4444-4444-4444-444444444444', 'Revocacion de poderes', 'Revocación de poder notarial', 'Corporativos y Mercantiles', 'ACTIVE')
ON CONFLICT (name) DO UPDATE SET category = EXCLUDED.category;

-- 5. Insert Seed Data for Requirements
INSERT INTO act_requirements (act_id, name)
VALUES
    ('a1111111-1111-1111-1111-111111111111', 'INE de los socios'),
    ('a1111111-1111-1111-1111-111111111111', 'Comprobante de domicilio social'),
    ('a1111111-1111-1111-1111-111111111111', 'Proyecto de estatutos'),
    ('a1111111-1111-1111-1111-111111111111', 'Permiso de la SE');
