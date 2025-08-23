-- Migration: Fix saved_stores table schema
-- This migration fixes the duplicate column and primary key issues in saved_stores table

-- First, backup any existing data
CREATE TEMP TABLE saved_stores_backup AS 
SELECT DISTINCT user_id, store_id, created_at 
FROM saved_stores 
WHERE user_id IS NOT NULL AND store_id IS NOT NULL;

-- Drop the problematic table
DROP TABLE IF EXISTS saved_stores CASCADE;

-- Recreate the table with correct schema
CREATE TABLE saved_stores (
    user_id TEXT NOT NULL,
    store_id VARCHAR(36) NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, store_id)
);

-- Restore data from backup (if any exists)
INSERT INTO saved_stores (user_id, store_id, created_at)
SELECT user_id, store_id, created_at 
FROM saved_stores_backup
ON CONFLICT (user_id, store_id) DO NOTHING;

-- Create indexes for performance
CREATE INDEX idx_saved_stores_user ON saved_stores(user_id);
CREATE INDEX idx_saved_stores_store ON saved_stores(store_id);

-- Clean up temp table
DROP TABLE saved_stores_backup;
