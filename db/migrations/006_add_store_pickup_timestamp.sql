-- Migration to add pickup_timestamp field to stores table
-- This adds a proper timestamp field for store pickup times

-- Add pickup_timestamp field to stores table
ALTER TABLE stores ADD COLUMN IF NOT EXISTS pickup_timestamp TIMESTAMP WITH TIME ZONE;

-- Set default pickup times for existing stores (2 PM today as default)
UPDATE stores 
SET pickup_timestamp = DATE_TRUNC('day', NOW()) + INTERVAL '14 hours'
WHERE pickup_timestamp IS NULL;

-- Set default value for new stores (2 PM today)
ALTER TABLE stores ALTER COLUMN pickup_timestamp SET DEFAULT DATE_TRUNC('day', NOW()) + INTERVAL '14 hours';

-- Create index for better query performance
CREATE INDEX IF NOT EXISTS idx_stores_pickup_timestamp ON stores(pickup_timestamp);
