-- Simple fix for saved_stores table
-- This handles the case where stores.id is already UUID type

-- Drop the problematic table completely
DROP TABLE IF EXISTS saved_stores CASCADE;

-- Recreate with correct UUID type (matching actual stores table)
CREATE TABLE saved_stores (
    user_id TEXT NOT NULL,
    store_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, store_id)
);

-- Add foreign key constraint separately (in case stores table doesn't exist yet)
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'stores') THEN
        ALTER TABLE saved_stores 
        ADD CONSTRAINT saved_stores_store_id_fkey 
        FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Create indexes for performance
CREATE INDEX idx_saved_stores_user ON saved_stores(user_id);
CREATE INDEX idx_saved_stores_store ON saved_stores(store_id);
