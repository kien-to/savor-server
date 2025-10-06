-- Migration to add pickup_timestamp field for proper timestamp handling
-- This adds a new pickup_timestamp field while keeping the old pickup_time field for backward compatibility

-- Add new timestamp field (nullable for existing records)
ALTER TABLE reservations ADD COLUMN IF NOT EXISTS pickup_timestamp TIMESTAMP WITH TIME ZONE;

-- Set default value for new records (1 hour from now)
ALTER TABLE reservations ALTER COLUMN pickup_timestamp SET DEFAULT NOW() + INTERVAL '1 hour';

-- Create index for better query performance
CREATE INDEX IF NOT EXISTS idx_reservations_pickup_timestamp ON reservations(pickup_timestamp);
