-- Simple migration to add pickup_timestamp to reservations table
ALTER TABLE reservations ADD COLUMN IF NOT EXISTS pickup_timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW() + INTERVAL '2 hours';
