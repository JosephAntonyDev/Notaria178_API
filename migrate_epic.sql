-- ==========================================
-- Migration: Epic - Requirements, Acts, Client CoW
-- ==========================================

-- Ad-hoc / extra requirements tied directly to a work (not via act_catalogs)
CREATE TABLE IF NOT EXISTS work_requirements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_id UUID NOT NULL REFERENCES works(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    document_id UUID REFERENCES documents(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add requirement_id to documents so we can link a document to a specific requirement
ALTER TABLE documents ADD COLUMN IF NOT EXISTS requirement_id UUID;
ALTER TABLE documents ADD COLUMN IF NOT EXISTS requirement_source VARCHAR(20) DEFAULT '';
-- requirement_source: 'act' (from act_requirements), 'work' (from work_requirements)
